package migrations

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func ApplyMigrations(pool *pgxpool.Pool) error {
	if os.Getenv("APPLY_MIGRATIONS") != "true" {
		slog.Warn("skipped applying migrations")
		return nil
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	return goose.Up(stdlib.OpenDBFromPool(pool), ".")
}
