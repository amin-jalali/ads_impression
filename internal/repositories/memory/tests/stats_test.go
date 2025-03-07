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
	// Initialize shared in-memory server
	memServer := memory.NewServer()

	// Pass shared memory to stats repository
	statsRepo := memory.NewInMemoryStatsRepository(memServer)
	statsHandler := handlers.NewStatsHandler(statsRepo)

	invalidCampaignIDs := []struct {
		campaignID     string
		expectedStatus int
	}{
		{"", http.StatusBadRequest},                          // Empty ID
		{"non-existent-id", http.StatusNotFound},             // ID that doesn't exist
		{"12345", http.StatusNotFound},                       // Too short ID
		{"aaaaaaaa-bbbb-cccc", http.StatusNotFound},          // Non-existent valid format ID
		{url.QueryEscape("!@#$%^&*()"), http.StatusNotFound}, // Invalid characters
	}

	for _, test := range invalidCampaignIDs {
		t.Run("Testing Campaign ID: "+test.campaignID, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/stats/"+test.campaignID, nil) // Fix path
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			statsHandler.GetCampaignStatsHandler(resp, req)

			// Check status code
			if resp.Code != test.expectedStatus {
				t.Errorf("❌ Expected status %d, got %d for campaign ID: %s", test.expectedStatus, resp.Code, test.campaignID)
			}
		})
	}
}
