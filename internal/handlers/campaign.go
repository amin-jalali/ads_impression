package handlers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CreateCampaignRequest struct {
	Name      string    `json:"name" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
}

func (s *Server) CreateCampaignHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateCampaignRequest
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

	id := uuid.New().String()
	campaign := Campaign{ID: id, Name: req.Name, StartTime: req.StartTime}

	s.initiator(id, campaign)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(campaign); err != nil {
		s.Logger.Error("Failed to encode response")
	}
}

func (s *Server) initiator(id string, campaign Campaign) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.campaigns[id] = campaign
	s.impressions[id] = make(map[string]time.Time)
	s.stats[id] = Stats{CampaignID: id, LastHour: 0, LastDay: 0, TotalCount: 0}
}
