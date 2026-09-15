package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alanraison/fleps/pkg/model"
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
	listResultsQuery = `
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
			round_id = $1
	`
)

func rollbackTransaction(tx *sql.Tx, cause error, context string) error {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return fmt.Errorf("%s: %w", context, errors.Join(cause, fmt.Errorf("rolling back transaction: %w", rollbackErr)))
	}

	return cause
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
	err := r.db.QueryRow("SELECT id FROM rounds WHERE id = $1", roundID).Scan(&foundRoundID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("finding round %q: %w", roundID, model.UnknownRoundErr)
	}
	if err != nil {
		return fmt.Errorf("finding round %q: %w", roundID, err)
	}

	return nil
}

func (r *fixtureRepository) GetLatestRoundWithNoResults() (model.RoundID, error) {
	var latestRoundID model.RoundID
	if err := r.db.QueryRow(`
			SELECT
				rounds.id
			FROM
				rounds 
			LEFT JOIN
				fixtures ON rounds.id = fixtures.round_id
			LEFT JOIN
				results ON fixtures.id = results.fixture_id
			GROUP BY
				rounds.id
			HAVING
				COUNT(results.fixture_id) = 0
			ORDER BY
				rounds.id ASC
			LIMIT 1
			`).Scan(&latestRoundID); err != nil {
		return "", fmt.Errorf("reading latest round id: %w", err)
	}

	return latestRoundID, nil
}

func (r *fixtureRepository) AddRound(id model.RoundID, season model.SeasonID) error {
	if _, err := r.db.Exec("INSERT INTO rounds (id, season_id) VALUES ($1, $2)", id, season); err != nil {
		return fmt.Errorf("inserting round record: %w", err)
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
		) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return rollbackTransaction(tx, fmt.Errorf("preparing add fixtures statement: %w", err), "rolling back add fixtures transaction")
	}
	defer stmt.Close()

	for _, fixture := range fixtures {
		res, err := stmt.Exec(fixture.RoundID, fixture.HomeTeam, fixture.AwayTeam, fixture.Date.UTC())
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

func (r *fixtureRepository) ListFixtures(roundId model.RoundID) ([]model.Fixture, error) {
	var rows *sql.Rows
	var err error
	if roundId == "" {
		roundId, err = r.GetLatestRoundWithNoResults()
		if err != nil {
			return nil, fmt.Errorf("getting latest round id: %w", err)
		}
	}
	rows, err = r.db.Query(`
		SELECT
			round_id,
			home_team,
			away_team,
			date_time
		FROM fixtures
		WHERE round_id = $1
		ORDER BY date_time ASC
	`, roundId)
	if err != nil {
		return nil, fmt.Errorf("querying fixtures for round %q: %w", roundId, err)
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
			FixtureKey: model.FixtureKey{
				RoundID:  model.RoundID(roundID),
				HomeTeam: model.TeamKey(homeTeamKey),
				AwayTeam: model.TeamKey(awayTeamKey),
			},
			Date: dateTime.Local(),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}
	return fixtures, nil
}
