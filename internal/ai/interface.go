package ai

import (
	"context"
	"ai-service/internal/models"
)

// Provider defines the interface that all AI providers must implement
type Provider interface {
	// Generate generates content using the AI provider
	Generate(ctx context.Context, req *models.GenerationRequest) (*models.GenerationResponse, error)
	
	// GetName returns the provider name
	GetName() string
	
	// IsAvailable checks if the provider is configured and available
	IsAvailable() bool
	
	// GetSupportedModels returns list of supported models
	GetSupportedModels() []string
	
	// ValidateRequest validates the generation request for this provider
	ValidateRequest(req *models.GenerationRequest) error
}