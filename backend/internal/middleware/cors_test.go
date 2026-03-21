package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TestCORSMiddleware_AllowedOrigin tests CORS with allowed origin
func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	app := fiber.New()
	app.Use(CORSMiddleware("http://localhost:3000"))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
}

// TestCORSMiddleware_PreflightRequest tests CORS preflight OPTIONS request
func TestCORSMiddleware_PreflightRequest(t *testing.T) {
	app := fiber.New()
	app.Use(CORSMiddleware("http://localhost:3000"))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type,Authorization")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Headers"), "Authorization")
}

// TestCORSMiddleware_AllowedMethods tests that allowed methods are configured
func TestCORSMiddleware_AllowedMethods(t *testing.T) {
	app := fiber.New()
	app.Use(CORSMiddleware("http://localhost:3000"))

	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")

	resp, err := app.Test(req)
	assert.NoError(t, err)

	allowedMethods := resp.Header.Get("Access-Control-Allow-Methods")
	assert.Contains(t, allowedMethods, "GET")
	assert.Contains(t, allowedMethods, "POST")
	assert.Contains(t, allowedMethods, "PUT")
	assert.Contains(t, allowedMethods, "DELETE")
	assert.Contains(t, allowedMethods, "OPTIONS")
}

// TestCORSMiddleware_AllowedHeaders tests that allowed headers are configured
func TestCORSMiddleware_AllowedHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(CORSMiddleware("http://localhost:3000"))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Headers", "Authorization,Content-Type")

	resp, err := app.Test(req)
	assert.NoError(t, err)

	allowedHeaders := resp.Header.Get("Access-Control-Allow-Headers")
	assert.Contains(t, allowedHeaders, "Authorization")
	assert.Contains(t, allowedHeaders, "Content-Type")
	assert.Contains(t, allowedHeaders, "Accept")
	assert.Contains(t, allowedHeaders, "Origin")
}

// TestCORSMiddleware_Credentials tests that credentials are allowed
func TestCORSMiddleware_Credentials(t *testing.T) {
	app := fiber.New()
	app.Use(CORSMiddleware("http://localhost:3000"))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
}

// TestCORSMiddleware_MaxAge tests that max age is set
func TestCORSMiddleware_MaxAge(t *testing.T) {
	app := fiber.New()
	app.Use(CORSMiddleware("http://localhost:3000"))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, "86400", resp.Header.Get("Access-Control-Max-Age"))
}

// TestCORSMiddleware_MultipleOrigins tests CORS with multiple origins
func TestCORSMiddleware_MultipleOrigins(t *testing.T) {
	app := fiber.New()
	app.Use(CORSMiddleware("http://localhost:3000,http://localhost:3001"))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	tests := []struct {
		name   string
		origin string
	}{
		{"first origin", "http://localhost:3000"},
		{"second origin", "http://localhost:3001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", tt.origin)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			assert.Contains(t, resp.Header.Get("Access-Control-Allow-Origin"), tt.origin)
		})
	}
}

// TestCORSMiddleware_WithAuthorizationHeader tests CORS with Authorization header
func TestCORSMiddleware_WithAuthorizationHeader(t *testing.T) {
	app := fiber.New()
	app.Use(CORSMiddleware("http://localhost:3000"))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Authorization", "Bearer token123")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
}

// TestCORSMiddleware_Integration tests full CORS flow
func TestCORSMiddleware_Integration(t *testing.T) {
	app := fiber.New()
	app.Use(CORSMiddleware("http://localhost:3000"))

	app.Post("/api/members", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})

	// Step 1: Preflight request
	preflightReq := httptest.NewRequest("OPTIONS", "/api/members", nil)
	preflightReq.Header.Set("Origin", "http://localhost:3000")
	preflightReq.Header.Set("Access-Control-Request-Method", "POST")
	preflightReq.Header.Set("Access-Control-Request-Headers", "Content-Type,Authorization")

	preflightResp, err := app.Test(preflightReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, preflightResp.StatusCode)

	// Step 2: Actual request
	actualReq := httptest.NewRequest("POST", "/api/members", nil)
	actualReq.Header.Set("Origin", "http://localhost:3000")
	actualReq.Header.Set("Content-Type", "application/json")
	actualReq.Header.Set("Authorization", "Bearer token123")

	actualResp, err := app.Test(actualReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, actualResp.StatusCode)
	assert.Equal(t, "http://localhost:3000", actualResp.Header.Get("Access-Control-Allow-Origin"))
}
