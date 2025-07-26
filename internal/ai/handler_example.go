package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ChatHandler example handler for AI chat functionality
type ChatHandler struct {
	client *Client
}

// NewChatHandler creates a new chat handler with AI client
func NewChatHandler(config ClientConfig) *ChatHandler {
	client := NewClient(config)
	return &ChatHandler{
		client: client,
	}
}

// ChatRequest represents the request structure for chat endpoint
type ChatRequest struct {
	Message     string  `json:"message" validate:"required"`
	Model       string  `json:"model,omitempty"`
	MaxTokens   *int    `json:"max_tokens,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	Stream      bool    `json:"stream,omitempty"`
}

// ChatResponse represents the response structure for chat endpoint
type ChatResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
	Usage    Usage  `json:"usage"`
}

// Chat handles chat completion requests
func (h *ChatHandler) Chat(c *fiber.Ctx) error {
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Set default model if not provided
	if req.Model == "" {
		req.Model = ModelGPT35Turbo
	}

	// Create chat completion request
	chatReq := ChatCompletionRequest{
		Model: req.Model,
		Messages: []ChatMessage{
			NewUserMessage(req.Message),
		},
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      req.Stream,
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	if req.Stream {
		return h.handleStreamingChat(c, ctx, chatReq)
	}

	return h.handleRegularChat(c, ctx, chatReq)
}

// handleRegularChat handles non-streaming chat completion
func (h *ChatHandler) handleRegularChat(c *fiber.Ctx, ctx context.Context, req ChatCompletionRequest) error {
	resp, err := h.client.CreateChatCompletion(ctx, req)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": apiErr.ErrorInfo.Message,
				"type":  apiErr.ErrorInfo.Type,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process request",
		})
	}

	if len(resp.Choices) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No response generated",
		})
	}

	// Handle content type assertion
	var content string
	if str, ok := resp.Choices[0].Message.Content.(string); ok {
		content = str
	} else {
		content = "Unable to parse response content"
	}

	response := ChatResponse{
		Response: content,
		Model:    resp.Model,
		Usage:    resp.Usage,
	}

	return c.JSON(response)
}

// handleStreamingChat handles streaming chat completion
func (h *ChatHandler) handleStreamingChat(c *fiber.Ctx, ctx context.Context, req ChatCompletionRequest) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Access-Control-Allow-Origin", "*")

	respChan, errChan := h.client.CreateChatCompletionStream(ctx, req)

	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				// Stream completed
				c.WriteString("data: [DONE]\n\n")
				return nil
			}

			if len(resp.Choices) > 0 && resp.Choices[0].Delta.Content != "" {
				data := map[string]interface{}{
					"content": resp.Choices[0].Delta.Content,
					"model":   resp.Model,
				}

				jsonData, _ := json.Marshal(data)
				c.WriteString(fmt.Sprintf("data: %s\n\n", jsonData))
			}

		case err := <-errChan:
			if err != nil {
				errorData := map[string]interface{}{
					"error": err.Error(),
				}
				jsonData, _ := json.Marshal(errorData)
				c.WriteString(fmt.Sprintf("data: %s\n\n", jsonData))
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// EmbeddingRequest represents the request structure for embedding endpoint
type EmbeddingRequest struct {
	Text  string `json:"text" validate:"required"`
	Model string `json:"model,omitempty"`
}

// EmbeddingResponse represents the response structure for embedding endpoint
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
	Model     string    `json:"model"`
	Usage     Usage     `json:"usage"`
}

// CreateEmbedding handles embedding creation requests
func (h *ChatHandler) CreateEmbedding(c *fiber.Ctx) error {
	var req EmbeddingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Set default model if not provided
	if req.Model == "" {
		req.Model = ModelTextEmbeddingAda002
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	resp, err := h.client.CreateEmbedding(ctx, req.Model, req.Text)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": apiErr.ErrorInfo.Message,
				"type":  apiErr.ErrorInfo.Type,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create embedding",
		})
	}

	if len(resp.Data) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No embedding generated",
		})
	}

	response := EmbeddingResponse{
		Embedding: resp.Data[0].Embedding,
		Model:     resp.Model,
		Usage:     resp.Usage,
	}

	return c.JSON(response)
}

// ListModels handles model listing requests
func (h *ChatHandler) ListModels(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.client.ListModels(ctx)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": apiErr.ErrorInfo.Message,
				"type":  apiErr.ErrorInfo.Type,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list models",
		})
	}

	return c.JSON(resp)
}

// RegisterRoutes registers AI-related routes to the Fiber app
func (h *ChatHandler) RegisterRoutes(app *fiber.App) {
	ai := app.Group("/api/v1/ai")

	ai.Post("/chat", h.Chat)
	ai.Post("/embeddings", h.CreateEmbedding)
	ai.Get("/models", h.ListModels)
}