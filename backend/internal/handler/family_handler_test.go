package handler

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/internal/service"
	"vamsasetu/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetFamilyTree_Success tests successful family tree retrieval
func TestGetFamilyTree_Success(t *testing.T) {
	// Setup
	app := fiber.New()
	
	// Create mock repositories
	memberRepo := repository.NewMemberRepository(nil) // Will use mock data
	relationshipRepo := repository.NewRelationshipRepository(nil)
	eventRepo := repository.NewEventRepository(nil)
	
	// Create tree builder
	treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)
	
	// Create handler without Redis (nil)
	handler := NewFamilyHandler(treeBuilder, nil)
	handler.RegisterRoutes(app)
	
	// Generate test JWT token
	token, err := utils.GenerateToken(1, "test@example.com", "owner")
	require.NoError(t, err)
	
	// Create request
	req := httptest.NewRequest("GET", "/api/family/tree", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	// Execute request
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	
	// Assert response
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	
	// Parse response body
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	
	// Verify response structure
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])
	assert.Equal(t, "", response["error"])
	
	// Verify data structure
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "nodes")
	assert.Contains(t, data, "edges")
}

// TestGetFamilyTree_Unauthorized tests request without authentication
func TestGetFamilyTree_Unauthorized(t *testing.T) {
	// Setup
	app := fiber.New()
	
	handler := NewFamilyHandler(nil, nil)
	handler.RegisterRoutes(app)
	
	// Create request without auth token
	req := httptest.NewRequest("GET", "/api/family/tree", nil)
	
	// Execute request
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	
	// Assert unauthorized response
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	
	// Parse response body
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	
	// Verify error response
	assert.False(t, response["success"].(bool))
	assert.Nil(t, response["data"])
	assert.NotEqual(t, "", response["error"])
}

// TestGetFamilyTree_InvalidToken tests request with invalid token
func TestGetFamilyTree_InvalidToken(t *testing.T) {
	// Setup
	app := fiber.New()
	
	handler := NewFamilyHandler(nil, nil)
	handler.RegisterRoutes(app)
	
	// Create request with invalid token
	req := httptest.NewRequest("GET", "/api/family/tree", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	
	// Execute request
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	
	// Assert unauthorized response
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// TestGetFamilyTree_EmptyTree tests retrieval of empty family tree
func TestGetFamilyTree_EmptyTree(t *testing.T) {
	// Setup
	app := fiber.New()
	
	// Create mock repositories with no data
	memberRepo := repository.NewMemberRepository(nil)
	relationshipRepo := repository.NewRelationshipRepository(nil)
	eventRepo := repository.NewEventRepository(nil)
	
	treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)
	handler := NewFamilyHandler(treeBuilder, nil)
	handler.RegisterRoutes(app)
	
	// Generate test JWT token
	token, err := utils.GenerateToken(1, "test@example.com", "owner")
	require.NoError(t, err)
	
	// Create request
	req := httptest.NewRequest("GET", "/api/family/tree", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	// Execute request
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	
	// Assert response
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	
	// Parse response body
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	
	// Verify empty tree structure
	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	nodes := data["nodes"].([]interface{})
	edges := data["edges"].([]interface{})
	
	// Empty tree should have empty arrays
	assert.Equal(t, 0, len(nodes))
	assert.Equal(t, 0, len(edges))
}

// TestGetFamilyTree_WithMembers tests tree with members but no relationships
func TestGetFamilyTree_WithMembers(t *testing.T) {
	// This test would require mock data setup
	// Skipping detailed implementation as it requires database mocking
	t.Skip("Requires database mocking")
}

// TestGetFamilyTree_WithRelationships tests tree with members and relationships
func TestGetFamilyTree_WithRelationships(t *testing.T) {
	// This test would require mock data setup
	// Skipping detailed implementation as it requires database mocking
	t.Skip("Requires database mocking")
}

// TestGetFamilyTree_WithUpcomingEvents tests tree with upcoming events
func TestGetFamilyTree_WithUpcomingEvents(t *testing.T) {
	// This test would require mock data setup
	// Skipping detailed implementation as it requires database mocking
	t.Skip("Requires database mocking")
}

// TestCacheKey tests cache key generation
func TestCacheKey(t *testing.T) {
	handler := NewFamilyHandler(nil, nil)
	
	// Test cache key generation
	key1 := handler.getCacheKey(1)
	key2 := handler.getCacheKey(2)
	key3 := handler.getCacheKey(1)
	
	// Same user ID should generate same key
	assert.Equal(t, key1, key3)
	
	// Different user IDs should generate different keys
	assert.NotEqual(t, key1, key2)
	
	// Key should have expected prefix
	assert.Contains(t, key1, "family_tree:")
}

// TestInvalidateCache tests cache invalidation
func TestInvalidateCache(t *testing.T) {
	// Test with nil Redis client (should not error)
	handler := NewFamilyHandler(nil, nil)
	err := handler.InvalidateCache(context.Background(), 1)
	assert.NoError(t, err)
}

// TestGetCacheKey tests the cache key format
func TestGetCacheKey(t *testing.T) {
	handler := NewFamilyHandler(nil, nil)
	
	tests := []struct {
		name   string
		userID uint
	}{
		{"User 1", 1},
		{"User 100", 100},
		{"User 999", 999},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := handler.getCacheKey(tt.userID)
			assert.NotEmpty(t, key)
			assert.Contains(t, key, "family_tree:")
		})
	}
}

// TestFamilyTreeResponseFormat tests the response format matches API spec
func TestFamilyTreeResponseFormat(t *testing.T) {
	// Create a sample family tree
	familyTree := &service.FamilyTree{
		Nodes: []service.ReactFlowNode{
			{
				ID:   "member-1",
				Type: "memberNode",
				Position: service.Position{
					X: 0,
					Y: 0,
				},
				Data: service.MemberNodeData{
					ID:               "member-1",
					Name:             "Test Member",
					AvatarURL:        "https://example.com/avatar.jpg",
					RelationBadge:    "Father",
					HasUpcomingEvent: true,
					Gender:           "male",
				},
			},
		},
		Edges: []service.ReactFlowEdge{
			{
				ID:       "edge-1",
				Source:   "member-1",
				Target:   "member-2",
				Type:     "bezier",
				Animated: false,
				Style: map[string]string{
					"stroke":      "#0D9488",
					"strokeWidth": "2",
				},
			},
		},
	}
	
	// Marshal to JSON
	data, err := json.Marshal(familyTree)
	require.NoError(t, err)
	
	// Unmarshal back
	var result service.FamilyTree
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)
	
	// Verify structure
	assert.Equal(t, 1, len(result.Nodes))
	assert.Equal(t, 1, len(result.Edges))
	assert.Equal(t, "member-1", result.Nodes[0].ID)
	assert.Equal(t, "memberNode", result.Nodes[0].Type)
	assert.Equal(t, "Test Member", result.Nodes[0].Data.Name)
}

// TestNodeDataStructure tests the node data structure
func TestNodeDataStructure(t *testing.T) {
	nodeData := service.MemberNodeData{
		ID:               "test-id",
		Name:             "Test Name",
		AvatarURL:        "https://example.com/avatar.jpg",
		RelationBadge:    "Father",
		HasUpcomingEvent: true,
		Gender:           "male",
	}
	
	// Verify all fields are set
	assert.Equal(t, "test-id", nodeData.ID)
	assert.Equal(t, "Test Name", nodeData.Name)
	assert.Equal(t, "https://example.com/avatar.jpg", nodeData.AvatarURL)
	assert.Equal(t, "Father", nodeData.RelationBadge)
	assert.True(t, nodeData.HasUpcomingEvent)
	assert.Equal(t, "male", nodeData.Gender)
}

// TestEdgeDataStructure tests the edge data structure
func TestEdgeDataStructure(t *testing.T) {
	edge := service.ReactFlowEdge{
		ID:       "edge-1",
		Source:   "member-1",
		Target:   "member-2",
		Type:     "bezier",
		Animated: false,
		Style: map[string]string{
			"stroke":      "#0D9488",
			"strokeWidth": "2",
		},
	}
	
	// Verify all fields are set
	assert.Equal(t, "edge-1", edge.ID)
	assert.Equal(t, "member-1", edge.Source)
	assert.Equal(t, "member-2", edge.Target)
	assert.Equal(t, "bezier", edge.Type)
	assert.False(t, edge.Animated)
	assert.Equal(t, "#0D9488", edge.Style["stroke"])
	assert.Equal(t, "2", edge.Style["strokeWidth"])
}

// TestRegisterRoutes tests route registration
func TestRegisterRoutes(t *testing.T) {
	app := fiber.New()
	handler := NewFamilyHandler(nil, nil)
	
	// Register routes
	handler.RegisterRoutes(app)
	
	// Test that route exists
	routes := app.GetRoutes()
	found := false
	for _, route := range routes {
		if route.Path == "/api/family/tree" && route.Method == "GET" {
			found = true
			break
		}
	}
	
	assert.True(t, found, "GET /api/family/tree route should be registered")
}

// TestNewFamilyHandler tests handler creation
func TestNewFamilyHandler(t *testing.T) {
	// Test with nil dependencies
	handler := NewFamilyHandler(nil, nil)
	assert.NotNil(t, handler)
	assert.Nil(t, handler.treeBuilder)
	assert.Nil(t, handler.redisClient)
	
	// Test with tree builder
	memberRepo := repository.NewMemberRepository(nil)
	relationshipRepo := repository.NewRelationshipRepository(nil)
	eventRepo := repository.NewEventRepository(nil)
	treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)
	
	handler = NewFamilyHandler(treeBuilder, nil)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.treeBuilder)
	assert.Nil(t, handler.redisClient)
}
