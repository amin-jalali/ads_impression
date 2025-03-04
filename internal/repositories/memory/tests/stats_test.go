package tests

import (
	"learning/internal/repositories/memory/handlers"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetCampaignStatsWithInvalidInput(t *testing.T) {
	s := handlers.NewServer()

	invalidCampaignIDs := []string{
		"",
		"non-existent-id",
		"12345",
		"aaaaaaaa-bbbb-cccc",
		url.QueryEscape("!@#$%^&*()"),
	}

	for _, campaignID := range invalidCampaignIDs {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/"+campaignID, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		s.GetCampaignStatsHandler(resp, req)

		if campaignID == "non-existent-id" || campaignID == "aaaaaaaa-bbbb-cccc" {
			if resp.Code != http.StatusNotFound {
				t.Errorf("expected status %d, got %d for campaign ID: %s", http.StatusNotFound, resp.Code, campaignID)
			}
		} else {
			if resp.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d for campaign ID: %s", http.StatusBadRequest, resp.Code, campaignID)
			}
		}
	}
}
