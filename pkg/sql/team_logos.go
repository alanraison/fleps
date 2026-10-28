package sql

import (
	"database/sql"
	"fmt"

	"github.com/alanraison/fleps/pkg/model"
)

type teamLogoRepository struct {
	slugs map[model.TeamKey]string
}

// NewTeamLogoRepository loads every team's logo slug from the database.
func NewTeamLogoRepository(db *sql.DB) (*teamLogoRepository, error) {
	rows, err := db.Query("SELECT team_key, logo_slug FROM team_logos")
	if err != nil {
		return nil, fmt.Errorf("loading all team logos: %w", err)
	}
	defer rows.Close()

	r := &teamLogoRepository{slugs: make(map[model.TeamKey]string)}
	for rows.Next() {
		var key model.TeamKey
		var slug string
		if err := rows.Scan(&key, &slug); err != nil {
			return nil, fmt.Errorf("scanning team logo: %w", err)
		}
		r.slugs[key] = slug
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over team logos: %w", err)
	}
	return r, nil
}

func (r *teamLogoRepository) FindLogoSlugByTeamKey(key model.TeamKey) (string, error) {
	slug, ok := r.slugs[key]
	if !ok {
		return "", fmt.Errorf("no logo found for team %q: %w", key, model.UnknownLogoErr)
	}
	return slug, nil
}
