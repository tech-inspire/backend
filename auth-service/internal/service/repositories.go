package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/tech-inspire/backend/auth-service/internal/clients/mail"
	"github.com/tech-inspire/backend/auth-service/internal/models"
	"github.com/tech-inspire/backend/auth-service/internal/service/dto"
)

type Generator interface {
	GenerateString(length int) string
	NewUUID() uuid.UUID
}

type UserRepository interface {
	GetUsersCount(ctx context.Context, params dto.GetUsersParams) (count int, err error)
	UpdateUserByID(ctx context.Context, userID uuid.UUID, params dto.UpdateUsersParams) error
	GetUsers(ctx context.Context, params dto.GetUsersParams) ([]models.User, error)
	GetUserByEmailWithHash(ctx context.Context, email string) (admin *models.User, hash []byte, err error)
	CreateUser(ctx context.Context, params dto.CreateUserParams) error

	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.User, error)

	GetUserByUsernameWithHash(ctx context.Context, username string) (admin *models.User, hash []byte, err error)
	GetUserByUsername(ctx context.Context, email string) (*models.User, error)

	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	DeleteUserByID(ctx context.Context, userID uuid.UUID) error

	ClearUserAvatarURL(ctx context.Context, userID uuid.UUID) error
}
type SessionRepository interface {
	GetUserSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (*models.Session, error)
	DeleteUserSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
	AddUserSession(ctx context.Context, session models.Session, sessionLimit int) error
}

type AvatarStorage interface {
	UploadUserAvatar(ctx context.Context, params dto.UploadUserAvatar) (path string, err error)
	GetUserAvatarURL(ctx context.Context, userID uuid.UUID) (string, error)
	DeleteUserAvatar(ctx context.Context, userID uuid.UUID) error
}

type FavoriteQuestionsRepository interface {
	MarkQuestionAsFavorite(ctx context.Context, questionID uuid.UUID, userID uuid.UUID) error
	UnmarkQuestionAsFavorite(ctx context.Context, questionID uuid.UUID, userID uuid.UUID) error
}

type MailClient interface {
	SendMail(to string, message mail.Message) error
}

type ConfirmationCodesRepository interface {
	StoreCode(ctx context.Context, data models.ConfirmationUserData) error
	CheckCode(ctx context.Context, email string, confirmationCode string) (*models.ConfirmationUserData, error)
	GetActiveCodesCount(ctx context.Context, email string) (int64, error)
	ClearAllCodes(ctx context.Context, email string) error
}

type ResetPasswordCodesRepository interface {
	StoreCode(ctx context.Context, data models.ResetPasswordData) error
	GetActiveCodesCount(ctx context.Context, email string) (int64, error)
	ClearAllCodes(ctx context.Context, email string) error
	CheckCode(ctx context.Context, email string, confirmationCode string) (*models.ResetPasswordData, error)
	DeleteCode(ctx context.Context, email, code string) error
}
