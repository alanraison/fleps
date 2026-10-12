package sql

import (
	"errors"
	"testing"

	"github.com/alanraison/fleps/pkg/model"
)

func TestShouldAddPrediction(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	player := "alan.raison@gmail.com"
	round := model.RoundID("R1")

	fr := NewFixtureRepository(db)
	repo := NewPredictionRepository(db)

	fs, err := fr.ListFixtures(round)
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}
	f := fs[0]

	err = repo.AddPredictions(player, round, model.GamePredictions{
		model.Game{
			HomeTeam: f.HomeTeam,
			AwayTeam: f.AwayTeam,
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("AddPrediction returned error: %v", err)
	}

	predictions, err := repo.ListPredictions(round)
	if err != nil {
		t.Fatalf("ListPredictions returned error: %v", err)
	}
	if len(predictions) != 1 {
		t.Fatalf("expected 1 prediction, got %d", len(predictions))
	}
	prediction := predictions["alan.raison@gmail.com"]
	p, ok := prediction[model.Game{
		HomeTeam: f.HomeTeam,
		AwayTeam: f.AwayTeam,
	}]
	if !ok {
		t.Fatalf("prediction for home team '%s'/away team '%s' not found", f.HomeTeam, f.AwayTeam)
	}
	if p.HomeGoals != 2 {
		t.Fatalf("expected home goals 2, got %d", p.HomeGoals)
	}
	if p.AwayGoals != 1 {
		t.Fatalf("expected away goals 1, got %d", p.AwayGoals)
	}
}

func TestShouldFailForUnknownFixture(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultPlayerData(t, db)

	player := "alan.raison@gmail.com"
	round := model.RoundID("R1")

	repo := NewPredictionRepository(db)

	err := repo.AddPredictions(player, round, model.GamePredictions{
		model.Game{
			HomeTeam: "LEE",
			AwayTeam: "MUN",
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
	})
	if err == nil {
		t.Fatalf("expected error when adding prediction for unknown fixture, got nil")
	}
	if !errors.Is(err, model.UnknownFixtureErr) {
		t.Fatalf("expected UnknownFixtureErr, got %v", err)
	}
}

func TestShouldAddMultiplePredictions(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	fr := NewFixtureRepository(db)
	repo := NewPredictionRepository(db)

	fs, err := fr.ListFixtures(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}
	f := fs[0]

	player1 := "alan.raison@gmail.com"
	player2 := "another.player@example.com"
	round := model.RoundID("R1")
	game := model.Game{
		HomeTeam: f.HomeTeam,
		AwayTeam: f.AwayTeam,
	}

	err = repo.AddPredictions(player1, round, model.GamePredictions{
		model.Game{
			HomeTeam: f.HomeTeam,
			AwayTeam: f.AwayTeam,
		}: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("AddPredictions returned error: %v", err)
	}
	err = repo.AddPredictions(player2, round, model.GamePredictions{
		model.Game{
			HomeTeam: f.HomeTeam,
			AwayTeam: f.AwayTeam,
		}: model.Prediction{
			HomeGoals: 3,
			AwayGoals: 2,
		},
	})
	if err != nil {
		t.Fatalf("AddPredictions returned error: %v", err)
	}

	predictions, err := repo.ListPredictions(round)
	if err != nil {
		t.Fatalf("ListPredictions returned error: %v", err)
	}
	if len(predictions) != 2 {
		t.Fatalf("expected 2 predictions, got %d", len(predictions))
	}
	if g, ok := predictions["alan.raison@gmail.com"]; !ok {
		t.Fatalf("expected prediction for alan.raison@gmail.com, got none")
	} else {
		if p, ok := g[game]; !ok {
			t.Fatalf("expected prediction for game %+v, got none", game)
		} else if p.HomeGoals != 2 || p.AwayGoals != 1 {
			t.Fatalf("unexpected prediction for alan.raison@gmail.com: %+v", g)
		}
	}
	if g, ok := predictions[player2]; !ok {
		t.Fatalf("expected prediction for player %s, got none", player2)
	} else {
		if p, ok := g[game]; !ok {
			t.Fatalf("expected prediction for game %+v, got none", game)
		} else if p.HomeGoals != 3 || p.AwayGoals != 2 {
			t.Fatalf("unexpected prediction for player %s: %+v", player2, g)
		}
	}
}

func TestShouldListNoPredictionsWhenNoPredictionsExist(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	repo := NewPredictionRepository(db)

	round := model.RoundID("R1")
	predictions, err := repo.ListPredictions(round)
	if err != nil {
		t.Fatalf("ListPredictions returned error: %v", err)
	}
	if len(predictions) != 0 {
		t.Fatalf("expected 0 predictions, got %d", len(predictions))
	}
}

func TestShouldListRoundPredictions(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	fr := NewFixtureRepository(db)
	repo := NewPredictionRepository(db)

	fs, err := fr.ListFixtures(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}
	f := fs[0]

	player := "alan.raison@gmail.com"
	round := model.RoundID("R1")
	game := model.Game{
		HomeTeam: f.HomeTeam,
		AwayTeam: f.AwayTeam,
	}

	err = repo.AddPredictions(player, round, model.GamePredictions{
		game: model.Prediction{
			HomeGoals: 2,
			AwayGoals: 1,
		},
	})
	if err != nil {
		t.Fatalf("AddPredictions returned error: %v", err)
	}

	predictions, err := repo.ListPredictions(round)
	if err != nil {
		t.Fatalf("ListPredictions returned error: %v", err)
	}
	if g, ok := predictions[player]; !ok {
		t.Fatalf("expected prediction for player %s, got none", player)
	} else {
		if p, ok := g[game]; !ok {
			t.Fatalf("expected prediction for game %+v, got none", game)
		} else if p.HomeGoals != 2 || p.AwayGoals != 1 {
			t.Fatalf("unexpected prediction for player %s: %+v", player, g)
		}
	}
}

func TestShouldUpdateExistingPrediction(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	player := "alan.raison@gmail.com"
	round := model.RoundID("R1")
	fs, err := NewFixtureRepository(db).ListFixtures(round)
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}
	game := model.Game{HomeTeam: fs[0].HomeTeam, AwayTeam: fs[0].AwayTeam}
	repo := NewPredictionRepository(db)
	for _, p := range []model.Prediction{{HomeGoals: 2, AwayGoals: 1}, {HomeGoals: 0, AwayGoals: 3}} {
		if err := repo.AddPredictions(player, round, model.GamePredictions{game: p}); err != nil {
			t.Fatalf("AddPredictions returned error: %v", err)
		}
	}
	predictions, err := repo.ListPredictions(round)
	if err != nil {
		t.Fatalf("ListPredictions returned error: %v", err)
	}
	if got := predictions[player][game]; got != (model.Prediction{HomeGoals: 0, AwayGoals: 3}) {
		t.Errorf("expected updated prediction, got %v", got)
	}
}

func TestShouldGetPredictionsForPlayer(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	player := "alan.raison@gmail.com"
	round := model.RoundID("R1")
	fs, err := NewFixtureRepository(db).ListFixtures(round)
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}
	game := model.Game{HomeTeam: fs[0].HomeTeam, AwayTeam: fs[0].AwayTeam}
	repo := NewPredictionRepository(db)

	none, err := repo.GetPredictionsForPlayer(player, round)
	if err != nil || len(none) != 0 {
		t.Fatalf("expected no predictions, got %v, %v", none, err)
	}
	want := model.Prediction{HomeGoals: 2, AwayGoals: 1}
	if err := repo.AddPredictions(player, round, model.GamePredictions{game: want}); err != nil {
		t.Fatalf("AddPredictions returned error: %v", err)
	}
	got, err := repo.GetPredictionsForPlayer(player, round)
	if err != nil {
		t.Fatalf("GetPredictionsForPlayer returned error: %v", err)
	}
	if len(got) != 1 || got[game] != want {
		t.Errorf("unexpected predictions: %v", got)
	}
}
