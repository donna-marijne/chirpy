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

type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (c *apiConfig) handlerUserCreate(writer http.ResponseWriter, req *http.Request) {
	body := UserCreateRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		sendErrorResponse(writer, "Malformed request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	createParams := database.CreateUserParams{
		Email:          body.Email,
		HashedPassword: hashedPassword,
	}
	user, err := c.dbQueries.CreateUser(req.Context(), createParams)
	if err != nil {
		log.Printf("Error from CreateUser: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	res := UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	sendResponse(writer, res, http.StatusCreated)
}
