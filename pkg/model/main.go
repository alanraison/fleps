package model

import "time"

type Team struct {
	Key       string
	FullName  string
	ShortName string
}
type Fixture struct {
	HomeTeam Team
	AwayTeam Team
	Date     time.Time
}
type FixtureFetcher interface {
	FetchFixtures(dateRange string, teams []string) []Fixture
}
type TeamRepository interface {
	FindByKey(key string) (Team, error)
}
type FixtureRepository interface {
	AddFixtures(fixtures []Fixture) error
}
