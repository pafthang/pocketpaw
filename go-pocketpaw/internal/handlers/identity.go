package handlers

import (
"net/http"

"github.com/labstack/echo/v4"
"gorm.io/gorm"
)

// IdentityHandler handles identity endpoints
type IdentityHandler struct {
db *gorm.DB
}

// NewIdentityHandler creates a new identity handler
func NewIdentityHandler(db *gorm.DB) *IdentityHandler {
return &IdentityHandler{db: db}
}

// RegisterRoutes registers identity routes
func (h *IdentityHandler) RegisterRoutes(g *echo.Group) {
g.GET("/identity", h.GetIdentity)
g.PUT("/identity", h.UpdateIdentity)
g.GET("/identity/avatar", h.GetAvatar)
g.POST("/identity/avatar", h.UploadAvatar)
}

// GetIdentity returns current identity
func (h *IdentityHandler) GetIdentity(c echo.Context) error {
return c.JSON(http.StatusOK, map[string]interface{}{
"name":        "PocketPaw",
"description": "Your AI assistant",
"avatar_url":  "/api/v1/identity/avatar",
})
}

// UpdateIdentity updates identity settings
func (h *IdentityHandler) UpdateIdentity(c echo.Context) error {
var req struct {
Name        string `json:"name"`
Description string `json:"description"`
}
if err := c.Bind(&req); err != nil {
return c.JSON(http.StatusBadRequest, map[string]string{"detail": "Invalid JSON body"})
}
return c.JSON(http.StatusOK, map[string]bool{"ok": true})
}

// GetAvatar returns avatar image
func (h *IdentityHandler) GetAvatar(c echo.Context) error {
return c.NoContent(http.StatusNotFound)
}

// UploadAvatar uploads new avatar
func (h *IdentityHandler) UploadAvatar(c echo.Context) error {
return c.JSON(http.StatusNotImplemented, map[string]string{"detail": "Not implemented"})
}
