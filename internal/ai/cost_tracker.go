package ai

import (
	"sync"
	"time"
)

// CostTracker tracks AI usage and costs
type CostTracker struct {
	mu sync.RWMutex

	// Total usage by provider
	usageByProvider map[string]*ProviderUsage

	// Daily usage
	dailyUsage map[string]map[string]*DailyUsage // date -> provider -> usage
}

// ProviderUsage tracks usage for a provider
type ProviderUsage struct {
	Provider         string
	TotalRequests    int64
	TotalTokens      int64
	PromptTokens     int64
	CompletionTokens int64
	TotalCostUSD     float64
	LastUpdated      time.Time
}

// DailyUsage tracks daily usage
type DailyUsage struct {
	Date             string
	Provider         string
	Requests         int64
	Tokens           int64
	PromptTokens     int64
	CompletionTokens int64
	CostUSD          float64
}

// NewCostTracker creates a new cost tracker
func NewCostTracker() *CostTracker {
	return &CostTracker{
		usageByProvider: make(map[string]*ProviderUsage),
		dailyUsage:      make(map[string]map[string]*DailyUsage),
	}
}

// Track tracks a request
func (ct *CostTracker) Track(provider string, tokens TokenUsage, costUSD float64) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	// Update provider usage
	if _, exists := ct.usageByProvider[provider]; !exists {
		ct.usageByProvider[provider] = &ProviderUsage{
			Provider: provider,
		}
	}

	usage := ct.usageByProvider[provider]
	usage.TotalRequests++
	usage.TotalTokens += int64(tokens.TotalTokens)
	usage.PromptTokens += int64(tokens.PromptTokens)
	usage.CompletionTokens += int64(tokens.CompletionTokens)
	usage.TotalCostUSD += costUSD
	usage.LastUpdated = time.Now()

	// Update daily usage
	date := time.Now().Format("2006-01-02")
	if _, exists := ct.dailyUsage[date]; !exists {
		ct.dailyUsage[date] = make(map[string]*DailyUsage)
	}

	if _, exists := ct.dailyUsage[date][provider]; !exists {
		ct.dailyUsage[date][provider] = &DailyUsage{
			Date:     date,
			Provider: provider,
		}
	}

	dailyUsage := ct.dailyUsage[date][provider]
	dailyUsage.Requests++
	dailyUsage.Tokens += int64(tokens.TotalTokens)
	dailyUsage.PromptTokens += int64(tokens.PromptTokens)
	dailyUsage.CompletionTokens += int64(tokens.CompletionTokens)
	dailyUsage.CostUSD += costUSD
}

// GetProviderUsage returns usage for a provider
func (ct *CostTracker) GetProviderUsage(provider string) *ProviderUsage {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	if usage, exists := ct.usageByProvider[provider]; exists {
		// Return a copy
		usageCopy := *usage
		return &usageCopy
	}

	return &ProviderUsage{Provider: provider}
}

// GetAllProviderUsage returns usage for all providers
func (ct *CostTracker) GetAllProviderUsage() []*ProviderUsage {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	result := make([]*ProviderUsage, 0, len(ct.usageByProvider))
	for _, usage := range ct.usageByProvider {
		usageCopy := *usage
		result = append(result, &usageCopy)
	}

	return result
}

// GetDailyUsage returns daily usage for a specific date
func (ct *CostTracker) GetDailyUsage(date string) []*DailyUsage {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	if daily, exists := ct.dailyUsage[date]; exists {
		result := make([]*DailyUsage, 0, len(daily))
		for _, usage := range daily {
			usageCopy := *usage
			result = append(result, &usageCopy)
		}
		return result
	}

	return []*DailyUsage{}
}

// GetTodayUsage returns today's usage
func (ct *CostTracker) GetTodayUsage() []*DailyUsage {
	today := time.Now().Format("2006-01-02")
	return ct.GetDailyUsage(today)
}

// CalculateCost calculates the cost based on token usage and model pricing
func CalculateCost(tokens TokenUsage, inputCostPer1M, outputCostPer1M float64) float64 {
	inputCost := (float64(tokens.PromptTokens) / 1_000_000) * inputCostPer1M
	outputCost := (float64(tokens.CompletionTokens) / 1_000_000) * outputCostPer1M
	return inputCost + outputCost
}

// ModelPricing holds pricing information for different models
var ModelPricing = map[string]struct {
	InputPer1M  float64
	OutputPer1M float64
}{
	// OpenAI models
	"gpt-4o": {
		InputPer1M:  5.00,
		OutputPer1M: 15.00,
	},
	"gpt-4-turbo": {
		InputPer1M:  10.00,
		OutputPer1M: 30.00,
	},
	"gpt-4": {
		InputPer1M:  30.00,
		OutputPer1M: 60.00,
	},
	"gpt-3.5-turbo": {
		InputPer1M:  0.50,
		OutputPer1M: 1.50,
	},

	// Anthropic models
	"claude-3-5-sonnet-20241022": {
		InputPer1M:  3.00,
		OutputPer1M: 15.00,
	},
	"claude-3-opus-20240229": {
		InputPer1M:  15.00,
		OutputPer1M: 75.00,
	},
	"claude-3-sonnet-20240229": {
		InputPer1M:  3.00,
		OutputPer1M: 15.00,
	},
	"claude-3-haiku-20240307": {
		InputPer1M:  0.25,
		OutputPer1M: 1.25,
	},
}

// GetModelPricing returns pricing for a model
func GetModelPricing(model string) (inputPer1M, outputPer1M float64, found bool) {
	if pricing, exists := ModelPricing[model]; exists {
		return pricing.InputPer1M, pricing.OutputPer1M, true
	}
	return 0, 0, false
}

