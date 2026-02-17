package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserCreateRequest struct {
	Email string `json:"email"`
}

type UserCreateResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (c *apiConfig) handlerUserCreate(writer http.ResponseWriter, req *http.Request) {
	body := UserCreateRequest{}
	// err := unmarshalRequestBody(req, &body)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		sendErrorResponse(writer, "Malformed request body", http.StatusBadRequest)
		return
	}

	user, err := c.dbQueries.CreateUser(req.Context(), body.Email)
	if err != nil {
		log.Printf("Error from CreateUser: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusInternalServerError)
		return
	}

	res := UserCreateResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	sendResponse(writer, res, http.StatusCreated)
}

func unmarshalRequestBody(req *http.Request, obj *any) error {
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(obj)
}
