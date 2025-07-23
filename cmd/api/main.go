package main

import (
	"flag"
	"log"
	"os"

	"apiserver/configs"
	"apiserver/docs"
	"apiserver/internal/database"
	"apiserver/internal/middleware"
	"apiserver/internal/modules/access"
	"apiserver/internal/modules/audit"
	"apiserver/internal/modules/group"
	"apiserver/internal/modules/permission"
	"apiserver/internal/modules/example"
	"apiserver/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"path/filepath"
)

// @title My API
// @version 1.0
// @description This is a modular REST API built with Go Fiber
// @contact.name API Support
// --@host localhost:3000
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your API token
func main() {
	// Parse command line flags
	seedFlag := flag.Bool("seed", false, "Run database seeding")
	flag.Parse()

	// Load configuration
	config := configs.LoadConfig()

	// Initialize docs
	docs.SwaggerInfo.Title = config.APIName
	docs.SwaggerInfo.Description = config.APIDescription
	docs.SwaggerInfo.Version = config.APIVersion
	docs.SwaggerInfo.Host = config.BaseURL
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Initialize database
	database.InitDatabase(config)

	// Auto-migrate models
	db := database.GetDB()
	err := db.AutoMigrate(&access.User{}, &example.Example{}, &permission.Permission{}, &group.Group{}, &audit.AuditLog{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed database with test data (if flag is provided or first run)
	if *seedFlag {
		database.SeedDatabase(db)
		log.Println("Seeding completed. Exiting...")
		os.Exit(0)
	}

	// Initialize repositories
	accessRepo := access.NewRepository(db)
	exampleRepo := example.NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	groupRepo := group.NewRepository(db)
	auditRepo := audit.NewRepository(db)

	// Initialize handlers
	accessHandler := access.NewHandler(accessRepo)
	exampleHandler := example.NewHandler(exampleRepo)
	permissionHandler := permission.NewHandler(permissionRepo)
	groupHandler := group.NewHandler(groupRepo)
	auditHandler := audit.NewHandler(auditRepo)

	// Initialize middleware with auth repository wrapper
	authRepo := access.NewAuthRepository(accessRepo)
	authMiddleware := middleware.NewAuthMiddleware(authRepo)
	auditMiddleware := audit.NewAuditMiddleware(auditRepo)

	// Initialize rate limiter middleware (default: 120 requests per minute)
	rateLimiter := middleware.NewRateLimiter(120)
	rateLimitMiddleware := middleware.RateLimitMiddleware(rateLimiter)

	// Initialize your custom module here

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		},
	})

	// Static Handler
	app.Static("/static", "./static")
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendFile("./static/favicon.ico")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return utils.Output(c, "OK")
	})

	app.Get("/rapidocs", func(c *fiber.Ctx) error {
		data, _ := os.ReadFile("./docs/rapidoc.html")
		html := string(data)
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	app.Get("/scalar", func(c *fiber.Ctx) error {
		data, _ := os.ReadFile("./docs/scalar.html")
		html := string(data)
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	app.Get("/docs", func(c *fiber.Ctx) error {
		data, _ := os.ReadFile("./docs/swagger.html")
		html := string(data)
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	app.Get("/docs/openapi.json", func(c *fiber.Ctx) error {
		baseDir, _ := filepath.Abs(".")
		jsonPath := filepath.Join(baseDir, "docs", "openapi.json")
		data, err := os.ReadFile(jsonPath)
		if err != nil {
			return utils.Output(c, "Failed to load OpenAPI spec", false, 500)
		}
		c.Set("Content-Type", "application/json")
		return c.Send(data)
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(auditMiddleware) // Add audit logging middleware

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
		})
	})

	// Show Version
	app.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app":           config.APIName,
			"build_version": config.Version,
			"build_date":    config.BuildDate,
		})
	})

	// Register routes with auth, rate limit, and permission middleware
	access.RegisterAccessRoutes(app, accessHandler, authMiddleware, rateLimitMiddleware, middleware.RequirePermission)
	example.RegisterExampleRoutes(app, exampleHandler, authMiddleware, rateLimitMiddleware, middleware.RequirePermission)
	permission.RegisterPermissionRoutes(app, permissionHandler, authMiddleware, rateLimitMiddleware, middleware.RequirePermission)
	group.RegisterGroupRoutes(app, groupHandler, authMiddleware, rateLimitMiddleware, middleware.RequirePermission)
	audit.RegisterAuditRoutes(app, auditHandler, authMiddleware, rateLimitMiddleware, middleware.RequirePermission)

	// Register your module route here

	// Start server
	log.Printf("Server starting on port %s", config.ServerPort)
	log.Fatal(app.Listen(":" + config.ServerPort))
}