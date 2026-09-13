package sql

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestDB(tb testing.TB) (func(tb testing.TB), *sql.DB) {
	tb.Helper()
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"docker.io/library/postgres:18-alpine",
		postgres.WithInitScripts(
			"../../sql/01_teams.sql",
			"../../sql/02_fixtures.sql",
			"../../sql/03_predictions.sql",
		),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		tb.Fatal(err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		tb.Fatal(err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		tb.Fatal(err)
	}

	// Keep one connection for PostgreSQL tests so schema and data are shared.
	db.SetMaxOpenConns(1)

	return func(tb testing.TB) {
		tb.Helper()

		_ = db.Close()

		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			tb.Fatal(err)
		}
	}, db
}
