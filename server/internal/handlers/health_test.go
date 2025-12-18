package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_ServeHTTP(t *testing.T) {
	handler := NewHealthHandler("1.0.0")

	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{"GET returns 200", http.MethodGet, http.StatusOK},
		{"POST returns 405", http.MethodPost, http.StatusMethodNotAllowed},
		{"PUT returns 405", http.MethodPut, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestHealthHandler_ResponseBody(t *testing.T) {
	version := "test-version-123"
	handler := NewHealthHandler(version)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Check content type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("got Content-Type %q, want %q", contentType, "application/json")
	}

	// Parse response
	var resp HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "healthy" {
		t.Errorf("got status %q, want %q", resp.Status, "healthy")
	}

	if resp.Version != version {
		t.Errorf("got version %q, want %q", resp.Version, version)
	}

	if resp.Timestamp == "" {
		t.Error("timestamp should not be empty")
	}
}
