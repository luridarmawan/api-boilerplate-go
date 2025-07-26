package ai

import (
	"context"
	"fmt"
	"log"
)

// ExampleUsage demonstrates how to use the AI client
func ExampleUsage() {
	// Load configuration from environment
	config := LoadConfigFromEnv()
	
	// Create client
	client := NewClient(config)
	
	ctx := context.Background()

	// Example 1: Simple chat completion
	fmt.Println("=== Example 1: Simple Chat Completion ===")
	chatReq := ChatCompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []ChatMessage{
			NewSystemMessage("You are a helpful assistant."),
			NewUserMessage("Hello, how are you?"),
		},
		MaxTokens:   func(i int) *int { return &i }(100),
		Temperature: func(f float64) *float64 { return &f }(0.7),
	}

	chatResp, err := client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Handle content type assertion
	var content string
	if str, ok := chatResp.Choices[0].Message.Content.(string); ok {
		content = str
	} else {
		content = "Unable to parse response content"
	}
	fmt.Printf("Response: %s\n", content)
	fmt.Printf("Usage: %+v\n", chatResp.Usage)

	// Example 2: Streaming chat completion
	fmt.Println("\n=== Example 2: Streaming Chat Completion ===")
	streamReq := ChatCompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []ChatMessage{
			NewUserMessage("Tell me a short story about a robot."),
		},
		MaxTokens:   func(i int) *int { return &i }(200),
		Temperature: func(f float64) *float64 { return &f }(0.8),
	}

	respChan, errChan := client.CreateChatCompletionStream(ctx, streamReq)

	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				fmt.Println("\nStream completed")
				goto nextExample
			}
			if len(resp.Choices) > 0 && resp.Choices[0].Delta.Content != "" {
				fmt.Print(resp.Choices[0].Delta.Content)
			}
		case err := <-errChan:
			if err != nil {
				log.Printf("Stream error: %v", err)
				goto nextExample
			}
		}
	}

nextExample:
	// Example 3: Create embeddings
	fmt.Println("\n=== Example 3: Create Embeddings ===")
	embeddingResp, err := client.CreateEmbedding(ctx, ModelTextEmbeddingAda002, "Hello world")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Embedding dimensions: %d\n", len(embeddingResp.Data[0].Embedding))
	fmt.Printf("First 5 values: %v\n", embeddingResp.Data[0].Embedding[:5])

	// Example 4: List models
	fmt.Println("\n=== Example 4: List Models ===")
	models, err := client.ListModels(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Available models: %d\n", len(models.Data))
	for i, model := range models.Data {
		if i < 5 { // Show first 5 models
			fmt.Printf("- %s (owned by: %s)\n", model.ID, model.OwnedBy)
		}
	}
}

// ExampleWithCustomProvider shows how to use with different AI providers
func ExampleWithCustomProvider() {
	// Example with Anthropic (Claude)
	anthropicConfig := ClientConfig{
		BaseURL: AnthropicProvider.BaseURL,
		APIKey:  "your-anthropic-api-key",
	}
	anthropicClient := NewClient(anthropicConfig)

	// Example with local Ollama
	localConfig := ClientConfig{
		BaseURL: LocalProvider.BaseURL,
		APIKey:  "", // Usually not needed for local
	}
	localClient := NewClient(localConfig)

	ctx := context.Background()

	// Use with different models
	req := ChatCompletionRequest{
		Model: "claude-3-haiku-20240307", // or "llama2" for local
		Messages: []ChatMessage{
			NewUserMessage("Explain quantum computing in simple terms."),
		},
	}

	// Choose client based on your needs
	_ = anthropicClient
	_ = localClient
	_ = req
	_ = ctx

	fmt.Println("Custom provider examples configured")
}