package memory

import (
	"learning/internal/entities"
	"time"
)

func NewServer() *entities.Server {
	return &entities.Server{
		Campaigns:   make(map[string]entities.Campaign),
		Impressions: make(map[string]map[string]time.Time),
		Stats:       make(map[string]entities.Stats),
	}
}
