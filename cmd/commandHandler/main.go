package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", helloWorld)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Command Handler"))
}