package outbound

import (
	"ai-service/internal/model"
	"context"
)

// Provider interface defines the contract for AI providers
type Provider interface {
	// Generate generates content using the AI provider
	Generate(ctx context.Context, req *model.GenerationRequest) (*model.GenerationResponse, error)

	// GetName returns the provider name
	GetName() string

	// IsAvailable returns whether the provider is available
	IsAvailable() bool

	// GetSupportedModels returns the list of supported models
	GetSupportedModels() []string

	// ValidateRequest validates the generation request
	ValidateRequest(req *model.GenerationRequest) error
}
