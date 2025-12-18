package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	stateExpiration   = 10 * time.Minute
)

// OAuthProvider handles Google OAuth authentication
type OAuthProvider struct {
	config     *oauth2.Config
	states     map[string]time.Time
	statesMu   sync.RWMutex
	devMode    bool
}

// NewOAuthProvider creates a new OAuth provider
func NewOAuthProvider(clientID, clientSecret, redirectURL string, devMode bool) *OAuthProvider {
	provider := &OAuthProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		states:  make(map[string]time.Time),
		devMode: devMode,
	}

	// Clean up expired states periodically
	go provider.cleanupStates()

	return provider
}

// GetAuthURL generates the OAuth authorization URL
func (p *OAuthProvider) GetAuthURL() (string, error) {
	state, err := generateState()
	if err != nil {
		return "", err
	}

	p.statesMu.Lock()
	p.states[state] = time.Now().Add(stateExpiration)
	p.statesMu.Unlock()

	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// ExchangeCode exchanges the authorization code for user info
func (p *OAuthProvider) ExchangeCode(ctx context.Context, code, state string) (*User, error) {
	// In dev mode, accept any user
	if p.devMode && code == "dev" {
		return &User{
			ID:    "dev-user",
			Email: "dev@localhost",
			Name:  "Development User",
		}, nil
	}

	// Verify state
	p.statesMu.RLock()
	expiry, exists := p.states[state]
	p.statesMu.RUnlock()

	if !exists || time.Now().After(expiry) {
		return nil, errors.New("invalid or expired state")
	}

	// Delete used state
	p.statesMu.Lock()
	delete(p.states, state)
	p.statesMu.Unlock()

	// Exchange code for token
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// Get user info
	client := p.config.Client(ctx, token)
	resp, err := client.Get(googleUserInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, err
	}

	return &User{
		ID:      googleUser.ID,
		Email:   googleUser.Email,
		Name:    googleUser.Name,
		Picture: googleUser.Picture,
	}, nil
}

func (p *OAuthProvider) cleanupStates() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		p.statesMu.Lock()
		for state, expiry := range p.states {
			if now.After(expiry) {
				delete(p.states, state)
			}
		}
		p.statesMu.Unlock()
	}
}

func generateState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
