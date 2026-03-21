// +build ignore

package middleware

import (
	"errors"
	"vamsasetu/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Example 1: Using ErrorHandler in main.go
func ExampleErrorHandler_Setup() {
	app := fiber.New(fiber.Config{
		// Set the custom error handler
		ErrorHandler: ErrorHandler(),
	})

	// Now all errors returned from handlers will be caught and formatted
	app.Get("/api/example", func(c *fiber.Ctx) error {
		// Return an error - it will be caught by ErrorHandler
		return BadRequest("Invalid input", nil)
	})

	// Start server
	_ = app.Listen(":8080")
}

// Example 2: Using AppError constructors in handlers
func ExampleAppError_InHandler() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Get("/api/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Validate input
		if id == "" {
			return BadRequest("User ID is required", nil)
		}

		// Simulate database lookup
		user, err := findUserByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NotFound("User not found", err)
			}
			return InternalServerError("Failed to fetch user", err)
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    user,
			"error":   "",
		})
	})
}

// Example 3: Using MapServiceError to convert service errors
func ExampleMapServiceError_InHandler() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Post("/api/users", func(c *fiber.Ctx) error {
		var req CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return BadRequest("Invalid request body", err)
		}

		// Call service layer
		user, err := createUser(req)
		if err != nil {
			// MapServiceError automatically converts service errors to AppErrors
			return MapServiceError(err)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"data":    user,
			"error":   "",
		})
	})
}

// Example 4: Returning validation errors
func ExampleValidationErrors_InHandler() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Post("/api/register", func(c *fiber.Ctx) error {
		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return BadRequest("Invalid request body", err)
		}

		// Validate request
		var errs utils.ValidationErrors

		if err := utils.ValidateEmail("email", req.Email); err != nil {
			errs = append(errs, *err)
		}

		if err := utils.ValidateRequired("name", req.Name); err != nil {
			errs = append(errs, *err)
		}

		if err := utils.ValidateMinLength("password", req.Password, 8); err != nil {
			errs = append(errs, *err)
		}

		// Return validation errors if any
		if errs.HasErrors() {
			return errs // ErrorHandler will catch this and return 400
		}

		// Process registration...
		return c.JSON(fiber.Map{
			"success": true,
			"data":    fiber.Map{"message": "Registration successful"},
			"error":   "",
		})
	})
}

// Example 5: Using response utilities with error handling
func ExampleResponseUtilities_WithErrorHandler() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Get("/api/members/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		member, err := getMemberByID(id)
		if err != nil {
			// Return error - ErrorHandler will format it
			return MapServiceError(err)
		}

		// Use response utilities for success
		return utils.SuccessResponse(c, member)
	})

	app.Post("/api/members", func(c *fiber.Ctx) error {
		var req CreateMemberRequest
		if err := c.BodyParser(&req); err != nil {
			return BadRequest("Invalid request body", err)
		}

		member, err := createMember(req)
		if err != nil {
			return MapServiceError(err)
		}

		// Use response utilities for created resource
		return utils.CreatedResponse(c, member)
	})
}

// Example 6: Custom error messages with underlying errors
func ExampleAppError_WithUnderlyingError() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Get("/api/data", func(c *fiber.Ctx) error {
		// Simulate a database error
		dbErr := errors.New("connection timeout")

		// Return user-friendly message with underlying error for logging
		// The underlying error will be logged but not exposed to the client
		return InternalServerError("Failed to fetch data. Please try again later.", dbErr)
	})
}

// Example 7: Different error types in a single handler
func ExampleMultipleErrorTypes() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Put("/api/resources/:id", func(c *fiber.Ctx) error {
		// Check authentication
		userID := c.Locals("userId")
		if userID == nil {
			return Unauthorized("Authentication required", nil)
		}

		// Validate input
		id := c.Params("id")
		if id == "" {
			return BadRequest("Resource ID is required", nil)
		}

		// Check permissions
		hasPermission := checkPermission(userID, id)
		if !hasPermission {
			return Forbidden("You don't have permission to modify this resource", nil)
		}

		// Check if resource exists
		resource, err := findResource(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NotFound("Resource not found", err)
			}
			return InternalServerError("Failed to fetch resource", err)
		}

		// Check for conflicts
		if resource.IsLocked {
			return Conflict("Resource is locked by another user", nil)
		}

		// Update resource...
		return utils.SuccessResponse(c, resource)
	})
}

// Mock types and functions for examples

type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
}

type RegisterRequest struct {
	Email    string
	Name     string
	Password string
}

type CreateMemberRequest struct {
	Name string
}

type User struct {
	ID    string
	Email string
	Name  string
}

type Resource struct {
	ID       string
	IsLocked bool
}

func findUserByID(id string) (*User, error) {
	return nil, gorm.ErrRecordNotFound
}

func createUser(req CreateUserRequest) (*User, error) {
	return &User{ID: "1", Email: req.Email, Name: req.Name}, nil
}

func getMemberByID(id string) (interface{}, error) {
	return nil, gorm.ErrRecordNotFound
}

func createMember(req CreateMemberRequest) (interface{}, error) {
	return fiber.Map{"id": "1", "name": req.Name}, nil
}

func checkPermission(userID interface{}, resourceID string) bool {
	return true
}

func findResource(id string) (*Resource, error) {
	return &Resource{ID: id, IsLocked: false}, nil
}

