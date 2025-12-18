package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project-template/server/internal/logging"
)

func TestLogsHandler_HandleUpload(t *testing.T) {
	var logBuf bytes.Buffer
	logger := logging.New("test")
	logger.SetOutput(&logBuf)

	handler := NewLogsHandler(logger)

	body := `{
		"entries": [
			{"level": "info", "message": "test log 1", "component": "ui"},
			{"level": "error", "message": "test log 2", "context": {"error": "something failed"}}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/logs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.HandleUpload(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var resp ClientLogResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Received != 2 {
		t.Errorf("got received %d, want %d", resp.Received, 2)
	}

	// Verify logs were written with source: client
	logLines := strings.Split(strings.TrimSpace(logBuf.String()), "\n")
	if len(logLines) != 2 {
		t.Errorf("expected 2 log entries, got %d", len(logLines))
	}

	var entry logging.Entry
	if err := json.Unmarshal([]byte(logLines[0]), &entry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}

	if entry.Source != "client" {
		t.Errorf("got source %q, want %q", entry.Source, "client")
	}

	if entry.Message != "test log 1" {
		t.Errorf("got message %q, want %q", entry.Message, "test log 1")
	}
}

func TestLogsHandler_MethodNotAllowed(t *testing.T) {
	logger := logging.New("test")
	handler := NewLogsHandler(logger)

	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/logs", nil)
			rec := httptest.NewRecorder()

			handler.HandleUpload(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestLogsHandler_InvalidJSON(t *testing.T) {
	logger := logging.New("test")
	handler := NewLogsHandler(logger)

	req := httptest.NewRequest(http.MethodPost, "/api/logs", strings.NewReader("not json"))
	rec := httptest.NewRecorder()

	handler.HandleUpload(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestLogsHandler_EmptyEntries(t *testing.T) {
	logger := logging.New("test")
	handler := NewLogsHandler(logger)

	req := httptest.NewRequest(http.MethodPost, "/api/logs", strings.NewReader(`{"entries": []}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.HandleUpload(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var resp ClientLogResponse
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.Received != 0 {
		t.Errorf("got received %d, want %d", resp.Received, 0)
	}
}
