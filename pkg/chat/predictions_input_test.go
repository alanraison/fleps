package chat

import "testing"

func TestGoalsInputNameRoundTrip(t *testing.T) {
	game, side, ok := ParseGoalsInputName(GoalsInputName("ARS", "LEE", AwaySide))
	if !ok || game.HomeTeam != "ARS" || game.AwayTeam != "LEE" || side != AwaySide {
		t.Errorf("unexpected parse: %v %v %v", game, side, ok)
	}
	if _, _, ok := ParseGoalsInputName("ARS"); ok {
		t.Error("expected non-prediction input to be rejected")
	}
}
