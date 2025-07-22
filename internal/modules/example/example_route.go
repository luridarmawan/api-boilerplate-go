package example

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterExampleRoutes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler, rateLimitMiddleware fiber.Handler, permissionMiddleware func(string, string) fiber.Handler) {
	v1 := app.Group("/v1")
	
	// Protected routes with auth, rate limit, and permission checking
	v1.Post("/examples", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("examples", "create"), 
		handler.CreateExample)
	v1.Get("/examples", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("examples", "read"), 
		handler.GetExamples)
	v1.Get("/examples/deleted", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("examples", "read"), 
		handler.GetDeletedExamples)
	v1.Get("/examples/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("examples", "read"), 
		handler.GetExample)
	v1.Put("/examples/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("examples", "update"), 
		handler.UpdateExample)
	v1.Delete("/examples/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("examples", "delete"), 
		handler.SoftDeleteExample)
	v1.Post("/examples/:id/restore", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("examples", "update"), 
		handler.RestoreExample)
}