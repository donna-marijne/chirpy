package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ValidateChirpRequest struct {
	Body string `json:"body"`
}

type ValidateChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerValidateChirp(writer http.ResponseWriter, req *http.Request) {
	chirp := ValidateChirpRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		sendErrorResponse(writer, "Something went wrong", http.StatusBadRequest)
		return
	}

	if len(chirp.Body) > 140 {
		sendErrorResponse(writer, "Chirp is too long", http.StatusBadRequest)
		return
	}

	cleanedBody := removeProfanity(chirp.Body)

	res := ValidateChirpResponse{
		CleanedBody: cleanedBody,
	}
	sendResponse(writer, res, http.StatusOK)
}
