package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"
	"github.com/alanraison/fleps/pkg/render"
	sqlpkg "github.com/alanraison/fleps/pkg/sql"
	"github.com/alanraison/fleps/pkg/storage"
)

func main() {
	cleanup, err := pgxv5.RegisterDriver("cloudsql-postgres", cloudsqlconn.WithIAMAuthN())
	if err != nil {
		fmt.Printf("Error registering driver: %v\n", err)
		return
	}
	defer cleanup()

	db, err := sql.Open("cloudsql-postgres", "")
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		return
	}
	teamRepo, err := sqlpkg.NewTeamRepository(db)
	if err != nil {
		fmt.Printf("Error creating team repository: %v\n", err)
		return
	}
	uploader, err := storage.NewImageUploader(context.Background(), os.Getenv("ASSETS_BUCKET"))
	if err != nil {
		fmt.Printf("Error creating image uploader: %v\n", err)
		return
	}
	http.Handle("/", &handler{
		results: &showResults{
			sources: render.Sources{
				Fixtures:    sqlpkg.NewFixtureRepository(db),
				Results:     sqlpkg.NewResultRepository(db),
				Predictions: sqlpkg.NewPredictionRepository(db),
				Players:     sqlpkg.NewPlayerRepository(db),
				Teams:       teamRepo,
				Scores:      sqlpkg.NewDBScoreService(db),
			},
			uploader: uploader,
		},
		fixtureRepository:    sqlpkg.NewFixtureRepository(db),
		teamRepository:       teamRepo,
		predictionRepository: sqlpkg.NewPredictionRepository(db),
		actionUrl:            os.Getenv("ACTION_URL"),
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
