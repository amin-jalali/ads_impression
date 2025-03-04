package entities

type Stats struct {
	CampaignID string `json:"campaign_id"`
	LastHour   int64  `json:"last_hour"`
	LastDay    int64  `json:"last_day"`
	TotalCount int64  `json:"total"`
}
