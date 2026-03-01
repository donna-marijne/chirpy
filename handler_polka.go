package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/donnamarijne/chirpy/internal/database"
	"github.com/google/uuid"
)

type PolkaWebhooksRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (c *apiConfig) handlerPolkaWebhooks(writer http.ResponseWriter, req *http.Request) {
	body := PolkaWebhooksRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		sendErrorResponse(writer, "Malformatted request body", http.StatusBadRequest)
		return
	}
	if body.Event != "user.upgraded" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	dbParams := database.UpdateUserSetChirpyRedParams{
		ID:          body.Data.UserID,
		IsChirpyRed: true,
	}
	_, err = c.dbQueries.UpdateUserSetChirpyRed(req.Context(), dbParams)
	if err != nil {
		log.Printf("Error from UpdateUserSetChirpyRed: %v", err)
		sendErrorResponse(writer, "Not found", http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
