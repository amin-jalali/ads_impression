package memory

import (
	"errors"
	"learning/cmd/config"
	"learning/internal/entities"
	logger2 "learning/internal/logger"
	"net/http"
	"sync"
	"time"
)

type InMemoryImpressionRepository struct {
	mu     sync.Mutex
	server *entities.Server // Use shared server instance
}

func NewInMemoryImpressionRepository(server *entities.Server) *InMemoryImpressionRepository {
	return &InMemoryImpressionRepository{
		server: server,
	}
}

func (r *InMemoryImpressionRepository) TrackImpression(req entities.TrackImpressionRequest) (error, int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Ensure campaign exists in shared storage
	if _, exists := r.server.Campaigns[req.CampaignID]; !exists {
		return errors.New("campaign not found"), http.StatusNotFound
	}

	now := time.Now()
	lastImpression, seen := r.server.Impressions[req.CampaignID][req.UserID]

	cfg, err := config.LoadConfig()
	logger := logger2.InitLogger()
	if err != nil {
		logger.Error(err.Error())
	}

	// Read port from config and convert to string
	ttl := cfg.App.TTL

	// Enforce TTL for impressions (1 hour)
	if seen && now.Sub(lastImpression) < time.Duration(ttl)*time.Second {
		return errors.New("duplicate impression"), http.StatusOK
	}

	// Store impression in shared memory
	if r.server.Impressions[req.CampaignID] == nil {
		r.server.Impressions[req.CampaignID] = make(map[string]time.Time)
	}
	r.server.Impressions[req.CampaignID][req.UserID] = now

	// Update stats in shared memory
	stats := r.server.Stats[req.CampaignID]
	stats.LastHour++
	stats.LastDay++
	stats.TotalCount++
	r.server.Stats[req.CampaignID] = stats

	return nil, http.StatusOK
}
