package cmd

import (
	"database/sql"
	"fmt"
	"os"

	sqlpkg "github.com/alanraison/predictions/pkg/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var (
	dbPath   string
	db       *sql.DB
	database *sqlpkg.Database
	rootCmd  = &cobra.Command{
		Use:   "predictions",
		Short: "Predictions is a CLI tool for making football predictions",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if dbPath == "" {
				return fmt.Errorf("db-path flag is required")
			}
			db, err := sql.Open("sqlite3", dbPath)
			if err != nil {
				return err
			}
			database = sqlpkg.NewDatabase(db)

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
	rootCmd.AddCommand(fixtures)
	rootCmd.AddCommand(players)
	rootCmd.AddCommand(predictions)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
