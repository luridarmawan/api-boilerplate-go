package configuration

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterConfigurationRoutes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler, rateLimitMiddleware fiber.Handler, permissionMiddleware func(string, string) fiber.Handler) {
	v1 := app.Group("/v1")
	
	// Protected routes with auth, rate limit, and permission checking
	v1.Post("/configurations", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("configurations", "create"), 
		handler.CreateConfiguration)
	v1.Get("/configurations", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("configurations", "read"), 
		handler.GetConfigurations)
	v1.Get("/configurations/deleted", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("configurations", "read"), 
		handler.GetDeletedConfigurations)
	v1.Get("/configurations/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("configurations", "read"), 
		handler.GetConfiguration)
	v1.Put("/configurations/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("configurations", "update"), 
		handler.UpdateConfiguration)
	v1.Delete("/configurations/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("configurations", "delete"), 
		handler.SoftDeleteConfiguration)
	v1.Post("/configurations/:id/restore", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("configurations", "update"), 
		handler.RestoreConfiguration)
}
