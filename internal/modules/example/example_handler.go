package example

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

// CreateExample godoc
// @Summary Create a new example
// @Description Create a new example with name and description
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param example body CreateExampleRequest true "Example data"
// @Success 201 {object} Example
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples [post]
func (h *Handler) CreateExample(c *fiber.Ctx) error {
	var req CreateExampleRequest
	
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

	example := &Example{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.repo.CreateExample(example); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create example",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   example,
	})
}

// GetExamples godoc
// @Summary Get all examples
// @Description Get list of all examples
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Example
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples [get]
func (h *Handler) GetExamples(c *fiber.Ctx) error {
	examples, err := h.repo.GetAllExamples()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch examples",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   examples,
	})
}

// GetExample godoc
// @Summary Get example by ID
// @Description Get a specific example by its ID
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Example ID"
// @Success 200 {object} Example
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/{id} [get]
func (h *Handler) GetExample(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid example ID",
		})
	}

	example, err := h.repo.GetExampleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Example not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   example,
	})
}

// UpdateExample godoc
// @Summary Update example
// @Description Update an existing example
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Example ID"
// @Param example body CreateExampleRequest true "Example data"
// @Success 200 {object} Example
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/{id} [put]
func (h *Handler) UpdateExample(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid example ID",
		})
	}

	// Check if example exists
	example, err := h.repo.GetExampleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Example not found",
		})
	}

	var req CreateExampleRequest
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

	// Update example
	example.Name = req.Name
	example.Description = req.Description

	if err := h.repo.UpdateExample(example); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update example",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   example,
	})
}

// SoftDeleteExample godoc
// @Summary Soft delete example
// @Description Soft delete an example (set status_id to 0)
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Example ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/{id} [delete]
func (h *Handler) SoftDeleteExample(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid example ID",
		})
	}

	// Check if example exists
	_, err := h.repo.GetExampleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Example not found",
		})
	}

	if err := h.repo.SoftDeleteExample(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete example",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Example deleted successfully",
	})
}

// RestoreExample godoc
// @Summary Restore example
// @Description Restore a soft deleted example (set status_id to 1)
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Example ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/{id}/restore [post]
func (h *Handler) RestoreExample(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid example ID",
		})
	}

	if err := h.repo.RestoreExample(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to restore example",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Example restored successfully",
	})
}

// GetDeletedExamples godoc
// @Summary Get deleted examples
// @Description Get list of all soft deleted examples
// @Tags Example
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Example
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/examples/deleted [get]
func (h *Handler) GetDeletedExamples(c *fiber.Ctx) error {
	examples, err := h.repo.GetDeletedExamples()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch deleted examples",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   examples,
	})
}