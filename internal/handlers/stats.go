package handlers

import (
	"learning/internal/repositories"
	"learning/internal/utils"
	"net/http"
	"strings"
)

type StatsHandler struct {
	Repo repositories.StatsRepository
}

func NewStatsHandler(repo repositories.StatsRepository) *StatsHandler {
	return &StatsHandler{Repo: repo}
}

func (h *StatsHandler) GetCampaignStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract campaign ID from URL
	campaignID := strings.TrimPrefix(r.URL.Path, "/api/v1/campaigns/stats/")
	if campaignID == "" || campaignID == "/api/v1/campaigns/stats" {
		utils.JSONError(w, "invalid campaign ID", http.StatusBadRequest)
		return
	}

	// Fetch stats
	stats, exists := h.Repo.GetCampaignStats(campaignID)
	if !exists {
		utils.JSONError(w, "campaign not found", http.StatusNotFound)
		return
	}

	// Return stats as JSON response
	utils.JSONSuccess(w, stats, http.StatusOK)
}
