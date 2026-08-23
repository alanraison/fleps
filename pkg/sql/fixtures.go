package sql

import (
	"database/sql"
	"github.com/alanraison/predictions/pkg/model"
)

type FixtureRepository struct {
	db *sql.DB
}

func (r *FixtureRepository) AddFixtures(fixtures []model.Fixture) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO fixtures (home_team, away_team, date) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, fixture := range fixtures {
		_, err = stmt.Exec(fixture.HomeTeam.Key, fixture.AwayTeam.Key, fixture.Date)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
