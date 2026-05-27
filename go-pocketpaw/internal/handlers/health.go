package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pocketpaw/pocketpaw/internal/services/health"
	"gorm.io/gorm"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	healthService *health.HealthService
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		healthService: health.NewHealthService(db),
	}
}

// RegisterRoutes registers health routes
func (h *HealthHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/version", h.GetVersion)
	g.GET("/health", h.GetHealthStatus)
	g.GET("/health/errors", h.GetHealthErrors)
	g.DELETE("/health/errors", h.ClearHealthErrors)
	g.POST("/health/check", h.TriggerHealthCheck)
	g.GET("/audit", h.GetAuditLog)
	g.DELETE("/audit", h.ClearAuditLog)
	g.POST("/security-audit", h.RunSecurityAudit)
	g.GET("/self-audit/reports", h.GetSelfAuditReports)
	g.GET("/self-audit/reports/:date", h.GetSelfAuditReport)
	g.POST("/self-audit/run", h.RunSelfAudit)
}

// GetVersion returns version info
func (h *HealthHandler) GetVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"version":       "0.1.0",
		"go_version":    "go1.19",
		"agent_backend": "claude_agent_sdk",
	})
}

// GetHealthStatus returns health summary
func (h *HealthHandler) GetHealthStatus(c echo.Context) error {
	summary := h.healthService.GetSummary()
	return c.JSON(http.StatusOK, summary)
}

// GetHealthErrors returns recent health errors
func (h *HealthHandler) GetHealthErrors(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}
	search := c.QueryParam("search")
	
	errors := h.healthService.GetRecentErrors(limit, search)
	return c.JSON(http.StatusOK, errors)
}

// ClearHealthErrors clears health errors
func (h *HealthHandler) ClearHealthErrors(c echo.Context) error {
	if err := h.healthService.ClearErrors(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"cleared": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]bool{"cleared": true})
}

// TriggerHealthCheck runs all health checks
func (h *HealthHandler) TriggerHealthCheck(c echo.Context) error {
	summary, err := h.healthService.RunChecks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": "unknown",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, summary)
}

// GetAuditLog returns audit log entries
func (h *HealthHandler) GetAuditLog(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}
	
	logs, err := h.healthService.GetAuditLogs(limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, logs)
}

// ClearAuditLog clears the audit log
func (h *HealthHandler) ClearAuditLog(c echo.Context) error {
	if err := h.healthService.ClearAuditLog(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]bool{"ok": true})
}

// RunSecurityAudit runs security audit checks
func (h *HealthHandler) RunSecurityAudit(c echo.Context) error {
	response := h.healthService.RunSecurityAudit()
	return c.JSON(http.StatusOK, response)
}

// GetSelfAuditReports returns recent self-audit reports
func (h *HealthHandler) GetSelfAuditReports(c echo.Context) error {
	reports := h.healthService.GetSelfAuditReports()
	return c.JSON(http.StatusOK, reports)
}

// GetSelfAuditReport returns a specific self-audit report
func (h *HealthHandler) GetSelfAuditReport(c echo.Context) error {
	date := c.Param("date")
	if date == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "date parameter required"})
	}
	
	report, err := h.healthService.GetSelfAuditReport(date)
	if err != nil {
		if err == health.ErrReportNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"detail": "Report not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, report)
}

// RunSelfAudit runs a self-audit
func (h *HealthHandler) RunSelfAudit(c echo.Context) error {
	report, err := h.healthService.RunSelfAudit()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, report)
}
