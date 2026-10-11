package sql

import (
	"database/sql"
	"fmt"

	"github.com/alanraison/fleps/pkg/model"
)

type teamRepository struct {
	db    *sql.DB
	cache map[model.TeamKey]*model.Team
}

func NewTeamRepository(db *sql.DB) (*teamRepository, error) {
	r := &teamRepository{
		db:    db,
		cache: make(map[model.TeamKey]*model.Team),
	}
	if err := r.loadAllTeams(); err != nil {
		return nil, err
	}
	return r, nil
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
	r.cache[newTeam.Key] = newTeam
	return nil
}

func (r *teamRepository) FindTeamByKey(key model.TeamKey) (team *model.Team, err error) {
	team = r.cache[key]
	if team == nil {
		return nil, fmt.Errorf("no team found with key %q: %w", key, model.UnknownTeamErr)
	}
	return team, nil
}

func (r *teamRepository) loadAllTeams() error {
	rows, err := r.db.Query("SELECT key, full_name, short_name, league FROM teams")
	if err != nil {
		return fmt.Errorf("loading all teams: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		team := &model.Team{}
		if err := rows.Scan(&team.Key, &team.FullName, &team.ShortName, &team.League); err != nil {
			return fmt.Errorf("scanning team: %w", err)
		}
		r.cache[team.Key] = team
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterating over teams: %w", err)
	}
	return nil
}
