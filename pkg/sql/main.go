package sql

import (
	"database/sql"
)

type Database struct {
	*teamRepository
	*fixtureRepository
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		teamRepository:    &teamRepository{db: db},
		fixtureRepository: &fixtureRepository{db: db},
	}
}
