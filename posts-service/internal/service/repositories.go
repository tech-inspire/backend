package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	"github.com/tech-inspire/backend/posts-service/internal/service/dto"
)

type Generator interface {
	GenerateString(length int) string
	NewUUID() uuid.UUID
}

type PostsRepository interface {
	UpdatePostByID(ctx context.Context, postID uuid.UUID, params dto.UpdatePostParams) error
	CreatePost(ctx context.Context, post *models.Post) error
	GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	GetPostsByIDs(ctx context.Context, postIDs []uuid.UUID) ([]*models.Post, error)
	DeletePostByID(ctx context.Context, postID uuid.UUID) error
}

type AvatarStorage interface {
	UploadUserAvatar(ctx context.Context, params dto.UploadUserAvatar) (path string, err error)
	GetUserAvatarURL(ctx context.Context, userID uuid.UUID) (string, error)
	DeleteUserAvatar(ctx context.Context, userID uuid.UUID) error
}

type EventDispatcher interface {
	DispatchPostDeletedEvent(ctx context.Context, post *models.Post) error
	DispatchPostCreatedEvent(ctx context.Context, post *models.Post) error
	DispatchPostUpdatedEvent(ctx context.Context, post *models.Post) error
}
