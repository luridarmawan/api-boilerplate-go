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

const permissionScriptTemplate = `package main

import (
	"fmt"
	"log"
	"time"

	"apiserver/configs"
	"apiserver/internal/database"
	"apiserver/internal/modules/group"
	"apiserver/internal/modules/permission"
)

// Helper function to create int16 pointer
func int16Ptr(v int16) *int16 {
	return &v
}

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Initialize database
	database.InitDatabase(config)
	db := database.GetDB()

	// Create {{.LowerModule}} permissions
	{{.LowerModule}}Permissions := []permission.Permission{
		{
			Name:        "Create {{.ModulePlural}}",
			Description: "Permission to create new {{.LowerModulePlural}}",
			Resource:    "{{.LowerModulePlural}}",
			Action:      "create",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Read {{.ModulePlural}}",
			Description: "Permission to read {{.LowerModulePlural}}",
			Resource:    "{{.LowerModulePlural}}",
			Action:      "read",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Update {{.ModulePlural}}",
			Description: "Permission to update {{.LowerModulePlural}}",
			Resource:    "{{.LowerModulePlural}}",
			Action:      "update",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Delete {{.ModulePlural}}",
			Description: "Permission to delete {{.LowerModulePlural}}",
			Resource:    "{{.LowerModulePlural}}",
			Action:      "delete",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Add permissions to database
	createdPermissions := make(map[string]uint)
	for _, p := range {{.LowerModule}}Permissions {
		var existingPermission permission.Permission
		result := db.Where("name = ?", p.Name).First(&existingPermission)

		if result.Error != nil {
			// Permission doesn't exist, create it
			if err := db.Create(&p).Error; err != nil {
				log.Printf("Failed to create permission %s: %v", p.Name, err)
			} else {
				log.Printf("Created permission: %s with ID: %d", p.Name, p.ID)
				createdPermissions[p.Action] = p.ID
			}
		} else {
			log.Printf("Permission %s already exists with ID: %d", p.Name, existingPermission.ID)
			createdPermissions[p.Action] = existingPermission.ID
		}
	}

	// Get groups
	var adminGroup, editorGroup, viewerGroup, generalClientGroup group.Group
	db.Where("name = ?", "Admin").First(&adminGroup)
	db.Where("name = ?", "Editor").First(&editorGroup)
	db.Where("name = ?", "Viewer").First(&viewerGroup)
	db.Where("name = ?", "General client").First(&generalClientGroup)

	// Assign permissions to groups
	if adminGroup.ID > 0 {
		// Admin gets all permissions
		for _, permID := range createdPermissions {
			var perm permission.Permission
			if err := db.First(&perm, permID).Error; err == nil {
				if err := db.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", adminGroup.ID, permID).Error; err != nil {
					log.Printf("Failed to assign permission %d to Admin group: %v", permID, err)
				} else {
					log.Printf("Assigned permission %d to Admin group", permID)
				}
			}
		}
	}

	if editorGroup.ID > 0 {
		// Editor gets create, read, update
		for _, action := range []string{"create", "read", "update"} {
			if permID, exists := createdPermissions[action]; exists {
				if err := db.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", editorGroup.ID, permID).Error; err != nil {
					log.Printf("Failed to assign permission %d to Editor group: %v", permID, err)
				} else {
					log.Printf("Assigned permission %d to Editor group", permID)
				}
			}
		}
	}

	if viewerGroup.ID > 0 {
		// Viewer gets read only
		if permID, exists := createdPermissions["read"]; exists {
			if err := db.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", viewerGroup.ID, permID).Error; err != nil {
				log.Printf("Failed to assign permission %d to Viewer group: %v", permID, err)
			} else {
				log.Printf("Assigned permission %d to Viewer group", permID)
			}
		}
	}

	if generalClientGroup.ID > 0 {
		// General client gets read only
		if permID, exists := createdPermissions["read"]; exists {
			if err := db.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", generalClientGroup.ID, permID).Error; err != nil {
				log.Printf("Failed to assign permission %d to General client group: %v", permID, err)
			} else {
				log.Printf("Assigned permission %d to General client group", permID)
			}
		}
	}

	fmt.Println("{{.Module}} permissions have been added and assigned to groups!")
	fmt.Println("\nPermissions created:")
	fmt.Println("- Create {{.ModulePlural}} (Admin, Editor)")
	fmt.Println("- Read {{.ModulePlural}} (Admin, Editor, Viewer, General client)")
	fmt.Println("- Update {{.ModulePlural}} (Admin, Editor)")
	fmt.Println("- Delete {{.ModulePlural}} (Admin only)")
}
`

const permissionSQLTemplate = `-- Add {{.LowerModule}} permissions
INSERT INTO permissions (name, description, resource, action, status_id, created_at, updated_at)
VALUES
('Create {{.ModulePlural}}', 'Permission to create new {{.LowerModulePlural}}', '{{.LowerModulePlural}}', 'create', 0, NOW(), NOW()),
('Read {{.ModulePlural}}', 'Permission to read {{.LowerModulePlural}}', '{{.LowerModulePlural}}', 'read', 0, NOW(), NOW()),
('Update {{.ModulePlural}}', 'Permission to update {{.LowerModulePlural}}', '{{.LowerModulePlural}}', 'update', 0, NOW(), NOW()),
('Delete {{.ModulePlural}}', 'Permission to delete {{.LowerModulePlural}}', '{{.LowerModulePlural}}', 'delete', 0, NOW(), NOW());

-- Assign permissions to Editor group (assuming group_id = 2)
INSERT INTO group_permissions (group_id, permission_id)
SELECT 2, id FROM permissions WHERE resource = '{{.LowerModulePlural}}' AND action IN ('create', 'read', 'update');

-- Assign read permission to Viewer group (assuming group_id = 3)
INSERT INTO group_permissions (group_id, permission_id)
SELECT 3, id FROM permissions WHERE resource = '{{.LowerModulePlural}}' AND action = 'read';

-- Assign read permission to General client group (assuming group_id = 4)
INSERT INTO group_permissions (group_id, permission_id)
SELECT 4, id FROM permissions WHERE resource = '{{.LowerModulePlural}}' AND action = 'read';

-- Assign all permissions to Admin group (assuming group_id = 1)
INSERT INTO group_permissions (group_id, permission_id)
SELECT 1, id FROM permissions WHERE resource = '{{.LowerModulePlural}}';
`

const testHTTPTemplate = `# {{.Module}} API Test File
# This file contains test requests for the {{.LowerModule}} module

### Create a new {{.LowerModule}} (Editor permission required)
POST http://localhost:3000/v1/{{.LowerModulePlural}}
Authorization: Bearer test-api-key-123
Content-Type: application/json

{
  "name": "Sample {{.Module}}",
  "description": "This is a sample {{.LowerModule}} for testing"
}

### Get all {{.LowerModulePlural}} (Editor permission required)
GET http://localhost:3000/v1/{{.LowerModulePlural}}
Authorization: Bearer test-api-key-123

### Get {{.LowerModule}} by ID (Editor permission required)
# Replace {id} with an actual {{.LowerModule}} ID
GET http://localhost:3000/v1/{{.LowerModulePlural}}/{id}
Authorization: Bearer test-api-key-123

### Update {{.LowerModule}} (Editor permission required)
# Replace {id} with an actual {{.LowerModule}} ID
PUT http://localhost:3000/v1/{{.LowerModulePlural}}/{id}
Authorization: Bearer test-api-key-123
Content-Type: application/json

{
  "name": "Updated {{.Module}}",
  "description": "This {{.LowerModule}} has been updated"
}

### Delete {{.LowerModule}} (Admin permission required - will fail with test-api-key-123)
# Replace {id} with an actual {{.LowerModule}} ID
DELETE http://localhost:3000/v1/{{.LowerModulePlural}}/{id}
Authorization: Bearer test-api-key-123

### Get deleted {{.LowerModulePlural}} (Editor permission required)
GET http://localhost:3000/v1/{{.LowerModulePlural}}/deleted
Authorization: Bearer test-api-key-123

### Restore deleted {{.LowerModule}} (Editor permission required)
# Replace {id} with an actual {{.LowerModule}} ID
POST http://localhost:3000/v1/{{.LowerModulePlural}}/{id}/restore
Authorization: Bearer test-api-key-123

### Try with Admin API key (all operations should work)
POST http://localhost:3000/v1/{{.LowerModulePlural}}
Authorization: Bearer admin-api-key-789
Content-Type: application/json

{
  "name": "Admin Created {{.Module}}",
  "description": "Created by admin user"
}

### Delete {{.LowerModule}} with Admin API key (should work)
# Replace {id} with an actual {{.LowerModule}} ID
DELETE http://localhost:3000/v1/{{.LowerModulePlural}}/{id}
Authorization: Bearer admin-api-key-789

### Test with Viewer API key (only read operations should work)
GET http://localhost:3000/v1/{{.LowerModulePlural}}
Authorization: Bearer viewer-api-key-456

### Try to create with Viewer API key (should fail)
POST http://localhost:3000/v1/{{.LowerModulePlural}}
Authorization: Bearer viewer-api-key-456
Content-Type: application/json

{
  "name": "Viewer Attempt",
  "description": "This should fail"
}
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
		fmt.Println("Usage: go run main.go <module-name> [--with-permissions]")
		fmt.Println("Example: go run main.go product")
		fmt.Println("Example: go run main.go product --with-permissions")
		os.Exit(1)
	}

	// Get module name from command line
	moduleName := os.Args[1]
	withPermissions := len(os.Args) > 2 && os.Args[2] == "--with-permissions"
	
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

	// Create permission script if requested
	if withPermissions {
		createPermissionScript(moduleName, data)
		createTestHTTPFile(moduleName, data)
	}

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

	if withPermissions {
		fmt.Println("\n🔐 Permission script created!")
		fmt.Printf("Run the following command to add permissions for %s module:\n", moduleName)
		fmt.Printf("go run scripts/add-%s-permissions.go\n", moduleName)
	} else {
		fmt.Println("\n💡 Tip: Use --with-permissions flag to automatically generate permission script")
		fmt.Printf("Example: go run tools/module-generator/main.go %s --with-permissions\n", moduleName)
	}
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

// createPermissionScript creates a permission script for the module
func createPermissionScript(moduleName string, data TemplateData) {
	// Create scripts directory if it doesn't exist
	scriptsDir := "scripts"
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		fmt.Printf("Error creating scripts directory: %v\n", err)
		return
	}

	// Create Go permission script
	permissionScriptPath := filepath.Join(scriptsDir, fmt.Sprintf("add-%s-permissions.go", moduleName))
	createFileFromTemplate(permissionScriptPath, permissionScriptTemplate, data)

	// Create SQL permission script
	sqlScriptPath := filepath.Join(scriptsDir, fmt.Sprintf("add-%s-permissions.sql", moduleName))
	createFileFromTemplate(sqlScriptPath, permissionSQLTemplate, data)
}

// createTestHTTPFile creates a test HTTP file for the module
func createTestHTTPFile(moduleName string, data TemplateData) {
	// Create test directory if it doesn't exist
	testDir := "test"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		fmt.Printf("Error creating test directory: %v\n", err)
		return
	}

	// Create HTTP test file
	testHTTPPath := filepath.Join(testDir, fmt.Sprintf("%s-api-test.http", moduleName))
	createFileFromTemplate(testHTTPPath, testHTTPTemplate, data)
}