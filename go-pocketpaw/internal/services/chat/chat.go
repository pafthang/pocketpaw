package chat

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatService handles chat operations
type ChatService struct {
	db *gorm.DB
}

// NewChatService creates a new chat service instance
func NewChatService(db *gorm.DB) *ChatService {
	return &ChatService{db: db}
}

// ChatRequest represents a chat request
type ChatRequest struct {
	Content     string       `json:"content"`
	SessionID   string       `json:"session_id,omitempty"`
	Media       []string     `json:"media,omitempty"`
	FileContext *FileContext `json:"file_context,omitempty"`
}

// FileContext represents file context in a chat request
type FileContext struct {
	Filename    string `json:"filename,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Size        int64  `json:"size,omitempty"`
	Data        string `json:"data,omitempty"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	SessionID string            `json:"session_id"`
	Content   string            `json:"content"`
	Usage     map[string]interface{} `json:"usage,omitempty"`
}

// StreamStartEvent represents the start of a stream
type StreamStartEvent struct {
	SessionID string `json:"session_id"`
}

// ChunkEvent represents a chunk of content
type ChunkEvent struct {
	Content string `json:"content"`
	Type    string `json:"type"`
}

// ToolStartEvent represents the start of a tool execution
type ToolStartEvent struct {
	Tool  string                 `json:"tool"`
	Input map[string]interface{} `json:"input"`
}

// ToolResultEvent represents the result of a tool execution
type ToolResultEvent struct {
	Tool   string      `json:"tool"`
	Output interface{} `json:"output"`
}

// ThinkingEvent represents thinking content
type ThinkingEvent struct {
	Content string `json:"content"`
}

// AskUserQuestionEvent represents a question for the user
type AskUserQuestionEvent struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

// PocketCreatedEvent represents a pocket creation event
type PocketCreatedEvent struct {
	Spec          map[string]interface{} `json:"spec"`
	SessionID     string                 `json:"session_id"`
	PocketCloudID string                 `json:"pocket_cloud_id,omitempty"`
}

// PocketMutationEvent represents a pocket mutation event
type PocketMutationEvent struct {
	Mutation map[string]interface{} `json:"mutation"`
}

// ErrorEvent represents an error event
type ErrorEvent struct {
	Detail string `json:"detail"`
}

// StreamEndEvent represents the end of a stream
type StreamEndEvent struct {
	SessionID string            `json:"session_id"`
	Usage     map[string]interface{} `json:"usage,omitempty"`
}

// StopResponse represents the response for stopping a chat
type StopResponse struct {
	Status    string `json:"status"`
	SessionID string `json:"session_id"`
}

// extractChatID converts a client-supplied session_id to a raw chat_id
func extractChatID(sessionID string) string {
	const wsPrefix = "websocket_"
	if sessionID == "" {
		return uuid.New().String()[:12]
	}
	if len(sessionID) > len(wsPrefix) && sessionID[:len(wsPrefix)] == wsPrefix {
		return sessionID[len(wsPrefix):]
	}
	return sessionID
}

// toSafeKey builds the safe_key that the client stores as its session identifier
func toSafeKey(chatID string) string {
	return "websocket_" + chatID
}

// SendMessage publishes a message and returns the chat_id
func (s *ChatService) SendMessage(req *ChatRequest) (string, error) {
	chatID := extractChatID(req.SessionID)
	// In production, this would publish to the message bus
	// For now, just return the chat_id
	return chatID, nil
}

// GetChatResponse collects the complete response for a chat message
func (s *ChatService) GetChatResponse(req *ChatRequest) (*ChatResponse, error) {
	chatID := extractChatID(req.SessionID)
	
	// In production, this would:
	// 1. Start a bridge to listen to the message bus
	// 2. Send the message via SendMessage
	// 3. Collect chunks until stream_end
	// 4. Return the aggregated response
	
	// For now, return a placeholder response
	return &ChatResponse{
		SessionID: toSafeKey(chatID),
		Content:   "",
		Usage:     map[string]interface{}{},
	}, nil
}

// StreamChat starts an SSE stream for chat responses
// This would typically return an io.Reader for SSE events
func (s *ChatService) StreamChat(req *ChatRequest) (string, error) {
	chatID := extractChatID(req.SessionID)
	safeKey := toSafeKey(chatID)
	
	// In production, this would:
	// 1. Cancel any existing stream for this session
	// 2. Start a new bridge to listen to the message bus
	// 3. Send the message
	// 4. Return an SSE stream generator
	
	return safeKey, nil
}

// StopChat cancels an in-flight chat response
func (s *ChatService) StopChat(sessionID string) (*StopResponse, error) {
	if sessionID == "" {
		return nil, &ChatError{Message: "session_id is required"}
	}
	
	// In production, this would:
	// 1. Find and cancel the active stream
	// 2. Cancel the agent loop task
	
	// Check if session exists (placeholder logic)
	cancelEvent := activeStreams[sessionID]
	if cancelEvent == nil && len(sessionID) > 10 && sessionID[:10] != "websocket_" {
		cancelEvent = activeStreams["websocket_"+sessionID]
	}
	if cancelEvent == nil {
		return nil, &ChatError{Message: "No active stream for this session"}
	}
	
	// Signal cancellation
	// cancelEvent.Set() // Would be called on actual event
	
	return &StopResponse{
		Status:    "ok",
		SessionID: sessionID,
	}, nil
}

// ChatError represents a chat service error
type ChatError struct {
	Message string
}

func (e *ChatError) Error() string {
	return e.Message
}

// Error definitions
var (
	ErrNoActiveStream = &ChatError{Message: "No active stream for this session"}
	ErrSessionRequired = &ChatError{Message: "session_id is required"}
)

// activeStreams tracks active SSE streams (in production, use sync.Map with proper types)
var activeStreams = make(map[string]interface{})
