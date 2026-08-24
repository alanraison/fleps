package sql

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTeamRepository(tb testing.TB) (func(tb testing.TB), *teamRepository) {
	tb.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		tb.Fatal(err)
	}

	if err := ApplyDefaultSchema(db); err != nil {
		tb.Fatal(err)
	}

	repo := &teamRepository{db: db}

	return func(tb testing.TB) {
		tb.Helper()
		db.Close()
	}, repo
}

func setupDefaultTeamData(tb testing.TB, repo *teamRepository) {
	tb.Helper()

	db := repo.db

	if _, err := db.Exec(`
		INSERT INTO teams (key, full_name, short_name, league) VALUES
		('LEE', 'Leeds United', 'Leeds', 1),
		('MUN', 'Manchester United', 'Man Utd', 1);
	`); err != nil {
		tb.Fatalf("failed to insert test data: %v", err)
	}
}

func TestTeamRepositoryFindTeamByKey_FindsExistingKey(t *testing.T) {
	teardown, repo := setupTeamRepository(t)
	defer teardown(t)

	setupDefaultTeamData(t, repo)
	team, ok, err := repo.FindByKey("LEE")
	if err != nil {
		t.Fatalf("FindByKey returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected team to be found")
	}
	if team == nil {
		t.Fatal("expected team, got nil")
	}
	if got, want := team.Key, "LEE"; got != want {
		t.Fatalf("Key = %q, want %q", got, want)
	}
	if got, want := team.FullName, "Leeds United"; got != want {
		t.Fatalf("FullName = %q, want %q", got, want)
	}
	if got, want := team.ShortName, "Leeds"; got != want {
		t.Fatalf("ShortName = %q, want %q", got, want)
	}
}

func TestTeamRepositoryFindTeamByKey_ReturnsNotFoundForNonExistingKey(t *testing.T) {
	teardown, repo := setupTeamRepository(t)
	defer teardown(t)

	setupDefaultTeamData(t, repo)

	team, ok, err := repo.FindByKey("XYZ")
	if err != nil {
		t.Fatalf("FindByKey returned error: %v", err)
	}
	if ok {
		t.Fatal("expected team not to be found")
	}
	if team != nil {
		t.Fatalf("expected nil team, got %+v", team)
	}
}
