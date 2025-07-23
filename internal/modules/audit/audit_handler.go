package audit

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// GetAuditLogs godoc
// SWAGGER_AUDIT_START
// @Summary Get audit logs
// @Description Get audit logs with filtering and pagination
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param access_id query string false "Filter by access ID (UUID)"
// @Param user_email query string false "Filter by user email"
// @Param method query string false "Filter by HTTP method"
// @Param path query string false "Filter by API path"
// @Param status_code query int false "Filter by status code"
// @Param date_from query string false "Filter from date (YYYY-MM-DD)"
// @Param date_to query string false "Filter to date (YYYY-MM-DD)"
// @Param limit query int false "Limit results (default: 50, max: 1000)"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/audit-logs [get]
// SWAGGER_AUDIT_END
func (h *Handler) GetAuditLogs(c *fiber.Ctx) error {
	filter := AuditLogFilter{
		AccessID:   c.Query("access_id"),
		UserEmail:  c.Query("user_email"),
		Method:     c.Query("method"),
		Path:       c.Query("path"),
		DateFrom:   c.Query("date_from"),
		DateTo:     c.Query("date_to"),
	}

	if statusCode := c.Query("status_code"); statusCode != "" {
		if code, err := strconv.Atoi(statusCode); err == nil {
			filter.StatusCode = code
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}

	logs, total, err := h.repo.GetAuditLogs(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch audit logs",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"logs":   logs,
			"total":  total,
			"limit":  filter.Limit,
			"offset": filter.Offset,
		},
	})
}

// GetAuditLog godoc
// SWAGGER_AUDIT_START
// @Summary Get audit log by ID
// @Description Get detailed audit log by ID including request/response body
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Audit Log ID"
// @Success 200 {object} AuditLog
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/audit-logs/{id} [get]
// SWAGGER_AUDIT_END
func (h *Handler) GetAuditLog(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid audit log ID",
		})
	}

	log, err := h.repo.GetAuditLogByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Audit log not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   log,
	})
}

// DeleteOldLogs godoc
// SWAGGER_AUDIT_START
// @Summary Delete old audit logs
// @Description Delete audit logs older than specified days
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int true "Delete logs older than this many days"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/audit-logs/cleanup [delete]
// SWAGGER_AUDIT_END
func (h *Handler) DeleteOldLogs(c *fiber.Ctx) error {
	daysStr := c.Query("days")
	if daysStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Days parameter is required",
		})
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid days parameter",
		})
	}

	if err := h.repo.DeleteOldLogs(days); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete old logs",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Old audit logs deleted successfully",
	})
}