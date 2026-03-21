# Event Handler Integration Guide

## Overview
This guide shows how to integrate the EventHandler into the VamsaSetu backend server.

## Integration Steps

### 1. Update main.go

Add the event handler initialization and route registration to `backend/cmd/server/main.go`:

```go
package main

import (
	"context"
	"log"
	"os"

	"vamsasetu/backend/internal/config"
	"vamsasetu/backend/internal/handler"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/internal/service"
	"vamsasetu/backend/pkg/neo4j"
	"vamsasetu/backend/pkg/postgres"
	"vamsasetu/backend/pkg/redis"

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
	pgClient, neo4jClient, redisClient, err := initializeDatabases(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer pgClient.Close()
	defer neo4jClient.Close(context.Background())

	// Initialize repositories
	eventRepo := repository.NewEventRepository(pgClient.DB)
	// ... other repositories

	// Initialize services
	eventService := service.NewEventService(eventRepo, redisClient)
	// ... other services

	// Initialize handlers
	eventHandler := handler.NewEventHandler(eventService)
	// ... other handlers

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(logger.New())

	// Register routes
	eventHandler.RegisterRoutes(app)
	// ... register other handlers

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"status": "healthy",
			},
			"error": "",
		})
	})

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on :%s", port)
	log.Fatal(app.Listen(":" + port))
}

func initializeDatabases(cfg *config.Config) (*postgres.Client, *neo4j.Client, *redis.Client, error) {
	ctx := context.Background()

	// Initialize PostgreSQL client
	log.Println("Connecting to PostgreSQL...")
	pgClient, err := postgres.NewClient(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	// Verify PostgreSQL connection
	if err := pgClient.HealthCheck(ctx); err != nil {
		return nil, nil, nil, err
	}
	log.Println("PostgreSQL connection established")

	// Run GORM auto-migration
	log.Println("Running GORM auto-migration...")
	if err := runMigrations(pgClient); err != nil {
		return nil, nil, nil, err
	}
	log.Println("Database migrations completed successfully")

	// Initialize Neo4j client
	log.Println("Connecting to Neo4j...")
	neo4jClient, err := neo4j.NewClient(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	// Verify Neo4j connection
	if err := neo4jClient.HealthCheck(ctx); err != nil {
		return nil, nil, nil, err
	}
	log.Println("Neo4j connection established")

	// Initialize Redis client
	log.Println("Connecting to Redis...")
	redisClient, err := redis.NewClient()
	if err != nil {
		return nil, nil, nil, err
	}
	log.Println("Redis connection established")

	return pgClient, neo4jClient, redisClient, nil
}

func runMigrations(pgClient *postgres.Client) error {
	// Auto-migrate all models
	if err := pgClient.DB.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Notification{},
		&models.AuditLog{},
	); err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("✓ Migrated User table")
	log.Println("✓ Migrated Event table")
	log.Println("✓ Migrated Notification table")
	log.Println("✓ Migrated AuditLog table")

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

### 2. Environment Variables

Ensure the following environment variables are set in your `.env` file:

```env
# Database connections
POSTGRES_URL=postgresql://vamsasetu:vamsasetu123@localhost:5432/vamsasetu
NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=vamsasetu123
REDIS_ADDR=localhost:6379

# JWT configuration
JWT_SECRET=your-secret-key-here

# Server configuration
PORT=8080
```

### 3. Testing the Integration

#### Start the server
```bash
cd backend
go run cmd/server/main.go
```

#### Test the endpoints

1. **Create an event** (requires authentication):
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Birthday Party",
    "description": "John'\''s 30th birthday",
    "eventDate": "2024-12-25T10:00:00Z",
    "eventType": "birthday",
    "memberIds": ["member-uuid-1"]
  }'
```

2. **Get all events**:
```bash
curl -X GET http://localhost:8080/api/events \
  -H "Authorization: Bearer <your-jwt-token>"
```

3. **Get upcoming events**:
```bash
curl -X GET http://localhost:8080/api/events/upcoming?days=30 \
  -H "Authorization: Bearer <your-jwt-token>"
```

4. **Filter by type**:
```bash
curl -X GET "http://localhost:8080/api/events?type=birthday" \
  -H "Authorization: Bearer <your-jwt-token>"
```

5. **Get event by ID**:
```bash
curl -X GET http://localhost:8080/api/events/1 \
  -H "Authorization: Bearer <your-jwt-token>"
```

6. **Update event**:
```bash
curl -X PUT http://localhost:8080/api/events/1 \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Birthday Party",
    "description": "Updated description",
    "eventDate": "2024-12-26T10:00:00Z",
    "eventType": "birthday",
    "memberIds": ["member-uuid-1", "member-uuid-2"]
  }'
```

7. **Delete event**:
```bash
curl -X DELETE http://localhost:8080/api/events/1 \
  -H "Authorization: Bearer <your-jwt-token>"
```

### 4. Docker Compose Integration

The event handler works seamlessly with the existing docker-compose setup. Simply run:

```bash
docker-compose up -d
```

The backend service will automatically:
- Connect to PostgreSQL, Neo4j, and Redis
- Run database migrations
- Register all routes including event endpoints
- Start listening on port 8080

### 5. Complete Handler Registration Example

Here's how to register all handlers together:

```go
// Initialize all repositories
eventRepo := repository.NewEventRepository(pgClient.DB)
memberRepo := repository.NewMemberRepository(neo4jClient)
relationshipRepo := repository.NewRelationshipRepository(neo4jClient)
userRepo := repository.NewUserRepository(pgClient.DB)

// Initialize all services
eventService := service.NewEventService(eventRepo, redisClient)
memberService := service.NewMemberService(memberRepo, redisClient)
relationshipService := service.NewRelationshipService(relationshipRepo, redisClient)
authService := service.NewAuthService(userRepo)
treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)

// Initialize all handlers
authHandler := handler.NewAuthHandler(authService)
memberHandler := handler.NewMemberHandler(memberService)
relationshipHandler := handler.NewRelationshipHandler(relationshipService)
eventHandler := handler.NewEventHandler(eventService)
familyHandler := handler.NewFamilyHandler(treeBuilder, redisClient)

// Register all routes
authHandler.RegisterRoutes(app)
memberHandler.RegisterRoutes(app)
relationshipHandler.RegisterRoutes(app)
eventHandler.RegisterRoutes(app)
familyHandler.RegisterRoutes(app)
```

## API Endpoints Summary

| Method | Endpoint | Auth | Role | Description |
|--------|----------|------|------|-------------|
| GET | /api/events | ✓ | Any | List all events with filters |
| POST | /api/events | ✓ | owner/admin | Create new event |
| GET | /api/events/upcoming | ✓ | Any | Get upcoming events |
| GET | /api/events/:id | ✓ | Any | Get event by ID |
| PUT | /api/events/:id | ✓ | owner/admin | Update event |
| DELETE | /api/events/:id | ✓ | owner/admin | Delete event |

## Query Parameters

### GET /api/events
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50, max: 100)
- `type`: Filter by event type (birthday, anniversary, ceremony, custom)
- `member`: Filter by member ID
- `startDate`: Filter by start date (RFC3339 format)
- `endDate`: Filter by end date (RFC3339 format)

### GET /api/events/upcoming
- `days`: Number of days to look ahead (default: 30, max: 365)

## Response Format

All endpoints return responses in this format:

```json
{
  "success": true,
  "data": {
    // Response data here
  },
  "error": ""
}
```

### Success Response Example
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Birthday Party",
    "description": "John's 30th birthday",
    "eventDate": "2024-12-25T10:00:00Z",
    "eventType": "birthday",
    "memberIds": ["member-uuid-1"],
    "createdBy": 1,
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  },
  "error": ""
}
```

### Error Response Example
```json
{
  "success": false,
  "data": null,
  "error": "Event not found"
}
```

## Troubleshooting

### Common Issues

1. **"Missing authorization header"**
   - Ensure you include the Authorization header with a valid JWT token
   - Format: `Authorization: Bearer <token>`

2. **"Insufficient permissions"**
   - POST, PUT, DELETE operations require owner or admin role
   - Check your user role in the JWT token

3. **"Invalid date format"**
   - Use RFC3339 format: `2024-12-25T10:00:00Z`
   - Include timezone information

4. **"Invalid event type"**
   - Valid types: birthday, anniversary, ceremony, custom
   - Check spelling and case sensitivity

5. **Database connection errors**
   - Verify PostgreSQL is running
   - Check connection string in .env file
   - Ensure database migrations have run

## Next Steps

After integrating the event handler:

1. Test all endpoints with Postman or curl
2. Verify authentication and authorization
3. Check database records in PostgreSQL
4. Monitor Redis cache performance
5. Integrate with frontend event components
6. Set up notification scheduling (Task 8.4)
7. Implement WebSocket updates (Task 8.5)

## Additional Resources

- [Event Service Documentation](../service/EVENT_SERVICE_SUMMARY.md)
- [Event Repository Documentation](../repository/TASK_5.2_SUMMARY.md)
- [Authentication Middleware](../middleware/auth.go)
- [API Design Document](../../../.kiro/specs/vamsasetu-full-system/design.md)
