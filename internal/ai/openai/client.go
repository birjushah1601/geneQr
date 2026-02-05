package openai

import (
	"context"
	"fmt"
	"time"

	"github.com/aby-med/medical-platform/internal/ai/types"
	"github.com/aby-med/medical-platform/internal/ai/aiconfig"
	aierrors "github.com/aby-med/medical-platform/pkg/ai"
	openai "github.com/sashabaranov/go-openai"
)

// Client implements the AI Provider interface for OpenAI
type Client struct {
	client       *openai.Client
	config       aiconfig.OpenAIConfig
	providerCode string
}

// NewClient creates a new OpenAI client
func NewClient(config aiconfig.OpenAIConfig) (*Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	clientConfig := openai.DefaultConfig(config.APIKey)
	
	if config.OrgID != "" {
		clientConfig.OrgID = config.OrgID
	}
	
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	return &Client{
		client:       openai.NewClientWithConfig(clientConfig),
		config:       config,
		providerCode: "openai-gpt",
	}, nil
}

// GetName returns the provider name
func (c *Client) GetName() string {
	return "openai"
}

// GetProviderCode returns the provider code
func (c *Client) GetProviderCode() string {
	return c.providerCode
}

// Chat sends a chat completion request
func (c *Client) Chat(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	startTime := time.Now()

	// Set default model if not specified
	model := req.Model
	if model == "" {
		model = c.config.DefaultModel
	}

	// Build OpenAI request
	openaiReq := openai.ChatCompletionRequest{
		Model:    model,
		Messages: convertMessages(req.Messages),
	}

	// Set optional parameters
	if req.Temperature != nil {
		openaiReq.Temperature = *req.Temperature
	}
	if req.TopP != nil {
		openaiReq.TopP = *req.TopP
	}
	if req.MaxTokens != nil {
		openaiReq.MaxTokens = *req.MaxTokens
	}
	if len(req.StopSequences) > 0 {
		openaiReq.Stop = req.StopSequences
	}
	if req.UserID != "" {
		openaiReq.User = req.UserID
	}

	// Function calling
	if req.FunctionCalling && len(req.Functions) > 0 {
		openaiReq.Functions = convertFunctions(req.Functions)
	}

	// Call OpenAI API
	resp, err := c.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, c.wrapError(err)
	}

	latency := time.Since(startTime)

	// Handle empty response
	if len(resp.Choices) == 0 {
		return nil, aierrors.NewProviderError("openai", fmt.Errorf("no choices returned"), "", false)
	}

	choice := resp.Choices[0]

	// Build response
	response := &types.ChatResponse{
		ID:      resp.ID,
		Content: choice.Message.Content,
		Usage: types.TokenUsage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		Provider:     "openai",
		Model:        resp.Model,
		Latency:      latency,
		FinishReason: string(choice.FinishReason),
	}

	// Function call if present
	if choice.Message.FunctionCall != nil {
		response.FunctionCall = &types.FunctionCall{
			Name:      choice.Message.FunctionCall.Name,
			Arguments: choice.Message.FunctionCall.Arguments,
		}
	}

	// Calculate cost
	inputCost, outputCost, found := types.GetModelPricing(model)
	if found {
		response.Cost = types.CalculateCost(response.Usage, inputCost, outputCost)
	}

	return response, nil
}

// ChatStream sends a streaming chat completion request
func (c *Client) ChatStream(ctx context.Context, req *types.ChatRequest, streamChan chan<- *types.ChatStreamResponse) error {
	defer close(streamChan)

	// Set default model if not specified
	model := req.Model
	if model == "" {
		model = c.config.DefaultModel
	}

	// Build OpenAI request
	openaiReq := openai.ChatCompletionRequest{
		Model:    model,
		Messages: convertMessages(req.Messages),
		Stream:   true,
	}

	// Set optional parameters
	if req.Temperature != nil {
		openaiReq.Temperature = *req.Temperature
	}
	if req.TopP != nil {
		openaiReq.TopP = *req.TopP
	}
	if req.MaxTokens != nil {
		openaiReq.MaxTokens = *req.MaxTokens
	}

	// Create stream
	stream, err := c.client.CreateChatCompletionStream(ctx, openaiReq)
	if err != nil {
		streamChan <- &types.ChatStreamResponse{
			Done:  true,
			Error: c.wrapError(err),
		}
		return c.wrapError(err)
	}
	defer stream.Close()

	// Stream responses
	for {
		response, err := stream.Recv()
		if err != nil {
			// Check if stream finished normally
			if err.Error() == "EOF" {
				streamChan <- &types.ChatStreamResponse{
					Done: true,
				}
				return nil
			}

			streamChan <- &types.ChatStreamResponse{
				Done:  true,
				Error: c.wrapError(err),
			}
			return c.wrapError(err)
		}

		// Send delta
		if len(response.Choices) > 0 {
			choice := response.Choices[0]
			streamChan <- &types.ChatStreamResponse{
				Delta:        choice.Delta.Content,
				Done:         false,
				FinishReason: string(choice.FinishReason),
			}

			// Check if done
			if choice.FinishReason != "" {
				streamChan <- &types.ChatStreamResponse{
					Done:         true,
					FinishReason: string(choice.FinishReason),
				}
				return nil
			}
		}
	}
}

// Analyze analyzes an image/video with vision capabilities
func (c *Client) Analyze(ctx context.Context, req *types.VisionRequest) (*types.VisionResponse, error) {
	startTime := time.Now()

	// Set default model if not specified
	model := req.Model
	if model == "" {
		model = "gpt-4o" // Default vision model
	}

	// Build vision messages
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			MultiContent: buildVisionContent(req.Prompt, req.ImageURLs, req.ImageData),
		},
	}

	// Build OpenAI request
	openaiReq := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	// Set optional parameters
	if req.Temperature != nil {
		openaiReq.Temperature = *req.Temperature
	}
	if req.MaxTokens != nil {
		openaiReq.MaxTokens = *req.MaxTokens
	}

	// Call OpenAI API
	resp, err := c.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, c.wrapError(err)
	}

	latency := time.Since(startTime)

	// Handle empty response
	if len(resp.Choices) == 0 {
		return nil, aierrors.NewProviderError("openai", fmt.Errorf("no choices returned"), "", false)
	}

	choice := resp.Choices[0]

	// Build response
	response := &types.VisionResponse{
		Analysis: choice.Message.Content,
		Usage: types.TokenUsage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		Provider: "openai",
		Model:    resp.Model,
		Latency:  latency,
	}

	// Calculate cost
	inputCost, outputCost, found := types.GetModelPricing(model)
	if found {
		response.Cost = types.CalculateCost(response.Usage, inputCost, outputCost)
	}

	return response, nil
}

// IsHealthy checks if the provider is healthy
func (c *Client) IsHealthy(ctx context.Context) bool {
	// Simple health check: try to list models
	_, err := c.client.ListModels(ctx)
	return err == nil
}

// GetCapabilities returns provider capabilities
func (c *Client) GetCapabilities() *types.ProviderCapabilities {
	return &types.ProviderCapabilities{
		SupportsChat:             true,
		SupportsStreaming:        true,
		SupportsVision:           true,
		SupportsFunctionCalling:  true,
		SupportsEmbeddings:       true,
		MaxContextTokens:         128000, // GPT-4o
		MaxOutputTokens:          16384,
		RateLimitRequestsPerMin:  10000,
		RateLimitTokensPerMin:    2000000,
	}
}

// Close closes any open connections
func (c *Client) Close() error {
	// OpenAI client doesn't require explicit closing
	return nil
}

// Helper functions

func convertMessages(messages []types.Message) []openai.ChatCompletionMessage {
	result := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		result[i] = openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
			Name:    msg.Name,
		}

		if msg.FunctionCall != nil {
			result[i].FunctionCall = &openai.FunctionCall{
				Name:      msg.FunctionCall.Name,
				Arguments: msg.FunctionCall.Arguments,
			}
		}
	}
	return result
}

func convertFunctions(functions []types.Function) []openai.FunctionDefinition {
	result := make([]openai.FunctionDefinition, len(functions))
	for i, fn := range functions {
		result[i] = openai.FunctionDefinition{
			Name:        fn.Name,
			Description: fn.Description,
			Parameters:  fn.Parameters,
		}
	}
	return result
}

func buildVisionContent(prompt string, imageURLs []string, imageData []string) []openai.ChatMessagePart {
	parts := []openai.ChatMessagePart{
		{
			Type: openai.ChatMessagePartTypeText,
			Text: prompt,
		},
	}

	// Default detail level
	detailLevel := string(openai.ImageURLDetailAuto)

	// Add images
	for _, url := range imageURLs {
		part := openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeImageURL,
			ImageURL: &openai.ChatMessageImageURL{
				Detail: openai.ImageURLDetail(detailLevel),
			},
		}

		// Use URL if provided, otherwise use base64
		if url != "" {
			part.ImageURL.URL = url
		
			
			 {
				
			}
			
		}

		parts = append(parts, part)
	}

	return parts
}

func (c *Client) wrapError(err error) error {
	// Wrap OpenAI errors with provider context
	errMsg := err.Error()

	// Check for specific error types
	if contains(errMsg, "rate limit") {
		return aierrors.NewProviderError("openai", aierrors.ErrRateLimitExceeded, errMsg, true)
	}
	if contains(errMsg, "invalid_api_key") || contains(errMsg, "authentication") {
		return aierrors.NewProviderError("openai", aierrors.ErrInvalidAPIKey, errMsg, false)
	}
	if contains(errMsg, "context_length_exceeded") {
		return aierrors.NewProviderError("openai", aierrors.ErrContextLengthExceeded, errMsg, false)
	}
	if contains(errMsg, "timeout") {
		return aierrors.NewProviderError("openai", aierrors.ErrTimeout, errMsg, true)
	}
	if contains(errMsg, "model_not_found") || contains(errMsg, "invalid_model") {
		return aierrors.NewProviderError("openai", aierrors.ErrModelNotSupported, errMsg, false)
	}

	// Default: provider unavailable (retryable)
	return aierrors.NewProviderError("openai", err, errMsg, true)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		len(s) > len(substr)*2 && s[len(s)/2-len(substr)/2:len(s)/2+len(substr)/2+len(substr)%2] == substr))
}












