package sql

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/alanraison/predictions/pkg/model"
)

// setupDefaultFixtureData inserts default fixtures into the database:
//
//  * R1: LEE vs MUN on 2023-10-01 15:00:00
//  * R1: ARS vs CHE on 2023-10-02 16:00:00
//  * R2: LEE vs CHE on 2023-10-03 15:00:00
//
// Depends on setupDefaultTeamData and setupDefaultRoundData being called before it.
func setupDefaultFixtureData(tb testing.TB, db *sql.DB) {
	tb.Helper()

	if _, err := db.Exec(`
		INSERT INTO fixtures (round_id, home_team, away_team, date_time) VALUES
		('R1', 'LEE', 'MUN', '2023-10-01 15:00:00'),
		('R1', 'ARS', 'CHE', '2023-10-02 16:00:00'),
		('R2', 'LEE', 'CHE', '2023-10-03 15:00:00');
	`); err != nil {
		tb.Fatalf("failed to insert test data: %v", err)
	}
}

// setupDefaultSeasonData inserts default seasons into the database:
//
//  * S1
func setupDefaultSeasonData(tb testing.TB, db *sql.DB) {
	tb.Helper()

	if _, err := db.Exec(`
		INSERT INTO seasons (id) VALUES
		('S1');
	`); err != nil {
		tb.Fatalf("failed to insert season test data: %v", err)
	}
}

// setupDefaultRoundData inserts default rounds into the database:
//
//  * R1: Round 1 of season S1
//  * R2: Round 2 of season S1
func setupDefaultRoundData(tb testing.TB, db *sql.DB) {
	tb.Helper()

	setupDefaultSeasonData(tb, db)
	if _, err := db.Exec(`
		INSERT INTO rounds (id, season_id) VALUES
		('R1', 'S1'),
		('R2', 'S1');
	`); err != nil {
		tb.Fatalf("failed to insert round test data: %v", err)
	}
}

func TestFixtureRepositoryReturnsEmptyListWhenNoFixtures(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	fixtures, err := repo.ListFixtures(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 0 {
		t.Fatalf("expected 0 fixtures, got %d", len(fixtures))
	}
}

func TestFixtureRepositoryListFixtures(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)

	fixtures, err := repo.ListFixtures(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}
	if len(fixtures) != 2 {
		t.Fatalf("expected 2 fixtures, got %d", len(fixtures))
	}
}

func TestFixtureRepositoryAddFixtures(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)

	fixturesToAdd := []model.Fixture{
		{FixtureKey: model.FixtureKey{RoundID: model.RoundID("R1"), HomeTeam: model.TeamKey("LEE"), AwayTeam: model.TeamKey("MUN")}, Date: time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC)},
		{FixtureKey: model.FixtureKey{RoundID: model.RoundID("R1"), HomeTeam: model.TeamKey("ARS"), AwayTeam: model.TeamKey("CHE")}, Date: time.Date(2023, 10, 2, 16, 0, 0, 0, time.UTC)},
	}

	err := repo.AddFixtures(fixturesToAdd)
	if err != nil {
		t.Fatalf("AddFixtures returned error: %v", err)
	}

	fixtures, err := repo.ListFixtures(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 2 {
		t.Fatalf("expected 2 fixtures, got %d", len(fixtures))
	}
	if got, want := fixtures[0].RoundID, model.RoundID("R1"); got != want {
		t.Fatalf("fixtures[0].RoundID = %q, want %q", got, want)
	}
}

func TestFixtureRepositoryShouldFailToAddUnknownTeam(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)

	fixturesToAdd := []model.Fixture{
		{FixtureKey: model.FixtureKey{RoundID: model.RoundID("R1"), HomeTeam: model.TeamKey("LEE"), AwayTeam: model.TeamKey("MUN")}, Date: time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC)},
		{FixtureKey: model.FixtureKey{RoundID: model.RoundID("R1"), HomeTeam: model.TeamKey("BBB"), AwayTeam: model.TeamKey("CCC")}, Date: time.Date(2023, 10, 2, 16, 0, 0, 0, time.UTC)},
	}

	err := repo.AddFixtures(fixturesToAdd)
	if err == nil {
		t.Fatalf("expected error when adding fixture with unknown team, got nil")
	}

	fixtures, err := repo.ListFixtures(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}
	if len(fixtures) != 0 {
		t.Fatalf("expected 0 fixtures after rollback, got %d", len(fixtures))
	}
}

func TestFixtureRepositoryShouldFailToAddUnknownRound(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)

	fixturesToAdd := []model.Fixture{
		{FixtureKey: model.FixtureKey{RoundID: model.RoundID("R9"), HomeTeam: model.TeamKey("LEE"), AwayTeam: model.TeamKey("MUN")}, Date: time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC)},
	}

	err := repo.AddFixtures(fixturesToAdd)
	if err == nil {
		t.Fatal("expected error when adding fixture with unknown round, got nil")
	}
	if !errors.Is(err, model.UnknownRoundErr) {
		t.Fatalf("expected UnknownRoundErr, got %v", err)
	}
}
