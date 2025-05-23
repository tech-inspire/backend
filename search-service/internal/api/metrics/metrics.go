package metrics

import (
	"github.com/IBM/pgxpoolprometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

func RegisterCollectors(pg *pgxpool.Pool) {
	prometheus.MustRegister(
		pgxpoolprometheus.NewCollector(pg, map[string]string{
			"database_name": pg.Config().ConnConfig.Database,
		}),
		// redisprometheus.NewCollector("redis", "client", redis),
	)
}
