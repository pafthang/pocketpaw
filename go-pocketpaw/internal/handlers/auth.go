package handlers

import (
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pocketpaw/pocketpaw/internal/config"
)

func generateSessionToken() string {
	return "ppat_" + uuid.New().String()
}

// AuthHandler handles authentication endpoints
type AuthHandler struct{}

// SessionTokenResponse represents session token response
type SessionTokenResponse struct {
	SessionToken   string `json:"session_token"`
	ExpiresInHours int    `json:"expires_in_hours"`
}

// TokenRegenerateResponse represents token regeneration response
type TokenRegenerateResponse struct {
	Token string `json:"token"`
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// RegisterRoutes registers auth routes
func (h *AuthHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/auth/session", h.ExchangeSessionToken)
	g.POST("/auth/login", h.CookieLogin)
	g.POST("/auth/logout", h.CookieLogout)
	g.GET("/qr", h.GetQRCode)
	g.POST("/token/regenerate", h.RegenerateAccessToken)
}

// ExchangeSessionToken exchanges master token for session token
func (h *AuthHandler) ExchangeSessionToken(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	var masterToken string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		masterToken = authHeader[7:]
	}

	currentToken, err := config.GetAccessToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if subtle.ConstantTimeCompare([]byte(masterToken), []byte(currentToken)) != 1 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"detail": "Invalid master token"})
	}

	settings, err := config.Load()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	ttlHours := settings.SessionTokenTTLHours
	if ttlHours <= 0 {
		ttlHours = 24
	}

	sessionToken := generateSessionToken()

	return c.JSON(http.StatusOK, SessionTokenResponse{
		SessionToken:   sessionToken,
		ExpiresInHours: ttlHours,
	})
}

// CookieLogin validates token and sets session cookie
func (h *AuthHandler) CookieLogin(c echo.Context) error {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"detail": "Invalid JSON body"})
	}

	submitted := req.Token
	currentToken, err := config.GetAccessToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	isValid := subtle.ConstantTimeCompare([]byte(submitted), []byte(currentToken)) == 1

	// Accept OAuth2 tokens (ppat_*)
	if !isValid && len(submitted) > 5 && submitted[:5] == "ppat_" {
		// TODO: Validate against OAuth server
		isValid = true
	}

	// Accept API keys (pp_*)
	if !isValid && len(submitted) > 3 && submitted[:3] == "pp_" && !(len(submitted) > 5 && submitted[:5] == "ppat_") {
		// TODO: Validate against API key manager
		isValid = true
	}

	if !isValid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"detail": "Invalid access token"})
	}

	settings, err := config.Load()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	sessionToken := generateSessionToken()
	maxAge := time.Duration(settings.SessionTokenTTLHours) * time.Hour

	cookie := &http.Cookie{
		Name:     "pocketpaw_session",
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   int(maxAge.Seconds()),
		Secure:   c.Request().TLS != nil,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]bool{"ok": true})
}

// CookieLogout clears session cookie
func (h *AuthHandler) CookieLogout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "pocketpaw_session",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]bool{"ok": true})
}

// GetQRCode generates QR login code
func (h *AuthHandler) GetQRCode(c echo.Context) error {
	// TODO: Implement QR code generation
	// Requires qrcode library integration
	return c.JSON(http.StatusNotImplemented, map[string]string{"detail": "Not implemented"})
}

// RegenerateAccessToken regenerates the access token
func (h *AuthHandler) RegenerateAccessToken(c echo.Context) error {
	newToken, err := config.RegenerateToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": newToken,
	})
}
