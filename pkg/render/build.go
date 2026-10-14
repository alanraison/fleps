package render

import (
	"fmt"
	"sort"

	"github.com/alanraison/fleps/pkg/model"
)

// Sources are the repositories the results table is built from.
type Sources struct {
	Fixtures    model.FixtureRepository
	Results     model.ResultRepository
	Predictions model.PredictionRepository
	Players     model.PlayerRepository
	Teams       model.TeamRepository
	Scores      model.ScoreService
}

// PredictionPoints returns the points a prediction earns against a result.
func PredictionPoints(p model.Prediction, r model.Result) int {
	switch {
	case p.HomeGoals == r.HomeGoals && p.AwayGoals == r.AwayGoals:
		return PointsExact
	case sign(p.HomeGoals-p.AwayGoals) == sign(r.HomeGoals-r.AwayGoals):
		return PointsOutcome
	}
	return PointsNone
}

func sign(n int) int {
	switch {
	case n > 0:
		return 1
	case n < 0:
		return -1
	}
	return 0
}

func score(h, a int) string { return fmt.Sprintf("%d-%d", h, a) }

// BuildTable assembles the table for a round, with players ordered by
// descending season score (ties broken by name).
func (s *Sources) BuildTable(roundID model.RoundID) (*Table, error) {
	fixtures, err := s.Fixtures.ListFixtures(roundID)
	if err != nil {
		return nil, fmt.Errorf("listing fixtures for round %v: %w", roundID, err)
	}
	results, err := s.Results.ListResults(roundID)
	if err != nil {
		return nil, fmt.Errorf("listing results for round %v: %w", roundID, err)
	}
	predictions, err := s.Predictions.ListPredictions(roundID)
	if err != nil {
		return nil, fmt.Errorf("listing predictions for round %v: %w", roundID, err)
	}
	players, err := s.Players.ListActivePlayers()
	if err != nil {
		return nil, fmt.Errorf("listing players: %w", err)
	}
	roundScores, err := s.Scores.CalculateRoundScores(roundID)
	if err != nil {
		return nil, fmt.Errorf("calculating scores for round %v: %w", roundID, err)
	}
	seasonScores, err := s.Scores.CalculateCurrentSeasonScores()
	if err != nil {
		return nil, fmt.Errorf("calculating current season scores: %w", err)
	}

	sort.SliceStable(players, func(i, j int) bool {
		a, b := seasonScores[players[i].Email], seasonScores[players[j].Email]
		if a != b {
			return a > b
		}
		return players[i].Name < players[j].Name
	})

	resultByGame := map[model.Game]model.Result{}
	for _, r := range results {
		resultByGame[model.Game{HomeTeam: r.HomeTeam, AwayTeam: r.AwayTeam}] = r
	}
	predByPlayer := map[string]model.GamePredictions{}
	for _, p := range predictions {
		if predByPlayer[p.Player] == nil {
			predByPlayer[p.Player] = model.GamePredictions{}
		}
		predByPlayer[p.Player][p.Game] = p.Prediction
	}

	t := &Table{}
	for _, p := range players {
		t.Players = append(t.Players, p.Name)
		t.RoundScores = append(t.RoundScores, roundScores[p.Email])
		t.SeasonScores = append(t.SeasonScores, seasonScores[p.Email])
	}
	for _, f := range fixtures {
		game := model.Game{HomeTeam: f.HomeTeam, AwayTeam: f.AwayTeam}
		home, err := s.Teams.FindTeamByKey(f.HomeTeam)
		if err != nil {
			return nil, fmt.Errorf("finding home team %q: %w", f.HomeTeam, err)
		}
		away, err := s.Teams.FindTeamByKey(f.AwayTeam)
		if err != nil {
			return nil, fmt.Errorf("finding away team %q: %w", f.AwayTeam, err)
		}
		row := Row{Label: home.ShortName + " v " + away.ShortName, Result: "-"}
		result, hasResult := resultByGame[game]
		if hasResult {
			row.Result = score(result.HomeGoals, result.AwayGoals)
		}
		for _, p := range players {
			cell := Cell{Text: "-"}
			if pred, ok := predByPlayer[p.Email][game]; ok {
				cell.Text = score(pred.HomeGoals, pred.AwayGoals)
				if hasResult {
					cell.Points = PredictionPoints(pred, result)
				}
			}
			row.Predictions = append(row.Predictions, cell)
		}
		t.Rows = append(t.Rows, row)
	}
	return t, nil
}
