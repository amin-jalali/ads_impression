package handlers

import (
	"encoding/json"
	"learning/internal/repositories"
	"net/http"
	"net/url"
	"regexp"
)

// StatsHandler uses a generic repository
type StatsHandler struct {
	Repo repositories.StatsRepository
}

// NewStatsHandler Constructor function
func NewStatsHandler(repo repositories.StatsRepository) *StatsHandler {
	return &StatsHandler{Repo: repo}
}

// Validate campaign ID
var validIDRegexp = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`).MatchString

func (h *StatsHandler) GetCampaignStatsHandler(w http.ResponseWriter, r *http.Request) {
	campaignID := r.URL.Path[len("/api/v1/campaigns/"):]

	decodedID, err := url.QueryUnescape(campaignID)
	if err != nil {
		http.Error(w, "Invalid campaign ID format", http.StatusBadRequest)
		return
	}

	if decodedID == "" || len(decodedID) < 8 || !validIDRegexp(decodedID) {
		http.Error(w, "Invalid campaign ID", http.StatusBadRequest)
		return
	}

	// Call repository to get stats
	stats, exists := h.Repo.GetCampaignStats(decodedID)
	if !exists {
		http.Error(w, "Campaign not found", http.StatusNotFound)
		return
	}

	// Respond with stats
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

//
//var validIDRegexp = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`).MatchString
//
//func (s *Server) GetCampaignStatsHandler(w http.ResponseWriter, r *http.Request) {
//	campaignID := r.URL.Path[len("/api/v1/campaigns/"):]
//
//	decodedID, err := url.QueryUnescape(campaignID)
//	if err != nil {
//		s.Logger.Error("Invalid campaign ID format")
//		http.Error(w, "Invalid campaign ID format", http.StatusBadRequest)
//		return
//	}
//
//	if decodedID == "" || len(decodedID) < 8 || !validIDRegexp(decodedID) {
//		s.Logger.Error("Invalid campaign ID")
//		http.Error(w, "Invalid campaign ID", http.StatusBadRequest)
//		return
//	}
//
//	stats, exists := func() (entities.Stats, bool) {
//		s.mu.Lock()
//		defer s.mu.Unlock()
//
//		stats, exists := s.stats[decodedID]
//
//		return stats, exists
//	}()
//
//	if !exists {
//		s.Logger.Error("Campaign not found")
//		http.Error(w, "Campaign not found", http.StatusNotFound)
//		return
//	}
//
//	err = json.NewEncoder(w).Encode(stats)
//	if err != nil {
//		s.Logger.Error("Failed To respond")
//		return
//	}
//}
