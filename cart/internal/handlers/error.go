package handlers

import (
	"net/http"
	"encoding/json"
	"log"
)

type handlerErrorResponse struct {
	Status int    `json:"status"`
	URL    string `json:"url"`
	Msg    string `json:"msg"`
}

func sendError(r handlerErrorResponse, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Printf("failed to send structured response: %w\n", err)
	}
}