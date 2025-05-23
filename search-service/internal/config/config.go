package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/go-errors/errors"
)

type Config struct {
	Server struct {
		Address            string   `env:"SERVER_ADDRESS,required"`
		MetricsAddress     string   `env:"SERVER_METRICS_ADDRESS,required"`
		ProxyHeader        string   `env:"SERVER_PROXY_HEADER,required"`
		CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	}

	DisableStackTrace bool `env:"DISABLE_STACK_TRACE"`

	AuthJWKSPath string `env:"JWKS_PATH,required"`

	Nats struct {
		URL             string `env:"NATS_URL,required"`
		PostsStreamName string `env:"POSTS_STREAM_NAME,required"`
	}

	ImageEmbeddings struct {
		ImageURLBasePath string `env:"IMAGE_URL_BASE_PATH,required"`
	}

	EmbeddingsClient struct {
		URL string `env:"EMBEDDINGS_CLIENT_URL,required"`
	}

	DB struct {
		PostgresDSN string `env:"POSTGRES_DSN,required"`
		// RedisDSN    string `env:"REDIS_DSN,required"`
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
