package sql

import (
	"errors"
	"testing"

	"github.com/alanraison/fleps/pkg/model"
)

func TestTeamLogoRepositoryFindLogoSlugByTeamKey(t *testing.T) {
	teardown, db := setupTestDB(t)
	defer teardown(t)
	setupDefaultTeamData(t, db)
	if _, err := db.Exec(`INSERT INTO team_logos (team_key, logo_slug) VALUES ('LEE', 'leeds-united')`); err != nil {
		t.Fatalf("failed to insert logo: %v", err)
	}

	repo, err := NewTeamLogoRepository(db)
	if err != nil {
		t.Fatalf("NewTeamLogoRepository: %v", err)
	}
	slug, err := repo.FindLogoSlugByTeamKey("LEE")
	if err != nil || slug != "leeds-united" {
		t.Errorf("got %q, %v; want leeds-united", slug, err)
	}
	if _, err := repo.FindLogoSlugByTeamKey("ARS"); !errors.Is(err, model.UnknownLogoErr) {
		t.Errorf("expected UnknownLogoErr, got %v", err)
	}
}
