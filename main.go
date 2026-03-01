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
	secret := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)

	config := apiConfig{
		dbQueries: database.New(db),
		platform:  platform,
		secret:    secret,
	}

	mux := http.NewServeMux()

	// File server
	handlerFileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", config.middlewareMetricsInc(handlerFileServer))

	// Metrics
	mux.HandleFunc("GET /admin/metrics", config.handlerMetrics)

	// Reset
	mux.HandleFunc("POST /admin/reset", config.handlerReset)

	// Chirps
	mux.HandleFunc("POST /api/chirps", config.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", config.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", config.handlerChirpsGetOne)
	mux.Handle(
		"DELETE /api/chirps/{chirpID}",
		config.middlewareAuthenticate(config.handlerChirpsDelete),
	)

	// Health
	mux.HandleFunc("GET /api/healthz", handlerHealth)

	// Login
	mux.HandleFunc("POST /api/login", config.handlerLogin)

	// Refresh
	mux.HandleFunc("POST /api/refresh", config.handlerRefresh)

	// Revoke
	mux.HandleFunc("POST /api/revoke", config.handlerRevoke)

	// Users
	mux.HandleFunc("POST /api/users", config.handlerUserCreate)
	mux.HandleFunc("PUT /api/users", config.handlerUserUpdate)

	// Webhooks
	mux.HandleFunc("POST /api/polka/webhooks", config.handlerPolkaWebhooks)

	server := http.Server{
		Addr:    ":" + port,
		Handler: middlewareLog(mux),
	}

	log.Printf("Starting Chirpy on port %s...", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
