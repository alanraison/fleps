package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/alanraison/predictions/pkg/csv"
	sqlpkg "github.com/alanraison/predictions/pkg/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"

	"github.com/alanraison/predictions/pkg/model"
)

type appContext struct {
	db              *sql.DB
	teamRepo        model.TeamRepository
	fixturesRepo    model.FixtureRepository
	predictionsRepo model.PredictionRepository
	resultsRepo     model.ResultRepository
	playerRepo      model.PlayerRepository
}

func newAppContext(dbUrl string) (*appContext, error) {
	db, err := sqlpkg.OpenDatabase(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	database := sqlpkg.NewDatabase(db)

	return &appContext{
		db:              db,
		teamRepo:        database,
		fixturesRepo:    database,
		predictionsRepo: database,
		resultsRepo:     database,
		playerRepo:      database,
	}, nil
}

func getAppContext(cmd *cobra.Command) *appContext {
	v := cmd.Context().Value("appContext")
	if v == nil {
		panic("appContext not found in command context")
	}
	return v.(*appContext)
}

var (
	dbPath  string
	c       *csv.Csv
	rootCmd = &cobra.Command{
		Use:   "predictions",
		Short: "Predictions is a CLI tool for making football predictions",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if dbPath == "" {
				return fmt.Errorf("db-path flag is required")
			}
			ctx, err := newAppContext(dbPath)
			if err != nil {
				return fmt.Errorf("creating app context: %w", err)
			}
			cmd.SetContext(context.WithValue(cmd.Context(), "appContext", ctx))
			c = csv.NewCsv(ctx.teamRepo, ctx.fixturesRepo)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			ctx := getAppContext(cmd)
			if ctx.db != nil {
				if err := ctx.db.Close(); err != nil {
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
