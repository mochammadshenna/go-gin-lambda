package controller

import (
	"ai-service/internal/repository"
	"ai-service/internal/util/template"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WebController interface {
	Home(c *gin.Context)
	History(c *gin.Context)
	Stats(c *gin.Context)
	Error(c *gin.Context)
}

type webController struct {
	generationRepo repository.GenerationRepository
}

func NewWebController(generationRepo repository.GenerationRepository) WebController {
	return &webController{
		generationRepo: generationRepo,
	}
}

// Home renders the home page with AI generation interface
func (c *webController) Home(ctx *gin.Context) {
	providers := map[string]gin.H{
		"openai": {
			"Name":      "OpenAI",
			"Models":    []string{"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo"},
			"Strengths": "Advanced language models for text generation",
			"BestFor":   "Text generation, coding, analysis",
			"MaxTokens": "4096",
			"Pricing":   "Pay per token",
		},
		"gemini": {
			"Name":      "Google Gemini",
			"Models":    []string{"gemini-1.5-flash", "gemini-1.5-pro", "gemini-2.0-flash"},
			"Strengths": "Multimodal AI with fast response times",
			"BestFor":   "Multimodal tasks, quick responses",
			"MaxTokens": "8192",
			"Pricing":   "Pay per token",
		},
		"anthropic": {
			"Name":      "Anthropic Claude",
			"Models":    []string{"claude-3-sonnet", "claude-3-opus", "claude-3-haiku"},
			"Strengths": "Constitutional AI for safe and helpful responses",
			"BestFor":   "Safe content, detailed analysis",
			"MaxTokens": "8192",
			"Pricing":   "Pay per token",
		},
	}

	// Create provider models mapping for JavaScript
	providerModels := make(map[string][]string)
	for key, provider := range providers {
		if models, ok := provider["Models"].([]string); ok {
			providerModels[key] = models
		}
	}

	// Convert to JSON for JavaScript
	providersJSON, _ := json.Marshal(providers)
	providerModelsJSON, _ := json.Marshal(providerModels)

	data := gin.H{
		"Title":              "Home",
		"Providers":          providers,
		"ProvidersJSON":      string(providersJSON),
		"ProviderModelsJSON": string(providerModelsJSON),
	}

	html, err := template.ExecuteTemplate("home_standalone.html", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to load home page: "+err.Error())
		return
	}

	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, html)
}

// History renders the history page with generation records
func (c *webController) History(ctx *gin.Context) {
	generations, err := c.generationRepo.GetRecent(ctx, 50, 0) // Get last 50 generations
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to load generation history: "+err.Error())
		return
	}

	data := gin.H{
		"Title":       "Generation History",
		"Generations": generations,
	}

	html, err := template.ExecuteTemplate("history_standalone.html", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to load history page: "+err.Error())
		return
	}

	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, html)
}

// Stats renders the statistics page with usage metrics
func (c *webController) Stats(ctx *gin.Context) {
	// Get stats for the last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	providerStats, err := c.generationRepo.GetStats(ctx, startDate, endDate)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to load usage statistics: "+err.Error())
		return
	}

	// Calculate totals
	var totalGenerations, totalTokens int
	providerStatsMap := make(map[string]gin.H)

	for _, stat := range providerStats {
		totalGenerations += stat.TotalGenerations
		totalTokens += stat.TotalTokens
		providerStatsMap[stat.Provider] = gin.H{
			"Count":  stat.TotalGenerations,
			"Tokens": stat.TotalTokens,
		}
	}

	// Format data for template
	statsData := gin.H{
		"TotalGenerations":           totalGenerations,
		"TotalTokensUsed":            totalTokens,
		"AverageDuration":            0, // TODO: Calculate from actual data
		"DaysActive":                 30,
		"ProviderStats":              providerStatsMap,
		"RecentActivity":             []gin.H{}, // TODO: Add recent activity
		"SuccessRate":                100,       // TODO: Calculate from actual data
		"AverageTokensPerGeneration": 0,         // TODO: Calculate from actual data
		"MostUsedProvider":           "N/A",     // TODO: Calculate from actual data
	}

	data := gin.H{
		"Title": "Usage Statistics",
		"Stats": statsData,
	}

	html, err := template.ExecuteTemplate("stats_standalone.html", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to load statistics page: "+err.Error())
		return
	}

	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, html)
}

// Error renders the error page
func (c *webController) Error(ctx *gin.Context) {
	data := gin.H{
		"Title":   "Error",
		"Message": "An error occurred while processing your request.",
	}

	html, err := template.ExecuteTemplate("error.html", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, html)
}
