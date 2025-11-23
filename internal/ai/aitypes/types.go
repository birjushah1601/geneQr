package aitypes

// ChatMessage represents a message in a conversation
type ChatMessage struct {
	Role    string
	Content string
}

// ChatRequest represents a request for chat completion
type ChatRequest struct {
	Messages    []ChatMessage
	Temperature float64
	MaxTokens   int
	Model       string
}

// ChatResponse represents a response from chat completion
type ChatResponse struct {
	Content      string
	FinishReason string
	TokensUsed   int
	Model        string
}

// CompletionRequest represents a request for text completion
type CompletionRequest struct {
	Prompt      string
	Temperature float64
	MaxTokens   int
	Model       string
}

// CompletionResponse represents a response from text completion
type CompletionResponse struct {
	Text         string
	FinishReason string
	TokensUsed   int
	Model        string
}

// ImageAnalysisRequest represents a request to analyze an image
type ImageAnalysisRequest struct {
	ImageURL    string
	ImageBase64 string
	Prompt      string
	MaxTokens   int
}

// ImageAnalysisResponse represents the analysis result
type ImageAnalysisResponse struct {
	Description string
	TokensUsed  int
}
