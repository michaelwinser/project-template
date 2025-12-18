package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID_GeneratesID(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request ID is in context
		id := GetRequestID(r.Context())
		if id == "" {
			t.Error("request ID should be set in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Verify request ID is in response header
	responseID := rec.Header().Get("X-Request-ID")
	if responseID == "" {
		t.Error("X-Request-ID header should be set")
	}

	// Verify ID format (should be 16 hex chars)
	if len(responseID) != 16 {
		t.Errorf("request ID should be 16 chars, got %d", len(responseID))
	}
}

func TestRequestID_PreservesExistingID(t *testing.T) {
	existingID := "existing-request-id-123"

	var capturedID string
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", existingID)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should use the existing ID
	if capturedID != existingID {
		t.Errorf("got ID %q, want %q", capturedID, existingID)
	}

	// Response should have the same ID
	if rec.Header().Get("X-Request-ID") != existingID {
		t.Errorf("response header should contain %q", existingID)
	}
}

func TestGetRequestID_EmptyContext(t *testing.T) {
	ctx := context.Background()
	id := GetRequestID(ctx)

	if id != "" {
		t.Errorf("expected empty string for context without request ID, got %q", id)
	}
}

func TestRequestID_UniquePerRequest(t *testing.T) {
	ids := make(map[string]bool)

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Make multiple requests and verify IDs are unique
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		id := rec.Header().Get("X-Request-ID")
		if ids[id] {
			t.Errorf("duplicate request ID generated: %s", id)
		}
		ids[id] = true
	}
}
