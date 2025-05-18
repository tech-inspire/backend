package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/tech-inspire/service/auth-service/internal/models"
	"github.com/tech-inspire/service/auth-service/internal/service/dto"
)

type AuthService interface {
	GetSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (*models.Session, error)
	DeleteSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
	Register(ctx context.Context, params dto.RegisterParams) (*dto.RegisterOutput, error)
	ConfirmRegistrationByCode(ctx context.Context, email string, code string) (*dto.LoginOutput, error)
	SendResetPasswordCode(ctx context.Context, email string) error
	CheckResetPasswordCode(ctx context.Context, email string, code string) error
	ConfirmResetPasswordByCode(ctx context.Context, email string, code string, password string) error
	LoginByEmail(ctx context.Context, email string, password string) (*dto.LoginOutput, error)
	LoginByUsername(ctx context.Context, username string, password string) (*dto.LoginOutput, error)
	RefreshSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, sessionToken string) (*models.User, error)
}

type UserService interface {
	GetUserInfoByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	DeleteUserByID(ctx context.Context, userID uuid.UUID) error
	UpdateUser(ctx context.Context, userID uuid.UUID, params dto.UpdateUsersInput) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.GetUserByIDOutput, error)
	GetCurrentUserByID(ctx context.Context, userID uuid.UUID) (*dto.GetCurrentUser, error)
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]dto.GetUserByIDOutput, error)
	GetUsersInfoByID(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error)
	GetUsers(ctx context.Context, params dto.GetUsersParams) (*dto.GetUsersOutput, error)
}

type AvatarService interface {
	UploadUserAvatar(ctx context.Context, params dto.UploadUserAvatar) error
	DeleteProfileAvatar(ctx context.Context, userID uuid.UUID) error
}
