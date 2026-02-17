package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/donnamarijne/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
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
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)

	body := fmt.Sprintf(
		`<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>`,
		c.fileserverHits.Load(),
	)
	writer.Write([]byte(body))
}

func (c *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {
	c.fileserverHits.Store(0)

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)

	writer.Write([]byte("OK"))
}
