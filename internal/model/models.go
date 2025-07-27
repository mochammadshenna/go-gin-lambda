package model

import (
	"time"

	"gorm.io/gorm"
)

// AIProvider represents the supported AI providers
type AIProvider string

const (
	OpenAI    AIProvider = "openai"
	Gemini    AIProvider = "gemini"
	Anthropic AIProvider = "anthropic"
)

// GenerationRequest represents the input for AI generation
type GenerationRequest struct {
	Provider    AIProvider `json:"provider" binding:"required" example:"openai"`
	Model       string     `json:"model" example:"gpt-3.5-turbo"`
	Prompt      string     `json:"prompt" binding:"required" example:"Write a hello world program in Go"`
	MaxTokens   int        `json:"max_tokens,omitempty" example:"1000"`
	Temperature float32    `json:"temperature,omitempty" example:"0.7"`
	SystemMsg   string     `json:"system_message,omitempty" example:"You are a helpful coding assistant"`
} // @name GenerationRequest

// GenerationResponse represents the AI generation output
type GenerationResponse struct {
	ID          string     `json:"id"`
	Provider    AIProvider `json:"provider"`
	Model       string     `json:"model"`
	Content     string     `json:"content"`
	TokensUsed  int        `json:"tokens_used"`
	GeneratedAt time.Time  `json:"generated_at"`
	Duration    string     `json:"duration"`
} // @name GenerationResponse

// ComparisonRequest for comparing AI providers
type ComparisonRequest struct {
	Prompt      string       `json:"prompt" binding:"required"`
	Providers   []AIProvider `json:"providers" binding:"required"`
	MaxTokens   int          `json:"max_tokens,omitempty"`
	Temperature float32      `json:"temperature,omitempty"`
} // @name ComparisonRequest

// ComparisonResponse contains results from multiple providers
type ComparisonResponse struct {
	Prompt    string               `json:"prompt"`
	Results   []GenerationResponse `json:"results"`
	CreatedAt time.Time            `json:"created_at"`
} // @name ComparisonResponse

// GenerationHistory stores generation history in database
type GenerationHistory struct {
	ID           string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Provider     string         `json:"provider"`
	Model        string         `json:"model"`
	Prompt       string         `json:"prompt"`
	Response     string         `json:"response"`
	TokensUsed   int            `json:"tokens_used"`
	Duration     int64          `json:"duration"` // milliseconds
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	Status       string         `json:"status"`
	ErrorMessage string         `json:"error_message,omitempty"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
} // @name ErrorResponse

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Providers map[string]string `json:"providers"`
} // @name HealthResponse

// AIProviderComparison contains comparison data for providers
type AIProviderComparison struct {
	Provider   AIProvider `json:"provider"`
	Name       string     `json:"name"`
	Strengths  []string   `json:"strengths"`
	Weaknesses []string   `json:"weaknesses"`
	BestFor    []string   `json:"best_for"`
	Pricing    string     `json:"pricing"`
	MaxTokens  int        `json:"max_tokens"`
	Available  bool       `json:"available"`
}

type ProviderStats struct {
	Provider         string  `json:"provider"`
	TotalGenerations int     `json:"total_generations"`
	TotalTokens      int     `json:"total_tokens"`
	AvgDuration      float64 `json:"avg_duration"`
	ErrorCount       int     `json:"error_count"`
}
