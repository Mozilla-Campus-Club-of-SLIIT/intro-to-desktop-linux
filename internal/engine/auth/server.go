package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine"
)

const (
	authServiceHost = "https://accounts.sliitmozilla.org"
	authServicePath = "/api/session"
)

// AuthService provides authentication related functionalities.
type AuthService struct {
	db         *engine.DB
	httpClient *http.Client
}

// NewAuthService creates a new AuthService.
func NewAuthService(db *engine.DB) *AuthService {
	return &AuthService{
		db:         db,
		httpClient: &http.Client{Timeout: 5 * time.Second}, // TODO: Make timeout configurable
	}
}

// VerifySession verifies a user's session using Valkey cache and external authentication service.
func (s *AuthService) VerifySession(ctx context.Context, userID, accessToken string) (bool, error) {
	// 1. Check Valkey cache
	session, err := s.db.GetSession(ctx, userID)
	if err != nil && err != engine.ErrSessionNotFound {
		// Log error, but proceed to external verification if it's not a "not found" error
		// In a real application, you might want more robust error handling/logging here.
		fmt.Printf("Error getting session from Valkey for user %s: %v\n", userID, err)
	}

	if session != nil {
		expiresAtStr, ok := session["expiresAt"] // Use "expiresAt" as the key for the timestamp in the hash
		if ok {
			expiresAt, parseErr := time.Parse(time.RFC3339, expiresAtStr)
			if parseErr != nil {
				fmt.Printf("Error parsing expiresAt for user %s: %v. Proceeding to external verification.\n", userID, parseErr)
			} else if time.Now().Before(expiresAt) {
				// Session is valid and not expired
				return true, nil
			}
		}
	}

	// 2. External API verification
	req, err := http.NewRequestWithContext(ctx, "GET", authServiceHost+authServicePath, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create auth request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Unauthorized, forbidden, or other error from auth service
		return false, fmt.Errorf("auth service returned non-OK status: %d", resp.StatusCode)
	}

	var authResponse struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return false, fmt.Errorf("failed to decode auth service response: %w", err)
	}

	if authResponse.Data.ID == userID {
		// User verified, update session in cache for 10 minutes
		if err := s.db.SetSession(ctx, userID, userID, 10*time.Minute); err != nil {
			fmt.Printf("Error setting session in Valkey for user %s: %v\n", userID, err)
		}
		return true, nil
	}

	return false, nil // User ID mismatch
}
