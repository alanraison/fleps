package sql

import (
	"database/sql"
	"fmt"

	"github.com/alanraison/predictions/pkg/model"
)

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *teamRepository {
	return &teamRepository{
		db: db,
	}
}

func (r *teamRepository) FindTeamByKey(key string) (team *model.Team, err error) {
	team = &model.Team{}
	err = r.db.
		QueryRow("SELECT key, full_name, short_name FROM teams WHERE key = ?", key).
		Scan(&team.Key, &team.FullName, &team.ShortName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("finding team by key: %w", err)
	}
	return team, nil
}
