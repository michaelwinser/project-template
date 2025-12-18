package logging

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	logger := New("test-component")

	if logger.component != "test-component" {
		t.Errorf("got component %q, want %q", logger.component, "test-component")
	}

	if logger.minLevel != LevelDebug {
		t.Errorf("got minLevel %q, want %q", logger.minLevel, LevelDebug)
	}
}

func TestNewWithLevel(t *testing.T) {
	logger := NewWithLevel("test", LevelWarn)

	if logger.minLevel != LevelWarn {
		t.Errorf("got minLevel %q, want %q", logger.minLevel, LevelWarn)
	}
}

func TestNewFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     Level
	}{
		{"debug level", "debug", LevelDebug},
		{"info level", "info", LevelInfo},
		{"warn level", "warn", LevelWarn},
		{"warning level", "warning", LevelWarn},
		{"error level", "error", LevelError},
		{"uppercase", "INFO", LevelInfo},
		{"unknown defaults to debug", "unknown", LevelDebug},
		{"empty defaults to debug", "", LevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("LOG_LEVEL", tt.envValue)
				defer os.Unsetenv("LOG_LEVEL")
			} else {
				os.Unsetenv("LOG_LEVEL")
			}

			logger := NewFromEnv("test")
			if logger.minLevel != tt.want {
				t.Errorf("got minLevel %q, want %q", logger.minLevel, tt.want)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input string
		want  Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"warn", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
		{"invalid", LevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseLevel(tt.input)
			if got != tt.want {
				t.Errorf("ParseLevel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithLevel("test", LevelWarn)
	logger.SetOutput(&buf)

	// Debug and Info should be filtered
	logger.Debug("debug message", nil)
	logger.Info("info message", nil)

	if buf.Len() > 0 {
		t.Error("debug and info messages should be filtered at warn level")
	}

	// Warn and Error should pass through
	logger.Warn("warn message", nil)
	logger.Error("error message", nil)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 log lines, got %d", len(lines))
	}
}

func TestLoggerJSONOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := New("test-component")
	logger.SetOutput(&buf)

	logger.Info("test message", map[string]interface{}{
		"key": "value",
	})

	var entry Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log output as JSON: %v", err)
	}

	if entry.Level != LevelInfo {
		t.Errorf("got level %q, want %q", entry.Level, LevelInfo)
	}

	if entry.Message != "test message" {
		t.Errorf("got message %q, want %q", entry.Message, "test message")
	}

	if entry.Component != "test-component" {
		t.Errorf("got component %q, want %q", entry.Component, "test-component")
	}

	if entry.Timestamp == "" {
		t.Error("timestamp should be set")
	}

	if entry.Context["key"] != "value" {
		t.Errorf("got context[key] %v, want %q", entry.Context["key"], "value")
	}
}

func TestContextLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := New("test")
	logger.SetOutput(&buf)

	ctx := logger.WithContext("req-123", "user-456")
	ctx.Info("context message", nil)

	var entry Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if entry.RequestID != "req-123" {
		t.Errorf("got request_id %q, want %q", entry.RequestID, "req-123")
	}

	if entry.UserID != "user-456" {
		t.Errorf("got user_id %q, want %q", entry.UserID, "user-456")
	}
}

func TestLogEntry(t *testing.T) {
	var buf bytes.Buffer
	logger := New("test")
	logger.SetOutput(&buf)

	entry := Entry{
		Level:   LevelWarn,
		Message: "client log",
		Source:  "client",
		Context: map[string]interface{}{"browser": "chrome"},
	}

	logger.LogEntry(entry)

	var output Entry
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if output.Source != "client" {
		t.Errorf("got source %q, want %q", output.Source, "client")
	}

	if output.Timestamp == "" {
		t.Error("timestamp should be auto-filled")
	}
}

func TestLogEntryRespectLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithLevel("test", LevelError)
	logger.SetOutput(&buf)

	// This warn entry should be filtered
	entry := Entry{
		Level:   LevelWarn,
		Message: "should be filtered",
	}

	logger.LogEntry(entry)

	if buf.Len() > 0 {
		t.Error("warn entry should be filtered at error level")
	}
}
