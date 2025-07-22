package middleware

import (
	"strconv"
	"sync"
	"time"

	"apiserver/internal/types"

	"github.com/gofiber/fiber/v2"
)

// RateLimiter struct to store rate limit data
type RateLimiter struct {
	sync.RWMutex
	requests     map[string][]time.Time // Map of API key to request timestamps
	defaultLimit int                    // Default requests per minute
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(defaultLimit int) *RateLimiter {
	return &RateLimiter{
		requests:     make(map[string][]time.Time),
		defaultLimit: defaultLimit,
	}
}

// RateLimitMiddleware creates a middleware that limits requests per minute based on API key
func RateLimitMiddleware(limiter *RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context (set by auth middleware)
		user, ok := c.Locals("user").(types.User)
		if !ok {
			// If no user, let auth middleware handle it
			return c.Next()
		}

		// Extract API key from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		// Remove "Bearer " prefix to get the actual API key
		apiKey := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			apiKey = authHeader[7:]
		}

		// Get rate limit for this user
		rateLimit := limiter.defaultLimit
		if userWithRateLimit, ok := user.(interface{ GetRateLimit() int }); ok {
			rateLimit = userWithRateLimit.GetRateLimit()
		}

		// Check rate limit
		limiter.Lock()
		defer limiter.Unlock()

		// Clean up old requests (older than 1 minute)
		now := time.Now()
		oneMinuteAgo := now.Add(-time.Minute)

		// Get or initialize request timestamps for this API key
		timestamps, exists := limiter.requests[apiKey]
		if !exists {
			timestamps = []time.Time{}
		}

		// Filter out old timestamps
		newTimestamps := []time.Time{}
		for _, t := range timestamps {
			if t.After(oneMinuteAgo) {
				newTimestamps = append(newTimestamps, t)
			}
		}

		// Check if rate limit exceeded
		if len(newTimestamps) >= rateLimit {
			// Add rate limit headers
			c.Set("X-RateLimit-Limit", strconv.Itoa(rateLimit))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", strconv.FormatInt(oneMinuteAgo.Add(time.Minute).Unix(), 10))
			
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Rate limit exceeded. Try again later.",
			})
		}

		// Add current request timestamp
		newTimestamps = append(newTimestamps, now)
		limiter.requests[apiKey] = newTimestamps

		// Add rate limit headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(rateLimit))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(rateLimit-len(newTimestamps)))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(oneMinuteAgo.Add(time.Minute).Unix(), 10))

		return c.Next()
	}
}