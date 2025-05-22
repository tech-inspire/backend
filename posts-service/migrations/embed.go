package migrations

import (
	"context"
	"embed"
	"fmt"

	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/migrate"
)

//go:embed *.cql
var embedMigrations embed.FS

func ApplyMigrations(ctx context.Context, session gocqlx.Session) error {
	err := migrate.FromFS(ctx, session, embedMigrations)
	if err != nil {
		return fmt.Errorf("apply cql migrations: %w", err)
	}

	return nil
}
