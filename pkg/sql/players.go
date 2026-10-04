package sql

import (
	"database/sql"
	"fmt"

	"github.com/alanraison/fleps/pkg/model"
)

type playerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *playerRepository {
	return &playerRepository{
		db: db,
	}
}

func (r *playerRepository) AddPlayer(name, email string) error {
	_, err := r.db.Exec("INSERT INTO players (name, email) VALUES ($1, $2)", name, email)
	if err != nil {
		return fmt.Errorf("adding player %q: %w", email, err)
	}

	return nil
}

func (r *playerRepository) GetPlayerByEmail(email string) (*model.Player, error) {
	row := r.db.QueryRow("SELECT email, name, active FROM players WHERE email = $1", email)
	var player model.Player
	if err := row.Scan(&player.Email, &player.Name, &player.Active); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getting player by email %q: %w", email, err)
	}
	return &player, nil
}

func (r *playerRepository) ListActivePlayers() ([]model.Player, error) {
	rows, err := r.db.Query("SELECT email, name FROM players WHERE active = TRUE")
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

func (r *playerRepository) DisablePlayerByEmail(email string) error {
	_, err := r.db.Exec("UPDATE players SET active = FALSE WHERE email = $1", email)
	if err != nil {
		return fmt.Errorf("disabling player by email %q: %w", email, err)
	}
	return nil
}

func (r *playerRepository) EnablePlayerByEmail(email string) error {
	_, err := r.db.Exec("UPDATE players SET active = TRUE WHERE email = $1", email)
	if err != nil {
		return fmt.Errorf("enabling player by email %q: %w", email, err)
	}
	return nil
}
