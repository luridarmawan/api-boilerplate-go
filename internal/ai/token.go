package ai

import "strings"

// EstimateOCRTokenUsage provides token usage estimation for OCR processing
// Based on image size and document type complexity
func EstimateOCRTokenUsage(imageData []byte, docType string) (tokenIn, tokenOut, tokenTotal int) {
	if len(imageData) == 0 {
		return 0, 0, 0
	}

	// Base estimation: file size in KB * 12 tokens per KB (rough estimate for images)
	fileSize := len(imageData)
	baseTokens := (fileSize / 1024) * 12

	// Minimum tokens for any image processing
	if baseTokens < 400 {
		baseTokens = 400
	}

	// Adjustment based on document type complexity
	multiplier := 1.0
	switch docType {
	case "receipt":
		multiplier = 1.0 // Normal complexity
	case "ktp", "passport", "sim", "driver_license":
		multiplier = 1.2 // More detailed structured documents
	case "invoice":
		multiplier = 1.3 // Most complex with tables and details
	default:
		multiplier = 1.1 // Default for unknown types
	}

	totalTokens := int(float64(baseTokens) * multiplier)

	// Cap the estimation (normal resolution images: 800-1500 tokens)
	if totalTokens > 3000 {
		totalTokens = 3000
	}

	// Estimate input/output split (typically 80% input, 20% output for OCR)
	inputTokens := int(float64(totalTokens) * 0.8)
	outputTokens := totalTokens - inputTokens

	return inputTokens, outputTokens, totalTokens
}

// EstimateTextTokenUsage provides token usage estimation for text-based AI operations
// This is a simple word-based estimation and should be replaced with proper tokenizer
func EstimateTextTokenUsage(text string) int {
	// Rough estimation: 1 token ≈ 0.75 words for English text
	words := len(strings.Fields(text))
	return int(float64(words) / 0.75)
}