package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alanraison/fleps/pkg/chat"
	chatmodel "github.com/alanraison/fleps/pkg/chat/model"
	"github.com/alanraison/fleps/pkg/model"
)

type handler struct {
	playerRepository     model.PlayerRepository
	predictionRepository model.PredictionRepository
}

var errNotAPlayer = errors.New("user is not a registered player")

// ServeHTTP handles card actions. Callers are authenticated by Cloud Run IAM
// before the request reaches this handler, and the user is taken from the event.
// Failures are reported as a Chat notification with a 200 status, as Chat only
// displays the content of successful responses.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event, err := chat.NewActionEvent(r.Body)
	if err != nil {
		fmt.Printf("Error parsing action event: %v\n", err)
		writeJSON(w, chat.NewErrorNotification("Sorry, Fleps couldn't understand that request."))
		return
	}
	switch event.Parameters["actionName"] {
	case chat.SubmitPredictionsAction:
		if err := h.submitPredictions(event); err != nil {
			fmt.Printf("Error submitting predictions: %v\n", err)
			msg := "Sorry, your predictions couldn't be saved. Please try again later."
			if errors.Is(err, errNotAPlayer) {
				msg = "Sorry, you need to be a registered player to make predictions."
			} else if errors.Is(err, errInvalidPrediction) {
				msg = "Please enter a whole number of goals for every match."
			}
			writeJSON(w, chat.NewErrorNotification(msg))
			return
		}
		writeJSON(w, chatmodel.RenderActions{
			Action: &chatmodel.RenderAction{
				Navigations: []chatmodel.Navigation{
					{
						EndNavigation: &chatmodel.EndNavigation{
							Action: chatmodel.EndNavigationCloseDialog,
						},
					},
				},
				Notification: &chatmodel.Notification{Text: "Predictions submitted"},
			},
		})
	default:
		fmt.Printf("Unhandled action: %q\n", event.Parameters["actionName"])
		writeJSON(w, chat.NewErrorNotification("Sorry, Fleps couldn't understand that request."))
	}
}

var errInvalidPrediction = errors.New("invalid prediction")

// submitPredictions saves the form inputs as predictions for the invoking user.
func (h *handler) submitPredictions(event *chat.ActionEvent) error {
	if event.UserEmail == "" {
		return fmt.Errorf("action has no user: %w", errNotAPlayer)
	}
	player, err := h.playerRepository.GetPlayerByEmail(event.UserEmail)
	if err != nil {
		return fmt.Errorf("getting player %q: %w", event.UserEmail, err)
	}
	if player == nil || !player.Active {
		return fmt.Errorf("player %q: %w", event.UserEmail, errNotAPlayer)
	}
	roundID := model.RoundID(event.Parameters["roundId"])
	if roundID == "" {
		return fmt.Errorf("action has no round: %w", errInvalidPrediction)
	}
	predictions := make(model.GamePredictions)
	for name, value := range event.FormInputs {
		game, side, ok := chat.ParseGoalsInputName(name)
		if !ok {
			continue
		}
		goals, err := parseGoals(value)
		if err != nil {
			return fmt.Errorf("%s goals for %v v %v: %w", side, game.HomeTeam, game.AwayTeam, err)
		}
		p := predictions[game]
		if side == chat.HomeSide {
			p.HomeGoals = goals
		} else {
			p.AwayGoals = goals
		}
		predictions[game] = p
	}
	if len(predictions) == 0 {
		return fmt.Errorf("no predictions submitted: %w", errInvalidPrediction)
	}
	if err := h.predictionRepository.AddPredictions(player.Email, roundID, predictions); err != nil {
		return fmt.Errorf("saving predictions for %q: %w", player.Email, err)
	}
	return nil
}

func parseGoals(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("%q: %w", s, errInvalidPrediction)
	}
	return n, nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	writeJSONBody(w, v)
}

func writeJSONBody(w http.ResponseWriter, v any) {
	bytes, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		return
	}
	w.Write(bytes)
}
