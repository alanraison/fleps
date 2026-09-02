package sql

import (
	"database/sql"
	"testing"
	"time"

	"github.com/alanraison/predictions/pkg/model"
)

func setupDefaultFixtureData(tb testing.TB, db *sql.DB) {
	tb.Helper()

	if _, err := db.Exec(`
		INSERT INTO fixtures (home_team, away_team, date_time) VALUES
		('LEE', 'MUN', '2023-10-01 15:00:00'),
		('ARS', 'CHE', '2023-10-02 16:00:00');
	`); err != nil {
		tb.Fatalf("failed to insert test data: %v", err)
	}
}

func TestFixtureRepositoryReturnsEmptyListWhenNoFixtures(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	from, _ := time.Parse("2006-01-02", "2023-10-01")
	to, _ := time.Parse("2006-01-02", "2023-10-02")

	fixtures, err := repo.ListFixtures(from, to, []string{"LEE", "MUN"})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 0 {
		t.Fatalf("expected 0 fixtures, got %d", len(fixtures))
	}
}

func TestFixtureRepositoryListFixturesWithHomeTeamFilter(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultFixtureData(t, db)

	from, _ := time.Parse("2006-01-02", "2023-10-01")
	to, _ := time.Parse("2006-01-02", "2023-10-02")

	fixtures, err := repo.ListFixtures(from, to, []string{"LEE"})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 1 {
		t.Fatalf("expected 1 fixture, got %d", len(fixtures))
	}

	fixture := fixtures[0]
	if got, want := fixture.HomeTeam, model.TeamKey("LEE"); got != want {
		t.Fatalf("HomeTeam.Key = %q, want %q", got, want)
	}
	if got, want := fixture.AwayTeam, model.TeamKey("MUN"); got != want {
		t.Fatalf("AwayTeam.Key = %q, want %q", got, want)
	}
	if got, want := fixture.Date.Format("2006-01-02"), "2023-10-01"; got != want {
		t.Fatalf("Date = %q, want %q", got, want)
	}
}

func TestFixtureRepositoryListFixturesWithAwayTeamFilter(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultFixtureData(t, db)

	from, _ := time.Parse("2006-01-02", "2023-10-01")
	to, _ := time.Parse("2006-01-02 15:04:05", "2023-10-02 23:59:59")

	fixtures, err := repo.ListFixtures(from, to, []string{"CHE"})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 1 {
		t.Fatalf("expected 1 fixture, got %d", len(fixtures))
	}

	fixture := fixtures[0]
	if got, want := fixture.HomeTeam, model.TeamKey("ARS"); got != want {
		t.Fatalf("HomeTeam.Key = %q, want %q", got, want)
	}
	if got, want := fixture.AwayTeam, model.TeamKey("CHE"); got != want {
		t.Fatalf("AwayTeam.Key = %q, want %q", got, want)
	}
	if got, want := fixture.Date.Format("2006-01-02"), "2023-10-02"; got != want {
		t.Fatalf("Date = %q, want %q", got, want)
	}
}

func TestFixtureRepositoryListFixturesWithNoTeamFilter(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultFixtureData(t, db)

	from, _ := time.Parse("2006-01-02", "2023-10-01")
	to, _ := time.Parse("2006-01-02 15:04:05", "2023-10-02 23:59:59")

	fixtures, err := repo.ListFixtures(from, to, []string{})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 2 {
		t.Fatalf("expected 2 fixtures, got %d", len(fixtures))
	}
}

func TestFixtureRepositoryListFixturesWithHomeAndAwayTeamFilter(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultFixtureData(t, db)

	from, _ := time.Parse("2006-01-02", "2023-10-01")
	to, _ := time.Parse("2006-01-02 15:04:05", "2023-10-02 23:59:59")

	fixtures, err := repo.ListFixtures(from, to, []string{"LEE", "MUN"})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 1 {
		t.Fatalf("expected 1 fixture, got %d", len(fixtures))
	}

	fixture := fixtures[0]
	if got, want := fixture.HomeTeam, model.TeamKey("LEE"); got != want {
		t.Fatalf("HomeTeam.Key = %q, want %q", got, want)
	}
	if got, want := fixture.AwayTeam, model.TeamKey("MUN"); got != want {
		t.Fatalf("AwayTeam.Key = %q, want %q", got, want)
	}
	if got, want := fixture.Date.Format("2006-01-02"), "2023-10-01"; got != want {
		t.Fatalf("Date = %q, want %q", got, want)
	}
}

func TestFixtureRepositoryListFixturesWithNoMatchingTeams(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultFixtureData(t, db)

	from, _ := time.Parse("2006-01-02", "2023-10-01")
	to, _ := time.Parse("2006-01-02 15:04:05", "2023-10-02 23:59:59")

	fixtures, err := repo.ListFixtures(from, to, []string{"XYZ"})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 0 {
		t.Fatalf("expected 0 fixtures, got %d", len(fixtures))
	}
}

func TestFixtureRepositoryListFixturesWithMatchedDateRange(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultFixtureData(t, db)

	from, _ := time.Parse("2006-01-02", "2023-10-02")
	to, _ := time.Parse("2006-01-02 15:04:05", "2023-10-03 23:59:59")

	fixtures, err := repo.ListFixtures(from, to, []string{})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 1 {
		t.Fatalf("expected 1 fixture, got %d", len(fixtures))
	}

	fixture := fixtures[0]
	if got, want := fixture.HomeTeam, model.TeamKey("ARS"); got != want {
		t.Fatalf("HomeTeam = %q, want %q", got, want)
	}
	if got, want := fixture.AwayTeam, model.TeamKey("CHE"); got != want {
		t.Fatalf("AwayTeam = %q, want %q", got, want)
	}
	if got, want := fixture.Date.Format("2006-01-02"), "2023-10-02"; got != want {
		t.Fatalf("Date = %q, want %q", got, want)
	}
}

func TestFixtureRepositoryListFixturesWithNoMatchingDateRange(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultFixtureData(t, db)

	from, _ := time.Parse("2006-01-02", "2023-10-03")
	to, _ := time.Parse("2006-01-02 15:04:05", "2023-10-04 23:59:59")

	fixtures, err := repo.ListFixtures(from, to, []string{})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 0 {
		t.Fatalf("expected 0 fixtures, got %d", len(fixtures))
	}
}

func TestFixtureRepositoryAddFixtures(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)

	fixturesToAdd := []model.Fixture{
		{HomeTeam: model.TeamKey("LEE"), AwayTeam: model.TeamKey("MUN"), Date: time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC)},
		{HomeTeam: model.TeamKey("ARS"), AwayTeam: model.TeamKey("CHE"), Date: time.Date(2023, 10, 2, 16, 0, 0, 0, time.UTC)},
	}

	err := repo.AddFixtures(fixturesToAdd)
	if err != nil {
		t.Fatalf("AddFixtures returned error: %v", err)
	}

	from, _ := time.Parse("2006-01-02", "2023-10-01")
	to, _ := time.Parse("2006-01-02 15:04:05", "2023-10-02 23:59:59")

	fixtures, err := repo.ListFixtures(from, to, []string{})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}

	if len(fixtures) != 2 {
		t.Fatalf("expected 2 fixtures, got %d", len(fixtures))
	}
}

func TestFixtureRepositoryShouldFailToAddUnknownTeam(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)

	fixturesToAdd := []model.Fixture{
		{HomeTeam: model.TeamKey("LEE"), AwayTeam: model.TeamKey("MUN"), Date: time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC)},
		{HomeTeam: model.TeamKey("BBB"), AwayTeam: model.TeamKey("CCC"), Date: time.Date(2023, 10, 2, 16, 0, 0, 0, time.UTC)},
	}

	err := repo.AddFixtures(fixturesToAdd)
	if err == nil {
		t.Fatalf("expected error when adding fixture with unknown team, got nil")
	}

	from, _ := time.Parse("2006-01-02", "2023-10-01")
	to, _ := time.Parse("2006-01-02 15:04:05", "2023-10-02 23:59:59")
	fixtures, err := repo.ListFixtures(from, to, []string{})
	if err != nil {
		t.Fatalf("ListFixtures returned error: %v", err)
	}
	if len(fixtures) != 0 {
		t.Fatalf("expected 0 fixtures after rollback, got %d", len(fixtures))
	}
}

func TestFixtureRepositoryAddResultShouldAddResult(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)
	setupDefaultFixtureData(t, db)

	matchDate := time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC)
	t.Logf("Match date: %v", matchDate)

	err := repo.AddResult("LEE", "MUN", matchDate, 2, 1)
	if err != nil {
		t.Fatalf("AddResult returned error: %v", err)
	}

	from, _ := time.Parse("2006-01-02", "2023-10-01")
	to, _ := time.Parse("2006-01-02 15:04:05", "2023-10-01 23:59:59")
	results, err := repo.ListResults(from, to, []string{})
	if err != nil {
		t.Fatalf("ListResults returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	fixture := results[0]
	if fixture.HomeTeam != model.TeamKey("LEE") || fixture.AwayTeam != model.TeamKey("MUN") {
		t.Fatalf("unexpected fixture teams: got %s vs %s", fixture.HomeTeam, fixture.AwayTeam)
	}
	if fixture.Date != matchDate {
		t.Fatalf("unexpected fixture date: got %v", fixture.Date)
	}
	if fixture.HomeScore != 2 || fixture.AwayScore != 1 {
		t.Fatalf("unexpected fixture score: got %d - %d", fixture.HomeScore, fixture.AwayScore)
	}
}

func TestFixtureRepositoryAddResultShouldFailForUnknownFixture(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	repo := NewFixtureRepository(db)

	setupDefaultTeamData(t, db)

	matchDate := time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC)
	err := repo.AddResult("LEE", "MUN", matchDate, 2, 1)
	if err == nil {
		t.Fatalf("expected error when adding result for unknown fixture, got nil")
	}
}
