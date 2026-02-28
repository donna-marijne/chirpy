package main

import (
	"context"
	"log"
	"net/http"

	"github.com/donnamarijne/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (c *apiConfig) middlewareAuthenticate(nextFunc http.HandlerFunc) http.Handler {
	next := http.HandlerFunc(nextFunc)

	return http.HandlerFunc(
		func(writer http.ResponseWriter, req *http.Request) {
			userID, err := c.getAuthenticatedUserID(req)
			if err != nil {
				log.Printf("Authorization failed: %v", err)
				sendErrorResponse(writer, "Unauthorized", http.StatusUnauthorized)
				return
			}

			newContext := context.WithValue(req.Context(), "UserID", userID)
			newReq := req.WithContext(newContext)
			next.ServeHTTP(writer, newReq)
		},
	)
}

func (c *apiConfig) getAuthenticatedUserID(req *http.Request) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := auth.ValidateJWT(token, c.secret)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}
