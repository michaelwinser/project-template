package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"project-template/server/internal/auth"
	"project-template/server/internal/logging"
)

func newTestAuthHandler() (*AuthHandler, *auth.SessionStore) {
	sessions := auth.NewSessionStore("test-secret-key")
	oauth := auth.NewOAuthProvider("test-client-id", "test-client-secret", "http://localhost/callback", true)
	logger := logging.New("test")
	return NewAuthHandler(oauth, sessions, logger), sessions
}

func TestAuthHandler_HandleMe_Unauthenticated(t *testing.T) {
	handler, _ := newTestAuthHandler()

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()

	handler.HandleMe(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != "not_authenticated" {
		t.Errorf("got error %q, want %q", resp.Error, "not_authenticated")
	}
}

func TestAuthHandler_HandleMe_Authenticated(t *testing.T) {
	handler, sessions := newTestAuthHandler()

	// Create a session
	user := &auth.User{ID: "user123", Email: "test@example.com", Name: "Test User"}
	session, err := sessions.CreateSession(user)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	// Create request with session cookie
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()

	// Set the session cookie (need to sign it properly)
	sessions.SetSessionCookie(rec, session)
	cookie := rec.Result().Cookies()[0]
	req.AddCookie(cookie)

	// Reset recorder for actual test
	rec = httptest.NewRecorder()
	handler.HandleMe(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var respUser auth.User
	if err := json.NewDecoder(rec.Body).Decode(&respUser); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if respUser.Email != user.Email {
		t.Errorf("got email %q, want %q", respUser.Email, user.Email)
	}
}

func TestAuthHandler_HandleMe_MethodNotAllowed(t *testing.T) {
	handler, _ := newTestAuthHandler()

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/auth/me", nil)
			rec := httptest.NewRecorder()

			handler.HandleMe(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestAuthHandler_HandleLogout_Unauthenticated(t *testing.T) {
	handler, _ := newTestAuthHandler()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rec := httptest.NewRecorder()

	handler.HandleLogout(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_HandleLogout_Authenticated(t *testing.T) {
	handler, sessions := newTestAuthHandler()

	// Create a session
	user := &auth.User{ID: "user123", Email: "test@example.com"}
	session, _ := sessions.CreateSession(user)

	// Get signed cookie
	rec := httptest.NewRecorder()
	sessions.SetSessionCookie(rec, session)
	cookie := rec.Result().Cookies()[0]

	// Make logout request
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()

	handler.HandleLogout(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var resp LogoutResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}

	// Verify session was deleted
	_, err := sessions.GetSession(session.ID)
	if err == nil {
		t.Error("session should have been deleted")
	}
}

func TestAuthHandler_HandleLogin_Redirect(t *testing.T) {
	handler, _ := newTestAuthHandler()

	req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
	rec := httptest.NewRecorder()

	handler.HandleLogin(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusFound)
	}

	location := rec.Header().Get("Location")
	if location == "" {
		t.Error("expected Location header to be set")
	}
}
