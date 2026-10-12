package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alanraison/fleps/pkg/model"
)

type fakePlayers struct {
	model.PlayerRepository
	player   *model.Player
	getErr   error
	added    []string
	enabled  []string
	disabled []string
	disErr   error
}

func (f *fakePlayers) GetPlayerByEmail(string) (*model.Player, error) { return f.player, f.getErr }
func (f *fakePlayers) AddPlayer(name, email string) error {
	f.added = append(f.added, email)
	return nil
}
func (f *fakePlayers) EnablePlayerByEmail(email string) error {
	f.enabled = append(f.enabled, email)
	return nil
}
func (f *fakePlayers) DisablePlayerByEmail(email string) error {
	f.disabled = append(f.disabled, email)
	return f.disErr
}

func body(kind, spaceType string) string {
	return `{"chat":{"user":{"displayName":"Ann","email":"ann@example.com"},"` + kind + `Payload":{"space":{"type":"` + spaceType + `"}}}}`
}

func do(h *handler, _, b string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(b)))
	return rec
}

func TestAddedNewPlayer(t *testing.T) {
	p := &fakePlayers{}
	rec := do(&handler{p}, "", body("addedToSpace", "DM"))
	if rec.Code != 200 || len(p.added) != 1 || !strings.Contains(rec.Body.String(), "Welcome, Ann!") {
		t.Fatalf("code=%d added=%v body=%s", rec.Code, p.added, rec.Body)
	}
}

func TestAddedExistingPlayer(t *testing.T) {
	p := &fakePlayers{player: &model.Player{}}
	rec := do(&handler{p}, "", body("addedToSpace", "DM"))
	if rec.Code != 200 || len(p.enabled) != 1 || !strings.Contains(rec.Body.String(), "Welcome back, Ann!") {
		t.Fatalf("code=%d enabled=%v body=%s", rec.Code, p.enabled, rec.Body)
	}
}

func TestAddedIgnoresRooms(t *testing.T) {
	p := &fakePlayers{}
	rec := do(&handler{p}, "", body("addedToSpace", "ROOM"))
	if rec.Code != 200 || len(p.added) != 0 || rec.Body.String() != "{}" {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body)
	}
}

const notification = `"notification":{"text":`

func TestAddedRepositoryError(t *testing.T) {
	rec := do(&handler{&fakePlayers{getErr: errors.New("boom")}}, "", body("addedToSpace", "DM"))
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), notification) || strings.Contains(rec.Body.String(), "boom") {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body)
	}
}

func TestRemoved(t *testing.T) {
	p := &fakePlayers{}
	rec := do(&handler{p}, "", body("removedFromSpace", "DM"))
	if rec.Code != 200 || len(p.disabled) != 1 || rec.Body.String() != "{}" {
		t.Fatalf("code=%d disabled=%v body=%s", rec.Code, p.disabled, rec.Body)
	}
	p = &fakePlayers{}
	do(&handler{p}, "", body("removedFromSpace", "ROOM"))
	if len(p.disabled) != 0 {
		t.Fatalf("room event disabled player")
	}
	rec = do(&handler{&fakePlayers{disErr: errors.New("boom")}}, "", body("removedFromSpace", "DM"))
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), notification) {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body)
	}
}

func TestBadBodyAndUnknownEvent(t *testing.T) {
	h := &handler{&fakePlayers{}}
	for _, b := range []string{"{", `{"chat":{"appCommandPayload":{}}}`} {
		if rec := do(h, "", b); rec.Code != 200 || !strings.Contains(rec.Body.String(), notification) {
			t.Fatalf("body %q: code=%d body=%s", b, rec.Code, rec.Body)
		}
	}
}
