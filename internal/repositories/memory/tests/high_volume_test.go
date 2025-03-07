package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"learning/internal/entities"
	"learning/internal/handlers"
	"learning/internal/repositories/memory"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHighVolumeConcurrentRequests(t *testing.T) {
	// Initialize a shared in-memory server
	memServer := memory.NewServer()

	// Pass shared memory to repositories
	campaignRepo := memory.NewInMemoryCampaignRepository(memServer)
	impressionRepo := memory.NewInMemoryImpressionRepository(memServer)
	statsRepo := memory.NewInMemoryStatsRepository(memServer)

	// Inject repositories into handlers
	campaignHandler := handlers.NewCampaignHandler(campaignRepo)
	impressionHandler := handlers.NewImpressionHandler(impressionRepo)
	statsHandler := handlers.NewStatsHandler(statsRepo)

	var wg sync.WaitGroup

	// Step 1: Create a new campaign
	campaignReq := entities.CreateCampaignRequest{Name: "High Volume Campaign", StartTime: time.Now()}
	jsonCampaign, _ := json.Marshal(campaignReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewBuffer(jsonCampaign))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	campaignHandler.CreateCampaignHandler(resp, req)

	var response struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Data    entities.Campaign `json:"data"` // Extract "data" field into Campaign struct
	}

	// Decode response to get the campaign ID
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("❌ Failed to decode response: %v", err)
	}
	campaign := response.Data

	// Step 2: Send concurrent impression requests
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()
			impReq := entities.TrackImpressionRequest{CampaignID: campaign.ID, UserID: userID, AdID: "123"}
			jsonImp, _ := json.Marshal(impReq)
			request := httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			impressionHandler.TrackImpressionHandler(response, request)
		}(fmt.Sprintf("user%d", i))
	}

	wg.Wait()

	// Step 3: Retrieve campaign stats
	req = httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/stats/"+campaign.ID, nil) // Fix path
	resp = httptest.NewRecorder()
	statsHandler.GetCampaignStatsHandler(resp, req)

	// Decode response to get the stats
	var res struct {
		Success bool           `json:"success"`
		Message string         `json:"message"`
		Data    entities.Stats `json:"data"` // Extract "data" field into Campaign struct
	}

	// Decode response to get the campaign ID
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("❌ Failed to decode response: %v", err)
	}
	stats := res.Data

	// Step 4: Validate results
	if stats.TotalCount != 100 {
		t.Errorf("❌ Expected total count 100, got %d", stats.TotalCount)
	}
}
