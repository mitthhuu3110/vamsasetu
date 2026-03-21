package utils

import "github.com/gofiber/fiber/v2"

// APIResponse represents the standard API response format
// All API responses follow this structure for consistency
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

// SuccessResponse creates a successful API response with data
// Returns HTTP 200 OK with success=true, populated data, and empty error
func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success: true,
		Data:    data,
		Error:   "",
	})
}

// CreatedResponse creates a successful creation response
// Returns HTTP 201 Created with success=true, populated data, and empty error
func CreatedResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Data:    data,
		Error:   "",
	})
}

// ErrorResponse creates an error response with a message
// Returns the specified HTTP status code with success=false, null data, and error message
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Data:    nil,
		Error:   message,
	})
}

// BadRequestResponse creates a 400 Bad Request error response
func BadRequestResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message)
}

// UnauthorizedResponse creates a 401 Unauthorized error response
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message)
}

// ForbiddenResponse creates a 403 Forbidden error response
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message)
}

// NotFoundResponse creates a 404 Not Found error response
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message)
}

// InternalServerErrorResponse creates a 500 Internal Server Error response
func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message)
}
