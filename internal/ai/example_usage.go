package ai

import (
	"context"
	"fmt"
	"log"
	"os"
)

// SimpleExample demonstrates basic usage of the AI client
func SimpleExample() {
	// Set environment variables (in real usage, put these in .env file)
	os.Setenv("AI_BASE_URL", "https://api.openai.com/v1")
	os.Setenv("AI_API_KEY", "your-openai-api-key-here")
	os.Setenv("AI_TIMEOUT", "30")

	// Load configuration from environment
	config := LoadConfigFromEnv()
	
	// Create AI client
	client := NewClient(config)
	
	ctx := context.Background()

	// Example 1: Simple chat
	fmt.Println("=== Simple Chat Example ===")
	chatRequest := ChatCompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []ChatMessage{
			NewSystemMessage("You are a helpful programming assistant."),
			NewUserMessage("Explain what is REST API in simple terms."),
		},
		MaxTokens:   func(i int) *int { return &i }(150),
		Temperature: func(f float64) *float64 { return &f }(0.7),
	}

	response, err := client.CreateChatCompletion(ctx, chatRequest)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("AI Response: %s\n", response.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n\n", response.Usage.TotalTokens)

	// Example 2: Create embedding
	fmt.Println("=== Embedding Example ===")
	embeddingResp, err := client.CreateEmbedding(ctx, ModelTextEmbeddingAda002, "Hello, world!")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Embedding created with %d dimensions\n", len(embeddingResp.Data[0].Embedding))
	fmt.Printf("First 3 values: [%.4f, %.4f, %.4f]\n\n", 
		embeddingResp.Data[0].Embedding[0],
		embeddingResp.Data[0].Embedding[1], 
		embeddingResp.Data[0].Embedding[2])

	// Example 3: List available models
	fmt.Println("=== Available Models ===")
	models, err := client.ListModels(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Found %d models\n", len(models.Data))
	for i, model := range models.Data {
		if i < 3 { // Show first 3 models
			fmt.Printf("- %s (owned by: %s)\n", model.ID, model.OwnedBy)
		}
	}
}

// StreamingExample demonstrates streaming chat completion
func StreamingExample() {
	config := LoadConfigFromEnv()
	client := NewClient(config)
	ctx := context.Background()

	fmt.Println("=== Streaming Chat Example ===")
	
	request := ChatCompletionRequest{
		Model: ModelGPT35Turbo,
		Messages: []ChatMessage{
			NewUserMessage("Write a short poem about programming."),
		},
		MaxTokens:   func(i int) *int { return &i }(100),
		Temperature: func(f float64) *float64 { return &f }(0.8),
	}

	respChan, errChan := client.CreateChatCompletionStream(ctx, request)

	fmt.Print("AI Response: ")
	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				fmt.Println("\n[Stream completed]")
				return
			}
			if len(resp.Choices) > 0 && resp.Choices[0].Delta.Content != "" {
				fmt.Print(resp.Choices[0].Delta.Content)
			}
		case err := <-errChan:
			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				return
			}
		}
	}
}

// MultiProviderExample shows how to use different AI providers
func MultiProviderExample() {
	fmt.Println("=== Multi-Provider Example ===")

	// OpenAI Configuration
	openaiConfig := ClientConfig{
		BaseURL: OpenAIProvider.BaseURL,
		APIKey:  "your-openai-key",
	}
	openaiClient := NewClient(openaiConfig)

	// Local Ollama Configuration
	localConfig := ClientConfig{
		BaseURL: LocalProvider.BaseURL,
		APIKey:  "", // Usually not needed for local
	}
	localClient := NewClient(localConfig)

	// Anthropic Configuration
	anthropicConfig := ClientConfig{
		BaseURL: AnthropicProvider.BaseURL,
		APIKey:  "your-anthropic-key",
	}
	anthropicClient := NewClient(anthropicConfig)

	ctx := context.Background()

	// Example request
	request := ChatCompletionRequest{
		Model: ModelGPT35Turbo, // Change model based on provider
		Messages: []ChatMessage{
			NewUserMessage("What is artificial intelligence?"),
		},
		MaxTokens: func(i int) *int { return &i }(50),
	}

	// You can switch between clients based on your needs
	fmt.Println("Using OpenAI client...")
	_, _ = openaiClient.CreateChatCompletion(ctx, request)

	fmt.Println("Using local client...")
	request.Model = "llama2" // Change to local model
	_, _ = localClient.CreateChatCompletion(ctx, request)

	fmt.Println("Using Anthropic client...")
	request.Model = "claude-3-haiku-20240307" // Change to Claude model
	_, _ = anthropicClient.CreateChatCompletion(ctx, request)

	fmt.Println("Multi-provider setup complete!")
}

// ErrorHandlingExample demonstrates proper error handling
func ErrorHandlingExample() {
	fmt.Println("=== Error Handling Example ===")

	// Create client with invalid configuration to trigger errors
	config := ClientConfig{
		BaseURL: "https://invalid-api-endpoint.com/v1",
		APIKey:  "invalid-key",
	}
	client := NewClient(config)

	ctx := context.Background()
	request := ChatCompletionRequest{
		Model: "invalid-model",
		Messages: []ChatMessage{
			NewUserMessage("This will fail"),
		},
	}

	_, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		// Check if it's an API error
		if apiErr, ok := err.(*APIError); ok {
			fmt.Printf("API Error: %s\n", apiErr.ErrorInfo.Message)
			fmt.Printf("Error Type: %s\n", apiErr.ErrorInfo.Type)
			if apiErr.ErrorInfo.Code != "" {
				fmt.Printf("Error Code: %s\n", apiErr.ErrorInfo.Code)
			}
		} else {
			// Network or other errors
			fmt.Printf("Network/Other Error: %v\n", err)
		}
	}
}