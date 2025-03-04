package entities

import "time"

type Campaign struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
}

type CreateCampaignRequest struct {
	Name      string    `json:"name" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
}
