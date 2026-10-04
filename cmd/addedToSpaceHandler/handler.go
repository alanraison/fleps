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
	event, err := chat.NewChatEvent(r.Body)
	if err != nil {
		fmt.Printf("Error parsing chat event: %v\n", err)
		return
	}
	if event.SpaceType == chat.TypeDirectMessage {
		isNewPlayer, err := h.isNewPlayer(event.UserEmail)
		if err != nil {
			fmt.Printf("Error getting player by email: %v\n", err)
			w.Write([]byte(`{"error": "could not determine if new player"}`))
			return
		}
		fmt.Printf("Is new player: %v, UserEmail: %s\n", isNewPlayer, event.UserEmail)
		if isNewPlayer {
			if err = h.playerRepository.AddPlayer(event.UserDisplayName, event.UserEmail); err != nil {
				fmt.Printf("Error creating new player by email: %v\n", err)
				w.Write([]byte(`{"error": "could not create new player"}`))
				return
			}
			message, err := json.Marshal(chat.NewResponseMessageWithText(
				fmt.Sprintf("Welcome, %s!", event.UserDisplayName),
			))
			if err != nil {
				fmt.Printf("Error encoding response message: %v\n", err)
				w.Write([]byte(`{"text": "could not encode response message"}`))
				return
			}
			w.Write(message)
		} else {
			if err = h.playerRepository.EnablePlayerByEmail(event.UserEmail); err != nil {
				fmt.Printf("Error enabling player by email: %v\n", err)
				errorMessage, err := json.Marshal(chat.NewResponseMessageWithText("could not enable player"))
				if err != nil {
					fmt.Printf("Error encoding error message: %v\n", err)
					return
				}
				w.Write(errorMessage)
				return
			}
			message, err := json.Marshal(chat.NewResponseMessageWithText(
				fmt.Sprintf("Welcome back, %s!", event.UserDisplayName),
			))
			if err != nil {
				fmt.Printf("Error encoding response message: %v\n", err)
				w.Write([]byte(`{"text": "could not encode response message"}`))
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
