package controller

import (
	"github.com/gin-gonic/gin"
)

type HealthController interface {
	GetHealthCheck(c *gin.Context)
}

type healthController struct {
	// service will be injected
}

func NewHealthController() HealthController {
	return &healthController{}
}

func (c *healthController) GetHealthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status":    "healthy",
		"service":   "ai-service",
		"version":   "1.0.0",
		"timestamp": "2024-01-15T10:30:00Z",
	})
}
