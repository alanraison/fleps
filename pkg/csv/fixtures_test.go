package csv

import (
	"errors"
	"strings"
	"testing"

	"github.com/alanraison/predictions/pkg/model"
)

type mockTeamRepository struct{}

func (m *mockTeamRepository) FindTeamByKey(key model.TeamKey) (*model.Team, error) {
	knownTeams := map[model.TeamKey]model.Team{
		"BHA": {
			Key:       "BHA",
			FullName:  "Brighton & Hove Albion",
			ShortName: "BHA",
		},
		"FUL": {
			Key:       "FUL",
			FullName:  "Fulham",
			ShortName: "FUL",
		},
		"LEE": {
			Key:       "LEE",
			FullName:  "Leeds United",
			ShortName: "LEE",
		},
		"BRE": {
			Key:       "BRE",
			FullName:  "Brentford",
			ShortName: "BRE",
		},
	}
	team, ok := knownTeams[key]
	if !ok {
		return nil, model.UnknownTeamErr
	}
	return &team, nil
}

func TestReadFixtureRows(t *testing.T) {
	r := strings.NewReader(`2026-08-30 15:00,BHA,FUL
2026-08-30 14:00,LEE,BRE`)
	csv := NewCsv(&mockTeamRepository{}, nil)
	fixtures, err := csv.ReadFixtureRows(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fixtures) != 2 {
		t.Fatalf("expected 2 fixtures, got %d", len(fixtures))
	}
	if fixtures[0].HomeTeam != "BHA" || fixtures[0].AwayTeam != "FUL" {
		t.Errorf("unexpected first fixture: %+v", fixtures[0])
	}
	if fixtures[1].HomeTeam != "LEE" || fixtures[1].AwayTeam != "BRE" {
		t.Errorf("unexpected second fixture: %+v", fixtures[1])
	}
}

func TestReadFixtureRows_UnknownTeam(t *testing.T) {
	r := strings.NewReader(`2026-08-30 15:00,XYZ,FUL`)
	csv := NewCsv(&mockTeamRepository{}, nil)
	_, err := csv.ReadFixtureRows(r)
	if err == nil {
		t.Fatalf("expected error for unknown team, got nil")
	}
	if !errors.Is(err, model.UnknownTeamErr) {
		t.Fatalf("expected UnknownTeamErr, got %+v", err)
	}
}

func TestReadFixtureRows_BadDate(t *testing.T) {
	r := strings.NewReader(`2026-08-30 15:00,BHA,FUL
2026-08-30 14:00,LEE,BRE
2026-08-30 99:99,BHA,FUL`)
	csv := NewCsv(&mockTeamRepository{}, nil)
	_, err := csv.ReadFixtureRows(r)
	if err == nil {
		t.Fatalf("expected error for bad date, got nil")
	}
}
