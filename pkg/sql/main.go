package sql

import (
	"database/sql"
	"fmt"
)

type Database struct {
	*fixtureRepository
	*playerRepository
	*predictionRepository
	*resultRepository
	*teamRepository
}

func OpenDatabase(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("opening database %q: %w", url, err)
	}
	return db, nil
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		fixtureRepository:    &fixtureRepository{db: db},
		playerRepository:     &playerRepository{db: db},
		predictionRepository: &predictionRepository{db: db},
		resultRepository:     &resultRepository{db: db},
		teamRepository:       &teamRepository{db: db},
	}
}
