package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"learning/internal/handlers"
	"learning/internal/logger"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHighVolumeConcurrentRequests(t *testing.T) {
	s := handlers.NewServer()
	var wg sync.WaitGroup

	campaignReq := handlers.CreateCampaignRequest{Name: "High Volume Campaign", StartTime: time.Now()}
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

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()
			impReq := handlers.TrackImpressionRequest{CampaignID: campaign.ID, UserID: userID, AdID: "ad456"}
			jsonImp, _ := json.Marshal(impReq)
			request := httptest.NewRequest(http.MethodPost, "/api/v1/impressions", bytes.NewBuffer(jsonImp))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			s.TrackImpressionHandler(response, request)
		}(fmt.Sprintf("user%d", i))
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

	if stats.TotalCount != 100 {
		t.Errorf("expected total count 100, got %d", stats.TotalCount)
	}
}
