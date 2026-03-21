# CORS and Logger Middleware Integration Guide

## Quick Start

### 1. Basic Integration in main.go

```go
package main

import (
    "log"
    "os"
    
    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/internal/middleware"
    "vamsasetu/backend/internal/handler"
    
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Create Fiber app with error handler
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.ErrorHandler(),
    })
    
    // Get frontend origin from environment
    frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
    if frontendOrigin == "" {
        frontendOrigin = "http://localhost:3000" // Default for development
    }
    
    // Apply global middleware
    app.Use(middleware.LoggerMiddleware())
    app.Use(middleware.CORSMiddleware(frontendOrigin))
    
    // Health check endpoint (no auth required)
    app.Get("/health", healthHandler)
    
    // Public API routes
    auth := app.Group("/api/auth")
    auth.Post("/register", authHandler.Register)
    auth.Post("/login", authHandler.Login)
    auth.Post("/refresh", authHandler.Refresh)
    
    // Protected API routes
    api := app.Group("/api")
    api.Use(middleware.AuthMiddleware())
    api.Use(middleware.DetailedLoggerMiddleware()) // Log with user context
    
    // Member routes
    api.Get("/members", memberHandler.GetAll)
    api.Post("/members", middleware.RequireRole("owner", "admin"), memberHandler.Create)
    api.Get("/members/:id", memberHandler.GetByID)
    api.Put("/members/:id", middleware.RequireRole("owner", "admin"), memberHandler.Update)
    api.Delete("/members/:id", middleware.RequireRole("owner", "admin"), memberHandler.Delete)
    
    // Relationship routes
    api.Get("/relationships", relationshipHandler.GetAll)
    api.Post("/relationships", middleware.RequireRole("owner", "admin"), relationshipHandler.Create)
    api.Delete("/relationships/:id", middleware.RequireRole("owner", "admin"), relationshipHandler.Delete)
    
    // Event routes
    api.Get("/events", eventHandler.GetAll)
    api.Post("/events", middleware.RequireRole("owner", "admin"), eventHandler.Create)
    api.Get("/events/:id", eventHandler.GetByID)
    api.Put("/events/:id", middleware.RequireRole("owner", "admin"), eventHandler.Update)
    api.Delete("/events/:id", middleware.RequireRole("owner", "admin"), eventHandler.Delete)
    
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    if err := app.Listen(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func healthHandler(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "healthy",
        "timestamp": time.Now(),
    })
}
```

### 2. Update .env File

Add the frontend origin configuration:

```env
# Existing configuration...
POSTGRES_URL=postgresql://vamsasetu:vamsasetu123@localhost:5432/vamsasetu
NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=vamsasetu123
REDIS_ADDR=localhost:6379
JWT_SECRET=your-secret-key-here

# CORS Configuration
FRONTEND_ORIGIN=http://localhost:3000

# Server Configuration
PORT=8080
```

### 3. Update .env.example

```env
# Database Configuration
POSTGRES_URL=postgresql://user:password@localhost:5432/vamsasetu
NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=password
REDIS_ADDR=localhost:6379

# JWT Configuration
JWT_SECRET=your-secret-key-here

# CORS Configuration
FRONTEND_ORIGIN=http://localhost:3000

# Server Configuration
PORT=8080

# Notification Services
SENDGRID_API_KEY=your-sendgrid-api-key
TWILIO_ACCOUNT_SID=your-twilio-account-sid
TWILIO_AUTH_TOKEN=your-twilio-auth-token
TWILIO_PHONE_NUMBER=your-twilio-phone-number
TWILIO_WHATSAPP_NUMBER=your-twilio-whatsapp-number
```

## Middleware Order

The order of middleware is crucial for correct behavior:

```go
// 1. Error Handler (in Fiber config)
app := fiber.New(fiber.Config{
    ErrorHandler: middleware.ErrorHandler(),
})

// 2. Logger - Log ALL requests (including failed auth)
app.Use(middleware.LoggerMiddleware())

// 3. CORS - Handle cross-origin requests
app.Use(middleware.CORSMiddleware(origin))

// 4. Public routes (no auth)
app.Post("/api/auth/login", authHandler.Login)

// 5. Protected routes
api := app.Group("/api")
api.Use(middleware.AuthMiddleware())           // Authenticate
api.Use(middleware.DetailedLoggerMiddleware()) // Log with user context
api.Use(middleware.RequireRole("owner"))       // Authorize (optional)
```

## Environment-Specific Configuration

### Development
```go
// Allow localhost origins
frontendOrigin := "http://localhost:3000"
app.Use(middleware.CORSMiddleware(frontendOrigin))
```

### Staging
```go
// Allow staging origin
frontendOrigin := "https://staging.vamsasetu.com"
app.Use(middleware.CORSMiddleware(frontendOrigin))
```

### Production
```go
// Allow production origin only
frontendOrigin := "https://vamsasetu.com"
app.Use(middleware.CORSMiddleware(frontendOrigin))
```

### Multiple Environments
```go
// Support multiple origins (comma-separated)
frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
if frontendOrigin == "" {
    // Default to all environments for development
    frontendOrigin = "http://localhost:3000,https://staging.vamsasetu.com,https://vamsasetu.com"
}
app.Use(middleware.CORSMiddleware(frontendOrigin))
```

## Advanced Usage

### Conditional Detailed Logging

Only use detailed logging for authenticated routes:

```go
// Basic logging for all requests
app.Use(middleware.LoggerMiddleware())

// Public routes
app.Post("/api/auth/login", authHandler.Login)

// Protected routes with detailed logging
api := app.Group("/api")
api.Use(middleware.AuthMiddleware())
api.Use(middleware.DetailedLoggerMiddleware()) // Includes user ID
```

### Custom Logger for Specific Routes

```go
// Global basic logger
app.Use(middleware.LoggerMiddleware())

// Admin routes with detailed logging
admin := app.Group("/api/admin")
admin.Use(middleware.AuthMiddleware())
admin.Use(middleware.RequireRole("admin"))
admin.Use(middleware.DetailedLoggerMiddleware())
```

### CORS for WebSocket

```go
// HTTP routes with CORS
app.Use(middleware.CORSMiddleware(frontendOrigin))

// WebSocket endpoint (CORS handled by upgrade)
app.Get("/ws", websocket.New(func(c *websocket.Conn) {
    // WebSocket handler
}))
```

## Testing the Integration

### 1. Test CORS

```bash
# Preflight request
curl -X OPTIONS http://localhost:8080/api/members \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type,Authorization" \
  -v

# Actual request
curl -X GET http://localhost:8080/api/members \
  -H "Origin: http://localhost:3000" \
  -H "Authorization: Bearer your-token" \
  -v
```

### 2. Test Logging

Start the server and make requests. You should see logs like:

```
2024/01/15 10:30:45 [GET] /api/members - 200 - 15.234ms
2024/01/15 10:30:46 [POST] /api/members - 201 - 45.123ms - user:1
```

### 3. Test with Frontend

```javascript
// In your React app
const response = await fetch('http://localhost:8080/api/members', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  },
  credentials: 'include', // Important for CORS with credentials
});
```

## Troubleshooting

### CORS Issues

**Problem**: "CORS policy: No 'Access-Control-Allow-Origin' header"

**Solution**: 
1. Check FRONTEND_ORIGIN environment variable
2. Ensure origin matches exactly (including protocol and port)
3. Verify CORS middleware is applied before routes

**Problem**: "CORS policy: Credentials flag is 'true', but the 'Access-Control-Allow-Credentials' header is ''"

**Solution**: 
- CORS middleware already sets `AllowCredentials: true`
- Ensure you're using the correct origin (not wildcard)

### Logging Issues

**Problem**: Logs not showing user ID

**Solution**:
1. Ensure DetailedLoggerMiddleware is after AuthMiddleware
2. Verify AuthMiddleware sets `c.Locals("userId", ...)`
3. Check that the route is protected (has AuthMiddleware)

**Problem**: Duplicate logs

**Solution**:
- Don't use both LoggerMiddleware and DetailedLoggerMiddleware on the same route
- Use LoggerMiddleware globally, DetailedLoggerMiddleware only on protected routes

## Performance Tips

1. **Logger Placement**: Place logger early in middleware chain to capture all requests
2. **CORS Caching**: Preflight responses are cached for 24 hours
3. **Conditional Logging**: Use basic logger for public routes, detailed for protected
4. **Log Levels**: Consider adding log levels (debug, info, error) for production

## Security Best Practices

1. **Exact Origins**: Never use wildcard (*) in production
2. **HTTPS Only**: Use HTTPS origins in production
3. **Credentials**: Only enable when necessary
4. **Log Sanitization**: Don't log sensitive data (passwords, tokens)
5. **Rate Limiting**: Add rate limiting middleware before CORS

## Docker Compose Integration

Update docker-compose.yml to pass FRONTEND_ORIGIN:

```yaml
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - FRONTEND_ORIGIN=http://localhost:3000
      - POSTGRES_URL=postgresql://vamsasetu:vamsasetu123@postgres:5432/vamsasetu
      # ... other env vars
```

## Monitoring and Observability

### Log Aggregation

For production, consider structured logging:

```go
// Use a structured logger like zerolog
import "github.com/rs/zerolog/log"

func StructuredLoggerMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        err := c.Next()
        duration := time.Since(start)
        
        log.Info().
            Str("method", c.Method()).
            Str("path", c.Path()).
            Int("status", c.Response().StatusCode()).
            Dur("duration", duration).
            Interface("userId", c.Locals("userId")).
            Msg("request")
        
        return err
    }
}
```

### Metrics

Add Prometheus metrics:

```go
import "github.com/gofiber/fiber/v2/middleware/monitor"

// Metrics endpoint
app.Get("/metrics", monitor.New())
```

## Summary

✅ CORS configured for frontend origin  
✅ Credentials enabled for authentication  
✅ All requests logged with method, path, status, duration  
✅ Detailed logging includes user context  
✅ Production-ready configuration  
✅ Environment-specific setup  
✅ Comprehensive testing guide  
✅ Troubleshooting documentation  

The middleware is now ready for integration into the main application!
