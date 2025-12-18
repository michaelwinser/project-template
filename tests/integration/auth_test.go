package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// These tests run against a test server
// To run: go test -v ./tests/integration/...

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	serverURL := getServerURL()

	resp, err := http.Get(serverURL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var health map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if health["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%v'", health["status"])
	}
}

// TestUnauthenticatedAccess tests that unauthenticated users cannot access protected endpoints
func TestUnauthenticatedAccess(t *testing.T) {
	serverURL := getServerURL()

	tests := []struct {
		name     string
		method   string
		path     string
		wantCode int
	}{
		{"GET /auth/me without auth", "GET", "/auth/me", http.StatusUnauthorized},
		{"POST /auth/logout without auth", "POST", "/auth/logout", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, serverURL+tt.path, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("Expected status %d, got %d", tt.wantCode, resp.StatusCode)
			}
		})
	}
}

// TestLoginRedirect tests that /auth/login redirects to OAuth provider
func TestLoginRedirect(t *testing.T) {
	serverURL := getServerURL()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	resp, err := client.Get(serverURL + "/auth/login")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("Expected status 302, got %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if !strings.Contains(location, "accounts.google.com") {
		t.Errorf("Expected redirect to Google, got: %s", location)
	}
}

// getServerURL returns the server URL from environment or default
func getServerURL() string {
	if url := os.Getenv("TEST_SERVER_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

// Helper for testing with a mock server
func setupTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}
