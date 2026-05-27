package auth

import (
	"crypto/subtle"
	"time"

	"github.com/google/uuid"
	"github.com/pocketpaw/pocketpaw/internal/config"
)

// AuthService handles authentication operations
type AuthService struct{}

// SessionTokenResponse represents the response for session token exchange
type SessionTokenResponse struct {
	SessionToken   string `json:"session_token"`
	ExpiresInHours int    `json:"expires_in_hours"`
}

// TokenRegenerateResponse represents the response for token regeneration
type TokenRegenerateResponse struct {
	Token string `json:"token"`
}

// ExchangeSessionToken exchanges a master token for a session token
func (s *AuthService) ExchangeSessionToken(masterToken string) (*SessionTokenResponse, error) {
	currentToken, err := config.GetAccessToken()
	if err != nil {
		return nil, err
	}

	// Validate master token using constant-time comparison
	if subtle.ConstantTimeCompare([]byte(masterToken), []byte(currentToken)) != 1 {
		return nil, ErrInvalidToken
	}

	settings, err := config.Load()
	if err != nil {
		return nil, err
	}

	ttlHours := settings.SessionTokenTTLHours
	if ttlHours <= 0 {
		ttlHours = 24
	}

	sessionToken := generateSessionToken()

	return &SessionTokenResponse{
		SessionToken:   sessionToken,
		ExpiresInHours: ttlHours,
	}, nil
}

// ValidateToken validates a token and returns whether it's valid
func (s *AuthService) ValidateToken(token string) bool {
	if token == "" {
		return false
	}

	masterToken, err := config.GetAccessToken()
	if err != nil {
		return false
	}

	// Check master token
	if subtle.ConstantTimeCompare([]byte(token), []byte(masterToken)) == 1 {
		return true
	}

	// Check session token format
	if len(token) > 5 && token[:5] == "ppat_" {
		// Valid session token format - in production, validate against stored sessions
		return true
	}

	// Check API key format
	if len(token) > 3 && token[:3] == "pp_" && !(len(token) > 5 && token[:5] == "ppat_") {
		// Valid API key format
		return true
	}

	return false
}

// RegenerateAccessToken generates a new access token
func (s *AuthService) RegenerateAccessToken() (string, error) {
	newToken, err := config.RegenerateToken()
	if err != nil {
		return "", err
	}

	return newToken, nil
}

// CreateSessionToken creates a new session token with the specified TTL
func (s *AuthService) CreateSessionToken(ttlHours int) string {
	return generateSessionToken()
}

func generateSessionToken() string {
	return "ppat_" + uuid.New().String()
}

// GetTokenExpiry calculates the expiry time for a session token
func (s *AuthService) GetTokenExpiry(ttlHours int) time.Time {
	return time.Now().Add(time.Duration(ttlHours) * time.Hour)
}

// Error definitions
var (
	ErrInvalidToken = &AuthError{Message: "Invalid token"}
)

// AuthError represents an authentication error
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// NewAuthService creates a new auth service instance
func NewAuthService() *AuthService {
	return &AuthService{}
}
