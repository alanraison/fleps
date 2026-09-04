package sql

import (
	"database/sql"
	"errors"
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
			round_id,
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
			round_id,
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
			round_id,
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
			round_id,
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

func rollbackTransaction(tx *sql.Tx, cause error, context string) error {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return fmt.Errorf("%s: %w", context, errors.Join(cause, fmt.Errorf("rolling back transaction: %w", rollbackErr)))
	}

	return cause
}

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

func (r *fixtureRepository) validateFixtureReferences(fixture model.Fixture) error {
	if err := r.ensureRoundExists(fixture.RoundID); err != nil {
		return fmt.Errorf("validating round %q for fixture %s vs %s: %w", fixture.RoundID, fixture.HomeTeam, fixture.AwayTeam, err)
	}

	teamRepo := NewTeamRepository(r.db)
	if _, err := teamRepo.FindTeamByKey(fixture.HomeTeam); err != nil {
		return fmt.Errorf("validating home team %q for fixture in round %q: %w", fixture.HomeTeam, fixture.RoundID, err)
	}
	if _, err := teamRepo.FindTeamByKey(fixture.AwayTeam); err != nil {
		return fmt.Errorf("validating away team %q for fixture in round %q: %w", fixture.AwayTeam, fixture.RoundID, err)
	}

	return nil
}

func (r *fixtureRepository) ensureRoundExists(roundID model.RoundID) error {
	var foundRoundID string
	err := r.db.QueryRow("SELECT id FROM rounds WHERE id = ?", roundID).Scan(&foundRoundID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("finding round %q: %w", roundID, model.UnknownRoundErr)
	}
	if err != nil {
		return fmt.Errorf("finding round %q: %w", roundID, err)
	}

	return nil
}

func (r *fixtureRepository) AddFixtures(fixtures []model.Fixture) error {
	for _, fixture := range fixtures {
		if err := r.validateFixtureReferences(fixture); err != nil {
			return fmt.Errorf("validating fixture references: %w", err)
		}
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	stmt, err := tx.Prepare(`
		INSERT INTO fixtures (
			round_id,
			home_team,
			away_team,
			date_time
		) VALUES (?1, ?2, ?3, ?4)`)
	if err != nil {
		return rollbackTransaction(tx, fmt.Errorf("preparing add fixtures statement: %w", err), "rolling back add fixtures transaction")
	}
	defer stmt.Close()

	for _, fixture := range fixtures {
		res, err := stmt.Exec(fixture.RoundID, fixture.HomeTeam, fixture.AwayTeam, fixture.Date)
		if err != nil {
			return rollbackTransaction(tx, fmt.Errorf("executing add fixture statement for round %q and teams %s vs %s: %w", fixture.RoundID, fixture.HomeTeam, fixture.AwayTeam, err), "rolling back add fixtures transaction")
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return rollbackTransaction(tx, fmt.Errorf("checking affected rows for round %q and teams %s vs %s: %w", fixture.RoundID, fixture.HomeTeam, fixture.AwayTeam, err), "rolling back add fixtures transaction")
		}
		if rows == 0 {
			return rollbackTransaction(tx, fmt.Errorf("no fixture inserted for round %q, home team %s, away team %s, and date %v", fixture.RoundID, fixture.HomeTeam, fixture.AwayTeam, fixture.Date), "rolling back add fixtures transaction")
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
		var roundID string
		var homeTeamKey string
		var awayTeamKey string
		var dateTime time.Time
		if err := rows.Scan(&roundID, &homeTeamKey, &awayTeamKey, &dateTime); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		fixtures = append(fixtures, model.Fixture{
			RoundID:  model.RoundID(roundID),
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
	roundID model.RoundID,
	homeTeam, awayTeam model.TeamKey,
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
			round_id = ?1
		AND home_team = ?2 
		AND away_team = ?3`)
	if err != nil {
		return rollbackTransaction(tx, fmt.Errorf("preparing add result statement: %w", err), "rolling back add result transaction")
	}
	defer stmt.Close()

	res, err := stmt.Exec(roundID, homeTeam, awayTeam, homeScore, awayScore)
	if err != nil {
		return rollbackTransaction(tx, fmt.Errorf("executing add result statement for round %q and teams %s vs %s: %w", roundID, homeTeam, awayTeam, err), "rolling back add result transaction")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return rollbackTransaction(tx, fmt.Errorf("checking affected rows when adding result for round %q and teams %s vs %s: %w", roundID, homeTeam, awayTeam, err), "rolling back add result transaction")
	}
	if rows == 0 {
		return rollbackTransaction(tx, fmt.Errorf("no fixture found for round %q and teams %s vs %s: %w", roundID, homeTeam, awayTeam, model.UnknownFixtureErr), "rolling back add result transaction")
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
		var roundID string
		var homeTeamKey string
		var awayTeamKey string
		var homeGoals, awayGoals int
		var dateTime time.Time
		if err := rows.Scan(
			&roundID,
			&homeTeamKey,
			&awayTeamKey,
			&dateTime,
			&homeGoals,
			&awayGoals); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		results = append(results, model.Result{
			Fixture: model.Fixture{
				RoundID:  model.RoundID(roundID),
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
