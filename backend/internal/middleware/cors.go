package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSMiddleware configures CORS for frontend origin
// Allows credentials for authentication and configures allowed methods and headers
func CORSMiddleware(frontendOrigin string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     frontendOrigin,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	})
}
