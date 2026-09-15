package sql

import (
	"database/sql"
	"fmt"

	"github.com/alanraison/predictions/pkg/model"
)

type DBScoreService struct {
	db *sql.DB
}

func NewDBScoreService(db *sql.DB) *DBScoreService {
	return &DBScoreService{db: db}
}

func (s *DBScoreService) CalculatePlayerRoundScore(player string, roundID model.RoundID) (int, error) {
	var score int
	err := s.db.QueryRow(`
		SELECT
			SUM(CASE
				WHEN results.home_goals = predictions.home_goals
				AND results.away_goals = predictions.away_goals
				THEN 3
				WHEN SIGN(results.home_goals - results.away_goals) = SIGN(predictions.home_goals - predictions.away_goals)
				THEN 1
			ELSE 0
			END) as score
		FROM
			results
		JOIN
			fixtures ON results.fixture_id = fixtures.id
		JOIN 
			predictions ON predictions.fixture_id = fixtures.id
		WHERE 
			predictions.player = $1 
		AND fixtures.round_id = $2`,
		player, roundID).Scan(&score)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate player round score: %w", err)
	}
	return score, nil
}

func (s *DBScoreService) CalculateRoundScores(roundID model.RoundID) (map[string]int, error) {
	rows, err := s.db.Query(`
		SELECT
			predictions.player,
			SUM(CASE
				WHEN results.home_goals = predictions.home_goals
				AND results.away_goals = predictions.away_goals
				THEN 3
				WHEN SIGN(results.home_goals - results.away_goals) = SIGN(predictions.home_goals - predictions.away_goals)
				THEN 1
			ELSE 0
			END) as score
		FROM
			results
		JOIN 
		  fixtures ON fixtures.id = results.fixture_id
		JOIN 
			predictions ON predictions.fixture_id = fixtures.id
		WHERE 
			fixtures.round_id = $1
		GROUP BY predictions.player
		ORDER BY score DESC`, roundID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate round scores: %w", err)
	}
	
	defer rows.Close()
	var scores = make(map[string]int)
	for rows.Next() {
		var player string
		var score int
		if err := rows.Scan(&player, &score); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		scores[player] = score
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return scores, nil
}

func (s *DBScoreService) CalculateSeasonScores(seasonID model.SeasonID) (map[string]int, error) {
	rows, err := s.db.Query(`
		SELECT
			predictions.player,
			SUM(CASE
				WHEN results.home_goals = predictions.home_goals
				AND results.away_goals = predictions.away_goals
				THEN 3
				WHEN SIGN(results.home_goals - results.away_goals) = SIGN(predictions.home_goals - predictions.away_goals)
				THEN 1
			ELSE 0
			END) as score
		FROM
			results
		JOIN fixtures ON fixtures.id = results.fixture_id
		JOIN 
			predictions ON predictions.fixture_id = fixtures.id
		JOIN
			rounds ON fixtures.round_id = rounds.id
		WHERE 
			rounds.season_id = $1
		GROUP BY predictions.player
		ORDER BY score DESC`, seasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate season scores: %w", err)
	}
	defer rows.Close()
	var scores = make(map[string]int)
	for rows.Next() {
		var player string
		var score int
		if err := rows.Scan(&player, &score); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		scores[player] = score
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return scores, nil
}
