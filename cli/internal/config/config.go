package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const defaultConfigFileName = ".project-template-cli.json"

// Config holds CLI configuration
type Config struct {
	ServerURL   string `json:"server_url"`
	SessionFile string `json:"session_file"`
	configPath  string // Path to the config file itself (not serialized)
}

// Session holds the stored session data
type Session struct {
	Cookie string `json:"cookie"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		ServerURL:   "http://localhost:8080",
		SessionFile: filepath.Join(homeDir, ".project-template-session.json"),
		configPath:  filepath.Join(homeDir, defaultConfigFileName),
	}
}

// Load loads configuration from the default config file
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return DefaultConfig(), nil
	}
	return LoadFromPath(filepath.Join(homeDir, defaultConfigFileName))
}

// LoadFromPath loads configuration from a specific path
func LoadFromPath(configPath string) (*Config, error) {
	// Determine session file path - same directory as config, with .session.json suffix
	sessionFile := deriveSessionPath(configPath)

	config := &Config{
		ServerURL:   "http://localhost:8080",
		SessionFile: sessionFile,
		configPath:  configPath,
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config file doesn't exist - return defaults with derived paths
		return config, nil
	}

	if err := json.Unmarshal(data, config); err != nil {
		return config, err
	}

	// If session file wasn't specified in config, use derived path
	if config.SessionFile == "" {
		config.SessionFile = sessionFile
	}

	config.configPath = configPath
	return config, nil
}

// deriveSessionPath creates a session file path based on the config path
// e.g., /path/to/myproject.json -> /path/to/myproject.session.json
// e.g., /path/to/.project-template-cli.json -> /path/to/.project-template-session.json
func deriveSessionPath(configPath string) string {
	dir := filepath.Dir(configPath)
	base := filepath.Base(configPath)

	// Remove .json extension if present
	base = strings.TrimSuffix(base, ".json")

	// Remove -cli suffix if present and add -session
	if strings.HasSuffix(base, "-cli") {
		base = strings.TrimSuffix(base, "-cli") + "-session"
	} else {
		base = base + ".session"
	}

	return filepath.Join(dir, base+".json")
}

// ConfigPath returns the path to the config file
func (c *Config) ConfigPath() string {
	return c.configPath
}

// Save saves configuration to the config file
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.configPath, data, 0600)
}

// LoadSession loads the stored session
func (c *Config) LoadSession() (*Session, error) {
	data, err := os.ReadFile(c.SessionFile)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// SaveSession saves a session to the session file
func (c *Config) SaveSession(session *Session) error {
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.SessionFile, data, 0600)
}

// ClearSession removes the session file
func (c *Config) ClearSession() error {
	if err := os.Remove(c.SessionFile); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
