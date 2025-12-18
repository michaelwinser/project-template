package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.ServerURL != "http://localhost:8080" {
		t.Errorf("got ServerURL %q, want %q", cfg.ServerURL, "http://localhost:8080")
	}

	if cfg.SessionFile == "" {
		t.Error("SessionFile should not be empty")
	}

	if cfg.ConfigPath() == "" {
		t.Error("ConfigPath should not be empty")
	}
}

func TestLoadFromPath_NonexistentFile(t *testing.T) {
	cfg, err := LoadFromPath("/nonexistent/path/config.json")
	if err != nil {
		t.Fatalf("should not error for nonexistent file: %v", err)
	}

	// Should return defaults
	if cfg.ServerURL != "http://localhost:8080" {
		t.Errorf("got ServerURL %q, want default", cfg.ServerURL)
	}
}

func TestLoadFromPath_ValidFile(t *testing.T) {
	// Create a temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	configJSON := `{
		"server_url": "http://custom-server:9090"
	}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.ServerURL != "http://custom-server:9090" {
		t.Errorf("got ServerURL %q, want %q", cfg.ServerURL, "http://custom-server:9090")
	}

	if cfg.ConfigPath() != configPath {
		t.Errorf("got ConfigPath %q, want %q", cfg.ConfigPath(), configPath)
	}
}

func TestLoadFromPath_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(configPath, []byte("not valid json"), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := LoadFromPath(configPath)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDeriveSessionPath(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		want       string
	}{
		{
			name:       "standard cli config",
			configPath: "/home/user/.project-template-cli.json",
			want:       "/home/user/.project-template-session.json",
		},
		{
			name:       "custom config name",
			configPath: "/path/to/myproject.json",
			want:       "/path/to/myproject.session.json",
		},
		{
			name:       "config without json extension",
			configPath: "/path/to/config",
			want:       "/path/to/config.session.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deriveSessionPath(tt.configPath)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConfig_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	cfg, _ := LoadFromPath(configPath)
	cfg.ServerURL = "http://saved-server:1234"

	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Load it again
	loaded, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("failed to reload config: %v", err)
	}

	if loaded.ServerURL != cfg.ServerURL {
		t.Errorf("got ServerURL %q, want %q", loaded.ServerURL, cfg.ServerURL)
	}
}

func TestConfig_SessionOperations(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	cfg, _ := LoadFromPath(configPath)

	// Initially no session
	_, err := cfg.LoadSession()
	if err == nil {
		t.Error("expected error when no session exists")
	}

	// Save a session
	session := &Session{
		Cookie: "test-cookie-value",
		UserID: "user123",
		Email:  "test@example.com",
	}

	if err := cfg.SaveSession(session); err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// Load it back
	loaded, err := cfg.LoadSession()
	if err != nil {
		t.Fatalf("failed to load session: %v", err)
	}

	if loaded.Cookie != session.Cookie {
		t.Errorf("got Cookie %q, want %q", loaded.Cookie, session.Cookie)
	}

	if loaded.Email != session.Email {
		t.Errorf("got Email %q, want %q", loaded.Email, session.Email)
	}

	// Clear session
	if err := cfg.ClearSession(); err != nil {
		t.Fatalf("failed to clear session: %v", err)
	}

	// Should be gone
	_, err = cfg.LoadSession()
	if err == nil {
		t.Error("expected error after clearing session")
	}
}

func TestConfig_ClearSession_Nonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	cfg, _ := LoadFromPath(configPath)

	// Should not error when clearing nonexistent session
	if err := cfg.ClearSession(); err != nil {
		t.Errorf("ClearSession should not error for nonexistent file: %v", err)
	}
}
