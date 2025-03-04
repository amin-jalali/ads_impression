package entities

import "time"

type Impression struct {
	CampaignID string    `json:"campaign_id"`
	Timestamp  time.Time `json:"timestamp"`
	UserID     string    `json:"user_id"`
	AdID       string    `json:"ad_id"`
}

type TrackImpressionRequest struct {
	CampaignID string `json:"campaign_id" validate:"required"`
	UserID     string `json:"user_id" validate:"required"`
	AdID       string `json:"ad_id" validate:"required"`
}
