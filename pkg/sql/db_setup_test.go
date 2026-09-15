package sql

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// testDB is shared across all tests in the package; the container backing it
// is started once in TestMain and torn down after all tests have run.
var testDB *sql.DB

func TestMain(m *testing.M) {
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
		log.Fatalf("failed to start postgres container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate postgres container: %v", err)
		}
	}()

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get postgres connection string: %v", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open postgres connection: %v", err)
	}
	defer db.Close()

	// Keep one connection for PostgreSQL tests so schema and data are shared.
	db.SetMaxOpenConns(1)
	testDB = db

	os.Exit(m.Run())
}

// setupTestDB returns the shared test database along with a teardown func
// that truncates all tables so each test starts with a clean, isolated state.
func setupTestDB(tb testing.TB) (func(tb testing.TB), *sql.DB) {
	tb.Helper()

	return func(tb testing.TB) {
		tb.Helper()

		if _, err := testDB.Exec(`TRUNCATE TABLE
			leagues, teams, seasons, rounds, fixtures, results, players, predictions
			RESTART IDENTITY CASCADE;`); err != nil {
			tb.Fatalf("failed to truncate tables: %v", err)
		}
	}, testDB
}
