package ai

import "context"

// CreateEmbeddings creates embeddings for the given input
func (c *Client) CreateEmbeddings(ctx context.Context, req EmbeddingAPIRequest) (*EmbeddingAPIResponse, error) {
	resp, err := c.makeRequest(ctx, "POST", "/embeddings", req)
	if err != nil {
		return nil, err
	}

	var result EmbeddingAPIResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateEmbedding creates a single embedding for text input
func (c *Client) CreateEmbedding(ctx context.Context, model, text string) (*EmbeddingAPIResponse, error) {
	req := EmbeddingAPIRequest{
		Model: model,
		Input: text,
	}

	return c.CreateEmbeddings(ctx, req)
}

// CreateEmbeddingsForTexts creates embeddings for multiple text inputs
func (c *Client) CreateEmbeddingsForTexts(ctx context.Context, model string, texts []string) (*EmbeddingAPIResponse, error) {
	req := EmbeddingAPIRequest{
		Model: model,
		Input: texts,
	}

	return c.CreateEmbeddings(ctx, req)
}