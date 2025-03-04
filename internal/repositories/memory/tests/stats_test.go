package tests

import (
	"learning/internal/handlers"
	"learning/internal/repositories/memory"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetCampaignStatsWithInvalidInput(t *testing.T) {
	// Initialize repositories and handlers
	statsRepo := memory.NewInMemoryStatsRepository()
	statsHandler := handlers.NewStatsHandler(statsRepo)

	invalidCampaignIDs := []struct {
		campaignID     string
		expectedStatus int
	}{
		{"", http.StatusBadRequest},                            // Empty ID
		{"non-existent-id", http.StatusNotFound},               // ID that doesn't exist
		{"12345", http.StatusBadRequest},                       // Too short ID
		{"aaaaaaaa-bbbb-cccc", http.StatusNotFound},            // Non-existent valid format ID
		{url.QueryEscape("!@#$%^&*()"), http.StatusBadRequest}, // Invalid characters
	}

	for _, test := range invalidCampaignIDs {
		t.Run("Testing Campaign ID: "+test.campaignID, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/"+test.campaignID, nil)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			statsHandler.GetCampaignStatsHandler(resp, req)

			// Check status code
			if resp.Code != test.expectedStatus {
				t.Errorf("expected status %d, got %d for campaign ID: %s", test.expectedStatus, resp.Code, test.campaignID)
			}
		})
	}
}
