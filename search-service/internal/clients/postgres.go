package clients

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
	"github.com/tech-inspire/backend/search-service/internal/config"
	"go.uber.org/fx"
)

func NewPostgres(lc fx.Lifecycle, cfg *config.Config) (*pgxpool.Pool, error) {
	connectCfg, err := pgxpool.ParseConfig(cfg.Database.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("parse postgres dsn: %w", err)
	}

	connectCfg.AfterConnect = pgxvec.RegisterTypes

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
