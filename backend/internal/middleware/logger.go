package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware logs all HTTP requests with method, path, status, and duration
// Format: [METHOD] /path - status - duration
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Record start time
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get status code
		status := c.Response().StatusCode()

		// Log request details
		log.Printf("[%s] %s - %d - %v",
			c.Method(),
			c.Path(),
			status,
			duration,
		)

		return err
	}
}

// DetailedLoggerMiddleware logs requests with additional details including user info
// Format: [METHOD] /path - status - duration - user:userId
func DetailedLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Record start time
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get status code
		status := c.Response().StatusCode()

		// Get user ID from context if available
		userID := c.Locals("userId")
		userInfo := ""
		if userID != nil {
			userInfo = fmt.Sprintf(" - user:%v", userID)
		}

		// Log request details
		log.Printf("[%s] %s - %d - %v%s",
			c.Method(),
			c.Path(),
			status,
			duration,
			userInfo,
		)

		return err
	}
}
