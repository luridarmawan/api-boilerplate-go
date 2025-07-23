package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Templates for each file
const modelTemplate = `package {{.Package}}

import (
	"time"

	"apiserver/internal/utils"
	"gorm.io/gorm"
)

type {{.Module}} struct {
	ID          string         ` + "`json:\"id\" gorm:\"type:uuid;primaryKey\"`" + `
	Name        string         ` + "`json:\"name\" gorm:\"not null\"`" + `
	Description string         ` + "`json:\"description\"`" + `
	CreatedAt   time.Time      ` + "`json:\"created_at\"`" + `
	UpdatedAt   time.Time      ` + "`json:\"updated_at\"`" + `
	DeletedAt   gorm.DeletedAt ` + "`json:\"-\" gorm:\"index\"`" + `
	StatusID    *int16         ` + "`json:\"status_id\" gorm:\"type:smallint;not null;default:1;index\"`" + `
}

type Create{{.Module}}Request struct {
	Name        string ` + "`json:\"name\" validate:\"required\"`" + `
	Description string ` + "`json:\"description\"`" + `
}

func ({{.Module}}) TableName() string {
	return "{{.TableName}}"
}

// BeforeCreate hook to generate UUIDv7 before creating a new {{.LowerModule}}
func (e *{{.Module}}) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = utils.GenerateUUIDv7()
	}
	return nil
}
`

const repositoryTemplate = `package {{.Package}}

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create{{.Module}}({{.LowerModule}} *{{.Module}}) error
	GetAll{{.ModulePlural}}() ([]{{.Module}}, error)
	Get{{.Module}}ByID(id string) (*{{.Module}}, error)
	Update{{.Module}}({{.LowerModule}} *{{.Module}}) error
	SoftDelete{{.Module}}(id string) error
	Restore{{.Module}}(id string) error
	GetDeleted{{.ModulePlural}}() ([]{{.Module}}, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create{{.Module}}({{.LowerModule}} *{{.Module}}) error {
	return r.db.Create({{.LowerModule}}).Error
}

func (r *repository) GetAll{{.ModulePlural}}() ([]{{.Module}}, error) {
	var {{.LowerModulePlural}} []{{.Module}}
	err := r.db.Where("status_id = ?", 0).Find(&{{.LowerModulePlural}}).Error
	return {{.LowerModulePlural}}, err
}

func (r *repository) Get{{.Module}}ByID(id string) (*{{.Module}}, error) {
	var {{.LowerModule}} {{.Module}}
	err := r.db.Where("id = ? AND status_id = ?", id, 0).First(&{{.LowerModule}}).Error
	if err != nil {
		return nil, err
	}
	return &{{.LowerModule}}, nil
}

func (r *repository) Update{{.Module}}({{.LowerModule}} *{{.Module}}) error {
	return r.db.Save({{.LowerModule}}).Error
}

func (r *repository) SoftDelete{{.Module}}(id string) error {
	return r.db.Model(&{{.Module}}{}).Where("id = ?", id).Update("status_id", 1).Error
}

func (r *repository) Restore{{.Module}}(id string) error {
	return r.db.Model(&{{.Module}}{}).Where("id = ?", id).Update("status_id", 0).Error
}

func (r *repository) GetDeleted{{.ModulePlural}}() ([]{{.Module}}, error) {
	var {{.LowerModulePlural}} []{{.Module}}
	err := r.db.Where("status_id = ?", 1).Find(&{{.LowerModulePlural}}).Error
	return {{.LowerModulePlural}}, err
}
`

const handlerTemplate = `package {{.Package}}

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// Create{{.Module}} godoc
// @Summary Create a new {{.LowerModule}}
// @Description Create a new {{.LowerModule}} with name and description
// @Tags {{.Module}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param {{.LowerModule}} body Create{{.Module}}Request true "{{.Module}} data"
// @Success 201 {object} {{.Module}}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/{{.LowerModulePlural}} [post]
func (h *Handler) Create{{.Module}}(c *fiber.Ctx) error {
	var req Create{{.Module}}Request
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Simple validation
	if strings.TrimSpace(req.Name) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Name is required",
		})
	}

	{{.LowerModule}} := &{{.Module}}{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.repo.Create{{.Module}}({{.LowerModule}}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create {{.LowerModule}}",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   {{.LowerModule}},
	})
}

// Get{{.ModulePlural}} godoc
// @Summary Get all {{.LowerModulePlural}}
// @Description Get list of all {{.LowerModulePlural}}
// @Tags {{.Module}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} {{.Module}}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/{{.LowerModulePlural}} [get]
func (h *Handler) Get{{.ModulePlural}}(c *fiber.Ctx) error {
	{{.LowerModulePlural}}, err := h.repo.GetAll{{.ModulePlural}}()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch {{.LowerModulePlural}}",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   {{.LowerModulePlural}},
	})
}

// Get{{.Module}} godoc
// @Summary Get {{.LowerModule}} by ID
// @Description Get a specific {{.LowerModule}} by its ID
// @Tags {{.Module}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "{{.Module}} ID"
// @Success 200 {object} {{.Module}}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/{{.LowerModulePlural}}/{id} [get]
func (h *Handler) Get{{.Module}}(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid {{.LowerModule}} ID",
		})
	}

	{{.LowerModule}}, err := h.repo.Get{{.Module}}ByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "{{.Module}} not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   {{.LowerModule}},
	})
}

// Update{{.Module}} godoc
// @Summary Update {{.LowerModule}}
// @Description Update an existing {{.LowerModule}}
// @Tags {{.Module}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "{{.Module}} ID"
// @Param {{.LowerModule}} body Create{{.Module}}Request true "{{.Module}} data"
// @Success 200 {object} {{.Module}}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/{{.LowerModulePlural}}/{id} [put]
func (h *Handler) Update{{.Module}}(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid {{.LowerModule}} ID",
		})
	}

	// Check if {{.LowerModule}} exists
	{{.LowerModule}}, err := h.repo.Get{{.Module}}ByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "{{.Module}} not found",
		})
	}

	var req Create{{.Module}}Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Simple validation
	if strings.TrimSpace(req.Name) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Name is required",
		})
	}

	// Update {{.LowerModule}}
	{{.LowerModule}}.Name = req.Name
	{{.LowerModule}}.Description = req.Description

	if err := h.repo.Update{{.Module}}({{.LowerModule}}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update {{.LowerModule}}",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   {{.LowerModule}},
	})
}

// SoftDelete{{.Module}} godoc
// @Summary Soft delete {{.LowerModule}}
// @Description Soft delete a {{.LowerModule}} (set status_id to 0)
// @Tags {{.Module}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "{{.Module}} ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/{{.LowerModulePlural}}/{id} [delete]
func (h *Handler) SoftDelete{{.Module}}(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid {{.LowerModule}} ID",
		})
	}

	// Check if {{.LowerModule}} exists
	_, err := h.repo.Get{{.Module}}ByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "{{.Module}} not found",
		})
	}

	if err := h.repo.SoftDelete{{.Module}}(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete {{.LowerModule}}",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "{{.Module}} deleted successfully",
	})
}

// Restore{{.Module}} godoc
// @Summary Restore {{.LowerModule}}
// @Description Restore a soft deleted {{.LowerModule}} (set status_id to 1)
// @Tags {{.Module}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "{{.Module}} ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/{{.LowerModulePlural}}/{id}/restore [post]
func (h *Handler) Restore{{.Module}}(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid {{.LowerModule}} ID",
		})
	}

	if err := h.repo.Restore{{.Module}}(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to restore {{.LowerModule}}",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "{{.Module}} restored successfully",
	})
}

// GetDeleted{{.ModulePlural}} godoc
// @Summary Get deleted {{.LowerModulePlural}}
// @Description Get list of all soft deleted {{.LowerModulePlural}}
// @Tags {{.Module}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} {{.Module}}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/{{.LowerModulePlural}}/deleted [get]
func (h *Handler) GetDeleted{{.ModulePlural}}(c *fiber.Ctx) error {
	{{.LowerModulePlural}}, err := h.repo.GetDeleted{{.ModulePlural}}()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch deleted {{.LowerModulePlural}}",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   {{.LowerModulePlural}},
	})
}
`

const routeTemplate = `package {{.Package}}

import (
	"github.com/gofiber/fiber/v2"
)

func Register{{.Module}}Routes(app *fiber.App, handler *Handler, authMiddleware fiber.Handler, rateLimitMiddleware fiber.Handler, permissionMiddleware func(string, string) fiber.Handler) {
	v1 := app.Group("/v1")
	
	// Protected routes with auth, rate limit, and permission checking
	v1.Post("/{{.LowerModulePlural}}", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("{{.LowerModulePlural}}", "create"), 
		handler.Create{{.Module}})
	v1.Get("/{{.LowerModulePlural}}", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("{{.LowerModulePlural}}", "read"), 
		handler.Get{{.ModulePlural}})
	v1.Get("/{{.LowerModulePlural}}/deleted", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("{{.LowerModulePlural}}", "read"), 
		handler.GetDeleted{{.ModulePlural}})
	v1.Get("/{{.LowerModulePlural}}/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("{{.LowerModulePlural}}", "read"), 
		handler.Get{{.Module}})
	v1.Put("/{{.LowerModulePlural}}/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("{{.LowerModulePlural}}", "update"), 
		handler.Update{{.Module}})
	v1.Delete("/{{.LowerModulePlural}}/:id", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("{{.LowerModulePlural}}", "delete"), 
		handler.SoftDelete{{.Module}})
	v1.Post("/{{.LowerModulePlural}}/:id/restore", 
		authMiddleware, 
		rateLimitMiddleware,
		permissionMiddleware("{{.LowerModulePlural}}", "update"), 
		handler.Restore{{.Module}})
}
`

const mainUpdateTemplate = `
// Add this to main.go in the appropriate section:

// Initialize {{.LowerModule}} module
{{.LowerModule}}Repo := {{.Package}}.NewRepository(db)
{{.LowerModule}}Handler := {{.Package}}.NewHandler({{.LowerModule}}Repo)
{{.Package}}.Register{{.Module}}Routes(app, {{.LowerModule}}Handler, authMiddleware, rateLimitMiddleware, middleware.RequirePermission)
`

// TemplateData holds the data for template rendering
type TemplateData struct {
	Package           string
	Module            string
	ModulePlural      string
	LowerModule       string
	LowerModulePlural string
	TableName         string
}

// Pluralize adds 's' to the end of a word, with special cases
func pluralize(word string) string {
	// Special cases
	switch {
	case strings.HasSuffix(word, "y"):
		return word[:len(word)-1] + "ies"
	case strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || 
		 strings.HasSuffix(word, "z") || strings.HasSuffix(word, "ch") || 
		 strings.HasSuffix(word, "sh"):
		return word + "es"
	default:
		return word + "s"
	}
}

// capitalize capitalizes the first letter of a string
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// uncapitalize makes the first letter of a string lowercase
func uncapitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func main() {
	// Check if module name is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <module-name>")
		fmt.Println("Example: go run main.go product")
		os.Exit(1)
	}

	// Get module name from command line
	moduleName := os.Args[1]
	
	// Prepare template data
	data := TemplateData{
		Package:           moduleName,
		Module:            capitalize(moduleName),
		ModulePlural:      capitalize(pluralize(moduleName)),
		LowerModule:       uncapitalize(moduleName),
		LowerModulePlural: uncapitalize(pluralize(moduleName)),
		TableName:         pluralize(moduleName),
	}

	// Create module directory
	moduleDir := filepath.Join("internal", "modules", moduleName)
	err := os.MkdirAll(moduleDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Create model file
	createFileFromTemplate(filepath.Join(moduleDir, fmt.Sprintf("%s_model.go", moduleName)), modelTemplate, data)
	
	// Create repository file
	createFileFromTemplate(filepath.Join(moduleDir, fmt.Sprintf("%s_repository.go", moduleName)), repositoryTemplate, data)
	
	// Create handler file
	createFileFromTemplate(filepath.Join(moduleDir, fmt.Sprintf("%s_handler.go", moduleName)), handlerTemplate, data)
	
	// Create route file
	createFileFromTemplate(filepath.Join(moduleDir, fmt.Sprintf("%s_route.go", moduleName)), routeTemplate, data)

	// Print instructions for main.go update
	fmt.Println("Module created successfully!")
	fmt.Println("\nTo complete the setup, add the following code to main.go:")
	
	// Execute the main update template
	tmpl, err := template.New("mainUpdate").Parse(mainUpdateTemplate)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}
	
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("\nDon't forget to add the module to the AutoMigrate function in main.go:")
	fmt.Printf("err := db.AutoMigrate(&%s.%s{}, ...)\n", moduleName, data.Module)
}

// createFileFromTemplate creates a file from a template
func createFileFromTemplate(filePath, templateContent string, data TemplateData) {
	// Parse template
	tmpl, err := template.New(filepath.Base(filePath)).Parse(templateContent)
	if err != nil {
		fmt.Printf("Error parsing template for %s: %v\n", filePath, err)
		os.Exit(1)
	}
	
	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()
	
	// Execute template
	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Printf("Error executing template for %s: %v\n", filePath, err)
		os.Exit(1)
	}
	
	fmt.Printf("Created %s\n", filePath)
}