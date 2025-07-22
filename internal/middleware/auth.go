package middleware

import (
	"strings"

	"apiserver/internal/types"

	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(authRepo types.AuthRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Authorization header is required",
			})
		}

		// Check if it's Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid authorization format. Use Bearer token",
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Token is required",
			})
		}

		// Validate token against database
		user, err := authRepo.FindByAPIKey(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid or expired token",
			})
		}

		// Store user information in context
		c.Locals("userID", user.GetID())
		c.Locals("user", user)

		return c.Next()
	}
}