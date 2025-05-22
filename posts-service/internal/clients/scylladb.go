package clients

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/tech-inspire/backend/posts-service/internal/config"
	"go.uber.org/fx"
)

func NewScyllaDBClient(lc fx.Lifecycle, config *config.Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(config.ScyllaDB.Hosts...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username:              config.ScyllaDB.Username,
		Password:              config.ScyllaDB.Password,
		AllowedAuthenticators: nil,
	}
	cluster.Keyspace = config.ScyllaDB.Keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(_ context.Context) error {
			session.Close()
			return nil
		},
	})

	return session, nil
}
