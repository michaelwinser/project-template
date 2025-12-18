package handlers

import (
	"encoding/json"
	"net/http"

	"project-template/server/internal/logging"
	"project-template/server/internal/middleware"
)

// ClientLogEntry represents a log entry from a client
type ClientLogEntry struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Component string                 `json:"component,omitempty"`
	Timestamp string                 `json:"timestamp,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// ClientLogRequest represents a batch of client log entries
type ClientLogRequest struct {
	Entries []ClientLogEntry `json:"entries"`
}

// ClientLogResponse represents the response to a log upload
type ClientLogResponse struct {
	Received int `json:"received"`
}

// LogsHandler handles client log uploads
type LogsHandler struct {
	logger *logging.Logger
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler(logger *logging.Logger) *LogsHandler {
	return &LogsHandler{logger: logger}
}

// HandleUpload handles POST /api/logs
func (h *LogsHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	requestID := middleware.GetRequestID(r.Context())

	// Limit request body size (100KB max for logs)
	r.Body = http.MaxBytesReader(w, r.Body, 100*1024)

	var req ClientLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	// Validate and forward each entry
	for _, entry := range req.Entries {
		level := logging.ParseLevel(entry.Level)

		logEntry := logging.Entry{
			Timestamp: entry.Timestamp,
			Level:     level,
			Message:   entry.Message,
			Source:    "client",
			Component: entry.Component,
			RequestID: requestID,
			Context:   entry.Context,
		}

		h.logger.LogEntry(logEntry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ClientLogResponse{
		Received: len(req.Entries),
	})
}

func (h *LogsHandler) writeError(w http.ResponseWriter, status int, errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errorCode,
		Message: message,
	})
}
