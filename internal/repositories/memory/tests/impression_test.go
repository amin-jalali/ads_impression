package tests

import (
	"bytes"
	"encoding/json"
	"learning/internal/entities"
	handlers2 "learning/internal/handlers"
	"learning/internal/repositories/memory"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestTrackImpressionWithInvalidInput(t *testing.T) {
	// Initialize repositories and handlers
	impressionRepo := memory.NewInMemoryImpressionRepository(nil)
	impressionHandler := handlers2.NewImpressionHandler(impressionRepo)

	invalidPayloads := []struct {
		payload        string
		expectedStatus int
	}{
		{`{}`, http.StatusBadRequest},
		{`{"campaign_id": "", "user_id": "", "ad_id": ""}`, http.StatusBadRequest},
		{`{"campaign_id": "123", "user_id": 456, "ad_id": true}`, http.StatusBadRequest},
		{`{"campaign_id": "non-existent", "user_id": "user123", "ad_id": "ad456"}`, http.StatusNotFound},
	}

	for _, test := range invalidPayloads {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer([]byte(test.payload)))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		impressionHandler.TrackImpressionHandler(resp, req)

		if resp.Code != test.expectedStatus {
			t.Errorf("expected status %d, got %d for payload %s", test.expectedStatus, resp.Code, test.payload)
		}
	}
}

func TestTrackImpressionWithTTL(t *testing.T) {
	// Initialize repositories and handlers
	campaignRepo := memory.NewInMemoryCampaignRepository()
	impressionRepo := memory.NewInMemoryImpressionRepository(nil)

	campaignHandler := handlers2.NewCampaignHandler(campaignRepo)
	impressionHandler := handlers2.NewImpressionHandler(impressionRepo)

	// Step 1: Create a campaign
	campaignReq := entities.CreateCampaignRequest{Name: "TTL Campaign", StartTime: time.Now()}
	jsonCampaign, _ := json.Marshal(campaignReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewBuffer(jsonCampaign))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	campaignHandler.CreateCampaignHandler(resp, req)

	// Decode response to get campaign ID
	var campaign entities.Campaign
	if err := json.NewDecoder(resp.Body).Decode(&campaign); err != nil {
		t.Fatalf("Failed to decode campaign response: %v", err)
	}

	// Step 2: Track an impression
	impReq := entities.TrackImpressionRequest{CampaignID: campaign.ID, UserID: "user123", AdID: "ad456"}
	jsonImp, _ := json.Marshal(impReq)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	impressionHandler.TrackImpressionHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	// Step 3: Track the same impression within TTL (should return HTTP 200 OK)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	impressionHandler.TrackImpressionHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestConcurrentImpressionTracking(t *testing.T) {
	// Initialize repositories and handlers
	campaignRepo := memory.NewInMemoryCampaignRepository()
	impressionRepo := memory.NewInMemoryImpressionRepository(nil)
	statsRepo := memory.NewInMemoryStatsRepository()

	campaignHandler := handlers2.NewCampaignHandler(campaignRepo)
	impressionHandler := handlers2.NewImpressionHandler(impressionRepo)
	statsHandler := handlers2.NewStatsHandler(statsRepo)

	var wg sync.WaitGroup

	// Step 1: Create a campaign
	campaignReq := entities.CreateCampaignRequest{Name: "Concurrent Campaign", StartTime: time.Now()}
	jsonCampaign, _ := json.Marshal(campaignReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewBuffer(jsonCampaign))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	campaignHandler.CreateCampaignHandler(resp, req)

	// Decode response to get campaign ID
	var campaign entities.Campaign
	if err := json.NewDecoder(resp.Body).Decode(&campaign); err != nil {
		t.Fatalf("Failed to decode campaign response: %v", err)
	}

	// Step 2: Track multiple impressions concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			impReq := entities.TrackImpressionRequest{CampaignID: campaign.ID, UserID: "user123", AdID: "ad456"}
			jsonImp, _ := json.Marshal(impReq)
			request := httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			impressionHandler.TrackImpressionHandler(response, request)
		}()
	}
	wg.Wait()

	// Step 3: Retrieve campaign stats
	req = httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/"+campaign.ID, nil)
	resp = httptest.NewRecorder()
	statsHandler.GetCampaignStatsHandler(resp, req)

	// Decode response to get stats
	var stats entities.Stats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode stats response: %v", err)
	}

	// Step 4: Validate results
	if stats.TotalCount == 0 {
		t.Errorf("expected total count greater than 0, got %d", stats.TotalCount)
	}
}

func TestTrackImpressionForNonExistentCampaign(t *testing.T) {
	// Initialize repositories and handlers
	impressionRepo := memory.NewInMemoryImpressionRepository(nil)
	impressionHandler := handlers2.NewImpressionHandler(impressionRepo)

	// Attempt to track an impression for a non-existent campaign
	impReq := entities.TrackImpressionRequest{CampaignID: "non-existent-id", UserID: "user123", AdID: "ad456"}
	jsonImp, _ := json.Marshal(impReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	impressionHandler.TrackImpressionHandler(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, resp.Code)
	}
}
