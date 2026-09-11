package model

import (
	"errors"
	"time"
)

var (
	UnknownTeamErr    = errors.New("unknown team")
	UnknownFixtureErr = errors.New("unknown fixture")
	UnknownRoundErr   = errors.New("unknown round")
)

type TeamKey string
type SeasonID string
type RoundID string

type Team struct {
	Key       TeamKey
	FullName  string
	ShortName string
	League    int
}

// FixtureKey uniquely identifies a fixture within a round by the home and away teams.
type FixtureKey struct {
	RoundID  RoundID
	HomeTeam TeamKey
	AwayTeam TeamKey
}

// Fixture represents a football match between two teams, including the round and fixture date.
type Fixture struct {
	FixtureKey
	Date time.Time
}

// Result represents the outcome of a football match, including the number of goals scored by each team.
type Result struct {
	FixtureKey
	HomeGoals int
	AwayGoals int
}
type TeamRepository interface {
	AddTeam(newTeam *Team) error
	FindTeamByKey(key TeamKey) (*Team, error)
}
type FixtureRepository interface {
	GetLatestRoundWithNoResults() (RoundID, error)
	AddRound(RoundID, SeasonID) error
	// AddFixtures stores the given fixtures in the repository.
	AddFixtures(fixtures []Fixture) error
	// ListFixtures returns a list of fixtures for a round, or the latest one if none supplied.
	ListFixtures(RoundID) ([]Fixture, error)
}
type ResultRepository interface {
	GetLatestRoundWithResults() (RoundID, error)
	AddResult(fixture FixtureKey, homeGoals, awayGoals int) error
	ListResults(roundID RoundID) ([]Result, error)
}
type Player struct {
	Email string
	Name  string
}
type PlayerRepository interface {
	AddPlayer(email string, name string) error
	ListPlayers() ([]Player, error)
}

// Game represents a football match between two teams. It does not contain a fixture date or round,
// or score
type Game struct {
	HomeTeam TeamKey
	AwayTeam TeamKey
}
type Prediction struct {
	HomeGoals int
	AwayGoals int
}
type GamePredictions map[Game]Prediction
type PlayerPredictions map[string]GamePredictions
type PredictionRepository interface {
	AddPredictions(player string, roundID RoundID, predictions GamePredictions) error
	ListPredictions(roundID RoundID) (PlayerPredictions, error)
}
