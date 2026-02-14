package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting Chirpy...")

	handler := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}

	log.Println("Chirpy started!")
}
