package tests

import (
	"errors"
	"learning/cmd/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func mockListenAndServe() error {
	return nil
}

func mockListenAndServeError() error {
	return errors.New("mock ListenAndServe error")
}

func TestRunFailure(t *testing.T) {

	var err = server.Run(mockListenAndServeError)
	expectedLog := "mock ListenAndServe error"

	if err.Error() != expectedLog {
		t.Errorf("expected log output to contain %q, but got: %q", expectedLog, err)
	}
}

func TestListenAndServeFailure(t *testing.T) {
	srv := &http.Server{
		Addr:    ":9999",
		Handler: nil,
	}

	go func() {
		err := srv.ListenAndServe()
		if err == nil {
			t.Errorf("expected ListenAndServe to return an error, but got nil")
		}
	}()
	// TODO: Check if system atency causes go routine closes befor execution
	time.Sleep(500 * time.Millisecond)
	err := srv.Close()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

// Mock ListenAndServe function
func mockListenAndServeSuccess() error {
	return nil // Simulate successful server start
}

func mockListenAndServeFailure() error {
	return errors.New("server failed to start") // Simulate failure
}

func TestRun_Success(t *testing.T) {
	err := server.Run(mockListenAndServeSuccess)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestRun_Failure(t *testing.T) {
	err := server.Run(mockListenAndServeFailure)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

func TestSetupServer(t *testing.T) {
	handler := server.SetupServer()

	tests := []struct {
		name           string
		url            string
		method         string
		body           string
		expectedStatus int
	}{
		{"CreateCampaignHandler",
			"/api/v1/campaigns",
			"POST",
			`{"name": "Test Campaign", "start_time": "2025-01-01T00:00:00Z"}`,
			http.StatusCreated,
		}, {"TrackImpressionHandler",
			"/api/v1/impressions",
			"POST",
			`{"campaign_id": 123}`,
			http.StatusBadRequest,
		}, {
			"GetCampaignStatsHandler",
			"/api/v1/campaigns/123",
			"GET",
			"",
			http.StatusBadRequest,
		}, {
			"NotFoundHandler",
			"/unknown",
			"GET",
			"",
			http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json") // Set appropriate header for POST requests
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d for URL %s with method %s", tt.expectedStatus, recorder.Code, tt.url, tt.method)
			}
		})
	}
}
