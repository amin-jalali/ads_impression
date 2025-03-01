package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
)

var validIDRegexp = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`).MatchString

func (s *Server) GetCampaignStatsHandler(w http.ResponseWriter, r *http.Request) {
	campaignID := r.URL.Path[len("/api/v1/campaigns/"):]

	decodedID, err := url.QueryUnescape(campaignID)
	if err != nil {
		s.Logger.Error("Invalid campaign ID format")
		http.Error(w, "Invalid campaign ID format", http.StatusBadRequest)
		return
	}

	if decodedID == "" || len(decodedID) < 8 || !validIDRegexp(decodedID) {
		s.Logger.Error("Invalid campaign ID")
		http.Error(w, "Invalid campaign ID", http.StatusBadRequest)
		return
	}

	stats, exists := func() (Stats, bool) {
		s.mu.Lock()
		defer s.mu.Unlock()

		stats, exists := s.stats[decodedID]

		return stats, exists
	}()

	if !exists {
		s.Logger.Error("Campaign not found")
		http.Error(w, "Campaign not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		s.Logger.Error("Failed To respond")
		return
	}
}
