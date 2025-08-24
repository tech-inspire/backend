package clients

import (
	"context"

	"github.com/go-errors/errors"
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
