package sql

import (
	"database/sql"
	"testing"
)

// setupDefaultPlayerData inserts default players into the database:
//
//   - Alan Raison - alan.raison@gmail.com
//   - Another Player - another.player@example.com
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

	players, err := repo.ListActivePlayers()
	if err != nil {
		t.Fatalf("ListActivePlayers returned error: %v", err)
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

	players, err := repo.ListActivePlayers()
	if err != nil {
		t.Fatalf("ListActivePlayers returned error: %v", err)
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

func TestShouldGetPlayerByEmail(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewPlayerRepository(db)

	setupDefaultPlayerData(t, db)

	player, err := repo.GetPlayerByEmail("alan.raison@gmail.com")
	if err != nil {
		t.Fatalf("GetPlayerByEmail returned error: %v", err)
	}
	if player == nil {
		t.Fatalf("expected to find player, got nil")
	}
	if player.Email != "alan.raison@gmail.com" {
		t.Fatalf("expected email 'alan.raison@gmail.com', got '%s'", player.Email)
	}
	if player.Name != "Alan Raison" {
		t.Fatalf("expected name 'Alan Raison', got '%s'", player.Name)
	}
}

func TestShouldReturnNilWhenPlayerNotFound(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewPlayerRepository(db)

	player, err := repo.GetPlayerByEmail("non.existent@example.com")
	if err != nil {
		t.Fatalf("GetPlayerByEmail returned error: %v", err)
	}
	if player != nil {
		t.Fatalf("expected nil, got player with email '%s'", player.Email)
	}
}

func TestShouldListActivePlayers(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewPlayerRepository(db)

	setupDefaultPlayerData(t, db)

	players, err := repo.ListActivePlayers()
	if err != nil {
		t.Fatalf("ListActivePlayers returned error: %v", err)
	}

	if len(players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(players))
	}

	expectedEmails := map[string]bool{
		"alan.raison@gmail.com":      true,
		"another.player@example.com": true,
	}
	for _, player := range players {
		if !expectedEmails[player.Email] {
			t.Fatalf("unexpected player email '%s'", player.Email)
		}
	}
}

func TestShouldDisablePlayerByEmail(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewPlayerRepository(db)

	setupDefaultPlayerData(t, db)

	err := repo.DisablePlayerByEmail("alan.raison@gmail.com")
	if err != nil {
		t.Fatalf("DisablePlayerByEmail returned error: %v", err)
	}

	player, err := repo.GetPlayerByEmail("alan.raison@gmail.com")
	if err != nil {
		t.Fatalf("GetPlayerByEmail returned error: %v", err)
	}
	if player == nil {
		t.Fatalf("expected to find player, got nil")
	}
	if player.Active {
		t.Fatalf("expected player to be disabled, but it is active")
	}
}

func TestShouldEnablePlayerByEmail(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewPlayerRepository(db)

	setupDefaultPlayerData(t, db)

	err := repo.DisablePlayerByEmail("alan.raison@gmail.com")
	if err != nil {
		t.Fatalf("DisablePlayerByEmail returned error: %v", err)
	}

	err = repo.EnablePlayerByEmail("alan.raison@gmail.com")
	if err != nil {
		t.Fatalf("EnablePlayerByEmail returned error: %v", err)
	}

	player, err := repo.GetPlayerByEmail("alan.raison@gmail.com")
	if err != nil {
		t.Fatalf("GetPlayerByEmail returned error: %v", err)
	}
	if player == nil {
		t.Fatalf("expected to find player, got nil")
	}
	if !player.Active {
		t.Fatalf("expected player to be enabled, but it is not active")
	}
}

func TestShouldNotListDisabledPlayers(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewPlayerRepository(db)

	setupDefaultPlayerData(t, db)

	err := repo.DisablePlayerByEmail("alan.raison@gmail.com")
	if err != nil {
		t.Fatalf("DisablePlayerByEmail returned error: %v", err)
	}

	players, err := repo.ListActivePlayers()
	if err != nil {
		t.Fatalf("ListActivePlayers returned error: %v", err)
	}

	for _, player := range players {
		if player.Email == "alan.raison@gmail.com" {
			t.Fatalf("disabled player should not be listed")
		}
	}
}
