package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"

	handler := http.NewServeMux()

	handler.Handle("/", http.FileServer(http.Dir(".")))

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Starting Chirpy on port %s...", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
