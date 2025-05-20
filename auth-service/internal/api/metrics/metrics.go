package metrics

import (
	"github.com/IBM/pgxpoolprometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/extra/redisprometheus/v9"
	"github.com/redis/go-redis/v9"
)

func RegisterCollectors(pg *pgxpool.Pool, redis *redis.Client) {
	prometheus.MustRegister(
		pgxpoolprometheus.NewCollector(pg, map[string]string{
			"database_name": pg.Config().ConnConfig.Database,
		}),
		redisprometheus.NewCollector("redis", "client", redis),
	)
}
