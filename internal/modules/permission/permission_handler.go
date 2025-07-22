package permission

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// CreatePermission godoc
// @Summary Create a new permission
// @Description Create a new permission with resource and action
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param permission body CreatePermissionRequest true "Permission data"
// @Success 201 {object} Permission
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/permissions [post]
func (h *Handler) CreatePermission(c *fiber.Ctx) error {
	var req CreatePermissionRequest
	
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

	if strings.TrimSpace(req.Resource) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Resource is required",
		})
	}

	if strings.TrimSpace(req.Action) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Action is required",
		})
	}

	permission := &Permission{
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
	}

	if err := h.repo.CreatePermission(permission); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create permission",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   permission,
	})
}

// GetPermissions godoc
// @Summary Get all permissions
// @Description Get list of all permissions
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Permission
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/permissions [get]
func (h *Handler) GetPermissions(c *fiber.Ctx) error {
	permissions, err := h.repo.GetAllPermissions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch permissions",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   permissions,
	})
}

// GetPermission godoc
// @Summary Get permission by ID
// @Description Get a specific permission by its ID
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Success 200 {object} Permission
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/permissions/{id} [get]
func (h *Handler) GetPermission(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid permission ID",
		})
	}

	permission, err := h.repo.GetPermissionByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   permission,
	})
}

// DeletePermission godoc
// @Summary Delete permission
// @Description Delete a permission by ID
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/permissions/{id} [delete]
func (h *Handler) DeletePermission(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid permission ID",
		})
	}

	if err := h.repo.DeletePermission(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete permission",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Permission deleted successfully",
	})
}