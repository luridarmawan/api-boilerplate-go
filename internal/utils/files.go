package utils

import (
	"os"
	"strings"
)

// ReadFile reads the content of a file and returns it as a string
// If the file doesn't exist or there's an error, it returns an empty string
func ReadFile(filepath string) string {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return ""
	}
	
	// Return the content as string, trimming any trailing whitespace
	return strings.TrimSpace(string(content))
}