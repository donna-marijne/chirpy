package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthHandler)
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting Chirpy on port %s...", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
