package handlers

import (
	"learning/internal/repositories"
	"learning/internal/utils"
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
	req, err := validators.ValidateCreateCampaign(r)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the repository to create a campaign
	campaign, err := h.Repo.CreateCampaign(*req)
	if err != nil {
		utils.JSONError(w, "Failed to create campaign", http.StatusInternalServerError)
		return
	}

	// Return created campaign
	utils.JSONSuccess(w, campaign, http.StatusCreated)
}
