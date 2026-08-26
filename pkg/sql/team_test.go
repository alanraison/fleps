package sql

import (
	"database/sql"
	"testing"
)

func setupDefaultTeamData(tb testing.TB, db *sql.DB) {
	tb.Helper()

	if _, err := db.Exec(`INSERT INTO leagues (rank, name) VALUES (1, 'Premier League')`); err != nil {
		tb.Fatalf("failed to insert league test data: %v", err)
	}

	if _, err := db.Exec(`
		INSERT INTO teams (key, full_name, short_name, league) VALUES
		('LEE', 'Leeds United', 'Leeds', 1),
		('MNU', 'Manchester United', 'Man Utd', 1),
		('ARS', 'Arsenal', 'Arsenal', 1),
		('CHE', 'Chelsea', 'Chelsea', 1);
	`); err != nil {
		tb.Fatalf("failed to insert test data: %v", err)
	}
}

func TestTeamRepositoryFindTeamByKey_FindsExistingKey(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := &teamRepository{db: db}

	setupDefaultTeamData(t, db)
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
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := &teamRepository{db: db}

	setupDefaultTeamData(t, db)

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
