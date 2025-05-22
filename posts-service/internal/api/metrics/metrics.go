package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/extra/redisprometheus/v9"
	"github.com/redis/go-redis/v9"
)

func RegisterCollectors(redis redis.UniversalClient) {
	prometheus.MustRegister(
		redisprometheus.NewCollector("redis", "client", redis),
	)
}
