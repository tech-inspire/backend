package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tech-inspire/backend/posts-service/internal/apperrors"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	"github.com/tech-inspire/backend/posts-service/internal/service/dto"
)

type PostsService struct {
	repo PostsRepository
}

func (p PostsService) UpdatePostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID, params dto.UpdatePostParams) error {
	// TODO implement me
	panic("implement me")
}

func (p PostsService) CreatePost(ctx context.Context, post *models.Post) error {
	// TODO implement me
	panic("implement me")
}

func (p PostsService) GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	// TODO implement me
	panic("implement me")
}

func (p PostsService) GetPostsByIDs(ctx context.Context, postIDs []uuid.UUID) ([]*models.Post, error) {
	// TODO implement me
	panic("implement me")
}

func (p PostsService) DeletePostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error {
	post, err := p.repo.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("get post: %w", err)
	}

	if post.AuthorID != userID {
		return apperrors.ErrForbidden
	}

	err = p.repo.DeletePostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
}
