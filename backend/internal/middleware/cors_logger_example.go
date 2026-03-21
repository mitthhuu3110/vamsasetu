// +build ignore

package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ExampleBasicSetup demonstrates basic CORS and logger setup
func ExampleBasicSetup() {
	app := fiber.New()

	// Add logger to log all requests
	app.Use(LoggerMiddleware())

	// Add CORS for frontend
	app.Use(CORSMiddleware("http://localhost:3000"))

	// Your routes
	app.Get("/api/members", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"members": []string{}})
	})

	log.Fatal(app.Listen(":8080"))
}

// ExampleWithAuthentication demonstrates CORS and logger with auth
func ExampleWithAuthentication() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	// Global middleware
	app.Use(LoggerMiddleware())
	app.Use(CORSMiddleware("http://localhost:3000"))

	// Public routes
	app.Post("/api/auth/login", func(c *fiber.Ctx) error {
		// Login logic
		return c.JSON(fiber.Map{"token": "example-token"})
	})

	// Protected routes with detailed logging
	api := app.Group("/api")
	api.Use(AuthMiddleware())
	api.Use(DetailedLoggerMiddleware()) // Logs with user ID

	api.Get("/members", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"members": []string{}})
	})

	log.Fatal(app.Listen(":8080"))
}

// ExampleMultipleOrigins demonstrates CORS with multiple origins
func ExampleMultipleOrigins() {
	app := fiber.New()

	// Support multiple origins (dev, staging, prod)
	origins := "http://localhost:3000,https://staging.example.com,https://example.com"
	app.Use(CORSMiddleware(origins))

	app.Get("/api/data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"data": "example"})
	})

	log.Fatal(app.Listen(":8080"))
}

// ExampleEnvironmentConfig demonstrates environment-based configuration
func ExampleEnvironmentConfig() {
	app := fiber.New()

	// Get frontend origin from environment
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:3000" // Default for development
	}

	app.Use(LoggerMiddleware())
	app.Use(CORSMiddleware(frontendOrigin))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now(),
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}

// ExampleConditionalLogging demonstrates conditional detailed logging
func ExampleConditionalLogging() {
	app := fiber.New()

	// Basic logging for all requests
	app.Use(LoggerMiddleware())
	app.Use(CORSMiddleware("http://localhost:3000"))

	// Public routes (basic logging only)
	app.Post("/api/auth/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"token": "example-token"})
	})

	// Protected routes with detailed logging
	api := app.Group("/api")
	api.Use(AuthMiddleware())
	api.Use(DetailedLoggerMiddleware()) // Adds user context to logs

	api.Get("/members", func(c *fiber.Ctx) error {
		userID := c.Locals("userId")
		return c.JSON(fiber.Map{
			"userId":  userID,
			"members": []string{},
		})
	})

	log.Fatal(app.Listen(":8080"))
}

// ExampleRoleBasedRoutes demonstrates CORS and logging with role-based access
func ExampleRoleBasedRoutes() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	app.Use(LoggerMiddleware())
	app.Use(CORSMiddleware("http://localhost:3000"))

	// Public routes
	app.Post("/api/auth/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"token": "example-token"})
	})

	// Protected routes
	api := app.Group("/api")
	api.Use(AuthMiddleware())
	api.Use(DetailedLoggerMiddleware())

	// Read-only routes (all authenticated users)
	api.Get("/members", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"members": []string{}})
	})

	// Write routes (owner and admin only)
	api.Post("/members",
		RequireRole("owner", "admin"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	api.Delete("/members/:id",
		RequireRole("owner", "admin"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	log.Fatal(app.Listen(":8080"))
}

// ExampleFullStack demonstrates complete middleware stack
func ExampleFullStack() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(),
	})

	// Get configuration from environment
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:3000"
	}

	// Global middleware (order matters!)
	app.Use(LoggerMiddleware())        // 1. Log all requests
	app.Use(CORSMiddleware(frontendOrigin)) // 2. Handle CORS

	// Health check (no auth)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now(),
		})
	})

	// Public auth routes
	auth := app.Group("/api/auth")
	auth.Post("/register", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})
	auth.Post("/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"token": "example-token"})
	})

	// Protected API routes
	api := app.Group("/api")
	api.Use(AuthMiddleware())           // 3. Authenticate
	api.Use(DetailedLoggerMiddleware()) // 4. Log with user context

	// Member routes
	members := api.Group("/members")
	members.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"members": []string{}})
	})
	members.Post("/",
		RequireRole("owner", "admin"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	// Event routes
	events := api.Group("/events")
	events.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"events": []string{}})
	})
	events.Post("/",
		RequireRole("owner", "admin"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Frontend origin: %s", frontendOrigin)
	log.Fatal(app.Listen(":" + port))
}
