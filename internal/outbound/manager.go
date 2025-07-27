package outbound

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"ai-service/cmd/config"
	"ai-service/internal/model"
)

type Manager struct {
	providers map[model.AIProvider]Provider
	config    *config.Config
	mu        sync.RWMutex
}

func NewManager(cfg *config.Config) *Manager {
	manager := &Manager{
		providers: make(map[model.AIProvider]Provider),
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
	if m.config.AIProviders.OpenAI.APIKey != "" && !isPlaceholderAPIKey(m.config.AIProviders.OpenAI.APIKey) {
		m.providers[model.OpenAI] = NewOpenAIProvider(m.config.AIProviders.OpenAI.APIKey)
	}

	// Initialize Gemini provider
	if m.config.AIProviders.Gemini.APIKey != "" && !isPlaceholderAPIKey(m.config.AIProviders.Gemini.APIKey) {
		m.providers[model.Gemini] = NewGeminiProvider(m.config.AIProviders.Gemini.APIKey)
	}

	// Note: Anthropic provider would be initialized here when implemented
	// if m.config.AIProviders.Anthropic.APIKey != "" && !isPlaceholderAPIKey(m.config.AIProviders.Anthropic.APIKey) {
	//     m.providers[model.Anthropic] = NewAnthropicProvider(m.config.AIProviders.Anthropic.APIKey)
	// }
}

// isPlaceholderAPIKey checks if the API key is a placeholder value
func isPlaceholderAPIKey(apiKey string) bool {
	placeholders := []string{
		"your_openai_api_key_here",
		"your_gemini_api_key_here",
		"your_anthropic_api_key_here",
		"your-api-key-here",
		"your_api_key_here",
		"placeholder",
		"",
	}

	for _, placeholder := range placeholders {
		if apiKey == placeholder {
			return true
		}
	}

	// Check if it contains placeholder-like patterns
	if len(apiKey) < 10 || strings.Contains(apiKey, "your_") || strings.Contains(apiKey, "placeholder") {
		return true
	}

	return false
}

func (m *Manager) Generate(ctx context.Context, req *model.GenerationRequest) (*model.GenerationResponse, error) {
	m.mu.RLock()
	provider, exists := m.providers[req.Provider]
	m.mu.RUnlock()

	if !exists {
		// Check if any providers are configured
		if len(m.providers) == 0 {
			return nil, fmt.Errorf("no AI providers configured. Please set at least one API key (OPENAI_API_KEY, GEMINI_API_KEY, or ANTHROPIC_API_KEY)")
		}
		return nil, fmt.Errorf("provider %s not found. Available providers: %v", req.Provider, m.getAvailableProviderNames())
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

func (m *Manager) Compare(ctx context.Context, req *model.ComparisonRequest) (*model.ComparisonResponse, error) {
	var results []model.GenerationResponse
	var wg sync.WaitGroup
	var mu sync.Mutex

	errChan := make(chan error, len(req.Providers))

	for _, providerType := range req.Providers {
		wg.Add(1)
		go func(pt model.AIProvider) {
			defer wg.Done()

			genReq := &model.GenerationRequest{
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

	return &model.ComparisonResponse{
		Prompt:    req.Prompt,
		Results:   results,
		CreatedAt: results[0].GeneratedAt,
	}, nil
}

func (m *Manager) GetAvailableProviders() map[string]model.AIProviderComparison {
	m.mu.RLock()
	defer m.mu.RUnlock()

	comparisons := make(map[string]model.AIProviderComparison)

	// OpenAI comparison
	openaiProvider, hasOpenAI := m.providers[model.OpenAI]
	comparisons["openai"] = model.AIProviderComparison{
		Provider: model.OpenAI,
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
	geminiProvider, hasGemini := m.providers[model.Gemini]
	comparisons["gemini"] = model.AIProviderComparison{
		Provider: model.Gemini,
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
	comparisons["anthropic"] = model.AIProviderComparison{
		Provider: model.Anthropic,
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

func (m *Manager) GetProvider(providerType model.AIProvider) (Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, exists := m.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerType)
	}

	return provider, nil
}

func (m *Manager) GetDefaultProvider() model.AIProvider {
	// For now, return OpenAI as default
	return model.OpenAI
}

// getAvailableProviderNames returns a list of available provider names
func (m *Manager) getAvailableProviderNames() []string {
	var names []string
	for provider := range m.providers {
		names = append(names, string(provider))
	}
	return names
}
