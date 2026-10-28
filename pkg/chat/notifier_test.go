package chat

import (
	"context"
	"errors"
	"testing"

	chatmodel "github.com/alanraison/fleps/pkg/chat/model"
	"github.com/alanraison/fleps/pkg/model"
)

type fakeTeams map[model.TeamKey]*model.Team

func (f fakeTeams) FindTeamByKey(k model.TeamKey) (*model.Team, error) {
	if t, ok := f[k]; ok {
		return t, nil
	}
	return nil, model.UnknownTeamErr
}

type fakeBadges map[string]string

type fakeLogos map[model.TeamKey]string

func (f fakeLogos) FindLogoSlugByTeamKey(k model.TeamKey) (string, error) {
	if s, ok := f[k]; ok {
		return s, nil
	}
	return "", model.UnknownLogoErr
}

func (f fakeBadges) BadgeURL(_ context.Context, t model.Team, _ string) (string, error) {
	return f[t.ShortName], nil
}

func newTestNotifier(create func(context.Context, string, *chatmodel.Message) error) *Notifier {
	return &Notifier{
		space: "spaces/X",
		teams: fakeTeams{
			"ARS": {Key: "ARS", ShortName: "Arsenal"},
			"CHE": {Key: "CHE", ShortName: "Chelsea"},
			"BMT": {Key: "BMT", ShortName: "B'mth"},
		},
		logos:  fakeLogos{"ARS": "arsenal"},
		badges: fakeBadges{"Arsenal": "http://img/ars.png", "Chelsea": "http://img/che.png"},
		create: create,
	}
}

func TestNotifyPredictionsBuildsCard(t *testing.T) {
	var space string
	var msg *chatmodel.Message
	n := newTestNotifier(func(_ context.Context, s string, m *chatmodel.Message) error {
		space, msg = s, m
		return nil
	})
	err := n.NotifyPredictions(context.Background(), model.Player{Name: "Alan"}, "R1", model.GamePredictions{
		{HomeTeam: "BMT", AwayTeam: "XXX"}: {HomeGoals: 0, AwayGoals: 0},
		{HomeTeam: "ARS", AwayTeam: "CHE"}: {HomeGoals: 2, AwayGoals: 1},
	})
	if err != nil || space != "spaces/X" {
		t.Fatalf("unexpected: %v %q", err, space)
	}
	card := msg.CardsV2[0].Card
	if err := card.Validate(); err != nil {
		t.Fatalf("invalid card: %v", err)
	}
	if card.Header.Title != "Alan has submitted predictions" || card.Header.Subtitle != "Round R1" {
		t.Errorf("unexpected header: %+v", card.Header)
	}
	ws := card.Sections[0].Widgets
	if len(ws) != 2 {
		t.Fatalf("expected 2 widgets, got %d", len(ws))
	}
	first := ws[0].DecoratedText
	if first.Text != "<b>Arsenal</b> 2 – 1 <b>Chelsea</b>" ||
		first.StartIcon.IconURL != "http://img/ars.png" || first.EndIcon.IconURL != "http://img/che.png" {
		t.Errorf("unexpected first widget: %+v", first)
	}
	// Unknown team falls back to its key with no badge; missing badge omits the icon.
	second := ws[1].DecoratedText
	if second.Text != "<b>B&#39;mth</b> 0 – 0 <b>XXX</b>" || second.StartIcon != nil || second.EndIcon != nil {
		t.Errorf("unexpected second widget: %+v", second)
	}
}

func TestNotifyPredictionsWrapsError(t *testing.T) {
	boom := errors.New("boom")
	n := newTestNotifier(func(context.Context, string, *chatmodel.Message) error { return boom })
	if err := n.NotifyPredictions(context.Background(), model.Player{}, "R1", nil); !errors.Is(err, boom) {
		t.Errorf("expected wrapped error, got %v", err)
	}
}
