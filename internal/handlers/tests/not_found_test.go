package tests

import (
	"learning/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFoundHandler(t *testing.T) {
	_ = handlers.NewServer()

	paths := []string{
		"/invalid-route",
		"/api/v1/unknown",
		"/wrongpath",
		"/random/endpoint",
	}

	for _, path := range paths {
		t.Run("Testing "+path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			resp := httptest.NewRecorder()

			handlers.NotFoundHandler(resp, req)

			if resp.Code != http.StatusNotFound {
				t.Errorf("expected status %d, got %d for path %s", http.StatusNotFound, resp.Code, path)
			}

			expectedBody := "404 not found\n"
			if resp.Body.String() != expectedBody {
				t.Errorf("expected response body %q, got %q", expectedBody, resp.Body.String())
			}
		})
	}
}
