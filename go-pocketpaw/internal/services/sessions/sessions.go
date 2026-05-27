package sessions

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/pocketpaw/pocketpaw/internal/models"
)

// SessionService handles session operations
type SessionService struct {
	db *gorm.DB
}

// NewSessionService creates a new session service instance
func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

// Session represents a chat session
type Session struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Channel      string     `json:"channel"`
	LastActivity *time.Time `json:"last_activity,omitempty"`
}

// SessionCreateResponse represents the response for session creation
type SessionCreateResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// SessionListResponse represents the response for session listing
type SessionListResponse struct {
	Sessions []Session `json:"sessions"`
	Total    int       `json:"total"`
}

// SessionSearchResult represents a search result item
type SessionSearchResult struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Channel      string `json:"channel"`
	Match        string `json:"match"`
	MatchRole    string `json:"match_role"`
	LastActivity string `json:"last_activity"`
}

// SessionSearchResponse represents the response for session search
type SessionSearchResponse struct {
	Sessions []SessionSearchResult `json:"sessions"`
}

// SessionTitleRequest represents a request to update session title
type SessionTitleRequest struct {
	Title string `json:"title"`
}

// CreateSession creates a new empty session
func (s *SessionService) CreateSession() (*SessionCreateResponse, error) {
	safeKey := "websocket_" + uuid.New().String()[:12]
	return &SessionCreateResponse{
		ID:    safeKey,
		Title: "New Chat",
	}, nil
}

// ListSessions lists sessions with optional limit
func (s *SessionService) ListSessions(limit int) (*SessionListResponse, error) {
	if s.db == nil {
		return &SessionListResponse{
			Sessions: []Session{},
			Total:    0,
		}, nil
	}

	var sessions []models.Session
	result := s.db.Order("last_activity DESC").Limit(limit).Find(&sessions)

	response := &SessionListResponse{
		Sessions: make([]Session, 0),
		Total:    int(result.RowsAffected),
	}

	for _, sess := range sessions {
		response.Sessions = append(response.Sessions, Session{
			ID:           uuid.New().String(), // In production, use proper ID mapping
			Title:        sess.Title,
			Channel:      sess.Channel,
			LastActivity: &sess.LastActivity,
		})
	}

	return response, nil
}

// DeleteSession deletes a session by ID
func (s *SessionService) DeleteSession(sessionID string) error {
	if s.db == nil {
		return ErrSessionNotFound
	}

	result := s.db.Where("token = ? OR id = ?", sessionID, sessionID).Delete(&models.Session{})
	if result.RowsAffected == 0 {
		return ErrSessionNotFound
	}

	return nil
}

// UpdateSessionTitle updates the title of a session
func (s *SessionService) UpdateSessionTitle(sessionID string, title string) error {
	if s.db == nil {
		return ErrSessionNotFound
	}

	result := s.db.Model(&models.Session{}).Where("token = ? OR id = ?", sessionID, sessionID).Update("title", title)
	if result.RowsAffected == 0 {
		return ErrSessionNotFound
	}

	return nil
}

// SearchSessions searches sessions by content
func (s *SessionService) SearchSessions(query string, limit int) (*SessionSearchResponse, error) {
	response := &SessionSearchResponse{
		Sessions: []SessionSearchResult{},
	}

	if query == "" {
		return response, nil
	}

	// In production, this would search through memory entries
	// For now, return empty results
	return response, nil
}

// GetSessionHistory gets message history for a session
func (s *SessionService) GetSessionHistory(sessionID string, limit int) ([]map[string]interface{}, error) {
	if s.db == nil {
		return []map[string]interface{}{}, nil
	}

	var entries []models.MemoryEntry
	result := s.db.Where("session_id = ?", sessionID).Order("created_at ASC").Limit(limit).Find(&entries)

	history := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		history = append(history, map[string]interface{}{
			"id":         entry.ID,
			"role":       entry.Role,
			"content":    entry.Content,
			"created_at": entry.CreatedAt,
			"metadata":   entry.Metadata,
		})
	}

	_ = result // Use result to avoid unused variable warning
	return history, nil
}

// ExportSession exports a session in the specified format
func (s *SessionService) ExportSession(sessionID string, format string) ([]byte, string, error) {
	if s.db == nil {
		return nil, "", ErrSessionNotFound
	}

	var entries []models.MemoryEntry
	result := s.db.Where("session_id = ?", sessionID).Order("created_at ASC").Find(&entries)
	if result.RowsAffected == 0 {
		return nil, "", ErrSessionNotFound
	}

	if format == "json" {
		// Return JSON export
		// In production, use proper JSON marshaling
		return []byte(`{"export_version":"1.0","session_id":"` + sessionID + `"}`), "application/json", nil
	} else if format == "md" {
		// Return Markdown export
		return []byte("# Conversation Export\n\nSession: " + sessionID), "text/markdown", nil
	}

	return nil, "", ErrInvalidFormat
}

// Error definitions
var (
	ErrSessionNotFound = &SessionError{Message: "Session not found"}
	ErrInvalidFormat   = &SessionError{Message: "Invalid export format"}
)

// SessionError represents a session service error
type SessionError struct {
	Message string
}

func (e *SessionError) Error() string {
	return e.Message
}
