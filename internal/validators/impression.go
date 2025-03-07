package validators

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"learning/internal/entities"
)

// ValidateTrackImpression extracts and validates the impression request
func ValidateTrackImpression(r *http.Request) (*entities.TrackImpressionRequest, error) {
	var req entities.TrackImpressionRequest

	// Ensure request body is not empty
	if r.Body == nil {
		return nil, errors.New("empty request body")
	}

	// Decode JSON and disallow unknown fields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		fmt.Println("❌ Invalid JSON format:", err)
		return nil, errors.New("invalid JSON format")
	}

	// Validate request struct
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		fmt.Println("❌ Validation failed:", err)
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	return &req, nil
}
