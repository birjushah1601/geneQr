// Package types provides shared AI types for provider abstraction
package types

import (
	"context"
	"time"
)

// Provider represents an AI provider interface that all AI providers must implement
type Provider interface {
	// GetName returns the provider name (e.g., "openai", "anthropic")
	GetName() string

	// GetProviderCode returns the provider code from database
	GetProviderCode() string

	// Chat sends a chat completion request
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// ChatStream sends a streaming chat completion request
	ChatStream(ctx context.Context, req *ChatRequest, streamChan chan<- *ChatStreamResponse) error

	// Analyze analyzes an image/video with vision capabilities
	Analyze(ctx context.Context, req *VisionRequest) (*VisionResponse, error)

	// IsHealthy checks if the provider is healthy and available
	IsHealthy(ctx context.Context) bool

	// GetCapabilities returns provider capabilities
	GetCapabilities() *ProviderCapabilities

	// Close closes any open connections
	Close() error
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	// Model to use (e.g., "gpt-4o", "claude-3-5-sonnet-20241022")
	Model string

	// Messages in the conversation
	Messages []Message

	// Temperature (0-2, controls randomness)
	Temperature *float32

	// TopP (0-1, nucleus sampling)
	TopP *float32

	// MaxTokens maximum tokens to generate
	MaxTokens *int

	// FunctionCalling enables function calling if supported
	FunctionCalling bool

	// Functions available for the model to call
	Functions []Function

	// StopSequences to end generation
	StopSequences []string

	// User identifier for tracking
	UserID string

	// Metadata for tracking
	Metadata map[string]interface{}
}

// Message represents a chat message
type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`

	// FunctionCall for function calling
	FunctionCall *FunctionCall `json:"function_call,omitempty"`

	// Name of the function (for function role)
	Name string `json:"name,omitempty"`
}

// MessageRole represents the role of a message
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleFunction  MessageRole = "function"
	RoleTool      MessageRole = "tool"
)

// Function represents a function available for calling
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// FunctionCall represents a function call made by the AI
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	// ID of the response
	ID string

	// Content of the response
	Content string

	// FunctionCall if model called a function
	FunctionCall *FunctionCall

	// Usage token usage statistics
	Usage TokenUsage

	// Cost estimated cost in USD
	Cost float64

	// Provider that generated the response
	Provider string

	// Model that generated the response
	Model string

	// Latency time taken to generate response
	Latency time.Duration

	// FinishReason why the generation stopped
	FinishReason string
}

// ChatStreamResponse represents a streaming chat response chunk
type ChatStreamResponse struct {
	// Delta content chunk
	Delta string

	// FunctionCallDelta if streaming function call
	FunctionCallDelta *FunctionCall

	// Done indicates if streaming is complete
	Done bool

	// Error if any error occurred
	Error error

	// FinishReason when done
	FinishReason string
}

// TokenUsage represents token usage statistics
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// VisionRequest represents a request to analyze images/video
type VisionRequest struct {
	// Model to use
	Model string

	// Prompt/question about the image
	Prompt string

	// ImageURLs list of image URLs to analyze
	ImageURLs []string

	// ImageData list of base64 encoded images
	ImageData []string

	// Temperature for response generation
	Temperature *float32

	// MaxTokens maximum tokens in response
	MaxTokens *int

	// Metadata for tracking
	Metadata map[string]interface{}
}

// VisionResponse represents the response from vision analysis
type VisionResponse struct {
	// Analysis result (alias for Content)
	Analysis string

	// Content analysis result
	Content string

	// Usage token usage
	Usage TokenUsage

	// Cost estimated cost
	Cost float64

	// Provider that generated response
	Provider string

	// Model used
	Model string

	// Latency response time
	Latency time.Duration
}

// ProviderCapabilities describes what a provider supports
type ProviderCapabilities struct {
	// SupportsChatCompletion whether provider supports chat (alias for SupportsChat)
	SupportsChatCompletion bool

	// SupportsChat whether provider supports chat
	SupportsChat bool

	// SupportsStreaming whether provider supports streaming
	SupportsStreaming bool

	// SupportsVision whether provider supports image analysis
	SupportsVision bool

	// SupportsFunctionCalling whether provider supports function calling
	SupportsFunctionCalling bool

	// SupportsEmbeddings whether provider supports embeddings
	SupportsEmbeddings bool

	// MaxTokens maximum tokens supported
	MaxTokens int

	// MaxContextTokens maximum context window
	MaxContextTokens int

	// MaxOutputTokens maximum output tokens
	MaxOutputTokens int

	// RateLimitRequestsPerMin rate limit on requests per minute
	RateLimitRequestsPerMin int

	// RateLimitTokensPerMin rate limit on tokens per minute
	RateLimitTokensPerMin int

	// Models list of available models
	Models []string
}

// ProviderHealth tracks provider health status
type ProviderHealth struct {
	IsHealthy           bool
	LastHealthCheck     time.Time
	ConsecutiveFailures int
	LastError           error
	HealthCheckFailures int
	AverageLatency      time.Duration
	RequestsLast24h     int
}

// Model pricing information
type ModelPricing struct {
	InputCostPer1K  float64
	OutputCostPer1K float64
}

// GetModelPricing returns pricing for a given model
func GetModelPricing(model string) (inputCost, outputCost float64, found bool) {
	pricing := map[string]ModelPricing{
		// OpenAI Models
		"gpt-4": {
			InputCostPer1K:  0.03,
			OutputCostPer1K: 0.06,
		},
		"gpt-4-turbo": {
			InputCostPer1K:  0.01,
			OutputCostPer1K: 0.03,
		},
		"gpt-4o": {
			InputCostPer1K:  0.005,
			OutputCostPer1K: 0.015,
		},
		"gpt-3.5-turbo": {
			InputCostPer1K:  0.0005,
			OutputCostPer1K: 0.0015,
		},
		// Anthropic Models
		"claude-3-opus-20240229": {
			InputCostPer1K:  0.015,
			OutputCostPer1K: 0.075,
		},
		"claude-3-sonnet-20240229": {
			InputCostPer1K:  0.003,
			OutputCostPer1K: 0.015,
		},
		"claude-3-5-sonnet-20241022": {
			InputCostPer1K:  0.003,
			OutputCostPer1K: 0.015,
		},
		"claude-3-haiku-20240307": {
			InputCostPer1K:  0.00025,
			OutputCostPer1K: 0.00125,
		},
	}

	p, ok := pricing[model]
	if !ok {
		return 0, 0, false
	}
	return p.InputCostPer1K, p.OutputCostPer1K, true
}

// CalculateCost calculates the cost based on token usage and pricing
func CalculateCost(usage TokenUsage, inputCostPer1K, outputCostPer1K float64) float64 {
	inputCost := float64(usage.PromptTokens) / 1000.0 * inputCostPer1K
	outputCost := float64(usage.CompletionTokens) / 1000.0 * outputCostPer1K
	return inputCost + outputCost
}

