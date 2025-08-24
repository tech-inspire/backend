package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/go-errors/errors"
)

type Server struct {
	Address            string   `env:"SERVER_ADDRESS,required"`
	MetricsAddress     string   `env:"SERVER_METRICS_ADDRESS,required"`
	ProxyHeader        string   `env:"SERVER_PROXY_HEADER,required"`
	CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
}

type Nats struct {
	URL             string `env:"NATS_URL,required"`
	PostsStreamName string `env:"POSTS_STREAM_NAME,required"`
}

type ImageEmbeddings struct {
	ImageURLBasePath string `env:"IMAGE_URL_BASE_PATH,required"`
}

type Database struct {
	PostgresDSN string `env:"POSTGRES_DSN,required"`
	// RedisDSN    string `env:"REDIS_DSN,required"`
}

type Config struct {
	Server

	Nats

	ImageEmbeddings

	Database

	EmbeddingsClient struct {
		URL string `env:"EMBEDDINGS_CLIENT_URL,required"`
	}

	DisableStackTrace bool   `env:"DISABLE_STACK_TRACE"`
	AuthJWKSPath      string `env:"JWKS_PATH,required"`
}

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, errors.Errorf("parse env: %w", err)
	}

	return &cfg, nil
}
