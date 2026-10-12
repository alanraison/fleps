package main

import (
	"encoding/json"
	"fmt"

	"github.com/alanraison/fleps/pkg/chat"
)

func (h *handler) handleMakePredictions(actionUrl, email string) ([]byte, error) {
	roundID, err := h.fixtureRepository.GetLatestRoundWithNoResults()
	if err != nil {
		return nil, fmt.Errorf("finding round to predict: %w", err)
	}
	fixtures, err := h.fixtureRepository.ListFixtures(roundID)
	if err != nil {
		return nil, fmt.Errorf("listing fixtures for round %v: %w", roundID, err)
	}
	predictions, err := h.predictionRepository.GetPredictionsForPlayer(email, roundID)
	if err != nil {
		return nil, fmt.Errorf("getting predictions for %q in round %v: %w", email, roundID, err)
	}
	dialog := chat.PredictionsDialogBuilder(h.teamRepository, actionUrl).
		MakePredictionsDialog(roundID, fixtures, predictions)
	bytes, err := json.Marshal(dialog)
	if err != nil {
		return nil, fmt.Errorf("encoding predictions dialog: %w", err)
	}
	return bytes, nil
}
