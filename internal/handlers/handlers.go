package handlers

import (
	"net/http"
	"strconv"
	"time"

	"ai-service/internal/ai"
	"ai-service/internal/config"
	"ai-service/internal/database"
	"ai-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	aiManager *ai.Manager
	db        *database.DB
	config    *config.Config
}

func New(aiManager *ai.Manager, db *database.DB, config *config.Config) *Handler {
	return &Handler{
		aiManager: aiManager,
		db:        db,
		config:    config,
	}
}

// @Summary Generate AI content
// @Description Generate content using specified AI provider
// @Tags AI Generation
// @Accept json
// @Produce json
// @Param request body models.GenerationRequest true "Generation request"
// @Success 200 {object} models.GenerationResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/generate [post]
func (h *Handler) Generate(c *gin.Context) {
	var req models.GenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// Use default provider if not specified
	if req.Provider == "" {
		req.Provider = h.aiManager.GetDefaultProvider()
	}

	startTime := time.Now()
	resp, err := h.aiManager.Generate(c.Request.Context(), &req)
	duration := time.Since(startTime)

	if err != nil {
		log.Error().Err(err).Str("provider", string(req.Provider)).Msg("Generation failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Generation failed",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	// Save to database
	history := &models.GenerationHistory{
		Provider:   string(resp.Provider),
		Model:      resp.Model,
		Prompt:     req.Prompt,
		Response:   resp.Content,
		TokensUsed: resp.TokensUsed,
		Duration:   duration.Milliseconds(),
		CreatedAt:  time.Now(),
	}

	if err := h.db.SaveGeneration(history); err != nil {
		log.Error().Err(err).Msg("Failed to save generation history")
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Compare AI providers
// @Description Generate content using multiple AI providers for comparison
// @Tags AI Generation
// @Accept json
// @Produce json
// @Param request body models.ComparisonRequest true "Comparison request"
// @Success 200 {object} models.ComparisonResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/compare [post]
func (h *Handler) Compare(c *gin.Context) {
	var req models.ComparisonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	resp, err := h.aiManager.Compare(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Comparison failed")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Comparison failed",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	// Save each result to database
	for _, result := range resp.Results {
		history := &models.GenerationHistory{
			Provider:   string(result.Provider),
			Model:      result.Model,
			Prompt:     req.Prompt,
			Response:   result.Content,
			TokensUsed: result.TokensUsed,
			Duration:   0, // Will be calculated from result.Duration if needed
			CreatedAt:  time.Now(),
		}

		if err := h.db.SaveGeneration(history); err != nil {
			log.Error().Err(err).Str("provider", string(result.Provider)).Msg("Failed to save generation history")
		}
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Get AI providers comparison
// @Description Get detailed comparison of available AI providers
// @Tags AI Providers
// @Produce json
// @Success 200 {object} map[string]models.AIProviderComparison
// @Router /api/v1/providers [get]
func (h *Handler) GetProviders(c *gin.Context) {
	providers := h.aiManager.GetAvailableProviders()
	c.JSON(http.StatusOK, providers)
}

// @Summary Get generation history
// @Description Get recent generation history with optional filtering
// @Tags History
// @Produce json
// @Param limit query int false "Number of records to return" default(50)
// @Param provider query string false "Filter by provider"
// @Success 200 {array} models.GenerationHistory
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/history [get]
func (h *Handler) GetHistory(c *gin.Context) {
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
		log.Error().Err(err).Msg("Failed to get generation history")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get history",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, histories)
}

// @Summary Get generation statistics
// @Description Get statistics about AI generations
// @Tags Statistics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.db.GetGenerationStats()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get generation stats")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get statistics",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary Health check
// @Description Check service health and provider availability
// @Tags Health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /api/v1/health [get]
func (h *Handler) Health(c *gin.Context) {
	providers := make(map[string]string)
	
	for providerType := range h.aiManager.GetAvailableProviders() {
		providers[providerType] = "available"
	}

	response := models.HealthResponse{
		Status:    "healthy",
		Version:   h.config.ServiceVersion,
		Timestamp: time.Now(),
		Providers: providers,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get specific AI commands
// @Description Get examples of specific commands for different AI tasks
// @Tags AI Commands
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/commands [get]
func (h *Handler) GetCommands(c *gin.Context) {
	commands := map[string]interface{}{
		"code_generation": map[string]interface{}{
			"description": "Generate code in various programming languages",
			"examples": []map[string]interface{}{
				{
					"prompt": "Write a REST API endpoint in Go using Gin framework for user registration",
					"system_message": "You are an expert Go developer. Write clean, efficient, and well-documented code.",
					"expected_output": "Complete Go code with proper error handling and validation",
				},
				{
					"prompt": "Create a React component for a todo list with add, delete, and toggle functionality",
					"system_message": "You are a senior frontend developer. Use modern React patterns and TypeScript.",
					"expected_output": "TypeScript React component with hooks and proper typing",
				},
			},
		},
		"data_analysis": map[string]interface{}{
			"description": "Analyze data and provide insights",
			"examples": []map[string]interface{}{
				{
					"prompt": "Analyze this sales data and provide insights: [JSON data]",
					"system_message": "You are a data analyst. Provide clear, actionable insights with supporting evidence.",
					"expected_output": "Structured analysis with key findings and recommendations",
				},
			},
		},
		"documentation": map[string]interface{}{
			"description": "Generate technical documentation",
			"examples": []map[string]interface{}{
				{
					"prompt": "Write API documentation for the AI service endpoints",
					"system_message": "You are a technical writer. Create comprehensive, user-friendly documentation.",
					"expected_output": "Well-structured API documentation with examples",
				},
			},
		},
		"problem_solving": map[string]interface{}{
			"description": "Help solve complex technical problems",
			"examples": []map[string]interface{}{
				{
					"prompt": "Debug this error: [error message and code context]",
					"system_message": "You are a debugging expert. Provide step-by-step solutions.",
					"expected_output": "Root cause analysis and solution steps",
				},
			},
		},
		"optimization": map[string]interface{}{
			"description": "Optimize code, queries, or processes",
			"examples": []map[string]interface{}{
				{
					"prompt": "Optimize this SQL query for better performance: [SQL query]",
					"system_message": "You are a database optimization expert. Focus on performance and maintainability.",
					"expected_output": "Optimized query with explanation of improvements",
				},
			},
		},
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"commands": commands,
		"usage_tips": []string{
			"Be specific in your prompts for better results",
			"Use system messages to set the AI's role and context",
			"Include relevant context and examples in your prompts",
			"Specify the desired output format (code, markdown, JSON, etc.)",
			"Use appropriate temperature settings (0.1 for factual, 0.7 for creative)",
		},
		"best_practices": []string{
			"Break complex tasks into smaller, focused prompts",
			"Provide clear requirements and constraints",
			"Use the comparison endpoint to evaluate different providers",
			"Review and validate AI-generated code before using in production",
			"Consider token limits when working with large inputs",
		},
	})
}