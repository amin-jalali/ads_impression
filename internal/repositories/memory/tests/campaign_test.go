package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"learning/internal/handlers"
	"learning/internal/repositories/memory"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Helper function to send a request and validate response
func sendTestRequest(t *testing.T, handler http.HandlerFunc, method, url, requestBody string, expectedStatus int, expectedError string) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	// Check HTTP status code
	if resp.Code != expectedStatus {
		t.Errorf("❌ [%s] Expected status %d, got %d. Request: %s", url, expectedStatus, resp.Code, requestBody)
	}

	// Check response body for expected error
	if expectedError != "" {
		var responseMap map[string]any
		_ = json.Unmarshal(resp.Body.Bytes(), &responseMap)

		if message, ok := responseMap["message"].(string); ok {
			if !strings.Contains(message, expectedError) {
				t.Errorf("❌ [%s] Expected error message to contain %q, but got %q. Request: %s", url, expectedError, message, requestBody)
			}
		} else {
			t.Errorf("❌ [%s] Expected a response with a 'message' field, but got %q", url, resp.Body.String())
		}
	}
}

func TestCreateCampaignHandler(t *testing.T) {
	// Initialize a shared in-memory server
	memServer := memory.NewServer()

	// Pass shared memory to the repository
	repo := memory.NewInMemoryCampaignRepository(memServer)
	handler := handlers.NewCampaignHandler(repo)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  string
	}{
		{"Valid Campaign", `{"name": "Test Campaign", "start_time": "2025-01-01T00:00:00Z"}`, http.StatusCreated, ""},
		{"❌ Malformed JSON", `{"name": "Test Campaign", "start_time": "invalid-date", }`, http.StatusBadRequest, "invalid JSON payload"},
		{"❌ Missing Name Field", `{"name": "", "start_time": "2025-01-01T00:00:00Z"}`, http.StatusBadRequest, "Field validation for 'Name' failed"},
		{"❌ Missing Start Time", `{"name": "Test Campaign", "start_time": ""}`, http.StatusBadRequest, "invalid JSON payload"},
		{"❌ Empty JSON Body", `{}`, http.StatusBadRequest, "Field validation"},
		{"❌ Unknown Field", `{"test": "test"}`, http.StatusBadRequest, "invalid JSON payload"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Println(test.name)
			sendTestRequest(t, handler.CreateCampaignHandler, http.MethodPost, "/api/v1/campaigns", test.requestBody, test.expectedStatus, test.expectedError)
		})
	}
}
