package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pocketpaw/pocketpaw/internal/services/sessions"
)

// SessionsHandler handles session-related HTTP requests
type SessionsHandler struct {
	service *sessions.SessionService
}

// NewSessionsHandler creates a new sessions handler
func NewSessionsHandler(service *sessions.SessionService) *SessionsHandler {
	return &SessionsHandler{service: service}
}

// RegisterRoutes registers session routes on the Echo router
func (h *SessionsHandler) RegisterRoutes(g *echo.Group) {
	g.POST("", h.CreateSession)
	g.GET("", h.ListSessions)
	g.DELETE("/:id", h.DeleteSession)
	g.POST("/:id/title", h.UpdateSessionTitle)
	g.GET("/search", h.SearchSessions)
	g.GET("/:id/history", h.GetSessionHistory)
	g.GET("/:id/export", h.ExportSession)
}

// CreateSession godoc
// @Summary Create a new empty session
// @Tags Sessions
// @Accept json
// @Produce json
// @Success 200 {object} sessions.SessionCreateResponse
// @Router /api/v1/sessions [post]
func (h *SessionsHandler) CreateSession(c echo.Context) error {
	resp, err := h.service.CreateSession()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

// ListSessions godoc
// @Summary List sessions
// @Tags Sessions
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Success 200 {object} sessions.SessionListResponse
// @Router /api/v1/sessions [get]
func (h *SessionsHandler) ListSessions(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	limit := 50
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
			if limit < 1 {
				limit = 1
			} else if limit > 500 {
				limit = 500
			}
		}
	}

	resp, err := h.service.ListSessions(limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

// DeleteSession godoc
// @Summary Delete a session
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/sessions/{id} [delete]
func (h *SessionsHandler) DeleteSession(c echo.Context) error {
	sessionID := c.Param("id")
	if err := h.service.DeleteSession(sessionID); err != nil {
		if err == sessions.ErrSessionNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"detail": "Session not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": true})
}

// UpdateSessionTitle godoc
// @Summary Update session title
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param body body sessions.SessionTitleRequest true "Title update request"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/sessions/{id}/title [post]
func (h *SessionsHandler) UpdateSessionTitle(c echo.Context) error {
	sessionID := c.Param("id")

	var req sessions.SessionTitleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"detail": "Invalid JSON body"})
	}

	if err := h.service.UpdateSessionTitle(sessionID, req.Title); err != nil {
		if err == sessions.ErrSessionNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"detail": "Session not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": true})
}

// SearchSessions godoc
// @Summary Search sessions
// @Tags Sessions
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} sessions.SessionSearchResponse
// @Router /api/v1/sessions/search [get]
func (h *SessionsHandler) SearchSessions(c echo.Context) error {
	query := c.QueryParam("q")
	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
			if limit < 1 {
				limit = 1
			} else if limit > 200 {
				limit = 200
			}
		}
	}

	resp, err := h.service.SearchSessions(query, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

// GetSessionHistory godoc
// @Summary Get session history
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param limit query int false "Limit" default(50)
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/sessions/{id}/history [get]
func (h *SessionsHandler) GetSessionHistory(c echo.Context) error {
	sessionID := c.Param("id")
	limitStr := c.QueryParam("limit")
	limit := 50
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
			if limit < 1 {
				limit = 1
			} else if limit > 500 {
				limit = 500
			}
		}
	}

	if sessionID == "" {
		return c.JSON(http.StatusOK, []map[string]interface{}{})
	}

	history, err := h.service.GetSessionHistory(sessionID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, history)
}

// ExportSession godoc
// @Summary Export session
// @Tags Sessions
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param format query string false "Export format" default(json)
// @Success 200 {string} string
// @Router /api/v1/sessions/{id}/export [get]
func (h *SessionsHandler) ExportSession(c echo.Context) error {
	sessionID := c.Param("id")
	format := c.QueryParam("format")
	if format == "" {
		format = "json"
	}

	if format != "json" && format != "md" {
		return c.JSON(http.StatusBadRequest, map[string]string{"detail": "Format must be 'json' or 'md'"})
	}

	data, mediaType, err := h.service.ExportSession(sessionID, format)
	if err != nil {
		if err == sessions.ErrSessionNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"detail": "Session not found: " + sessionID})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	filename := "pocketpaw-session-" + sessionID[:20] + "." + format
	c.Response().Header().Set(echo.HeaderContentType, mediaType)
	c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="`+filename+`"`)
	return c.Blob(http.StatusOK, mediaType, data)
}
