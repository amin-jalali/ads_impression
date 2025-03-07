package tests

import (
	"encoding/json"
	"learning/internal/entities"
	"net/http/httptest"
	"testing"
)

func GetCampaignCreateResponse(resp *httptest.ResponseRecorder, t *testing.T) entities.Campaign {
	var response struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Data    entities.Campaign `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("❌ Failed to decode response: %v", err)
	}
	campaign := response.Data
	return campaign
}

func GetStatsResponse(resp *httptest.ResponseRecorder, t *testing.T) entities.Stats {
	var response struct {
		Success bool           `json:"success"`
		Message string         `json:"message"`
		Data    entities.Stats `json:"data"`
	}

	// Decode response to get the campaign ID
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("❌ Failed to decode response: %v", err)
	}
	stats := response.Data
	return stats
}
