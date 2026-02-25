package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/donnamarijne/chirpy/internal/auth"
	"github.com/donnamarijne/chirpy/internal/database"
	"github.com/google/uuid"
)

type ChirpCreateRequest struct {
	Body string `json:"body"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (c *apiConfig) handlerChirpsCreate(writer http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		sendErrorResponse(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, c.secret)
	if err != nil {
		sendErrorResponse(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chirpCreate := ChirpCreateRequest{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&chirpCreate)
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

	chirp, err := c.dbQueries.CreateChirp(
		req.Context(),
		database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: userID,
		},
	)
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

func (c *apiConfig) handlerChirpsGet(writer http.ResponseWriter, req *http.Request) {
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

func (c *apiConfig) handlerChirpsGetOne(writer http.ResponseWriter, req *http.Request) {
	chirpIDStr := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		sendErrorResponse(writer, "Invalid chirpID format", http.StatusBadRequest)
		return
	}

	chirp, err := c.dbQueries.GetChirp(req.Context(), chirpID)
	if err != nil {
		sendErrorResponse(writer, "Chirp not found", http.StatusNotFound)
		return
	}

	res := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	sendResponse(writer, res, http.StatusOK)
}
