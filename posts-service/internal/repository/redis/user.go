package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tech-inspire/backend/posts-service/internal/apperrors"
	"github.com/tech-inspire/backend/posts-service/internal/config"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	"github.com/tech-inspire/backend/posts-service/pkg/generics"
)

type PostRepository struct {
	client redis.UniversalClient
	ttl    time.Duration
}

func NewPostRepository(client redis.UniversalClient, cfg *config.Config) *PostRepository {
	return &PostRepository{
		client: client,
		ttl:    cfg.Redis.PostsCacheTTL,
	}
}

func (r *PostRepository) key(postID uuid.UUID) string {
	return fmt.Sprintf("post:%s", postID.String())
}

func (r *PostRepository) GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	key := r.key(postID)

	data, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, apperrors.ErrPostNotFound // cache miss
	}
	if err != nil {
		return nil, errors.Errorf("redis: get '%s': %w", key, err)
	}

	var Post models.Post
	if err := json.Unmarshal([]byte(data), &Post); err != nil {
		return nil, errors.Errorf("unmarshal Post from json: %w", err)
	}

	return &Post, nil
}

type GetPostsByIDsOutput struct {
	Posts          []*models.Post
	MissingPostIDs []uuid.UUID
}

func (r *PostRepository) GetPostsByIDs(ctx context.Context, PostIDs []uuid.UUID) (*GetPostsByIDsOutput, error) {
	keys := generics.Convert(PostIDs, r.key)

	results, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &GetPostsByIDsOutput{
				MissingPostIDs: PostIDs,
			}, nil
		}
		return nil, errors.Errorf("redis: get '%v': %w", keys, err)
	}

	var out GetPostsByIDsOutput

	for i, data := range results {
		if data == nil {
			out.MissingPostIDs = append(out.MissingPostIDs, PostIDs[i])
			continue
		}

		str, ok := data.(string)
		if !ok {
			return nil, errors.Errorf("redis: data '%s': expected string, got %T", PostIDs[i], data)
		}

		Post := new(models.Post)
		if err = json.Unmarshal([]byte(str), Post); err != nil {
			return nil, errors.Errorf("unmarshal Post from json: %w", err)
		}

		out.Posts = append(out.Posts, Post)
	}

	return &out, nil
}

func (r *PostRepository) SetPostByID(ctx context.Context, post *models.Post) error {
	key := r.key(post.PostID)

	data, err := json.Marshal(post)
	if err != nil {
		return errors.Errorf("marshal post to json: %w", err)
	}

	err = r.client.Set(ctx, key, data, r.ttl).Err()
	if err != nil {
		return errors.Errorf("redis: set '%s': %w", key, err)
	}

	return nil
}

func (r *PostRepository) SetPosts(ctx context.Context, posts []*models.Post) error {
	for _, Post := range posts {
		err := r.SetPostByID(ctx, Post)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostRepository) DeletePostByID(ctx context.Context, id uuid.UUID) error {
	key := r.key(id)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return errors.Errorf("redis: delete '%s': %w", key, err)
	}
	return nil
}
