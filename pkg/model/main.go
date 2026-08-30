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
type TeamRepository interface {
	FindTeamByKey(key string) (*Team, error)
}
type FixtureRepository interface {
	// AddFixtures stores the given fixtures in the repository.
	AddFixtures(fixtures []Fixture) error
	// ListFixtures returns a list of fixtures between the given dates and for the given teams.
	// If no teams are provided, all fixtures are returned.
	ListFixtures(fromDate time.Time, toDate time.Time, teams []string) ([]Fixture, error)
}
type Player struct {
	Email string
	Name  string
}
type PlayerRepository interface {
	AddPlayer(email string, name string) error
	ListPlayers() ([]Player, error)
}
type Prediction struct {
	Fixture
	Player    string
	HomeScore int
	AwayScore int
}
type PredictionRepository interface {
	AddPrediction(player, homeTeam, awayTeam string, date time.Time, homeScore, awayScore int) error
	ListPredictions(from, to time.Time) ([]Prediction, error)
}
