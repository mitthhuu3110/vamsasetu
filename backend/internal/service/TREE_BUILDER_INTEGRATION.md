# Tree Builder Integration Guide

## Overview
This guide explains how to integrate the TreeBuilder service into the VamsaSetu backend API.

## Architecture

```
┌─────────────────┐
│  HTTP Handler   │  (Task 7.8 - To be implemented)
│  /api/family/   │
│      tree       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Cache Service  │  (Optional - Redis)
│  TTL: 5 minutes │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Tree Builder   │  ✅ Implemented (Task 7.6)
│    Service      │
└────────┬────────┘
         │
         ├──────────────┬──────────────┐
         ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│   Member     │ │ Relationship │ │    Event     │
│  Repository  │ │  Repository  │ │  Repository  │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │
       ▼                ▼                ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│    Neo4j     │ │    Neo4j     │ │  PostgreSQL  │
│   (Members)  │ │(Relationships)│ │   (Events)   │
└──────────────┘ └──────────────┘ └──────────────┘
```

## Step 1: Initialize TreeBuilder in main.go

```go
package main

import (
    "vamsasetu/backend/internal/repository"
    "vamsasetu/backend/internal/service"
    "vamsasetu/backend/pkg/neo4j"
    "gorm.io/gorm"
)

func main() {
    // ... existing initialization code ...

    // Initialize Neo4j client
    neo4jClient, err := neo4j.NewClient(cfg.Neo4jURI, cfg.Neo4jUsername, cfg.Neo4jPassword)
    if err != nil {
        log.Fatal("Failed to connect to Neo4j:", err)
    }
    defer neo4jClient.Close()

    // Initialize PostgreSQL
    db, err := postgres.NewClient(cfg.PostgresURL)
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }

    // Initialize repositories
    memberRepo := repository.NewMemberRepository(neo4jClient)
    relationshipRepo := repository.NewRelationshipRepository(neo4jClient)
    eventRepo := repository.NewEventRepository(db)

    // Initialize TreeBuilder service
    treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)

    // Pass to handler (Task 7.8)
    // familyHandler := handler.NewFamilyHandler(treeBuilder)
    
    // ... rest of initialization ...
}
```

## Step 2: Create Family Handler (Task 7.8)

**File**: `backend/internal/handler/family_handler.go`

```go
package handler

import (
    "vamsasetu/backend/internal/service"
    "github.com/gofiber/fiber/v2"
)

type FamilyHandler struct {
    treeBuilder *service.TreeBuilder
}

func NewFamilyHandler(treeBuilder *service.TreeBuilder) *FamilyHandler {
    return &FamilyHandler{
        treeBuilder: treeBuilder,
    }
}

// GetFamilyTree handles GET /api/family/tree
func (h *FamilyHandler) GetFamilyTree(c *fiber.Ctx) error {
    ctx := c.Context()

    // Build the family tree
    tree, err := h.treeBuilder.BuildTree(ctx)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "data":    nil,
            "error":   "Failed to build family tree: " + err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    tree,
        "error":   "",
    })
}
```

## Step 3: Register Route

```go
// In main.go or routes setup
app.Get("/api/family/tree", authMiddleware, familyHandler.GetFamilyTree)
```

## Step 4: Add Caching (Optional but Recommended)

**Enhanced Handler with Caching**:

```go
package handler

import (
    "encoding/json"
    "fmt"
    "time"
    
    "vamsasetu/backend/internal/service"
    "github.com/gofiber/fiber/v2"
)

type FamilyHandler struct {
    treeBuilder  *service.TreeBuilder
    cacheService *service.CacheService
}

func NewFamilyHandler(treeBuilder *service.TreeBuilder, cacheService *service.CacheService) *FamilyHandler {
    return &FamilyHandler{
        treeBuilder:  treeBuilder,
        cacheService: cacheService,
    }
}

func (h *FamilyHandler) GetFamilyTree(c *fiber.Ctx) error {
    ctx := c.Context()
    
    // Get user ID from JWT context
    userID := c.Locals("userId").(uint)
    cacheKey := fmt.Sprintf("family_tree:%d", userID)

    // Try cache first
    var cachedTree service.FamilyTree
    err := h.cacheService.Get(ctx, cacheKey, &cachedTree)
    if err == nil {
        return c.JSON(fiber.Map{
            "success": true,
            "data":    cachedTree,
            "error":   "",
        })
    }

    // Cache miss - build tree
    tree, err := h.treeBuilder.BuildTree(ctx)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "data":    nil,
            "error":   "Failed to build family tree: " + err.Error(),
        })
    }

    // Cache the result (5 minutes TTL)
    _ = h.cacheService.Set(ctx, cacheKey, tree, 5*60)

    return c.JSON(fiber.Map{
        "success": true,
        "data":    tree,
        "error":   "",
    })
}
```

## Step 5: Cache Invalidation

When members or relationships are modified, invalidate the cache:

```go
// In member handler after create/update/delete
func (h *MemberHandler) Create(c *fiber.Ctx) error {
    // ... create member logic ...
    
    // Invalidate family tree cache
    userID := c.Locals("userId").(uint)
    cacheKey := fmt.Sprintf("family_tree:%d", userID)
    h.cacheService.Delete(c.Context(), cacheKey)
    
    // ... return response ...
}

// In relationship handler after create/delete
func (h *RelationshipHandler) Create(c *fiber.Ctx) error {
    // ... create relationship logic ...
    
    // Invalidate family tree cache
    userID := c.Locals("userId").(uint)
    cacheKey := fmt.Sprintf("family_tree:%d", userID)
    h.cacheService.Delete(c.Context(), cacheKey)
    
    // ... return response ...
}
```

## API Response Format

**Endpoint**: `GET /api/family/tree`

**Headers**:
```
Authorization: Bearer <jwt_token>
```

**Response** (200 OK):
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
        "id": "edge-1",
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

**Response** (500 Internal Server Error):
```json
{
  "success": false,
  "data": null,
  "error": "Failed to build family tree: <error details>"
}
```

## Frontend Integration (React Flow)

```typescript
// In FamilyTreePage.tsx
import ReactFlow from 'reactflow';
import { useFamilyTree } from '../hooks/useFamilyTree';

export function FamilyTreePage() {
  const { data: treeData, isLoading, error } = useFamilyTree();

  if (isLoading) return <LoadingSpinner />;
  if (error) return <ErrorMessage error={error} />;

  return (
    <div className="w-full h-screen">
      <ReactFlow
        nodes={treeData.nodes}
        edges={treeData.edges}
        nodeTypes={{ memberNode: MemberNode }}
        fitView
      >
        <Background />
        <Controls />
        <MiniMap />
      </ReactFlow>
    </div>
  );
}

// In hooks/useFamilyTree.ts
import { useQuery } from '@tanstack/react-query';
import { familyService } from '../services/familyService';

export function useFamilyTree() {
  return useQuery({
    queryKey: ['familyTree'],
    queryFn: familyService.getTree,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

// In services/familyService.ts
export const familyService = {
  async getTree() {
    const response = await api.get('/api/family/tree');
    return response.data.data;
  },
};
```

## Performance Considerations

1. **Caching**: Always use Redis caching with 5-minute TTL
2. **Cache Invalidation**: Invalidate on member/relationship changes
3. **Query Optimization**: Neo4j indexes on member.id and member.name
4. **Pagination**: For very large families (>1000 members), consider pagination
5. **Lazy Loading**: Load tree incrementally (root → children → grandchildren)

## Testing

```bash
# Run unit tests
cd backend
go test ./internal/service/tree_builder_test.go ./internal/service/tree_builder.go -v

# Test API endpoint (after implementing handler)
curl -X GET http://localhost:8080/api/family/tree \
  -H "Authorization: Bearer <jwt_token>"
```

## Monitoring

Add logging for performance monitoring:

```go
func (h *FamilyHandler) GetFamilyTree(c *fiber.Ctx) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        log.Printf("Family tree generation took %v", duration)
    }()
    
    // ... handler logic ...
}
```

## Next Steps

1. ✅ Task 7.6: Tree builder implementation (COMPLETED)
2. ⏳ Task 7.7: Write property tests for tree builder
3. ⏳ Task 7.8: Implement family tree handler
4. ⏳ Integrate with caching service
5. ⏳ Add WebSocket broadcasting for real-time updates
6. ⏳ Frontend React Flow integration

## Troubleshooting

**Issue**: Tree is empty
- Check Neo4j connection
- Verify members exist in database
- Check member.isDeleted = false

**Issue**: Nodes overlap
- Verify collision detection is working
- Check NodeWidth and HorizontalSpace constants
- Increase spacing constants if needed

**Issue**: Slow performance
- Enable Redis caching
- Check Neo4j indexes
- Monitor query execution time
- Consider pagination for large families

**Issue**: Missing relationships
- Verify relationship types match constants
- Check bidirectional relationship handling
- Ensure relationships are not soft-deleted
