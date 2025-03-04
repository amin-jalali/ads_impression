package validators

import (
	"encoding/json"
	"errors"
	"learning/internal/entities"
	"net/http"
)

func ValidateTrackImpression(r *http.Request) (*entities.TrackImpressionRequest, error) {
	var req entities.TrackImpressionRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return nil, errors.New("invalid JSON format")
	}

	// Use shared validate instance
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	return &req, nil
}
