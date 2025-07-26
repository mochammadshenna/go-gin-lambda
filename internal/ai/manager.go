package ai

import (
	"context"
	"fmt"
	"sync"

	"ai-service/internal/config"
	"ai-service/internal/models"
)

type Manager struct {
	providers map[models.AIProvider]Provider
	config    *config.Config
	mu        sync.RWMutex
}

func NewManager(cfg *config.Config) *Manager {
	manager := &Manager{
		providers: make(map[models.AIProvider]Provider),
		config:    cfg,
	}
	
	// Initialize providers
	manager.initProviders()
	
	return manager
}

func (m *Manager) initProviders() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Initialize OpenAI provider
	if m.config.OpenAIAPIKey != "" {
		m.providers[models.OpenAI] = NewOpenAIProvider(m.config.OpenAIAPIKey)
	}
	
	// Initialize Gemini provider
	if m.config.GeminiAPIKey != "" {
		m.providers[models.Gemini] = NewGeminiProvider(m.config.GeminiAPIKey)
	}
	
	// Note: Anthropic provider would be initialized here when implemented
	// if m.config.AnthropicAPIKey != "" {
	//     m.providers[models.Anthropic] = NewAnthropicProvider(m.config.AnthropicAPIKey)
	// }
}

func (m *Manager) Generate(ctx context.Context, req *models.GenerationRequest) (*models.GenerationResponse, error) {
	m.mu.RLock()
	provider, exists := m.providers[req.Provider]
	m.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("provider %s not found", req.Provider)
	}
	
	if !provider.IsAvailable() {
		return nil, fmt.Errorf("provider %s is not available", req.Provider)
	}
	
	// Validate request
	if err := provider.ValidateRequest(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	// Generate content
	return provider.Generate(ctx, req)
}

func (m *Manager) Compare(ctx context.Context, req *models.ComparisonRequest) (*models.ComparisonResponse, error) {
	var results []models.GenerationResponse
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	errChan := make(chan error, len(req.Providers))
	
	for _, providerType := range req.Providers {
		wg.Add(1)
		go func(pt models.AIProvider) {
			defer wg.Done()
			
			genReq := &models.GenerationRequest{
				Provider:    pt,
				Prompt:      req.Prompt,
				MaxTokens:   req.MaxTokens,
				Temperature: req.Temperature,
			}
			
			resp, err := m.Generate(ctx, genReq)
			if err != nil {
				errChan <- fmt.Errorf("provider %s failed: %w", pt, err)
				return
			}
			
			mu.Lock()
			results = append(results, *resp)
			mu.Unlock()
		}(providerType)
	}
	
	wg.Wait()
	close(errChan)
	
	// Check for errors
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}
	
	if len(errs) > 0 && len(results) == 0 {
		return nil, fmt.Errorf("all providers failed: %v", errs)
	}
	
	return &models.ComparisonResponse{
		Prompt:    req.Prompt,
		Results:   results,
		CreatedAt: results[0].GeneratedAt,
	}, nil
}

func (m *Manager) GetAvailableProviders() map[string]models.AIProviderComparison {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	comparisons := make(map[string]models.AIProviderComparison)
	
	// OpenAI comparison
	openaiProvider, hasOpenAI := m.providers[models.OpenAI]
	comparisons["openai"] = models.AIProviderComparison{
		Provider: models.OpenAI,
		Name:     "OpenAI GPT",
		Strengths: []string{
			"Excellent general knowledge",
			"Strong reasoning capabilities",
			"Good code generation",
			"Wide range of models",
			"Reliable API",
		},
		Weaknesses: []string{
			"Can be expensive for high usage",
			"Knowledge cutoff limitations",
			"Rate limiting on free tier",
		},
		BestFor: []string{
			"General text generation",
			"Code completion and debugging",
			"Creative writing",
			"Question answering",
		},
		Pricing:   "Pay per token (~$0.002/1K tokens for GPT-3.5)",
		MaxTokens: 4096,
		Available: hasOpenAI && openaiProvider.IsAvailable(),
	}
	
	// Gemini comparison
	geminiProvider, hasGemini := m.providers[models.Gemini]
	comparisons["gemini"] = models.AIProviderComparison{
		Provider: models.Gemini,
		Name:     "Google Gemini",
		Strengths: []string{
			"Multimodal capabilities",
			"Large context window",
			"Good at reasoning tasks",
			"Free tier available",
			"Fast response times",
		},
		Weaknesses: []string{
			"Less mature than OpenAI",
			"Limited third-party integrations",
			"Newer model with less community support",
		},
		BestFor: []string{
			"Multimodal tasks",
			"Long document analysis",
			"Research and analysis",
			"Cost-effective solutions",
		},
		Pricing:   "Free tier available, pay per token for pro usage",
		MaxTokens: 8192,
		Available: hasGemini && geminiProvider.IsAvailable(),
	}
	
	// Anthropic comparison (placeholder for future implementation)
	comparisons["anthropic"] = models.AIProviderComparison{
		Provider: models.Anthropic,
		Name:     "Anthropic Claude",
		Strengths: []string{
			"Strong safety focus",
			"Excellent at analysis",
			"Good refusal mechanisms",
			"Thoughtful responses",
		},
		Weaknesses: []string{
			"More conservative responses",
			"Limited availability",
			"Higher cost",
		},
		BestFor: []string{
			"Safety-critical applications",
			"Research and analysis",
			"Ethical AI use cases",
		},
		Pricing:   "Pay per token (premium pricing)",
		MaxTokens: 8192,
		Available: false, // Not implemented yet
	}
	
	return comparisons
}

func (m *Manager) GetProvider(providerType models.AIProvider) (Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	provider, exists := m.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerType)
	}
	
	return provider, nil
}

func (m *Manager) GetDefaultProvider() models.AIProvider {
	switch m.config.DefaultAIProvider {
	case "openai":
		return models.OpenAI
	case "gemini":
		return models.Gemini
	case "anthropic":
		return models.Anthropic
	default:
		return models.OpenAI
	}
}