package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	config := ClientConfig{
		BaseURL: "https://api.openai.com/v1",
		APIKey:  "test-key",
		Timeout: 10 * time.Second,
	}

	client := NewClient(config)

	if client.BaseURL != config.BaseURL {
		t.Errorf("Expected BaseURL %s, got %s", config.BaseURL, client.BaseURL)
	}

	if client.APIKey != config.APIKey {
		t.Errorf("Expected APIKey %s, got %s", config.APIKey, client.APIKey)
	}

	if client.timeout != config.Timeout {
		t.Errorf("Expected timeout %v, got %v", config.Timeout, client.timeout)
	}
}

func TestNewClientDefaults(t *testing.T) {
	config := ClientConfig{
		APIKey: "test-key",
	}

	client := NewClient(config)

	expectedBaseURL := "https://api.openai.com/v1"
	if client.BaseURL != expectedBaseURL {
		t.Errorf("Expected default BaseURL %s, got %s", expectedBaseURL, client.BaseURL)
	}

	expectedTimeout := 30 * time.Second
	if client.timeout != expectedTimeout {
		t.Errorf("Expected default timeout %v, got %v", expectedTimeout, client.timeout)
	}
}

func TestCreateChatCompletion(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/chat/completions" {
			t.Errorf("Expected path /chat/completions, got %s", r.URL.Path)
		}

		// Check authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got %s", auth)
		}

		// Mock response
		response := ChatCompletionResponse{
			ID:      "chatcmpl-123",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   "gpt-3.5-turbo",
			Choices: []ChatCompletionChoice{
				{
					Index: 0,
					Message: ChatMessage{
						Role:    "assistant",
						Content: "Hello! How can I help you today?",
					},
					FinishReason: "stop",
				},
			},
			Usage: Usage{
				PromptTokens:     10,
				CompletionTokens: 9,
				TotalTokens:      19,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	config := ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
	}
	client := NewClient(config)

	// Test request
	req := ChatCompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []ChatMessage{
			NewUserMessage("Hello"),
		},
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.ID != "chatcmpl-123" {
		t.Errorf("Expected ID 'chatcmpl-123', got %s", resp.ID)
	}

	if len(resp.Choices) != 1 {
		t.Errorf("Expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Content != "Hello! How can I help you today?" {
		t.Errorf("Unexpected response content: %s", resp.Choices[0].Message.Content)
	}
}

func TestCreateEmbedding(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/embeddings" {
			t.Errorf("Expected path /embeddings, got %s", r.URL.Path)
		}

		// Mock response
		response := EmbeddingAPIResponse{
			Object: "list",
			Data: []Embedding{
				{
					Object:    "embedding",
					Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
					Index:     0,
				},
			},
			Model: "text-embedding-ada-002",
			Usage: Usage{
				PromptTokens: 5,
				TotalTokens:  5,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	config := ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
	}
	client := NewClient(config)

	ctx := context.Background()
	resp, err := client.CreateEmbedding(ctx, ModelTextEmbeddingAda002, "Hello world")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 embedding, got %d", len(resp.Data))
	}

	if len(resp.Data[0].Embedding) != 5 {
		t.Errorf("Expected embedding length 5, got %d", len(resp.Data[0].Embedding))
	}

	expectedEmbedding := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	for i, val := range resp.Data[0].Embedding {
		if val != expectedEmbedding[i] {
			t.Errorf("Expected embedding[%d] = %f, got %f", i, expectedEmbedding[i], val)
		}
	}
}

func TestListModels(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/models" {
			t.Errorf("Expected path /models, got %s", r.URL.Path)
		}

		// Mock response
		response := ModelsResponse{
			Object: "list",
			Data: []Model{
				{
					ID:      "gpt-3.5-turbo",
					Object:  "model",
					Created: time.Now().Unix(),
					OwnedBy: "openai",
				},
				{
					ID:      "gpt-4",
					Object:  "model",
					Created: time.Now().Unix(),
					OwnedBy: "openai",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	config := ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
	}
	client := NewClient(config)

	ctx := context.Background()
	resp, err := client.ListModels(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 models, got %d", len(resp.Data))
	}

	if resp.Data[0].ID != "gpt-3.5-turbo" {
		t.Errorf("Expected first model ID 'gpt-3.5-turbo', got %s", resp.Data[0].ID)
	}

	if resp.Data[1].ID != "gpt-4" {
		t.Errorf("Expected second model ID 'gpt-4', got %s", resp.Data[1].ID)
	}
}

func TestAPIError(t *testing.T) {
	// Mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := APIError{
			ErrorInfo: struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Param   string `json:"param,omitempty"`
				Code    string `json:"code,omitempty"`
			}{
				Message: "Invalid request",
				Type:    "invalid_request_error",
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
	}))
	defer server.Close()

	// Create client with mock server
	config := ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
	}
	client := NewClient(config)

	req := ChatCompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []ChatMessage{
			NewUserMessage("Hello"),
		},
	}

	ctx := context.Background()
	_, err := client.CreateChatCompletion(ctx, req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}

	if apiErr.ErrorInfo.Message != "Invalid request" {
		t.Errorf("Expected error message 'Invalid request', got %s", apiErr.ErrorInfo.Message)
	}

	if apiErr.ErrorInfo.Type != "invalid_request_error" {
		t.Errorf("Expected error type 'invalid_request_error', got %s", apiErr.ErrorInfo.Type)
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test NewSystemMessage
	sysMsg := NewSystemMessage("You are a helpful assistant")
	if sysMsg.Role != "system" {
		t.Errorf("Expected role 'system', got %s", sysMsg.Role)
	}
	if sysMsg.Content != "You are a helpful assistant" {
		t.Errorf("Expected content 'You are a helpful assistant', got %s", sysMsg.Content)
	}

	// Test NewUserMessage
	userMsg := NewUserMessage("Hello")
	if userMsg.Role != "user" {
		t.Errorf("Expected role 'user', got %s", userMsg.Role)
	}
	if userMsg.Content != "Hello" {
		t.Errorf("Expected content 'Hello', got %s", userMsg.Content)
	}

	// Test NewAssistantMessage
	assistantMsg := NewAssistantMessage("Hi there!")
	if assistantMsg.Role != "assistant" {
		t.Errorf("Expected role 'assistant', got %s", assistantMsg.Role)
	}
	if assistantMsg.Content != "Hi there!" {
		t.Errorf("Expected content 'Hi there!', got %s", assistantMsg.Content)
	}

	// Test NewToolMessage
	toolMsg := NewToolMessage("call_123", "Function result")
	if toolMsg.Role != "tool" {
		t.Errorf("Expected role 'tool', got %s", toolMsg.Role)
	}
	if toolMsg.Content != "Function result" {
		t.Errorf("Expected content 'Function result', got %s", toolMsg.Content)
	}
	if toolMsg.ToolCallID != "call_123" {
		t.Errorf("Expected tool call ID 'call_123', got %s", toolMsg.ToolCallID)
	}
}