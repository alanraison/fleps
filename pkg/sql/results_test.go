package sql

import (
	"errors"
	"testing"

	"github.com/alanraison/predictions/pkg/model"
)

func TestFixtureRepositoryAddResultShouldAddResult(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewResultRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)

	err := repo.AddResult(model.FixtureKey{
		RoundID:  model.RoundID("R1"),
		HomeTeam: model.TeamKey("LEE"),
		AwayTeam: model.TeamKey("MUN"),
	}, 2, 1)
	if err != nil {
		t.Fatalf("AddResult returned error: %v", err)
	}

	results, err := repo.ListResults(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListResults returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	fixture := results[0]
	if fixture.RoundID != model.RoundID("R1") {
		t.Fatalf("unexpected fixture round: got %s", fixture.RoundID)
	}
	if fixture.HomeTeam != model.TeamKey("LEE") || fixture.AwayTeam != model.TeamKey("MUN") {
		t.Fatalf("unexpected fixture teams: got %s vs %s", fixture.HomeTeam, fixture.AwayTeam)
	}
	if fixture.HomeGoals != 2 || fixture.AwayGoals != 1 {
		t.Fatalf("unexpected fixture score: got %d - %d", fixture.HomeGoals, fixture.AwayGoals)
	}
}

func TestResultRepositoryAddResultShouldFailForUnknownFixture(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewResultRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)

	err := repo.AddResult(model.FixtureKey{
		RoundID:  model.RoundID("R1"),
		HomeTeam: model.TeamKey("LEE"),
		AwayTeam: model.TeamKey("MUN"),
	}, 2, 1)
	if err == nil {
		t.Fatalf("expected error when adding result for unknown fixture, got nil")
	}
	if !errors.Is(err, model.UnknownFixtureErr) {
		t.Fatalf("expected UnknownFixtureErr, got %v", err)
	}
}

func TestResultRepositoryListResultsShouldReturnEmptyWhenNoResultsExist(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewResultRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)

	results, err := repo.ListResults(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListResults returned error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestResultRepositoryListResultsShouldReturnResultsForGivenRound(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewResultRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)

	err := repo.AddResult(model.FixtureKey{
		RoundID:  model.RoundID("R1"),
		HomeTeam: model.TeamKey("LEE"),
		AwayTeam: model.TeamKey("MUN"),
	}, 2, 1)
	if err != nil {
		t.Fatalf("AddResult returned error: %v", err)
	}
	err = repo.AddResult(model.FixtureKey{
		RoundID:  model.RoundID("R2"),
		HomeTeam: model.TeamKey("LEE"),
		AwayTeam: model.TeamKey("MUN"),
	}, 3, 2)
	if err != nil {
		t.Fatalf("AddResult returned error: %v", err)
	}

	results, err := repo.ListResults(model.RoundID("R1"))
	if err != nil {
		t.Fatalf("ListResults returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	fixture := results[0]
	if fixture.RoundID != model.RoundID("R1") {
		t.Fatalf("unexpected fixture round: got %s", fixture.RoundID)
	}
	if fixture.HomeTeam != model.TeamKey("LEE") || fixture.AwayTeam != model.TeamKey("MUN") {
		t.Fatalf("unexpected fixture teams: got %s vs %s", fixture.HomeTeam, fixture.AwayTeam)
	}
	if fixture.HomeGoals != 2 || fixture.AwayGoals != 1 {
		t.Fatalf("unexpected fixture score: got %d - %d", fixture.HomeGoals, fixture.AwayGoals)
	}
}

func TestResultRepositoryListResultsShouldReturnErrorForUnknownRound(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewResultRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultRoundData(t, db)
	setupDefaultFixtureData(t, db)

	_, err := repo.ListResults(model.RoundID("R999"))
	if err == nil {
		t.Fatalf("expected error when listing results for unknown round, got nil")
	}
	if !errors.Is(err, model.UnknownRoundErr) {
		t.Fatalf("expected UnknownRoundErr, got %v", err)
	}
}
