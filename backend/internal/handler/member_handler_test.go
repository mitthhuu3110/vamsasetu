package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMemberService is a mock implementation of MemberService
type MockMemberService struct {
	mock.Mock
}

func (m *MockMemberService) Create(ctx context.Context, member *models.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockMemberService) GetByID(ctx context.Context, id string) (*models.Member, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Member), args.Error(1)
}

func (m *MockMemberService) GetAll(ctx context.Context) ([]*models.Member, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Member), args.Error(1)
}

func (m *MockMemberService) Update(ctx context.Context, member *models.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockMemberService) SoftDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMemberService) Search(ctx context.Context, query string) ([]*models.Member, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Member), args.Error(1)
}

// setupTestApp creates a Fiber app with the member handler for testing
func setupTestApp(mockService *MockMemberService) *fiber.App {
	app := fiber.New()
	handler := NewMemberHandler(mockService)
	
	// Register routes without auth middleware for testing
	members := app.Group("/api/members")
	members.Get("/", handler.ListMembers)
	members.Post("/", handler.CreateMember)
	members.Get("/:id", handler.GetMember)
	members.Put("/:id", handler.UpdateMember)
	members.Delete("/:id", handler.DeleteMember)
	
	return app
}

func TestCreateMember_Success(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Prepare request
	reqBody := map[string]interface{}{
		"name":        "John Doe",
		"dateOfBirth": "1990-01-15T00:00:00Z",
		"gender":      "male",
		"email":       "john@example.com",
		"phone":       "+1234567890",
		"avatarUrl":   "https://example.com/avatar.jpg",
	}
	body, _ := json.Marshal(reqBody)

	// Mock service call
	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Member")).Return(nil)

	// Make request
	req := httptest.NewRequest("POST", "/api/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])
	mockService.AssertExpectations(t)
}

func TestCreateMember_InvalidDateFormat(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Prepare request with invalid date
	reqBody := map[string]interface{}{
		"name":        "John Doe",
		"dateOfBirth": "invalid-date",
		"gender":      "male",
	}
	body, _ := json.Marshal(reqBody)

	// Make request
	req := httptest.NewRequest("POST", "/api/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["error"].(string), "Invalid date format")
}

func TestCreateMember_ValidationError(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Prepare request
	reqBody := map[string]interface{}{
		"name":        "",
		"dateOfBirth": "1990-01-15T00:00:00Z",
		"gender":      "male",
	}
	body, _ := json.Marshal(reqBody)

	// Mock service call to return validation error
	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Member")).
		Return(errors.New("validation failed: name is required"))

	// Make request
	req := httptest.NewRequest("POST", "/api/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["error"].(string), "validation failed")
	mockService.AssertExpectations(t)
}

func TestGetMember_Success(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Create test member
	testMember := &models.Member{
		ID:          "test-id-123",
		Name:        "John Doe",
		DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		Gender:      "male",
		Email:       "john@example.com",
		Phone:       "+1234567890",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	// Mock service call
	mockService.On("GetByID", mock.Anything, "test-id-123").Return(testMember, nil)

	// Make request
	req := httptest.NewRequest("GET", "/api/members/test-id-123", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])
	mockService.AssertExpectations(t)
}

func TestGetMember_NotFound(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Mock service call to return error
	mockService.On("GetByID", mock.Anything, "nonexistent-id").
		Return(nil, errors.New("member not found"))

	// Make request
	req := httptest.NewRequest("GET", "/api/members/nonexistent-id", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Member not found", response["error"].(string))
	mockService.AssertExpectations(t)
}

func TestUpdateMember_Success(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Create existing member
	existingMember := &models.Member{
		ID:          "test-id-123",
		Name:        "John Doe",
		DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		Gender:      "male",
		Email:       "john@example.com",
		Phone:       "+1234567890",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	// Prepare update request
	reqBody := map[string]interface{}{
		"name":        "John Updated",
		"dateOfBirth": "1990-01-15T00:00:00Z",
		"gender":      "male",
		"email":       "john.updated@example.com",
		"phone":       "+9876543210",
		"avatarUrl":   "https://example.com/new-avatar.jpg",
	}
	body, _ := json.Marshal(reqBody)

	// Mock service calls
	mockService.On("GetByID", mock.Anything, "test-id-123").Return(existingMember, nil)
	mockService.On("Update", mock.Anything, mock.AnythingOfType("*models.Member")).Return(nil)

	// Make request
	req := httptest.NewRequest("PUT", "/api/members/test-id-123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])
	mockService.AssertExpectations(t)
}

func TestUpdateMember_NotFound(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Prepare update request
	reqBody := map[string]interface{}{
		"name":        "John Updated",
		"dateOfBirth": "1990-01-15T00:00:00Z",
		"gender":      "male",
	}
	body, _ := json.Marshal(reqBody)

	// Mock service call to return error
	mockService.On("GetByID", mock.Anything, "nonexistent-id").
		Return(nil, errors.New("member not found"))

	// Make request
	req := httptest.NewRequest("PUT", "/api/members/nonexistent-id", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.False(t, response["success"].(bool))
	mockService.AssertExpectations(t)
}

func TestDeleteMember_Success(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Mock service call
	mockService.On("SoftDelete", mock.Anything, "test-id-123").Return(nil)

	// Make request
	req := httptest.NewRequest("DELETE", "/api/members/test-id-123", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])
	mockService.AssertExpectations(t)
}

func TestDeleteMember_NotFound(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Mock service call to return error
	mockService.On("SoftDelete", mock.Anything, "nonexistent-id").
		Return(errors.New("member not found"))

	// Make request
	req := httptest.NewRequest("DELETE", "/api/members/nonexistent-id", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.False(t, response["success"].(bool))
	mockService.AssertExpectations(t)
}

func TestListMembers_Success(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Create test members
	testMembers := []*models.Member{
		{
			ID:          "id-1",
			Name:        "John Doe",
			DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			Gender:      "male",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "id-2",
			Name:        "Jane Smith",
			DateOfBirth: time.Date(1992, 5, 20, 0, 0, 0, 0, time.UTC),
			Gender:      "female",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Mock service call
	mockService.On("GetAll", mock.Anything).Return(testMembers, nil)

	// Make request
	req := httptest.NewRequest("GET", "/api/members", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response["success"].(bool))
	
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["members"])
	assert.Equal(t, float64(2), data["total"])
	mockService.AssertExpectations(t)
}

func TestListMembers_WithSearch(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Create test members
	testMembers := []*models.Member{
		{
			ID:          "id-1",
			Name:        "John Doe",
			DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			Gender:      "male",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Mock service call
	mockService.On("Search", mock.Anything, "John").Return(testMembers, nil)

	// Make request
	req := httptest.NewRequest("GET", "/api/members?search=John", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response["success"].(bool))
	
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["members"])
	assert.Equal(t, float64(1), data["total"])
	mockService.AssertExpectations(t)
}

func TestListMembers_WithGenderFilter(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Create test members with different genders
	testMembers := []*models.Member{
		{
			ID:          "id-1",
			Name:        "John Doe",
			DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			Gender:      "male",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "id-2",
			Name:        "Jane Smith",
			DateOfBirth: time.Date(1992, 5, 20, 0, 0, 0, 0, time.UTC),
			Gender:      "female",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Mock service call
	mockService.On("GetAll", mock.Anything).Return(testMembers, nil)

	// Make request with gender filter
	req := httptest.NewRequest("GET", "/api/members?gender=male", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response["success"].(bool))
	
	data := response["data"].(map[string]interface{})
	members := data["members"].([]interface{})
	assert.Equal(t, 1, len(members))
	assert.Equal(t, float64(1), data["total"])
	mockService.AssertExpectations(t)
}

func TestListMembers_WithPagination(t *testing.T) {
	mockService := new(MockMemberService)
	app := setupTestApp(mockService)

	// Create test members
	testMembers := make([]*models.Member, 15)
	for i := 0; i < 15; i++ {
		testMembers[i] = &models.Member{
			ID:          "id-" + string(rune(i)),
			Name:        "Member " + string(rune(i)),
			DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			Gender:      "male",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}

	// Mock service call
	mockService.On("GetAll", mock.Anything).Return(testMembers, nil)

	// Make request with pagination
	req := httptest.NewRequest("GET", "/api/members?page=2&limit=10", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response["success"].(bool))
	
	data := response["data"].(map[string]interface{})
	members := data["members"].([]interface{})
	assert.Equal(t, 5, len(members)) // 15 total, page 2 with limit 10 should have 5
	assert.Equal(t, float64(15), data["total"])
	assert.Equal(t, float64(2), data["page"])
	mockService.AssertExpectations(t)
}
