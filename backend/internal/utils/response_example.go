// +build ignore

package utils

import "github.com/gofiber/fiber/v2"

// Example usage of response utilities in handlers

// ExampleSuccessHandler demonstrates using SuccessResponse
func ExampleSuccessHandler(c *fiber.Ctx) error {
	// Return a successful response with data
	return SuccessResponse(c, fiber.Map{
		"message": "Operation completed successfully",
		"userId":  123,
	})
}

// ExampleCreatedHandler demonstrates using CreatedResponse
func ExampleCreatedHandler(c *fiber.Ctx) error {
	// Return a 201 Created response after creating a resource
	newUser := fiber.Map{
		"id":    "uuid-123",
		"name":  "John Doe",
		"email": "john@example.com",
	}
	return CreatedResponse(c, newUser)
}

// ExampleBadRequestHandler demonstrates using BadRequestResponse
func ExampleBadRequestHandler(c *fiber.Ctx) error {
	// Return a 400 Bad Request when validation fails
	return BadRequestResponse(c, "Email is required")
}

// ExampleUnauthorizedHandler demonstrates using UnauthorizedResponse
func ExampleUnauthorizedHandler(c *fiber.Ctx) error {
	// Return a 401 Unauthorized when authentication fails
	return UnauthorizedResponse(c, "Invalid or expired token")
}

// ExampleForbiddenHandler demonstrates using ForbiddenResponse
func ExampleForbiddenHandler(c *fiber.Ctx) error {
	// Return a 403 Forbidden when user lacks permissions
	return ForbiddenResponse(c, "Insufficient permissions to access this resource")
}

// ExampleNotFoundHandler demonstrates using NotFoundResponse
func ExampleNotFoundHandler(c *fiber.Ctx) error {
	// Return a 404 Not Found when resource doesn't exist
	return NotFoundResponse(c, "User not found")
}

// ExampleInternalServerErrorHandler demonstrates using InternalServerErrorResponse
func ExampleInternalServerErrorHandler(c *fiber.Ctx) error {
	// Return a 500 Internal Server Error when something goes wrong
	return InternalServerErrorResponse(c, "Database connection failed")
}

// ExampleCustomErrorHandler demonstrates using ErrorResponse with custom status code
func ExampleCustomErrorHandler(c *fiber.Ctx) error {
	// Return a custom error response with specific status code
	return ErrorResponse(c, fiber.StatusConflict, "Resource already exists")
}

// ExampleComplexDataHandler demonstrates returning complex data structures
func ExampleComplexDataHandler(c *fiber.Ctx) error {
	// Return paginated data with metadata
	return SuccessResponse(c, fiber.Map{
		"users": []fiber.Map{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
		},
		"pagination": fiber.Map{
			"total":       100,
			"page":        1,
			"limit":       10,
			"totalPages":  10,
		},
	})
}

// ExampleConditionalResponseHandler demonstrates conditional responses
func ExampleConditionalResponseHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	
	// Simulate database lookup
	user, err := findUserByID(userID)
	if err != nil {
		return InternalServerErrorResponse(c, "Failed to fetch user")
	}
	
	if user == nil {
		return NotFoundResponse(c, "User not found")
	}
	
	return SuccessResponse(c, user)
}

// Mock function for example
func findUserByID(id string) (interface{}, error) {
	// This is just a placeholder for the example
	return nil, nil
}
