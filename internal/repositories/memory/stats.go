package memory

import (
	"learning/internal/entities"
	"sync"
)

type InMemoryStatsRepository struct {
	mu    sync.Mutex
	stats map[string]entities.Stats
}

func NewInMemoryStatsRepository() *InMemoryStatsRepository {
	return &InMemoryStatsRepository{
		stats: make(map[string]entities.Stats),
	}
}

func (r *InMemoryStatsRepository) GetCampaignStats(campaignID string) (entities.Stats, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stats, exists := r.stats[campaignID]
	return stats, exists
}
