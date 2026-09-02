package sql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/alanraison/predictions/pkg/model"
)

type fixtureRepository struct {
	db *sql.DB
}

func NewFixtureRepository(db *sql.DB) *fixtureRepository {
	return &fixtureRepository{
		db: db,
	}
}

var (
	noTeamsQuery = `
		SELECT 
			home_team,
			away_team,
			date_time
		FROM 
			fixtures
		JOIN
		  teams h ON fixtures.home_team = h.key
		JOIN
		  teams a ON fixtures.away_team = a.key
		WHERE 
			date_time BETWEEN ? AND ?`
	teamsQuery = `
		SELECT 
			home_team,
			away_team,
			date_time
		FROM 
			fixtures
		JOIN
		  teams h ON fixtures.home_team = h.key
		JOIN
		  teams a ON fixtures.away_team = a.key
		WHERE 
			date_time BETWEEN ? AND ? AND (home_team IN (%s) OR away_team IN (%s))`
	noTeamsResultsQuery = `
		SELECT
			home_team,
			away_team,
			date_time,
			home_goals,
			away_goals
		FROM 
			fixtures
		JOIN
		  teams h ON fixtures.home_team = h.key
		JOIN
		  teams a ON fixtures.away_team = a.key
		JOIN
			results ON fixtures.id = results.fixture_id
		WHERE 
			date_time BETWEEN ? AND ?`
	teamsResultsQuery = `
		SELECT
			home_team,
			away_team,
			date_time,
			home_goals,
			away_goals
		FROM 
			fixtures
		JOIN
		  teams h ON fixtures.home_team = h.key
		JOIN
		  teams a ON fixtures.away_team = a.key
		JOIN
			results ON fixtures.id = results.fixture_id
		WHERE 
			date_time BETWEEN ? AND ? AND (home_team IN (%s) OR away_team IN (%s))`
)

func fmtTeamsQuery(teams []string) string {
	n := len(teams)
	if n <= 0 {
		return ""
	}
	ph := strings.Repeat("?,", n-1) + "?"
	return fmt.Sprintf(teamsQuery, ph, ph)
}

func fmtTeamsResultsQuery(teams []string) string {
	n := len(teams)
	if n <= 0 {
		return ""
	}
	ph := strings.Repeat("?,", n-1) + "?"
	return fmt.Sprintf(teamsResultsQuery, ph, ph)
}

func (r *fixtureRepository) AddFixtures(fixtures []model.Fixture) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	stmt, err := tx.Prepare(`
	INSERT INTO
		fixtures (
			home_team, 
			away_team, 
			date_time
		) 
	SELECT
		 	?1, ?2, ?3
	FROM
		teams h,
		teams a
	WHERE 
		h.key = ?1 
	AND a.key = ?2`)
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, fixture := range fixtures {
		res, err := stmt.Exec(fixture.HomeTeam, fixture.AwayTeam, fixture.Date)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("executing statement: %w", err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("checking rows affected: %w", err)
		}
		if rows == 0 {
			tx.Rollback()
			return fmt.Errorf("no fixture inserted for the given teams and date; home team: %s, away team: %s, date: %v. Operation has been rolled back", fixture.HomeTeam, fixture.AwayTeam, fixture.Date)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (r *fixtureRepository) ListFixtures(fromDate time.Time, toDate time.Time, teams []string) ([]model.Fixture, error) {
	var rows *sql.Rows
	var err error
	if len(teams) == 0 {
		rows, err = r.db.Query(noTeamsQuery, fromDate, toDate)
		if err != nil {
			return nil, fmt.Errorf("reading fixtures: %w", err)
		}
	} else {
		args := make([]any, 2+len(teams)*2)
		args[0] = fromDate
		args[1] = toDate
		for i, team := range teams {
			args[2+i] = team
			args[2+len(teams)+i] = team
		}
		rows, err = r.db.Query(fmtTeamsQuery(teams), args...)
		if err != nil {
			return nil, fmt.Errorf("reading fixtures: %w", err)
		}
	}
	defer rows.Close()

	var fixtures []model.Fixture
	for rows.Next() {
		var homeTeamKey string
		var awayTeamKey string
		var dateTime time.Time
		if err := rows.Scan(&homeTeamKey, &awayTeamKey, &dateTime); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		fixtures = append(fixtures, model.Fixture{
			HomeTeam: model.TeamKey(homeTeamKey),
			AwayTeam: model.TeamKey(awayTeamKey),
			Date:     dateTime,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}
	return fixtures, nil
}

func (r *fixtureRepository) AddResult(
	homeTeam, awayTeam model.TeamKey,
	date time.Time,
	homeScore, awayScore int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	stmt, err := tx.Prepare(`
		INSERT INTO results (
				fixture_id, 
				home_goals, 
				away_goals
			)
		SELECT
			id,
			?4,
			?5
		FROM fixtures
		WHERE 
			home_team = ?1 
		AND away_team = ?2 
		AND datetime(date_time) = datetime(?3)`)
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(homeTeam, awayTeam, date, homeScore, awayScore)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("executing statement: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rows == 0 {
		tx.Rollback()
		return fmt.Errorf("no fixture found for the given teams and date")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (r *fixtureRepository) ListResults(fromDate time.Time, toDate time.Time, teams []string) ([]model.Result, error) {
	var rows *sql.Rows
	var err error
	if len(teams) == 0 {
		rows, err = r.db.Query(noTeamsResultsQuery, fromDate, toDate)
		if err != nil {
			return nil, fmt.Errorf("reading results: %w", err)
		}
	} else {
		args := make([]any, 2+len(teams)*2)
		args[0] = fromDate
		args[1] = toDate
		for i, team := range teams {
			args[2+i] = team
			args[2+len(teams)+i] = team
		}
		rows, err = r.db.Query(fmtTeamsResultsQuery(teams), args...)
		if err != nil {
			return nil, fmt.Errorf("reading results: %w", err)
		}
	}
	defer rows.Close()

	var results []model.Result
	for rows.Next() {
		var homeTeamKey string
		var awayTeamKey string
		var homeGoals, awayGoals int
		var dateTime time.Time
		if err := rows.Scan(
			&homeTeamKey,
			&awayTeamKey,
			&dateTime,
			&homeGoals,
			&awayGoals); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		results = append(results, model.Result{
			Fixture: model.Fixture{
				HomeTeam: model.TeamKey(homeTeamKey),
				AwayTeam: model.TeamKey(awayTeamKey),
				Date:     dateTime,
			},
			HomeScore: homeGoals,
			AwayScore: awayGoals,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}
	return results, nil
}
