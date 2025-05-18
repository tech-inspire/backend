package clients

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tech-inspire/service/auth-service/internal/config"
	"go.uber.org/fx"
)

func NewPostgres(lc fx.Lifecycle, cfg *config.Config) (*pgxpool.Pool, error) {
	connectCfg, err := pgxpool.ParseConfig(cfg.DB.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("parse postgres dsn: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), connectCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return pool.Ping(ctx)
		},
		OnStop: func(_ context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}
