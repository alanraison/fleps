package csv

import (
	"strings"
	"testing"
)

func TestReadPredictionRows(t *testing.T) {
	csv := `match,player1,player2
ARSBOU,2-0,1-1,
MUNMCI,1-1,0-0,
`
	rows, err := readPredictionRows(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	row := rows[0]
	if row.Match != "ARSBOU" {
		t.Errorf("expected match 'ARSBOU', got '%s'", row.Match)
	}
	if row.Predictions["player1"] != "2-0" {
		t.Errorf("expected player1 '2-0', got '%s'", row.Predictions["player1"])
	}
	if row.Predictions["player2"] != "1-1" {
		t.Errorf("expected player2 '1-1', got '%s'", row.Predictions["player2"])
	}
	row = rows[1]
	if row.Match != "MUNMCI" {
		t.Errorf("expected match 'MUNMCI', got '%s'", row.Match)
	}
	if row.Predictions["player1"] != "1-1" {
		t.Errorf("expected player1 '1-1', got '%s'", row.Predictions["player1"])
	}
	if row.Predictions["player2"] != "0-0" {
		t.Errorf("expected player2 '0-0', got '%s'", row.Predictions["player2"])
	}
}
