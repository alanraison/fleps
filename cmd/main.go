package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/alanraison/predictions/pkg/csv"
	sqlpkg "github.com/alanraison/predictions/pkg/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var (
	dbPath   string
	db       *sql.DB
	database *sqlpkg.Database
	c        *csv.Csv
	rootCmd  = &cobra.Command{
		Use:   "predictions",
		Short: "Predictions is a CLI tool for making football predictions",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if dbPath == "" {
				return fmt.Errorf("db-path flag is required")
			}
			db, err := sqlpkg.OpenDatabase(dbPath)
			if err != nil {
				return fmt.Errorf("opening database: %w", err)
			}
			database = sqlpkg.NewDatabase(db)
			c = csv.NewCsv(database, database)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if db != nil {
				if err := db.Close(); err != nil {
					return fmt.Errorf("closing database: %w", err)
				}
			}
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db-path", "football.db", "Path to the SQLite database file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
