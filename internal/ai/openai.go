package ai

import (
	"context"
	"fmt"
	"time"

	"ai-service/internal/models"
	"github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
	apiKey string
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	var client *openai.Client
	if apiKey != "" {
		client = openai.NewClient(apiKey)
	}
	
	return &OpenAIProvider{
		client: client,
		apiKey: apiKey,
	}
}

func (p *OpenAIProvider) Generate(ctx context.Context, req *models.GenerationRequest) (*models.GenerationResponse, error) {
	if !p.IsAvailable() {
		return nil, fmt.Errorf("OpenAI provider not available - check API key configuration")
	}

	startTime := time.Now()
	
	// Set default model if not specified
	model := req.Model
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}

	// Prepare messages
	messages := []openai.ChatCompletionMessage{}
	
	if req.SystemMsg != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemMsg,
		})
	}
	
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Prompt,
	})

	// Prepare request
	chatReq := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	if req.MaxTokens > 0 {
		chatReq.MaxTokens = req.MaxTokens
	}
	
	if req.Temperature > 0 {
		chatReq.Temperature = req.Temperature
	}

	// Make API call
	resp, err := p.client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	duration := time.Since(startTime)
	
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return &models.GenerationResponse{
		ID:          resp.ID,
		Provider:    models.OpenAI,
		Model:       model,
		Content:     resp.Choices[0].Message.Content,
		TokensUsed:  resp.Usage.TotalTokens,
		GeneratedAt: time.Now(),
		Duration:    duration.String(),
	}, nil
}

func (p *OpenAIProvider) GetName() string {
	return "OpenAI"
}

func (p *OpenAIProvider) IsAvailable() bool {
	return p.client != nil && p.apiKey != ""
}

func (p *OpenAIProvider) GetSupportedModels() []string {
	return []string{
		openai.GPT3Dot5Turbo,
		openai.GPT3Dot5Turbo16K,
		openai.GPT4,
		openai.GPT4TurboPreview,
		"gpt-4o",
		"gpt-4o-mini",
	}
}

func (p *OpenAIProvider) ValidateRequest(req *models.GenerationRequest) error {
	if req.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	
	if req.MaxTokens > 4096 {
		return fmt.Errorf("max_tokens cannot exceed 4096 for most OpenAI models")
	}
	
	if req.Temperature < 0 || req.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}
	
	return nil
}