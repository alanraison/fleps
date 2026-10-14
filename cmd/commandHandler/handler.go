package main

import (
	"fmt"
	"net/http"

	"github.com/alanraison/fleps/pkg/chat"
	"github.com/alanraison/fleps/pkg/model"
)

type handler struct {
	fixtureRepository    model.FixtureRepository
	teamRepository       model.TeamRepository
	predictionRepository model.PredictionRepository
	actionUrl            string
	results              resultsSource
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
		bytes, err := h.handleMakePredictions(h.actionUrl, event.UserEmail)
		if err != nil {
			fmt.Printf("Error handling MakePredictionsCommand: %v\n", err)
			w.WriteHeader(500)
			w.Write([]byte(`{}`))
			return
		}
		fmt.Printf("Response bytes: %s\n", string(bytes))
		w.Write(bytes)
	case chat.ShowResultsCommand:
		bytes, err := h.handleShowResults(r.Context())
		if err != nil {
			fmt.Printf("Error handling ShowResultsCommand: %v\n", err)
			w.WriteHeader(500)
			w.Write([]byte(`{}`))
			return
		}
		w.Write(bytes)
	default:
		fmt.Printf("CommandId %v not handled", event.CommandId)
		w.WriteHeader(400)
		w.Write([]byte("{}"))
		// Handle unknown commands
	}
}
