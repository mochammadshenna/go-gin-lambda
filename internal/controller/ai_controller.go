package controller

import (
	"ai-service/internal/model"
	"ai-service/internal/outbound"
	"ai-service/internal/repository"
	"fmt"
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
	generations, err := c.generationRepo.GetRecent(ctx, 50, 0)
	if err != nil {
		log.Printf("Failed to load generation history: %v", err)
		ctx.JSON(500, gin.H{
			"error":   "Failed to load generation history",
			"details": err.Error(),
		})
		return
	}

	// Convert to API response format
	history := make([]gin.H, len(generations))
	for i, gen := range generations {
		history[i] = gin.H{
			"id":          gen.ID,
			"provider":    gen.Provider,
			"model":       gen.Model,
			"prompt":      gen.Prompt,
			"response":    gen.Response,
			"tokens_used": gen.TokensUsed,
			"duration":    fmt.Sprintf("%dms", gen.Duration),
			"created_at":  gen.CreatedAt.Format(time.RFC3339),
			"status":      gen.Status,
		}
	}

	ctx.JSON(200, gin.H{
		"history": history,
		"total":   len(history),
	})
}

func (c *aiController) GetStats(ctx *gin.Context) {
	// Get stats for the last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	providerStats, err := c.generationRepo.GetStats(ctx, startDate, endDate)
	if err != nil {
		log.Printf("Failed to load usage statistics: %v", err)
		ctx.JSON(500, gin.H{
			"error":   "Failed to load usage statistics",
			"details": err.Error(),
		})
		return
	}

	// Convert to API response format
	stats := make(map[string]gin.H)
	var totalGenerations, totalTokens, totalErrors int
	var totalDuration int64

	for _, stat := range providerStats {
		totalGenerations += stat.TotalGenerations
		totalTokens += stat.TotalTokens
		totalErrors += stat.ErrorCount
		totalDuration += int64(stat.AvgDuration * float64(stat.TotalGenerations))

		successRate := 100.0
		if stat.TotalGenerations > 0 {
			successRate = float64(stat.TotalGenerations-stat.ErrorCount) / float64(stat.TotalGenerations) * 100
		}

		stats[stat.Provider] = gin.H{
			"provider":          stat.Provider,
			"total_generations": stat.TotalGenerations,
			"total_tokens":      stat.TotalTokens,
			"avg_duration":      stat.AvgDuration,
			"error_count":       stat.ErrorCount,
			"success_rate":      successRate,
		}
	}

	avgDuration := float64(0)
	if totalGenerations > 0 {
		avgDuration = float64(totalDuration) / float64(totalGenerations)
	}

	successRate := 100.0
	if totalGenerations > 0 {
		successRate = float64(totalGenerations-totalErrors) / float64(totalGenerations) * 100
	}

	ctx.JSON(200, gin.H{
		"stats": stats,
		"summary": gin.H{
			"total_generations": totalGenerations,
			"total_tokens":      totalTokens,
			"avg_duration":      avgDuration,
			"total_errors":      totalErrors,
			"success_rate":      successRate,
		},
	})
}
