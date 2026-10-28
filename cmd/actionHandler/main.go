package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"
	"github.com/alanraison/fleps/pkg/chat"
	"github.com/alanraison/fleps/pkg/logos"
	sqlpkg "github.com/alanraison/fleps/pkg/sql"
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
	h := &handler{
		playerRepository:     sqlpkg.NewPlayerRepository(db),
		predictionRepository: sqlpkg.NewPredictionRepository(db),
	}
	if space := os.Getenv("MAIN_SPACE"); space != "" {
		teams, err := sqlpkg.NewTeamRepository(db)
		if err != nil {
			fmt.Printf("Error loading teams: %v\n", err)
			return
		}
		teamLogos, err := sqlpkg.NewTeamLogoRepository(db)
		if err != nil {
			fmt.Printf("Error loading team logos: %v\n", err)
			return
		}
		n, err := chat.NewNotifier(context.Background(), space, teams, teamLogos, logos.NewResolver())
		if err != nil {
			fmt.Printf("Error creating notifier: %v\n", err)
			return
		}
		h.notifier = n
	}
	http.Handle("/", h)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
