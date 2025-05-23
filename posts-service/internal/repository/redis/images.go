package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tech-inspire/backend/posts-service/internal/apperrors"
	"github.com/tech-inspire/backend/posts-service/internal/config"
)

// PendingImageUploadsRepository manages temporary image keys in a Redis sorted set.
// Each member's score represents its expiry timestamp (Unix seconds).
type PendingImageUploadsRepository struct {
	client redis.UniversalClient
	// key under which all temp image entries are stored
	setKey string
}

// NewPendingImageUploadsRepository creates a new TempImageRepository.
// setKey is the Redis key for the sorted set (e.g., "pending_uploads").
func NewPendingImageUploadsRepository(client redis.UniversalClient, cfg *config.Config) *PendingImageUploadsRepository {
	return &PendingImageUploadsRepository{client: client, setKey: cfg.Redis.PendingImagesSetKey}
}

// Add registers a new temporary image key with the given expiry time.
// s3Key is the object key in S3, expiry is when the key should be considered expired.
func (r *PendingImageUploadsRepository) Add(ctx context.Context, s3Key string, expiry time.Time) error {
	score := float64(expiry.Unix())
	z := redis.Z{
		Score:  score,
		Member: s3Key,
	}

	if err := r.client.ZAddNX(ctx, r.setKey, z).Err(); err != nil {
		return err
	}
	return nil
}

// Remove deletes one or more image keys from the pending set.
func (r *PendingImageUploadsRepository) Remove(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	removed, err := r.client.ZRem(ctx, r.setKey, toInterfaceSlice(keys)...).Result()
	if err != nil {
		return fmt.Errorf("zrem: %w", err)
	}

	if removed == 0 {
		return fmt.Errorf("%w: temp image not found", apperrors.ErrForbidden)
	}

	return nil
}

// GetExpired returns all image keys whose expiry is at or before "now".
func (r *PendingImageUploadsRepository) GetExpired(ctx context.Context, now time.Time) ([]string, error) {
	keys, err := r.client.ZRangeByScore(ctx, r.setKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: strconv.FormatInt(now.Unix(), 10),
	}).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func toInterfaceSlice(keys []string) []interface{} {
	r := make([]interface{}, len(keys))
	for i, k := range keys {
		r[i] = k
	}
	return r
}
