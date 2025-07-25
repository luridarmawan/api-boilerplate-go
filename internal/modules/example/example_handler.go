package example

import (
	"context"
	"strings"
	"time"

	"apiserver/internal/ai"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	repo     Repository
	aiClient *ai.Client
}

func NewHandler(repo Repository) *Handler {
	// Load AI configuration
	aiConfig := ai.LoadConfigFromEnv()
	aiClient := ai.NewClient(aiConfig)

	return &Handler{
		repo:     repo,
		aiClient: aiClient,
	}
}

// Helper function to create AI client based on request
func (h *Handler) getAIClient(req ChatCompletionRequest, timeout time.Duration) *ai.Client {
	if req.CustomEndpoint != "" && req.CustomAPIKey != "" {
		// Use custom endpoint and API key
		customConfig := ai.ClientConfig{
			BaseURL: req.CustomEndpoint,
			APIKey:  req.CustomAPIKey,
			Timeout: timeout,
		}
		return ai.NewClient(customConfig)
	}
	// Use default client
	return h.aiClient
}

// Helper function to validate chat completion request
func (h *Handler) validateChatRequest(req ChatCompletionRequest) error {
	if strings.TrimSpace(req.Message) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Message is required")
	}

	// Validate custom endpoint and API key
	if (req.CustomEndpoint != "" && req.CustomAPIKey == "") || (req.CustomEndpoint == "" && req.CustomAPIKey != "") {
		return fiber.NewError(fiber.StatusBadRequest, "Both custom_endpoint and custom_api_key must be provided together")
	}

	return nil
}

// CreateExample godoc
// @Summary Create a new example
// @Description Create a new example with name and description
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param example body CreateExampleRequest true "Example data"
// @Success 201 {object} Example
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples [post]
func (h *Handler) CreateExample(c *fiber.Ctx) error {
	var req CreateExampleRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Simple validation
	if strings.TrimSpace(req.Name) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Name is required",
		})
	}

	example := &Example{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.repo.CreateExample(example); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create example",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   example,
	})
}

// GetExamples godoc
// @Summary Get all examples
// @Description Get list of all examples
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Example
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples [get]
func (h *Handler) GetExamples(c *fiber.Ctx) error {
	examples, err := h.repo.GetAllExamples()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch examples",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   examples,
	})
}

// GetExample godoc
// @Summary Get example by ID
// @Description Get a specific example by its ID
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Example ID"
// @Success 200 {object} Example
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/{id} [get]
func (h *Handler) GetExample(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid example ID",
		})
	}

	example, err := h.repo.GetExampleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Example not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   example,
	})
}

// UpdateExample godoc
// @Summary Update example
// @Description Update an existing example
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Example ID"
// @Param example body CreateExampleRequest true "Example data"
// @Success 200 {object} Example
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/{id} [put]
func (h *Handler) UpdateExample(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid example ID",
		})
	}

	// Check if example exists
	example, err := h.repo.GetExampleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Example not found",
		})
	}

	var req CreateExampleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Simple validation
	if strings.TrimSpace(req.Name) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Name is required",
		})
	}

	// Update example
	example.Name = req.Name
	example.Description = req.Description

	if err := h.repo.UpdateExample(example); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update example",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   example,
	})
}

// SoftDeleteExample godoc
// @Summary Soft delete example
// @Description Soft delete an example (set status_id to 0)
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Example ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/{id} [delete]
func (h *Handler) SoftDeleteExample(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid example ID",
		})
	}

	// Check if example exists
	_, err := h.repo.GetExampleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Example not found",
		})
	}

	if err := h.repo.SoftDeleteExample(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete example",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Example deleted successfully",
	})
}

// RestoreExample godoc
// @Summary Restore example
// @Description Restore a soft deleted example (set status_id to 1)
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Example ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/{id}/restore [post]
func (h *Handler) RestoreExample(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid example ID",
		})
	}

	if err := h.repo.RestoreExample(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to restore example",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Example restored successfully",
	})
}

// GetDeletedExamples godoc
// @Summary Get deleted examples
// @Description Get list of all soft deleted examples
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Example
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/deleted [get]
func (h *Handler) GetDeletedExamples(c *fiber.Ctx) error {
	examples, err := h.repo.GetDeletedExamples()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch deleted examples",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   examples,
	})
}

// ChatCompletion godoc
// @Summary AI Chat Completion
// @Description Create a chat completion using AI (OpenAI compatible API)
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chat body ChatCompletionRequest true "Chat completion request"
// @Success 200 {object} ChatCompletionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/chat/completion [post]
func (h *Handler) ChatCompletion(c *fiber.Ctx) error {
	var req ChatCompletionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validation
	if err := h.validateChatRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Set default values
	if req.Model == "" {
		req.Model = ai.ModelGPT35Turbo
	}
	if req.MaxTokens == nil {
		defaultMaxTokens := 500
		req.MaxTokens = &defaultMaxTokens
	}
	if req.Temperature == nil {
		defaultTemperature := 0.7
		req.Temperature = &defaultTemperature
	}
	if req.SystemPrompt == "" {
		req.SystemPrompt = "You are a helpful assistant from CARIK.id that provides clear and concise answers."
	}

	// Create AI request
	messages := []ai.ChatMessage{
		ai.NewSystemMessage(req.SystemPrompt),
		ai.NewUserMessage(req.Message),
	}

	aiRequest := ai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Record start time for processing time calculation
	startTime := time.Now()

	// Get appropriate AI client
	aiClient := h.getAIClient(req, 30*time.Second)

	// Call AI API
	aiResponse, err := aiClient.CreateChatCompletion(ctx, aiRequest)
	if err != nil {
		// Handle AI API errors
		if apiErr, ok := err.(*ai.APIError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "AI API Error: " + apiErr.ErrorInfo.Message,
				"type":    apiErr.ErrorInfo.Type,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to process AI request: " + err.Error(),
		})
	}

	// Calculate processing time
	processingTime := time.Since(startTime)

	// Check if we have a response
	if len(aiResponse.Choices) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "No response generated from AI",
		})
	}

	// Build response
	response := ChatCompletionResponse{
		Response: aiResponse.Choices[0].Message.Content,
		Model:    aiResponse.Model,
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     aiResponse.Usage.PromptTokens,
			CompletionTokens: aiResponse.Usage.CompletionTokens,
			TotalTokens:      aiResponse.Usage.TotalTokens,
		},
		ProcessingTime: processingTime.String(),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

// ChatCompletionStream godoc
// @Summary AI Chat Completion Stream
// @Description Create a streaming chat completion using AI (OpenAI compatible API)
// @Tags Example
// @Accept json
// @Produce text/event-stream
// @Security BearerAuth
// @Param chat body ChatCompletionRequest true "Chat completion request"
// @Success 200 {string} string "Server-Sent Events stream"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/chat/completion/stream [post]
func (h *Handler) ChatCompletionStream(c *fiber.Ctx) error {
	var req ChatCompletionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validation
	if err := h.validateChatRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Set default values
	if req.Model == "" {
		req.Model = ai.ModelGPT35Turbo
	}
	if req.MaxTokens == nil {
		defaultMaxTokens := 500
		req.MaxTokens = &defaultMaxTokens
	}
	if req.Temperature == nil {
		defaultTemperature := 0.7
		req.Temperature = &defaultTemperature
	}
	if req.SystemPrompt == "" {
		req.SystemPrompt = "You are a helpful assistant from CARIK.id that provides clear and concise answers."
	}

	// Set headers for Server-Sent Events
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Access-Control-Allow-Origin", "*")

	// Create AI request
	messages := []ai.ChatMessage{
		ai.NewSystemMessage(req.SystemPrompt),
		ai.NewUserMessage(req.Message),
	}

	aiRequest := ai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      true, // Enable streaming
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(c.Context(), 60*time.Second)
	defer cancel()

	// Get appropriate AI client
	aiClient := h.getAIClient(req, 60*time.Second)

	// Call AI API for streaming
	respChan, errChan := aiClient.CreateChatCompletionStream(ctx, aiRequest)

	// Send initial event
	c.WriteString("data: {\"type\":\"start\",\"message\":\"Starting AI response...\"}\n\n")

	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				// Stream completed
				c.WriteString("data: {\"type\":\"end\",\"message\":\"Stream completed\"}\n\n")
				c.WriteString("data: [DONE]\n\n")
				return nil
			}

			if len(resp.Choices) > 0 && resp.Choices[0].Delta.Content != "" {
				// Send content chunk
				c.WriteString("data: {\"type\":\"content\",\"content\":\"" +
					strings.ReplaceAll(resp.Choices[0].Delta.Content, "\"", "\\\"") +
					"\",\"model\":\"" + resp.Model + "\"}\n\n")
			}

		case err := <-errChan:
			if err != nil {
				// Send error event
				c.WriteString("data: {\"type\":\"error\",\"message\":\"" +
					strings.ReplaceAll(err.Error(), "\"", "\\\"") + "\"}\n\n")
				return err
			}

		case <-ctx.Done():
			// Timeout or cancellation
			c.WriteString("data: {\"type\":\"error\",\"message\":\"Request timeout or cancelled\"}\n\n")
			return ctx.Err()
		}
	}
}