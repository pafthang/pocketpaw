package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex" json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	Active    bool           `json:"active"`
}

// Session represents a user session
type Session struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	UserID         uint           `json:"user_id"`
	Token          string         `gorm:"uniqueIndex" json:"-"`
	ExpiresAt      time.Time      `json:"expires_at"`
	LastActivity   time.Time      `json:"last_activity"`
	Channel        string         `json:"channel"`
	Title          string         `json:"title"`
	Metadata       string         `gorm:"type:json" json:"metadata,omitempty"`
}

// APIKey represents an API key for authentication
type APIKey struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name"`
	Key         string         `gorm:"uniqueIndex" json:"-"`
	KeyPrefix   string         `gorm:"index" json:"key_prefix"`
	Scopes      string         `gorm:"type:json" json:"scopes"`
	LastUsedAt  *time.Time     `json:"last_used_at,omitempty"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	Active      bool           `json:"active"`
}

// OAuthToken represents an OAuth2 token
type OAuthToken struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	ClientID     string         `gorm:"index" json:"client_id"`
	AccessToken  string         `gorm:"uniqueIndex" json:"-"`
	RefreshToken string         `json:"-"`
	Scope        string         `json:"scope"`
	ExpiresAt    time.Time      `json:"expires_at"`
	UserID       uint           `json:"user_id"`
}

// OAuthClient represents an OAuth2 client
type OAuthClient struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	ClientID      string         `gorm:"uniqueIndex" json:"client_id"`
	ClientSecret  string         `json:"-"`
	Name          string         `json:"name"`
	RedirectURIs  string         `gorm:"type:json" json:"redirect_uris"`
	Scopes        string         `gorm:"type:json" json:"scopes"`
	Active        bool           `json:"active"`
}

// Channel represents a communication channel
type Channel struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Name       string         `gorm:"uniqueIndex" json:"name"`
	Type       string         `json:"type"` // telegram, whatsapp, slack, discord, etc.
	Config     string         `gorm:"type:json" json:"config,omitempty"`
	Active     bool           `json:"active"`
	LastSyncAt *time.Time     `json:"last_sync_at,omitempty"`
}

// MemoryEntry represents a memory/message entry
type MemoryEntry struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	SessionID string         `gorm:"index" json:"session_id"`
	Role      string         `json:"role"` // user, assistant, system
	Content   string         `gorm:"type:text" json:"content"`
	Metadata  string         `gorm:"type:json" json:"metadata,omitempty"`
}

// Settings represents application settings
type Setting struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Key       string         `gorm:"uniqueIndex" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
}

// Skill represents a skill/plugin
type Skill struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Name       string        `gorm:"uniqueIndex" json:"name"`
	Description string       `json:"description"`
	Enabled    bool          `json:"enabled"`
	Config     string        `gorm:"type:json" json:"config,omitempty"`
}

// Webhook represents a webhook endpoint
type Webhook struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name"`
	URL       string         `json:"url"`
	Secret    string         `json:"-"`
	Events    string         `gorm:"type:json" json:"events"`
	Active    bool           `json:"active"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Action    string    `json:"action"`
	UserID    uint      `json:"user_id,omitempty"`
	Resource  string    `json:"resource"`
	Details   string    `gorm:"type:json" json:"details,omitempty"`
	IPAddress string    `json:"ip_address,omitempty"`
}

// TableName overrides the table name used by Session to avoid conflicts
func (Session) TableName() string {
	return "sessions"
}
