# AI Client Library

Library Go untuk koneksi ke penyedia API AI yang kompatibel dengan OpenAI. Library ini mendukung berbagai provider AI seperti OpenAI, Anthropic, dan server lokal seperti Ollama.

## Fitur

- ✅ Chat Completions (streaming dan non-streaming)
- ✅ Text Embeddings
- ✅ Model Management
- ✅ Text Completions (legacy)
- ✅ Konfigurasi fleksibel untuk berbagai provider
- ✅ Error handling yang robust
- ✅ Context support untuk timeout dan cancellation
- ✅ Type-safe dengan struct yang lengkap

## Instalasi

Library ini sudah terintegrasi dalam project. Tidak perlu instalasi tambahan.

## Konfigurasi

### Environment Variables

Tambahkan konfigurasi berikut ke file `.env`:

```env
# AI Configuration
AI_BASE_URL=https://api.openai.com/v1
AI_API_KEY=your-openai-api-key-here
AI_TIMEOUT=30
```

### Konfigurasi Manual

```go
import "apiserver/internal/ai"

config := ai.ClientConfig{
    BaseURL: "https://api.openai.com/v1",
    APIKey:  "your-api-key",
    Timeout: 30 * time.Second,
}

client := ai.NewClient(config)
```

## Penggunaan

### 1. Chat Completion

```go
package main

import (
    "context"
    "fmt"
    "apiserver/internal/ai"
)

func main() {
    // Load dari environment
    config := ai.LoadConfigFromEnv()
    client := ai.NewClient(config)
    
    ctx := context.Background()
    
    req := ai.ChatCompletionRequest{
        Model: ai.ModelGPT35Turbo,
        Messages: []ai.ChatMessage{
            ai.NewSystemMessage("You are a helpful assistant."),
            ai.NewUserMessage("Hello, how are you?"),
        },
        MaxTokens:   &[]int{100}[0],
        Temperature: &[]float64{0.7}[0],
    }
    
    resp, err := client.CreateChatCompletion(ctx, req)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(resp.Choices[0].Message.Content)
}
```

### 2. Streaming Chat Completion

```go
req := ai.ChatCompletionRequest{
    Model: ai.ModelGPT35Turbo,
    Messages: []ai.ChatMessage{
        ai.NewUserMessage("Tell me a story"),
    },
}

respChan, errChan := client.CreateChatCompletionStream(ctx, req)

for {
    select {
    case resp, ok := <-respChan:
        if !ok {
            return // Stream selesai
        }
        if len(resp.Choices) > 0 {
            fmt.Print(resp.Choices[0].Delta.Content)
        }
    case err := <-errChan:
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }
    }
}
```

### 3. Embeddings

```go
// Single text
resp, err := client.CreateEmbedding(ctx, ai.ModelTextEmbeddingAda002, "Hello world")
if err != nil {
    panic(err)
}

fmt.Printf("Embedding: %v\n", resp.Data[0].Embedding)

// Multiple texts
texts := []string{"Hello", "World", "AI"}
resp, err := client.CreateEmbeddingsForTexts(ctx, ai.ModelTextEmbeddingAda002, texts)
```

### 4. List Models

```go
models, err := client.ListModels(ctx)
if err != nil {
    panic(err)
}

for _, model := range models.Data {
    fmt.Printf("Model: %s\n", model.ID)
}
```

## Provider yang Didukung

### OpenAI
```go
config := ai.ClientConfig{
    BaseURL: "https://api.openai.com/v1",
    APIKey:  "sk-...",
}
```

### Anthropic (Claude)
```go
config := ai.ClientConfig{
    BaseURL: "https://api.anthropic.com/v1",
    APIKey:  "your-anthropic-key",
}
```

### Ollama (Local)
```go
config := ai.ClientConfig{
    BaseURL: "http://localhost:11434/v1",
    APIKey:  "", // Biasanya tidak diperlukan
}
```

### Azure OpenAI
```go
config := ai.ClientConfig{
    BaseURL: "https://your-resource.openai.azure.com/openai/deployments/your-deployment",
    APIKey:  "your-azure-key",
}
```

## Model Constants

Library menyediakan konstanta untuk model yang umum digunakan:

```go
ai.ModelGPT4o              // gpt-4o
ai.ModelGPT4oMini          // gpt-4o-mini
ai.ModelGPT4Turbo          // gpt-4-turbo
ai.ModelGPT4               // gpt-4
ai.ModelGPT35Turbo         // gpt-3.5-turbo
ai.ModelTextEmbedding3Small // text-embedding-3-small
ai.ModelTextEmbedding3Large // text-embedding-3-large
ai.ModelTextEmbeddingAda002 // text-embedding-ada-002
```

## Error Handling

Library menggunakan custom error types untuk handling yang lebih baik:

```go
resp, err := client.CreateChatCompletion(ctx, req)
if err != nil {
    if apiErr, ok := err.(*ai.APIError); ok {
        fmt.Printf("API Error: %s (Type: %s)\n", apiErr.Error.Message, apiErr.Error.Type)
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
    return
}
```

## Helper Functions

Library menyediakan helper functions untuk membuat message:

```go
// System message
systemMsg := ai.NewSystemMessage("You are a helpful assistant")

// User message  
userMsg := ai.NewUserMessage("Hello!")

// Assistant message
assistantMsg := ai.NewAssistantMessage("Hi there!")

// Tool message
toolMsg := ai.NewToolMessage("call_123", "Function result")
```

## Contoh Penggunaan dalam Handler

```go
package handlers

import (
    "apiserver/internal/ai"
    "github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
    aiClient *ai.Client
}

func NewChatHandler() *ChatHandler {
    config := ai.LoadConfigFromEnv()
    client := ai.NewClient(config)
    
    return &ChatHandler{
        aiClient: client,
    }
}

func (h *ChatHandler) Chat(c *fiber.Ctx) error {
    var req struct {
        Message string `json:"message"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    chatReq := ai.ChatCompletionRequest{
        Model: ai.ModelGPT35Turbo,
        Messages: []ai.ChatMessage{
            ai.NewUserMessage(req.Message),
        },
    }
    
    resp, err := h.aiClient.CreateChatCompletion(c.Context(), chatReq)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(fiber.Map{
        "response": resp.Choices[0].Message.Content,
        "usage":    resp.Usage,
    })
}
```

## Testing

Untuk testing, Anda dapat menggunakan mock client atau test server:

```go
// Test dengan mock
func TestChatCompletion(t *testing.T) {
    config := ai.ClientConfig{
        BaseURL: "http://localhost:8080", // Test server
        APIKey:  "test-key",
    }
    
    client := ai.NewClient(config)
    // ... test logic
}
```

## Best Practices

1. **Gunakan Context**: Selalu gunakan context untuk timeout dan cancellation
2. **Handle Errors**: Periksa dan handle error dengan proper
3. **Rate Limiting**: Implementasikan rate limiting untuk production
4. **Caching**: Cache embeddings dan responses jika memungkinkan
5. **Monitoring**: Monitor usage dan costs
6. **Security**: Jangan hardcode API keys, gunakan environment variables

## Troubleshooting

### Common Issues

1. **Invalid API Key**: Pastikan API key valid dan memiliki permissions yang tepat
2. **Rate Limiting**: Implementasikan retry logic dengan exponential backoff
3. **Timeout**: Sesuaikan timeout berdasarkan kebutuhan
4. **Model Not Found**: Pastikan model tersedia di provider yang digunakan

### Debug Mode

Untuk debugging, Anda dapat menambahkan logging:

```go
import "log"

// Log request/response untuk debugging
resp, err := client.CreateChatCompletion(ctx, req)
if err != nil {
    log.Printf("AI API Error: %v", err)
}
log.Printf("AI Response: %+v", resp)
```