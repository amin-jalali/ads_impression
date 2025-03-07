package memory

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"learning/internal/entities"
)

type InMemoryCampaignRepository struct {
	mu     sync.Mutex
	server *entities.Server // Use shared server instance
}

// NewInMemoryCampaignRepository Use shared `Server` storage instead of creating new maps
func NewInMemoryCampaignRepository(server *entities.Server) *InMemoryCampaignRepository {
	return &InMemoryCampaignRepository{
		server: server,
	}
}

// CreateCampaign Store campaigns in shared memory
func (r *InMemoryCampaignRepository) CreateCampaign(req entities.CreateCampaignRequest) (entities.Campaign, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()
	campaign := entities.Campaign{
		ID:        id,
		Name:      req.Name,
		StartTime: req.StartTime,
	}

	// Store in shared memory
	r.server.Campaigns[id] = campaign
	r.server.Impressions[id] = make(map[string]time.Time)                                       // Initialize impressions tracking
	r.server.Stats[id] = entities.Stats{CampaignID: id, LastHour: 0, LastDay: 0, TotalCount: 0} // Initialize stats

	return campaign, nil
}

// GetCampaigns Return campaigns from shared memory
func (r *InMemoryCampaignRepository) GetCampaigns() map[string]entities.Campaign {
	return r.server.Campaigns
}

// GetStats Return stats from shared memory
func (r *InMemoryCampaignRepository) GetStats() map[string]entities.Stats {
	return r.server.Stats
}
