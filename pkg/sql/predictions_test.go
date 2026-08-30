package sql

import (
	"database/sql"
	"testing"
	"time"
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
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	fr := NewFixtureRepository(db)
	repo := NewPredictionRepository(db)

	from, _ := time.Parse("2006-01-02 15:04", "2023-10-01 00:00")
	to, _ := time.Parse("2006-01-02 15:04", "2023-10-02 23:59")

	fs, err := fr.ListFixtures(from, to, []string{"LEE"})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}
	f := fs[0]
	t.Log("finished listing fixtures")

	err = repo.AddPrediction("alan.raison@gmail.com", f.HomeTeam.Key, f.AwayTeam.Key, f.Date, 2, 1)
	if err != nil {
		t.Fatalf("AddPrediction returned error: %v", err)
	}

	predictions, err := repo.ListPredictions(from, to)
	if err != nil {
		t.Fatalf("ListPredictions returned error: %v", err)
	}
	if len(predictions) != 1 {
		t.Fatalf("expected 1 prediction, got %d", len(predictions))
	}
	prediction := predictions[0]
	if prediction.HomeTeam.Key != f.HomeTeam.Key {
		t.Fatalf("expected home team key '%s', got '%s'", f.HomeTeam.Key, prediction.HomeTeam.Key)
	}
	if prediction.AwayTeam.Key != f.AwayTeam.Key {
		t.Fatalf("expected away team key '%s', got '%s'", f.AwayTeam.Key, prediction.AwayTeam.Key)
	}
	if prediction.Date != f.Date {
		t.Fatalf("expected date '%v', got '%v'", f.Date, prediction.Date)
	}
	if prediction.HomeScore != 2 {
		t.Fatalf("expected home score 2, got %d", prediction.HomeScore)
	}
	if prediction.AwayScore != 1 {
		t.Fatalf("expected away score 1, got %d", prediction.AwayScore)
	}
}
