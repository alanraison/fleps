package main

import (
	"encoding/json"
	"net/http"

	chatmodel "github.com/alanraison/fleps/pkg/chat/model"
	"github.com/alanraison/fleps/pkg/model"
)

type handler struct {
	playerRepository model.PlayerRepository
}

type action struct {
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := chatmodel.RenderActions{
		Action: &chatmodel.RenderAction{
			Navigations: []chatmodel.Navigation{
				{
					EndNavigation: &chatmodel.EndNavigation{
						Action: chatmodel.EndNavigationCloseDialog,
					},
				},
			},
			Notification: &chatmodel.Notification{
				Text: "Predictions submitted",
			},
		},
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
