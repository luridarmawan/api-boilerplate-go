package middleware

import (
	"apiserver/internal/types"

	"github.com/gofiber/fiber/v2"
)

// RequirePermission creates a middleware that checks if the user has the required permission
func RequirePermission(resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context (set by auth middleware)
		user, ok := c.Locals("user").(types.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "User not authenticated",
			})
		}

		// Check if user has a group
		group := user.GetGroup()
		if group == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Access denied: No group assigned",
			})
		}

		// Check if user's group has the required permission
		hasPermission := false
		for _, permission := range group.Permissions {
			if permission.Resource == resource && permission.Action == action {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Access denied: Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireAnyPermission creates a middleware that checks if the user has any of the required permissions
func RequireAnyPermission(permissions []struct{ Resource, Action string }) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context (set by auth middleware)
		user, ok := c.Locals("user").(types.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "User not authenticated",
			})
		}

		// Check if user has a group
		group := user.GetGroup()
		if group == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Access denied: No group assigned",
			})
		}

		// Check if user's group has any of the required permissions
		hasPermission := false
		for _, userPerm := range group.Permissions {
			for _, reqPerm := range permissions {
				if userPerm.Resource == reqPerm.Resource && userPerm.Action == reqPerm.Action {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Access denied: Insufficient permissions",
			})
		}

		return c.Next()
	}
}