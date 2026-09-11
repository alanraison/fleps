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

func OpenDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database %q: %w", dbPath, err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enabling foreign keys for %q: %w", dbPath, err)
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
