# Panduan Lengkap AI Client Library

Library Go untuk koneksi ke penyedia API AI yang kompatibel dengan OpenAI. Library ini dirancang untuk fleksibilitas dan kemudahan penggunaan dalam aplikasi Go.

## 📋 Daftar Isi

1. [Instalasi dan Setup](#instalasi-dan-setup)
2. [Konfigurasi](#konfigurasi)
3. [Penggunaan Dasar](#penggunaan-dasar)
4. [Fitur Lanjutan](#fitur-lanjutan)
5. [Integrasi dengan Fiber](#integrasi-dengan-fiber)
6. [Provider yang Didukung](#provider-yang-didukung)
7. [Error Handling](#error-handling)
8. [Best Practices](#best-practices)
9. [Troubleshooting](#troubleshooting)

## 🚀 Instalasi dan Setup

Library ini sudah terintegrasi dalam project. Tidak perlu instalasi tambahan.

### Struktur File

```
internal/ai/
├── client.go          # Core client implementation
├── types.go           # Type definitions
├── chat.go            # Chat completion functions
├── embeddings.go      # Embedding functions
├── models.go          # Model management
├── completions.go     # Text completion (legacy)
├── config.go          # Configuration helpers
├── handler_example.go # Contoh handler untuk Fiber
├── examples.go        # Contoh penggunaan
├── example_usage.go   # Contoh penggunaan sederhana
├── client_test.go     # Unit tests
├── README.md          # Dokumentasi bahasa Inggris
└── PANDUAN.md         # Dokumentasi bahasa Indonesia
```

## ⚙️ Konfigurasi

### Environment Variables

Tambahkan konfigurasi berikut ke file `.env`:

```env
# Konfigurasi AI
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

### Load dari Environment

```go
config := ai.LoadConfigFromEnv()
client := ai.NewClient(config)
```

## 📖 Penggunaan Dasar

### 1. Chat Completion Sederhana

```go
package main

import (
    "context"
    "fmt"
    "log"
    "apiserver/internal/ai"
)

func main() {
    // Load konfigurasi
    config := ai.LoadConfigFromEnv()
    client := ai.NewClient(config)
    
    ctx := context.Background()
    
    // Buat request
    request := ai.ChatCompletionRequest{
        Model: ai.ModelGPT35Turbo,
        Messages: []ai.ChatMessage{
            ai.NewSystemMessage("Kamu adalah asisten yang membantu."),
            ai.NewUserMessage("Halo, apa kabar?"),
        },
        MaxTokens:   &[]int{100}[0],
        Temperature: &[]float64{0.7}[0],
    }
    
    // Kirim request
    response, err := client.CreateChatCompletion(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Respons AI:", response.Choices[0].Message.Content)
    fmt.Printf("Token digunakan: %d\n", response.Usage.TotalTokens)
}
```

### 2. Streaming Chat Completion

```go
func streamingChat() {
    config := ai.LoadConfigFromEnv()
    client := ai.NewClient(config)
    ctx := context.Background()
    
    request := ai.ChatCompletionRequest{
        Model: ai.ModelGPT35Turbo,
        Messages: []ai.ChatMessage{
            ai.NewUserMessage("Ceritakan tentang pemrograman Go"),
        },
        MaxTokens: &[]int{200}[0],
    }
    
    respChan, errChan := client.CreateChatCompletionStream(ctx, request)
    
    fmt.Print("Respons AI: ")
    for {
        select {
        case resp, ok := <-respChan:
            if !ok {
                fmt.Println("\n[Selesai]")
                return
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
}
```

### 3. Membuat Embeddings

```go
func createEmbedding() {
    config := ai.LoadConfigFromEnv()
    client := ai.NewClient(config)
    ctx := context.Background()
    
    // Single text
    resp, err := client.CreateEmbedding(ctx, ai.ModelTextEmbeddingAda002, "Halo dunia")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Dimensi embedding: %d\n", len(resp.Data[0].Embedding))
    fmt.Printf("3 nilai pertama: %v\n", resp.Data[0].Embedding[:3])
    
    // Multiple texts
    texts := []string{"Halo", "Dunia", "AI"}
    respMultiple, err := client.CreateEmbeddingsForTexts(ctx, ai.ModelTextEmbeddingAda002, texts)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Jumlah embeddings: %d\n", len(respMultiple.Data))
}
```

### 4. List Models

```go
func listModels() {
    config := ai.LoadConfigFromEnv()
    client := ai.NewClient(config)
    ctx := context.Background()
    
    models, err := client.ListModels(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Model tersedia: %d\n", len(models.Data))
    for _, model := range models.Data {
        fmt.Printf("- %s (pemilik: %s)\n", model.ID, model.OwnedBy)
    }
}
```

## 🔧 Fitur Lanjutan

### Helper Functions untuk Messages

```go
// System message
systemMsg := ai.NewSystemMessage("Kamu adalah asisten AI yang membantu")

// User message  
userMsg := ai.NewUserMessage("Pertanyaan saya adalah...")

// Assistant message
assistantMsg := ai.NewAssistantMessage("Jawaban saya adalah...")

// Tool message (untuk function calling)
toolMsg := ai.NewToolMessage("call_123", "Hasil function")
```

### Model Constants

```go
// Chat models
ai.ModelGPT4o              // gpt-4o
ai.ModelGPT4oMini          // gpt-4o-mini
ai.ModelGPT4Turbo          // gpt-4-turbo
ai.ModelGPT4               // gpt-4
ai.ModelGPT35Turbo         // gpt-3.5-turbo

// Embedding models
ai.ModelTextEmbedding3Small // text-embedding-3-small
ai.ModelTextEmbedding3Large // text-embedding-3-large
ai.ModelTextEmbeddingAda002 // text-embedding-ada-002
```

### Advanced Chat Parameters

```go
request := ai.ChatCompletionRequest{
    Model: ai.ModelGPT4,
    Messages: []ai.ChatMessage{
        ai.NewSystemMessage("Kamu adalah expert programmer"),
        ai.NewUserMessage("Jelaskan design patterns"),
    },
    MaxTokens:        &[]int{500}[0],
    Temperature:      &[]float64{0.7}[0],
    TopP:             &[]float64{0.9}[0],
    PresencePenalty:  &[]float64{0.1}[0],
    FrequencyPenalty: &[]float64{0.1}[0],
    Stop:             []string{"\n\n", "###"},
    User:             "user123",
}
```

## 🌐 Integrasi dengan Fiber

### Setup Handler

```go
package main

import (
    "apiserver/internal/ai"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()
    
    // Load AI configuration
    aiConfig := ai.LoadConfigFromEnv()
    
    // Create AI handler
    aiHandler := ai.NewChatHandler(aiConfig)
    
    // Register routes
    aiHandler.RegisterRoutes(app)
    
    app.Listen(":3000")
}
```

### Endpoints yang Tersedia

```bash
# Chat completion
POST /api/v1/ai/chat
{
    "message": "Halo, apa kabar?",
    "model": "gpt-3.5-turbo",
    "max_tokens": 100,
    "temperature": 0.7,
    "stream": false
}

# Create embedding
POST /api/v1/ai/embeddings
{
    "text": "Teks untuk di-embed",
    "model": "text-embedding-ada-002"
}

# List models
GET /api/v1/ai/models
```

### Custom Handler

```go
type MyAIHandler struct {
    aiClient *ai.Client
}

func NewMyAIHandler() *MyAIHandler {
    config := ai.LoadConfigFromEnv()
    client := ai.NewClient(config)
    
    return &MyAIHandler{
        aiClient: client,
    }
}

func (h *MyAIHandler) CustomChat(c *fiber.Ctx) error {
    var req struct {
        Question string `json:"question"`
        Context  string `json:"context,omitempty"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    messages := []ai.ChatMessage{
        ai.NewSystemMessage("Kamu adalah asisten yang membantu menjawab pertanyaan berdasarkan konteks yang diberikan."),
    }
    
    if req.Context != "" {
        messages = append(messages, ai.NewSystemMessage("Konteks: "+req.Context))
    }
    
    messages = append(messages, ai.NewUserMessage(req.Question))
    
    chatReq := ai.ChatCompletionRequest{
        Model:    ai.ModelGPT35Turbo,
        Messages: messages,
        MaxTokens: &[]int{300}[0],
    }
    
    resp, err := h.aiClient.CreateChatCompletion(c.Context(), chatReq)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(fiber.Map{
        "answer": resp.Choices[0].Message.Content,
        "usage":  resp.Usage,
    })
}
```

## 🔌 Provider yang Didukung

### OpenAI

```go
config := ai.ClientConfig{
    BaseURL: "https://api.openai.com/v1",
    APIKey:  "sk-...",
}
```

### Azure OpenAI

```go
config := ai.ClientConfig{
    BaseURL: "https://your-resource.openai.azure.com/openai/deployments/your-deployment/",
    APIKey:  "your-azure-key",
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

### Custom Provider

```go
config := ai.ClientConfig{
    BaseURL: "https://your-custom-api.com/v1",
    APIKey:  "your-custom-key",
    Timeout: 60 * time.Second,
}
```

## ❌ Error Handling

### Tipe Error

```go
resp, err := client.CreateChatCompletion(ctx, req)
if err != nil {
    // Cek apakah API error
    if apiErr, ok := err.(*ai.APIError); ok {
        fmt.Printf("API Error: %s\n", apiErr.ErrorInfo.Message)
        fmt.Printf("Type: %s\n", apiErr.ErrorInfo.Type)
        
        // Handle berdasarkan tipe error
        switch apiErr.ErrorInfo.Type {
        case "invalid_request_error":
            // Request tidak valid
        case "authentication_error":
            // API key salah
        case "rate_limit_exceeded":
            // Rate limit terlampaui
        case "insufficient_quota":
            // Quota habis
        }
    } else {
        // Network atau error lainnya
        fmt.Printf("Network Error: %v\n", err)
    }
}
```

### Retry Logic

```go
func chatWithRetry(client *ai.Client, req ai.ChatCompletionRequest, maxRetries int) (*ai.ChatCompletionResponse, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        resp, err := client.CreateChatCompletion(ctx, req)
        cancel()
        
        if err == nil {
            return resp, nil
        }
        
        lastErr = err
        
        // Cek apakah perlu retry
        if apiErr, ok := err.(*ai.APIError); ok {
            switch apiErr.ErrorInfo.Type {
            case "rate_limit_exceeded":
                // Wait dan retry
                time.Sleep(time.Duration(i+1) * time.Second)
                continue
            case "server_error":
                // Server error, retry
                time.Sleep(time.Duration(i+1) * time.Second)
                continue
            default:
                // Error lain, jangan retry
                return nil, err
            }
        }
        
        // Network error, retry
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

## 💡 Best Practices

### 1. Gunakan Context dengan Timeout

```go
// Selalu gunakan timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := client.CreateChatCompletion(ctx, req)
```

### 2. Handle Rate Limiting

```go
type RateLimitedClient struct {
    client *ai.Client
    limiter *rate.Limiter
}

func NewRateLimitedClient(config ai.ClientConfig, requestsPerSecond float64) *RateLimitedClient {
    return &RateLimitedClient{
        client:  ai.NewClient(config),
        limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), 1),
    }
}

func (r *RateLimitedClient) CreateChatCompletion(ctx context.Context, req ai.ChatCompletionRequest) (*ai.ChatCompletionResponse, error) {
    if err := r.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    return r.client.CreateChatCompletion(ctx, req)
}
```

### 3. Cache Embeddings

```go
type EmbeddingCache struct {
    client *ai.Client
    cache  map[string][]float64
    mutex  sync.RWMutex
}

func (e *EmbeddingCache) GetEmbedding(ctx context.Context, text string) ([]float64, error) {
    e.mutex.RLock()
    if embedding, exists := e.cache[text]; exists {
        e.mutex.RUnlock()
        return embedding, nil
    }
    e.mutex.RUnlock()
    
    resp, err := e.client.CreateEmbedding(ctx, ai.ModelTextEmbeddingAda002, text)
    if err != nil {
        return nil, err
    }
    
    embedding := resp.Data[0].Embedding
    
    e.mutex.Lock()
    e.cache[text] = embedding
    e.mutex.Unlock()
    
    return embedding, nil
}
```

### 4. Monitoring dan Logging

```go
type MonitoredClient struct {
    client *ai.Client
    logger *log.Logger
}

func (m *MonitoredClient) CreateChatCompletion(ctx context.Context, req ai.ChatCompletionRequest) (*ai.ChatCompletionResponse, error) {
    start := time.Now()
    
    resp, err := m.client.CreateChatCompletion(ctx, req)
    
    duration := time.Since(start)
    
    if err != nil {
        m.logger.Printf("Chat completion failed: %v (duration: %v)", err, duration)
    } else {
        m.logger.Printf("Chat completion success: %d tokens (duration: %v)", resp.Usage.TotalTokens, duration)
    }
    
    return resp, err
}
```

### 5. Configuration Management

```go
type AIService struct {
    clients map[string]*ai.Client
}

func NewAIService() *AIService {
    return &AIService{
        clients: map[string]*ai.Client{
            "openai": ai.NewClient(ai.ClientConfig{
                BaseURL: "https://api.openai.com/v1",
                APIKey:  os.Getenv("OPENAI_API_KEY"),
            }),
            "anthropic": ai.NewClient(ai.ClientConfig{
                BaseURL: "https://api.anthropic.com/v1",
                APIKey:  os.Getenv("ANTHROPIC_API_KEY"),
            }),
            "local": ai.NewClient(ai.ClientConfig{
                BaseURL: "http://localhost:11434/v1",
                APIKey:  "",
            }),
        },
    }
}

func (s *AIService) GetClient(provider string) *ai.Client {
    return s.clients[provider]
}
```

## 🔍 Troubleshooting

### Common Issues

#### 1. API Key Invalid
```
Error: API error: Incorrect API key provided (type: invalid_request_error)
```
**Solusi**: Pastikan API key valid dan memiliki permissions yang tepat.

#### 2. Rate Limit Exceeded
```
Error: API error: Rate limit reached (type: rate_limit_exceeded)
```
**Solusi**: Implementasikan retry logic dengan exponential backoff.

#### 3. Model Not Found
```
Error: API error: The model 'gpt-5' does not exist (type: invalid_request_error)
```
**Solusi**: Gunakan model yang tersedia. Cek dengan `ListModels()`.

#### 4. Timeout
```
Error: context deadline exceeded
```
**Solusi**: Tingkatkan timeout atau optimasi request.

### Debug Mode

```go
import "log"

// Enable debug logging
client := ai.NewClient(config)

// Log semua request/response
resp, err := client.CreateChatCompletion(ctx, req)
if err != nil {
    log.Printf("AI Request failed: %+v", req)
    log.Printf("AI Error: %v", err)
} else {
    log.Printf("AI Request: %+v", req)
    log.Printf("AI Response: %+v", resp)
}
```

### Testing dengan Mock

```go
func TestWithMockServer(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        response := ai.ChatCompletionResponse{
            Choices: []ai.ChatCompletionChoice{
                {Message: ai.ChatMessage{Content: "Test response"}},
            },
        }
        json.NewEncoder(w).Encode(response)
    }))
    defer server.Close()
    
    config := ai.ClientConfig{
        BaseURL: server.URL,
        APIKey:  "test-key",
    }
    client := ai.NewClient(config)
    
    // Test your code here
}
```

## 📝 Contoh Lengkap

Lihat file `example_usage.go` untuk contoh penggunaan lengkap yang mencakup:
- Chat completion sederhana
- Streaming chat
- Embeddings
- Multi-provider setup
- Error handling

## 🤝 Kontribusi

Untuk menambah fitur atau memperbaiki bug:

1. Buat test untuk fitur baru
2. Pastikan semua test pass: `go test ./internal/ai -v`
3. Update dokumentasi jika diperlukan
4. Follow Go best practices

## 📄 Lisensi

Library ini mengikuti lisensi project utama.