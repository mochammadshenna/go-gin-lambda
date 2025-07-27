package outbound

import (
	"ai-service/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAIProvider struct {
	apiKey string
	client *http.Client
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *OpenAIProvider) Generate(ctx context.Context, req *model.GenerationRequest) (*model.GenerationResponse, error) {
	startTime := time.Now()

	// Set default model if not specified
	modelName := req.Model
	if modelName == "" {
		modelName = "gpt-3.5-turbo"
	}

	// Prepare request payload
	payload := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": req.Prompt,
			},
		},
	}

	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}

	if req.Temperature > 0 {
		payload["temperature"] = req.Temperature
	}

	if req.SystemMsg != "" {
		// Add system message as the first message
		systemMsg := map[string]string{
			"role":    "system",
			"content": req.SystemMsg,
		}
		payload["messages"] = append([]map[string]string{systemMsg}, payload["messages"].([]map[string]string)...)
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	// Make request
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	duration := time.Since(startTime)

	return &model.GenerationResponse{
		ID:          fmt.Sprintf("openai-%d", time.Now().UnixNano()),
		Provider:    model.OpenAI,
		Model:       modelName,
		Content:     openAIResp.Choices[0].Message.Content,
		TokensUsed:  openAIResp.Usage.TotalTokens,
		GeneratedAt: time.Now(),
		Duration:    duration.String(),
	}, nil
}

func (p *OpenAIProvider) GetName() string {
	return "OpenAI"
}

func (p *OpenAIProvider) IsAvailable() bool {
	return p.apiKey != ""
}

func (p *OpenAIProvider) GetSupportedModels() []string {
	return []string{
		"gpt-4",
		"gpt-4-turbo",
		"gpt-3.5-turbo",
		"gpt-3.5-turbo-16k",
	}
}

func (p *OpenAIProvider) ValidateRequest(req *model.GenerationRequest) error {
	if req.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}

	if req.MaxTokens > 4096 {
		return fmt.Errorf("max_tokens cannot exceed 4096 for OpenAI models")
	}

	if req.Temperature < 0 || req.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2 for OpenAI")
	}

	return nil
}
