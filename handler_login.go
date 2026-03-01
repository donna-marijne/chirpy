package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/donnamarijne/chirpy/internal/auth"
	"github.com/donnamarijne/chirpy/internal/database"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	token, err := auth.MakeJWT(
		user.ID,
		c.secret,
		time.Hour,
	)
	if err != nil {
		log.Printf("Error making a JWT: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}
	refreshToken, err := c.dbQueries.CreateRefreshToken(req.Context(), createRefreshTokenParams)
	if err != nil {
		log.Printf("Error creating a refresh token: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	res := LoginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        token,
		RefreshToken: refreshToken.Token,
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
