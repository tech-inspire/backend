package clients

import (
	"context"
	"time"

	"github.com/go-errors/errors"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"github.com/tech-inspire/backend/posts-service/internal/config"
	"go.uber.org/fx"
)

func NewRedis(lc fx.Lifecycle, cfg *config.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.Redis.DSN)
	if err != nil {
		return nil, errors.Errorf("parse redis dsn: %w", err)
	}

	client := redis.NewClient(opts)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx).Err()
		},
		OnStop: func(_ context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}

func FiberRedisStorage(client *redis.Client) fiber.Storage {
	return &storage{client}
}

type storage struct {
	db *redis.Client
}

// Get value by key
func (s *storage) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	val, err := s.db.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

// Set key with value
func (s *storage) Set(key string, val []byte, exp time.Duration) error {
	// Ain't Nobody Got Time For That
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}
	return s.db.Set(context.Background(), key, val, exp).Err()
}

// Delete key by key
func (s *storage) Delete(key string) error {
	// Ain't Nobody Got Time For That
	if len(key) <= 0 {
		return nil
	}
	return s.db.Del(context.Background(), key).Err()
}

// Reset all keys
func (s *storage) Reset() error {
	return s.db.FlushDB(context.Background()).Err()
}

// Close the database
func (s *storage) Close() error {
	return s.db.Close()
}

// Return database client
func (s *storage) Conn() *redis.Client {
	return s.db
}
