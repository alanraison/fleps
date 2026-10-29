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

	_, err = repo.AddPredictions(player, round, model.GamePredictions{
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
	result := predictions[0]
	if result.Player != player {
		t.Fatalf("expected prediction for player %q, got %q", player, result.Player)
	}
	if result.Game != (model.Game{
		HomeTeam: f.HomeTeam,
		AwayTeam: f.AwayTeam,
	}) {
		t.Fatalf("unexpected game: %+v", result.Game)
	}
	p := result.Prediction
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

	_, err := repo.AddPredictions(player, round, model.GamePredictions{
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

	_, err = repo.AddPredictions(player1, round, model.GamePredictions{
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
	_, err = repo.AddPredictions(player2, round, model.GamePredictions{
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
	gotByPlayer := make(map[string]model.Prediction)
	for _, result := range predictions {
		gotByPlayer[result.Player] = result.Prediction
		if result.Game != game {
			t.Fatalf("unexpected game for player %s: %+v", result.Player, result.Game)
		}
	}
	if p := gotByPlayer["alan.raison@gmail.com"]; p.HomeGoals != 2 || p.AwayGoals != 1 {
		t.Fatalf("unexpected prediction for alan.raison@gmail.com: %+v", p)
	}
	if p := gotByPlayer[player2]; p.HomeGoals != 3 || p.AwayGoals != 2 {
		t.Fatalf("unexpected prediction for player %s: %+v", player2, p)
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

	_, err = repo.AddPredictions(player, round, model.GamePredictions{
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
	if len(predictions) != 1 {
		t.Fatalf("expected 1 prediction, got %d", len(predictions))
	}
	if predictions[0].Player != player || predictions[0].Game != game {
		t.Fatalf("unexpected prediction result: %+v", predictions[0])
	}
	if p := predictions[0].Prediction; p.HomeGoals != 2 || p.AwayGoals != 1 {
		t.Fatalf("unexpected prediction for player %s: %+v", player, p)
	}
}

func TestListPredictionsOrdersByKickoffLeagueAndHomeTeamName(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultPlayerData(t, db)
	if _, err := db.Exec(`
		INSERT INTO leagues (rank, name) VALUES (2, 'Championship');
		INSERT INTO teams (key, full_name, short_name, league)
		VALUES ('NEW', 'Newcastle United', 'Newcastle', 2);
		INSERT INTO fixtures (round_id, home_team, away_team, date_time) VALUES
			('R1', 'ARS', 'CHE', (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + INTERVAL '3 days'),
			('R1', 'NEW', 'ARS', (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + INTERVAL '2 days'),
			('R1', 'LEE', 'MUN', (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + INTERVAL '2 days'),
			('R1', 'CHE', 'ARS', (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + INTERVAL '2 days'),
			('R1', 'LEE', 'CHE', (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + INTERVAL '1 day');
	`); err != nil {
		t.Fatalf("failed to insert ordered prediction fixtures: %v", err)
	}

	games := model.GamePredictions{
		{HomeTeam: "ARS", AwayTeam: "CHE"}: {HomeGoals: 1, AwayGoals: 0},
		{HomeTeam: "NEW", AwayTeam: "ARS"}: {HomeGoals: 1, AwayGoals: 0},
		{HomeTeam: "LEE", AwayTeam: "MUN"}: {HomeGoals: 1, AwayGoals: 0},
		{HomeTeam: "CHE", AwayTeam: "ARS"}: {HomeGoals: 1, AwayGoals: 0},
		{HomeTeam: "LEE", AwayTeam: "CHE"}: {HomeGoals: 1, AwayGoals: 0},
	}
	repo := NewPredictionRepository(db)
	player := "alan.raison@gmail.com"
	if _, err := repo.AddPredictions(player, "R1", games); err != nil {
		t.Fatalf("AddPredictions returned error: %v", err)
	}

	predictions, err := repo.ListPredictions("R1")
	if err != nil {
		t.Fatalf("ListPredictions returned error: %v", err)
	}
	wantHomeTeams := []model.TeamKey{"LEE", "CHE", "LEE", "NEW", "ARS"}
	if len(predictions) != len(wantHomeTeams) {
		t.Fatalf("expected %d predictions, got %d", len(wantHomeTeams), len(predictions))
	}
	for i, want := range wantHomeTeams {
		if got := predictions[i].Game.HomeTeam; got != want {
			t.Errorf("predictions[%d].Game.HomeTeam = %q, want %q", i, got, want)
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
		if _, err := repo.AddPredictions(player, round, model.GamePredictions{game: p}); err != nil {
			t.Fatalf("AddPredictions returned error: %v", err)
		}
	}
	predictions, err := repo.ListPredictions(round)
	if err != nil {
		t.Fatalf("ListPredictions returned error: %v", err)
	}
	if got := predictions[0].Prediction; got != (model.Prediction{HomeGoals: 0, AwayGoals: 3}) {
		t.Errorf("expected updated prediction, got %v", got)
	}
}

func TestAddPredictionsSkipsStartedFixtureWithoutExistingPrediction(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	game := model.Game{HomeTeam: "LEE", AwayTeam: "MUN"}
	if _, err := db.Exec(`
		UPDATE fixtures
		SET date_time = (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') - INTERVAL '1 minute'
		WHERE round_id = 'R1' AND home_team = 'LEE' AND away_team = 'MUN'
	`); err != nil {
		t.Fatalf("marking fixture as started: %v", err)
	}

	repo := NewPredictionRepository(db)
	saved, err := repo.AddPredictions("alan.raison@gmail.com", "R1", model.GamePredictions{
		game: {HomeGoals: 2, AwayGoals: 1},
	})
	if err != nil {
		t.Fatalf("AddPredictions returned error: %v", err)
	}
	if len(saved) != 0 {
		t.Fatalf("expected no saved predictions, got %v", saved)
	}
	got, err := repo.GetPredictionsForPlayer("alan.raison@gmail.com", "R1")
	if err != nil {
		t.Fatalf("GetPredictionsForPlayer returned error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected no stored predictions, got %v", got)
	}
}

func TestAddPredictionsKeepsExistingPredictionForStartedFixture(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	player := "alan.raison@gmail.com"
	game := model.Game{HomeTeam: "LEE", AwayTeam: "MUN"}
	repo := NewPredictionRepository(db)
	existing := model.Prediction{HomeGoals: 2, AwayGoals: 1}
	if _, err := repo.AddPredictions(player, "R1", model.GamePredictions{game: existing}); err != nil {
		t.Fatalf("adding initial prediction: %v", err)
	}
	if _, err := db.Exec(`
		UPDATE fixtures
		SET date_time = (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') - INTERVAL '1 minute'
		WHERE round_id = 'R1' AND home_team = 'LEE' AND away_team = 'MUN'
	`); err != nil {
		t.Fatalf("marking fixture as started: %v", err)
	}

	saved, err := repo.AddPredictions(player, "R1", model.GamePredictions{
		game: {HomeGoals: 0, AwayGoals: 4},
	})
	if err != nil {
		t.Fatalf("AddPredictions returned error: %v", err)
	}
	if len(saved) != 0 {
		t.Fatalf("expected no saved predictions for started fixture, got %v", saved)
	}
	got, err := repo.GetPredictionsForPlayer(player, "R1")
	if err != nil {
		t.Fatalf("GetPredictionsForPlayer returned error: %v", err)
	}
	if got[game] != existing {
		t.Fatalf("existing prediction changed: got %v, want %v", got[game], existing)
	}
}

func TestAddPredictionsSavesUpcomingFixturesWhenAnotherHasStarted(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)
	setupDefaultPlayerData(t, db)

	started := model.Game{HomeTeam: "LEE", AwayTeam: "MUN"}
	upcoming := model.Game{HomeTeam: "ARS", AwayTeam: "CHE"}
	if _, err := db.Exec(`
		UPDATE fixtures
		SET date_time = (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') - INTERVAL '1 minute'
		WHERE round_id = 'R1' AND home_team = 'LEE' AND away_team = 'MUN'
	`); err != nil {
		t.Fatalf("marking fixture as started: %v", err)
	}

	repo := NewPredictionRepository(db)
	predictions := model.GamePredictions{
		started:  {HomeGoals: 2, AwayGoals: 1},
		upcoming: {HomeGoals: 3, AwayGoals: 0},
	}
	saved, err := repo.AddPredictions("alan.raison@gmail.com", "R1", predictions)
	if err != nil {
		t.Fatalf("AddPredictions returned error: %v", err)
	}
	if len(saved) != 1 || saved[upcoming] != predictions[upcoming] {
		t.Fatalf("unexpected saved predictions: %v", saved)
	}
	got, err := repo.GetPredictionsForPlayer("alan.raison@gmail.com", "R1")
	if err != nil {
		t.Fatalf("GetPredictionsForPlayer returned error: %v", err)
	}
	if len(got) != 1 || got[upcoming] != predictions[upcoming] {
		t.Fatalf("unexpected stored predictions: %v", got)
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
	if _, err := repo.AddPredictions(player, round, model.GamePredictions{game: want}); err != nil {
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
