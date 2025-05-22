package service

import (
	"context"
	"time"

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

type ImageStorage interface {
	GenerateTempImageUpload(ctx context.Context, params dto.GenerateImageUploadURLParams, expire time.Duration) (*dto.GeneratedImageUpload, error)
	CreatePostImage(ctx context.Context, tempImageKey string, postID uuid.UUID) (*dto.CreatedPostImage, error)
	DeletePostImage(ctx context.Context, postID uuid.UUID) error
}

type PostsEventDispatcher interface {
	DispatchPostCreatedEvent(ctx context.Context, post *models.Post) error
	DispatchPostUpdatedEvent(ctx context.Context, post *models.Post, updatedAt time.Time) error
	DispatchPostDeletedEvent(ctx context.Context, post *models.Post, deletedAt time.Time) error
}

type PendingImagesRepository interface {
	Add(ctx context.Context, s3Key string, expiry time.Time) error
	Remove(ctx context.Context, keys ...string) error
	GetExpired(ctx context.Context, now time.Time) ([]string, error)
}
