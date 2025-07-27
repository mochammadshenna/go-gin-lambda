package entity

import (
	"context"
	"time"
)

// AIProvider represents the core domain entity for AI providers
type AIProvider struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	APIKey      string    `json:"-"` // Never expose in JSON
	IsAvailable bool      `json:"is_available"`
	MaxTokens   int       `json:"max_tokens"`
	Models      []string  `json:"models"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GenerationRequest represents a request to generate content
type GenerationRequest struct {
	ID           string    `json:"id"`
	ProviderType string    `json:"provider_type"`
	Model        string    `json:"model"`
	Prompt       string    `json:"prompt"`
	SystemMsg    string    `json:"system_msg,omitempty"`
	MaxTokens    int       `json:"max_tokens"`
	Temperature  float64   `json:"temperature"`
	UserID       string    `json:"user_id,omitempty"`
	RequestID    string    `json:"request_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// GenerationResponse represents the response from AI generation
type GenerationResponse struct {
	ID          string        `json:"id"`
	RequestID   string        `json:"request_id"`
	Provider    string        `json:"provider"`
	Model       string        `json:"model"`
	Content     string        `json:"content"`
	TokensUsed  int           `json:"tokens_used"`
	Duration    time.Duration `json:"duration"`
	GeneratedAt time.Time     `json:"generated_at"`
	Error       string        `json:"error,omitempty"`
}

// GenerationHistory represents the history of generations
type GenerationHistory struct {
	ID         string    `json:"id"`
	RequestID  string    `json:"request_id"`
	Provider   string    `json:"provider"`
	Model      string    `json:"model"`
	Prompt     string    `json:"prompt"`
	Response   string    `json:"response"`
	TokensUsed int       `json:"tokens_used"`
	Duration   int64     `json:"duration"` // in milliseconds
	UserID     string    `json:"user_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ProviderStats represents statistics for a provider
type ProviderStats struct {
	Provider string `json:"provider"`
	Count    int64  `json:"count"`
	Tokens   int64  `json:"tokens"`
	Duration int64  `json:"duration"` // average in milliseconds
}

// ServiceStats represents overall service statistics
type ServiceStats struct {
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

// RecentActivity represents recent activity items
type RecentActivity struct {
	Date        string `json:"date"`
	Description string `json:"description"`
	Provider    string `json:"provider"`
	Tokens      int    `json:"tokens"`
}

// AIProviderRepository defines the interface for AI provider data access
type AIProviderRepository interface {
	GetProvider(ctx context.Context, providerType string) (*AIProvider, error)
	GetAllProviders(ctx context.Context) ([]*AIProvider, error)
	SaveProvider(ctx context.Context, provider *AIProvider) error
	UpdateProvider(ctx context.Context, provider *AIProvider) error
	DeleteProvider(ctx context.Context, id string) error
}

// GenerationRepository defines the interface for generation data access
type GenerationRepository interface {
	SaveGeneration(ctx context.Context, history *GenerationHistory) error
	GetGenerationHistory(ctx context.Context, limit int) ([]*GenerationHistory, error)
	GetGenerationsByProvider(ctx context.Context, provider string, limit int) ([]*GenerationHistory, error)
	GetGenerationStats(ctx context.Context) (*ServiceStats, error)
	GetGenerationByID(ctx context.Context, id string) (*GenerationHistory, error)
}

// AIProviderService defines the interface for AI provider business logic
type AIProviderService interface {
	GenerateContent(ctx context.Context, req *GenerationRequest) (*GenerationResponse, error)
	CompareProviders(ctx context.Context, req *GenerationRequest, providers []string) (map[string]*GenerationResponse, error)
	GetAvailableProviders(ctx context.Context) ([]*AIProvider, error)
	GetProviderModels(ctx context.Context, providerType string) ([]string, error)
	ValidateRequest(ctx context.Context, req *GenerationRequest) error
}
