package main

import (
	"net/http"

	"github.com/donnamarijne/chirpy/internal/auth"
)

func (c *apiConfig) handlerRevoke(writer http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		sendErrorResponse(writer, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	_, err = c.dbQueries.RevokeRefreshToken(req.Context(), bearerToken)
	if err != nil {
		sendErrorResponse(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
