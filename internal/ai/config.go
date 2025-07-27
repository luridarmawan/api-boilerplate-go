package ai

import (
	"os"
	"strconv"
	"time"
)

// LoadConfigFromEnv loads AI client configuration from environment variables
func LoadConfigFromEnv() ClientConfig {
	config := ClientConfig{
		BaseURL: getEnvOrDefault("AI_BASE_URL", "https://api.openai.com/v1"),
		APIKey:  getEnvOrDefault("AI_API_KEY", ""),
	}

	// Parse timeout from environment
	if timeoutStr := os.Getenv("AI_TIMEOUT"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			config.Timeout = time.Duration(timeout) * time.Second
		}
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return config
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Common model constants
const (
	ModelGPT4o        = "gpt-4o"
	ModelGPT4oMini    = "gpt-4o-mini"
	ModelGPT4Turbo    = "gpt-4-turbo"
	ModelGPT4         = "gpt-4"
	ModelGPT35Turbo   = "gpt-3.5-turbo"
	ModelTextEmbedding3Small = "text-embedding-3-small"
	ModelTextEmbedding3Large = "text-embedding-3-large"
	ModelTextEmbeddingAda002 = "text-embedding-ada-002"
)

// Provider-specific configurations
type ProviderConfig struct {
	Name     string
	BaseURL  string
	Models   []string
}

// Common AI providers
var (
	OpenAIProvider = ProviderConfig{
		Name:    "OpenAI",
		BaseURL: "https://api.openai.com/v1",
		Models:  []string{ModelGPT4o, ModelGPT4oMini, ModelGPT4Turbo, ModelGPT4, ModelGPT35Turbo},
	}

	AnthropicProvider = ProviderConfig{
		Name:    "Anthropic",
		BaseURL: "https://api.anthropic.com/v1",
		Models:  []string{"claude-3-5-sonnet-20241022", "claude-3-haiku-20240307"},
	}

	LocalProvider = ProviderConfig{
		Name:    "Local",
		BaseURL: "http://localhost:11434/v1",
		Models:  []string{"llama2", "codellama", "mistral"},
	}
)