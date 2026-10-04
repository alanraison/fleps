package main

import (
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
	if event.SpaceType != chat.TypeDirectMessage {
		fmt.Printf("Ignoring non-direct message event. SpaceType: %s\n", event.SpaceType)
		w.Write([]byte(`{}`))
		return
	}
	fmt.Printf("Removing player with email: %s\n", event.UserEmail)
	err = h.playerRepository.DisablePlayerByEmail(event.UserEmail)
	if err != nil {
		fmt.Printf("Error removing player by email: %v\n", err)
		w.Write([]byte(`{"error": "could not remove player"}`))
		return
	}
	w.Write([]byte(`{"status": "player removed"}`))
}
