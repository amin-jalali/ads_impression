package repositories

import "learning/internal/entities"

type CampaignRepository interface {
	CreateCampaign(req entities.CreateCampaignRequest) (entities.Campaign, error)
}
