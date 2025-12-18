package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSessionStore_CreateSession(t *testing.T) {
	store := NewSessionStore("test-secret")
	user := &User{ID: "user123", Email: "test@example.com", Name: "Test"}

	session, err := store.CreateSession(user)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	if session.ID == "" {
		t.Error("session ID should not be empty")
	}

	if session.User.ID != user.ID {
		t.Errorf("got user ID %q, want %q", session.User.ID, user.ID)
	}

	if session.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}

	if session.ExpiresAt.Before(time.Now()) {
		t.Error("ExpiresAt should be in the future")
	}
}

func TestSessionStore_GetSession(t *testing.T) {
	store := NewSessionStore("test-secret")
	user := &User{ID: "user123", Email: "test@example.com"}

	created, _ := store.CreateSession(user)

	retrieved, err := store.GetSession(created.ID)
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("got session ID %q, want %q", retrieved.ID, created.ID)
	}

	if retrieved.User.Email != user.Email {
		t.Errorf("got email %q, want %q", retrieved.User.Email, user.Email)
	}
}

func TestSessionStore_GetSession_NotFound(t *testing.T) {
	store := NewSessionStore("test-secret")

	_, err := store.GetSession("nonexistent-session-id")
	if err == nil {
		t.Error("expected error for nonexistent session")
	}
}

func TestSessionStore_DeleteSession(t *testing.T) {
	store := NewSessionStore("test-secret")
	user := &User{ID: "user123", Email: "test@example.com"}

	session, _ := store.CreateSession(user)

	// Verify it exists
	_, err := store.GetSession(session.ID)
	if err != nil {
		t.Fatalf("session should exist: %v", err)
	}

	// Delete it
	store.DeleteSession(session.ID)

	// Verify it's gone
	_, err = store.GetSession(session.ID)
	if err == nil {
		t.Error("session should be deleted")
	}
}

func TestSessionStore_SetAndGetSessionCookie(t *testing.T) {
	store := NewSessionStore("test-secret")
	user := &User{ID: "user123", Email: "test@example.com"}

	session, _ := store.CreateSession(user)

	// Set cookie
	rec := httptest.NewRecorder()
	store.SetSessionCookie(rec, session)

	// Extract cookie
	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie to be set")
	}

	cookie := cookies[0]
	if cookie.Name != "session" {
		t.Errorf("got cookie name %q, want %q", cookie.Name, "session")
	}

	if !cookie.HttpOnly {
		t.Error("cookie should be HttpOnly")
	}

	// Verify we can retrieve session from request with this cookie
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)

	retrieved, err := store.GetSessionFromRequest(req)
	if err != nil {
		t.Fatalf("failed to get session from request: %v", err)
	}

	if retrieved.ID != session.ID {
		t.Errorf("got session ID %q, want %q", retrieved.ID, session.ID)
	}
}

func TestSessionStore_GetSessionFromRequest_NoCookie(t *testing.T) {
	store := NewSessionStore("test-secret")

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := store.GetSessionFromRequest(req)
	if err == nil {
		t.Error("expected error when no cookie present")
	}
}

func TestSessionStore_GetSessionFromRequest_InvalidSignature(t *testing.T) {
	store := NewSessionStore("test-secret")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: "tampered-value.invalid-signature",
	})

	_, err := store.GetSessionFromRequest(req)
	if err == nil {
		t.Error("expected error for invalid signature")
	}
}

func TestSessionStore_ClearSessionCookie(t *testing.T) {
	store := NewSessionStore("test-secret")

	rec := httptest.NewRecorder()
	store.ClearSessionCookie(rec)

	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected cookie to be set for clearing")
	}

	cookie := cookies[0]
	if cookie.MaxAge != -1 {
		t.Errorf("MaxAge should be -1 to clear cookie, got %d", cookie.MaxAge)
	}
}

func TestSessionStore_SignatureVerification(t *testing.T) {
	store1 := NewSessionStore("secret-1")
	store2 := NewSessionStore("secret-2")

	user := &User{ID: "user123", Email: "test@example.com"}
	session, _ := store1.CreateSession(user)

	// Get cookie from store1
	rec := httptest.NewRecorder()
	store1.SetSessionCookie(rec, session)
	cookie := rec.Result().Cookies()[0]

	// Try to use it with store2 (different secret)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)

	_, err := store2.GetSessionFromRequest(req)
	if err == nil {
		t.Error("expected error when using cookie with different secret")
	}
}
