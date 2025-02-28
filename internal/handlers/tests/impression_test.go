package tests

import (
	"bytes"
	"encoding/json"
	"learning/internal/handlers"
	"learning/internal/logger"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestTrackImpressionWithInvalidInput(t *testing.T) {
	s := handlers.NewServer()

	invalidPayloads := []string{
		`{}`,
		`{"campaign_id": "", "user_id": "", "ad_id": ""}`,
		`{"campaign_id": "123", "user_id": 456, "ad_id": true}`,
		`{"campaign_id": "non-existent", "user_id": "user123", "ad_id": "ad456"}`,
	}

	for _, payload := range invalidPayloads {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		s.TrackImpressionHandler(resp, req)

		if payload == `{"campaign_id": "non-existent", "user_id": "user123", "ad_id": "ad456"}` {
			if resp.Code != http.StatusNotFound {
				t.Errorf("expected status %d, got %d for payload %s", http.StatusNotFound, resp.Code, payload)
			}
		} else {
			if resp.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d for payload %s", http.StatusBadRequest, resp.Code, payload)
			}
		}
	}
}

func TestTrackImpressionWithTTL(t *testing.T) {
	s := handlers.NewServer()

	campaignReq := handlers.CreateCampaignRequest{Name: "TTL Campaign", StartTime: time.Now()}
	jsonCampaign, _ := json.Marshal(campaignReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewBuffer(jsonCampaign))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.CreateCampaignHandler(resp, req)

	var campaign handlers.Campaign
	err := json.NewDecoder(resp.Body).Decode(&campaign)
	if err != nil {
		logger.Log.Error("invalid request. unable to decode")
		return
	}

	impReq := handlers.TrackImpressionRequest{CampaignID: campaign.ID, UserID: "user123", AdID: "ad456"}
	jsonImp, _ := json.Marshal(impReq)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	s.TrackImpressionHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	s.TrackImpressionHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestConcurrentImpressionTracking(t *testing.T) {
	s := handlers.NewServer()
	var wg sync.WaitGroup

	campaignReq := handlers.CreateCampaignRequest{Name: "Concurrent Campaign", StartTime: time.Now()}
	jsonCampaign, _ := json.Marshal(campaignReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewBuffer(jsonCampaign))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.CreateCampaignHandler(resp, req)

	var campaign handlers.Campaign
	err := json.NewDecoder(resp.Body).Decode(&campaign)
	if err != nil {
		return
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			impReq := handlers.TrackImpressionRequest{CampaignID: campaign.ID, UserID: "user123", AdID: "ad456"}
			jsonImp, _ := json.Marshal(impReq)
			request := httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			s.TrackImpressionHandler(response, request)
		}()
	}
	wg.Wait()

	req = httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/"+campaign.ID, nil)
	resp = httptest.NewRecorder()
	s.GetCampaignStatsHandler(resp, req)

	var stats handlers.Stats
	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		logger.Log.Error("invalid request. unable to decode")
		return
	}

	if stats.TotalCount == 0 {
		t.Errorf("expected total count greater than 0, got %d", stats.TotalCount)
	}
}

func TestTrackImpressionForNonExistentCampaign(t *testing.T) {
	s := handlers.NewServer()
	impReq := handlers.TrackImpressionRequest{CampaignID: "non-existent-id", UserID: "user123", AdID: "ad456"}
	jsonImp, _ := json.Marshal(impReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.TrackImpressionHandler(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, resp.Code)
	}
}
