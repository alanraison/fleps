package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/alanraison/fleps/pkg/csv"
	sqlpkg "github.com/alanraison/fleps/pkg/sql"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"

	"github.com/alanraison/fleps/pkg/model"
)

type appContext struct {
	db             *sql.DB
	fixtureRepo    model.FixtureRepository
	playerRepo     model.PlayerRepository
	predictionRepo model.PredictionRepository
	resultRepo     model.ResultRepository
	teamRepo       model.TeamRepository
}

func newAppContext(dbUrl string) (*appContext, error) {
	db, err := sqlpkg.OpenDatabase(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	database := sqlpkg.NewDatabase(db)

	return &appContext{
		db:             db,
		fixtureRepo:    database,
		playerRepo:     database,
		predictionRepo: database,
		resultRepo:     database,
		teamRepo:       database,
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
	dbURL  string
	c       *csv.Csv
	rootCmd = &cobra.Command{
		Use:   "fleps",
		Short: "Fleps is a CLI tool for making football predictions",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if dbURL == "" {
				return fmt.Errorf("db-url flag is required")
			}
			ctx, err := newAppContext(dbURL)
			if err != nil {
				return fmt.Errorf("creating app context: %w", err)
			}
			cmd.SetContext(context.WithValue(cmd.Context(), "appContext", ctx))
			c = csv.NewCsv(ctx.teamRepo, ctx.fixtureRepo)

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
	rootCmd.PersistentFlags().StringVar(&dbURL, "db-url", "postgres://postgres@localhost/postgres?sslmode=disable", "Database connection URL")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
