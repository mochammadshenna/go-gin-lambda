package main

import (
	"os"

	"ai-service/internal/ai"
	"ai-service/internal/config"
	"ai-service/internal/database"
	"ai-service/internal/handlers"
	"ai-service/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title AI Service API
// @version 1.0
// @description A comprehensive AI service that supports multiple providers (OpenAI, Gemini, Anthropic) for content generation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Initialize configuration
	cfg := config.Load()

	// Setup logging
	setupLogging(cfg.LogLevel)

	log.Info().
		Str("service", cfg.ServiceName).
		Str("version", cfg.ServiceVersion).
		Str("port", cfg.Port).
		Msg("Starting AI Service")

	// Initialize database
	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	log.Info().Str("path", cfg.DBPath).Msg("Database initialized")

	// Initialize AI manager
	aiManager := ai.NewManager(cfg)
	log.Info().Msg("AI Manager initialized")

	// Initialize handlers
	handler := handlers.New(aiManager, db, cfg)

	// Setup Gin router
	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.New()

	// Add middleware
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	
	// Add rate limiting
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRequestsPerMin)
	router.Use(rateLimiter.Middleware())

	// API routes
	v1 := router.Group("/api/v1")
	{
		// AI Generation endpoints
		v1.POST("/generate", handler.Generate)
		v1.POST("/compare", handler.Compare)
		
		// Provider information
		v1.GET("/providers", handler.GetProviders)
		v1.GET("/commands", handler.GetCommands)
		
		// History and statistics
		v1.GET("/history", handler.GetHistory)
		v1.GET("/stats", handler.GetStats)
		
		// Health check
		v1.GET("/health", handler.Health)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": cfg.ServiceName,
			"version": cfg.ServiceVersion,
			"status":  "running",
			"docs":    "/swagger/index.html",
			"endpoints": map[string]string{
				"generate":  "/api/v1/generate",
				"compare":   "/api/v1/compare",
				"providers": "/api/v1/providers",
				"commands":  "/api/v1/commands",
				"history":   "/api/v1/history",
				"stats":     "/api/v1/stats",
				"health":    "/api/v1/health",
			},
		})
	})

	// Start server
	log.Info().
		Str("port", cfg.Port).
		Str("swagger", "http://localhost:"+cfg.Port+"/swagger/index.html").
		Msg("Server starting")

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func setupLogging(level string) {
	// Set global log level
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Configure console output
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	})
}