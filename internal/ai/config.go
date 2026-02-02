package ai

import (
	"fmt"
	"os"
	"strconv"
	"time"
	
	"github.com/aby-med/medical-platform/internal/ai/aiconfig"
)

// Config holds AI provider configuration
type Config struct {
	// OpenAI configuration
	OpenAI aiconfig.OpenAIConfig

	// Anthropic configuration
	Anthropic aiconfig.AnthropicConfig

	// Default provider to use
	DefaultProvider string

	// Enable fallback to secondary provider
	EnableFallback bool

	// Max retries per provider
	MaxRetries int

	// Retry backoff multiplier
	RetryBackoffMultiplier float64

	// Default timeout for requests
	DefaultTimeout time.Duration

	// Enable cost tracking
	EnableCostTracking bool

	// Enable health checks
	EnableHealthChecks bool

	// Health check interval
	HealthCheckInterval time.Duration
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() (*Config, error) {
	config := &Config{
		DefaultProvider:        getEnvOrDefault("AI_DEFAULT_PROVIDER", "openai"),
		EnableFallback:         getEnvBool("AI_ENABLE_FALLBACK", true),
		MaxRetries:             getEnvInt("AI_MAX_RETRIES", 3),
		RetryBackoffMultiplier: getEnvFloat("AI_RETRY_BACKOFF", 2.0),
		DefaultTimeout:         getEnvDuration("AI_DEFAULT_TIMEOUT", 30*time.Second),
		EnableCostTracking:     getEnvBool("AI_ENABLE_COST_TRACKING", true),
		EnableHealthChecks:     getEnvBool("AI_ENABLE_HEALTH_CHECKS", true),
		HealthCheckInterval:    getEnvDuration("AI_HEALTH_CHECK_INTERVAL", 5*time.Minute),
	}

	// OpenAI configuration
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" && config.DefaultProvider == "openai" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	config.OpenAI = aiconfig.OpenAIConfig{
		APIKey:          openaiAPIKey,
		OrgID:           os.Getenv("OPENAI_ORG_ID"),
		DefaultModel:    getEnvOrDefault("OPENAI_DEFAULT_MODEL", "gpt-4o"),
		BaseURL:         getEnvOrDefault("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		Timeout:         getEnvDuration("OPENAI_TIMEOUT", 30*time.Second),
		MaxRetries:      getEnvInt("OPENAI_MAX_RETRIES", 3),
		EnableStreaming: getEnvBool("OPENAI_ENABLE_STREAMING", true),
	}

	// Anthropic configuration
	anthropicAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicAPIKey == "" && config.EnableFallback {
		// Warning: fallback provider not configured, but not fatal
		fmt.Println("WARNING: ANTHROPIC_API_KEY not set, fallback disabled")
		config.EnableFallback = false
	}

	config.Anthropic = aiconfig.AnthropicConfig{
		APIKey:       anthropicAPIKey,
		DefaultModel: getEnvOrDefault("ANTHROPIC_DEFAULT_MODEL", "claude-3-5-sonnet-20241022"),
		BaseURL:      getEnvOrDefault("ANTHROPIC_BASE_URL", "https://api.anthropic.com"),
		Timeout:      getEnvDuration("ANTHROPIC_TIMEOUT", 30*time.Second),
		MaxRetries:   getEnvInt("ANTHROPIC_MAX_RETRIES", 3),
		APIVersion:   getEnvOrDefault("ANTHROPIC_API_VERSION", "2023-06-01"),
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.DefaultProvider != "openai" && c.DefaultProvider != "anthropic" {
		return fmt.Errorf("invalid default provider: %s (must be 'openai' or 'anthropic')", c.DefaultProvider)
	}

	if c.DefaultProvider == "openai" && c.OpenAI.APIKey == "" {
		return fmt.Errorf("OpenAI API key is required when using OpenAI as default provider")
	}

	if c.EnableFallback && c.Anthropic.APIKey == "" && c.DefaultProvider == "openai" {
		return fmt.Errorf("Anthropic API key is required when fallback is enabled")
	}

	if c.MaxRetries < 0 || c.MaxRetries > 10 {
		return fmt.Errorf("max retries must be between 0 and 10")
	}

	if c.RetryBackoffMultiplier < 1.0 {
		return fmt.Errorf("retry backoff multiplier must be >= 1.0")
	}

	return nil
}

// Helper functions

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}
	return defaultValue
}



