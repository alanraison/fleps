package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alanraison/fleps/pkg/chat"
	"github.com/alanraison/fleps/pkg/model"
)

const emptyResponse = `{}`

type handler struct {
	playerRepository model.PlayerRepository
}

// ServeHTTP decodes the event and dispatches on the kind of event in the
// payload. Failures are reported as a Chat notification with a 200 status,
// as Chat only displays the content of successful responses.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event, err := chat.NewSpaceEvent(r.Body)
	if err != nil {
		fmt.Printf("Error parsing chat event: %v\n", err)
		writeError(w, "Sorry, Fleps couldn't understand that request.")
		return
	}
	switch event.Kind {
	case chat.SpaceEventAdded:
		h.addedToSpace(w, event)
	case chat.SpaceEventRemoved:
		h.removedFromSpace(w, event)
	default:
		fmt.Printf("Unhandled space event kind: %v\n", event.Kind)
		writeError(w, "Sorry, Fleps couldn't understand that request.")
	}
}

func (h *handler) addedToSpace(w http.ResponseWriter, event *chat.SpaceEvent) {
	if event.SpaceType != chat.TypeDirectMessage {
		writeJSON(w, []byte(emptyResponse))
		return
	}
	text, err := h.welcome(event)
	if err != nil {
		fmt.Printf("Error handling added-to-space event: %v\n", err)
		writeError(w, "Sorry, Fleps couldn't set up your account. Please try again later.")
		return
	}
	message, err := json.Marshal(chat.NewResponseMessageWithText(text))
	if err != nil {
		fmt.Printf("Error encoding response message: %v\n", err)
		writeError(w, "Sorry, something went wrong. Please try again later.")
		return
	}
	writeJSON(w, message)
}

// welcome registers or re-enables the player and returns the greeting to send.
func (h *handler) welcome(event *chat.SpaceEvent) (string, error) {
	player, err := h.playerRepository.GetPlayerByEmail(event.UserEmail)
	if err != nil {
		return "", fmt.Errorf("getting player by email: %w", err)
	}
	if player == nil {
		if err := h.playerRepository.AddPlayer(event.UserDisplayName, event.UserEmail); err != nil {
			return "", fmt.Errorf("creating new player %q: %w", event.UserEmail, err)
		}
		return fmt.Sprintf("Welcome, %s!", event.UserDisplayName), nil
	}
	if err := h.playerRepository.EnablePlayerByEmail(event.UserEmail); err != nil {
		return "", fmt.Errorf("enabling player %q: %w", event.UserEmail, err)
	}
	return fmt.Sprintf("Welcome back, %s!", event.UserDisplayName), nil
}

func (h *handler) removedFromSpace(w http.ResponseWriter, event *chat.SpaceEvent) {
	if event.SpaceType != chat.TypeDirectMessage {
		fmt.Printf("Ignoring non-direct message event. SpaceType: %s\n", event.SpaceType)
		writeJSON(w, []byte(emptyResponse))
		return
	}
	fmt.Printf("Removing player with email: %s\n", event.UserEmail)
	if err := h.playerRepository.DisablePlayerByEmail(event.UserEmail); err != nil {
		fmt.Printf("Error disabling player %q: %v\n", event.UserEmail, err)
		writeError(w, "Sorry, Fleps couldn't update your account. Please try again later.")
		return
	}
	writeJSON(w, []byte(emptyResponse))
}

func writeError(w http.ResponseWriter, text string) {
	body, err := json.Marshal(chat.NewErrorNotification(text))
	if err != nil {
		fmt.Printf("Error encoding error notification: %v\n", err)
		body = []byte(emptyResponse)
	}
	writeJSON(w, body)
}

func writeJSON(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
