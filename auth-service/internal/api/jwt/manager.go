package jwt

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tech-inspire/service/auth-service/internal/config"
	"github.com/tech-inspire/service/auth-service/internal/models"
	"golang.org/x/crypto/ed25519"
)

type Manager struct {
	privateKey ed25519.PrivateKey
	jwtKID     string
	keyFunc    jwt.Keyfunc
	jwks       jwkset.Storage

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func mustGenerateKIDFromPublicKey(publicKey ed25519.PublicKey) string {
	hash := sha256.New()
	_, err := hash.Write(publicKey)
	if err != nil {
		panic(err)
	}

	kid := uuid.UUID(hash.Sum(nil)[:16]).String()
	return kid
}

func NewManager(
	cfg *config.Config,
) (*Manager, error) {
	data, err := os.ReadFile(cfg.JWT.UserJWKPath)
	if err != nil {
		return nil, fmt.Errorf("read JWT key: %w", err)
	}

	key, err := jwt.ParseEdPrivateKeyFromPEM(data)
	if err != nil {
		return nil, fmt.Errorf("parse ed private key: %w", err)
	}

	privateKey := key.(ed25519.PrivateKey)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	kid := mustGenerateKIDFromPublicKey(publicKey)

	options := jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{
			KID: kid,
		},
	}

	// Create the JWK from the key and options.
	jwk, err := jwkset.NewJWKFromKey(publicKey, options)
	if err != nil {
		return nil, fmt.Errorf("create jwk from key: %w", err)
	}

	storage := jwkset.NewMemoryStorage()
	if err = storage.KeyWrite(context.TODO(), jwk); err != nil {
		return nil, fmt.Errorf("add key to storage: %w", err)
	}

	kf, err := keyfunc.New(keyfunc.Options{
		Storage: storage,
	})
	if err != nil {
		return nil, fmt.Errorf("create keyfunc: %w", err)
	}

	return &Manager{
		privateKey: privateKey,
		jwtKID:     kid,
		keyFunc:    kf.Keyfunc,
		jwks:       storage,

		accessTokenDuration:  cfg.JWT.AccessTokenDuration,
		refreshTokenDuration: cfg.JWT.RefreshTokenDuration,
	}, nil
}

type BuildOutput struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

func (j Manager) BuildTokens(user models.User, sessionID uuid.UUID, sessionToken string) (*BuildOutput, error) {
	accessToken, accessExpiresAt, err := j.BuildUserAccessToken(user, sessionID)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := j.BuildUserRefreshToken(user.ID, sessionID, sessionToken)
	if err != nil {
		return nil, err
	}

	return &BuildOutput{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}

func (j Manager) PublicUsersJWKS() (json.RawMessage, error) {
	return j.jwks.JSONPublic(context.TODO())
}
