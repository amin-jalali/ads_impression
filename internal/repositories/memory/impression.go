package memory

import (
	"errors"
	"learning/internal/entities"
	"net/http"
	"sync"
	"time"
)

type InMemoryImpressionRepository struct {
	mu          sync.Mutex
	campaigns   map[string]entities.Campaign
	impressions map[string]map[string]time.Time
	stats       map[string]entities.Stats
}

func NewInMemoryImpressionRepository(sharedCampaigns map[string]entities.Campaign) *InMemoryImpressionRepository {
	return &InMemoryImpressionRepository{
		campaigns:   sharedCampaigns, // Share storage with campaign repository
		impressions: make(map[string]map[string]time.Time),
		stats:       make(map[string]entities.Stats),
	}
}

func (r *InMemoryImpressionRepository) TrackImpression(req entities.TrackImpressionRequest) (error, int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Ensure campaign exists
	if _, exists := r.campaigns[req.CampaignID]; !exists {
		return errors.New("campaign not found"), http.StatusNotFound
	}

	now := time.Now()
	lastImpression, seen := r.impressions[req.CampaignID][req.UserID]

	// Enforce TTL for impressions (1 hour)
	if seen && now.Sub(lastImpression) < time.Hour {
		return errors.New("duplicate impression"), http.StatusOK
	}

	// Store impression
	if r.impressions[req.CampaignID] == nil {
		r.impressions[req.CampaignID] = make(map[string]time.Time)
	}
	r.impressions[req.CampaignID][req.UserID] = now

	// Update stats
	stats := r.stats[req.CampaignID]
	stats.LastHour++
	stats.LastDay++
	stats.TotalCount++
	r.stats[req.CampaignID] = stats

	return nil, http.StatusOK
}
