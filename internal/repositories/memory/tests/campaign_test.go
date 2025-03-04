package tests

import (
	"bytes"
	"encoding/json"
	"learning/internal/handlers"
	"learning/internal/repositories/memory"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateCampaignHandler(t *testing.T) {
	// Initialize the in-memory repository
	repo := memory.NewInMemoryCampaignRepository()
	handler := handlers.NewCampaignHandler(repo)

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
			expectedError:  "invalid JSON payload",
		}, {
			name:           "Missing Name Field",
			requestBody:    `{"name": "", "start_time": "2025-01-01T00:00:00Z"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Field validation for 'Name' failed",
		}, {
			name:           "Missing Start Time",
			requestBody:    `{"name": "Test Campaign", "start_time": ""}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid JSON payload",
		}, {
			name:           "Empty JSON Body",
			requestBody:    `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Field validation",
		}, {
			name:           "Unknown Field",
			requestBody:    `{"test": "test"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid JSON payload",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewBuffer([]byte(test.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			handler.CreateCampaignHandler(resp, req)

			// Check HTTP status code
			if resp.Code != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, resp.Code)
			}

			// Check response body for expected error
			if test.expectedError != "" {
				var responseMap map[string]any
				_ = json.Unmarshal(resp.Body.Bytes(), &responseMap)

				if message, ok := responseMap["message"].(string); ok {
					if !strings.Contains(message, test.expectedError) {
						t.Errorf("expected error message to contain %q, but got %q", test.expectedError, message)
					}
				} else {
					t.Errorf("expected a response with a 'message' field, but got %q", resp.Body.String())
				}
			}
		})
	}
}

func TestCreateCampaignWithInvalidInput(t *testing.T) {
	// Initialize the in-memory repository
	repo := memory.NewInMemoryCampaignRepository()
	handler := handlers.NewCampaignHandler(repo)

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

		handler.CreateCampaignHandler(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d for payload %s", http.StatusBadRequest, resp.Code, payload)
		}
	}
}
