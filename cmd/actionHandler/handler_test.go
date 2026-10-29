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
	skip   bool
}

func (f *fakePredictions) AddPredictions(p string, _ model.RoundID, g model.GamePredictions) (model.GamePredictions, error) {
	f.player, f.saved = p, g
	if f.skip {
		return make(model.GamePredictions), nil
	}
	return g, nil
}

func post(h *handler, body string) string {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body)))
	return rec.Body.String()
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

type fakeNotifier struct {
	calls  int
	player model.Player
	preds  model.GamePredictions
	err    error
}

func (f *fakeNotifier) NotifyPredictions(_ context.Context, p model.Player, _ model.RoundID, g model.GamePredictions) error {
	f.calls++
	f.player, f.preds = p, g
	return f.err
}

func TestSubmitPredictionsNotifies(t *testing.T) {
	n := &fakeNotifier{}
	h := &handler{
		playerRepository:     fakePlayers{player: &model.Player{Email: "a@b.c", Name: "A", Active: true}},
		predictionRepository: &fakePredictions{},
		notifier:             n,
	}
	post(h, submitBody)
	if n.calls != 1 || n.player.Name != "A" || len(n.preds) != 1 {
		t.Errorf("unexpected notification: %+v", n)
	}
}

func TestSubmitPredictionsNotifierFailureIsNonFatal(t *testing.T) {
	n := &fakeNotifier{err: errors.New("boom")}
	h := &handler{
		playerRepository:     fakePlayers{player: &model.Player{Email: "a@b.c", Active: true}},
		predictionRepository: &fakePredictions{},
		notifier:             n,
	}
	if out := post(h, submitBody); !strings.Contains(out, "Predictions submitted") {
		t.Errorf("unexpected response: %s", out)
	}
}

func TestSubmitPredictionsExplainsStartedMatchesWereSkipped(t *testing.T) {
	preds := &fakePredictions{skip: true}
	n := &fakeNotifier{}
	h := &handler{
		playerRepository:     fakePlayers{player: &model.Player{Email: "a@b.c", Active: true}},
		predictionRepository: preds,
		notifier:             n,
	}
	out := post(h, submitBody)
	if !strings.Contains(out, "Started matches are locked") || !strings.Contains(out, "Existing predictions were kept") {
		t.Errorf("expected started-match notification, got %s", out)
	}
	if n.calls != 0 {
		t.Errorf("expected no notification for skipped predictions, got %d calls", n.calls)
	}
}

func TestNoNotificationOnRejection(t *testing.T) {
	n := &fakeNotifier{}
	h := &handler{playerRepository: fakePlayers{}, predictionRepository: &fakePredictions{}, notifier: n}
	post(h, submitBody)
	if n.calls != 0 {
		t.Errorf("expected no notification")
	}
}
