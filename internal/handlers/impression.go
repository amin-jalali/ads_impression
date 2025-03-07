package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"learning/internal/repositories"
	"learning/internal/validators"
)

// ImpressionHandler uses a generic repository
type ImpressionHandler struct {
	Repo repositories.ImpressionRepository
}

// NewImpressionHandler Constructor function
func NewImpressionHandler(repo repositories.ImpressionRepository) *ImpressionHandler {
	return &ImpressionHandler{Repo: repo}
}

func (h *ImpressionHandler) TrackImpressionHandler(w http.ResponseWriter, r *http.Request) {
	// Validate request using validators.ValidateTrackImpression
	req, err := validators.ValidateTrackImpression(r)
	if err != nil {
		log.Printf("❌ Request validation failed: %v", err)
		http.Error(w, "Request validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call the repository method to track the impression
	err, status := h.Repo.TrackImpression(*req)
	if err != nil {
		log.Printf("❌ Impression set failed: %v", err)
		http.Error(w, "Impression set failed: "+err.Error(), status)
		return
	}

	// Send success response
	response := map[string]string{"message": "Impression saved successfully"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("❌ Response encoding failed: %v", err)
		http.Error(w, "Response failed: "+err.Error(), http.StatusInternalServerError)
	}
}
