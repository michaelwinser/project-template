package config

import (
	"os"
)

// Config holds all configuration for the server
type Config struct {
	Port               string
	Host               string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	SessionSecret      string
	LogLevel           string
	LogFormat          string
	Env                string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		Host:               getEnv("HOST", "localhost"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/callback"),
		SessionSecret:      getEnv("SESSION_SECRET", "dev-secret-change-in-production"),
		LogLevel:           getEnv("LOG_LEVEL", "debug"),
		LogFormat:          getEnv("LOG_FORMAT", "json"),
		Env:                getEnv("ENV", "development"),
	}
}

// Address returns the full address to listen on
func (c *Config) Address() string {
	return c.Host + ":" + c.Port
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
