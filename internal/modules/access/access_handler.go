package access

import (
	"time"

	"apiserver/internal/utils"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	repo      Repository
	validator *validator.Validate
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo:      repo,
		validator: validator.New(),
	}
}

// GetProfile godoc
// SWAGGER_ACCESS_START
// @Summary Get user profile
// @Description Get current user profile information
// @Tags Access
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/profile [get]
// SWAGGER_ACCESS_END
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	// Get user data from middleware
	user := c.Locals("user").(*User)
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

// UpdateExpiredDate godoc
// SWAGGER_ACCESS_START
// @Summary Update API key expiration date
// @Description Update the expiration date for a user's API key
// @Tags Access
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param data body UpdateExpiredDateRequest true "Expiration date data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/access/{id}/expired-date [put]
// SWAGGER_ACCESS_END
func (h *Handler) UpdateExpiredDate(c *fiber.Ctx) error {
	// Parse user ID from path
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user ID",
		})
	}

	// Parse request body
	var req UpdateExpiredDateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Check if user exists
	user, err := h.repo.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	// Handle the expiration date
	var expiredDate *time.Time
	if req.ExpiredDate != nil {
		// Validate date is in the future
		if req.ExpiredDate.Before(time.Now()) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Expiration date must be in the future",
			})
		}
		expiredDate = req.ExpiredDate
	}

	// Update the expiration date
	if err := h.repo.UpdateExpiredDate(user.ID, expiredDate); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update expiration date",
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user_id":      user.ID,
			"email":        user.Email,
			"expired_date": expiredDate,
		},
	})
}

// RemoveExpiredDate godoc
// SWAGGER_ACCESS_START
// @Summary Remove API key expiration date
// @Description Remove the expiration date for a user's API key (never expires)
// @Tags Access
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/access/{id}/expired-date [delete]
// SWAGGER_ACCESS_END
func (h *Handler) RemoveExpiredDate(c *fiber.Ctx) error {
	// Parse user ID from path
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user ID",
		})
	}

	// Check if user exists
	user, err := h.repo.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	// Update the expiration date to NULL (never expires)
	if err := h.repo.UpdateExpiredDate(user.ID, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to remove expiration date",
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user_id":      user.ID,
			"email":        user.Email,
			"expired_date": nil,
			"message":      "API key will never expire",
		},
	})
}

// UpdateRateLimit godoc
// SWAGGER_ACCESS_START
// @Summary Update API key rate limit
// @Description Update the rate limit for a user's API key (requests per minute)
// @Tags Access
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param data body UpdateRateLimitRequest true "Rate limit data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/access/{id}/rate-limit [put]
// SWAGGER_ACCESS_END
func (h *Handler) UpdateRateLimit(c *fiber.Ctx) error {
	// Parse user ID from path
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user ID",
		})
	}

	// Parse request body
	var req UpdateRateLimitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validate rate limit
	if req.RateLimit < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Rate limit must be at least 1",
		})
	}

	// Check if user exists
	user, err := h.repo.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	// Update the rate limit
	if err := h.repo.UpdateRateLimit(user.ID, req.RateLimit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update rate limit",
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user_id":    user.ID,
			"email":      user.Email,
			"rate_limit": req.RateLimit,
		},
	})
}

// CreateAccess godoc
// SWAGGER_ACCESS_START
// @Summary Create new access
// @Description Create new access with API key, default expiration (6 months), group_id 4, and rate limit 120
// @Tags Access
// @Accept json
// @Produce json
// @Param data body CreateAccessRequest true "Access creation data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/access [post]
// SWAGGER_ACCESS_END
func (h *Handler) CreateAccess(c *fiber.Ctx) error {
	// Parse request body
	var req CreateAccessRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		return utils.HandleError(c, err)
	}

	// Check if email already exists
	if existingUser, err := h.repo.FindByEmail(req.Email); err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Email already exists",
		})
	}

	// Generate API key
	apiKey := utils.GenerateAPIKey()

	// Set default expiration date (6 months from now)
	expiredDate := time.Now().AddDate(0, 6, 0)

	// Create new user
	user := &User{
		Name:        req.FullName,
		Email:       req.Email,
		APIKey:      apiKey,
		GroupID:     utils.UintPtr(4), // Default group_id: 4 (generic client)
		ExpiredDate: &expiredDate,
		RateLimit:   120, // Default rate limit: 120
		StatusID:    utils.Int16Ptr(0), // Active status
	}

	// Save to database
	if err := h.repo.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create access",
		})
	}

	// Prepare response
	response := CreateAccessResponse{
		APIKey:      apiKey,
		ExpiredDate: &expiredDate,
		RateLimit:   120,
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}