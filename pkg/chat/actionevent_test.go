package chat

import (
	"strings"
	"testing"
)

func TestNewActionEvent(t *testing.T) {
	body := `{"commonEventObject":{"parameters":{"actionName":"submitPredictions","roundId":"R1"},
	"formInputs":{"ARS":{"stringInputs":{"value":["2"]}}}},
	"chat":{"user":{"email":"a@b.c"}}}`
	e, err := NewActionEvent(strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.UserEmail != "a@b.c" || e.Parameters["roundId"] != "R1" || e.FormInputs["ARS"] != "2" {
		t.Errorf("unexpected event: %+v", e)
	}
}

func TestNewActionEventInvalidJSON(t *testing.T) {
	if _, err := NewActionEvent(strings.NewReader("{")); err == nil {
		t.Fatal("expected error")
	}
}
