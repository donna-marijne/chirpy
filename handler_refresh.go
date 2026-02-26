package main

import (
	"log"
	"net/http"
	"time"

	"github.com/donnamarijne/chirpy/internal/auth"
	"github.com/donnamarijne/chirpy/internal/database"
)

type RefreshResponse struct {
	Token string `json:"token"`
}

func (c *apiConfig) handlerRefresh(writer http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		sendErrorResponse(writer, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	refreshToken, err := c.dbQueries.GetRefreshToken(req.Context(), bearerToken)
	if err != nil {
		sendErrorResponse(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !isRefreshTokenValid(refreshToken) {
		sendErrorResponse(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(
		refreshToken.UserID,
		c.secret,
		time.Hour,
	)
	if err != nil {
		log.Printf("Error making a JWT: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	res := RefreshResponse{
		Token: token,
	}
	sendResponse(writer, res, http.StatusOK)
}

func isRefreshTokenValid(refreshToken database.RefreshToken) bool {
	return !refreshToken.RevokedAt.Valid && refreshToken.ExpiresAt.After(time.Now())
}
