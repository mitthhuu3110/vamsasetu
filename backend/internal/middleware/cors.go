package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSMiddleware configures CORS for frontend origin
// Allows credentials for authentication and configures allowed methods and headers
func CORSMiddleware(frontendOrigin string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*", // Allow all origins in development
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false, // Must be false when AllowOrigins is *
		MaxAge:           86400, // 24 hours
	})
}
