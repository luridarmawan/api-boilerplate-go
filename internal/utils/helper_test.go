// USAGE
//   go test ./internal/utils -v -run TestFormatMessage

package utils

import (
	"testing"
)

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Replace OpenAI with vendor",
			input:    "OpenAI API error occurred",
			expected: "vendor API error occurred",
		},
		{
			name:     "Replace Anthropic with vendor",
			input:    "Anthropic service is unavailable",
			expected: "vendor service is unavailable",
		},
		{
			name:     "Replace Gemini with vendor",
			input:    "Gemini model not found",
			expected: "vendor model not found",
		},
		{
			name:     "Case insensitive replacement",
			input:    "OPENAI and anthropic and GeMiNi",
			expected: "vendor and vendor and vendor",
		},
		{
			name:     "Remove HTTP URL",
			input:    "Error from https://api.openai.com/v1/chat/completions",
			expected: "Error from",
		},
		{
			name:     "Remove HTTPS URL",
			input:    "Failed to connect to http://localhost:8080/api",
			expected: "Failed to connect to",
		},
		{
			name:     "Remove www URL",
			input:    "Visit www.example.com for more info",
			expected: "Visit for more info",
		},
		{
			name:     "Complex message with vendor and URL",
			input:    "OpenAI API error: failed to connect to https://api.openai.com/v1/completions",
			expected: "vendor API error: failed to connect to",
		},
		{
			name:     "Multiple URLs in message",
			input:    "Check https://docs.openai.com or http://help.anthropic.com",
			expected: "Check or",
		},
		{
			name:     "No changes needed",
			input:    "This is a normal error message",
			expected: "This is a normal error message",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "Word boundaries test",
			input:    "reopenai and anthropics and geminis should not be replaced",
			expected: "reopenai and anthropics and geminis should not be replaced",
		},
		{
			name:     "Mixed case preservation",
			input:    "OpenAI Error: ANTHROPIC failed, Gemini works",
			expected: "vendor Error: vendor failed, vendor works",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatMessage(tt.input)
			if result != tt.expected {
				t.Errorf("FormatMessage(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}