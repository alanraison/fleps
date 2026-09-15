package csv

import (
	"github.com/alanraison/fleps/pkg/model"
)

type Csv struct {
	teamRepo    model.TeamRepository
	fixtureRepo model.FixtureRepository
}

type PredictionRow struct {
	HomeTeam  string
	AwayTeam  string
	HomeGoals int
	AwayGoals int
}

func NewCsv(teamRepo model.TeamRepository, fixtureRepo model.FixtureRepository) *Csv {
	return &Csv{
		teamRepo:    teamRepo,
		fixtureRepo: fixtureRepo,
	}
}
