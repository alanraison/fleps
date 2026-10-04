package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"
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

	addedToSpaceHandler := &handler{
		playerRepository: sqlpkg.NewPlayerRepository(db),
	}
	http.Handle("/", addedToSpaceHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
