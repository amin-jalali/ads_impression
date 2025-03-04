package repositories

import "learning/internal/entities"

type StatsRepository interface {
	GetCampaignStats(campaignID string) (entities.Stats, bool)
}
