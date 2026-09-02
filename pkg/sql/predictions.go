package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alanraison/predictions/pkg/model"
)

type playerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *playerRepository {
	return &playerRepository{
		db: db,
	}
}

func (r *playerRepository) AddPlayer(email, name string) error {
	_, err := r.db.Exec("INSERT INTO players (email, name) VALUES (?, ?)", email, name)
	return err
}

func (r *playerRepository) ListPlayers() ([]model.Player, error) {
	rows, err := r.db.Query("SELECT email, name FROM players")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []model.Player
	for rows.Next() {
		var player model.Player
		if err := rows.Scan(&player.Email, &player.Name); err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return players, nil
}

type predictionRepository struct {
	db *sql.DB
}

func NewPredictionRepository(db *sql.DB) *predictionRepository {
	return &predictionRepository{
		db: db,
	}
}

// AddPrediction adds a prediction to the repository.
func (r *predictionRepository) AddPrediction(player string, homeTeam, awayTeam model.TeamKey, date time.Time, homeScore, awayScore int) error {
	var fixtureID int
	var dbDate string
	err := r.db.QueryRow(`SELECT ? AS date_time`, date.Local()).Scan(&dbDate)
	if err != nil {
		return fmt.Errorf("querying date: %w", err)
	}
	err = r.db.QueryRow(`
		SELECT 
			id 
		FROM 
			fixtures 
		WHERE 
			home_team = ? 
		AND away_team = ? 
		AND datetime(date_time) = datetime(?)`,
		homeTeam, awayTeam, date.Local(),
	).Scan(&fixtureID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("unknown fixture: %w", model.UnknownFixtureErr)
	}
	if err != nil {
		return fmt.Errorf("finding fixture id: %w", err)
	}
	_, err = r.db.Exec(`
		INSERT INTO 
			predictions (
				player, 
				fixture_id, 
				home_goals, 
				away_goals
			)
		VALUES (?, ?, ?, ?)`,
		player, fixtureID, homeScore, awayScore)
	if err != nil {
		return fmt.Errorf("inserting prediction: %w", err)
	}
	return nil
}

func (r *predictionRepository) ListPredictions(from, to time.Time) ([]model.Prediction, error) {
	rows, err := r.db.Query(`
		SELECT 
			p.player, 
			f.home_team,
			f.away_team, 
			f.date_time, 
			p.home_goals, 
			p.away_goals
		FROM predictions p
		JOIN fixtures f ON p.fixture_id = f.id
		WHERE f.date_time BETWEEN ? AND ?`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []model.Prediction
	for rows.Next() {
		var pred model.Prediction
		var homeTeamKey string
		var awayTeamKey string
		if err := rows.Scan(
			&pred.Player,
			&homeTeamKey,
			&awayTeamKey,
			&pred.Date,
			&pred.HomeGoals,
			&pred.AwayGoals,
		); err != nil {
			return nil, err
		}
		if homeTeamKey == "" {
			return nil, fmt.Errorf("home team key is empty")
		}
		if awayTeamKey == "" {
			return nil, fmt.Errorf("away team key is empty")
		}
		pred.HomeTeam = model.TeamKey(homeTeamKey)
		pred.AwayTeam = model.TeamKey(awayTeamKey)
		predictions = append(predictions, pred)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return predictions, nil
}
