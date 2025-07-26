package ai

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// CreateChatCompletion creates a chat completion
func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	resp, err := c.makeRequest(ctx, "POST", "/chat/completions", req)
	if err != nil {
		return nil, err
	}

	var result ChatCompletionResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateChatCompletionStream creates a streaming chat completion
func (c *Client) CreateChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (<-chan StreamResponse, <-chan error) {
	req.Stream = true
	
	respChan := make(chan StreamResponse)
	errChan := make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		resp, err := c.makeRequest(ctx, "POST", "/chat/completions", req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			var apiError APIError
			if err := c.handleResponse(resp, &apiError); err != nil {
				errChan <- err
				return
			}
			errChan <- &apiError
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			
			// Skip empty lines and comments
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			// Remove "data: " prefix
			if strings.HasPrefix(line, "data: ") {
				line = strings.TrimPrefix(line, "data: ")
			}

			// Check for end of stream
			if line == "[DONE]" {
				return
			}

			var streamResp StreamResponse
			if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
				// Skip malformed JSON lines
				continue
			}

			select {
			case respChan <- streamResp:
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("error reading stream: %w", err)
		}
	}()

	return respChan, errChan
}

// Helper functions for creating common message types

// NewSystemMessage creates a system message
func NewSystemMessage(content string) ChatMessage {
	return ChatMessage{
		Role:    "system",
		Content: content,
	}
}

// NewUserMessage creates a user message
func NewUserMessage(content string) ChatMessage {
	return ChatMessage{
		Role:    "user",
		Content: content,
	}
}

// NewUserMessageWithImage creates a user message with image
func NewUserMessageWithImage(text, imageURL string) ChatMessage {
	return ChatMessage{
		Role: "user",
		Content: []MessageContent{
			{
				Type: "text",
				Text: text,
			},
			{
				Type: "image_url",
				ImageURL: &ImageURL{
					URL:    imageURL,
					Detail: "high",
				},
			},
		},
	}
}

// NewAssistantMessage creates an assistant message
func NewAssistantMessage(content string) ChatMessage {
	return ChatMessage{
		Role:    "assistant",
		Content: content,
	}
}

// NewToolMessage creates a tool message
func NewToolMessage(toolCallID, content string) ChatMessage {
	return ChatMessage{
		Role:       "tool",
		Content:    content,
		ToolCallID: toolCallID,
	}
}