package jwt

import (
	"time"

	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tech-inspire/service/auth-service/internal/models"
)

type TokenUse string

const (
	AccessToken  TokenUse = "access"
	RefreshToken TokenUse = "refresh"
)

type UserAccessTokenClaims struct {
	jwt.RegisteredClaims
	TokenUse TokenUse `json:"token_use"`

	IsAdmin bool `json:"is_admin,omitempty"`

	SessionID uuid.UUID `json:"session_id"`
}

const (
	Issuer   = "inspire-auth"
	Audience = "inspire-web"
)

func (j Manager) BuildUserAccessToken(user models.User, sessionID uuid.UUID) (token string, expiresAt time.Time, err error) {
	var (
		now             = time.Now()
		accessExpiresAt = now.Add(j.accessTokenDuration)
	)

	accessToken := jwt.NewWithClaims(new(jwt.SigningMethodEd25519), UserAccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   user.ID.String(),
			Audience:  []string{Audience},
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		TokenUse:  AccessToken,
		SessionID: sessionID,
		IsAdmin:   user.IsAdmin,
	})

	accessToken.Header["kid"] = j.jwtKID

	signedAccessToken, err := accessToken.SignedString(j.privateKey)
	if err != nil {
		return "", accessExpiresAt, errors.Errorf("failed to sign access token: %w", err)
	}

	return signedAccessToken, accessExpiresAt, nil
}

type UserRefreshTokenClaims struct {
	jwt.RegisteredClaims
	TokenUse TokenUse `json:"token_use"`

	SessionID    uuid.UUID `json:"session_id"`
	SessionToken string    `json:"session_token"`
}

func (j Manager) BuildUserRefreshToken(userID, sessionID uuid.UUID, sessionToken string) (token string, expiresAt time.Time, err error) {
	var (
		now              = time.Now()
		refreshExpiresAt = now.Add(j.refreshTokenDuration)
	)

	refreshToken := jwt.NewWithClaims(new(jwt.SigningMethodEd25519), UserRefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   userID.String(),
			Audience:  []string{Audience},
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		TokenUse:     RefreshToken,
		SessionID:    sessionID,
		SessionToken: sessionToken,
	})
	refreshToken.Header["kid"] = j.jwtKID

	signedRefreshToken, err := refreshToken.SignedString(j.privateKey)
	if err != nil {
		return "", refreshExpiresAt, errors.Errorf("failed to sign refresh token: %w", err)
	}

	return signedRefreshToken, refreshExpiresAt, nil
}
