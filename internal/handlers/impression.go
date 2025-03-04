package handlers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"learning/internal/entities"
	"learning/internal/repositories"
	"net/http"
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
	var req entities.TrackImpressionRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate request struct
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	err, status := h.Repo.TrackImpression(req)
	if err != nil {
		http.Error(w, "impression set failed: "+err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Impression saved successfully"})
	if err != nil {
		http.Error(w, "response failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

//func (s *Server) TrackImpressionHandler(w http.ResponseWriter, r *http.Request) {
//	var req entities.TrackImpressionRequest
//
//	decoder := json.NewDecoder(r.Body)
//	decoder.DisallowUnknownFields()
//
//	if err := decoder.Decode(&req); err != nil {
//		s.Logger.Error("Invalid JSON format")
//		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
//		return
//	}
//
//	// Validate request struct
//	validate := validator.New()
//	if err := validate.Struct(req); err != nil {
//		s.Logger.Error("Validation failed")
//		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	err, status := s.updateImpression(req)
//	if err != nil {
//		s.Logger.Error("Validation failed: " + err.Error())
//		http.Error(w, "impression set failed: "+err.Error(), status)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	err = json.NewEncoder(w).Encode(map[string]string{"message": "Impression saved successfully"})
//	if err != nil {
//		s.Logger.Error("response failed: " + err.Error())
//		return
//	}
//}
//
//func (s *Server) updateImpression(req entities.TrackImpressionRequest) (error, int) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	// Check campaign exist
//	if _, exists := s.campaigns[req.CampaignID]; !exists {
//		return errors.New("campaign not found"), http.StatusNotFound
//	}
//
//	now := time.Now()
//	lastImpression, seen := s.impressions[req.CampaignID][req.UserID]
//
//	// Prevent duplicate impressions in TTL
//	if seen && now.Sub(lastImpression) < time.Hour {
//		return errors.New("duplicate impression"), http.StatusOK
//	}
//
//	// Update impression data
//	s.impressions[req.CampaignID][req.UserID] = now
//	stats := s.stats[req.CampaignID]
//	stats.LastHour++
//	stats.LastDay++
//	stats.TotalCount++
//	s.stats[req.CampaignID] = stats
//	return nil, http.StatusOK
//}
