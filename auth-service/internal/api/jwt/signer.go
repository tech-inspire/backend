package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tech-inspire/backend/auth-service/internal/config"
	authjwt "github.com/tech-inspire/backend/auth-service/pkg/jwt"
	"golang.org/x/crypto/ed25519"
)

type Signer struct {
	privateKey ed25519.PrivateKey
	jwtKID     string
	jwks       jwkset.Storage

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewSigner(
	cfg *config.Config,
) (*Signer, error) {
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
	kid := mustGenerateKID(publicKey)

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

	return &Signer{
		privateKey: privateKey,
		jwtKID:     kid,
		jwks:       storage,

		accessTokenDuration:  cfg.JWT.AccessTokenDuration,
		refreshTokenDuration: cfg.JWT.RefreshTokenDuration,
	}, nil
}

func (j Signer) Validator() (*authjwt.Validator, error) {
	kf, err := keyfunc.New(keyfunc.Options{
		Storage: j.jwks,
	})
	if err != nil {
		return nil, fmt.Errorf("create keyfunc: %w", err)
	}

	return authjwt.NewValidatorFromKeyFunc(kf.Keyfunc), nil
}

func (j Signer) PublicUsersJWKS() (json.RawMessage, error) {
	return j.jwks.JSONPublic(context.TODO())
}
