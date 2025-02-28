package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type TrackImpressionRequest struct {
	CampaignID string `json:"campaign_id" validate:"required"`
	UserID     string `json:"user_id" validate:"required"`
	AdID       string `json:"ad_id" validate:"required"`
}

func (s *Server) TrackImpressionHandler(w http.ResponseWriter, r *http.Request) {
	var req TrackImpressionRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		s.Logger.Error("Invalid JSON format")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate request struct
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		s.Logger.Error("Validation failed")
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	err, status := s.updateImpression(req)
	if err != nil {
		s.Logger.Error("Validation failed: " + err.Error())
		http.Error(w, "impression set failed: "+err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Impression saved successfully"})
	if err != nil {
		s.Logger.Error("response failed: " + err.Error())
		return
	}
}

func (s *Server) updateImpression(req TrackImpressionRequest) (error, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check campaign exist
	if _, exists := s.campaigns[req.CampaignID]; !exists {
		return errors.New("campaign not found"), http.StatusNotFound
	}

	now := time.Now()
	lastImpression, seen := s.impressions[req.CampaignID][req.UserID]

	// Prevent duplicate impressions in TTL
	if seen && now.Sub(lastImpression) < time.Hour {
		return errors.New("duplicate impression"), http.StatusOK
	}

	// Update impression data
	s.impressions[req.CampaignID][req.UserID] = now
	stats := s.stats[req.CampaignID]
	stats.LastHour++
	stats.LastDay++
	stats.TotalCount++
	s.stats[req.CampaignID] = stats
	return nil, http.StatusOK
}
