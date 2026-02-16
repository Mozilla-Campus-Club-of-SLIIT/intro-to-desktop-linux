package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

const (
	BaseURL      = "https://accounts.sliitmozilla.org/api"
	CallbackPort = "6767"
	RedirectURL  = "http://localhost:" + CallbackPort + "/"
)

// Session represents the cached authentication session
type Session struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	UserID       string    `json:"user_id"`
	UserEmail    string    `json:"user_email"`
	UserName     string    `json:"user_name"`
	Roles        []string  `json:"roles"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// TokenResponse represents the response from /api/token
type TokenResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

// SessionResponse represents the response from /api/session
type SessionResponse struct {
	Data struct {
		ID    string   `json:"id"`
		Roles []string `json:"roles"`
	} `json:"data"`
}

// UserResponse represents the response from /api/users/me
type UserResponse struct {
	Data struct {
		ID        string    `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		Roles     []string  `json:"roles"`
		CreatedAt time.Time `json:"createdAt"`
	} `json:"data"`
}

// AuthClient manages authentication state
type AuthClient struct {
	session    *Session
	sessionMu  sync.RWMutex
	httpClient *http.Client
}

var (
	globalClient     *AuthClient
	clientOnce       sync.Once
	isSessionHost    bool
	sessionHostMutex sync.RWMutex
)

// GetClient returns the singleton auth client
func GetClient() *AuthClient {
	clientOnce.Do(func() {
		globalClient = &AuthClient{
			httpClient: &http.Client{
				Timeout: 10 * time.Second,
			},
		}
		// Try to load existing session
		globalClient.loadSession()
	})
	return globalClient
}

// VerifyAuth checks if user is authenticated and refreshes token if needed
func VerifyAuth() bool {
	client := GetClient()
	client.sessionMu.RLock()

	if client.session == nil {
		client.sessionMu.RUnlock()
		return false
	}

	// Check if token is expired
	isExpired := time.Now().After(client.session.ExpiresAt)
	hasRefreshToken := client.session.RefreshToken != ""
	client.sessionMu.RUnlock()

	// If token is expired but we have a refresh token, try to refresh
	if isExpired && hasRefreshToken {
		// Try to refresh the token
		if err := client.RefreshToken(); err != nil {
			// Refresh failed, user needs to log in again
			return false
		}
		// Refresh succeeded
		return true
	}

	// Token is valid
	return !isExpired
}

// AuthUser initiates the OAuth flow
func AuthUser() error {
	client := GetClient()

	// Start callback server
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Create a new ServeMux to avoid conflicts with global mux
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + CallbackPort,
		Handler: mux,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no code received")
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`
				<!DOCTYPE html>
				<html>
				<head><title>Authentication Failed</title></head>
				<body>
					<h1>❌ Authentication Failed</h1>
					<p>No code received. This window will close automatically...</p>
					<script>setTimeout(function(){ window.close(); }, 2000);</script>
				</body>
				</html>
			`))
			return
		}

		codeChan <- code
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head><title>Authentication Successful</title></head>
			<body style="font-family: system-ui, -apple-system, sans-serif; text-align: center; padding: 50px;">
				<h1>✅ Authentication Successful!</h1>
				<p>You can now return to the application.</p>
				<p style="color: #666;">This window will close automatically...</p>
				<script>
					setTimeout(function(){
						window.close();
						// Fallback if window.close() doesn't work
						setTimeout(function() {
							document.body.innerHTML = '<h2>Please close this tab manually</h2>';
						}, 500);
					}, 1500);
				</script>
			</body>
			</html>
		`))
	})

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Open browser
	authURL := fmt.Sprintf("%s/authorize?redirect=%s", BaseURL, url.QueryEscape(RedirectURL))
	if err := openBrowser(authURL); err != nil {
		server.Shutdown(context.Background())
		return fmt.Errorf("failed to open browser: %w", err)
	}

	// Wait for code or error with timeout
	select {
	case code := <-codeChan:
		// Shutdown server
		server.Shutdown(context.Background())

		// Exchange code for tokens
		return client.exchangeCodeForToken(code)

	case err := <-errChan:
		server.Shutdown(context.Background())
		return err

	case <-time.After(5 * time.Minute):
		server.Shutdown(context.Background())
		return fmt.Errorf("authentication timeout")
	}
}

// exchangeCodeForToken exchanges the authorization code for access token
func (c *AuthClient) exchangeCodeForToken(code string) error {
	tokenURL := fmt.Sprintf("%s/token?code=%s", BaseURL, code)

	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Extract refresh token from cookie
	var refreshToken string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "refreshToken" {
			refreshToken = cookie.Value
			break
		}
	}

	// Fetch user info
	userInfo, err := c.fetchUserInfo(tokenResp.Data.Token)
	if err != nil {
		return fmt.Errorf("failed to fetch user info: %w", err)
	}

	// Create session
	c.sessionMu.Lock()
	c.session = &Session{
		AccessToken:  tokenResp.Data.Token,
		RefreshToken: refreshToken,
		UserID:       userInfo.Data.ID,
		UserEmail:    userInfo.Data.Email,
		UserName:     userInfo.Data.Name,
		Roles:        userInfo.Data.Roles,
		ExpiresAt:    time.Now().Add(15 * time.Minute), // Assume 15 min expiry
	}
	c.sessionMu.Unlock()

	setSessionHost(slices.Contains(userInfo.Data.Roles, "session-host"))

	// Save session to disk
	return c.saveSession()

}

// RefreshToken refreshes the access token using refresh token
func (c *AuthClient) RefreshToken() error {
	c.sessionMu.RLock()
	if c.session == nil || c.session.RefreshToken == "" {
		c.sessionMu.RUnlock()
		return fmt.Errorf("no refresh token available")
	}
	refreshToken := c.session.RefreshToken
	c.sessionMu.RUnlock()

	refreshURL := fmt.Sprintf("%s/token/refresh", BaseURL)

	req, err := http.NewRequest("POST", refreshURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add refresh token as cookie
	req.AddCookie(&http.Cookie{
		Name:  "refreshToken",
		Value: refreshToken,
	})

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		// If refresh fails, clear the session so user has to log in again
		if resp.StatusCode == http.StatusUnauthorized {
			c.Logout()
		}
		return fmt.Errorf("token refresh failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Check if a new refresh token was provided in the response cookies
	var newRefreshToken string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "refreshToken" {
			newRefreshToken = cookie.Value
			break
		}
	}

	// Update session with new token
	c.sessionMu.Lock()
	c.session.AccessToken = tokenResp.Data.Token
	c.session.ExpiresAt = time.Now().Add(15 * time.Minute)
	// Update refresh token if a new one was provided
	if newRefreshToken != "" {
		c.session.RefreshToken = newRefreshToken
	}
	c.sessionMu.Unlock()

	return c.saveSession()
}

// ensureValidToken checks if token is valid and refreshes if needed
func (c *AuthClient) ensureValidToken() error {
	c.sessionMu.RLock()
	isExpired := c.session != nil && time.Now().After(c.session.ExpiresAt)
	hasRefreshToken := c.session != nil && c.session.RefreshToken != ""
	c.sessionMu.RUnlock()

	if isExpired && hasRefreshToken {
		return c.RefreshToken()
	}
	return nil
}

// GetSession fetches the current session from the API
func (c *AuthClient) GetSession() (*SessionResponse, error) {
	// Ensure we have a valid token
	if err := c.ensureValidToken(); err != nil {
		return nil, fmt.Errorf("failed to ensure valid token: %w", err)
	}

	c.sessionMu.RLock()
	if c.session == nil {
		c.sessionMu.RUnlock()
		return nil, fmt.Errorf("not authenticated")
	}
	accessToken := c.session.AccessToken
	c.sessionMu.RUnlock()

	sessionURL := fmt.Sprintf("%s/session", BaseURL)

	req, err := http.NewRequest("GET", sessionURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get session failed (status %d): %s", resp.StatusCode, string(body))
	}

	var sessionResp SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
		return nil, fmt.Errorf("failed to parse session response: %w", err)
	}

	return &sessionResp, nil
}

// fetchUserInfo fetches user information from the API
func (c *AuthClient) fetchUserInfo(accessToken string) (*UserResponse, error) {
	userURL := fmt.Sprintf("%s/users/me", BaseURL)

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get user info failed (status %d): %s", resp.StatusCode, string(body))
	}

	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	return &userResp, nil
}

// GetCurrentUserInfo fetches current user info with auto token refresh
func (c *AuthClient) GetCurrentUserInfo() (*UserResponse, error) {
	// Ensure we have a valid token
	if err := c.ensureValidToken(); err != nil {
		return nil, fmt.Errorf("failed to ensure valid token: %w", err)
	}

	c.sessionMu.RLock()
	if c.session == nil {
		c.sessionMu.RUnlock()
		return nil, fmt.Errorf("not authenticated")
	}
	accessToken := c.session.AccessToken
	c.sessionMu.RUnlock()

	return c.fetchUserInfo(accessToken)
}

// GetUserInfo returns cached user information with roles
func (c *AuthClient) GetUserInfo() (string, string, []string, error) {
	c.sessionMu.RLock()
	defer c.sessionMu.RUnlock()

	if c.session == nil {
		return "", "", nil, fmt.Errorf("not authenticated")
	}

	return c.session.UserName, c.session.UserEmail, c.session.Roles, nil
}

// HasRole checks if the user has a specific role
func (c *AuthClient) HasRole(role string) bool {
	c.sessionMu.RLock()
	defer c.sessionMu.RUnlock()

	if c.session == nil {
		return false
	}

	return slices.Contains(c.session.Roles, role)

}

// IsSessionHost checks if user has session-host role
func IsSessionHost() bool {
	sessionHostMutex.RLock()
	defer sessionHostMutex.RUnlock()
	return isSessionHost
}

// setSessionHost sets the session-host flag (internal use)
func setSessionHost(value bool) {
	sessionHostMutex.Lock()
	defer sessionHostMutex.Unlock()
	isSessionHost = value
}

// saveSession saves the session to disk
func (c *AuthClient) saveSession() error {
	c.sessionMu.RLock()
	defer c.sessionMu.RUnlock()

	if c.session == nil {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	sessionDir := filepath.Join(homeDir, ".config", "intro-to-desktop-linux")
	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	sessionFile := filepath.Join(sessionDir, "session.json")
	data, err := json.Marshal(c.session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := os.WriteFile(sessionFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// loadSession loads the session from disk
func (c *AuthClient) loadSession() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	sessionFile := filepath.Join(homeDir, ".config", "intro-to-desktop-linux", "session.json")
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return err
	}

	c.sessionMu.Lock()
	c.session = &session
	c.sessionMu.Unlock()

	setSessionHost(slices.Contains(session.Roles, "session-host"))

	return nil

}

// Logout clears the session
func (c *AuthClient) Logout() error {
	c.sessionMu.Lock()
	c.session = nil
	c.sessionMu.Unlock()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	sessionFile := filepath.Join(homeDir, ".config", "intro-to-desktop-linux", "session.json")
	os.Remove(sessionFile)

	return nil
}

// openBrowser opens the default browser with the given URL
func openBrowser(url string) error {
	cmd := exec.Command("xdg-open", url)
	return cmd.Start()
}
