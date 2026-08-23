package csv

import (
	"fmt"
	"strings"
	"testing"

	"github.com/alanraison/predictions/pkg/model"
)

var (
	c = &Csv{
		teamRepo:    &mockTeamRepository{},
		fixtureRepo: &mockFixtureRepository{},
	}
	teams = map[string]model.Team{
		"ARS": {Key: "ARS", FullName: "Arsenal", ShortName: "Arsenal"},
		"BOU": {Key: "BOU", FullName: "Bournemouth", ShortName: "Bournemouth"},
		"LEE": {Key: "LEE", FullName: "Leeds United", ShortName: "Leeds"},
		"LIV": {Key: "LIV", FullName: "Liverpool", ShortName: "Liverpool"},
	}
)

type mockTeamRepository struct{}

type mockFixtureRepository struct{
	fixtures []model.Fixture
}

func (m *mockTeamRepository) FindByKey(key string) (model.Team, error) {
	if team, ok := teams[key]; ok {
		return team, nil
	}
	return model.Team{}, fmt.Errorf("team %s not found", key)
}

func (m *mockFixtureRepository) AddFixtures(fixtures []model.Fixture) error {
	m.fixtures = append(m.fixtures, fixtures...)
	return nil
}

func TestReadFixtures(t *testing.T) {
	csvData := `home_team,away_team,date
ARS,BOU,2026-09-01T15:00
LEE,LIV,2026-09-02T16:30
`
	fixtures, err := c.readFixtures(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fixtures) != 2 {
		t.Fatalf("expected 2 fixtures, got %d", len(fixtures))
	}
	fixture := fixtures[0]
	if fixture.HomeTeam.Key != "ARS" {
		t.Errorf("expected home team 'ARS', got '%s'", fixture.HomeTeam.Key)
	}
	if fixture.AwayTeam.Key != "BOU" {
		t.Errorf("expected away team 'BOU', got '%s'", fixture.AwayTeam.Key)
	}
	if fixture.Date.Format("2006-01-02T15:04") != "2026-09-01T15:00" {
		t.Errorf("expected date '2026-09-01T15:00', got '%s'", fixture.Date.Format("2006-01-02T15:04"))
	}
	fixture = fixtures[1]
	if fixture.HomeTeam.Key != "LEE" {
		t.Errorf("expected home team 'LEE', got '%s'", fixture.HomeTeam.Key)
	}
	if fixture.AwayTeam.Key != "LIV" {
		t.Errorf("expected away team 'LIV', got '%s'", fixture.AwayTeam.Key)
	}
	if fixture.Date.Format("2006-01-02T15:04") != "2026-09-02T16:30" {
		t.Errorf("expected date '2026-09-02T16:30', got '%s'", fixture.Date.Format("2006-01-02T15:04"))
	}
}

func TestReadFixturesInvalidTeam(t *testing.T) {
	csvData := `home_team,away_team,date
ARS,XYZ,2026-09-01T15:00
`
	_, err := c.readFixtures(strings.NewReader(csvData))
	if err == nil {
		t.Fatal("expected error for invalid team, got nil")
	}
}

func TestReadFixturesInvalidDate(t *testing.T) {
	csvData := `home_team,away_team,date
ARS,BOU,invalid-date
`
	_, err := c.readFixtures(strings.NewReader(csvData))
	if err == nil {
		t.Fatal("expected error for invalid date, got nil")
	}
}
