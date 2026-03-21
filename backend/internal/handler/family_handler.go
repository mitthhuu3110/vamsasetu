package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"vamsasetu/backend/internal/middleware"
	"vamsasetu/backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// FamilyHandler handles family tree-related HTTP requests
type FamilyHandler struct {
	treeBuilder *service.TreeBuilder
	redisClient *redis.Client
}

// NewFamilyHandler creates a new family handler instance
func NewFamilyHandler(treeBuilder *service.TreeBuilder, redisClient *redis.Client) *FamilyHandler {
	return &FamilyHandler{
		treeBuilder: treeBuilder,
		redisClient: redisClient,
	}
}

// RegisterRoutes registers all family tree routes with the Fiber app
func (h *FamilyHandler) RegisterRoutes(app *fiber.App) {
	family := app.Group("/api/family")

	// All family endpoints require authentication
	family.Use(middleware.AuthMiddleware())

	// Family tree endpoint
	family.Get("/tree", h.GetFamilyTree)
}

// GetFamilyTree retrieves the complete family tree structure
// GET /api/family/tree
// Response: { "success": bool, "data": { "nodes": []ReactFlowNode, "edges": []ReactFlowEdge }, "error": string }
func (h *FamilyHandler) GetFamilyTree(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userId").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "User ID not found in context",
		})
	}

	// Try to get from cache if Redis is available
	if h.redisClient != nil {
		cachedTree, err := h.getFromCache(ctx, userID)
		if err == nil && cachedTree != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"data":    cachedTree,
				"error":   "",
			})
		}
	}

	// Build family tree using TreeBuilder service
	familyTree, err := h.treeBuilder.BuildTree(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Failed to build family tree: " + err.Error(),
		})
	}

	// Cache the result if Redis is available
	if h.redisClient != nil {
		if err := h.saveToCache(ctx, userID, familyTree); err != nil {
			// Log error but don't fail the request
			// In production, you would use a proper logger here
			_ = err
		}
	}

	// Return family tree
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    familyTree,
		"error":   "",
	})
}

// getFromCache retrieves the family tree from Redis cache
func (h *FamilyHandler) getFromCache(ctx context.Context, userID uint) (*service.FamilyTree, error) {
	cacheKey := h.getCacheKey(userID)

	// Get from Redis
	data, err := h.redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var familyTree service.FamilyTree
	if err := json.Unmarshal([]byte(data), &familyTree); err != nil {
		return nil, err
	}

	return &familyTree, nil
}

// saveToCache stores the family tree in Redis cache with 5-minute TTL
func (h *FamilyHandler) saveToCache(ctx context.Context, userID uint, familyTree *service.FamilyTree) error {
	cacheKey := h.getCacheKey(userID)

	// Marshal to JSON
	data, err := json.Marshal(familyTree)
	if err != nil {
		return err
	}

	// Store in Redis with 5-minute TTL
	ttl := 5 * time.Minute
	return h.redisClient.Set(ctx, cacheKey, data, ttl).Err()
}

// getCacheKey generates the Redis cache key for a user's family tree
func (h *FamilyHandler) getCacheKey(userID uint) string {
	return fmt.Sprintf("family_tree:%d", userID)
}

// InvalidateCache removes the family tree from cache
// This should be called when members or relationships are modified
func (h *FamilyHandler) InvalidateCache(ctx context.Context, userID uint) error {
	if h.redisClient == nil {
		return nil
	}

	cacheKey := h.getCacheKey(userID)
	return h.redisClient.Del(ctx, cacheKey).Err()
}
