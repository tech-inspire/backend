package service

import (
	"context"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/tech-inspire/service/auth-service/internal/models"
)

func (a AuthService) createSession(ctx context.Context, userID, sessionID uuid.UUID) (*models.Session, error) {
	const refreshTokenLength = 64
	var (
		createdAt    = time.Now()
		expiresAt    = createdAt.Add(a.refreshTokenDuration)
		refreshToken = a.generator.GenerateString(refreshTokenLength)
	)

	session := models.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     refreshToken,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}

	err := a.sessionRepository.AddUserSession(ctx, session, a.sessionsLimitPerUser)
	if err != nil {
		return nil, errors.Errorf("create user session: %w", err)
	}

	return &session, nil
}
