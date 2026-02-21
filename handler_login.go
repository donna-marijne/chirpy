package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/donnamarijne/chirpy/internal/auth"
	"github.com/donnamarijne/chirpy/internal/database"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *apiConfig) handlerLogin(writer http.ResponseWriter, req *http.Request) {
	body := LoginRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		sendErrorResponse(writer, "Malformed request body", http.StatusBadRequest)
		return
	}

	user, err := login(req.Context(), c.dbQueries, body.Email, body.Password)
	if err != nil {
		sendErrorResponse(writer, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	res := UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	sendResponse(writer, res, http.StatusOK)
}

func login(context context.Context, dbQueries *database.Queries, email string, password string) (database.User, error) {
	user, err := dbQueries.GetUserByEmail(context, email)
	if err != nil {
		log.Printf("Failed login attempt from %s: no such email: %v", email, err)
		return database.User{}, err
	}

	ok, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil {
		log.Printf("Failed login attempt from %s: password check failed: %v", email, err)
		return database.User{}, err
	}
	if !ok {
		log.Printf("Failed login attempt from %s: incorrect password", email)
		return database.User{}, errors.New("incorrect password")
	}

	return user, nil
}
