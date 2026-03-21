package handler

import (
	"vamsasetu/backend/internal/middleware"
	"vamsasetu/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new authentication handler instance
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRoutes registers all authentication routes with the Fiber app
func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	auth := app.Group("/api/auth")

	// Public routes
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.RefreshToken)

	// Protected routes
	auth.Get("/profile", middleware.AuthMiddleware(), h.GetProfile)
}

// Register handles user registration
// POST /api/auth/register
// Request body: { "email": string, "password": string, "name": string, "role": string }
// Response: { "success": bool, "data": { "accessToken": string, "refreshToken": string, "user": User }, "error": string }
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req service.RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	// Call auth service
	authResponse, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    authResponse,
		"error":   "",
	})
}

// Login handles user authentication
// POST /api/auth/login
// Request body: { "email": string, "password": string }
// Response: { "success": bool, "data": { "accessToken": string, "refreshToken": string, "user": User }, "error": string }
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req service.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	// Call auth service
	authResponse, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    authResponse,
		"error":   "",
	})
}

// RefreshToken handles token refresh
// POST /api/auth/refresh
// Request body: { "refreshToken": string }
// Response: { "success": bool, "data": { "accessToken": string, "refreshToken": string, "user": User }, "error": string }
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	// Call auth service
	authResponse, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    authResponse,
		"error":   "",
	})
}

// GetProfile returns the current authenticated user's profile
// GET /api/auth/profile
// Requires: Authorization header with valid JWT token
// Response: { "success": bool, "data": User, "error": string }
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	// Extract user ID from context (set by AuthMiddleware)
	userID, ok := c.Locals("userId").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "User ID not found in context",
		})
	}

	// Get user from auth service
	user, err := h.authService.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "User not found",
		})
	}

	// Return user profile
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    user,
		"error":   "",
	})
}
