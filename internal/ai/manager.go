package ai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aby-med/medical-platform/internal/ai/anthropic"
	"github.com/aby-med/medical-platform/internal/ai/openai"
	aitypes "github.com/aby-med/medical-platform/internal/ai/types"
	aierrors "github.com/aby-med/medical-platform/pkg/ai"
)

// Manager manages multiple AI providers with fallback logic
type Manager struct {
	config      *Config
	providers   map[string]aitypes.Provider
	costTracker *CostTracker

	// Health tracking
	healthMu     sync.RWMutex
	providerHealth map[string]*aitypes.ProviderHealth

	// Retry configuration
	maxRetries         int
	retryBackoff       time.Duration
	retryBackoffMult   float64
}

// NewManager creates a new AI provider manager
func NewManager(config *Config) (*Manager, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	manager := &Manager{
		config:           config,
		providers:        make(map[string]aitypes.Provider),
		costTracker:      NewCostTracker(),
		providerHealth:   make(map[string]*aitypes.ProviderHealth),
		maxRetries:       config.MaxRetries,
		retryBackoff:     1 * time.Second,
		retryBackoffMult: config.RetryBackoffMultiplier,
	}

	// Initialize OpenAI provider
	if config.OpenAI.APIKey != "" {
		openaiClient, err := openai.NewClient(config.OpenAI)
		if err != nil {
			return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
		}
		manager.providers["openai"] = openaiClient
		manager.providerHealth["openai"] = &aitypes.ProviderHealth{
			IsHealthy:       true,
			LastHealthCheck: time.Now(),
		}
	}

	// Initialize Anthropic provider
	if config.Anthropic.APIKey != "" && config.EnableFallback {
		anthropicClient, err := anthropic.NewClient(config.Anthropic)
		if err != nil {
			return nil, fmt.Errorf("failed to create Anthropic client: %w", err)
		}
		manager.providers["anthropic"] = anthropicClient
		manager.providerHealth["anthropic"] = &aitypes.ProviderHealth{
			IsHealthy:       true,
			LastHealthCheck: time.Now(),
		}
	}

	if len(manager.providers) == 0 {
		return nil, fmt.Errorf("no providers configured")
	}

	// Start health check routine if enabled
	if config.EnableHealthChecks {
		go manager.healthCheckLoop(context.Background())
	}

	return manager, nil
}

// Chat sends a chat completion request with automatic fallback
func (m *Manager) Chat(ctx context.Context, req *aitypes.ChatRequest) (*aitypes.ChatResponse, error) {
	providers := m.getProviderOrder()
	
	var lastErr error
	for _, providerName := range providers {
		provider, exists := m.providers[providerName]
		if !exists {
			continue
		}

		// Check if provider is healthy
		if !m.isProviderHealthy(providerName) {
			lastErr = fmt.Errorf("provider %s is unhealthy", providerName)
			continue
		}

		// Try with retries
		resp, err := m.chatWithRetry(ctx, provider, req)
		if err != nil {
			lastErr = err
			m.recordError(providerName)
			
			// If error is not retryable and we have fallback, try next provider
			if !aierrors.IsRetryable(err) && m.config.EnableFallback {
				continue
			}
			
			// If retryable error exhausted retries, try fallback
			if m.config.EnableFallback {
				continue
			}
			
			return nil, err
		}

		// Track cost if enabled
		if m.config.EnableCostTracking {
			m.costTracker.Track(providerName, resp.Usage, resp.Cost)
		}

		m.recordSuccess(providerName, resp.Latency)
		return resp, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("all providers failed: %w", lastErr)
	}

	return nil, aierrors.ErrNoProvidersAvailable
}

// ChatStream sends a streaming chat completion request with automatic fallback
func (m *Manager) ChatStream(ctx context.Context, req *aitypes.ChatRequest, streamChan chan<- *ChatStreamResponse) error {
	providers := m.getProviderOrder()
	
	var lastErr error
	for _, providerName := range providers {
		provider, exists := m.providers[providerName]
		if !exists {
			continue
		}

		// Check if provider is healthy
		if !m.isProviderHealthy(providerName) {
			lastErr = fmt.Errorf("provider %s is unhealthy", providerName)
			continue
		}

		// Try streaming
		err := provider.ChatStream(ctx, req, streamChan)
		if err != nil {
			lastErr = err
			m.recordError(providerName)
			
			// Try fallback on error
			if m.config.EnableFallback {
				continue
			}
			
			return err
		}

		return nil
	}

	if lastErr != nil {
		return fmt.Errorf("all providers failed: %w", lastErr)
	}

	return aierrors.ErrNoProvidersAvailable
}

// Analyze sends a vision analysis request with automatic fallback
func (m *Manager) Analyze(ctx context.Context, req *VisionRequest) (*VisionResponse, error) {
	providers := m.getProviderOrder()
	
	var lastErr error
	for _, providerName := range providers {
		provider, exists := m.providers[providerName]
		if !exists {
			continue
		}

		// Check if provider supports vision
		caps := provider.GetCapabilities()
		if !caps.SupportsVision {
			continue
		}

		// Check if provider is healthy
		if !m.isProviderHealthy(providerName) {
			lastErr = fmt.Errorf("provider %s is unhealthy", providerName)
			continue
		}

		// Try with retries
		resp, err := m.analyzeWithRetry(ctx, provider, req)
		if err != nil {
			lastErr = err
			m.recordError(providerName)
			
			// Try fallback
			if m.config.EnableFallback {
				continue
			}
			
			return nil, err
		}

		// Track cost
		if m.config.EnableCostTracking {
			m.costTracker.Track(providerName, resp.Usage, resp.Cost)
		}

		m.recordSuccess(providerName, resp.Latency)
		return resp, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("all providers failed: %w", lastErr)
	}

	return nil, aierrors.ErrNoProvidersAvailable
}

// GetCostTracker returns the cost tracker
func (m *Manager) GetCostTracker() *CostTracker {
	return m.costTracker
}

// GetProviderHealth returns health status for a provider
func (m *Manager) GetProviderHealth(providerName string) *aitypes.ProviderHealth {
	m.healthMu.RLock()
	defer m.healthMu.RUnlock()

	if health, exists := m.providerHealth[providerName]; exists {
		healthCopy := *health
		return &healthCopy
	}

	return nil
}

// GetAllProviderHealth returns health status for all providers
func (m *Manager) GetAllProviderHealth() map[string]*aitypes.ProviderHealth {
	m.healthMu.RLock()
	defer m.healthMu.RUnlock()

	result := make(map[string]*aitypes.ProviderHealth, len(m.providerHealth))
	for name, health := range m.providerHealth {
		healthCopy := *health
		result[name] = &healthCopy
	}

	return result
}

// Close closes all provider connections
func (m *Manager) Close() error {
	for _, provider := range m.providers {
		if err := provider.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Internal helper methods

func (m *Manager) getProviderOrder() []string {
	// Start with default provider
	order := []string{m.config.DefaultProvider}

	// Add fallback providers
	if m.config.EnableFallback {
		for name := range m.providers {
			if name != m.config.DefaultProvider {
				order = append(order, name)
			}
		}
	}

	return order
}

func (m *Manager) isProviderHealthy(providerName string) bool {
	m.healthMu.RLock()
	defer m.healthMu.RUnlock()

	if health, exists := m.providerHealth[providerName]; exists {
		return health.IsHealthy
	}

	return false
}

func (m *Manager) recordError(providerName string) {
	m.healthMu.Lock()
	defer m.healthMu.Unlock()

	if health, exists := m.providerHealth[providerName]; exists {
		health.HealthCheckFailures++
		
		// Mark as unhealthy after 3 consecutive failures
		if health.HealthCheckFailures >= 3 {
			health.IsHealthy = false
		}
	}
}

func (m *Manager) recordSuccess(providerName string, latency time.Duration) {
	m.healthMu.Lock()
	defer m.healthMu.Unlock()

	if health, exists := m.providerHealth[providerName]; exists {
		health.IsHealthy = true
		health.HealthCheckFailures = 0
		health.LastHealthCheck = time.Now()
		
		// Update average latency (simple moving average)
		if health.AverageLatency == 0 {
			health.AverageLatency = latency
		} else {
			health.AverageLatency = (health.AverageLatency + latency) / 2
		}
		
		health.RequestsLast24h++
	}
}

func (m *Manager) chatWithRetry(ctx context.Context, provider aitypes.Provider, req *aitypes.ChatRequest) (*aitypes.ChatResponse, error) {
	var lastErr error
	backoff := m.retryBackoff

	for attempt := 0; attempt <= m.maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
			
			backoff = time.Duration(float64(backoff) * m.retryBackoffMult)
		}

		resp, err := provider.Chat(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry non-retryable errors
		if !aierrors.IsRetryable(err) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func (m *Manager) analyzeWithRetry(ctx context.Context, provider aitypes.Provider, req *VisionRequest) (*VisionResponse, error) {
	var lastErr error
	backoff := m.retryBackoff

	for attempt := 0; attempt <= m.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
			
			backoff = time.Duration(float64(backoff) * m.retryBackoffMult)
		}

		resp, err := provider.Analyze(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		if !aierrors.IsRetryable(err) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func (m *Manager) healthCheckLoop(ctx context.Context) {
	ticker := time.NewTicker(m.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.performHealthChecks(ctx)
		}
	}
}

func (m *Manager) performHealthChecks(ctx context.Context) {
	for name, provider := range m.providers {
		go func(providerName string, prov aitypes.Provider) {
			healthy := prov.IsHealthy(ctx)

			m.healthMu.Lock()
			defer m.healthMu.Unlock()

			if health, exists := m.providerHealth[providerName]; exists {
				health.LastHealthCheck = time.Now()
				
				if healthy {
					health.IsHealthy = true
					health.HealthCheckFailures = 0
				} else {
					health.HealthCheckFailures++
					if health.HealthCheckFailures >= 3 {
						health.IsHealthy = false
					}
				}
			}
		}(name, provider)
	}
}




