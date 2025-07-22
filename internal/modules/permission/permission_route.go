package permission

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterPermissionRoutes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler, rateLimitMiddleware fiber.Handler, permissionMiddleware func(string, string) fiber.Handler) {
	v1 := app.Group("/v1")
	
	// Protected routes with auth, rate limit, and permission checking
	v1.Post("/permissions", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("permissions", "manage"), 
		handler.CreatePermission)
	v1.Get("/permissions", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("permissions", "manage"), 
		handler.GetPermissions)
	v1.Get("/permissions/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("permissions", "manage"), 
		handler.GetPermission)
	v1.Delete("/permissions/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("permissions", "manage"), 
		handler.DeletePermission)
}