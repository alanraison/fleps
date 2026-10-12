package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alanraison/fleps/pkg/model"
)

type fakePlayers struct {
	model.PlayerRepository
	player *model.Player
}

func (f fakePlayers) GetPlayerByEmail(string) (*model.Player, error) { return f.player, nil }

type fakePredictions struct {
	model.PredictionRepository
	player string
	saved  model.GamePredictions
}

func (f *fakePredictions) AddPredictions(p string, _ model.RoundID, g model.GamePredictions) error {
	f.player, f.saved = p, g
	return nil
}

type fakeVerifier struct {
	email string
	err   error
}

func (f fakeVerifier) VerifiedEmail(context.Context, string, string) (string, error) {
	return f.email, f.err
}

func post(h *handler, body string) string {
	rec, _ := postAuth(h, body, "Bearer tok")
	return rec.Body.String()
}

func postAuth(h *handler, body, auth string) (*httptest.ResponseRecorder, error) {
	if h.verifier == nil {
		h.verifier = fakeVerifier{email: "a@b.c"}
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	h.ServeHTTP(rec, req)
	return rec, nil
}

func TestRejectsRequestsThatFailAuthentication(t *testing.T) {
	player := fakePlayers{player: &model.Player{Email: "a@b.c", Active: true}}
	tests := map[string]struct {
		verifier fakeVerifier
		auth     string
	}{
		"no token":         {fakeVerifier{email: "a@b.c"}, ""},
		"invalid token":    {fakeVerifier{err: errors.New("bad")}, "Bearer tok"},
		"mismatched email": {fakeVerifier{email: "someone@else.com"}, "Bearer tok"},
		"no email":         {fakeVerifier{}, "Bearer tok"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			preds := &fakePredictions{}
			h := &handler{verifier: tc.verifier, playerRepository: player, predictionRepository: preds}
			rec, _ := postAuth(h, submitBody, tc.auth)
			if rec.Code != http.StatusUnauthorized || preds.saved != nil {
				t.Errorf("expected 401 and nothing saved, got %d, saved %v", rec.Code, preds.saved)
			}
		})
	}
}

const submitBody = `{"commonEventObject":{"parameters":{"actionName":"submitPredictions","roundId":"R1"},
"formInputs":{"ARS_LEE_home":{"stringInputs":{"value":["2"]}},"ARS_LEE_away":{"stringInputs":{"value":["1"]}}}},
"chat":{"user":{"email":"a@b.c"}}}`

func TestSubmitPredictionsSavesForInvokingUser(t *testing.T) {
	preds := &fakePredictions{}
	h := &handler{
		playerRepository:     fakePlayers{player: &model.Player{Email: "a@b.c", Active: true}},
		predictionRepository: preds,
	}
	out := post(h, submitBody)
	if !strings.Contains(out, "Predictions submitted") {
		t.Errorf("unexpected response: %s", out)
	}
	got := preds.saved[model.Game{HomeTeam: "ARS", AwayTeam: "LEE"}]
	if preds.player != "a@b.c" || got != (model.Prediction{HomeGoals: 2, AwayGoals: 1}) {
		t.Errorf("unexpected save: %v %v", preds.player, preds.saved)
	}
}

func TestSubmitPredictionsRejectsUnknownPlayer(t *testing.T) {
	preds := &fakePredictions{}
	h := &handler{playerRepository: fakePlayers{}, predictionRepository: preds}
	out := post(h, submitBody)
	if !strings.Contains(out, "registered player") || preds.saved != nil {
		t.Errorf("expected rejection, got %s", out)
	}
}

func TestSubmitPredictionsRejectsInvalidGoals(t *testing.T) {
	preds := &fakePredictions{}
	h := &handler{
		playerRepository:     fakePlayers{player: &model.Player{Email: "a@b.c", Active: true}},
		predictionRepository: preds,
	}
	out := post(h, strings.Replace(submitBody, `["2"]`, `["x"]`, 1))
	if !strings.Contains(out, "whole number") || preds.saved != nil {
		t.Errorf("expected rejection, got %s", out)
	}
}
