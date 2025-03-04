package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// JSONSuccess Success response
func JSONSuccess(w http.ResponseWriter, data any, statusCode int) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: "Request successful",
		Data:    data,
	})
	if err != nil {
		// TODO
		return
	}
}

// JSONError Error response
func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Message: message,
	})
	if err != nil {
		//TODO
		return
	}
}
