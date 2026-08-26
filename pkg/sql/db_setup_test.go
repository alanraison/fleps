package sql

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(tb testing.TB) (func(tb testing.TB), *sql.DB) {
	tb.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		tb.Fatal(err)
	}

	// Keep one connection for SQLite in-memory tests so schema and data are shared.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		tb.Fatal(err)
	}

	if err := ApplyDefaultSchema(db); err != nil {
		_ = db.Close()
		tb.Fatal(err)
	}

	return func(tb testing.TB) {
		tb.Helper()
		_ = db.Close()
	}, db
}

