package handlers

import (
	"encoding/json"
	"learning/internal/entities"
	"learning/internal/repositories"
	"learning/internal/validators"
	"net/http"
)

// CampaignHandler uses a generic repository
type CampaignHandler struct {
	Repo repositories.CampaignRepository
}

// NewCampaignHandler Constructor function
func NewCampaignHandler(repo repositories.CampaignRepository) *CampaignHandler {
	return &CampaignHandler{Repo: repo}
}

func (h *CampaignHandler) CreateCampaignHandler(w http.ResponseWriter, r *http.Request) {
	// Validate input
	if err := validators.CreateCampaignValidator(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode request body
	var req entities.CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Call the repository to create a campaign
	campaign, err := h.Repo.CreateCampaign(req)
	if err != nil {
		http.Error(w, "Failed to create campaign", http.StatusInternalServerError)
		return
	}

	// Return created campaign
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(campaign); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

//func (s *Server) CreateCampaignHandler(w http.ResponseWriter, r *http.Request) {
//	var req entities.CreateCampaignRequest
//
//	err := validators.CreateCampaignValidator(r)
//	if err != nil {
//		return
//	}
//
//	id := uuid.New().String()
//	campaign := entities.Campaign{ID: id, Name: req.Name, StartTime: req.StartTime}
//
//	s.initiator(id, campaign)
//
//	w.WriteHeader(http.StatusCreated)
//	if err := json.NewEncoder(w).Encode(campaign); err != nil {
//		s.Logger.Error("Failed to encode response")
//	}
//}
//
//func (s *Server) initiator(id string, campaign entities.Campaign) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	s.campaigns[id] = campaign
//	s.impressions[id] = make(map[string]time.Time)
//	s.stats[id] = entities.Stats{CampaignID: id, LastHour: 0, LastDay: 0, TotalCount: 0}
//}
