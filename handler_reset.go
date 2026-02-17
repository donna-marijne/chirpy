package main

import (
	"log"
	"net/http"
)

func (c *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {
	if c.platform != "dev" {
		log.Printf("Attempt to reset on a non-dev environment blocked!")
		sendErrorResponse(writer, "Forbidden", http.StatusForbidden)
		return
	}

	c.fileserverHits.Store(0)

	err := c.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		log.Printf("Error from DeleteUsers: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

	writer.Write([]byte("OK"))
}
