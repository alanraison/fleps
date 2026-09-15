package csv

import (
	"strings"
	"testing"

	"github.com/alanraison/fleps/pkg/model"
)

const mockRound = model.RoundID("R1")

type mockFixtureRepo struct {
}

func (m *mockFixtureRepo) GetLatestRoundWithNoResults() (model.RoundID, error) {
	return mockRound, nil
}
func (m *mockFixtureRepo) AddRound(roundID model.RoundID, seasonID model.SeasonID) error {
	return nil
}
func (m *mockFixtureRepo) ListFixtures(roundID model.RoundID) ([]model.Fixture, error) {
	return []model.Fixture{
		{
			FixtureKey: model.FixtureKey{
				RoundID:  mockRound,
				HomeTeam: "BHA",
				AwayTeam: "LEE",
			},
		},
		{
			FixtureKey: model.FixtureKey{
				RoundID:  mockRound,
				HomeTeam: "FUL",
				AwayTeam: "BRE",
			},
		},
	}, nil
}
func (m *mockFixtureRepo) AddFixtures(fixtures []model.Fixture) error {
	return nil
}

func TestReadPredictionRows(t *testing.T) {
	csv := `BHA,LEE,2,0
FUL,BRE,1,1
`

	c := &Csv{
		teamRepo:    &mockTeamRepository{},
		fixtureRepo: &mockFixtureRepo{},
	}
	rows, err := c.ReadPredictionRows(strings.NewReader(csv), "R1", "player1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	row, ok := rows[model.Game{
		HomeTeam: "BHA",
		AwayTeam: "LEE",
	}]
	if !ok {
		t.Fatalf("expected prediction for game BHA vs LEE, got none")
	}
	if row.HomeGoals != 2 || row.AwayGoals != 0 {
		t.Errorf("expected prediction '2-0', got '%d-%d'", row.HomeGoals, row.AwayGoals)
	}
	row, ok = rows[model.Game{
		HomeTeam: "FUL",
		AwayTeam: "BRE",
	}]
	if !ok {
		t.Fatalf("expected prediction for game FUL vs BRE, got none")
	}
	if row.HomeGoals != 1 || row.AwayGoals != 1 {
		t.Errorf("expected prediction '1-1', got '%d-%d'", row.HomeGoals, row.AwayGoals)
	}
}
