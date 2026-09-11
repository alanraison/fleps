package sql

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/alanraison/predictions/pkg/model"
)

func setupDefaultTeamData(tb testing.TB, db *sql.DB) {
	tb.Helper()

	if _, err := db.Exec(`INSERT INTO leagues (rank, name) VALUES (1, 'Premier League')`); err != nil {
		tb.Fatalf("failed to insert league test data: %v", err)
	}

	if _, err := db.Exec(`
		INSERT INTO teams (key, full_name, short_name, league) VALUES
		('LEE', 'Leeds United', 'Leeds', 1),
		('MUN', 'Manchester United', 'Man Utd', 1),
		('ARS', 'Arsenal', 'Arsenal', 1),
		('CHE', 'Chelsea', 'Chelsea', 1);
	`); err != nil {
		tb.Fatalf("failed to insert test data: %v", err)
	}
}

func TestTeamRepositoryFindTeamByKey_FindsExistingKey(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewTeamRepository(db)

	setupDefaultTeamData(t, db)
	team, err := repo.FindTeamByKey("LEE")
	if err != nil {
		t.Fatalf("FindTeamByKey returned error: %v", err)
	}
	if team == nil {
		t.Fatal("expected team, got nil")
	}
	if got, want := team.Key, model.TeamKey("LEE"); got != want {
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
	repo := NewTeamRepository(db)

	setupDefaultTeamData(t, db)

	team, err := repo.FindTeamByKey("XYZ")
	if err == nil {
		t.Fatal("expected error for unknown team, got nil")
	}
	if !errors.Is(err, model.UnknownTeamErr) {
		t.Fatalf("expected UnknownTeamErr, got %v", err)
	}
	if team != nil {
		t.Fatal("expected team not to be found")
	}
}

func TestTeamRepositoryAddTeam(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewTeamRepository(db)

	setupDefaultTeamData(t, db)

	newTeam := &model.Team{
		Key:       "TOT",
		FullName:  "Tottenham Hotspur",
		ShortName: "Spurs",
		League:    1,
	}
	if err := repo.AddTeam(newTeam); err != nil {
		t.Fatalf("AddTeam returned error: %v", err)
	}

	team, err := repo.FindTeamByKey("TOT")
	if err != nil {
		t.Fatalf("FindTeamByKey returned error: %v", err)
	}
	if team == nil {
		t.Fatal("expected team, got nil")
	}
	if got, want := team.Key, model.TeamKey("TOT"); got != want {
		t.Fatalf("Key = %q, want %q", got, want)
	}
	if got, want := team.FullName, "Tottenham Hotspur"; got != want {
		t.Fatalf("FullName = %q, want %q", got, want)
	}
	if got, want := team.ShortName, "Spurs"; got != want {
		t.Fatalf("ShortName = %q, want %q", got, want)
	}
}
