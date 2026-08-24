package sql

import (
	"database/sql"

	"github.com/alanraison/predictions/pkg/model"
)

type playerRepository struct {
	db *sql.DB
}

func (r *playerRepository) AddPlayer(email string, name string) error {
	_, err := r.db.Exec("INSERT INTO players (email, name) VALUES (?, ?)", email, name)
	return err
}

func (r *playerRepository) ListPlayers() ([]model.Player, error) {
	rows, err := r.db.Query("SELECT email, name FROM players")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []model.Player
	for rows.Next() {
		var player model.Player
		if err := rows.Scan(&player.Email, &player.Name); err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return players, nil
}
