package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/internal/service"
	"vamsasetu/backend/pkg/postgres"
	"vamsasetu/backend/pkg/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupEventTestApp creates a test Fiber app with event routes
func setupEventTestApp(t *testing.T) (*fiber.App, *service.EventService) {
	// Initialize test database
	pgClient, err := postgres.NewTestClient()
	require.NoError(t, err)

	// Initialize Redis client
	redisClient, err := redis.NewClient()
	require.NoError(t, err)

	// Auto-migrate models
	err = pgClient.DB.AutoMigrate(&models.User{}, &models.Event{})
	require.NoError(t, err)

	// Create test user
	testUser := &models.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		Role:         "owner",
	}
	err = pgClient.DB.Create(testUser).Error
	require.NoError(t, err)

	// Create event repository and service
	eventRepo := repository.NewEventRepository(pgClient.DB)
	eventService := service.NewEventService(eventRepo, redisClient, nil)

	// Create handler
	handler := NewEventHandler(eventService)

	// Create Fiber app
	app := fiber.New()

	// Register routes
	handler.RegisterRoutes(app)

	return app, eventService
}

// TestCreateEvent tests event creation
func TestCreateEvent(t *testing.T) {
	app, _ := setupEventTestApp(t)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid event creation",
			requestBody: map[string]interface{}{
				"title":       "Birthday Party",
				"description": "John's 30th birthday",
				"eventDate":   "2024-12-25T10:00:00Z",
				"eventType":   "birthday",
				"memberIds":   []string{"member-uuid-1", "member-uuid-2"},
			},
			expectedStatus: fiber.StatusCreated,
		},
		{
			name: "Missing title",
			requestBody: map[string]interface{}{
				"description": "Test event",
				"eventDate":   "2024-12-25T10:00:00Z",
				"eventType":   "birthday",
				"memberIds":   []string{},
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Title is required",
		},
		{
			name: "Missing event date",
			requestBody: map[string]interface{}{
				"title":       "Test Event",
				"description": "Test description",
				"eventType":   "birthday",
				"memberIds":   []string{},
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Event date is required",
		},
		{
			name: "Invalid event type",
			requestBody: map[string]interface{}{
				"title":       "Test Event",
				"description": "Test description",
				"eventDate":   "2024-12-25T10:00:00Z",
				"eventType":   "invalid_type",
				"memberIds":   []string{},
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid event type",
		},
		{
			name: "Invalid date format",
			requestBody: map[string]interface{}{
				"title":       "Test Event",
				"description": "Test description",
				"eventDate":   "2024-12-25",
				"eventType":   "birthday",
				"memberIds":   []string{},
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid date format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			// Create request
			req := httptest.NewRequest("POST", "/api/events", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid-token")

			// Mock auth middleware by setting locals
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("userId", uint(1))
				c.Locals("userRole", "owner")
				return c.Next()
			})

			// Execute request
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Parse response
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			if tt.expectedError != "" {
				assert.False(t, response["success"].(bool))
				assert.Contains(t, response["error"].(string), tt.expectedError)
			} else {
				assert.True(t, response["success"].(bool))
				assert.NotNil(t, response["data"])
			}
		})
	}
}

// TestGetEvent tests retrieving an event by ID
func TestGetEvent(t *testing.T) {
	app, eventService := setupEventTestApp(t)

	// Create a test event
	testEvent := &models.Event{
		Title:       "Test Event",
		Description: "Test Description",
		EventDate:   time.Now().Add(24 * time.Hour),
		EventType:   "birthday",
		MemberIDs:   []string{"member-1"},
		CreatedBy:   1,
	}
	err := eventService.Create(context.Background(), testEvent)
	require.NoError(t, err)

	tests := []struct {
		name           string
		eventID        string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid event ID",
			eventID:        fmt.Sprintf("%d", testEvent.ID),
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Invalid event ID",
			eventID:        "invalid",
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid event ID",
		},
		{
			name:           "Non-existent event ID",
			eventID:        "99999",
			expectedStatus: fiber.StatusNotFound,
			expectedError:  "Event not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/api/events/"+tt.eventID, nil)
			req.Header.Set("Authorization", "Bearer valid-token")

			// Mock auth middleware
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("userId", uint(1))
				c.Locals("userRole", "owner")
				return c.Next()
			})

			// Execute request
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Parse response
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			if tt.expectedError != "" {
				assert.False(t, response["success"].(bool))
				assert.Contains(t, response["error"].(string), tt.expectedError)
			} else {
				assert.True(t, response["success"].(bool))
				assert.NotNil(t, response["data"])
			}
		})
	}
}

// TestUpdateEvent tests updating an event
func TestUpdateEvent(t *testing.T) {
	app, eventService := setupEventTestApp(t)

	// Create a test event
	testEvent := &models.Event{
		Title:       "Original Title",
		Description: "Original Description",
		EventDate:   time.Now().Add(24 * time.Hour),
		EventType:   "birthday",
		MemberIDs:   []string{"member-1"},
		CreatedBy:   1,
	}
	err := eventService.Create(context.Background(), testEvent)
	require.NoError(t, err)

	tests := []struct {
		name           string
		eventID        string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "Valid update",
			eventID: fmt.Sprintf("%d", testEvent.ID),
			requestBody: map[string]interface{}{
				"title":       "Updated Title",
				"description": "Updated Description",
				"eventDate":   "2024-12-31T23:59:59Z",
				"eventType":   "anniversary",
				"memberIds":   []string{"member-1", "member-2"},
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:    "Invalid event type",
			eventID: fmt.Sprintf("%d", testEvent.ID),
			requestBody: map[string]interface{}{
				"eventType": "invalid_type",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid event type",
		},
		{
			name:           "Non-existent event",
			eventID:        "99999",
			requestBody:    map[string]interface{}{"title": "New Title"},
			expectedStatus: fiber.StatusNotFound,
			expectedError:  "Event not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			// Create request
			req := httptest.NewRequest("PUT", "/api/events/"+tt.eventID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid-token")

			// Mock auth middleware
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("userId", uint(1))
				c.Locals("userRole", "owner")
				return c.Next()
			})

			// Execute request
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Parse response
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			if tt.expectedError != "" {
				assert.False(t, response["success"].(bool))
				assert.Contains(t, response["error"].(string), tt.expectedError)
			} else {
				assert.True(t, response["success"].(bool))
				assert.NotNil(t, response["data"])
			}
		})
	}
}

// TestDeleteEvent tests deleting an event
func TestDeleteEvent(t *testing.T) {
	app, eventService := setupEventTestApp(t)

	// Create a test event
	testEvent := &models.Event{
		Title:       "Test Event",
		Description: "Test Description",
		EventDate:   time.Now().Add(24 * time.Hour),
		EventType:   "birthday",
		MemberIDs:   []string{"member-1"},
		CreatedBy:   1,
	}
	err := eventService.Create(context.Background(), testEvent)
	require.NoError(t, err)

	tests := []struct {
		name           string
		eventID        string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid deletion",
			eventID:        fmt.Sprintf("%d", testEvent.ID),
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Invalid event ID",
			eventID:        "invalid",
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid event ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("DELETE", "/api/events/"+tt.eventID, nil)
			req.Header.Set("Authorization", "Bearer valid-token")

			// Mock auth middleware
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("userId", uint(1))
				c.Locals("userRole", "owner")
				return c.Next()
			})

			// Execute request
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Parse response
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			if tt.expectedError != "" {
				assert.False(t, response["success"].(bool))
				assert.Contains(t, response["error"].(string), tt.expectedError)
			} else {
				assert.True(t, response["success"].(bool))
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "Event deleted successfully", data["message"])
			}
		})
	}
}

// TestListEvents tests listing events with filters
func TestListEvents(t *testing.T) {
	app, eventService := setupEventTestApp(t)

	// Create test events
	events := []*models.Event{
		{
			Title:       "Birthday 1",
			Description: "Test",
			EventDate:   time.Now().Add(24 * time.Hour),
			EventType:   "birthday",
			MemberIDs:   []string{"member-1"},
			CreatedBy:   1,
		},
		{
			Title:       "Anniversary 1",
			Description: "Test",
			EventDate:   time.Now().Add(48 * time.Hour),
			EventType:   "anniversary",
			MemberIDs:   []string{"member-2"},
			CreatedBy:   1,
		},
		{
			Title:       "Ceremony 1",
			Description: "Test",
			EventDate:   time.Now().Add(72 * time.Hour),
			EventType:   "ceremony",
			MemberIDs:   []string{"member-1", "member-2"},
			CreatedBy:   1,
		},
	}

	for _, event := range events {
		err := eventService.Create(context.Background(), event)
		require.NoError(t, err)
	}

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		minEvents      int
	}{
		{
			name:           "List all events",
			queryParams:    "",
			expectedStatus: fiber.StatusOK,
			minEvents:      3,
		},
		{
			name:           "Filter by type",
			queryParams:    "?type=birthday",
			expectedStatus: fiber.StatusOK,
			minEvents:      1,
		},
		{
			name:           "Filter by member",
			queryParams:    "?member=member-1",
			expectedStatus: fiber.StatusOK,
			minEvents:      2,
		},
		{
			name:           "Pagination",
			queryParams:    "?page=1&limit=2",
			expectedStatus: fiber.StatusOK,
			minEvents:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/api/events"+tt.queryParams, nil)
			req.Header.Set("Authorization", "Bearer valid-token")

			// Mock auth middleware
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("userId", uint(1))
				c.Locals("userRole", "owner")
				return c.Next()
			})

			// Execute request
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Parse response
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.True(t, response["success"].(bool))
			data := response["data"].(map[string]interface{})
			events := data["events"].([]interface{})
			assert.GreaterOrEqual(t, len(events), tt.minEvents)
		})
	}
}

// TestGetUpcomingEvents tests retrieving upcoming events
func TestGetUpcomingEvents(t *testing.T) {
	app, eventService := setupEventTestApp(t)

	// Create test events
	upcomingEvent := &models.Event{
		Title:       "Upcoming Event",
		Description: "Test",
		EventDate:   time.Now().Add(5 * 24 * time.Hour),
		EventType:   "birthday",
		MemberIDs:   []string{"member-1"},
		CreatedBy:   1,
	}
	err := eventService.Create(context.Background(), upcomingEvent)
	require.NoError(t, err)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
	}{
		{
			name:           "Default 30 days",
			queryParams:    "",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Custom days",
			queryParams:    "?days=7",
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/api/events/upcoming"+tt.queryParams, nil)
			req.Header.Set("Authorization", "Bearer valid-token")

			// Mock auth middleware
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("userId", uint(1))
				c.Locals("userRole", "owner")
				return c.Next()
			})

			// Execute request
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Parse response
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.True(t, response["success"].(bool))
			assert.NotNil(t, response["data"])
		})
	}
}
