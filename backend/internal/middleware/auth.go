package middleware

import (
	"strings"
	"vamsasetu/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens and stores user context
// Extracts Bearer token from Authorization header, validates it,
// and stores user ID and role in Fiber context locals
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Missing authorization header",
			})
		}

		// Extract Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// "Bearer " prefix not found
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Invalid authorization header format. Expected: Bearer <token>",
			})
		}

		// Validate token using JWT utilities
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Invalid or expired token",
			})
		}

		// Store user ID and role in context locals
		c.Locals("userId", claims.UserID)
		c.Locals("userRole", claims.Role)
		c.Locals("userEmail", claims.Email)

		return c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has one of the allowed roles
// Must be used after AuthMiddleware
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("userRole").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "User role not found in context",
			})
		}

		// Check if user role is in allowed roles
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Insufficient permissions",
		})
	}
}
