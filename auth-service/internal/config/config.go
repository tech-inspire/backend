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
		Endpoint        string `env:"S3_ENDPOINT,required"`
		AccessKeyID     string `env:"S3_ACCESS_KEY_ID,required"`
		SecretAccessKey string `env:"S3_SECRET_ACCESS_KEY,required"`
		BucketName      string `env:"S3_BUCKET_NAME,required"`
		Region          string `env:"S3_REGION,required"`
		UseSSL          bool   `env:"S3_USE_SSL,required"`
	}

	ApplicationURL string `env:"APPLICATION_URL,required"`

	DisableStackTrace bool `env:"DISABLE_STACK_TRACE"`

	JWT struct {
		UserJWKPath          string        `env:"JWT_USER_KEY_PATH,required"`
		AccessTokenDuration  time.Duration `env:"JWT_ACCESS_TOKEN_DURATION,required"`
		RefreshTokenDuration time.Duration `env:"JWT_REFRESH_TOKEN_DURATION,required"`
	}

	Session struct {
		MaxAllowedSessionsPerUser int `env:"MAX_ALLOWED_SESSIONS_PER_USER,required"`
	}

	DB struct {
		PostgresDSN string `env:"POSTGRES_DSN,required"`
		RedisDSN    string `env:"REDIS_DSN,required"`
	}

	TestMode bool `env:"TEST_MODE"`

	SMTP struct {
		From     string `env:"EMAILS_SMTP_FROM,required"`
		Password string `env:"EMAILS_SMTP_PASSWORD,required"`
		Host     string `env:"EMAILS_SMTP_HOST,required"`
		Port     string `env:"EMAILS_SMTP_PORT,required"`
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
