package handlers

import (
	"learning/internal/repositories"
	"learning/internal/utils"
	"learning/internal/validators"
	"net/http"
)

// StatsHandler uses a generic repository
type StatsHandler struct {
	Repo repositories.StatsRepository
}

// NewStatsHandler Constructor function
func NewStatsHandler(repo repositories.StatsRepository) *StatsHandler {
	return &StatsHandler{Repo: repo}
}

func (h *StatsHandler) GetCampaignStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Validate campaign ID
	campaignID, err := validators.ValidateCampaignID(r)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call repository to get stats
	stats, exists := h.Repo.GetCampaignStats(campaignID)
	if !exists {
		utils.JSONError(w, "Campaign not found", http.StatusNotFound)
		return
	}

	// Return campaign stats
	utils.JSONSuccess(w, stats, http.StatusOK)
}
