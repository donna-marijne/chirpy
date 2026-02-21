package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/donnamarijne/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	port := "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)

	config := apiConfig{
		dbQueries: database.New(db),
		platform:  platform,
	}

	mux := http.NewServeMux()

	// File server
	handlerFileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", middlewareLog(config.middlewareMetricsInc(handlerFileServer)))

	// Metrics
	mux.HandleFunc("GET /admin/metrics", config.handlerMetrics)

	// Reset
	mux.HandleFunc("POST /admin/reset", config.handlerReset)

	// Chirps
	mux.HandleFunc("POST /api/chirps", config.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", config.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", config.handlerChirpsGetOne)

	// Health
	mux.HandleFunc("GET /api/healthz", handlerHealth)

	// Login
	mux.HandleFunc("POST /api/login", config.handlerLogin)

	// Users
	mux.HandleFunc("POST /api/users", config.handlerUserCreate)

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
