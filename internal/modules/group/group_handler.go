package group

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

// CreateGroup godoc
// @Summary Create a new group
// @Description Create a new group with permissions
// @Tags Group
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param group body CreateGroupRequest true "Group data"
// @Success 201 {object} Group
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/groups [post]
func (h *Handler) CreateGroup(c *fiber.Ctx) error {
	var req CreateGroupRequest
	
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

	group := &Group{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.repo.CreateGroup(group); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create group",
		})
	}

	// Add permissions if provided
	if len(req.PermissionIDs) > 0 {
		if err := h.repo.UpdateGroupPermissions(group.ID, req.PermissionIDs); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Group created but failed to assign permissions",
			})
		}
	}

	// Fetch the group with permissions
	groupWithPermissions, err := h.repo.GetGroupWithPermissions(group.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Group created but failed to fetch details",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   groupWithPermissions,
	})
}

// GetGroups godoc
// @Summary Get all groups
// @Description Get list of all groups with their permissions
// @Tags Group
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Group
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/groups [get]
func (h *Handler) GetGroups(c *fiber.Ctx) error {
	groups, err := h.repo.GetAllGroups()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch groups",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   groups,
	})
}

// GetGroup godoc
// @Summary Get group by ID
// @Description Get a specific group by its ID with permissions
// @Tags Group
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Success 200 {object} Group
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/groups/{id} [get]
func (h *Handler) GetGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid group ID",
		})
	}

	group, err := h.repo.GetGroupWithPermissions(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Group not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   group,
	})
}

// UpdateGroupPermissions godoc
// @Summary Update group permissions
// @Description Update permissions assigned to a group
// @Tags Group
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Param permissions body UpdateGroupPermissionsRequest true "Permission IDs"
// @Success 200 {object} Group
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/groups/{id}/permissions [put]
func (h *Handler) UpdateGroupPermissions(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid group ID",
		})
	}

	var req UpdateGroupPermissionsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if err := h.repo.UpdateGroupPermissions(uint(id), req.PermissionIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update group permissions",
		})
	}

	// Fetch updated group with permissions
	group, err := h.repo.GetGroupWithPermissions(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Permissions updated but failed to fetch group details",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   group,
	})
}

// DeleteGroup godoc
// @Summary Delete group
// @Description Delete a group by ID
// @Tags Group
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/groups/{id} [delete]
func (h *Handler) DeleteGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid group ID",
		})
	}

	if err := h.repo.DeleteGroup(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete group",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Group deleted successfully",
	})
}