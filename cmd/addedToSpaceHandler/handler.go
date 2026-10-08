package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alanraison/fleps/pkg/chat"
	"github.com/alanraison/fleps/pkg/model"
)

type handler struct {
	playerRepository model.PlayerRepository
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event, err := chat.NewAddedToSpaceEvent(r.Body)
	if err != nil {
		fmt.Printf("Error parsing chat event: %v\n", err)
		return
	}
	if event.SpaceType == chat.TypeDirectMessage {
		isNewPlayer, err := h.isNewPlayer(event.UserEmail)
		if err != nil {
			fmt.Printf("Error getting player by email: %v\n", err)
			w.WriteHeader(500)
			w.Write([]byte(`{}`))
			return
		}
		fmt.Printf("Is new player: %v, UserEmail: %s\n", isNewPlayer, event.UserEmail)
		if isNewPlayer {
			if err = h.playerRepository.AddPlayer(event.UserDisplayName, event.UserEmail); err != nil {
				fmt.Printf("Error creating new player by email: %v\n", err)
				w.WriteHeader(500)
				w.Write([]byte(`{}`))
				return
			}
			message, err := json.Marshal(chat.NewResponseMessageWithText(
				fmt.Sprintf("Welcome, %s!", event.UserDisplayName),
			))
			if err != nil {
				fmt.Printf("Error encoding response message: %v\n", err)
				w.WriteHeader(500)
				w.Write([]byte(`{}`))
				return
			}
			w.Write(message)
		} else {
			if err = h.playerRepository.EnablePlayerByEmail(event.UserEmail); err != nil {
				fmt.Printf("Error enabling player %v by email: %v\n", event.UserEmail, err)
				w.WriteHeader(500)
				w.Write([]byte(`{}`))
				return
			}
			message, err := json.Marshal(chat.NewResponseMessageWithText(
				fmt.Sprintf("Welcome back, %s!", event.UserDisplayName),
			))
			if err != nil {
				fmt.Printf("Error encoding response message: %v\n", err)
				w.WriteHeader(500)
				w.Write([]byte(`{}`))
				return
			}
			fmt.Printf("Sending response message: %s\n", message)
			w.Write(message)
		}
	}
	w.Write([]byte("{}"))
}

func (h *handler) isNewPlayer(email string) (bool, error) {
	player, err := h.playerRepository.GetPlayerByEmail(email)
	if err != nil {
		return false, err
	}
	return player == nil, nil
}
