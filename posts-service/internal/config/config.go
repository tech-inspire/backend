package config

import (
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/go-errors/errors"
)

type Config struct {
	Server struct {
		Address            string   `env:"SERVER_ADDRESS,required"`
		MetricsAddress     string   `env:"SERVER_METRICS_ADDRESS,required"`
		ProxyHeader        string   `env:"SERVER_PROXY_HEADER,required"`
		CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
		DebugCORS          bool     `env:"CORS_DEBUG" envDefault:"false"`
	}

	S3 struct {
		Endpoint         string `env:"S3_ENDPOINT,required"`
		BucketName       string `env:"S3_BUCKET_NAME,required"`
		MinioResolveMode bool   `env:"S3_MINIO_RESOLVE_MODE" envDefault:"false"`
	}

	DisableStackTrace bool `env:"DISABLE_STACK_TRACE"`

	AuthJWKSPath string `env:"JWKS_PATH,required"`

	ScyllaDB struct {
		Hosts    []string `env:"SCYLLA_HOSTS,required"`
		Username string   `env:"SCYLLA_USERNAME,required"`
		Password string   `env:"SCYLLA_PASSWORD,required"`
		Keyspace string   `env:"SCYLLA_KEYSPACE,required"`
	}

	Nats struct {
		URL             string `env:"NATS_URL,required"`
		PostsStreamName string `env:"POSTS_STREAM_NAME,required"`
	}

	Redis struct {
		DSN                 string        `env:"REDIS_DSN,required"`
		PendingImagesSetKey string        `env:"REDIS_PENDING_IMAGES_SET_KEY" envDefault:"pending_uploads"`
		PostsCacheTTL       time.Duration `env:"REDIS_POSTS_CACHE_TTL" envDefault:"15m"`
	}
}

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, errors.Errorf("parse env: %w", err)
	}

	return &cfg, nil
}
