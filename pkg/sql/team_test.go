package sql

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestTeamRepositoryFindByKey(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close database: %v", err)
		}
	})

	if err := ApplyDefaultSchema(db); err != nil {
		t.Fatalf("failed to apply schema: %v", err)
	}

	if _, err := db.Exec(`
		INSERT INTO teams (key, full_name, short_name, league) VALUES
		('LEE', 'Leeds United', 'Leeds', 1),
		('MUN', 'Manchester United', 'Man Utd', 1);
	`); err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	repo := &TeamRepository{db: db}

	t.Run("returns team for existing key", func(t *testing.T) {
		team, ok, err := repo.FindByKey("LEE")
		if err != nil {
			t.Fatalf("FindByKey returned error: %v", err)
		}
		if !ok {
			t.Fatal("expected team to be found")
		}
		if team == nil {
			t.Fatal("expected team, got nil")
		}
		if got, want := team.Key, "LEE"; got != want {
			t.Fatalf("Key = %q, want %q", got, want)
		}
		if got, want := team.FullName, "Leeds United"; got != want {
			t.Fatalf("FullName = %q, want %q", got, want)
		}
		if got, want := team.ShortName, "Leeds"; got != want {
			t.Fatalf("ShortName = %q, want %q", got, want)
		}
	})

	t.Run("returns no result for missing key", func(t *testing.T) {
		team, ok, err := repo.FindByKey("XYZ")
		if err != nil {
			t.Fatalf("FindByKey returned error: %v", err)
		}
		if ok {
			t.Fatal("expected team not to be found")
		}
		if team != nil {
			t.Fatalf("expected nil team, got %+v", team)
		}
	})
}
