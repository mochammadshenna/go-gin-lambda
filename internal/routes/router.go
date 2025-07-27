package routes

import (
	"ai-service/internal/app/middleware"
	"ai-service/internal/controller"
	"ai-service/internal/outbound"

	"github.com/gin-gonic/gin"
)

func NewRouters(aiManager *outbound.Manager) *gin.Engine {
	// Initialize controllers with AI manager
	aiController := controller.NewAIController(aiManager)
	webController := controller.NewWebController(nil)
	healthController := controller.NewHealthController()

	router := router(
		aiController,
		webController,
		healthController,
	)

	return router
}

func router(
	aiController controller.AIController,
	webController controller.WebController,
	healthController controller.HealthController,
) *gin.Engine {
	// set gin mode
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()

	router := gin.New()

	// global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORSMiddleware())

	router.HandleMethodNotAllowed = true

	// API routes
	api := router.Group("/api")
	{
		// AI Generation endpoints
		api.POST("/generate", aiController.GenerateContent)
		api.POST("/compare", aiController.CompareProviders)
		api.GET("/providers", aiController.GetProviders)
		api.GET("/history", aiController.GetHistory)
		api.GET("/stats", aiController.GetStats)

		// Health check
		api.GET("/health", healthController.GetHealthCheck)
	}

	// Web UI routes
	web := router.Group("/")
	{
		web.GET("/", webController.Home)
		web.GET("/history", webController.History)
		web.GET("/stats", webController.Stats)
		web.GET("/error", webController.Error)
	}

	return router
}
