package sql

import (
	"database/sql"
	"fmt"

	"github.com/alanraison/predictions/pkg/model"
)

type TeamRepository struct {
	db *sql.DB
}

func (r *TeamRepository) FindByKey(key string) (team *model.Team, ok bool, err error) {
	team = &model.Team{}
	err = r.db.
		QueryRow("SELECT key, full_name, short_name FROM teams WHERE key = ?", key).
		Scan(&team.Key, &team.FullName, &team.ShortName)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("finding team by key: %w", err)
	}
	return team, true, nil
}
