package jwt

import (
	"time"

	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tech-inspire/service/auth-service/internal/models"
	authjwt "github.com/tech-inspire/service/auth-service/pkg/jwt"
)

func (j Signer) BuildUserAccessToken(user models.User, sessionID uuid.UUID) (token string, expiresAt time.Time, err error) {
	var (
		now             = time.Now()
		accessExpiresAt = now.Add(j.accessTokenDuration)
	)

	accessToken := jwt.NewWithClaims(new(jwt.SigningMethodEd25519), authjwt.UserAccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    authjwt.Issuer,
			Subject:   user.ID.String(),
			Audience:  []string{authjwt.Audience},
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		TokenUse:  authjwt.AccessToken,
		SessionID: sessionID,
		IsAdmin:   user.IsAdmin,
	})

	accessToken.Header["kid"] = j.jwtKID

	signedAccessToken, err := accessToken.SignedString(j.privateKey)
	if err != nil {
		return "", accessExpiresAt, errors.Errorf("sign access token: %w", err)
	}

	return signedAccessToken, accessExpiresAt, nil
}

func (j Signer) BuildUserRefreshToken(userID, sessionID uuid.UUID, sessionToken string) (token string, expiresAt time.Time, err error) {
	var (
		now              = time.Now()
		refreshExpiresAt = now.Add(j.refreshTokenDuration)
	)

	refreshToken := jwt.NewWithClaims(new(jwt.SigningMethodEd25519), authjwt.UserRefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    authjwt.Issuer,
			Subject:   userID.String(),
			Audience:  []string{authjwt.Audience},
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		TokenUse:     authjwt.RefreshToken,
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
