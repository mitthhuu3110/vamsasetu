package middleware

import (
	"bytes"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// captureLog captures log output for testing
func captureLog(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)
	f()
	return buf.String()
}

// TestLoggerMiddleware_BasicLogging tests basic request logging
func TestLoggerMiddleware_BasicLogging(t *testing.T) {
	app := fiber.New()
	app.Use(LoggerMiddleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	logOutput := captureLog(func() {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	assert.Contains(t, logOutput, "[GET]")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "200")
}

// TestLoggerMiddleware_LogsMethod tests that HTTP method is logged
func TestLoggerMiddleware_LogsMethod(t *testing.T) {
	app := fiber.New()
	app.Use(LoggerMiddleware())

	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		app.Add(method, "/test", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			logOutput := captureLog(func() {
				req := httptest.NewRequest(method, "/test", nil)
				_, err := app.Test(req)
				assert.NoError(t, err)
			})

			assert.Contains(t, logOutput, "["+method+"]")
		})
	}
}

// TestLoggerMiddleware_LogsPath tests that request path is logged
func TestLoggerMiddleware_LogsPath(t *testing.T) {
	app := fiber.New()
	app.Use(LoggerMiddleware())

	paths := []string{"/api/members", "/api/events", "/api/relationships"}

	for _, path := range paths {
		app.Get(path, func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			logOutput := captureLog(func() {
				req := httptest.NewRequest("GET", path, nil)
				_, err := app.Test(req)
				assert.NoError(t, err)
			})

			assert.Contains(t, logOutput, path)
		})
	}
}

// TestLoggerMiddleware_LogsStatusCode tests that status code is logged
func TestLoggerMiddleware_LogsStatusCode(t *testing.T) {
	app := fiber.New()
	app.Use(LoggerMiddleware())

	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/created", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusCreated)
	})

	app.Get("/notfound", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusInternalServerError)
	})

	tests := []struct {
		path         string
		expectedCode string
	}{
		{"/ok", "200"},
		{"/created", "201"},
		{"/notfound", "404"},
		{"/error", "500"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			logOutput := captureLog(func() {
				req := httptest.NewRequest("GET", tt.path, nil)
				_, err := app.Test(req)
				assert.NoError(t, err)
			})

			assert.Contains(t, logOutput, tt.expectedCode)
		})
	}
}

// TestLoggerMiddleware_LogsDuration tests that request duration is logged
func TestLoggerMiddleware_LogsDuration(t *testing.T) {
	app := fiber.New()
	app.Use(LoggerMiddleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		time.Sleep(10 * time.Millisecond)
		return c.SendString("OK")
	})

	logOutput := captureLog(func() {
		req := httptest.NewRequest("GET", "/test", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	// Check that duration is logged (should contain time units like µs, ms, or s)
	assert.True(t,
		strings.Contains(logOutput, "µs") ||
			strings.Contains(logOutput, "ms") ||
			strings.Contains(logOutput, "s"),
		"Log should contain duration with time units")
}

// TestLoggerMiddleware_LogFormat tests the log format
func TestLoggerMiddleware_LogFormat(t *testing.T) {
	app := fiber.New()
	app.Use(LoggerMiddleware())

	app.Get("/api/members", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	logOutput := captureLog(func() {
		req := httptest.NewRequest("GET", "/api/members", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	// Verify format: [METHOD] /path - status - duration
	assert.Contains(t, logOutput, "[GET]")
	assert.Contains(t, logOutput, "/api/members")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, "-")
}

// TestLoggerMiddleware_WithError tests logging when handler returns error
func TestLoggerMiddleware_WithError(t *testing.T) {
	app := fiber.New()
	app.Use(LoggerMiddleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "Bad request")
	})

	logOutput := captureLog(func() {
		req := httptest.NewRequest("GET", "/test", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	assert.Contains(t, logOutput, "[GET]")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "400")
}

// TestDetailedLoggerMiddleware_WithoutUser tests detailed logger without user context
func TestDetailedLoggerMiddleware_WithoutUser(t *testing.T) {
	app := fiber.New()
	app.Use(DetailedLoggerMiddleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	logOutput := captureLog(func() {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	assert.Contains(t, logOutput, "[GET]")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "200")
	assert.NotContains(t, logOutput, "user:")
}

// TestDetailedLoggerMiddleware_WithUser tests detailed logger with user context
func TestDetailedLoggerMiddleware_WithUser(t *testing.T) {
	app := fiber.New()
	app.Use(DetailedLoggerMiddleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		c.Locals("userId", uint(123))
		return c.SendString("OK")
	})

	logOutput := captureLog(func() {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	assert.Contains(t, logOutput, "[GET]")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, "user:123")
}

// TestDetailedLoggerMiddleware_WithAuthMiddleware tests detailed logger after auth
func TestDetailedLoggerMiddleware_WithAuthMiddleware(t *testing.T) {
	app := fiber.New()

	// Simulate auth middleware setting user context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userId", uint(456))
		return c.Next()
	})

	app.Use(DetailedLoggerMiddleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	logOutput := captureLog(func() {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	assert.Contains(t, logOutput, "user:456")
}

// TestLoggerMiddleware_MultipleRequests tests logging multiple requests
func TestLoggerMiddleware_MultipleRequests(t *testing.T) {
	app := fiber.New()
	app.Use(LoggerMiddleware())

	app.Get("/test1", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/test2", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	logOutput := captureLog(func() {
		req1 := httptest.NewRequest("GET", "/test1", nil)
		_, err := app.Test(req1)
		assert.NoError(t, err)

		req2 := httptest.NewRequest("GET", "/test2", nil)
		_, err = app.Test(req2)
		assert.NoError(t, err)
	})

	assert.Contains(t, logOutput, "/test1")
	assert.Contains(t, logOutput, "/test2")
}

// TestLoggerMiddleware_Integration tests logger with other middleware
func TestLoggerMiddleware_Integration(t *testing.T) {
	app := fiber.New()

	// Add logger first
	app.Use(LoggerMiddleware())

	// Add CORS
	app.Use(CORSMiddleware("http://localhost:3000"))

	app.Post("/api/members", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})

	logOutput := captureLog(func() {
		req := httptest.NewRequest("POST", "/api/members", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	assert.Contains(t, logOutput, "[POST]")
	assert.Contains(t, logOutput, "/api/members")
	assert.Contains(t, logOutput, "200")
}

// TestLoggerMiddleware_WithQueryParams tests logging with query parameters
func TestLoggerMiddleware_WithQueryParams(t *testing.T) {
	app := fiber.New()
	app.Use(LoggerMiddleware())

	app.Get("/api/members", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	logOutput := captureLog(func() {
		req := httptest.NewRequest("GET", "/api/members?page=1&limit=10", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	// Path should be logged (query params may or may not be included depending on implementation)
	assert.Contains(t, logOutput, "/api/members")
}

// TestDetailedLoggerMiddleware_Format tests detailed logger format
func TestDetailedLoggerMiddleware_Format(t *testing.T) {
	app := fiber.New()
	app.Use(DetailedLoggerMiddleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		c.Locals("userId", uint(789))
		return c.SendStatus(fiber.StatusOK)
	})

	logOutput := captureLog(func() {
		req := httptest.NewRequest("GET", "/test", nil)
		_, err := app.Test(req)
		assert.NoError(t, err)
	})

	// Verify format: [METHOD] /path - status - duration - user:userId
	assert.Contains(t, logOutput, "[GET]")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, "user:789")
	assert.Contains(t, logOutput, "-")
}
