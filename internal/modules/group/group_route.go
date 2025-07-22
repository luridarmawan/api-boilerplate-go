package group

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterGroupRoutes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler, rateLimitMiddleware fiber.Handler, permissionMiddleware func(string, string) fiber.Handler) {
	v1 := app.Group("/v1")
	
	// Protected routes with auth, rate limit, and permission checking
	v1.Post("/groups", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("groups", "manage"), 
		handler.CreateGroup)
	v1.Get("/groups", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("groups", "manage"), 
		handler.GetGroups)
	v1.Get("/groups/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("groups", "manage"), 
		handler.GetGroup)
	v1.Put("/groups/:id/permissions", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("groups", "manage"), 
		handler.UpdateGroupPermissions)
	v1.Delete("/groups/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("groups", "manage"), 
		handler.DeleteGroup)
}