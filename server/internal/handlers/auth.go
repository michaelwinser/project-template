package handlers

import (
	"encoding/json"
	"net/http"

	"project-template/server/internal/auth"
	"project-template/server/internal/logging"
	"project-template/server/internal/middleware"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// LogoutResponse represents the logout response
type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	oauth    *auth.OAuthProvider
	sessions *auth.SessionStore
	logger   *logging.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(oauth *auth.OAuthProvider, sessions *auth.SessionStore, logger *logging.Logger) *AuthHandler {
	return &AuthHandler{
		oauth:    oauth,
		sessions: sessions,
		logger:   logger,
	}
}

// HandleLogin initiates the OAuth flow
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	requestID := middleware.GetRequestID(r.Context())
	log := h.logger.WithContext(requestID, "")

	authURL, err := h.oauth.GetAuthURL()
	if err != nil {
		log.Error("failed to generate auth URL", map[string]interface{}{"error": err.Error()})
		h.writeError(w, http.StatusInternalServerError, "auth_error", "Failed to initiate login")
		return
	}

	log.Info("redirecting to OAuth provider", nil)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// HandleCallback handles the OAuth callback
func (h *AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	requestID := middleware.GetRequestID(r.Context())
	log := h.logger.WithContext(requestID, "")

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		log.Warn("missing authorization code", nil)
		h.writeError(w, http.StatusBadRequest, "invalid_request", "Missing authorization code")
		return
	}

	user, err := h.oauth.ExchangeCode(r.Context(), code, state)
	if err != nil {
		log.Error("failed to exchange code", map[string]interface{}{"error": err.Error()})
		h.writeError(w, http.StatusUnauthorized, "auth_failed", "Authentication failed")
		return
	}

	session, err := h.sessions.CreateSession(user)
	if err != nil {
		log.Error("failed to create session", map[string]interface{}{"error": err.Error()})
		h.writeError(w, http.StatusInternalServerError, "session_error", "Failed to create session")
		return
	}

	h.sessions.SetSessionCookie(w, session)
	log.Info("user logged in", map[string]interface{}{"user_id": user.ID, "email": user.Email})

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// HandleLogout logs out the user
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	requestID := middleware.GetRequestID(r.Context())

	session, err := h.sessions.GetSessionFromRequest(r)
	if err != nil {
		log := h.logger.WithContext(requestID, "")
		log.Warn("logout attempt without valid session", map[string]interface{}{"error": err.Error()})
		h.writeError(w, http.StatusUnauthorized, "not_authenticated", "Not authenticated")
		return
	}

	log := h.logger.WithContext(requestID, session.User.ID)

	h.sessions.DeleteSession(session.ID)
	h.sessions.ClearSessionCookie(w)

	log.Info("user logged out", nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LogoutResponse{
		Success: true,
		Message: "Successfully logged out",
	})
}

// HandleMe returns the current user
func (h *AuthHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	session, err := h.sessions.GetSessionFromRequest(r)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "not_authenticated", "Not authenticated")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.User)
}

func (h *AuthHandler) writeError(w http.ResponseWriter, status int, errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errorCode,
		Message: message,
	})
}
