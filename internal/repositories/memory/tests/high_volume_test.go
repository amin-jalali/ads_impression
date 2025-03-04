package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"learning/internal/entities"
	handlers2 "learning/internal/handlers"
	"learning/internal/repositories/memory"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHighVolumeConcurrentRequests(t *testing.T) {
	// Initialize in-memory repositories
	campaignRepo := memory.NewInMemoryCampaignRepository()
	impressionRepo := memory.NewInMemoryImpressionRepository(nil)
	statsRepo := memory.NewInMemoryStatsRepository()

	// Inject repositories into handlers
	campaignHandler := handlers2.NewCampaignHandler(campaignRepo)
	impressionHandler := handlers2.NewImpressionHandler(impressionRepo)
	statsHandler := handlers2.NewStatsHandler(statsRepo)

	var wg sync.WaitGroup

	// Step 1: Create a new campaign
	campaignReq := entities.CreateCampaignRequest{Name: "High Volume Campaign", StartTime: time.Now()}
	jsonCampaign, _ := json.Marshal(campaignReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewBuffer(jsonCampaign))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	campaignHandler.CreateCampaignHandler(resp, req)

	// Decode response to get the campaign ID
	var campaign entities.Campaign
	if err := json.NewDecoder(resp.Body).Decode(&campaign); err != nil {
		t.Fatalf("Failed to decode campaign response: %v", err)
	}

	// Step 2: Send concurrent impression requests
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()
			impReq := entities.TrackImpressionRequest{CampaignID: campaign.ID, UserID: userID, AdID: "ad456"}
			jsonImp, _ := json.Marshal(impReq)
			request := httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			impressionHandler.TrackImpressionHandler(response, request)
		}(fmt.Sprintf("user%d", i))
	}

	wg.Wait()

	// Step 3: Retrieve campaign stats
	req = httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/"+campaign.ID, nil)
	resp = httptest.NewRecorder()
	statsHandler.GetCampaignStatsHandler(resp, req)

	// Decode response to get the stats
	var stats entities.Stats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode stats response: %v", err)
	}

	// Step 4: Validate results
	if stats.TotalCount != 100 {
		t.Errorf("Expected total count 100, got %d", stats.TotalCount)
	}
}
