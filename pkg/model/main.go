package model

import (
	"errors"
	"time"
)

var (
	UnknownTeamErr    = errors.New("unknown team")
	UnknownFixtureErr = errors.New("unknown fixture")
)

type TeamKey string
type Team struct {
	Key       TeamKey
	FullName  string
	ShortName string
}
type Fixture struct {
	HomeTeam TeamKey
	AwayTeam TeamKey
	Date     time.Time
}
type Result struct {
	Fixture
	HomeScore int
	AwayScore int
}
type TeamRepository interface {
	FindTeamByKey(key TeamKey) (*Team, error)
}
type FixtureRepository interface {
	// AddFixtures stores the given fixtures in the repository.
	AddFixtures(fixtures []Fixture) error
	// ListFixtures returns a list of fixtures between the given dates and for the given teams.
	// If no teams are provided, all fixtures are returned.
	ListFixtures(fromDate time.Time, toDate time.Time, teams []string) ([]Fixture, error)
	AddResult(homeTeam, awayTeam TeamKey, date time.Time, homeGoals, awayGoals int) error
	ListResults(fromDate time.Time, toDate time.Time, teams []string) ([]Result, error)
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
	HomeGoals int
	AwayGoals int
}
type PredictionRepository interface {
	AddPrediction(player string, homeTeam, awayTeam TeamKey, date time.Time, homeGoals, awayGoals int) error
	ListPredictions(from, to time.Time) ([]Prediction, error)
}
