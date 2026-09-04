package csv

import (
	"errors"
	"strings"
	"testing"

	"github.com/alanraison/predictions/pkg/model"
)

func TestReadResultRows(t *testing.T) {
	r := strings.NewReader(`BHA,FUL,2,1
LEE,BRE,0,0`)
	csv := NewCsv(&mockTeamRepository{}, nil)

	results, err := csv.ReadResultRows(r, "R1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].RoundID != "R1" || results[0].HomeTeam != "BHA" || results[0].AwayTeam != "FUL" {
		t.Fatalf("unexpected first result fixture: %+v", results[0].Fixture)
	}
	if results[0].HomeScore != 2 || results[0].AwayScore != 1 {
		t.Fatalf("unexpected first result score: %+v", results[0])
	}
}

func TestReadResultRows_UnknownTeam(t *testing.T) {
	r := strings.NewReader(`XYZ,FUL,2,1`)
	csv := NewCsv(&mockTeamRepository{}, nil)

	_, err := csv.ReadResultRows(r, "R1")
	if err == nil {
		t.Fatal("expected error for unknown team, got nil")
	}
	if !errors.Is(err, model.UnknownTeamErr) {
		t.Fatalf("expected UnknownTeamErr, got %v", err)
	}
}

func TestReadResultRows_BadScore(t *testing.T) {
	r := strings.NewReader(`BHA,FUL,abc,1`)
	csv := NewCsv(&mockTeamRepository{}, nil)

	_, err := csv.ReadResultRows(r, "R1")
	if err == nil {
		t.Fatal("expected error for invalid score, got nil")
	}
}
