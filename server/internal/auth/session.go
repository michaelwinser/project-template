package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	sessionCookieName = "session"
	sessionDuration   = 24 * time.Hour
)

// User represents an authenticated user
type User struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionStore manages user sessions
type SessionStore struct {
	sessions map[string]*Session
	secret   []byte
	mu       sync.RWMutex
}

// NewSessionStore creates a new session store
func NewSessionStore(secret string) *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
		secret:   []byte(secret),
	}
}

// CreateSession creates a new session for a user
func (s *SessionStore) CreateSession(user *User) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session := &Session{
		ID:        sessionID,
		User:      user,
		CreatedAt: now,
		ExpiresAt: now.Add(sessionDuration),
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()

	return session, nil
}

// GetSession retrieves a session by ID
func (s *SessionStore) GetSession(sessionID string) (*Session, error) {
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists {
		return nil, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(sessionID)
		return nil, errors.New("session expired")
	}

	return session, nil
}

// DeleteSession removes a session
func (s *SessionStore) DeleteSession(sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
}

// SetSessionCookie sets the session cookie on the response
func (s *SessionStore) SetSessionCookie(w http.ResponseWriter, session *Session) {
	signedValue := s.signValue(session.ID)

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    signedValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	})
}

// ClearSessionCookie clears the session cookie
func (s *SessionStore) ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// GetSessionFromRequest retrieves the session from a request cookie
func (s *SessionStore) GetSessionFromRequest(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, errors.New("no session cookie")
	}

	sessionID, err := s.verifyValue(cookie.Value)
	if err != nil {
		return nil, err
	}

	return s.GetSession(sessionID)
}

func (s *SessionStore) signValue(value string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(value))
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	return value + "." + signature
}

func (s *SessionStore) verifyValue(signedValue string) (string, error) {
	parts := strings.SplitN(signedValue, ".", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid signed value")
	}

	value := parts[0]
	signature := parts[1]

	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(value))
	expectedSig := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return "", errors.New("invalid signature")
	}

	return value, nil
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// ToJSON serializes a user to JSON
func (u *User) ToJSON() ([]byte, error) {
	return json.Marshal(u)
}
