package utils

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors"`
	Data    any    `json:"data"`
}

type EmptyResponse struct{}

// Output is a helper to return a JSON response with optional success and status code.
// Usage:
// utils.Output(c, "ok")
// utils.Output(c, "created", true, 201)
// utils.Output(c, "error", false, 400)
func Output(c *fiber.Ctx, message string, successAndStatus ...interface{}) error {
	isSuccess := true
	statusCode := 200

	if len(successAndStatus) > 0 {
		// First optional param: success (bool)
		if s, ok := successAndStatus[0].(bool); ok {
			isSuccess = s
		}
	}
	if len(successAndStatus) > 1 {
		// Second optional param: status (int)
		if code, ok := successAndStatus[1].(int); ok {
			statusCode = code
		}
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"success": isSuccess,
		"message": message,
	})
}

type ApiError struct {
	Field string `json:"field"`
	Msg   string `json:"message"`
}

func SuccessResponse(message string, data any, errors any) Response {
	res := Response{
		Success: true,
		Message: message,
		Errors:  errors,
		Data:    data,
	}

	return res
}

func ErrorResponse(message string, err ...error) Response {
	var ve validator.ValidationErrors
	var finalMessage string

	if len(err) == 0 {
		return Response{
			Success: false,
			Message: message,
			Errors:  nil,
			Data:    nil,
		}
	}

	if errors.As(err[0], &ve) {
		errorsMap := map[string][]string{}
		for i, fe := range ve {
			field := toSnakeCase(fe.Field())
			msg := formatValidationMessage(fe)
			errorsMap[field] = append(errorsMap[field], msg)
			if i == 0 {
				finalMessage = msg
			}
		}

		if len(errorsMap) > 1 {
			finalMessage = fmt.Sprintf("%s and %d other error(s)", finalMessage, len(errorsMap)-1)
		}

		return Response{
			Success: false,
			Message: finalMessage,
			Errors:  errorsMap,
			Data:    nil,
		}
	}

	// fallback error
	return Response{
		Success: false,
		Message: message,
		Errors:  err,
		Data:    nil,
	}
}

func formatValidationMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required.", field)
	case "email":
		return fmt.Sprintf("The %s must be a valid email address.", field)
	case "min":
		return fmt.Sprintf("The %s field must be at least %s characters.", field, fe.Param())
	case "oneof":
		options := strings.ReplaceAll(fe.Param(), " ", ", ")
		return fmt.Sprintf("The %s field must be one of: %s.", field, options)
	default:
		return fmt.Sprintf("The %s field is invalid.", field)
	}
}

func toSnakeCase(str string) string {
	var sb strings.Builder
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			sb.WriteByte('_')
		}
		sb.WriteRune(r)
	}
	return strings.ToLower(sb.String())
}

func HandleError(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	switch e := err.(type) {
	case validator.ValidationErrors:
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			err.Error(),
			err,
		))

	case *fiber.Error:
		code = e.Code
		if code == fiber.StatusInternalServerError {
			// TODO: Implement proper error logging
			return c.Status(code).JSON(ErrorResponse(
				"Internal Server Error",
				fiber.NewError(code, "Internal Server Error"),
			))
		}
		return c.Status(code).JSON(ErrorResponse(
			e.Error(),
			e,
		))

	default:
		return c.Status(code).JSON(ErrorResponse(
			err.Error(),
			err,
		))
	}
}
