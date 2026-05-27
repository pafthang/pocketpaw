package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var (
	configDir string
	settings  *Settings
	mu        sync.RWMutex
)

// Settings represents PocketPaw settings
type Settings struct {
	// Telegram
	TelegramBotToken string `json:"telegram_bot_token,omitempty"`
	AllowedUserID    *int   `json:"allowed_user_id,omitempty"`

	// Agent Backend
	AgentBackend    string   `json:"agent_backend"`
	FallbackBackends []string `json:"fallback_backends,omitempty"`

	// Claude Agent SDK Settings
	ClaudeSDKProvider  string `json:"claude_sdk_provider"`
	ClaudeSDKModel     string `json:"claude_sdk_model"`
	ClaudeSDKMaxTurns  int    `json:"claude_sdk_max_turns"`

	// OpenAI Agents SDK Settings
	OpenAIAgentsProvider string `json:"openai_agents_provider"`
	OpenAIAgentsModel    string `json:"openai_agents_model"`
	OpenAIAgentsMaxTurns int    `json:"openai_agents_max_turns"`

	// Google ADK Settings
	GoogleADKProvider string `json:"google_adk_provider"`
	GoogleADKModel    string `json:"google_adk_model"`
	GoogleADKMaxTurns int    `json:"google_adk_max_turns"`

	// Codex CLI Settings
	CodexCLIModel    string `json:"codex_cli_model"`
	CodexCLIMaxTurns int    `json:"codex_cli_max_turns"`

	// Copilot SDK Settings
	CopilotSDKProvider string `json:"copilot_sdk_provider"`
	CopilotSDKModel    string `json:"copilot_sdk_model"`
	CopilotSDKMaxTurns int    `json:"copilot_sdk_max_turns"`

	// DeepAgents Settings
	DeepAgentsModel    string `json:"deep_agents_model"`
	DeepAgentsMaxTurns int    `json:"deep_agents_max_turns"`

	// Opencode Settings
	OpencodeBaseURL    string `json:"opencode_base_url,omitempty"`
	OpencodeModel      string `json:"opencode_model"`
	OpencodeMaxTurns   int    `json:"opencode_max_turns"`

	// LiteLLM Settings
	LiteLLMBaseURL   string `json:"litellm_api_base,omitempty"`
	LiteLLMAPIKey    string `json:"litellm_api_key,omitempty"`
	LiteLLMModel     string `json:"litellm_model"`
	LiteLLMMaxTurns  int    `json:"litellm_max_turns"`

	// Ollama Settings
	OllamaBaseURL   string `json:"ollama_base_url,omitempty"`
	OllamaEmbedding string `json:"ollama_embedding,omitempty"`

	// OpenAI Compatible Settings
	OpenAICompatibleBaseURL string `json:"openai_compatible_base_url,omitempty"`
	OpenAICompatibleAPIKey  string `json:"openai_compatible_api_key,omitempty"`
	OpenAICompatibleModel   string `json:"openai_compatible_model"`

	// Anthropic
	AnthropicAPIKey string `json:"anthropic_api_key,omitempty"`

	// OpenAI
	OpenAIAPIKey       string `json:"openai_api_key,omitempty"`
	OpenAIOrgID        string `json:"openai_org_id,omitempty"`
	OpenAIProjectID    string `json:"openai_project_id,omitempty"`

	// OpenRouter
	OpenRouterAPIKey string `json:"openrouter_api_key,omitempty"`

	// Groq
	GroqAPIKey string `json:"groq_api_key,omitempty"`

	// Gemini
	GeminiAPIKey string `json:"gemini_api_key,omitempty"`

	// Cohere
	CohereAPIKey string `json:"cohere_api_key,omitempty"`

	// HuggingFace
	HuggingFaceAPIKey string `json:"huggingface_api_key,omitempty"`

	// Together
	TogetherAPIKey string `json:"together_api_key,omitempty"`

	// Anyscale
	AnyscaleAPIKey string `json:"anyscale_api_key,omitempty"`

	// Perplexity
	PerplexityAPIKey string `json:"perplexity_api_key,omitempty"`

	// Mistral
	MistralAPIKey string `json:"mistral_api_key,omitempty"`

	// DeepSeek
	DeepSeekAPIKey string `json:"deepseek_api_key,omitempty"`

	// Embedding
	EmbeddingProvider   string `json:"embedding_provider"`
	EmbeddingBaseURL    string `json:"embedding_base_url,omitempty"`
	EmbeddingAPIKey     string `json:"embedding_api_key,omitempty"`
	EmbeddingModel      string `json:"embedding_model"`
	EmbeddingDimensions int    `json:"embedding_dimensions"`

	// Mem0 Settings
	Mem0OllamaBaseURL string `json:"mem0_ollama_base_url,omitempty"`
	Mem0User          string `json:"mem0_user"`

	// Signal Settings
	SignalAPIURL     string `json:"signal_api_url,omitempty"`
	SignalPhone      string `json:"signal_phone,omitempty"`
	SignalRecipients []string `json:"signal_recipients,omitempty"`

	// WhatsApp Settings
	WhatsAppMode         string `json:"whatsapp_mode"`
	WhatsAppCloudAPIKey  string `json:"whatsapp_cloud_api_key,omitempty"`
	WhatsAppPhoneNumber  string `json:"whatsapp_phone_number,omitempty"`
	WhatsAppBusinessID   string `json:"whatsapp_business_id,omitempty"`

	// TTS/STT
	TTSProvider  string `json:"tts_provider"`
	STTProvider  string `json:"stt_provider"`
	ElevenLabsAPIKey string `json:"elevenlabs_api_key,omitempty"`

	// Soul Settings
	SoulCognitiveModel string `json:"soul_cognitive_model"`
	SoulChatModel      string `json:"soul_chat_model"`

	// Health
	HealthCheckOnStartup bool `json:"health_check_on_startup"`

	// Session Token
	SessionTokenTTLHours int `json:"session_token_ttl_hours"`

	// CORS
	APICORSAllowedOrigins []string `json:"api_cors_allowed_origins,omitempty"`

	// MCP Client
	MCPClientMetadataURL string `json:"mcp_client_metadata_url,omitempty"`

	// Guardian
	GuardianAgentEnabled bool   `json:"guardian_agent_enabled"`
	GuardianAgentModel   string `json:"guardian_agent_model"`

	// File Jail
	FileJailEnabled bool     `json:"file_jail_enabled"`
	FileJailPaths   []string `json:"file_jail_paths,omitempty"`

	// Rate Limiting
	RateLimitEnabled       bool `json:"rate_limit_enabled"`
	RateLimitRequestsPerMin int `json:"rate_limit_requests_per_min"`

	// Audit
	AuditLogEnabled bool `json:"audit_log_enabled"`

	// Remote Access
	RemoteAccessEnabled bool   `json:"remote_access_enabled"`
	RemoteAccessPassword string `json:"remote_access_password,omitempty"`

	// Cloudflare Tunnel
	CloudflareTunnelEnabled bool   `json:"cloudflare_tunnel_enabled"`
	CloudflareTunnelSubdomain string `json:"cloudflare_tunnel_subdomain,omitempty"`

	// Bot Gateway
	BotGatewayHost string `json:"bot_gateway_host,omitempty"`
	BotGatewayPort int    `json:"bot_gateway_port"`

	// Analytics
	AnalyticsEnabled bool `json:"analytics_enabled"`

	// Budget
	BudgetDailyLimitUSD float64 `json:"budget_daily_limit_usd"`

	// Vector DB
	VectorDBProvider string `json:"vector_db_provider"`
	VectorDBPath     string `json:"vector_db_path,omitempty"`

	// Browser
	BrowserHeadless bool `json:"browser_headless"`

	// Skills
	SkillsEnabled bool `json:"skills_enabled"`

	// Knowledge
	KnowledgeEnabled bool   `json:"knowledge_enabled"`
	KnowledgePath    string `json:"knowledge_path,omitempty"`

	// A2A
	A2AEnabled bool `json:"a2a_enabled"`

	// Daemon
	DaemonEnabled bool `json:"daemon_enabled"`

	// Bootstrap
	BootstrapEnabled bool `json:"bootstrap_enabled"`

	// Fleet
	FleetEnabled bool `json:"fleet_enabled"`

	// Instinct
	InstinctEnabled bool `json:"instinct_enabled"`

	// Retrieval
	RetrievalEnabled bool `json:"retrieval_enabled"`

	// Widget
	WidgetEnabled bool `json:"widget_enabled"`

	// Automations
	AutomationsEnabled bool `json:"automations_enabled"`
}

// GetConfigDir returns the config directory path
func GetConfigDir() (string, error) {
	if configDir != "" {
		return configDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir = filepath.Join(homeDir, ".pocketpaw")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}

	return configDir, nil
}

// GetConfigPath returns the config file path
func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// GetTokenPath returns the access token file path
func GetTokenPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "access_token"), nil
}

// Load loads settings from config file
func Load() (*Settings, error) {
	mu.Lock()
	defer mu.Unlock()

	if settings != nil {
		return settings, nil
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default settings
			settings = getDefaultSettings()
			return settings, nil
		}
		return nil, err
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	// Apply defaults for zero values
	applyDefaults(&s)
	settings = &s
	return settings, nil
}

// Save saves settings to config file
func (s *Settings) Save() error {
	mu.Lock()
	defer mu.Unlock()

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return err
	}

	settings = s
	return nil
}

// GetAccessToken reads the access token from file
func GetAccessToken() (string, error) {
	tokenPath, err := GetTokenPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	return string(data), nil
}

// RegenerateToken generates a new access token
func RegenerateToken() (string, error) {
	token := generateSecureToken()

	tokenPath, err := GetTokenPath()
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		return "", err
	}

	// Invalidate old sessions - would need session service integration
	mu.Lock()
	settings = nil // Force reload
	mu.Unlock()

	return token, nil
}

func getDefaultSettings() *Settings {
	return &Settings{
		AgentBackend:           "claude_agent_sdk",
		ClaudeSDKProvider:      "anthropic",
		ClaudeSDKMaxTurns:      100,
		OpenAIAgentsProvider:   "openai",
		OpenAIAgentsMaxTurns:   100,
		GoogleADKProvider:      "google",
		GoogleADKModel:         "gemini-3-pro-preview",
		GoogleADKMaxTurns:      100,
		CodexCLIModel:          "gpt-5.3-codex",
		CodexCLIMaxTurns:       100,
		CopilotSDKProvider:     "copilot",
		CopilotSDKMaxTurns:     100,
		DeepAgentsMaxTurns:     100,
		LiteLLMMaxTurns:        100,
		EmbeddingProvider:      "ollama",
		EmbeddingDimensions:    768,
		Mem0User:               "default",
		WhatsAppMode:           "cloud",
		TTSProvider:            "elevenlabs",
		STTProvider:            "whisper",
		SoulCognitiveModel:     "haiku",
		SoulChatModel:          "sonnet",
		HealthCheckOnStartup:   true,
		SessionTokenTTLHours:   24,
		GuardianAgentEnabled:   true,
		GuardianAgentModel:     "haiku",
		FileJailEnabled:        true,
		RateLimitEnabled:       true,
		RateLimitRequestsPerMin: 60,
		AuditLogEnabled:        true,
		BotGatewayPort:         8889,
		AnalyticsEnabled:       true,
		BudgetDailyLimitUSD:    10.0,
		VectorDBProvider:       "chroma",
		BrowserHeadless:        true,
		SkillsEnabled:          true,
		KnowledgeEnabled:       true,
		A2AEnabled:             true,
		DaemonEnabled:          true,
		BootstrapEnabled:       true,
		FleetEnabled:           true,
		InstinctEnabled:        true,
		RetrievalEnabled:       true,
		WidgetEnabled:          true,
		AutomationsEnabled:     true,
	}
}

func generateSecureToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to simple token if random fails
		return "pocketpaw_token_" + hex.EncodeToString([]byte("fallback"))
	}
	return "pp_" + hex.EncodeToString(bytes)
}

func applyDefaults(s *Settings) {
	defaults := getDefaultSettings()
	
	if s.AgentBackend == "" {
		s.AgentBackend = defaults.AgentBackend
	}
	if s.ClaudeSDKProvider == "" {
		s.ClaudeSDKProvider = defaults.ClaudeSDKProvider
	}
	if s.ClaudeSDKMaxTurns == 0 {
		s.ClaudeSDKMaxTurns = defaults.ClaudeSDKMaxTurns
	}
	if s.SessionTokenTTLHours == 0 {
		s.SessionTokenTTLHours = defaults.SessionTokenTTLHours
	}
	// Add more default applications as needed
}
