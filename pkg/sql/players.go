package sql

import (
	"database/sql"
	"fmt"

	"github.com/alanraison/predictions/pkg/model"
)

type playerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *playerRepository {
	return &playerRepository{
		db: db,
	}
}

func (r *playerRepository) AddPlayer(email, name string) error {
	_, err := r.db.Exec("INSERT INTO players (email, name) VALUES (?, ?)", email, name)
	if err != nil {
		return fmt.Errorf("adding player %q: %w", email, err)
	}

	return nil
}

func (r *playerRepository) ListPlayers() ([]model.Player, error) {
	rows, err := r.db.Query("SELECT email, name FROM players")
	if err != nil {
		return nil, fmt.Errorf("listing players: %w", err)
	}
	defer rows.Close()

	var players []model.Player
	for rows.Next() {
		var player model.Player
		if err := rows.Scan(&player.Email, &player.Name); err != nil {
			return nil, fmt.Errorf("scanning player row: %w", err)
		}
		players = append(players, player)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating player rows: %w", err)
	}
	return players, nil
}
