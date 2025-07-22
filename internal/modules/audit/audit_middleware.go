package audit

import (
	"encoding/json"
	"strings"
	"time"

	"apiserver/internal/modules/access"

	"github.com/gofiber/fiber/v2"
)

// NewAuditMiddleware creates a middleware that logs all API requests and responses
func NewAuditMiddleware(auditRepo Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip health check and docs endpoints
		if strings.HasPrefix(c.Path(), "/health") || 
		   strings.HasPrefix(c.Path(), "/docs") ||
		   strings.HasPrefix(c.Path(), "/swagger") {
			return c.Next()
		}

		start := time.Now()

		// Capture request body
		var requestBody string
		if c.Body() != nil {
			requestBody = string(c.Body())
		}

		// Capture request headers (excluding sensitive ones)
		requestHeaders := make(map[string]string)
		c.Request().Header.VisitAll(func(key, value []byte) {
			keyStr := string(key)
			// Skip sensitive headers
			if !isSensitiveHeader(keyStr) {
				requestHeaders[keyStr] = string(value)
			}
		})
		requestHeadersJSON, _ := json.Marshal(requestHeaders)

		// Process request
		err := c.Next()

		// Calculate response time
		responseTime := time.Since(start).Milliseconds()

		// Get user information if available
		var userID *string
		var userEmail, apiKey string
		
		if user, ok := c.Locals("user").(*access.User); ok {
			userID = &user.ID
			userEmail = user.Email
		}
		
		// Try to get API key from Authorization header
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			// Mask the API key for security (show only first 8 chars)
			if len(apiKey) > 8 {
				apiKey = apiKey[:8] + "****"
			}
		}

		// Get response body from Fiber context
		responseBody := string(c.Response().Body())
		if len(responseBody) > 10000 { // Limit to 10KB
			responseBody = responseBody[:10000] + "... [truncated]"
		}

		// Create audit log entry
		auditLog := &AuditLog{
			UserID:         userID,
			UserEmail:      userEmail,
			APIKey:         apiKey,
			Method:         c.Method(),
			Path:           c.Path(),
			StatusCode:     c.Response().StatusCode(),
			RequestHeaders: string(requestHeadersJSON),
			RequestBody:    requestBody,
			ResponseBody:   responseBody,
			ResponseTime:   responseTime,
			IPAddress:      c.IP(),
			UserAgent:      c.Get("User-Agent"),
			StatusID:       func() *int16 { v := int16(0); return &v }(), // Active
		}

		// Save audit log asynchronously to avoid blocking the response
		go func() {
			if saveErr := auditRepo.CreateAuditLog(auditLog); saveErr != nil {
				// Log error but don't fail the request
				// In production, you might want to use a proper logger here
				// log.Printf("Failed to save audit log: %v", saveErr)
			}
		}()

		return err
	}
}

// isSensitiveHeader checks if a header contains sensitive information
func isSensitiveHeader(header string) bool {
	sensitiveHeaders := []string{
		"authorization",
		"cookie",
		"set-cookie",
		"x-api-key",
		"x-auth-token",
	}
	
	headerLower := strings.ToLower(header)
	for _, sensitive := range sensitiveHeaders {
		if headerLower == sensitive {
			return true
		}
	}
	return false
}