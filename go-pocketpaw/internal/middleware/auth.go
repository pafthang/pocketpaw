package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pocketpaw/pocketpaw/internal/config"
)

// AuthMiddleware handles authentication for API requests
type AuthMiddleware struct {
	GetAccessTokenFn func() (string, error)
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		GetAccessTokenFn: config.GetAccessToken,
	}
}

// Middleware returns the Echo middleware function
func (m *AuthMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Skip auth for certain paths
		skipPaths := []string{
			"/api/v1/docs",
			"/api/v1/redoc",
			"/api/v1/openapi.json",
			"/health",
			"/api/v1/health",
			"/api/v1/version",
		}

		path := c.Request().URL.Path
		for _, skip := range skipPaths {
			if strings.HasPrefix(path, skip) {
				return next(c)
			}
		}

		// Check for localhost access
		if isLocalhost(c.Request()) {
			c.Set("full_access", true)
			return next(c)
		}

		// Get authorization header
		authHeader := c.Request().Header.Get("Authorization")
		var token string

		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Check cookie if no bearer token
		if token == "" {
			cookie, err := c.Cookie("pocketpaw_session")
			if err == nil && cookie != nil {
				token = cookie.Value
			}
		}

		// Validate token
		if token != "" {
			masterToken, err := m.GetAccessTokenFn()
			if err == nil {
				// Check master token
				if subtle.ConstantTimeCompare([]byte(token), []byte(masterToken)) == 1 {
					c.Set("full_access", true)
					return next(c)
				}

				// Check session token (ppat_*)
				if strings.HasPrefix(token, "ppat_") {
					// Session tokens start with ppat_ - validate against stored sessions
					// For now, accept any valid session token format
					c.Set("full_access", true)
					return next(c)
				}

				// Check API key (pp_*)
				if strings.HasPrefix(token, "pp_") && !strings.HasPrefix(token, "ppat_") {
					// API key validation would go here
					// For now, accept valid API key format
					c.Set("api_key_validated", true)
					return next(c)
				}
			}
		}

		// No valid auth found - return 401 for protected routes
		// Allow OPTIONS for CORS preflight
		if c.Request().Method == http.MethodOptions {
			return next(c)
		}

		return echo.ErrUnauthorized
	}
}

// isLocalhost checks if the request is from localhost
func isLocalhost(r *http.Request) bool {
	host := r.RemoteAddr
	// Strip port if present
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	return host == "127.0.0.1" || host == "::1" || host == "localhost"
}
