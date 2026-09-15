package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alanraison/fleps/pkg/model"
)

type resultRepository struct {
	db *sql.DB
}

func NewResultRepository(db *sql.DB) *resultRepository {
	return &resultRepository{
		db: db,
	}
}

func (r *resultRepository) GetLatestRoundWithResults() (model.RoundID, error) {
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
				COUNT(results.fixture_id) > 0
			ORDER BY
				rounds.id DESC
			LIMIT 1
			`).Scan(&latestRoundID); err != nil {
		return "", fmt.Errorf("reading latest round with results id: %w", err)
	}

	return latestRoundID, nil
}

func (r *resultRepository) AddResult(
	fixtureKey model.FixtureKey,
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
			$4,
			$5
		FROM fixtures
		WHERE 
			round_id = $1
		AND home_team = $2 
		AND away_team = $3`)
	if err != nil {
		return rollbackTransaction(tx, fmt.Errorf("preparing add result statement: %w", err),
			"rolling back add result transaction")
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		fixtureKey.RoundID,
		fixtureKey.HomeTeam,
		fixtureKey.AwayTeam,
		homeScore,
		awayScore,
	)
	if err != nil {
		return rollbackTransaction(
			tx,
			fmt.Errorf(
				"executing add result statement for round %q and teams %s vs %s: %w",
				fixtureKey.RoundID,
				fixtureKey.HomeTeam,
				fixtureKey.AwayTeam,
				err,
			),
			"rolling back add result transaction",
		)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return rollbackTransaction(
			tx,
			fmt.Errorf(
				"checking affected rows when adding result for round %q and teams %s vs %s: %w",
				fixtureKey.RoundID,
				fixtureKey.HomeTeam,
				fixtureKey.AwayTeam,
				err,
			),
			"rolling back add result transaction",
		)
	}
	if rows == 0 {
		return rollbackTransaction(
			tx,
			fmt.Errorf(
				"no fixture found for round %q and teams %s vs %s: %w",
				fixtureKey.RoundID,
				fixtureKey.HomeTeam,
				fixtureKey.AwayTeam,
				model.UnknownFixtureErr,
			),
			"rolling back add result transaction",
		)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (r *resultRepository) ListResults(round model.RoundID) ([]model.Result, error) {
	var rows *sql.Rows
	var err error
	if err = r.db.QueryRow(
		"SELECT true FROM rounds where id = $1", round,
	).Scan(&[]byte{}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("unknown round: %w", model.UnknownRoundErr)
		}
		return nil, fmt.Errorf("checking round existence: %w", err)
	}
	rows, err = r.db.Query(listResultsQuery, round)
	if err != nil {
		return nil, fmt.Errorf("reading results: %w", err)
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
			FixtureKey: model.FixtureKey{
				RoundID:  model.RoundID(roundID),
				HomeTeam: model.TeamKey(homeTeamKey),
				AwayTeam: model.TeamKey(awayTeamKey),
			},
			HomeGoals: homeGoals,
			AwayGoals: awayGoals,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}
	return results, nil
}
