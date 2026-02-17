package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/donnamarijne/chirpy/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	port := "8080"

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	config := apiConfig{
		dbQueries: database.New(db),
	}

	mux := http.NewServeMux()

	// File server
	handlerFileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", middlewareLog(config.middlewareMetricsInc(handlerFileServer)))

	// Metrics
	mux.HandleFunc("GET /admin/metrics", config.handlerMetrics)

	// Reset
	mux.HandleFunc("POST /admin/reset", config.handlerReset)

	// Health
	mux.HandleFunc("GET /api/healthz", handlerHealth)

	// Validate chirp
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting Chirpy on port %s...", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
