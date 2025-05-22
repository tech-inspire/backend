package cache

import (
	"context"
	"fmt"
	"slices"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/tech-inspire/backend/posts-service/internal/apperrors"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	"github.com/tech-inspire/backend/posts-service/internal/repository/redis"
	"github.com/tech-inspire/backend/posts-service/internal/repository/scylla"
	"github.com/tech-inspire/backend/posts-service/internal/service/dto"
)

type PostsRepository struct {
	main  *scylla.PostsRepository
	cache *redis.PostRepository
}

func NewPostsRepository(main *scylla.PostsRepository, cache *redis.PostRepository) *PostsRepository {
	return &PostsRepository{main: main, cache: cache}
}

func (r PostsRepository) UpdatePostByID(ctx context.Context, postID uuid.UUID, params dto.UpdatePostParams) error {
	updatedPost, err := r.main.Update(ctx, postID, params)
	if err != nil {
		return fmt.Errorf("postgres: update post by id: %w", err)
	}

	err = r.cache.SetPostByID(ctx, updatedPost)
	if err != nil {
		return err
	}

	return nil
}

func (r PostsRepository) CreatePost(ctx context.Context, post *models.Post) error {
	err := r.main.Create(ctx, post)
	if err != nil {
		return fmt.Errorf("scylla: create post: %w", err)
	}

	err = r.cache.SetPostByID(ctx, post)
	if err != nil {
		return fmt.Errorf("redis: set post: %w", err)
	}

	return nil
}

func (r PostsRepository) GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	cachedPost, err := r.cache.GetPostByID(ctx, postID)
	if err != nil && !errors.Is(err, apperrors.ErrPostNotFound) {
		return nil, fmt.Errorf("cache: get post by id: %w", err)
	}
	if err == nil {
		return cachedPost, nil
	}

	post, err := r.main.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("scylla: get post by id: %w", err)
	}

	err = r.cache.SetPostByID(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("cache: set post: %w", err)
	}

	return post, nil
}

func (r PostsRepository) GetPostsByIDs(ctx context.Context, postIDs []uuid.UUID) ([]*models.Post, error) {
	res, err := r.cache.GetPostsByIDs(ctx, postIDs)
	if err != nil {
		return nil, fmt.Errorf("cache: get posts by ids: %w", err)
	}

	if len(res.MissingPostIDs) == 0 {
		return res.Posts, nil
	}

	missingPosts, err := r.main.GetMany(ctx, res.MissingPostIDs)
	if err != nil {
		return nil, fmt.Errorf("scylla: get posts by ids: %w", err)
	}

	err = r.cache.SetPosts(ctx, missingPosts)
	if err != nil {
		return nil, fmt.Errorf("cache: set posts: %w", err)
	}

	out := slices.Concat(res.Posts, missingPosts)
	return out, nil
}

func (r PostsRepository) DeletePostByID(ctx context.Context, postID uuid.UUID) error {
	err := r.main.Delete(ctx, postID)
	if err != nil {
		return fmt.Errorf("postgress: delete post by id: %w", err)
	}

	err = r.cache.DeletePostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("redis: delete post by id: %w", err)
	}

	return err
}
