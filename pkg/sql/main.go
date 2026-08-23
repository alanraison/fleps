package sql

import (
	"database/sql"
)

type Database struct {
	TeamRepo    *TeamRepository
	FixtureRepo *FixtureRepository
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		TeamRepo:    &TeamRepository{db: db},
		FixtureRepo: &FixtureRepository{db: db},
	}
}
