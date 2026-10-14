package render

import (
	"reflect"
	"testing"

	"github.com/alanraison/fleps/pkg/model"
)

type fakeFixtures struct{ model.FixtureRepository }

func (fakeFixtures) ListFixtures(model.RoundID) ([]model.Fixture, error) {
	return []model.Fixture{{FixtureKey: model.FixtureKey{HomeTeam: "ars", AwayTeam: "che"}}}, nil
}

type fakeResults struct{ model.ResultRepository }

func (fakeResults) ListResults(model.RoundID) ([]model.Result, error) {
	return []model.Result{{FixtureKey: model.FixtureKey{HomeTeam: "ars", AwayTeam: "che"}, HomeGoals: 2, AwayGoals: 1}}, nil
}

type fakePredictions struct{ model.PredictionRepository }

func (fakePredictions) ListPredictions(model.RoundID) (model.PlayerPredictions, error) {
	g := model.Game{HomeTeam: "ars", AwayTeam: "che"}
	return model.PlayerPredictions{
		{Player: "a@x", Game: g, Prediction: model.Prediction{HomeGoals: 1, AwayGoals: 0}},
		{Player: "b@x", Game: g, Prediction: model.Prediction{HomeGoals: 2, AwayGoals: 1}},
	}, nil
}

type fakePlayers struct{ model.PlayerRepository }

func (fakePlayers) ListActivePlayers() ([]model.Player, error) {
	return []model.Player{{Email: "a@x", Name: "Ann"}, {Email: "b@x", Name: "Bob"}, {Email: "c@x", Name: "Cy"}}, nil
}

type fakeTeams struct{ model.TeamRepository }

func (fakeTeams) FindTeamByKey(k model.TeamKey) (*model.Team, error) {
	return &model.Team{Key: k, ShortName: string(k)}, nil
}

type fakeScores struct{ model.ScoreService }

func (fakeScores) CalculateRoundScores(model.RoundID) (map[string]int, error) {
	return map[string]int{"a@x": 1, "b@x": 3}, nil
}
func (fakeScores) CalculateCurrentSeasonScores() (map[string]int, error) {
	return map[string]int{"a@x": 5, "b@x": 9, "c@x": 5}, nil
}

func TestBuildTable(t *testing.T) {
	s := Sources{fakeFixtures{}, fakeResults{}, fakePredictions{}, fakePlayers{}, fakeTeams{}, fakeScores{}}
	table, err := s.BuildTable("r1")
	if err != nil {
		t.Fatalf("BuildTable: %v", err)
	}
	if want := []string{"Bob", "Ann", "Cy"}; !reflect.DeepEqual(table.Players, want) {
		t.Errorf("players = %v, want %v", table.Players, want)
	}
	if want := []int{9, 5, 5}; !reflect.DeepEqual(table.SeasonScores, want) {
		t.Errorf("season scores = %v, want %v", table.SeasonScores, want)
	}
	if want := []int{3, 1, 0}; !reflect.DeepEqual(table.RoundScores, want) {
		t.Errorf("round scores = %v, want %v", table.RoundScores, want)
	}
	want := Row{
		Label:  "ars v che",
		Result: "2-1",
		Predictions: []Cell{
			{Text: "2-1", Points: PointsExact},
			{Text: "1-0", Points: PointsOutcome},
			{Text: "-"},
		},
	}
	if !reflect.DeepEqual(table.Rows[0], want) {
		t.Errorf("row = %+v, want %+v", table.Rows[0], want)
	}
}

func TestPredictionPoints(t *testing.T) {
	r := model.Result{HomeGoals: 1, AwayGoals: 1}
	for _, tc := range []struct {
		p    model.Prediction
		want int
	}{
		{model.Prediction{HomeGoals: 1, AwayGoals: 1}, PointsExact},
		{model.Prediction{HomeGoals: 0, AwayGoals: 0}, PointsOutcome},
		{model.Prediction{HomeGoals: 1, AwayGoals: 0}, PointsNone},
	} {
		if got := PredictionPoints(tc.p, r); got != tc.want {
			t.Errorf("PredictionPoints(%+v) = %d, want %d", tc.p, got, tc.want)
		}
	}
}
