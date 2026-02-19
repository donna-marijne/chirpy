package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/donnamarijne/chirpy/internal/database"
	"github.com/google/uuid"
)

type ChirpCreateRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (c *apiConfig) handlerChirpCreate(writer http.ResponseWriter, req *http.Request) {
	chirpCreate := ChirpCreateRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirpCreate)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusBadRequest)
		return
	}

	if len(chirpCreate.Body) > 140 {
		sendErrorResponse(writer, "Chirp is too long", http.StatusBadRequest)
		return
	}

	cleanedBody := removeProfanity(chirpCreate.Body)

	chirp, err := c.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: chirpCreate.UserID,
	})
	if err != nil {
		log.Printf("Error from CreateChirp: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	res := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	sendResponse(writer, res, http.StatusCreated)
}

func (c *apiConfig) handlerChirpGet(writer http.ResponseWriter, req *http.Request) {
	chirps, err := c.dbQueries.GetChirps(req.Context())
	if err != nil {
		log.Printf("Error from CreateChirp: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	res := make([]ChirpResponse, len(chirps))
	for i, chirp := range chirps {
		res[i] = ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	sendResponse(writer, res, http.StatusOK)
}
