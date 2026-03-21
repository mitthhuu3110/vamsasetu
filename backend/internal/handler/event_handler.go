package handler

import (
	"strconv"
	"time"
	"vamsasetu/backend/internal/middleware"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

// EventHandler handles event-related HTTP requests
type EventHandler struct {
	eventService *service.EventService
}

// NewEventHandler creates a new event handler instance
func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

// RegisterRoutes registers all event routes with the Fiber app
func (h *EventHandler) RegisterRoutes(app *fiber.App) {
	events := app.Group("/api/events")

	// All event endpoints require authentication
	events.Use(middleware.AuthMiddleware())

	// Event CRUD endpoints
	events.Get("/", h.ListEvents)
	events.Post("/", middleware.RequireRole("owner", "admin"), h.CreateEvent)
	events.Get("/upcoming", h.GetUpcomingEvents)
	events.Get("/:id", h.GetEvent)
	events.Put("/:id", middleware.RequireRole("owner", "admin"), h.UpdateEvent)
	events.Delete("/:id", middleware.RequireRole("owner", "admin"), h.DeleteEvent)
}

// CreateEvent handles event creation
// POST /api/events
// Request body: { "title": string, "description": string, "eventDate": string, "eventType": string, "memberIds": []string }
// Response: { "success": bool, "data": Event, "error": string }
func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		EventDate   string   `json:"eventDate"`
		EventType   string   `json:"eventType"`
		MemberIDs   []string `json:"memberIds"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	// Validate required fields
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Title is required",
		})
	}

	if req.EventDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Event date is required",
		})
	}

	if req.EventType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Event type is required",
		})
	}

	// Validate event type
	validTypes := map[string]bool{
		"birthday":    true,
		"anniversary": true,
		"ceremony":    true,
		"custom":      true,
	}
	if !validTypes[req.EventType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid event type. Must be one of: birthday, anniversary, ceremony, custom",
		})
	}

	// Parse event date
	eventDate, err := time.Parse(time.RFC3339, req.EventDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid date format. Expected RFC3339 format (e.g., 2006-01-02T15:04:05Z)",
		})
	}

	// Get user ID from context
	userID, ok := c.Locals("userId").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "User ID not found in context",
		})
	}

	// Create event model
	event := &models.Event{
		Title:       req.Title,
		Description: req.Description,
		EventDate:   eventDate,
		EventType:   req.EventType,
		CreatedBy:   userID,
	}
	event.SetMemberIDs(req.MemberIDs)

	// Call event service
	if err := h.eventService.Create(c.Context(), event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    event,
		"error":   "",
	})
}

// GetEvent retrieves an event by ID
// GET /api/events/:id
// Response: { "success": bool, "data": Event, "error": string }
func (h *EventHandler) GetEvent(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Event ID is required",
		})
	}

	// Parse ID
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid event ID",
		})
	}

	// Get event from service
	event, err := h.eventService.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Event not found",
		})
	}

	// Return event
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    event,
		"error":   "",
	})
}

// UpdateEvent updates an existing event
// PUT /api/events/:id
// Request body: { "title": string, "description": string, "eventDate": string, "eventType": string, "memberIds": []string }
// Response: { "success": bool, "data": Event, "error": string }
func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Event ID is required",
		})
	}

	// Parse ID
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid event ID",
		})
	}

	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		EventDate   string   `json:"eventDate"`
		EventType   string   `json:"eventType"`
		MemberIDs   []string `json:"memberIds"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	// Get existing event
	event, err := h.eventService.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Event not found",
		})
	}

	// Update event fields
	if req.Title != "" {
		event.Title = req.Title
	}
	event.Description = req.Description

	if req.EventDate != "" {
		eventDate, err := time.Parse(time.RFC3339, req.EventDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Invalid date format. Expected RFC3339 format (e.g., 2006-01-02T15:04:05Z)",
			})
		}
		event.EventDate = eventDate
	}

	if req.EventType != "" {
		// Validate event type
		validTypes := map[string]bool{
			"birthday":    true,
			"anniversary": true,
			"ceremony":    true,
			"custom":      true,
		}
		if !validTypes[req.EventType] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Invalid event type. Must be one of: birthday, anniversary, ceremony, custom",
			})
		}
		event.EventType = req.EventType
	}

	if req.MemberIDs != nil {
		event.SetMemberIDs(req.MemberIDs)
	}

	// Call event service
	if err := h.eventService.Update(c.Context(), event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return updated event
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    event,
		"error":   "",
	})
}

// DeleteEvent deletes an event
// DELETE /api/events/:id
// Response: { "success": bool, "data": { "message": string }, "error": string }
func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Event ID is required",
		})
	}

	// Parse ID
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid event ID",
		})
	}

	// Call event service
	if err := h.eventService.Delete(c.Context(), uint(id)); err != nil {
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
			"message": "Event deleted successfully",
		},
		"error": "",
	})
}

// ListEvents retrieves all events with filters and pagination
// GET /api/events?page=1&limit=10&type=birthday&member=uuid&startDate=2024-01-01T00:00:00Z&endDate=2024-12-31T23:59:59Z
// Response: { "success": bool, "data": { "events": []Event, "total": int, "page": int, "limit": int }, "error": string }
func (h *EventHandler) ListEvents(c *fiber.Ctx) error {
	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	eventType := c.Query("type", "")
	memberID := c.Query("member", "")
	startDateStr := c.Query("startDate", "")
	endDateStr := c.Query("endDate", "")

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// Get events from service based on filters
	var events []*models.Event
	var err error

	if startDateStr != "" && endDateStr != "" {
		// Filter by date range
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Invalid startDate format. Expected RFC3339 format",
			})
		}

		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Invalid endDate format. Expected RFC3339 format",
			})
		}

		events, err = h.eventService.GetByDateRange(c.Context(), startDate, endDate)
	} else if eventType != "" {
		// Filter by event type
		events, err = h.eventService.GetByType(c.Context(), eventType)
	} else if memberID != "" {
		// Filter by member
		events, err = h.eventService.GetByMember(c.Context(), memberID)
	} else {
		// Get all events
		events, err = h.eventService.GetAll(c.Context())
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Calculate pagination
	total := len(events)
	start := (page - 1) * limit
	end := start + limit

	// Handle out of bounds
	if start >= total {
		events = []*models.Event{}
	} else {
		if end > total {
			end = total
		}
		events = events[start:end]
	}

	// Return paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"events": events,
			"total":  total,
			"page":   page,
			"limit":  limit,
		},
		"error": "",
	})
}

// GetUpcomingEvents retrieves upcoming events within the next 30 days
// GET /api/events/upcoming?days=30
// Response: { "success": bool, "data": []Event, "error": string }
func (h *EventHandler) GetUpcomingEvents(c *fiber.Ctx) error {
	// Parse days parameter
	days := c.QueryInt("days", 30)

	// Validate days parameter
	if days < 1 {
		days = 30
	}
	if days > 365 {
		days = 365
	}

	// Get upcoming events from service
	events, err := h.eventService.GetUpcoming(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return events
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    events,
		"error":   "",
	})
}
