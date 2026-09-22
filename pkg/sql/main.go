package sql

import (
	"database/sql"
	"fmt"
)

func OpenDatabase(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("opening database %q: %w", url, err)
	}
	return db, nil
}
