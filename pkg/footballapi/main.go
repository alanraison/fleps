package footballapi

import "github.com/alanraison/predictions/pkg/model"

type FootballAPI struct {
	APIKey string
}

func (f *FootballAPI) FetchFixtures(dateRange string, teams []string) []model.Fixture {
	// Implementation to fetch fixtures from the football API using the APIKey
	// This is a placeholder for the actual implementation
	return []model.Fixture{}
}
