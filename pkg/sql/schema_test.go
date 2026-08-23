package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func ApplySchema(db *sql.DB, schemaDir string) error {
	entries, err := os.ReadDir(schemaDir)
	if err != nil {
		return fmt.Errorf("reading schema directory %q: %w", schemaDir, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("starting schema transaction: %w", err)
	}
	defer tx.Rollback()

	for _, name := range files {
		content, err := os.ReadFile(filepath.Join(schemaDir, name))
		if err != nil {
			return fmt.Errorf("reading schema file %q: %w", name, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			return fmt.Errorf("applying schema file %q: %w", name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing schema transaction: %w", err)
	}

	return nil
}

func ApplyDefaultSchema(db *sql.DB) error {
	dir, err := defaultSchemaDir()
	if err != nil {
		return err
	}

	return ApplySchema(db, dir)
}

func defaultSchemaDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("locating schema helper source file")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "sql")), nil
}
