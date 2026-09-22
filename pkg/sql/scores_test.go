package sql

import (
	"testing"

	"github.com/alanraison/fleps/pkg/model"
)

func TestCalculatePlayerRoundScore_AllWrong(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	pr := NewPredictionRepository(db)
	rr := NewResultRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 1,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	if err = rr.AddResult(model.FixtureKey{
		RoundID:  "R1",
		HomeTeam: "LEE",
		AwayTeam: "MUN",
	}, 2, 2); err != nil {
		t.Fatalf("failed to add result: %v", err)
	}
	if err = rr.AddResult(model.FixtureKey{
		RoundID:  "R1",
		HomeTeam: "ARS",
		AwayTeam: "CHE",
	}, 2, 0); err != nil {
		t.Fatalf("failed to add result: %v", err)
	}

	score, err := ss.CalculatePlayerRoundScore("alan.raison@gmail.com", "R1")
	if err != nil {
		t.Fatalf("failed to calculate player round score: %v", err)
	}
	if score != 0 {
		t.Fatalf("expected score to be 0, got %d", score)
	}
}

func TestCalculatePlayerRoundScore_CorrectResult_HomeWin(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	pr := NewPredictionRepository(db)
	rr := NewResultRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 0,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	if err = rr.AddResult(model.FixtureKey{
		RoundID:  "R1",
		HomeTeam: "LEE",
		AwayTeam: "MUN",
	}, 1, 0); err != nil {
		t.Fatalf("failed to add result: %v", err)
	}

	score, err := ss.CalculatePlayerRoundScore("alan.raison@gmail.com", "R1")
	if err != nil {
		t.Fatalf("failed to calculate player round score: %v", err)
	}
	if score != 1 {
		t.Fatalf("expected score to be 1, got %d", score)
	}
}

func TestCalculatePlayerRoundScore_CorrectResult_Draw(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	pr := NewPredictionRepository(db)
	rr := NewResultRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 0,
			AwayGoals: 0,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	if err = rr.AddResult(model.FixtureKey{
		RoundID:  "R1",
		HomeTeam: "LEE",
		AwayTeam: "MUN",
	}, 1, 1); err != nil {
		t.Fatalf("failed to add result: %v", err)
	}

	score, err := ss.CalculatePlayerRoundScore("alan.raison@gmail.com", "R1")
	if err != nil {
		t.Fatalf("failed to calculate player round score: %v", err)
	}
	if score != 1 {
		t.Fatalf("expected score to be 1, got %d", score)
	}
}

func TestCalculatePlayerRoundScore_CorrectResult_AwayWin(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	pr := NewPredictionRepository(db)
	rr := NewResultRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 0,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	if err = rr.AddResult(model.FixtureKey{
		RoundID:  "R1",
		HomeTeam: "LEE",
		AwayTeam: "MUN",
	}, 0, 2); err != nil {
		t.Fatalf("failed to add result: %v", err)
	}

	score, err := ss.CalculatePlayerRoundScore("alan.raison@gmail.com", "R1")
	if err != nil {
		t.Fatalf("failed to calculate player round score: %v", err)
	}
	if score != 1 {
		t.Fatalf("expected score to be 1, got %d", score)
	}
}

func TestCalculatePlayerRoundScore_CorrectScore(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	pr := NewPredictionRepository(db)
	rr := NewResultRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	if err = rr.AddResult(model.FixtureKey{
		RoundID:  "R1",
		HomeTeam: "LEE",
		AwayTeam: "MUN",
	}, 2, 1); err != nil {
		t.Fatalf("failed to add result: %v", err)
	}

	score, err := ss.CalculatePlayerRoundScore("alan.raison@gmail.com", "R1")
	if err != nil {
		t.Fatalf("failed to calculate player round score: %v", err)
	}
	if score != 3 {
		t.Fatalf("expected score to be 3, got %d", score)
	}
}

func TestCalculatePlayerRoundScore_MultiplePoints(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)
	setupDefaultResultData(t, db)

	pr := NewPredictionRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 1,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	score, err := ss.CalculatePlayerRoundScore("alan.raison@gmail.com", "R1")
	if err != nil {
		t.Fatalf("failed to calculate player round score: %v", err)
	}
	if score != 4 {
		t.Fatalf("expected score to be 4, got %d", score)
	}
}

func TestCalculateRoundScores(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)
	setupDefaultResultData(t, db)

	pr := NewPredictionRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 1,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}
	err = pr.AddPredictions("another.player@example.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 0,
			AwayGoals: 0,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions for another player: %v", err)
	}

	scores, err := ss.CalculateRoundScores("R1")
	if err != nil {
		t.Fatalf("failed to calculate round scores: %v", err)
	}
	expectedScores := map[string]int{
		"alan.raison@gmail.com":      4,
		"another.player@example.com": 6,
	}
	for player, expectedScore := range expectedScores {
		score, ok := scores[player]
		if !ok {
			t.Fatalf("expected score for player %s, but not found", player)
		}
		if score != expectedScore {
			t.Fatalf("expected score for player %s to be %d, got %d", player, expectedScore, score)
		}
	}
}

func TestCalculateRoundScoresWithMissingPredictions(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)
	setupDefaultResultData(t, db)

	pr := NewPredictionRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 1,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}
	err = pr.AddPredictions("another.player@example.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 0,
			AwayGoals: 0,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions for another player: %v", err)
	}

	scores, err := ss.CalculateRoundScores("R1")
	if err != nil {
		t.Fatalf("failed to calculate round scores: %v", err)
	}
	expectedScores := map[string]int{
		"alan.raison@gmail.com": 4,
		"another.player@example.com": 3,
	}
	for player, expectedScore := range expectedScores {
		score, ok := scores[player]
		if !ok {
			t.Fatalf("expected score for player %s, but not found", player)
		}
		if score != expectedScore {
			t.Fatalf("expected score for player %s to be %d, got %d", player, expectedScore, score)
		}
	}
}

func TestCalculateSeasonScores_FirstRound_OnePlayer(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)
	setupDefaultResultData(t, db)

	pr := NewPredictionRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 1,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	scores, err := ss.CalculateSeasonScores("S1")
	if err != nil {
		t.Fatalf("failed to calculate season scores: %v", err)
	}
	expectedScores := map[string]int{
		"alan.raison@gmail.com": 4,
	}
	for player, expectedScore := range expectedScores {
		if scores[player] != expectedScore {
			t.Fatalf("expected score for player %s to be %d, got %d", player, expectedScore, scores[player])
		}
	}
}

func TestCalculateSeasonScores_FirstRound_TwoPlayers(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)
	setupDefaultResultData(t, db)

	pr := NewPredictionRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	err = pr.AddPredictions("another.player@example.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 0,
			AwayGoals: 0,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions for another player: %v", err)
	}
	err = pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 1,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	scores, err := ss.CalculateSeasonScores("S1")
	if err != nil {
		t.Fatalf("failed to calculate season scores: %v", err)
	}
	expectedScores := map[string]int{
		"alan.raison@gmail.com":      4,
		"another.player@example.com": 6,
	}
	for player, expectedScore := range expectedScores {
		if scores[player] != expectedScore {
			t.Fatalf("expected score for player %s to be %d, got %d", player, expectedScore, scores[player])
		}
	}
}

func TestCalculateSeasonScores_TwoRounds_TwoPlayers(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)
	setupDefaultResultData(t, db)

	rr := NewResultRepository(db)
	pr := NewPredictionRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 1,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	err = pr.AddPredictions("another.player@example.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 0,
			AwayGoals: 2,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions for another player: %v", err)
	}

	err = pr.AddPredictions("alan.raison@gmail.com", "R2", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 3,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	err = pr.AddPredictions("another.player@example.com", "R2", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 1,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions for another player: %v", err)
	}

	rr.AddResult(model.FixtureKey{
		RoundID:  "R2",
		HomeTeam: "LEE",
		AwayTeam: "CHE",
	}, 1, 1)

	scores, err := ss.CalculateSeasonScores("S1")
	if err != nil {
		t.Fatalf("failed to calculate season scores: %v", err)
	}
	expectedScores := map[string]int{
		"alan.raison@gmail.com":      4,
		"another.player@example.com": 6,
	}
	for player, expectedScore := range expectedScores {
		if scores[player] != expectedScore {
			t.Fatalf("expected score for player %s to be %d, got %d", player, expectedScore, scores[player])
		}
	}
}

func TestCalculateCurrentSeasonScores_TwoRounds_TwoPlayers(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)
	setupDefaultResultData(t, db)

	pr := NewPredictionRepository(db)
	ss := NewDBScoreService(db)

	err := pr.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 1,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	err = pr.AddPredictions("another.player@example.com", "R1", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
		model.Game{
			HomeTeam: "ARS",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 0,
			AwayGoals: 2,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions for another player: %v", err)
	}

	err = pr.AddPredictions("alan.raison@gmail.com", "R2", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 3,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions: %v", err)
	}

	err = pr.AddPredictions("another.player@example.com", "R2", model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "CHE",
		}: model.Prediction{
			HomeGoals: 1,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to add predictions for another player: %v", err)
	}

	rr := NewResultRepository(db)
	rr.AddResult(model.FixtureKey{
		RoundID:  "R2",
		HomeTeam: "LEE",
		AwayTeam: "CHE",
	}, 1, 1)

	scores, err := ss.CalculateCurrentSeasonScores()
	if err != nil {
		t.Fatalf("failed to calculate current season scores: %v", err)
	}
	expectedScores := map[string]int{
		"alan.raison@gmail.com":      4,
		"another.player@example.com": 6,
	}
	for player, expectedScore := range expectedScores {
		if scores[player] != expectedScore {
			t.Fatalf("expected score for player %s to be %d, got %d", player, expectedScore, scores[player])
		}
	}
}
