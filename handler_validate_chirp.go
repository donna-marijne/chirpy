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

type ValidateChirpErrorResponse struct {
	Error string `json:"error"`
}

func handlerValidateChirp(writer http.ResponseWriter, req *http.Request) {
	chirp := ValidateChirpRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		sendErrorResponse(writer, "Something went wrong", 400)
		return
	}

	if len(chirp.Body) > 140 {
		sendErrorResponse(writer, "Chirp is too long", 400)
		return
	}

	cleanedBody := removeProfanity(chirp.Body)

	res := ValidateChirpResponse{
		CleanedBody: cleanedBody,
	}
	sendResponse(writer, res, 200)
}

func sendErrorResponse(writer http.ResponseWriter, error string, status int) {
	sendResponse(writer, ValidateChirpErrorResponse{Error: error}, status)
}

func sendResponse(writer http.ResponseWriter, obj any, status int) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Error marshalling response: %+v", obj)
		writer.WriteHeader(500)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(bytes)
}
