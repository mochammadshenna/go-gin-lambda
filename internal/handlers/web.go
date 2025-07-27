package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"

	"ai-service/internal/ai"
	"ai-service/internal/config"
	"ai-service/internal/database"
	"ai-service/internal/models"

	"github.com/gin-gonic/gin"
)

type WebHandler struct {
	aiManager *ai.Manager
	db        *database.DB
	config    *config.Config
	templates *template.Template
}

func NewWebHandler(aiManager *ai.Manager, db *database.DB, config *config.Config) *WebHandler {
	// Create template with custom functions first
	tmpl := template.New("").Funcs(template.FuncMap{
		"percent": func(count, total int64) float64 {
			if total == 0 {
				return 0
			}
			return float64(count) / float64(total) * 100
		},
	})

	// Parse templates with error handling
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse templates: %v", err))
	}

	return &WebHandler{
		aiManager: aiManager,
		db:        db,
		config:    config,
		templates: tmpl,
	}
}

// Home page handler
func (h *WebHandler) Home(c *gin.Context) {
	// Get actual provider data
	providers := h.aiManager.GetAvailableProviders()

	// Convert providers to template-friendly format
	providerData := make(map[string]map[string]interface{})
	providerModels := make(map[string][]string)

	for providerType, providerInfo := range providers {
		// Convert slices to strings for display
		strengths := ""
		if len(providerInfo.Strengths) > 0 {
			strengths = providerInfo.Strengths[0] // Use first strength for display
		}

		bestFor := ""
		if len(providerInfo.BestFor) > 0 {
			bestFor = providerInfo.BestFor[0] // Use first best for use case
		}

		providerData[providerType] = map[string]interface{}{
			"Name":       providerInfo.Name,
			"Strengths":  strengths,
			"Weaknesses": providerInfo.Weaknesses,
			"BestFor":    bestFor,
			"Pricing":    providerInfo.Pricing,
			"MaxTokens":  providerInfo.MaxTokens,
			"Available":  providerInfo.Available,
		}

		// Get provider models
		if provider, err := h.aiManager.GetProvider(providerInfo.Provider); err == nil {
			providerModels[providerType] = provider.GetSupportedModels()
		}
	}

	// Convert to JSON for JavaScript
	providersJSON, _ := json.Marshal(providerData)
	providerModelsJSON, _ := json.Marshal(providerModels)

	err := h.templates.ExecuteTemplate(c.Writer, "home_standalone.html", gin.H{
		"Title":              "AI Content Generation",
		"Providers":          providerData,
		"ProvidersJSON":      string(providersJSON),
		"ProviderModelsJSON": string(providerModelsJSON),
	})
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		panic(err)
	}
}

// History page handler
func (h *WebHandler) History(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	provider := c.Query("provider")

	var histories []models.GenerationHistory
	if provider != "" {
		histories, err = h.db.GetGenerationsByProvider(provider, limit)
	} else {
		histories, err = h.db.GetGenerationHistory(limit)
	}

	if err != nil {
		err := h.templates.ExecuteTemplate(c.Writer, "error.html", gin.H{
			"Title": "Error",
			"Error": "Failed to load history: " + err.Error(),
		})
		if err != nil {
			fmt.Printf("Template execution error: %v\n", err)
			panic(err)
		}
		return
	}

	// Get provider names for display
	providers := h.aiManager.GetAvailableProviders()
	providerNames := make(map[string]string)
	for providerType, providerInfo := range providers {
		providerNames[providerType] = providerInfo.Name
	}

	err = h.templates.ExecuteTemplate(c.Writer, "history_standalone.html", gin.H{
		"Title":         "Generation History",
		"Histories":     histories,
		"Providers":     providerNames,
		"CurrentFilter": provider,
	})
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		panic(err)
	}
}

// Stats page handler
func (h *WebHandler) Stats(c *gin.Context) {
	rawStats, err := h.db.GetGenerationStats()
	if err != nil {
		err := h.templates.ExecuteTemplate(c.Writer, "error.html", gin.H{
			"Title": "Error",
			"Error": "Failed to load statistics: " + err.Error(),
		})
		if err != nil {
			fmt.Printf("Template execution error: %v\n", err)
			panic(err)
		}
		return
	}

	// Transform stats to match template expectations
	totalGenerations := rawStats["total_generations"].(int64)
	avgTokensUsed := rawStats["avg_tokens_used"].(float64)
	avgDuration := rawStats["avg_duration_ms"].(float64)

	stats := gin.H{
		"TotalGenerations":           totalGenerations,
		"TotalTokensUsed":            int(avgTokensUsed * float64(totalGenerations)),
		"AverageDuration":            int(avgDuration),
		"DaysActive":                 1,   // Default to 1 for now
		"SuccessRate":                100, // Default to 100% for now
		"AverageTokensPerGeneration": int(avgTokensUsed),
		"MostUsedProvider":           "gemini", // Default for now
		"ProviderStats":              make(map[string]gin.H),
		"RecentActivity":             []gin.H{},
	}

	// Transform provider stats
	if byProvider, ok := rawStats["by_provider"].([]struct {
		Provider string
		Count    int64
	}); ok {
		for _, provider := range byProvider {
			stats["ProviderStats"].(map[string]gin.H)[provider.Provider] = gin.H{
				"Count": provider.Count,
			}
		}
	}

	err = h.templates.ExecuteTemplate(c.Writer, "stats_standalone.html", gin.H{
		"Title": "Statistics",
		"Stats": stats,
	})
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		panic(err)
	}
}

// Test page handler
func (h *WebHandler) Test(c *gin.Context) {
	fmt.Println("Test handler called")
	err := h.templates.ExecuteTemplate(c.Writer, "test.html", gin.H{
		"Title": "Test Page",
	})
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		panic(err)
	}
}

// Error page handler
func (h *WebHandler) Error(c *gin.Context) {
	errorMsg := c.Query("error")
	if errorMsg == "" {
		errorMsg = "An unknown error occurred"
	}

	err := h.templates.ExecuteTemplate(c.Writer, "error.html", gin.H{
		"Title": "Error",
		"Error": errorMsg,
	})
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		panic(err)
	}
}
