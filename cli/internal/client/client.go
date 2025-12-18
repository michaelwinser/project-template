package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// User represents an authenticated user
type User struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version,omitempty"`
}

// LogoutResponse represents the logout response
type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ApiError represents an API error with status code
type ApiError struct {
	StatusCode int
	ErrorCode  string
	Message    string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Message)
}

// Client is the API client for the CLI
type Client struct {
	baseURL    string
	httpClient *http.Client
	cookie     string
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{},
	}
}

// SetCookie sets the session cookie for authenticated requests
func (c *Client) SetCookie(cookie string) {
	c.cookie = cookie
}

// GetHealth checks the server health
func (c *Client) GetHealth() (*HealthResponse, error) {
	var response HealthResponse
	if err := c.request("GET", "/health", nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetLoginURL returns the OAuth login URL
func (c *Client) GetLoginURL() string {
	return c.baseURL + "/auth/login"
}

// GetCurrentUser gets the current authenticated user
func (c *Client) GetCurrentUser() (*User, error) {
	var user User
	if err := c.request("GET", "/auth/me", nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Logout logs out the current user
func (c *Client) Logout() (*LogoutResponse, error) {
	var response LogoutResponse
	if err := c.request("POST", "/auth/logout", nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) request(method, path string, body io.Reader, result interface{}) error {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.cookie != "" {
		req.Header.Set("Cookie", "session="+c.cookie)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return &ApiError{
				StatusCode: resp.StatusCode,
				ErrorCode:  "unknown_error",
				Message:    resp.Status,
			}
		}
		return &ApiError{
			StatusCode: resp.StatusCode,
			ErrorCode:  errResp.Error,
			Message:    errResp.Message,
		}
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return err
		}
	}

	return nil
}

// ExtractSessionCookie extracts the session cookie from a Set-Cookie header
func ExtractSessionCookie(setCookie string) (string, error) {
	parts := strings.Split(setCookie, ";")
	if len(parts) == 0 {
		return "", errors.New("invalid cookie")
	}

	cookiePart := parts[0]
	if !strings.HasPrefix(cookiePart, "session=") {
		return "", errors.New("session cookie not found")
	}

	return strings.TrimPrefix(cookiePart, "session="), nil
}
