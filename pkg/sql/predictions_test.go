package sql

import (
	"database/sql"
	"testing"

	"github.com/alanraison/predictions/pkg/model"
	_ "github.com/mattn/go-sqlite3"
)

func setupPlayerRepository(tb testing.TB) (func(tb testing.TB), model.PlayerRepository) {
	tb.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		tb.Fatal(err)
	}

	if err := ApplyDefaultSchema(db); err != nil {
		tb.Fatal(err)
	}

	repo := &playerRepository{db: db}

	return func(tb testing.TB) {
		tb.Helper()
		db.Close()
	}, repo
}

func TestPlayerRepositoryReturnsEmptyListWhenNoPlayers(t *testing.T) {
	teardown, repo := setupPlayerRepository(t)
	defer teardown(t)

	players, err := repo.ListPlayers()
	if err != nil {
		t.Fatalf("ListPlayers returned error: %v", err)
	}

	if len(players) != 0 {
		t.Fatalf("expected 0 players, got %d", len(players))
	}
}

func TestShouldAddPlayer(t *testing.T) {
	teardown, repo := setupPlayerRepository(t)
	defer teardown(t)

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
