package sql

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/alanraison/predictions/pkg/model"
)

func setupDefaultPlayerData(tb testing.TB, db *sql.DB) {
	tb.Helper()
	pr := NewPlayerRepository(db)
	if err := pr.AddPlayer("alan.raison@gmail.com", "Alan Raison"); err != nil {
		tb.Fatalf("setupDefaultPlayerData returned error: %v", err)
	}
}

func TestPlayerRepositoryReturnsEmptyListWhenNoPlayers(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewPlayerRepository(db)

	players, err := repo.ListPlayers()
	if err != nil {
		t.Fatalf("ListPlayers returned error: %v", err)
	}

	if len(players) != 0 {
		t.Fatalf("expected 0 players, got %d", len(players))
	}
}

func TestShouldAddPlayer(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewPlayerRepository(db)

	err := repo.AddPlayer("alan.raison@gmail.com", "Alan Raison")
	if err != nil {
		t.Fatalf("AddPlayer returned error: %v", err)
	}

	players, err := repo.ListPlayers()
	if err != nil {
		t.Fatalf("ListPlayers returned error: %v", err)
	}

	if len(players) != 1 {
		t.Fatalf("expected 1 player, got %d", len(players))
	}

	player := players[0]
	if player.Email != "alan.raison@gmail.com" {
		t.Fatalf("expected email 'alan.raison@gmail.com', got '%s'", player.Email)
	}
	if player.Name != "Alan Raison" {
		t.Fatalf("expected name 'Alan Raison', got '%s'", player.Name)
	}
}

func TestShouldAddPrediction(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	fr := NewFixtureRepository(db)
	repo := NewPredictionRepository(db)

	fs, err := fr.ListFixtures(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}
	f := fs[0]

	err = repo.AddPrediction("alan.raison@gmail.com", model.TeamKey(f.HomeTeam), model.TeamKey(f.AwayTeam), f.Date, 2, 1)
	if err != nil {
		t.Fatalf("AddPrediction returned error: %v", err)
	}

	predictions, err := repo.ListPredictions(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListPredictions returned error: %v", err)
	}
	if len(predictions) != 1 {
		t.Fatalf("expected 1 prediction, got %d", len(predictions))
	}
	prediction := predictions[0]
	if prediction.HomeTeam != f.HomeTeam {
		t.Fatalf("expected home team key '%s', got '%s'", f.HomeTeam, prediction.HomeTeam)
	}
	if prediction.RoundID != f.RoundID {
		t.Fatalf("expected round id '%s', got '%s'", f.RoundID, prediction.RoundID)
	}
	if prediction.AwayTeam != f.AwayTeam {
		t.Fatalf("expected away team key '%s', got '%s'", f.AwayTeam, prediction.AwayTeam)
	}
	if prediction.Date != f.Date {
		t.Fatalf("expected date '%v', got '%v'", f.Date, prediction.Date)
	}
	if prediction.HomeGoals != 2 {
		t.Fatalf("expected home goals 2, got %d", prediction.HomeGoals)
	}
	if prediction.AwayGoals != 1 {
		t.Fatalf("expected away goals 1, got %d", prediction.AwayGoals)
	}
}

func TestPredicitionsRepositoryShouldFailForUnknownFixture(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultPlayerData(t, db)

	repo := NewPredictionRepository(db)

	matchDate := time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC)
	err := repo.AddPrediction("alan.raison@gmail.com", "LEE", "MUN", matchDate, 2, 1)
	if err == nil {
		t.Fatalf("expected error when adding prediction for unknown fixture, got nil")
	}
	if !errors.Is(err, model.UnknownFixtureErr) {
		t.Fatalf("expected UnknownFixtureErr, got %v", err)
	}
}
