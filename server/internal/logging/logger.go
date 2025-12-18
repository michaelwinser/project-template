package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Level represents the log level
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Entry represents a structured log entry
type Entry struct {
	Timestamp string                 `json:"timestamp"`
	Level     Level                  `json:"level"`
	Message   string                 `json:"message"`
	Source    string                 `json:"source,omitempty"`
	Component string                 `json:"component,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// Logger provides structured logging
type Logger struct {
	component string
	minLevel  Level
	output    io.Writer
}

// New creates a new logger for a component
func New(component string) *Logger {
	return &Logger{
		component: component,
		minLevel:  LevelDebug,
		output:    os.Stdout,
	}
}

// NewWithLevel creates a logger with a specific level
func NewWithLevel(component string, level Level) *Logger {
	return &Logger{
		component: component,
		minLevel:  level,
		output:    os.Stdout,
	}
}

// NewFromEnv creates a logger with level from LOG_LEVEL environment variable
func NewFromEnv(component string) *Logger {
	level := LevelDebug // Default to debug
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		level = ParseLevel(envLevel)
	}
	return NewWithLevel(component, level)
}

// ParseLevel converts a string to a Level, defaulting to debug for unknown values
func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelDebug
	}
}

// SetOutput sets the output writer (useful for testing)
func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.minLevel = level
}

// WithContext returns a context-aware log function
func (l *Logger) WithContext(requestID, userID string) *ContextLogger {
	return &ContextLogger{
		logger:    l,
		requestID: requestID,
		userID:    userID,
	}
}

func (l *Logger) log(level Level, message string, ctx map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	entry := Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Message:   message,
		Component: l.component,
		Context:   ctx,
	}

	l.writeEntry(entry)
}

// LogEntry writes a pre-formed entry (used for client log forwarding)
func (l *Logger) LogEntry(entry Entry) {
	if !l.shouldLog(entry.Level) {
		return
	}
	// Ensure timestamp if not set
	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	}
	l.writeEntry(entry)
}

func (l *Logger) writeEntry(entry Entry) {
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal log entry: %v\n", err)
		return
	}

	fmt.Fprintln(l.output, string(data))
}

func (l *Logger) shouldLog(level Level) bool {
	levels := map[Level]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
	}
	return levels[level] >= levels[l.minLevel]
}

// Debug logs a debug message
func (l *Logger) Debug(message string, ctx map[string]interface{}) {
	l.log(LevelDebug, message, ctx)
}

// Info logs an info message
func (l *Logger) Info(message string, ctx map[string]interface{}) {
	l.log(LevelInfo, message, ctx)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, ctx map[string]interface{}) {
	l.log(LevelWarn, message, ctx)
}

// Error logs an error message
func (l *Logger) Error(message string, ctx map[string]interface{}) {
	l.log(LevelError, message, ctx)
}

// ContextLogger is a logger with request context
type ContextLogger struct {
	logger    *Logger
	requestID string
	userID    string
}

func (cl *ContextLogger) log(level Level, message string, ctx map[string]interface{}) {
	if !cl.logger.shouldLog(level) {
		return
	}

	entry := Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Message:   message,
		Component: cl.logger.component,
		RequestID: cl.requestID,
		UserID:    cl.userID,
		Context:   ctx,
	}

	cl.logger.writeEntry(entry)
}

// Debug logs a debug message with context
func (cl *ContextLogger) Debug(message string, ctx map[string]interface{}) {
	cl.log(LevelDebug, message, ctx)
}

// Info logs an info message with context
func (cl *ContextLogger) Info(message string, ctx map[string]interface{}) {
	cl.log(LevelInfo, message, ctx)
}

// Warn logs a warning message with context
func (cl *ContextLogger) Warn(message string, ctx map[string]interface{}) {
	cl.log(LevelWarn, message, ctx)
}

// Error logs an error message with context
func (cl *ContextLogger) Error(message string, ctx map[string]interface{}) {
	cl.log(LevelError, message, ctx)
}
