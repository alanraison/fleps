package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"sort"

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

// Notifier posts messages to a Chat space as the Fleps app.
type Notifier struct {
	space  string
	teams  TeamFinder
	logos  LogoFinder
	badges BadgeFinder
	create func(ctx context.Context, space string, msg *chatmodel.Message) error
}

// NewNotifier creates a Notifier posting to the given space (e.g. "spaces/AAAA")
// using application default credentials.
func NewNotifier(ctx context.Context, space string, teams TeamFinder, logos LogoFinder, badges BadgeFinder) (*Notifier, error) {
	svc, err := chatapi.NewService(ctx, option.WithScopes("https://www.googleapis.com/auth/chat.bot"))
	if err != nil {
		return nil, fmt.Errorf("creating chat service: %w", err)
	}
	return &Notifier{
		space:  space,
		teams:  teams,
		logos:  logos,
		badges: badges,
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
	msg := n.predictionsMessage(ctx, player.Name, round, predictions)
	if err := n.create(ctx, n.space, msg); err != nil {
		return fmt.Errorf("posting predictions for %q to %q: %w", player.Email, n.space, err)
	}
	return nil
}

// side is a team as displayed in the card.
type side struct {
	name  string
	badge string
}

func (n *Notifier) side(ctx context.Context, key model.TeamKey) side {
	s := side{name: string(key)}
	team, err := n.teams.FindTeamByKey(key)
	if err != nil {
		fmt.Printf("Error finding team %q: %v\n", key, err)
		return s
	}
	s.name = team.ShortName
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

func (n *Notifier) predictionsMessage(ctx context.Context, name string, round model.RoundID, predictions model.GamePredictions) *chatmodel.Message {
	games := make([]model.Game, 0, len(predictions))
	for g := range predictions {
		games = append(games, g)
	}
	sort.Slice(games, func(i, j int) bool {
		if games[i].HomeTeam != games[j].HomeTeam {
			return games[i].HomeTeam < games[j].HomeTeam
		}
		return games[i].AwayTeam < games[j].AwayTeam
	})
	widgets := make([]chatmodel.Widget, 0, len(games))
	for _, g := range games {
		widgets = append(widgets, predictionWidget(n.side(ctx, g.HomeTeam), n.side(ctx, g.AwayTeam), predictions[g]))
	}
	return &chatmodel.Message{
		Text: fmt.Sprintf("%s has submitted predictions for round %s", name, round),
		CardsV2: []chatmodel.CardWithID{{
			CardID: "predictions",
			Card: &chatmodel.Card{
				Header:   &chatmodel.CardHeader{Title: name + " has submitted predictions", Subtitle: fmt.Sprintf("Round %s", round)},
				Sections: []chatmodel.Section{{Widgets: widgets}},
			},
		}},
	}
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
