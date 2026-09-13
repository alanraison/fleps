package sql

import (
	"database/sql"
	"testing"
)

func setupDefaultPlayerData(tb testing.TB, db *sql.DB) {
	tb.Helper()
	pr := NewPlayerRepository(db)
	if err := pr.AddPlayer("Alan Raison", "alan.raison@gmail.com"); err != nil {
		tb.Fatalf("setupDefaultPlayerData returned error: %v", err)
	}
	if err := pr.AddPlayer("Another Player", "another.player@example.com"); err != nil {
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

	err := repo.AddPlayer("Alan Raison", "alan.raison@gmail.com")
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
