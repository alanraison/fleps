package main

import (
	"fmt"
	"net/http"

	"github.com/alanraison/fleps/pkg/chat"
	"github.com/alanraison/fleps/pkg/model"
)

type handler struct {
	fixtureRepository model.FixtureRepository
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event, err := chat.NewAppCommandEvent(r.Body)
	if err != nil {
		fmt.Printf("Error creating app command event: %v\n", err)
		w.WriteHeader(500)
		w.Write([]byte(`{}`))
		return
	}
	switch event.CommandId {
	case chat.MakePredictionsCommand:
		// Handle the MakePredictionsCommand
	default:
		// Handle unknown commands
	}
	w.Write([]byte("{}"))
}
