package ai

import "context"

// ListModels retrieves available models
func (c *Client) ListModels(ctx context.Context) (*ModelsResponse, error) {
	resp, err := c.makeRequest(ctx, "GET", "/models", nil)
	if err != nil {
		return nil, err
	}

	var result ModelsResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetModel retrieves information about a specific model
func (c *Client) GetModel(ctx context.Context, modelID string) (*Model, error) {
	endpoint := "/models/" + modelID
	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var result Model
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}