package memory

import (
	"learning/internal/entities"
	"sync"
)

type InMemoryStatsRepository struct {
	mu     sync.Mutex
	server *entities.Server // Add shared server instance
}

// NewInMemoryStatsRepository Accept shared `Server` instance
func NewInMemoryStatsRepository(server *entities.Server) *InMemoryStatsRepository {
	return &InMemoryStatsRepository{
		server: server, // Assign server instance
	}
}

// GetCampaignStats Fetch stats from shared memory
func (r *InMemoryStatsRepository) GetCampaignStats(campaignID string) (entities.Stats, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stats, exists := r.server.Stats[campaignID] // Fetch from shared memory
	return stats, exists
}
