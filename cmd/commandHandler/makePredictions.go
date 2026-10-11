package main

import (
	"encoding/json"
	"time"

	"github.com/alanraison/fleps/pkg/chat"
	"github.com/alanraison/fleps/pkg/model"
)

func (h *handler) handleMakePredictions(actionUrl string) ([]byte, error) {
	dialog := chat.PredictionsDialogBuilder(h.teamRepository, actionUrl).MakePredictionsDialog(
		"alan.raison@newredo.com",
		[]model.Fixture{
			{
				FixtureKey: model.FixtureKey{
					RoundID:  "26/27R5",
					HomeTeam: "ARS",
					AwayTeam: "LEE",
				},
				Date: time.Date(2026, time.October, 10, 11, 30, 0, 0, time.UTC),
			},
			{
				FixtureKey: model.FixtureKey{
					RoundID:  "26/27R5",
					HomeTeam: "WBA",
					AwayTeam: "BIR",
				},
				Date: time.Date(2026, time.October, 10, 11, 30, 0, 0, time.UTC),
			},
			{
				FixtureKey: model.FixtureKey{
					RoundID:  "26/27R5",
					HomeTeam: "CHE",
					AwayTeam: "BOU",
				},
				Date: time.Date(2026, time.October, 10, 14, 00, 0, 0, time.UTC),
			},
			{
				FixtureKey: model.FixtureKey{
					RoundID:  "26/27R5",
					HomeTeam: "BAR",
					AwayTeam: "STP",
				},
				Date: time.Date(2026, time.October, 10, 14, 00, 0, 0, time.UTC),
			},
			{
				FixtureKey: model.FixtureKey{
					RoundID:  "26/27R5",
					HomeTeam: "MUN",
					AwayTeam: "TOT",
				},
				Date: time.Date(2026, time.October, 10, 16, 30, 0, 0, time.UTC),
			},
			{
				FixtureKey: model.FixtureKey{
					RoundID:  "26/27R5",
					HomeTeam: "LIV",
					AwayTeam: "MCI",
				},
				Date: time.Date(2026, time.October, 11, 15, 30, 0, 0, time.UTC),
			},
		},
		nil,
	)
	bytes, err := json.Marshal(dialog)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
