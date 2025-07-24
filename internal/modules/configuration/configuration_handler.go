package configuration

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

// CreateConfiguration godoc
// @Summary Create a new configuration
// @Description Create a new configuration with name and description
// @Tags Configuration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param configuration body CreateConfigurationRequest true "Configuration data"
// @Success 201 {object} Configuration
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/configurations [post]
func (h *Handler) CreateConfiguration(c *fiber.Ctx) error {
	var req CreateConfigurationRequest
	
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

	configuration := &Configuration{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.repo.CreateConfiguration(configuration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create configuration",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   configuration,
	})
}

// GetConfigurations godoc
// @Summary Get all configurations
// @Description Get list of all configurations
// @Tags Configuration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Configuration
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/configurations [get]
func (h *Handler) GetConfigurations(c *fiber.Ctx) error {
	configurations, err := h.repo.GetAllConfigurations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch configurations",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   configurations,
	})
}

// GetConfiguration godoc
// @Summary Get configuration by ID
// @Description Get a specific configuration by its ID
// @Tags Configuration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Configuration ID"
// @Success 200 {object} Configuration
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/configurations/{id} [get]
func (h *Handler) GetConfiguration(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid configuration ID",
		})
	}

	configuration, err := h.repo.GetConfigurationByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Configuration not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   configuration,
	})
}

// UpdateConfiguration godoc
// @Summary Update configuration
// @Description Update an existing configuration
// @Tags Configuration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Configuration ID"
// @Param configuration body CreateConfigurationRequest true "Configuration data"
// @Success 200 {object} Configuration
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/configurations/{id} [put]
func (h *Handler) UpdateConfiguration(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid configuration ID",
		})
	}

	// Check if configuration exists
	configuration, err := h.repo.GetConfigurationByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Configuration not found",
		})
	}

	var req CreateConfigurationRequest
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

	// Update configuration
	configuration.Name = req.Name
	configuration.Description = req.Description

	if err := h.repo.UpdateConfiguration(configuration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update configuration",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   configuration,
	})
}

// SoftDeleteConfiguration godoc
// @Summary Soft delete configuration
// @Description Soft delete a configuration (set status_id to 0)
// @Tags Configuration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Configuration ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/configurations/{id} [delete]
func (h *Handler) SoftDeleteConfiguration(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid configuration ID",
		})
	}

	// Check if configuration exists
	_, err := h.repo.GetConfigurationByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Configuration not found",
		})
	}

	if err := h.repo.SoftDeleteConfiguration(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete configuration",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Configuration deleted successfully",
	})
}

// RestoreConfiguration godoc
// @Summary Restore configuration
// @Description Restore a soft deleted configuration (set status_id to 1)
// @Tags Configuration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Configuration ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/configurations/{id}/restore [post]
func (h *Handler) RestoreConfiguration(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid configuration ID",
		})
	}

	if err := h.repo.RestoreConfiguration(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to restore configuration",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Configuration restored successfully",
	})
}

// GetDeletedConfigurations godoc
// @Summary Get deleted configurations
// @Description Get list of all soft deleted configurations
// @Tags Configuration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Configuration
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/configurations/deleted [get]
func (h *Handler) GetDeletedConfigurations(c *fiber.Ctx) error {
	configurations, err := h.repo.GetDeletedConfigurations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch deleted configurations",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   configurations,
	})
}
