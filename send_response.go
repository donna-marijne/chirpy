package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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

type ErrorResponse struct {
	Error string `json:"error"`
}

func sendErrorResponse(writer http.ResponseWriter, error string, status int) {
	sendResponse(writer, ErrorResponse{Error: error}, status)
}
