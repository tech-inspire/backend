package jwt

import (
	"time"

	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ValidateUserAccessTokenOutput struct {
	SessionID uuid.UUID
	UserID    uuid.UUID

	IsAdmin bool
}

func (j Manager) ValidateUserAccessToken(accessToken string) (*ValidateUserAccessTokenOutput, error) {
	token, err := jwt.ParseWithClaims(accessToken, new(UserAccessTokenClaims), j.keyFunc)
	if err != nil {
		return nil, errors.Errorf("jwt: parse token: %w", err)
	}

	claims := token.Claims.(*UserAccessTokenClaims)

	if claims.TokenUse != AccessToken {
		return nil, errors.Errorf("jwt: invalid token used: expected '%s', got '%s'",
			AccessToken, claims.TokenUse,
		)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.Errorf("jwt: parse 'sub' into uuid: %w", err)
	}

	return &ValidateUserAccessTokenOutput{
		SessionID: claims.SessionID,
		UserID:    userID,
		IsAdmin:   claims.IsAdmin,
	}, nil
}

type ValidateAdminRefreshTokenOutput struct {
	UserID       uuid.UUID
	SessionID    uuid.UUID
	SessionToken string
	ExpiresAt    time.Time
}

func (j Manager) ValidateUserRefreshToken(refreshToken string) (*ValidateAdminRefreshTokenOutput, error) {
	token, err := jwt.ParseWithClaims(refreshToken, new(UserRefreshTokenClaims), j.keyFunc)
	if err != nil {
		return nil, errors.Errorf("jwt: parse token: %w", err)
	}

	claims := token.Claims.(*UserRefreshTokenClaims)

	if claims.TokenUse != RefreshToken {
		return nil, errors.Errorf("jwt: invalid token used: expected '%s', got '%s'",
			RefreshToken, claims.TokenUse,
		)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.Errorf("jwt: parse 'sub' into uuid: %w", err)
	}

	return &ValidateAdminRefreshTokenOutput{
		UserID:       userID,
		SessionID:    claims.SessionID,
		SessionToken: claims.SessionToken,
		ExpiresAt:    claims.ExpiresAt.Time,
	}, nil
}
