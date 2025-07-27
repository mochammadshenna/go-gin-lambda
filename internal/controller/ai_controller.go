package controller

import (
	"ai-service/internal/model"
	"ai-service/internal/outbound"
	"ai-service/internal/repository"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type AIController interface {
	GenerateContent(c *gin.Context)
	CompareProviders(c *gin.Context)
	GetProviders(c *gin.Context)
	GetHistory(c *gin.Context)
	GetStats(c *gin.Context)
}

type aiController struct {
	aiManager      *outbound.Manager
	generationRepo repository.GenerationRepository
}

func NewAIController(aiManager *outbound.Manager, generationRepo repository.GenerationRepository) AIController {
	return &aiController{
		aiManager:      aiManager,
		generationRepo: generationRepo,
	}
}

func (c *aiController) GenerateContent(ctx *gin.Context) {
	var request struct {
		Provider    string  `json:"provider" binding:"required"`
		Model       string  `json:"model" binding:"required"`
		Prompt      string  `json:"prompt" binding:"required"`
		SystemMsg   string  `json:"systemMsg"`
		Temperature float64 `json:"temperature"`
		MaxTokens   int     `json:"maxTokens"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Convert provider string to AIProvider type
	var provider model.AIProvider
	switch request.Provider {
	case "openai":
		provider = model.OpenAI
	case "gemini":
		provider = model.Gemini
	case "anthropic":
		provider = model.Anthropic
	default:
		ctx.JSON(400, gin.H{"error": "Unsupported provider", "provider": request.Provider})
		return
	}

	// Create generation request
	genReq := &model.GenerationRequest{
		Provider:    provider,
		Model:       request.Model,
		Prompt:      request.Prompt,
		SystemMsg:   request.SystemMsg,
		Temperature: float32(request.Temperature),
		MaxTokens:   request.MaxTokens,
	}

	// Generate content using AI manager
	startTime := time.Now()
	response, err := c.aiManager.Generate(ctx, genReq)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error":   "Failed to generate content",
			"details": err.Error(),
		})
		return
	}
	duration := time.Since(startTime)

	// Save generation record to database
	generationRecord := &model.GenerationHistory{
		Provider:     string(response.Provider),
		Model:        response.Model,
		Prompt:       genReq.Prompt,
		Response:     response.Content,
		TokensUsed:   response.TokensUsed,
		Duration:     int64(duration.Milliseconds()),
		Status:       "success",
		ErrorMessage: "",
	}

	err = c.generationRepo.Create(ctx, generationRecord)
	if err != nil {
		// Log the error but don't fail the request
		log.Printf("Failed to save generation record: %v", err)
	}

	// Return the response
	ctx.JSON(200, gin.H{
		"content":     response.Content,
		"provider":    string(response.Provider),
		"model":       response.Model,
		"tokens_used": response.TokensUsed,
		"duration":    duration.String(),
		"status":      "success",
	})
}

func (c *aiController) CompareProviders(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "Compare providers endpoint"})
}

func (c *aiController) GetProviders(ctx *gin.Context) {
	providers := []gin.H{
		{
			"id":          "openai",
			"name":        "OpenAI",
			"description": "Advanced language models for text generation",
			"models":      []string{"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo"},
			"available":   true,
		},
		{
			"id":          "gemini",
			"name":        "Google Gemini",
			"description": "Google's multimodal AI model",
			"models":      []string{"gemini-1.5-flash", "gemini-1.5-pro", "gemini-2.0-flash"},
			"available":   true,
		},
		{
			"id":          "anthropic",
			"name":        "Anthropic Claude",
			"description": "Constitutional AI for safe and helpful responses",
			"models":      []string{"claude-3-sonnet", "claude-3-opus", "claude-3-haiku"},
			"available":   false,
		},
	}

	ctx.JSON(200, gin.H{
		"providers": providers,
		"total":     len(providers),
	})
}

func (c *aiController) GetHistory(ctx *gin.Context) {
	history := []gin.H{
		{
			"id":          "gen-123",
			"provider":    "openai",
			"model":       "gpt-3.5-turbo",
			"prompt":      "Write a hello world program in Go",
			"response":    "Here's a simple hello world program in Go...",
			"tokens_used": 150,
			"duration":    "2.5s",
			"created_at":  "2024-01-15T10:30:00Z",
			"status":      "success",
		},
		{
			"id":          "gen-124",
			"provider":    "gemini",
			"model":       "gemini-1.5-flash",
			"prompt":      "Explain quantum computing",
			"response":    "Quantum computing is a revolutionary technology...",
			"tokens_used": 300,
			"duration":    "1.8s",
			"created_at":  "2024-01-15T09:15:00Z",
			"status":      "success",
		},
	}

	ctx.JSON(200, gin.H{
		"history": history,
		"total":   len(history),
	})
}

func (c *aiController) GetStats(ctx *gin.Context) {
	stats := map[string]gin.H{
		"openai": {
			"provider":          "openai",
			"total_generations": 150,
			"total_tokens":      45000,
			"avg_duration":      2.5,
			"error_count":       3,
			"success_rate":      98.0,
		},
		"gemini": {
			"provider":          "gemini",
			"total_generations": 75,
			"total_tokens":      25000,
			"avg_duration":      1.8,
			"error_count":       1,
			"success_rate":      98.7,
		},
		"anthropic": {
			"provider":          "anthropic",
			"total_generations": 25,
			"total_tokens":      8000,
			"avg_duration":      3.2,
			"error_count":       0,
			"success_rate":      100.0,
		},
	}

	ctx.JSON(200, gin.H{
		"stats": stats,
		"summary": gin.H{
			"total_generations": 250,
			"total_tokens":      78000,
			"avg_duration":      2.3,
			"total_errors":      4,
		},
	})
}
