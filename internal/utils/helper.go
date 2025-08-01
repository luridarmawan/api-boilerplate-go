package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"regexp"
	"strings"
)

// Helper function to create int16 pointer
func Int16Ptr(v int16) *int16 {
	return &v
}

// Helper function to create uint pointer
func UintPtr(v uint) *uint {
	return &v
}

// FormatMessage formats error messages by replacing vendor names and removing URLs
// while preserving the original message format and case
func FormatMessage(message string) string {
	result := message

	// Replace vendor names with "vendor" using case-insensitive regex
	openaiPattern := regexp.MustCompile(`(?i)\bopenai\b`)
	result = openaiPattern.ReplaceAllString(result, "vendor")

	anthropicPattern := regexp.MustCompile(`(?i)\banthropic\b`)
	result = anthropicPattern.ReplaceAllString(result, "vendor")

	geminiPattern := regexp.MustCompile(`(?i)\bgemini\b`)
	result = geminiPattern.ReplaceAllString(result, "vendor")

	// Remove URLs using regex pattern
	// This pattern matches http://, https://, ftp://, and www. URLs
	urlPattern := regexp.MustCompile(`(?i)\b(?:https?://|ftp://|www\.)[^\s<>"{}|\\^` + "`" + `\[\]]*`)
	result = urlPattern.ReplaceAllString(result, "")

	// Remove extra spaces that might be left after URL removal
	spacePattern := regexp.MustCompile(`\s+`)
	result = spacePattern.ReplaceAllString(result, " ")

	// Trim leading and trailing spaces
	result = strings.TrimSpace(result)

	return result
}

// GenerateAPIKey generates a secure random API key
func GenerateAPIKey() string {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	rand.Read(bytes)

	// Encode to base64 and remove padding
	apiKey := base64.URLEncoding.EncodeToString(bytes)
	apiKey = strings.TrimRight(apiKey, "=")

	return "sk-" + apiKey
}

// GetTableName returns the table name with prefix from environment variable
func GetTableName(tableName string) string {
	prefix := os.Getenv("DB_TABLE_PREFIX")
	if prefix == "" {
		return tableName
	}
	return prefix + tableName
}

// GetJoinTableName returns the join table name with prefix for many-to-many relationships
func GetJoinTableName(joinTableName string) string {
	prefix := os.Getenv("DB_TABLE_PREFIX")
	if prefix == "" {
		return joinTableName
	}
	return prefix + joinTableName
}
