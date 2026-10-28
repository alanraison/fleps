package chat

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"slices"
	"time"

	chatmodel "github.com/alanraison/fleps/pkg/chat/model"
	"github.com/alanraison/fleps/pkg/model"
	chatapi "google.golang.org/api/chat/v1"
	"google.golang.org/api/option"
)

// TeamFinder looks up teams by key.
type TeamFinder interface {
	FindTeamByKey(key model.TeamKey) (*model.Team, error)
}

// LogoFinder looks up the logo slug for a team by key.
type LogoFinder interface {
	FindLogoSlugByTeamKey(key model.TeamKey) (string, error)
}

// BadgeFinder looks up the badge image URL for a team and its logo slug
// (which may be empty), returning "" if there is none.
type BadgeFinder interface {
	BadgeURL(ctx context.Context, team model.Team, slug string) (string, error)
}

// FixtureFinder lists the fixtures of a round.
type FixtureFinder interface {
	ListFixtures(round model.RoundID) ([]model.Fixture, error)
}

// Notifier posts messages to a Chat space as the Fleps app.
type Notifier struct {
	space    string
	teams    TeamFinder
	fixtures FixtureFinder
	logos    LogoFinder
	badges   BadgeFinder
	create   func(ctx context.Context, space string, msg *chatmodel.Message) error
}

// NewNotifier creates a Notifier posting to the given space (e.g. "spaces/AAAA")
// using application default credentials.
func NewNotifier(ctx context.Context, space string, teams TeamFinder, fixtures FixtureFinder, logos LogoFinder, badges BadgeFinder) (*Notifier, error) {
	svc, err := chatapi.NewService(ctx, option.WithScopes("https://www.googleapis.com/auth/chat.bot"))
	if err != nil {
		return nil, fmt.Errorf("creating chat service: %w", err)
	}
	return &Notifier{
		space:    space,
		teams:    teams,
		fixtures: fixtures,
		logos:    logos,
		badges:   badges,
		create: func(ctx context.Context, space string, msg *chatmodel.Message) error {
			raw, err := json.Marshal(msg)
			if err != nil {
				return fmt.Errorf("encoding message: %w", err)
			}
			var apiMsg chatapi.Message
			if err := json.Unmarshal(raw, &apiMsg); err != nil {
				return fmt.Errorf("converting message: %w", err)
			}
			if _, err := svc.Spaces.Messages.Create(space, &apiMsg).Context(ctx).Do(); err != nil {
				return fmt.Errorf("calling chat API: %w", err)
			}
			return nil
		},
	}, nil
}

// NotifyPredictions announces a player's submitted predictions in the space.
func (n *Notifier) NotifyPredictions(ctx context.Context, player model.Player, round model.RoundID, predictions model.GamePredictions) error {
	msg, err := n.predictionsMessage(ctx, player.Name, round, predictions)
	if err != nil {
		return fmt.Errorf("building predictions message for %q: %w", player.Email, err)
	}
	if err := n.create(ctx, n.space, msg); err != nil {
		return fmt.Errorf("posting predictions for %q to %q: %w", player.Email, n.space, err)
	}
	return nil
}

// side is a team as displayed in the card.
type side struct {
	name   string
	badge  string
	league int
}

func (n *Notifier) side(ctx context.Context, key model.TeamKey) side {
	s := side{name: string(key)}
	team, err := n.teams.FindTeamByKey(key)
	if err != nil {
		fmt.Printf("Error finding team %q: %v\n", key, err)
		return s
	}
	s.name = team.ShortName
	s.league = team.League
	slug, err := n.logos.FindLogoSlugByTeamKey(key)
	if err != nil && !errors.Is(err, model.UnknownLogoErr) {
		fmt.Printf("Error finding logo for team %q: %v\n", key, err)
	}
	badge, err := n.badges.BadgeURL(ctx, *team, slug)
	if err != nil {
		fmt.Printf("Error finding badge for team %q: %v\n", key, err)
		return s
	}
	s.badge = badge
	return s
}

func (n *Notifier) predictionsMessage(ctx context.Context, name string, round model.RoundID, predictions model.GamePredictions) (*chatmodel.Message, error) {
	fixtures, err := n.fixtures.ListFixtures(round)
	if err != nil {
		return nil, fmt.Errorf("listing fixtures for round %q: %w", round, err)
	}
	var entries []entry
	for _, f := range fixtures {
		if p, ok := predictions[model.Game{HomeTeam: f.HomeTeam, AwayTeam: f.AwayTeam}]; ok {
			entries = append(entries, entry{
				date: f.Date, home: n.side(ctx, f.HomeTeam), away: n.side(ctx, f.AwayTeam), prediction: p,
			})
		}
	}
	// Kick-off time, then the home team's league rank, then team names.
	slices.SortStableFunc(entries, func(a, b entry) int {
		return cmp.Or(
			a.date.Compare(b.date),
			cmp.Compare(a.home.league, b.home.league),
			cmp.Compare(a.home.name, b.home.name),
			cmp.Compare(a.away.name, b.away.name),
		)
	})
	widgets := make([]chatmodel.Widget, 0, len(entries))
	for _, e := range entries {
		widgets = append(widgets, predictionWidget(e.home, e.away, e.prediction))
	}
	return &chatmodel.Message{
		CardsV2: []chatmodel.CardWithID{{
			CardID: "predictions",
			Card: &chatmodel.Card{
				Header:   &chatmodel.CardHeader{Title: name + " has submitted predictions", Subtitle: fmt.Sprintf("Round %s", round)},
				Sections: []chatmodel.Section{{Widgets: widgets}},
			},
		}},
	}, nil
}

func predictionWidget(home, away side, p model.Prediction) chatmodel.Widget {
	dt := &chatmodel.DecoratedText{
		Text: fmt.Sprintf("<b>%s</b> %d – %d <b>%s</b>", html.EscapeString(home.name), p.HomeGoals, p.AwayGoals, html.EscapeString(away.name)),
	}
	if home.badge != "" {
		dt.StartIcon = &chatmodel.Icon{IconURL: home.badge, AltText: home.name}
	}
	if away.badge != "" {
		dt.EndIcon = &chatmodel.Icon{IconURL: away.badge, AltText: away.name}
	}
	return chatmodel.Widget{DecoratedText: dt}
}

// entry is a predicted game along with what is needed to order it.
type entry struct {
	date       time.Time
	home, away side
	prediction model.Prediction
}
