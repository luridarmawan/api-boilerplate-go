package access

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterAccessRoutes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler, rateLimitMiddleware fiber.Handler, permissionMiddleware func(string, string) fiber.Handler) {
	v1 := app.Group("/v1")
	
	// Public route for creating access
	v1.Post("/access",
		authMiddleware,
		rateLimitMiddleware,
		permissionMiddleware("access", "manage"),
		handler.CreateAccess)

	// Protected routes with auth, rate limit, and permission checking
	v1.Get("/profile",
		authMiddleware,
		rateLimitMiddleware,
		permissionMiddleware("profile", "read"),
		handler.GetProfile)

	// API key expiration management routes
	v1.Put("/access/:id/expired-date",
		authMiddleware,
		rateLimitMiddleware,
		permissionMiddleware("access", "manage"),
		handler.UpdateExpiredDate)

	v1.Delete("/access/:id/expired-date",
		authMiddleware,
		rateLimitMiddleware,
		permissionMiddleware("access", "manage"),
		handler.RemoveExpiredDate)

	// API key rate limit management routes
	v1.Put("/access/:id/rate-limit",
		authMiddleware,
		rateLimitMiddleware,
		permissionMiddleware("access", "manage"),
		handler.UpdateRateLimit)
}