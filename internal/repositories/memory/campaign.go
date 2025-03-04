package memory

import (
	"github.com/google/uuid"
	"learning/internal/entities"
	"sync"
)

type InMemoryCampaignRepository struct {
	mu        sync.Mutex
	campaigns map[string]entities.Campaign
}

func NewInMemoryCampaignRepository() *InMemoryCampaignRepository {
	return &InMemoryCampaignRepository{
		campaigns: make(map[string]entities.Campaign),
	}
}

// CreateCampaign Implement
func (r *InMemoryCampaignRepository) CreateCampaign(req entities.CreateCampaignRequest) (entities.Campaign, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()
	campaign := entities.Campaign{
		ID:        id,
		Name:      req.Name,
		StartTime: req.StartTime,
	}

	r.campaigns[id] = campaign
	return campaign, nil
}

//func (s *Server) CreateCampaign() {
//
//	id := uuid.New().String()
//	campaign := entities.Campaign{ID: id, Name: req.Name, StartTime: req.StartTime}
//
//	s.initiator(id, campaign)
//
//}
//
//func (s *Server) initiator(id string, campaign entities.Campaign) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	s.campaigns[id] = campaign
//	s.impressions[id] = make(map[string]time.Time)
//	s.stats[id] = entities.Stats{CampaignID: id, LastHour: 0, LastDay: 0, TotalCount: 0}
//}
