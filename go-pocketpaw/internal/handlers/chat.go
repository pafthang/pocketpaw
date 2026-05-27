package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pocketpaw/pocketpaw/internal/services/chat"
)

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	service *chat.ChatService
}

// NewChatHandler creates a new chat handler
func NewChatHandler(service *chat.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

// RegisterRoutes registers chat routes on the Echo router
func (h *ChatHandler) RegisterRoutes(g *echo.Group) {
	g.POST("", h.SendChat)
	g.POST("/stream", h.StreamChat)
	g.POST("/stop", h.StopChat)
}

// SendChat godoc
// @Summary Send a message and get complete response
// @Tags Chat
// @Accept json
// @Produce json
// @Param body body chat.ChatRequest true "Chat request"
// @Success 200 {object} chat.ChatResponse
// @Router /api/v1/chat [post]
func (h *ChatHandler) SendChat(c echo.Context) error {
	var req chat.ChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"detail": "Invalid JSON body"})
	}

	resp, err := h.service.GetChatResponse(&req)
	if err != nil {
		if chatErr, ok := err.(*chat.ChatError); ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"detail": chatErr.Message})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

// StreamChat godoc
// @Summary Send a message and receive SSE stream
// @Tags Chat
// @Accept json
// @Produce text/event-stream
// @Param body body chat.ChatRequest true "Chat request"
// @Success 200 {string} string
// @Router /api/v1/chat/stream [post]
func (h *ChatHandler) StreamChat(c echo.Context) error {
	var req chat.ChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"detail": "Invalid JSON body"})
	}

	safeKey, err := h.service.StreamChat(&req)
	if err != nil {
		if chatErr, ok := err.(*chat.ChatError); ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"detail": chatErr.Message})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Set SSE headers
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")

	// Send initial event
	initialEvent := `event: stream_start` + "\ndata: " + `{"session_id":"` + safeKey + `"}` + "\n\n"
	if _, err := c.Response().Write([]byte(initialEvent)); err != nil {
		return err
	}
	c.Response().Flush()

	// In production, this would stream events from the service
	// For now, just send a stream_end event
	endEvent := `event: stream_end` + "\ndata: " + `{"session_id":"` + safeKey + `","usage":{}}` + "\n\n"
	if _, err := c.Response().Write([]byte(endEvent)); err != nil {
		return err
	}
	c.Response().Flush()

	return nil
}

// StopChat godoc
// @Summary Cancel an in-flight chat
// @Tags Chat
// @Accept json
// @Produce json
// @Param session_id query string true "Session ID"
// @Success 200 {object} chat.StopResponse
// @Router /api/v1/chat/stop [post]
func (h *ChatHandler) StopChat(c echo.Context) error {
	sessionID := c.QueryParam("session_id")
	
	resp, err := h.service.StopChat(sessionID)
	if err != nil {
		if chatErr, ok := err.(*chat.ChatError); ok {
			if chatErr.Message == "session_id is required" {
				return c.JSON(http.StatusBadRequest, map[string]string{"detail": chatErr.Message})
			}
			if chatErr.Message == "No active stream for this session" {
				return c.JSON(http.StatusNotFound, map[string]string{"detail": chatErr.Message})
			}
			return c.JSON(http.StatusBadRequest, map[string]string{"detail": chatErr.Message})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}
