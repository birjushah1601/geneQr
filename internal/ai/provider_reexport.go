// Package ai re-exports types from internal/ai/types for backward compatibility
package ai

import (
	"github.com/aby-med/medical-platform/internal/ai/types"
)

// Re-export types for backward compatibility
type (
	Provider               = types.Provider
	ChatRequest            = types.ChatRequest
	ChatResponse           = types.ChatResponse
	ChatStreamResponse     = types.ChatStreamResponse
	Message                = types.Message
	MessageRole            = types.MessageRole
	Function               = types.Function
	FunctionCall           = types.FunctionCall
	TokenUsage             = types.TokenUsage
	VisionRequest          = types.VisionRequest
	VisionResponse         = types.VisionResponse
	ProviderCapabilities   = types.ProviderCapabilities
	ProviderHealth         = types.ProviderHealth
)

// Re-export constants
const (
	RoleSystem    = types.RoleSystem
	RoleUser      = types.RoleUser
	RoleAssistant = types.RoleAssistant
	RoleFunction  = types.RoleFunction
	RoleTool      = types.RoleTool
)

// Note: GetModelPricing and CalculateCost are in cost_tracker.go, not re-exported here
