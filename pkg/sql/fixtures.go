package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alanraison/predictions/pkg/model"
)

type fixtureRepository struct {
	db *sql.DB
}

func (r *fixtureRepository) AddFixtures(fixtures []model.Fixture) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	stmt, err := tx.Prepare("INSERT INTO fixtures (home_team, away_team, date_time) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, fixture := range fixtures {
		_, err = stmt.Exec(fixture.HomeTeam.Key, fixture.AwayTeam.Key, fixture.Date)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("executing statement: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (r *fixtureRepository) ListFixtures(fromDate time.Time, toDate time.Time, teams []string) ([]model.Fixture, error) {
	rows, err := r.db.Query(`
		SELECT 
			home_team,
			home.full_name AS home_team_name,
			home.short_name AS home_team_short_name,
			away_team,
			away.full_name AS away_team_name,
			away.short_name AS away_team_short_name,
			date_time
		FROM 
			fixtures 
		JOIN 
			teams AS home ON fixtures.home_team = home.key 
		JOIN 
			teams AS away ON fixtures.away_team = away.key
		WHERE 
			date_time BETWEEN ? AND ?`,
		fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("reading fixtures: %w", err)
	}
	defer rows.Close()

	var fixtures []model.Fixture
	for rows.Next() {
		var homeTeamKey, homeTeamName, homeTeamShortName string
		var awayTeamKey, awayTeamName, awayTeamShortName string
		var dateTime time.Time
		if err := rows.Scan(&homeTeamKey, &homeTeamName, &homeTeamShortName, &awayTeamKey, &awayTeamName, &awayTeamShortName, &dateTime); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		fixtures = append(fixtures, model.Fixture{
			HomeTeam: model.Team{Key: homeTeamKey, FullName: homeTeamName, ShortName: homeTeamShortName},
			AwayTeam: model.Team{Key: awayTeamKey, FullName: awayTeamName, ShortName: awayTeamShortName},
			Date:     dateTime,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}
	return fixtures, nil
}
