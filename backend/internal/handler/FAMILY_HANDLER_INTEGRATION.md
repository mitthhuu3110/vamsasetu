# Family Handler Integration Guide

## Overview

The `FamilyHandler` provides the `/api/family/tree` endpoint that returns the complete family tree structure in React Flow compatible format. It integrates with the `TreeBuilder` service and supports optional Redis caching.

## Handler Creation

```go
import (
    "vamsasetu/backend/internal/handler"
    "vamsasetu/backend/internal/repository"
    "vamsasetu/backend/internal/service"
    "vamsasetu/backend/pkg/neo4j"
    "vamsasetu/backend/pkg/postgres"
    "vamsasetu/backend/pkg/redis"
)

// Initialize database clients
neo4jClient, err := neo4j.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}

pgClient, err := postgres.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}

redisClient, err := redis.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}

// Create repositories
memberRepo := repository.NewMemberRepository(neo4jClient.Driver)
relationshipRepo := repository.NewRelationshipRepository(neo4jClient.Driver)
eventRepo := repository.NewEventRepository(pgClient.DB)

// Create TreeBuilder service
treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)

// Create FamilyHandler with Redis caching
familyHandler := handler.NewFamilyHandler(treeBuilder, redisClient.Client)

// Or without caching (pass nil for Redis)
familyHandler := handler.NewFamilyHandler(treeBuilder, nil)
```

## Route Registration

```go
import (
    "github.com/gofiber/fiber/v2"
)

app := fiber.New()

// Register family routes
familyHandler.RegisterRoutes(app)

// This registers:
// GET /api/family/tree (requires authentication)
```

## Complete Server Setup Example

```go
package main

import (
    "log"
    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/internal/handler"
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
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }
    
    // Initialize database clients
    neo4jClient, err := neo4j.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer neo4jClient.Close(context.Background())
    
    pgClient, err := postgres.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer pgClient.Close()
    
    redisClient, err := redis.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer redisClient.Close()
    
    // Create repositories
    memberRepo := repository.NewMemberRepository(neo4jClient.Driver)
    relationshipRepo := repository.NewRelationshipRepository(neo4jClient.Driver)
    eventRepo := repository.NewEventRepository(pgClient.DB)
    
    // Create services
    treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)
    
    // Create handlers
    familyHandler := handler.NewFamilyHandler(treeBuilder, redisClient.Client)
    
    // Setup Fiber app
    app := fiber.New()
    app.Use(cors.New())
    app.Use(logger.New())
    
    // Register routes
    familyHandler.RegisterRoutes(app)
    
    // Start server
    log.Fatal(app.Listen(":8080"))
}
```

## API Endpoint

### GET /api/family/tree

Returns the complete family tree structure with nodes and edges in React Flow format.

**Authentication:** Required (JWT Bearer token)

**Request:**
```http
GET /api/family/tree HTTP/1.1
Host: localhost:8080
Authorization: Bearer <jwt-token>
```

**Response (Success):**
```json
{
  "success": true,
  "data": {
    "nodes": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "type": "memberNode",
        "position": {
          "x": 0,
          "y": 0
        },
        "data": {
          "id": "550e8400-e29b-41d4-a716-446655440000",
          "name": "Rajesh Kumar",
          "avatarUrl": "https://example.com/avatar.jpg",
          "relationBadge": "Father",
          "hasUpcomingEvent": true,
          "gender": "male"
        }
      }
    ],
    "edges": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000-660e8400-e29b-41d4-a716-446655440001-PARENT_OF",
        "source": "550e8400-e29b-41d4-a716-446655440000",
        "target": "660e8400-e29b-41d4-a716-446655440001",
        "type": "bezier",
        "animated": false,
        "style": {
          "stroke": "#0D9488",
          "strokeWidth": "2"
        }
      }
    ]
  },
  "error": ""
}
```

**Response (Unauthorized):**
```json
{
  "success": false,
  "data": null,
  "error": "Missing authorization header"
}
```

**Response (Error):**
```json
{
  "success": false,
  "data": null,
  "error": "Failed to build family tree: <error details>"
}
```

## Caching Behavior

When Redis is configured:

1. **Cache Hit:** Returns cached tree within ~50ms
2. **Cache Miss:** Builds tree from database, caches result with 5-minute TTL
3. **Cache Key Format:** `family_tree:{userId}`
4. **TTL:** 5 minutes (300 seconds)

### Cache Invalidation

The handler provides an `InvalidateCache` method that should be called when members or relationships are modified:

```go
// In member handler after create/update/delete
ctx := c.Context()
userID := c.Locals("userId").(uint)
familyHandler.InvalidateCache(ctx, userID)

// In relationship handler after create/delete
ctx := c.Context()
userID := c.Locals("userId").(uint)
familyHandler.InvalidateCache(ctx, userID)
```

## Frontend Integration

### React Flow Setup

```typescript
import ReactFlow, { Node, Edge } from 'reactflow';
import 'reactflow/dist/style.css';

interface FamilyTreeResponse {
  success: boolean;
  data: {
    nodes: Node[];
    edges: Edge[];
  };
  error: string;
}

async function fetchFamilyTree(): Promise<FamilyTreeResponse> {
  const response = await fetch('http://localhost:8080/api/family/tree', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  return response.json();
}

function FamilyTreeCanvas() {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  useEffect(() => {
    fetchFamilyTree().then(response => {
      if (response.success) {
        setNodes(response.data.nodes);
        setEdges(response.data.edges);
      }
    });
  }, []);

  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      fitView
    />
  );
}
```

## Testing

Run the handler tests:

```bash
cd backend
go test -v ./internal/handler -run TestGetFamilyTree
```

## Performance Considerations

1. **Without Cache:** ~200-500ms for medium trees (50-100 members)
2. **With Cache:** ~50ms for cached responses
3. **Cache Invalidation:** Automatic on member/relationship changes
4. **Concurrent Requests:** Safe with Redis caching

## Error Handling

The handler returns appropriate HTTP status codes:

- `200 OK`: Successful tree retrieval
- `401 Unauthorized`: Missing or invalid JWT token
- `500 Internal Server Error`: Database or service errors

All errors follow the consistent API response format with `success: false` and descriptive error messages.

## Security

- **Authentication Required:** All endpoints require valid JWT token
- **User Context:** User ID extracted from JWT claims
- **Authorization:** Future enhancement could add family-level permissions

## Future Enhancements

1. **Pagination:** For very large trees (1000+ members)
2. **Filtering:** By generation, branch, or relationship type
3. **Partial Updates:** WebSocket support for real-time tree updates
4. **Export:** PDF/PNG export of tree visualization
5. **Permissions:** Family-level access control
