package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"project-template/server/internal/auth"
	"project-template/server/internal/config"
	"project-template/server/internal/handlers"
	"project-template/server/internal/logging"
	"project-template/server/internal/middleware"
)

const version = "0.1.0"

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logging.New("server")
	if cfg.LogLevel == "info" {
		logger.SetLevel(logging.LevelInfo)
	}

	logger.Info("starting server", map[string]interface{}{
		"version": version,
		"address": cfg.Address(),
		"env":     cfg.Env,
	})

	// Initialize auth components
	oauth := auth.NewOAuthProvider(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
		cfg.IsDevelopment(),
	)
	sessions := auth.NewSessionStore(cfg.SessionSecret)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(version)
	authHandler := handlers.NewAuthHandler(oauth, sessions, logging.New("auth"))

	// Set up routes
	mux := http.NewServeMux()

	// Health check
	mux.Handle("/health", healthHandler)

	// Auth routes
	mux.HandleFunc("/auth/login", authHandler.HandleLogin)
	mux.HandleFunc("/auth/callback", authHandler.HandleCallback)
	mux.HandleFunc("/auth/logout", authHandler.HandleLogout)
	mux.HandleFunc("/auth/me", authHandler.HandleMe)

	// Apply middleware
	handler := middleware.RequestID(mux)
	handler = middleware.Logging(logger)(handler)

	// Start server
	server := &http.Server{
		Addr:    cfg.Address(),
		Handler: handler,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", map[string]interface{}{"error": err.Error()})
			os.Exit(1)
		}
	}()

	logger.Info("server started", map[string]interface{}{"address": cfg.Address()})

	<-done
	logger.Info("shutting down server", nil)
}
