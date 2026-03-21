package utils

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Helper function to parse response body
func parseResponse(t *testing.T, body io.Reader) APIResponse {
	var response APIResponse
	err := json.NewDecoder(body).Decode(&response)
	assert.NoError(t, err, "Failed to parse response body")
	return response
}

func TestSuccessResponse(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return SuccessResponse(c, fiber.Map{
			"message": "Operation successful",
			"id":      123,
		})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.Empty(t, response.Error)
	
	// Verify data structure
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Operation successful", data["message"])
	assert.Equal(t, float64(123), data["id"]) // JSON numbers are float64
}

func TestCreatedResponse(t *testing.T) {
	app := fiber.New()
	
	app.Post("/test", func(c *fiber.Ctx) error {
		return CreatedResponse(c, fiber.Map{
			"id":   "uuid-123",
			"name": "Test User",
		})
	})
	
	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.Empty(t, response.Error)
}

func TestErrorResponse(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Equal(t, "Invalid input", response.Error)
}

func TestBadRequestResponse(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return BadRequestResponse(c, "Missing required field")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Equal(t, "Missing required field", response.Error)
}

func TestUnauthorizedResponse(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return UnauthorizedResponse(c, "Invalid token")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Equal(t, "Invalid token", response.Error)
}

func TestForbiddenResponse(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return ForbiddenResponse(c, "Insufficient permissions")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Equal(t, "Insufficient permissions", response.Error)
}

func TestNotFoundResponse(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return NotFoundResponse(c, "Resource not found")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Equal(t, "Resource not found", response.Error)
}

func TestInternalServerErrorResponse(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return InternalServerErrorResponse(c, "Database connection failed")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Equal(t, "Database connection failed", response.Error)
}

func TestAPIResponseStructure(t *testing.T) {
	// Test that APIResponse struct has correct JSON tags
	response := APIResponse{
		Success: true,
		Data:    "test data",
		Error:   "",
	}
	
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)
	
	var parsed map[string]interface{}
	err = json.Unmarshal(jsonData, &parsed)
	assert.NoError(t, err)
	
	// Verify JSON field names
	assert.Contains(t, parsed, "success")
	assert.Contains(t, parsed, "data")
	assert.Contains(t, parsed, "error")
	
	assert.Equal(t, true, parsed["success"])
	assert.Equal(t, "test data", parsed["data"])
	assert.Equal(t, "", parsed["error"])
}

func TestSuccessResponseWithNilData(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return SuccessResponse(c, nil)
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.True(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Empty(t, response.Error)
}

func TestSuccessResponseWithComplexData(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return SuccessResponse(c, fiber.Map{
			"users": []fiber.Map{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
			"total": 2,
			"page":  1,
		})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(2), data["total"])
	assert.Equal(t, float64(1), data["page"])
	
	users, ok := data["users"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, users, 2)
}

func TestErrorResponseWithEmptyMessage(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return ErrorResponse(c, fiber.StatusBadRequest, "")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	
	response := parseResponse(t, resp.Body)
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Empty(t, response.Error)
}

func TestMultipleResponsesInSequence(t *testing.T) {
	app := fiber.New()
	
	app.Get("/success", func(c *fiber.Ctx) error {
		return SuccessResponse(c, fiber.Map{"status": "ok"})
	})
	
	app.Get("/error", func(c *fiber.Ctx) error {
		return BadRequestResponse(c, "error occurred")
	})
	
	// Test success endpoint
	req1 := httptest.NewRequest("GET", "/success", nil)
	resp1, err := app.Test(req1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)
	
	response1 := parseResponse(t, resp1.Body)
	assert.True(t, response1.Success)
	
	// Test error endpoint
	req2 := httptest.NewRequest("GET", "/error", nil)
	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp2.StatusCode)
	
	response2 := parseResponse(t, resp2.Body)
	assert.False(t, response2.Success)
}
