# Middleware Package

This package contains all HTTP middleware for the VamsaSetu backend API.

## Available Middleware

### 1. Authentication & Authorization
- **AuthMiddleware**: Validates JWT tokens and stores user context
- **RequireRole**: Checks if user has required role(s)

📖 See: `auth.go`

### 2. Error Handling
- **ErrorHandler**: Global error handler with consistent response format
- **AppError**: Structured error type with HTTP status codes
- Helper functions: `BadRequest`, `Unauthorized`, `Forbidden`, `NotFound`, `Conflict`, `InternalServerError`

📖 See: `error.go`, `ERROR_HANDLER_GUIDE.md`

### 3. CORS
- **CORSMiddleware**: Configures CORS for frontend origin
  - Allows credentials for authentication
  - Supports multiple origins
  - 24-hour preflight cache

📖 See: `cors.go`, `CORS_LOGGER_INTEGRATION.md`

### 4. Logging
- **LoggerMiddleware**: Basic request logging (method, path, status, duration)
- **DetailedLoggerMiddleware**: Enhanced logging with user context

📖 See: `logger.go`, `CORS_LOGGER_INTEGRATION.md`

## Quick Start

```go
import (
    "vamsasetu/backend/internal/middleware"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Create app with error handler
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.ErrorHandler(),
    })
    
    // Global middleware
    app.Use(middleware.LoggerMiddleware())
    app.Use(middleware.CORSMiddleware("http://localhost:3000"))
    
    // Public routes
    app.Post("/api/auth/login", authHandler.Login)
    
    // Protected routes
    api := app.Group("/api")
    api.Use(middleware.AuthMiddleware())
    api.Use(middleware.DetailedLoggerMiddleware())
    
    api.Get("/members", memberHandler.GetAll)
    api.Post("/members", 
        middleware.RequireRole("owner", "admin"),
        memberHandler.Create,
    )
    
    app.Listen(":8080")
}
```

## Middleware Order

The order of middleware is important:

1. **Error Handler** (in Fiber config)
2. **Logger** - Log all requests
3. **CORS** - Handle cross-origin requests
4. **Auth** - Authenticate user (protected routes only)
5. **Detailed Logger** - Log with user context (protected routes only)
6. **Role Check** - Authorize user (specific routes only)

## Testing

All middleware have comprehensive test coverage:

```bash
# Run all middleware tests
go test -v ./internal/middleware/

# Run specific middleware tests
go test -v ./internal/middleware/auth_test.go
go test -v ./internal/middleware/error_test.go
go test -v ./internal/middleware/cors_test.go
go test -v ./internal/middleware/logger_test.go

# Run with coverage
go test -v -coverprofile=coverage.out ./internal/middleware/
go tool cover -html=coverage.out
```

## Documentation

- **INTEGRATION_EXAMPLE.md**: General middleware integration guide
- **ERROR_HANDLER_GUIDE.md**: Error handling patterns and best practices
- **CORS_LOGGER_INTEGRATION.md**: CORS and logging setup guide
- **TASK_10.5_SUMMARY.md**: Error handler implementation summary
- **TASK_10.7_SUMMARY.md**: CORS and logger implementation summary

## Examples

See `cors_logger_example.go` and `error_example.go` for complete working examples.

## Environment Variables

```env
# CORS Configuration
FRONTEND_ORIGIN=http://localhost:3000

# JWT Configuration (for auth middleware)
JWT_SECRET=your-secret-key-here

# Server Configuration
PORT=8080
```

## Features

✅ JWT authentication with role-based access control  
✅ Consistent error handling with APIResponse format  
✅ CORS with credentials support  
✅ Request logging with duration tracking  
✅ User context logging for authenticated requests  
✅ Comprehensive test coverage  
✅ Production-ready implementations  

## Next Steps

1. Integrate middleware into `cmd/server/main.go`
2. Add FRONTEND_ORIGIN to `.env` file
3. Configure error handler in Fiber app
4. Apply middleware to routes
5. Test CORS with frontend
6. Monitor logs for debugging

## Support

For questions or issues:
- Check the documentation files in this directory
- Review the example files
- Run the tests to verify behavior
- See the integration guides for setup instructions
