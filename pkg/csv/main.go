package csv

import (
	"github.com/alanraison/predictions/pkg/model"
)

type Csv struct {
	teamRepo model.TeamRepository
}

type PredictionRow struct {
	Match       string
	Predictions map[string]string
}

func NewCsv(teamRepo model.TeamRepository, fixtureRepo model.FixtureRepository) *Csv {
	return &Csv{
		teamRepo: teamRepo,
	}
}
