package migrations

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/migrate"
	"go.uber.org/fx"
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

func ApplyMigrationsFX(lc fx.Lifecycle, session gocqlx.Session) error {
	if os.Getenv("APPLY_MIGRATIONS") != "true" {
		slog.Warn("skipped applying migrations")
		return nil
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return ApplyMigrations(ctx, session)
		},
	})
	return nil
}
