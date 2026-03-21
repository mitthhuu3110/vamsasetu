# Relationship Handler Integration Example

## How to Integrate into main.go

When you're ready to integrate the relationship handler into your Fiber application, follow this pattern:

```go
package main

import (
    "context"
    "log"
    "os"

    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/internal/handler"
    "vamsasetu/backend/internal/repository"
    "vamsasetu/backend/internal/service"
    "vamsasetu/backend/pkg/neo4j"
    "vamsasetu/backend/pkg/postgres"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
    log.Println("VamsaSetu Backend Server - Starting...")

    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    log.Println("Configuration loaded successfully")

    // Initialize database clients
    ctx := context.Background()
    
    // PostgreSQL
    pgClient, err := postgres.NewClient(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to PostgreSQL: %v", err)
    }
    defer pgClient.Close()

    // Neo4j
    neo4jClient, err := neo4j.NewClient(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to Neo4j: %v", err)
    }
    defer neo4jClient.Close(ctx)

    // Initialize Fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: customErrorHandler,
    })

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:3000",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
        AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
    }))

    // Health check endpoint
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "healthy",
            "timestamp": time.Now(),
        })
    })

    // Initialize repositories
    memberRepo := repository.NewMemberRepository(neo4jClient)
    relationshipRepo := repository.NewRelationshipRepository(neo4jClient)
    userRepo := repository.NewUserRepository(pgClient)

    // Initialize services
    memberService := service.NewMemberService(memberRepo)
    relationshipService := service.NewRelationshipService(relationshipRepo)
    authService := service.NewAuthService(userRepo)

    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService)
    memberHandler := handler.NewMemberHandler(memberService)
    relationshipHandler := handler.NewRelationshipHandler(relationshipService)

    // Register routes
    authHandler.RegisterRoutes(app)
    memberHandler.RegisterRoutes(app)
    relationshipHandler.RegisterRoutes(app)

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

func customErrorHandler(c *fiber.Ctx, err error) error {
    code := fiber.StatusInternalServerError
    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
    }

    return c.Status(code).JSON(fiber.Map{
        "success": false,
        "data":    nil,
        "error":   err.Error(),
    })
}
```

## Route Registration

The `RegisterRoutes` method in `relationship_handler.go` automatically sets up all the following routes:

```
POST   /api/relationships              - Create relationship (owner/admin only)
GET    /api/relationships              - List all relationships
GET    /api/relationships/:id          - Get relationship by ID (not implemented)
PUT    /api/relationships/:id          - Update relationship (not implemented)
DELETE /api/relationships/:id          - Delete relationship (owner/admin only)
GET    /api/members/:id/relationships  - Get member's relationships
```

## Testing the Endpoints

### 1. Create a Relationship
```bash
curl -X POST http://localhost:8080/api/relationships \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "PARENT_OF",
    "fromId": "parent-uuid",
    "toId": "child-uuid"
  }'
```

### 2. List All Relationships
```bash
curl -X GET http://localhost:8080/api/relationships \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. Delete a Relationship
```bash
curl -X DELETE "http://localhost:8080/api/relationships/dummy?fromId=parent-uuid&toId=child-uuid&type=PARENT_OF" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Get Member's Relationships
```bash
curl -X GET http://localhost:8080/api/members/member-uuid/relationships \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Dependencies

Make sure these packages are in your `go.mod`:

```go
require (
    github.com/gofiber/fiber/v2 v2.52.0
    github.com/neo4j/neo4j-go-driver/v5 v5.15.0
    github.com/google/uuid v1.5.0
)
```

## Environment Variables

Ensure these are set in your `.env` file:

```env
# Neo4j Configuration
NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=your_password

# JWT Configuration
JWT_SECRET=your_secret_key

# Server Configuration
PORT=8080
```

## Next Steps

1. Update `main.go` with the integration code above
2. Run the server: `go run cmd/server/main.go`
3. Test the endpoints using curl or Postman
4. Integrate with frontend React application
5. Add WebSocket support for real-time relationship updates (future task)

## Notes

- All relationship endpoints require authentication via JWT token
- Create and Delete operations require owner or admin role
- The handler properly handles Neo4j's bidirectional relationships
- Validation is performed at both the model and handler levels
- Error responses follow the consistent API format
