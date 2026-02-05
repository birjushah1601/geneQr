package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aby-med/medical-platform/internal/ai/types"
	"github.com/aby-med/medical-platform/internal/ai/aiconfig"
	aierrors "github.com/aby-med/medical-platform/pkg/ai"
)

// Client implements the AI Provider interface for Anthropic
type Client struct {
	httpClient   *http.Client
	config       aiconfig.AnthropicConfig
	providerCode string
}

// NewClient creates a new Anthropic client
func NewClient(config aiconfig.AnthropicConfig) (*Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Anthropic API key is required")
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		config:       config,
		providerCode: "anthropic-claude",
	}, nil
}

// GetName returns the provider name
func (c *Client) GetName() string {
	return "anthropic"
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

	// Build Anthropic request
	anthropicReq := &anthropicChatRequest{
		Model:      model,
		Messages:   convertMessages(req.Messages),
		MaxTokens:  1024, // Default
		Stream:     false,
	}

	// Set optional parameters
	if req.Temperature != nil {
		anthropicReq.Temperature = req.Temperature
	}
	if req.TopP != nil {
		anthropicReq.TopP = req.TopP
	}
	if req.MaxTokens != nil {
		anthropicReq.MaxTokens = *req.MaxTokens
	}
	if len(req.StopSequences) > 0 {
		anthropicReq.StopSequences = req.StopSequences
	}

	// Extract system message if present
	if len(req.Messages) > 0 && req.Messages[0].Role == types.RoleSystem {
		anthropicReq.System = req.Messages[0].Content
		anthropicReq.Messages = anthropicReq.Messages[1:] // Remove system message from messages
	}

	// Call Anthropic API
	respData, err := c.makeRequest(ctx, "/v1/messages", anthropicReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	var anthropicResp anthropicChatResponse
	if err := json.Unmarshal(respData, &anthropicResp); err != nil {
		return nil, aierrors.NewProviderError("anthropic", err, "failed to parse response", false)
	}

	latency := time.Since(startTime)

	// Handle empty response
	if len(anthropicResp.Content) == 0 {
		return nil, aierrors.NewProviderError("anthropic", fmt.Errorf("no content returned"), "", false)
	}

	// Extract text content
	var content string
	for _, block := range anthropicResp.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	// Build response
	response := &types.ChatResponse{
		ID:      anthropicResp.ID,
		Content: content,
		Usage: types.TokenUsage{
			PromptTokens:     anthropicResp.Usage.InputTokens,
			CompletionTokens: anthropicResp.Usage.OutputTokens,
			TotalTokens:      anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
		},
		Provider:     "anthropic",
		Model:        anthropicResp.Model,
		Latency:      latency,
		FinishReason: anthropicResp.StopReason,
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

	// Build Anthropic request
	anthropicReq := &anthropicChatRequest{
		Model:     model,
		Messages:  convertMessages(req.Messages),
		MaxTokens: 1024,
		Stream:    true,
	}

	// Set optional parameters
	if req.Temperature != nil {
		anthropicReq.Temperature = req.Temperature
	}
	if req.TopP != nil {
		anthropicReq.TopP = req.TopP
	}
	if req.MaxTokens != nil {
		anthropicReq.MaxTokens = *req.MaxTokens
	}

	// Extract system message
	if len(req.Messages) > 0 && req.Messages[0].Role == types.RoleSystem {
		anthropicReq.System = req.Messages[0].Content
		anthropicReq.Messages = anthropicReq.Messages[1:]
	}

	// Make streaming request
	reqBody, err := json.Marshal(anthropicReq)
	if err != nil {
		streamChan <- &types.ChatStreamResponse{Done: true, Error: err}
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/v1/messages", bytes.NewBuffer(reqBody))
	if err != nil {
		streamChan <- &types.ChatStreamResponse{Done: true, Error: err}
		return err
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		wrappedErr := c.wrapError(err)
		streamChan <- &types.ChatStreamResponse{Done: true, Error: wrappedErr}
		return wrappedErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("API error: %s - %s", resp.Status, string(body))
		wrappedErr := c.wrapError(err)
		streamChan <- &types.ChatStreamResponse{Done: true, Error: wrappedErr}
		return wrappedErr
	}

	// Read stream line by line
	decoder := json.NewDecoder(resp.Body)
	for {
		var event anthropicStreamEvent
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				streamChan <- &types.ChatStreamResponse{Done: true}
				return nil
			}
			wrappedErr := c.wrapError(err)
			streamChan <- &types.ChatStreamResponse{Done: true, Error: wrappedErr}
			return wrappedErr
		}

		// Handle different event types
		switch event.Type {
		case "content_block_delta":
			if event.Delta.Type == "text_delta" {
				streamChan <- &types.ChatStreamResponse{
					Delta: event.Delta.Text,
					Done:  false,
				}
			}
		case "message_stop":
			streamChan <- &types.ChatStreamResponse{
				Done:         true,
				FinishReason: "stop",
			}
			return nil
		case "error":
			err := fmt.Errorf("stream error: %s", event.Error.Message)
			wrappedErr := c.wrapError(err)
			streamChan <- &types.ChatStreamResponse{Done: true, Error: wrappedErr}
			return wrappedErr
		}
	}
}

// Analyze analyzes an image/video with vision capabilities
func (c *Client) Analyze(ctx context.Context, req *types.VisionRequest) (*types.VisionResponse, error) {
	startTime := time.Now()

	// Set default model if not specified
	model := req.Model
	if model == "" {
		model = c.config.DefaultModel
	}

	// Build vision content
	content := buildVisionContent(req.Prompt, req.ImageURLs, req.ImageData)

	// Build Anthropic request
	anthropicReq := &anthropicChatRequest{
		Model: model,
		Messages: []anthropicMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
		MaxTokens: 1024,
	}

	if req.Temperature != nil {
		anthropicReq.Temperature = req.Temperature
	}
	if req.MaxTokens != nil {
		anthropicReq.MaxTokens = *req.MaxTokens
	}

	// Call Anthropic API
	respData, err := c.makeRequest(ctx, "/v1/messages", anthropicReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	var anthropicResp anthropicChatResponse
	if err := json.Unmarshal(respData, &anthropicResp); err != nil {
		return nil, aierrors.NewProviderError("anthropic", err, "failed to parse response", false)
	}

	latency := time.Since(startTime)

	// Extract text content
	var analysis string
	for _, block := range anthropicResp.Content {
		if block.Type == "text" {
			analysis += block.Text
		}
	}

	// Build response
	response := &types.VisionResponse{
		Analysis: analysis,
		Usage: types.TokenUsage{
			PromptTokens:     anthropicResp.Usage.InputTokens,
			CompletionTokens: anthropicResp.Usage.OutputTokens,
			TotalTokens:      anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
		},
		Provider: "anthropic",
		Model:    anthropicResp.Model,
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
	// Simple health check: try a minimal request
	req := &types.ChatRequest{
		Model: c.config.DefaultModel,
		Messages: []types.Message{
			{Role: types.RoleUser, Content: "test"},
		},
	}

	_, err := c.Chat(ctx, req)
	return err == nil
}

// GetCapabilities returns provider capabilities
func (c *Client) GetCapabilities() *types.ProviderCapabilities {
	return &types.ProviderCapabilities{
		SupportsChat:             true,
		SupportsStreaming:        true,
		SupportsVision:           true,
		SupportsFunctionCalling:  false, // Claude doesn't support function calling directly
		SupportsEmbeddings:       false,
		MaxContextTokens:         200000, // Claude 3.5 Sonnet
		MaxOutputTokens:          8192,
		RateLimitRequestsPerMin:  4000,
		RateLimitTokensPerMin:    400000,
	}
}

// Close closes any open connections
func (c *Client) Close() error {
	// HTTP client doesn't require explicit closing
	return nil
}

// Helper functions

func (c *Client) makeRequest(ctx context.Context, endpoint string, body interface{}) ([]byte, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, aierrors.NewProviderError("anthropic", err, "failed to marshal request", false)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, aierrors.NewProviderError("anthropic", err, "failed to create request", false)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, c.wrapError(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, aierrors.NewProviderError("anthropic", err, "failed to read response", false)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleAPIError(resp.StatusCode, respBody)
	}

	return respBody, nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("anthropic-version", c.config.APIVersion)
}

func (c *Client) handleAPIError(statusCode int, body []byte) error {
	var errResp anthropicErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return aierrors.NewProviderError("anthropic", fmt.Errorf("API error: %d - %s", statusCode, string(body)), "", false)
	}

	errMsg := fmt.Sprintf("%s: %s", errResp.Error.Type, errResp.Error.Message)

	// Check for specific error types
	switch errResp.Error.Type {
	case "rate_limit_error":
		return aierrors.NewProviderError("anthropic", aierrors.ErrRateLimitExceeded, errMsg, true)
	case "authentication_error", "permission_error":
		return aierrors.NewProviderError("anthropic", aierrors.ErrInvalidAPIKey, errMsg, false)
	case "invalid_request_error":
		return aierrors.NewProviderError("anthropic", aierrors.ErrInvalidRequest, errMsg, false)
	default:
		return aierrors.NewProviderError("anthropic", fmt.Errorf("API error"), errMsg, true)
	}
}

func (c *Client) wrapError(err error) error {
	errMsg := err.Error()

	if strings.Contains(errMsg, "timeout") {
		return aierrors.NewProviderError("anthropic", aierrors.ErrTimeout, errMsg, true)
	}

	return aierrors.NewProviderError("anthropic", err, errMsg, true)
}

func convertMessages(messages []types.Message) []anthropicMessage {
	result := make([]anthropicMessage, 0, len(messages))
	
	for _, msg := range messages {
		// Skip system messages (handled separately)
		if msg.Role == types.RoleSystem {
			continue
		}

		result = append(result, anthropicMessage{
			Role: string(msg.Role),
			Content: []anthropicContent{
				{
					Type: "text",
					Text: msg.Content,
				},
			},
		})
	}
	
	return result
}

func buildVisionContent(prompt string, imageURLs []string, imageData []string) []anthropicContent {
	content := make([]anthropicContent, 0, len(imageData)+1)

	// Add images first
	for _, data := range imageData {
		if data != "" {
			content = append(content, anthropicContent{
				Type: "image",
				Source: &anthropicImageSource{
					Type:      "base64",
					MediaType: "image/jpeg",
					Data:      data,
				},
			})
		}
	}


	// Add prompt
	content = append(content, anthropicContent{
		Type: "text",
		Text: prompt,
	})

	return content
}

// Anthropic API types

type anthropicChatRequest struct {
	Model         string              `json:"model"`
	Messages      []anthropicMessage  `json:"messages"`
	MaxTokens     int                 `json:"max_tokens"`
	Temperature   *float32            `json:"temperature,omitempty"`
	TopP          *float32            `json:"top_p,omitempty"`
	StopSequences []string            `json:"stop_sequences,omitempty"`
	Stream        bool                `json:"stream,omitempty"`
	System        string              `json:"system,omitempty"`
}

type anthropicMessage struct {
	Role    string             `json:"role"`
	Content []anthropicContent `json:"content"`
}

type anthropicContent struct {
	Type   string                `json:"type"`
	Text   string                `json:"text,omitempty"`
	Source *anthropicImageSource `json:"source,omitempty"`
}

type anthropicImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type anthropicChatResponse struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Role       string `json:"role"`
	Content    []anthropicContentBlock `json:"content"`
	Model      string `json:"model"`
	StopReason string `json:"stop_reason"`
	Usage      anthropicUsage `json:"usage"`
}

type anthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type anthropicStreamEvent struct {
	Type  string                  `json:"type"`
	Delta anthropicStreamDelta    `json:"delta,omitempty"`
	Error anthropicStreamError    `json:"error,omitempty"`
}

type anthropicStreamDelta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicStreamError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type anthropicErrorResponse struct {
	Type  string               `json:"type"`
	Error anthropicErrorDetail `json:"error"`
}

type anthropicErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}









