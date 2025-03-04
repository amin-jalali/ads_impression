package repositories

import (
	"learning/internal/entities"
)

type ImpressionRepository interface {
	TrackImpression(req entities.TrackImpressionRequest) (error, int)
}
