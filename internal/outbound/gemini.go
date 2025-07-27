package outbound

import (
	"context"
	"fmt"
	"time"

	"ai-service/internal/model"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiProvider struct {
	client *genai.Client
	apiKey string
}

func NewGeminiProvider(apiKey string) *GeminiProvider {
	return &GeminiProvider{
		apiKey: apiKey,
	}
}

func (p *GeminiProvider) initClient(ctx context.Context) error {
	if p.client != nil {
		return nil
	}

	if p.apiKey == "" {
		return fmt.Errorf("Gemini API key not configured")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(p.apiKey))
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}

	p.client = client
	return nil
}

func (p *GeminiProvider) Generate(ctx context.Context, req *model.GenerationRequest) (*model.GenerationResponse, error) {
	if err := p.initClient(ctx); err != nil {
		return nil, err
	}

	startTime := time.Now()

	// Set default model if not specified
	modelName := req.Model
	if modelName == "" {
		modelName = "gemini-1.5-flash"
	}

	// Get the model
	geminiModel := p.client.GenerativeModel(modelName)

	// Configure generation parameters
	if req.MaxTokens > 0 {
		geminiModel.SetMaxOutputTokens(int32(req.MaxTokens))
	}

	if req.Temperature > 0 {
		geminiModel.SetTemperature(req.Temperature)
	}

	// Prepare prompt
	prompt := req.Prompt
	if req.SystemMsg != "" {
		prompt = fmt.Sprintf("System: %s\n\nUser: %s", req.SystemMsg, req.Prompt)
	}

	// Generate content
	resp, err := geminiModel.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("Gemini API error: %w", err)
	}

	duration := time.Since(startTime)

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	// Extract content
	var content string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			content += string(text)
		}
	}

	// Calculate tokens used (approximate - field may not be available in all versions)
	tokensUsed := 0
	// Note: UsageMetadata field may not be available in all versions of the Gemini API
	// For now, we'll estimate based on content length
	tokensUsed = len(content) / 4 // Rough estimation: 1 token ≈ 4 characters

	return &model.GenerationResponse{
		ID:          fmt.Sprintf("gemini-%d", time.Now().UnixNano()),
		Provider:    model.Gemini,
		Model:       modelName,
		Content:     content,
		TokensUsed:  tokensUsed,
		GeneratedAt: time.Now(),
		Duration:    duration.String(),
	}, nil
}

func (p *GeminiProvider) GetName() string {
	return "Google Gemini"
}

func (p *GeminiProvider) IsAvailable() bool {
	return p.apiKey != ""
}

func (p *GeminiProvider) GetSupportedModels() []string {
	return []string{
		"gemini-1.5-flash",
		"gemini-1.5-pro",
		"gemini-1.0-pro",
		"gemini-pro-vision",
	}
}

func (p *GeminiProvider) ValidateRequest(req *model.GenerationRequest) error {
	if req.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}

	if req.MaxTokens > 8192 {
		return fmt.Errorf("max_tokens cannot exceed 8192 for Gemini models")
	}

	if req.Temperature < 0 || req.Temperature > 1 {
		return fmt.Errorf("temperature must be between 0 and 1 for Gemini")
	}

	return nil
}
