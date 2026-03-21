package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// MockRelationshipRepository is a mock implementation of RelationshipRepository
type MockRelationshipRepository struct {
	relationships []*models.Relationship
	createErr     error
	getAllErr     error
	deleteErr     error
}

func (m *MockRelationshipRepository) Create(ctx context.Context, rel *models.Relationship) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.relationships = append(m.relationships, rel)
	return nil
}

func (m *MockRelationshipRepository) GetAll(ctx context.Context) ([]*models.Relationship, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.relationships, nil
}

func (m *MockRelationshipRepository) Delete(ctx context.Context, fromID, toID, relType string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	// Remove the relationship from the mock
	for i, rel := range m.relationships {
		if rel.FromID == fromID && rel.ToID == toID && rel.Type == relType {
			m.relationships = append(m.relationships[:i], m.relationships[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockRelationshipRepository) FindPath(ctx context.Context, fromID, toID string) (*repository.RelationshipPath, error) {
	return nil, nil
}

func setupRelationshipTestApp() (*fiber.App, *MockRelationshipRepository) {
	app := fiber.New()

	// Create mock repository
	mockRepo := &MockRelationshipRepository{
		relationships: []*models.Relationship{},
	}

	// Create service and handler
	relationshipService := service.NewRelationshipService(mockRepo, nil)
	relationshipHandler := NewRelationshipHandler(relationshipService)

	// Register routes
	relationshipHandler.RegisterRoutes(app)

	return app, mockRepo
}

func TestCreateRelationship(t *testing.T) {
	app, _ := setupRelationshipTestApp()

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid relationship creation",
			requestBody: map[string]interface{}{
				"type":   "PARENT_OF",
				"fromId": "parent-id",
				"toId":   "child-id",
			},
			expectedStatus: fiber.StatusCreated,
		},
		{
			name: "Invalid relationship type",
			requestBody: map[string]interface{}{
				"type":   "INVALID_TYPE",
				"fromId": "parent-id",
				"toId":   "child-id",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "invalid relationship type",
		},
		{
			name: "Missing fromId",
			requestBody: map[string]interface{}{
				"type": "PARENT_OF",
				"toId": "child-id",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "fromId is required",
		},
		{
			name: "Self relationship",
			requestBody: map[string]interface{}{
				"type":   "PARENT_OF",
				"fromId": "same-id",
				"toId":   "same-id",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "a member cannot have a relationship with themselves",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/relationships", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != "" {
				var response map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&response)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

func TestListRelationships(t *testing.T) {
	app, mockRepo := setupRelationshipTestApp()

	// Add test relationships
	mockRepo.relationships = []*models.Relationship{
		{
			Type:      "PARENT_OF",
			FromID:    "parent-1",
			ToID:      "child-1",
			CreatedAt: time.Now(),
		},
		{
			Type:      "SPOUSE_OF",
			FromID:    "spouse-1",
			ToID:      "spouse-2",
			CreatedAt: time.Now(),
		},
	}

	req := httptest.NewRequest("GET", "/api/relationships", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])
}

func TestDeleteRelationship(t *testing.T) {
	app, mockRepo := setupRelationshipTestApp()

	// Add a test relationship
	mockRepo.relationships = []*models.Relationship{
		{
			Type:      "PARENT_OF",
			FromID:    "parent-1",
			ToID:      "child-1",
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid deletion",
			queryParams:    "?fromId=parent-1&toId=child-1&type=PARENT_OF",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Missing query parameters",
			queryParams:    "",
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Query parameters fromId, toId, and type are required",
		},
		{
			name:           "Invalid relationship type",
			queryParams:    "?fromId=parent-1&toId=child-1&type=INVALID",
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid relationship type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/relationships/dummy"+tt.queryParams, nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != "" {
				var response map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&response)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

func TestGetMemberRelationships(t *testing.T) {
	app, mockRepo := setupRelationshipTestApp()

	// Add test relationships
	mockRepo.relationships = []*models.Relationship{
		{
			Type:      "PARENT_OF",
			FromID:    "member-1",
			ToID:      "child-1",
			CreatedAt: time.Now(),
		},
		{
			Type:      "SPOUSE_OF",
			FromID:    "member-1",
			ToID:      "spouse-1",
			CreatedAt: time.Now(),
		},
		{
			Type:      "PARENT_OF",
			FromID:    "parent-1",
			ToID:      "member-2",
			CreatedAt: time.Now(),
		},
	}

	req := httptest.NewRequest("GET", "/api/members/member-1/relationships", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	relationships := data["relationships"].([]interface{})
	
	// Should return 2 relationships for member-1
	assert.Equal(t, 2, len(relationships))
}

func TestGetRelationship(t *testing.T) {
	app, _ := setupRelationshipTestApp()

	req := httptest.NewRequest("GET", "/api/relationships/some-id", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotImplemented, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.False(t, response["success"].(bool))
}

func TestUpdateRelationship(t *testing.T) {
	app, _ := setupRelationshipTestApp()

	body := map[string]interface{}{
		"type":   "PARENT_OF",
		"fromId": "parent-1",
		"toId":   "child-1",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/api/relationships/some-id", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotImplemented, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.False(t, response["success"].(bool))
}
