package main

import (
	"fmt"
	"io"
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
	body, err := io.ReadAll(io.LimitReader(r.Body, 10240))
	if err != nil {
		fmt.Printf("Error reading request body: %v\n", err)
		return
	}
	fmt.Printf("Received request %+v\n", string(body))
	w.Write([]byte("{}"))
}
