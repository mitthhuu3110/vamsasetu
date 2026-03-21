package middleware

import (
	"errors"
	"fmt"
	"log"
	"vamsasetu/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AppError represents a structured application error with HTTP status code
type AppError struct {
	Code    int    // HTTP status code
	Message string // User-friendly error message
	Err     error  // Original error (for logging)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error for errors.Is and errors.As
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError with the given code, message, and underlying error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error constructors for convenience

// BadRequest creates a 400 Bad Request error
func BadRequest(message string, err error) *AppError {
	return NewAppError(fiber.StatusBadRequest, message, err)
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message string, err error) *AppError {
	return NewAppError(fiber.StatusUnauthorized, message, err)
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message string, err error) *AppError {
	return NewAppError(fiber.StatusForbidden, message, err)
}

// NotFound creates a 404 Not Found error
func NotFound(message string, err error) *AppError {
	return NewAppError(fiber.StatusNotFound, message, err)
}

// Conflict creates a 409 Conflict error
func Conflict(message string, err error) *AppError {
	return NewAppError(fiber.StatusConflict, message, err)
}

// InternalServerError creates a 500 Internal Server Error
func InternalServerError(message string, err error) *AppError {
	return NewAppError(fiber.StatusInternalServerError, message, err)
}

// ErrorHandler is a Fiber error handler middleware that catches errors
// and formats them consistently using the APIResponse format
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Default to 500 Internal Server Error
		code := fiber.StatusInternalServerError
		message := "An unexpected error occurred"

		// Check if it's an AppError
		var appErr *AppError
		if errors.As(err, &appErr) {
			code = appErr.Code
			message = appErr.Message

			// Log the underlying error if present
			if appErr.Err != nil {
				log.Printf("[ERROR] %s: %v", message, appErr.Err)
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// Handle GORM not found errors
			code = fiber.StatusNotFound
			message = "Resource not found"
			log.Printf("[ERROR] Record not found: %v", err)
		} else if validationErrs, ok := err.(utils.ValidationErrors); ok {
			// Handle validation errors
			code = fiber.StatusBadRequest
			message = validationErrs.Error()
			log.Printf("[VALIDATION] %s", message)
		} else if validationErr, ok := err.(*utils.ValidationError); ok {
			// Handle single validation error
			code = fiber.StatusBadRequest
			message = validationErr.Error()
			log.Printf("[VALIDATION] %s", message)
		} else if fiberErr, ok := err.(*fiber.Error); ok {
			// Handle Fiber's built-in errors
			code = fiberErr.Code
			message = fiberErr.Message
			log.Printf("[FIBER] %d: %s", code, message)
		} else {
			// Log unexpected errors with full details
			log.Printf("[ERROR] Unexpected error: %v", err)
		}

		// Return consistent error response using response utilities
		return utils.ErrorResponse(c, code, message)
	}
}

// MapServiceError maps common service layer errors to appropriate HTTP errors
// This helper function can be used in handlers to convert service errors to AppErrors
func MapServiceError(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's already an AppError
	var appErr *AppError
	if errors.As(err, &appErr) {
		return err
	}

	// Map GORM errors
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound("Resource not found", err)
	}

	// Map validation errors
	if _, ok := err.(utils.ValidationErrors); ok {
		return BadRequest(err.Error(), err)
	}
	if _, ok := err.(*utils.ValidationError); ok {
		return BadRequest(err.Error(), err)
	}

	// Default to internal server error for unknown errors
	return InternalServerError("An unexpected error occurred", err)
}

