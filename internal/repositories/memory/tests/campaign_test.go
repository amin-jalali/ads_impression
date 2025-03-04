package tests

import (
	"bytes"
	"learning/internal/repositories/memory/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateCampaignHandler(t *testing.T) {
	s := handlers.NewServer()

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid Campaign Creation",
			requestBody:    `{"name": "Test Campaign", "start_time": "2025-01-01T00:00:00Z"}`,
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		}, {
			name:           "Malformed JSON",
			requestBody:    `{"name": "Test Campaign", "start_time": "invalid-date", }`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON format",
		}, {
			name:           "Missing Name Field",
			requestBody:    `{"name": "", "start_time": "2025-01-01T00:00:00Z"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Error:Field validation for 'Name' failed on the 'required'",
		}, {
			name:           "Missing Start Time",
			requestBody:    `{"name": "Test Campaign", "start_time": ""}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON format",
		}, {
			name:           "Empty JSON Body",
			requestBody:    `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
		}, {
			name:           "Unknown Field",
			requestBody:    `{"test": "test"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON format",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewBuffer([]byte(test.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			s.CreateCampaignHandler(resp, req)

			if resp.Code != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, resp.Code)
			}

			if test.expectedError != "" {
				body := resp.Body.String()
				if !strings.Contains(body, test.expectedError) {
					t.Errorf("expected error message to contain %q, but got %q", test.expectedError, body)
				}
			}
		})
	}
}

func TestCreateCampaignWithInvalidInput(t *testing.T) {
	s := handlers.NewServer()

	invalidPayloads := []string{
		`{}`,
		`{"name": "", "start_time": ""}`,
		`{"name": 123, "start_time": 456}`,
		`{"name": "Valid Name", "start_time": "invalid-date"}`,
	}

	for _, payload := range invalidPayloads {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		s.CreateCampaignHandler(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d for payload %s", http.StatusBadRequest, resp.Code, payload)
		}
	}
}
