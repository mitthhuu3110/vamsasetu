# Database Clients

This package provides database connection clients for VamsaSetu's three data stores: Neo4j, PostgreSQL, and Redis.

## Overview

Each client:
- Accepts the `Config` struct from `internal/config`
- Initializes the respective database connection
- Provides a health check function
- Handles connection errors gracefully
- Returns the initialized client for use by repositories

## Neo4j Client

**Location:** `pkg/neo4j/client.go`

**Purpose:** Graph database for storing family members and relationships.

### Usage

```go
import (
    "context"
    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/pkg/neo4j"
)

// Initialize client
cfg, _ := config.Load()
neo4jClient, err := neo4j.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer neo4jClient.Close(context.Background())

// Health check
if err := neo4jClient.HealthCheck(context.Background()); err != nil {
    log.Printf("Neo4j unhealthy: %v", err)
}

// Use the driver
session := neo4jClient.Driver.NewSession(context.Background(), neo4j.SessionConfig{})
defer session.Close(context.Background())
```

### Configuration Requirements

- `NEO4J_URI`: Bolt connection URI (e.g., `bolt://localhost:7687`)
- `NEO4J_USERNAME`: Database username
- `NEO4J_PASSWORD`: Database password

### Features

- Automatic connectivity verification on initialization
- Context-aware operations
- Graceful error handling with descriptive messages

## PostgreSQL Client

**Location:** `pkg/postgres/client.go`

**Purpose:** Relational database for storing users, events, notifications, and audit logs.

### Usage

```go
import (
    "context"
    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/pkg/postgres"
)

// Initialize client
cfg, _ := config.Load()
postgresClient, err := postgres.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer postgresClient.Close()

// Health check
if err := postgresClient.HealthCheck(context.Background()); err != nil {
    log.Printf("PostgreSQL unhealthy: %v", err)
}

// Use GORM
var users []User
postgresClient.DB.Find(&users)
```

### Configuration Requirements

- `POSTGRES_URL`: PostgreSQL connection string (e.g., `postgresql://user:pass@localhost:5432/dbname`)

### Features

- GORM integration for ORM operations
- Connection pooling via underlying `sql.DB`
- Ping verification on initialization
- Context-aware health checks

## Redis Client

**Location:** `pkg/redis/client.go`

**Purpose:** Caching layer for performance optimization and session management.

### Usage

```go
import (
    "context"
    "time"
    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/pkg/redis"
)

// Initialize client
cfg, _ := config.Load()
redisClient, err := redis.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer redisClient.Close()

// Health check
if err := redisClient.HealthCheck(context.Background()); err != nil {
    log.Printf("Redis unhealthy: %v", err)
}

// Use Redis
ctx := context.Background()
redisClient.Client.Set(ctx, "key", "value", 5*time.Minute)
val, _ := redisClient.Client.Get(ctx, "key").Result()
```

### Configuration Requirements

- `REDIS_ADDR`: Redis server address (e.g., `localhost:6379`)

### Features

- go-redis/v9 client integration
- Automatic ping verification on initialization
- Context-aware operations
- Connection pooling

## Complete Example

See `pkg/example_usage.go` for a complete example of initializing all three clients and performing health checks.

```go
package main

import (
    "context"
    "log"
    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/pkg/neo4j"
    "vamsasetu/backend/pkg/postgres"
    "vamsasetu/backend/pkg/redis"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // Initialize all clients
    neo4jClient, _ := neo4j.NewClient(cfg)
    defer neo4jClient.Close(context.Background())

    postgresClient, _ := postgres.NewClient(cfg)
    defer postgresClient.Close()

    redisClient, _ := redis.NewClient(cfg)
    defer redisClient.Close()

    // Perform health checks
    ctx := context.Background()
    neo4jClient.HealthCheck(ctx)
    postgresClient.HealthCheck(ctx)
    redisClient.HealthCheck(ctx)
}
```

## Error Handling

All clients follow consistent error handling patterns:

1. **Configuration Validation**: Returns descriptive errors for missing required config
2. **Connection Errors**: Wraps underlying errors with context using `fmt.Errorf`
3. **Health Check Failures**: Returns errors that can be logged or used for monitoring

Example error messages:
- `"neo4j URI is required"`
- `"failed to connect to postgres: connection refused"`
- `"redis health check failed: connection timeout"`

## Testing

Each client includes unit tests for validation logic:

- `pkg/neo4j/client_test.go`
- `pkg/postgres/client_test.go`
- `pkg/redis/client_test.go`

Run tests:
```bash
go test ./pkg/neo4j -v
go test ./pkg/postgres -v
go test ./pkg/redis -v
```

## Integration with Repositories

These clients are designed to be injected into repository layers:

```go
type MemberRepository struct {
    neo4j *neo4j.Client
    cache *redis.Client
}

func NewMemberRepository(neo4jClient *neo4j.Client, cacheClient *redis.Client) *MemberRepository {
    return &MemberRepository{
        neo4j: neo4jClient,
        cache: cacheClient,
    }
}
```

## Health Check Endpoint

Use the health check functions to implement a `/health` endpoint:

```go
app.Get("/health", func(c *fiber.Ctx) error {
    ctx := context.Background()
    
    health := map[string]string{
        "neo4j":    "healthy",
        "postgres": "healthy",
        "redis":    "healthy",
    }
    
    if err := neo4jClient.HealthCheck(ctx); err != nil {
        health["neo4j"] = "unhealthy"
    }
    if err := postgresClient.HealthCheck(ctx); err != nil {
        health["postgres"] = "unhealthy"
    }
    if err := redisClient.HealthCheck(ctx); err != nil {
        health["redis"] = "unhealthy"
    }
    
    return c.JSON(health)
})
```

## Dependencies

- Neo4j: `github.com/neo4j/neo4j-go-driver/v5 v5.15.0`
- PostgreSQL: `gorm.io/gorm v1.25.5` and `gorm.io/driver/postgres v1.5.4`
- Redis: `github.com/redis/go-redis/v9 v9.4.0`

See `go.mod` for complete dependency list.
