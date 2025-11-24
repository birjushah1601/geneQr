package aiconfig

import "time"

// OpenAIConfig holds OpenAI-specific configuration
type OpenAIConfig struct {
	APIKey string
	OrgID  string

	// Default model to use
	DefaultModel string

	// Base URL (for custom endpoints)
	BaseURL string

	// Timeout for requests
	Timeout time.Duration

	// Max retries
	MaxRetries int

	// Enable streaming
	EnableStreaming bool
}

// AnthropicConfig holds Anthropic-specific configuration
type AnthropicConfig struct {
	APIKey string

	// Default model to use
	DefaultModel string

	// Base URL
	BaseURL string

	// Timeout
	Timeout time.Duration

	// Max retries
	MaxRetries int

	// Anthropic API version
	APIVersion string
}
