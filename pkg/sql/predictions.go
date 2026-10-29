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

// AddPredictions saves predictions for fixtures that have not kicked off and returns the
// predictions that were written. Started fixtures are omitted and any existing predictions
// for them remain unchanged.
//
// player is the email of the player making the predictions and must already have been added.
//
// roundID is the ID of the round for which the predictions are being made.
//
// predictions is a map of games to the predicted outcomes. Existing predictions for the
// same player and fixture are replaced only before kickoff; predictions for started fixtures
// remain unchanged.
func (r *predictionRepository) AddPredictions(player string, roundID model.RoundID, predictions model.GamePredictions) (saved model.GamePredictions, err error) {
	saved = make(model.GamePredictions)
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		if cerr := tx.Commit(); cerr != nil {
			err = fmt.Errorf("committing predictions: %w", cerr)
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
		return nil, fmt.Errorf("preparing check fixture statement: %w", err)
	}
	insertPredictionStmt, err := tx.Prepare(`
		INSERT INTO 
			predictions (
				player, 
				fixture_id, 
				home_goals, 
				away_goals
			)
		SELECT $1, fixtures.id, $3, $4
		FROM fixtures
		WHERE fixtures.id = $2
		AND fixtures.date_time > (clock_timestamp() AT TIME ZONE 'UTC')
		ON CONFLICT (player, fixture_id) DO UPDATE
		SET home_goals = EXCLUDED.home_goals, away_goals = EXCLUDED.away_goals`)
	if err != nil {
		return nil, fmt.Errorf("preparing insert prediction statement: %w", err)
	}
	defer insertPredictionStmt.Close()
	defer checkFixtureStmt.Close()

	for game, pred := range predictions {
		var fixtureID int
		err = checkFixtureStmt.QueryRow(game.HomeTeam, game.AwayTeam, roundID).Scan(&fixtureID)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unknown fixture: %v: %v v %v %w", roundID, game.HomeTeam, game.AwayTeam, model.UnknownFixtureErr)
		}
		if err != nil {
			return nil, fmt.Errorf("finding fixture id: %w", err)
		}
		result, err := insertPredictionStmt.Exec(
			player, fixtureID, pred.HomeGoals, pred.AwayGoals)
		if err != nil {
			return nil, fmt.Errorf("inserting prediction: %w", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("checking inserted prediction rows: %w", err)
		}
		if rowsAffected == 1 {
			saved[game] = pred
		}
	}
	return saved, nil
}

func (r *predictionRepository) ListPredictions(roundID model.RoundID) (model.PlayerPredictions, error) {
	rows, err := r.db.Query(`
		SELECT 
			p.player, 
			f.home_team,
			f.away_team, 
			p.home_goals, 
			p.away_goals
		FROM predictions p
		JOIN fixtures f ON p.fixture_id = f.id
		JOIN teams h ON f.home_team = h.key
		WHERE f.round_id = $1
		ORDER BY f.date_time ASC, h.league ASC, h.full_name ASC, p.player ASC`, roundID)
	if err != nil {
		return nil, fmt.Errorf("listing predictions: %w", err)
	}
	defer rows.Close()

	var predictions model.PlayerPredictions
	for rows.Next() {
		var pred model.Prediction
		var player string
		var homeTeamKey string
		var awayTeamKey string
		if err := rows.Scan(
			&player,
			&homeTeamKey,
			&awayTeamKey,
			&pred.HomeGoals,
			&pred.AwayGoals,
		); err != nil {
			return nil, fmt.Errorf("scanning prediction row: %w", err)
		}
		if homeTeamKey == "" {
			return nil, fmt.Errorf("home team key is empty")
		}
		if awayTeamKey == "" {
			return nil, fmt.Errorf("away team key is empty")
		}
		predictions = append(predictions, model.PlayerPrediction{
			Player: player,
			Game: model.Game{
				HomeTeam: model.TeamKey(homeTeamKey),
				AwayTeam: model.TeamKey(awayTeamKey),
			},
			Prediction: pred,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating prediction rows: %w", err)
	}
	return predictions, nil
}

// GetPredictionsForPlayer returns the predictions made by the player (identified by email) for a
// round. The result is empty if the player has made no predictions.
func (r *predictionRepository) GetPredictionsForPlayer(player string, roundID model.RoundID) (model.GamePredictions, error) {
	rows, err := r.db.Query(`
		SELECT
			f.home_team,
			f.away_team,
			p.home_goals,
			p.away_goals
		FROM predictions p
		JOIN fixtures f ON p.fixture_id = f.id
		WHERE p.player = $1
		AND f.round_id = $2`, player, roundID)
	if err != nil {
		return nil, fmt.Errorf("querying predictions for player %q in round %v: %w", player, roundID, err)
	}
	defer rows.Close()

	predictions := make(model.GamePredictions)
	for rows.Next() {
		var homeTeam, awayTeam string
		var pred model.Prediction
		if err := rows.Scan(&homeTeam, &awayTeam, &pred.HomeGoals, &pred.AwayGoals); err != nil {
			return nil, fmt.Errorf("scanning prediction row: %w", err)
		}
		predictions[model.Game{
			HomeTeam: model.TeamKey(homeTeam),
			AwayTeam: model.TeamKey(awayTeam),
		}] = pred
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating prediction rows: %w", err)
	}
	return predictions, nil
}
