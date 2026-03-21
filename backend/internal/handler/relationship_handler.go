package handler

import (
	"vamsasetu/backend/internal/middleware"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

// RelationshipHandler handles relationship-related HTTP requests
type RelationshipHandler struct {
	relationshipService *service.RelationshipService
}

// NewRelationshipHandler creates a new relationship handler instance
func NewRelationshipHandler(relationshipService *service.RelationshipService) *RelationshipHandler {
	return &RelationshipHandler{
		relationshipService: relationshipService,
	}
}

// RegisterRoutes registers all relationship routes with the Fiber app
func (h *RelationshipHandler) RegisterRoutes(app *fiber.App) {
	relationships := app.Group("/api/relationships")

	// All relationship endpoints require authentication
	relationships.Use(middleware.AuthMiddleware())

	// Relationship CRUD endpoints
	relationships.Get("/", h.ListRelationships)
	relationships.Post("/", middleware.RequireRole("owner", "admin"), h.CreateRelationship)
	relationships.Get("/:id", h.GetRelationship)
	relationships.Put("/:id", middleware.RequireRole("owner", "admin"), h.UpdateRelationship)
	relationships.Delete("/:id", middleware.RequireRole("owner", "admin"), h.DeleteRelationship)

	// Member-specific relationships endpoint
	app.Get("/api/members/:id/relationships", middleware.AuthMiddleware(), h.GetMemberRelationships)
}

// CreateRelationship handles relationship creation
// POST /api/relationships
// Request body: { "type": string, "fromId": string, "toId": string }
// Response: { "success": bool, "data": Relationship, "error": string }
func (h *RelationshipHandler) CreateRelationship(c *fiber.Ctx) error {
	var req struct {
		Type   string `json:"type"`
		FromID string `json:"fromId"`
		ToID   string `json:"toId"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	// Create relationship model
	relationship := models.NewRelationship(req.Type, req.FromID, req.ToID)

	// Validate relationship
	if err := relationship.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Call relationship service
	if err := h.relationshipService.Create(c.Context(), relationship); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    relationship,
		"error":   "",
	})
}

// GetRelationship retrieves a relationship by ID
// GET /api/relationships/:id
// Response: { "success": bool, "data": Relationship, "error": string }
func (h *RelationshipHandler) GetRelationship(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Relationship ID is required",
		})
	}

	// Note: Neo4j relationships don't have IDs in the same way as nodes
	// This endpoint would need to be implemented differently or removed
	// For now, return a not implemented response
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"error":   "Get relationship by ID is not implemented. Use GET /api/relationships to list all relationships.",
	})
}

// UpdateRelationship updates an existing relationship
// PUT /api/relationships/:id
// Request body: { "type": string, "fromId": string, "toId": string }
// Response: { "success": bool, "data": Relationship, "error": string }
func (h *RelationshipHandler) UpdateRelationship(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Relationship ID is required",
		})
	}

	// Note: Updating relationships in Neo4j typically involves deleting and recreating
	// This endpoint would need custom implementation
	// For now, return a not implemented response
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"error":   "Update relationship is not implemented. Delete and recreate the relationship instead.",
	})
}

// DeleteRelationship deletes a relationship
// DELETE /api/relationships/:id
// Response: { "success": bool, "data": { "message": string }, "error": string }
func (h *RelationshipHandler) DeleteRelationship(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Relationship ID is required",
		})
	}

	// Parse query parameters for fromId, toId, and type
	// Since Neo4j relationships are identified by their endpoints and type
	fromID := c.Query("fromId")
	toID := c.Query("toId")
	relType := c.Query("type")

	if fromID == "" || toID == "" || relType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Query parameters fromId, toId, and type are required",
		})
	}

	// Validate relationship type
	if !models.IsValidRelationshipType(relType) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid relationship type",
		})
	}

	// Call service to delete relationship
	if err := h.relationshipService.Delete(c.Context(), fromID, toID, relType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "Relationship deleted successfully",
		},
		"error": "",
	})
}

// ListRelationships retrieves all relationships
// GET /api/relationships
// Response: { "success": bool, "data": []Relationship, "error": string }
func (h *RelationshipHandler) ListRelationships(c *fiber.Ctx) error {
	// Get all relationships
	relationships, err := h.relationshipService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return relationships
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    relationships,
		"error":   "",
	})
}

// GetMemberRelationships retrieves all relationships for a specific member
// GET /api/members/:id/relationships
// Response: { "success": bool, "data": { "relationships": []Relationship }, "error": string }
func (h *RelationshipHandler) GetMemberRelationships(c *fiber.Ctx) error {
	memberID := c.Params("id")
	if memberID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Member ID is required",
		})
	}

	// Get all relationships for the member
	relationships, err := h.relationshipService.GetByMemberID(c.Context(), memberID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return relationships
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"relationships": relationships,
		},
		"error": "",
	})
}
