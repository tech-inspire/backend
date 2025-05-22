package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	"github.com/tech-inspire/backend/posts-service/internal/service/dto"
)

type PostsService interface {
	GenerateTempImageUpload(ctx context.Context, params dto.GenerateImageUploadURLParams) (*dto.GeneratedImageUpload, error)
	UpdatePostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID, params dto.UpdatePostParams) error
	CreatePost(ctx context.Context, userID uuid.UUID, params dto.CreatePostParams) (*models.Post, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	GetPostsByIDs(ctx context.Context, postIDs []uuid.UUID) ([]*models.Post, error)
	DeletePostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error
}
