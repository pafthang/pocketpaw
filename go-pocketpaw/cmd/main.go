package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/pocketpaw/pocketpaw/internal/config"
	"github.com/pocketpaw/pocketpaw/internal/database"
	"github.com/pocketpaw/pocketpaw/internal/handlers"
	"github.com/pocketpaw/pocketpaw/internal/middleware"
	"github.com/pocketpaw/pocketpaw/internal/models"
	"github.com/pocketpaw/pocketpaw/internal/services/chat"
	"github.com/pocketpaw/pocketpaw/internal/services/sessions"
)

func main() {
	// Load configuration
	if _, err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Get config dir for database path
	configDir, err := config.GetConfigDir()
	if err != nil {
		log.Fatalf("Failed to get config dir: %v", err)
	}
	dbPath := filepath.Join(configDir, "pocketpaw.db")

	// Initialize database
	if err := database.Init(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	db := database.GetDB()

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.APIKey{},
		&models.OAuthToken{},
		&models.OAuthClient{},
		&models.Channel{},
		&models.MemoryEntry{},
		&models.Setting{},
		&models.Skill{},
		&models.Webhook{},
		&models.AuditLog{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	// Initialize services
	sessionService := sessions.NewSessionService(db)
	chatService := chat.NewChatService(db)

	// Initialize handlers
	sessionsHandler := handlers.NewSessionsHandler(sessionService)
	chatHandler := handlers.NewChatHandler(chatService)
	authHandler := handlers.NewAuthHandler()
	healthHandler := handlers.NewHealthHandler(db)

	// Auth middleware
	authMW := middleware.NewAuthMiddleware()

	// API v1 group
	v1 := e.Group("/api/v1")

	// Public auth routes
	authGroup := v1.Group("")
	authHandler.RegisterRoutes(authGroup)

	// Public health routes
	healthGroup := v1.Group("")
	healthHandler.RegisterRoutes(healthGroup)

	// Protected routes
	protected := v1.Group("")
	protected.Use(authMW.Middleware)

	// Sessions routes
	sessionsGroup := protected.Group("/sessions")
	sessionsHandler.RegisterRoutes(sessionsGroup)

	// Chat routes
	chatGroup := protected.Group("/chat")
	chatHandler.RegisterRoutes(chatGroup)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Starting PocketPaw server on %s", addr)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
