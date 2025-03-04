package validators

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func CreateCampaignValidator(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return errors.New("Invalid JSON format: " + err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}
