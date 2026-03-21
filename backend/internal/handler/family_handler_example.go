// +build ignore

package handler

// This file contains example usage of the FamilyHandler
// It is not meant to be compiled, but serves as documentation

/*
Example 1: Basic Handler Setup (No Caching)

package main

import (
    "context"
    "log"
    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/internal/handler"
    "vamsasetu/backend/internal/repository"
    "vamsasetu/backend/internal/service"
    "vamsasetu/backend/pkg/neo4j"
    "vamsasetu/backend/pkg/postgres"
    
    "github.com/gofiber/fiber/v2"
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
    
    // Create repositories
    memberRepo := repository.NewMemberRepository(neo4jClient.Driver)
    relationshipRepo := repository.NewRelationshipRepository(neo4jClient.Driver)
    eventRepo := repository.NewEventRepository(pgClient.DB)
    
    // Create TreeBuilder service
    treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)
    
    // Create FamilyHandler WITHOUT caching (pass nil for Redis)
    familyHandler := handler.NewFamilyHandler(treeBuilder, nil)
    
    // Setup Fiber app
    app := fiber.New()
    
    // Register routes
    familyHandler.RegisterRoutes(app)
    
    // Start server
    log.Fatal(app.Listen(":8080"))
}

---

Example 2: Handler Setup WITH Redis Caching

package main

import (
    "context"
    "log"
    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/internal/handler"
    "vamsasetu/backend/internal/repository"
    "vamsasetu/backend/internal/service"
    "vamsasetu/backend/pkg/neo4j"
    "vamsasetu/backend/pkg/postgres"
    "vamsasetu/backend/pkg/redis"
    
    "github.com/gofiber/fiber/v2"
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
    
    // Initialize Redis client
    redisClient, err := redis.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer redisClient.Close()
    
    // Create repositories
    memberRepo := repository.NewMemberRepository(neo4jClient.Driver)
    relationshipRepo := repository.NewRelationshipRepository(neo4jClient.Driver)
    eventRepo := repository.NewEventRepository(pgClient.DB)
    
    // Create TreeBuilder service
    treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)
    
    // Create FamilyHandler WITH caching
    familyHandler := handler.NewFamilyHandler(treeBuilder, redisClient.Client)
    
    // Setup Fiber app
    app := fiber.New()
    
    // Register routes
    familyHandler.RegisterRoutes(app)
    
    // Start server
    log.Fatal(app.Listen(":8080"))
}

---

Example 3: Cache Invalidation in Member Handler

package handler

import (
    "vamsasetu/backend/internal/middleware"
    "vamsasetu/backend/internal/service"
    
    "github.com/gofiber/fiber/v2"
)

type MemberHandler struct {
    memberService *service.MemberService
    familyHandler *FamilyHandler  // Add reference to family handler
}

func (h *MemberHandler) CreateMember(c *fiber.Ctx) error {
    // ... existing member creation logic ...
    
    // After successful member creation, invalidate family tree cache
    ctx := c.Context()
    userID := c.Locals("userId").(uint)
    
    if err := h.familyHandler.InvalidateCache(ctx, userID); err != nil {
        // Log error but don't fail the request
        log.Printf("Failed to invalidate cache: %v", err)
    }
    
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "data":    member,
        "error":   "",
    })
}

func (h *MemberHandler) UpdateMember(c *fiber.Ctx) error {
    // ... existing member update logic ...
    
    // After successful member update, invalidate family tree cache
    ctx := c.Context()
    userID := c.Locals("userId").(uint)
    h.familyHandler.InvalidateCache(ctx, userID)
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "data":    member,
        "error":   "",
    })
}

func (h *MemberHandler) DeleteMember(c *fiber.Ctx) error {
    // ... existing member deletion logic ...
    
    // After successful member deletion, invalidate family tree cache
    ctx := c.Context()
    userID := c.Locals("userId").(uint)
    h.familyHandler.InvalidateCache(ctx, userID)
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "data": fiber.Map{
            "message": "Member deleted successfully",
        },
        "error": "",
    })
}

---

Example 4: Cache Invalidation in Relationship Handler

package handler

import (
    "vamsasetu/backend/internal/middleware"
    "vamsasetu/backend/internal/service"
    
    "github.com/gofiber/fiber/v2"
)

type RelationshipHandler struct {
    relationshipService *service.RelationshipService
    familyHandler       *FamilyHandler  // Add reference to family handler
}

func (h *RelationshipHandler) CreateRelationship(c *fiber.Ctx) error {
    // ... existing relationship creation logic ...
    
    // After successful relationship creation, invalidate family tree cache
    ctx := c.Context()
    userID := c.Locals("userId").(uint)
    h.familyHandler.InvalidateCache(ctx, userID)
    
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "data":    relationship,
        "error":   "",
    })
}

func (h *RelationshipHandler) DeleteRelationship(c *fiber.Ctx) error {
    // ... existing relationship deletion logic ...
    
    // After successful relationship deletion, invalidate family tree cache
    ctx := c.Context()
    userID := c.Locals("userId").(uint)
    h.familyHandler.InvalidateCache(ctx, userID)
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "data": fiber.Map{
            "message": "Relationship deleted successfully",
        },
        "error": "",
    })
}

---

Example 5: Frontend Integration with React

import React, { useEffect, useState } from 'react';
import ReactFlow, { Node, Edge, Controls, MiniMap, Background } from 'reactflow';
import 'reactflow/dist/style.css';

interface FamilyTreeResponse {
  success: boolean;
  data: {
    nodes: Node[];
    edges: Edge[];
  };
  error: string;
}

function FamilyTreeCanvas() {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchFamilyTree();
  }, []);

  const fetchFamilyTree = async () => {
    try {
      setLoading(true);
      const token = localStorage.getItem('authToken');
      
      const response = await fetch('http://localhost:8080/api/family/tree', {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      
      const data: FamilyTreeResponse = await response.json();
      
      if (data.success) {
        setNodes(data.data.nodes);
        setEdges(data.data.edges);
        setError(null);
      } else {
        setError(data.error);
      }
    } catch (err) {
      setError('Failed to fetch family tree');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div>Loading family tree...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        fitView
      >
        <Controls />
        <MiniMap />
        <Background />
      </ReactFlow>
    </div>
  );
}

export default FamilyTreeCanvas;

---

Example 6: Testing the Endpoint with curl

# Get family tree (requires valid JWT token)
curl -X GET http://localhost:8080/api/family/tree \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Expected response:
{
  "success": true,
  "data": {
    "nodes": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "type": "memberNode",
        "position": {"x": 0, "y": 0},
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

---

Example 7: Complete Server with All Handlers

package main

import (
    "context"
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
    userRepo := repository.NewUserRepository(pgClient.DB)
    
    // Create services
    memberService := service.NewMemberService(memberRepo)
    relationshipService := service.NewRelationshipService(relationshipRepo)
    authService := service.NewAuthService(userRepo)
    treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)
    
    // Create handlers
    authHandler := handler.NewAuthHandler(authService)
    familyHandler := handler.NewFamilyHandler(treeBuilder, redisClient.Client)
    memberHandler := handler.NewMemberHandler(memberService)
    relationshipHandler := handler.NewRelationshipHandler(relationshipService)
    
    // Setup Fiber app
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
    authHandler.RegisterRoutes(app)
    familyHandler.RegisterRoutes(app)
    memberHandler.RegisterRoutes(app)
    relationshipHandler.RegisterRoutes(app)
    
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
    log.Printf("Server starting on :8080")
    log.Fatal(app.Listen(":8080"))
}

*/
