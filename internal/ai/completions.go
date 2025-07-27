package ai

import "context"

// CreateCompletion creates a text completion (legacy endpoint)
func (c *Client) CreateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	resp, err := c.makeRequest(ctx, "POST", "/completions", req)
	if err != nil {
		return nil, err
	}

	var result CompletionResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}