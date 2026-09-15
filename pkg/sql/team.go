package sql

import (
	"database/sql"
	"fmt"

	"github.com/alanraison/fleps/pkg/model"
)

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *teamRepository {
	return &teamRepository{
		db: db,
	}
}

func (r *teamRepository) AddTeam(newTeam *model.Team) error {
	_, err := r.db.Exec(
		`INSERT INTO 
			teams (
				key, 
				full_name, 
				short_name, 
				league
			) VALUES ($1, $2, $3, $4)`,
		newTeam.Key, newTeam.FullName, newTeam.ShortName, newTeam.League,
	)
	if err != nil {
		return fmt.Errorf("adding team: %w", err)
	}
	return nil
}

func (r *teamRepository) FindTeamByKey(key model.TeamKey) (team *model.Team, err error) {
	team = &model.Team{}
	err = r.db.
		QueryRow("SELECT key, full_name, short_name, league FROM teams WHERE key = $1", key).
		Scan(&team.Key, &team.FullName, &team.ShortName, &team.League)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no team found with key %q: %w", key, model.UnknownTeamErr)
	}
	if err != nil {
		return nil, fmt.Errorf("finding team by key: %w", err)
	}
	return team, nil
}
