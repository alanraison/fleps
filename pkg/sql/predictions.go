package sql

import (
	"database/sql"
	"fmt"

	"github.com/alanraison/fleps/pkg/model"
)

type predictionRepository struct {
	db *sql.DB
}

func NewPredictionRepository(db *sql.DB) *predictionRepository {
	return &predictionRepository{
		db: db,
	}
}

// AddPredictions adds predictions to the repository.
//
// player is the email of the player making the predictions and must already have been added.
//
// roundID is the ID of the round for which the predictions are being made.
//
// predictions is a map of games to the predicted outcomes.
func (r *predictionRepository) AddPredictions(player string, roundID model.RoundID, predictions model.GamePredictions) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	checkFixtureStmt, err := tx.Prepare(`
		SELECT 
			id 
		FROM 
			fixtures 
		WHERE 
			home_team = $1
		AND away_team = $2
		AND round_id = $3`)
	if err != nil {
		return fmt.Errorf("preparing check fixture statement: %w", err)
	}
	insertPredictionStmt, err := tx.Prepare(`
		INSERT INTO 
			predictions (
				player, 
				fixture_id, 
				home_goals, 
				away_goals
			)
		VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return fmt.Errorf("preparing insert prediction statement: %w", err)
	}
	defer insertPredictionStmt.Close()
	defer checkFixtureStmt.Close()

	for game, pred := range predictions {
		var fixtureID int
		err := checkFixtureStmt.QueryRow(game.HomeTeam, game.AwayTeam, roundID).Scan(&fixtureID)
		if err == sql.ErrNoRows {
			return fmt.Errorf("unknown fixture: %v: %v v %v %w", roundID, game.HomeTeam, game.AwayTeam, model.UnknownFixtureErr)
		}
		if err != nil {
			return fmt.Errorf("finding fixture id: %w", err)
		}
		_, err = insertPredictionStmt.Exec(
			player, fixtureID, pred.HomeGoals, pred.AwayGoals)
		if err != nil {
			return fmt.Errorf("inserting prediction: %w", err)
		}
	}
	return nil
}

func (r *predictionRepository) ListPredictions(roundID model.RoundID) (model.PlayerPredictions, error) {
	rows, err := r.db.Query(`
		SELECT 
			p.player, 
			f.round_id,
			f.home_team,
			f.away_team, 
			p.home_goals, 
			p.away_goals
		FROM predictions p
		JOIN fixtures f ON p.fixture_id = f.id
		WHERE f.round_id = $1`, roundID)
	if err != nil {
		return nil, fmt.Errorf("listing predictions: %w", err)
	}
	defer rows.Close()

	var predictions model.PlayerPredictions = make(model.PlayerPredictions)
	for rows.Next() {
		var pred model.Prediction
		var roundID string
		var player string
		var homeTeamKey string
		var awayTeamKey string
		if err := rows.Scan(
			&player,
			&roundID,
			&homeTeamKey,
			&awayTeamKey,
			&pred.HomeGoals,
			&pred.AwayGoals,
		); err != nil {
			return nil, fmt.Errorf("scanning prediction row: %w", err)
		}
		if _, ok := predictions[player]; !ok {
			predictions[player] = make(model.GamePredictions)
		}
		if homeTeamKey == "" {
			return nil, fmt.Errorf("home team key is empty")
		}
		if awayTeamKey == "" {
			return nil, fmt.Errorf("away team key is empty")
		}
		predictions[player][model.Game{
			HomeTeam: model.TeamKey(homeTeamKey),
			AwayTeam: model.TeamKey(awayTeamKey),
		}] = pred
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating prediction rows: %w", err)
	}
	return predictions, nil
}
