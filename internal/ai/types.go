package ai

import "fmt"

// ChatCompletionRequest represents a chat completion request
type ChatCompletionRequest struct {
	Model            string                 `json:"model"`
	Messages         []ChatMessage          `json:"messages"`
	MaxTokens        *int                   `json:"max_tokens,omitempty"`
	Temperature      *float64               `json:"temperature,omitempty"`
	TopP             *float64               `json:"top_p,omitempty"`
	N                *int                   `json:"n,omitempty"`
	Stream           bool                   `json:"stream,omitempty"`
	Stop             []string               `json:"stop,omitempty"`
	PresencePenalty  *float64               `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64               `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int         `json:"logit_bias,omitempty"`
	User             string                 `json:"user,omitempty"`
	Tools            []Tool                 `json:"tools,omitempty"`
	ToolChoice       interface{}            `json:"tool_choice,omitempty"`
	ResponseFormat   *ResponseFormat        `json:"response_format,omitempty"`
}

// ChatMessage represents a message in chat completion
type ChatMessage struct {
	Role         string       `json:"role"`
	Content      interface{}  `json:"content"` // Can be string or array for vision models
	Name         string       `json:"name,omitempty"`
	ToolCalls    []ToolCall   `json:"tool_calls,omitempty"`
	ToolCallID   string       `json:"tool_call_id,omitempty"`
}

// MessageContent represents content for vision models
type MessageContent struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL represents an image URL for vision models
type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// Tool represents a function tool
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents a function definition
type Function struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"`
}

// ToolCall represents a tool call in response
type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ResponseFormat specifies the format of the response
type ResponseFormat struct {
	Type string `json:"type"`
}

// ChatCompletionResponse represents the response from chat completion
type ChatCompletionResponse struct {
	ID                string                        `json:"id"`
	Object            string                        `json:"object"`
	Created           int64                         `json:"created"`
	Model             string                        `json:"model"`
	SystemFingerprint string                        `json:"system_fingerprint,omitempty"`
	Choices           []ChatCompletionChoice        `json:"choices"`
	Usage             Usage                         `json:"usage"`
}

// ChatCompletionChoice represents a choice in chat completion response
type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// EmbeddingAPIRequest represents an embedding request
type EmbeddingAPIRequest struct {
	Model          string      `json:"model"`
	Input          interface{} `json:"input"`
	EncodingFormat string      `json:"encoding_format,omitempty"`
	Dimensions     *int        `json:"dimensions,omitempty"`
	User           string      `json:"user,omitempty"`
}

// EmbeddingAPIResponse represents the response from embedding
type EmbeddingAPIResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

// Embedding represents a single embedding
type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// Model represents an AI model
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ModelsResponse represents the response from models endpoint
type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// APIError represents an error from the API
type APIError struct {
	ErrorInfo struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param,omitempty"`
		Code    string `json:"code,omitempty"`
	} `json:"error"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s (type: %s)", e.ErrorInfo.Message, e.ErrorInfo.Type)
}

// StreamResponse represents a streaming response chunk
type StreamResponse struct {
	ID                string                     `json:"id"`
	Object            string                     `json:"object"`
	Created           int64                      `json:"created"`
	Model             string                     `json:"model"`
	SystemFingerprint string                     `json:"system_fingerprint,omitempty"`
	Choices           []StreamChoice             `json:"choices"`
}

// StreamChoice represents a choice in streaming response
type StreamChoice struct {
	Index        int              `json:"index"`
	Delta        ChatMessage      `json:"delta"`
	FinishReason *string          `json:"finish_reason"`
}

// CompletionRequest represents a text completion request (legacy)
type CompletionRequest struct {
	Model            string             `json:"model"`
	Prompt           interface{}        `json:"prompt"`
	MaxTokens        *int               `json:"max_tokens,omitempty"`
	Temperature      *float64           `json:"temperature,omitempty"`
	TopP             *float64           `json:"top_p,omitempty"`
	N                *int               `json:"n,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	Logprobs         *int               `json:"logprobs,omitempty"`
	Echo             bool               `json:"echo,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	PresencePenalty  *float64           `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64           `json:"frequency_penalty,omitempty"`
	BestOf           *int               `json:"best_of,omitempty"`
	LogitBias        map[string]int     `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
	Suffix           string             `json:"suffix,omitempty"`
}

// CompletionResponse represents the response from text completion
type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`
}

// CompletionChoice represents a choice in completion response
type CompletionChoice struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}