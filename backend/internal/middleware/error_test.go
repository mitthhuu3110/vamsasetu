package middleware

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"vamsasetu/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestAppError_Error tests the Error method of AppError
func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		expected string
	}{
		{
			name: "error with underlying error",
			appError: &AppError{
				Code:    400,
				Message: "Invalid input",
				Err:     errors.New("field validation failed"),
			},
			expected: "Invalid input: field validation failed",
		},
		{
			name: "error without underlying error",
			appError: &AppError{
				Code:    404,
				Message: "Resource not found",
				Err:     nil,
			},
			expected: "Resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appError.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAppError_Unwrap tests the Unwrap method of AppError
func TestAppError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	appErr := &AppError{
		Code:    500,
		Message: "Server error",
		Err:     underlyingErr,
	}

	unwrapped := appErr.Unwrap()
	assert.Equal(t, underlyingErr, unwrapped)
}

// TestNewAppError tests the NewAppError constructor
func TestNewAppError(t *testing.T) {
	code := 400
	message := "Bad request"
	err := errors.New("validation failed")

	appErr := NewAppError(code, message, err)

	assert.Equal(t, code, appErr.Code)
	assert.Equal(t, message, appErr.Message)
	assert.Equal(t, err, appErr.Err)
}

// TestErrorConstructors tests all error constructor functions
func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name         string
		constructor  func(string, error) *AppError
		expectedCode int
	}{
		{"BadRequest", BadRequest, fiber.StatusBadRequest},
		{"Unauthorized", Unauthorized, fiber.StatusUnauthorized},
		{"Forbidden", Forbidden, fiber.StatusForbidden},
		{"NotFound", NotFound, fiber.StatusNotFound},
		{"Conflict", Conflict, fiber.StatusConflict},
		{"InternalServerError", InternalServerError, fiber.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := "test message"
			err := errors.New("test error")

			appErr := tt.constructor(message, err)

			assert.Equal(t, tt.expectedCode, appErr.Code)
			assert.Equal(t, message, appErr.Message)
			assert.Equal(t, err, appErr.Err)
		})
	}
}

// TestErrorHandler_AppError tests ErrorHandler with AppError
func TestErrorHandler_AppError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return BadRequest("Invalid input", errors.New("validation failed"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid input")
	assert.Contains(t, string(body), `"success":false`)
}

// TestErrorHandler_GORMNotFound tests ErrorHandler with GORM not found error
func TestErrorHandler_GORMNotFound(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return gorm.ErrRecordNotFound
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Resource not found")
	assert.Contains(t, string(body), `"success":false`)
}

// TestErrorHandler_ValidationErrors tests ErrorHandler with ValidationErrors
func TestErrorHandler_ValidationErrors(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		var errs utils.ValidationErrors
		errs.Add("email", "email is required")
		errs.Add("name", "name must be at least 3 characters")
		return errs
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "email")
	assert.Contains(t, string(body), `"success":false`)
}

// TestErrorHandler_SingleValidationError tests ErrorHandler with single ValidationError
func TestErrorHandler_SingleValidationError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return &utils.ValidationError{
			Field:   "email",
			Message: "email format is invalid",
		}
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "email")
	assert.Contains(t, string(body), "invalid")
	assert.Contains(t, string(body), `"success":false`)
}

// TestErrorHandler_FiberError tests ErrorHandler with Fiber's built-in error
func TestErrorHandler_FiberError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusTeapot, "I'm a teapot")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusTeapot, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "I'm a teapot")
	assert.Contains(t, string(body), `"success":false`)
}

// TestErrorHandler_UnexpectedError tests ErrorHandler with unexpected error
func TestErrorHandler_UnexpectedError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return errors.New("unexpected error")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "An unexpected error occurred")
	assert.Contains(t, string(body), `"success":false`)
}

// TestMapServiceError tests the MapServiceError helper function
func TestMapServiceError(t *testing.T) {
	tests := []struct {
		name         string
		inputError   error
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "nil error",
			inputError:   nil,
			expectedCode: 0,
			expectedMsg:  "",
		},
		{
			name:         "already an AppError",
			inputError:   BadRequest("Invalid input", nil),
			expectedCode: fiber.StatusBadRequest,
			expectedMsg:  "Invalid input",
		},
		{
			name:         "GORM not found error",
			inputError:   gorm.ErrRecordNotFound,
			expectedCode: fiber.StatusNotFound,
			expectedMsg:  "Resource not found",
		},
		{
			name: "ValidationErrors",
			inputError: utils.ValidationErrors{
				{Field: "email", Message: "email is required"},
			},
			expectedCode: fiber.StatusBadRequest,
			expectedMsg:  "email: email is required",
		},
		{
			name: "single ValidationError",
			inputError: &utils.ValidationError{
				Field:   "name",
				Message: "name is required",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedMsg:  "name: name is required",
		},
		{
			name:         "unknown error",
			inputError:   errors.New("unknown error"),
			expectedCode: fiber.StatusInternalServerError,
			expectedMsg:  "An unexpected error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapServiceError(tt.inputError)

			if tt.inputError == nil {
				assert.Nil(t, result)
				return
			}

			var appErr *AppError
			assert.True(t, errors.As(result, &appErr))
			assert.Equal(t, tt.expectedCode, appErr.Code)
			assert.Equal(t, tt.expectedMsg, appErr.Message)
		})
	}
}

// TestErrorHandler_Integration tests the full error handling flow
func TestErrorHandler_Integration(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	// Endpoint that returns different errors based on query parameter
	app.Get("/test", func(c *fiber.Ctx) error {
		errorType := c.Query("type")

		switch errorType {
		case "bad_request":
			return BadRequest("Invalid request", nil)
		case "unauthorized":
			return Unauthorized("Authentication required", nil)
		case "forbidden":
			return Forbidden("Access denied", nil)
		case "not_found":
			return NotFound("Resource not found", nil)
		case "conflict":
			return Conflict("Resource already exists", nil)
		case "internal":
			return InternalServerError("Server error", nil)
		case "validation":
			var errs utils.ValidationErrors
			errs.Add("field1", "error1")
			return errs
		default:
			return errors.New("unknown error")
		}
	})

	tests := []struct {
		errorType    string
		expectedCode int
	}{
		{"bad_request", fiber.StatusBadRequest},
		{"unauthorized", fiber.StatusUnauthorized},
		{"forbidden", fiber.StatusForbidden},
		{"not_found", fiber.StatusNotFound},
		{"conflict", fiber.StatusConflict},
		{"internal", fiber.StatusInternalServerError},
		{"validation", fiber.StatusBadRequest},
		{"unknown", fiber.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.errorType, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test?type="+tt.errorType, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			assert.Contains(t, string(body), `"success":false`)
			assert.Contains(t, string(body), `"error"`)
		})
	}
}

// TestErrorHandler_WithUnderlyingError tests that underlying errors are properly logged
func TestErrorHandler_WithUnderlyingError(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	underlyingErr := errors.New("database connection failed")
	app.Get("/test", func(c *fiber.Ctx) error {
		return InternalServerError("Failed to fetch data", underlyingErr)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	// User-friendly message should be in response
	assert.Contains(t, string(body), "Failed to fetch data")
	// Underlying error should NOT be exposed to client
	assert.NotContains(t, string(body), "database connection failed")
}

// TestErrorHandler_ConsistentResponseFormat tests that all errors return consistent format
func TestErrorHandler_ConsistentResponseFormat(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return BadRequest("Test error", nil)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	// Verify APIResponse format
	assert.Contains(t, bodyStr, `"success"`)
	assert.Contains(t, bodyStr, `"data"`)
	assert.Contains(t, bodyStr, `"error"`)
	assert.Contains(t, bodyStr, `"success":false`)
	assert.Contains(t, bodyStr, `"data":null`)
}

