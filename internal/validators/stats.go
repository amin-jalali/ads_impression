package validators

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
)

// Regular expression for valid campaign IDs
var validIDRegexp = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`).MatchString

func ValidateCampaignID(r *http.Request) (string, error) {
	campaignID := r.URL.Path[len("/api/v1/campaigns/"):]

	decodedID, err := url.QueryUnescape(campaignID)
	if err != nil {
		return "", errors.New("invalid campaign ID format")
	}

	if decodedID == "" || len(decodedID) < 8 || !validIDRegexp(decodedID) {
		return "", errors.New("invalid campaign ID")
	}

	return decodedID, nil
}
