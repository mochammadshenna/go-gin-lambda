package api

import (
	"time"
)

// GenerationRequest represents a request to generate content
type GenerationRequest struct {
	Provider    string  `json:"provider" validate:"required"`
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt" validate:"required"`
	SystemMsg   string  `json:"system_msg,omitempty"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	UserID      string  `json:"user_id,omitempty"`
}

// GenerationResponse represents the response from AI generation
type GenerationResponse struct {
	ID          string        `json:"id"`
	Provider    string        `json:"provider"`
	Model       string        `json:"model"`
	Content     string        `json:"content"`
	TokensUsed  int           `json:"tokens_used"`
	Duration    time.Duration `json:"duration"`
	GeneratedAt time.Time     `json:"generated_at"`
	Error       string        `json:"error,omitempty"`
}

// CompareRequest represents a request to compare providers
type CompareRequest struct {
	Prompt      string   `json:"prompt" validate:"required"`
	Providers   []string `json:"providers" validate:"required"`
	MaxTokens   int      `json:"max_tokens"`
	Model       string   `json:"model"`
	SystemMsg   string   `json:"system_msg,omitempty"`
	Temperature float64  `json:"temperature"`
}

// CompareResponse represents the comparison response
type CompareResponse struct {
	Results map[string]GenerationResponse `json:"results"`
}

// ProviderInfo represents information about an AI provider
type ProviderInfo struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Models      []string `json:"models"`
	IsAvailable bool     `json:"is_available"`
	MaxTokens   int      `json:"max_tokens"`
}

// ProvidersResponse represents the providers response
type ProvidersResponse struct {
	Providers []ProviderInfo `json:"providers"`
}

// HistoryRequest represents a request for generation history
type HistoryRequest struct {
	Limit    int    `json:"limit"`
	Provider string `json:"provider,omitempty"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// HistoryItem represents a history item
type HistoryItem struct {
	ID         string    `json:"id"`
	Provider   string    `json:"provider"`
	Model      string    `json:"model"`
	Prompt     string    `json:"prompt"`
	Response   string    `json:"response"`
	TokensUsed int       `json:"tokens_used"`
	Duration   int64     `json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
}

// HistoryResponse represents the history response
type HistoryResponse struct {
	Items []HistoryItem `json:"items"`
	Total int           `json:"total"`
}

// StatsResponse represents the statistics response
type StatsResponse struct {
	TotalGenerations           int64                    `json:"total_generations"`
	TotalTokensUsed            int64                    `json:"total_tokens_used"`
	AverageDuration            int64                    `json:"average_duration"`
	DaysActive                 int                      `json:"days_active"`
	SuccessRate                float64                  `json:"success_rate"`
	AverageTokensPerGeneration int                      `json:"average_tokens_per_generation"`
	MostUsedProvider           string                   `json:"most_used_provider"`
	ProviderStats              map[string]ProviderStats `json:"provider_stats"`
	RecentActivity             []RecentActivity         `json:"recent_activity"`
}

// ProviderStats represents statistics for a provider
type ProviderStats struct {
	Provider string `json:"provider"`
	Count    int64  `json:"count"`
	Tokens   int64  `json:"tokens"`
	Duration int64  `json:"duration"`
}

// RecentActivity represents recent activity items
type RecentActivity struct {
	Date        string `json:"date"`
	Description string `json:"description"`
	Provider    string `json:"provider"`
	Tokens      int    `json:"tokens"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}
