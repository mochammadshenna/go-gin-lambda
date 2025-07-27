package controller

import (
	"ai-service/internal/util/template"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebController interface {
	Home(c *gin.Context)
	History(c *gin.Context)
	Stats(c *gin.Context)
	Error(c *gin.Context)
}

type webController struct {
	// service will be injected
}

func NewWebController(service interface{}) WebController {
	return &webController{}
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
	data := gin.H{
		"Title": "Generation History",
		"Generations": []gin.H{
			{
				"ID":         "gen-123",
				"Provider":   "OpenAI",
				"Model":      "gpt-3.5-turbo",
				"Prompt":     "Write a hello world program",
				"Response":   "Here's a simple hello world program...",
				"TokensUsed": 150,
				"Duration":   "2.5s",
				"CreatedAt":  "2024-01-15 10:30:00",
				"Status":     "success",
			},
		},
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
	data := gin.H{
		"Title": "Usage Statistics",
		"Stats": map[string]gin.H{
			"openai": {
				"TotalGenerations": 150,
				"TotalTokens":      45000,
				"AvgDuration":      2.5,
				"ErrorCount":       3,
			},
			"gemini": {
				"TotalGenerations": 75,
				"TotalTokens":      25000,
				"AvgDuration":      1.8,
				"ErrorCount":       1,
			},
			"anthropic": {
				"TotalGenerations": 25,
				"TotalTokens":      8000,
				"AvgDuration":      3.2,
				"ErrorCount":       0,
			},
		},
	}

	html, err := template.ExecuteTemplate("stats_standalone.html", data)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title":   "Error",
			"Message": "Failed to load statistics page",
		})
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
