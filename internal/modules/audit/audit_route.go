package audit

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterAuditRoutes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler, rateLimitMiddleware fiber.Handler, requirePermission func(string, string) fiber.Handler) {
	v1 := app.Group("/v1")
	
	// Protected routes with auth, rate limit, and permission checking
	v1.Get("/audit-logs", 
		authMiddleware, 
		rateLimitMiddleware,
		requirePermission("audit", "read"), 
		handler.GetAuditLogs)
	v1.Get("/audit-logs/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		requirePermission("audit", "read"), 
		handler.GetAuditLog)
	v1.Delete("/audit-logs/cleanup", 
		authMiddleware, 
		rateLimitMiddleware,
		requirePermission("audit", "manage"), 
		handler.DeleteOldLogs)
}