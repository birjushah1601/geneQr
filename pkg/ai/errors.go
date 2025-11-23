// Package ai provides AI-specific error types
package ai

import (
	"errors"
	"fmt"
)

var (
	// ErrProviderUnavailable indicates the provider is not available
	ErrProviderUnavailable = errors.New("ai provider unavailable")

	// ErrRateLimitExceeded indicates rate limit has been exceeded
	ErrRateLimitExceeded = errors.New("ai provider rate limit exceeded")

	// ErrInvalidAPIKey indicates the API key is invalid
	ErrInvalidAPIKey = errors.New("invalid API key")

	// ErrModelNotSupported indicates the model is not supported
	ErrModelNotSupported = errors.New("model not supported")

	// ErrInvalidRequest indicates the request is invalid
	ErrInvalidRequest = errors.New("invalid request")

	// ErrContextLengthExceeded indicates context length exceeded
	ErrContextLengthExceeded = errors.New("context length exceeded")

	// ErrTimeout indicates request timed out
	ErrTimeout = errors.New("request timeout")

	// ErrNoProvidersAvailable indicates no providers are available
	ErrNoProvidersAvailable = errors.New("no providers available")

	// ErrAllProvidersFailed indicates all providers failed
	ErrAllProvidersFailed = errors.New("all providers failed")
)

// ProviderError wraps an error with provider context
type ProviderError struct {
	Provider string
	Err      error
	Message  string
	Retryable bool
}

func (e *ProviderError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("[%s] %s: %v", e.Provider, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %v", e.Provider, e.Err)
}

func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError creates a new provider error
func NewProviderError(provider string, err error, message string, retryable bool) *ProviderError {
	return &ProviderError{
		Provider:  provider,
		Err:       err,
		Message:   message,
		Retryable: retryable,
	}
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	var provErr *ProviderError
	if errors.As(err, &provErr) {
		return provErr.Retryable
	}

	// Default retryable errors
	return errors.Is(err, ErrProviderUnavailable) ||
		errors.Is(err, ErrRateLimitExceeded) ||
		errors.Is(err, ErrTimeout)
}
