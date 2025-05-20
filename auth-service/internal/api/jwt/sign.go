package jwt

import (
	"time"

	"github.com/google/uuid"
	"github.com/tech-inspire/service/auth-service/internal/models"
)

type SignOutput struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

func (j Signer) SignTokens(user models.User, sessionID uuid.UUID, sessionToken string) (*SignOutput, error) {
	accessToken, accessExpiresAt, err := j.BuildUserAccessToken(user, sessionID)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := j.BuildUserRefreshToken(user.ID, sessionID, sessionToken)
	if err != nil {
		return nil, err
	}

	return &SignOutput{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}
