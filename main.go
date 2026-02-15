package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, req *http.Request) {
			c.fileserverHits.Add(1)
			next.ServeHTTP(writer, req)
		},
	)
}

func (c *apiConfig) handlerMetrics(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)

	body := fmt.Sprintf("Hits: %d", c.fileserverHits.Load())
	writer.Write([]byte(body))
}

func (c *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {
	c.fileserverHits.Store(0)

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)

	writer.Write([]byte("OK"))
}

func main() {
	port := "8080"

	config := apiConfig{}

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/healthz", healthHandler)

	// File server
	handlerFileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", middlewareLog(config.middlewareMetricsInc(handlerFileServer)))

	// Metrics
	mux.HandleFunc("/metrics", config.handlerMetrics)

	// Reset
	mux.HandleFunc("/reset", config.handlerReset)

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
