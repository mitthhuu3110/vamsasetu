// +build ignore

package handler

import (
	"context"
	"log"
	"time"

	"vamsasetu/backend/internal/config"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/internal/service"
	"vamsasetu/backend/pkg/postgres"
	"vamsasetu/backend/pkg/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ExampleEventHandler demonstrates how to use the EventHandler
func ExampleEventHandler() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize PostgreSQL client
	pgClient, err := postgres.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgClient.Close()

	// Initialize Redis client
	redisClient, err := redis.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Auto-migrate models
	if err := pgClient.DB.AutoMigrate(&models.User{}, &models.Event{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create event repository and service
	eventRepo := repository.NewEventRepository(pgClient.DB)
	eventService := service.NewEventService(eventRepo, redisClient)

	// Create event handler
	eventHandler := NewEventHandler(eventService)

	// Create Fiber app
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
	eventHandler.RegisterRoutes(app)

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

// ExampleCreateEvent demonstrates creating an event
func ExampleCreateEvent() {
	// Setup (same as above)
	cfg, _ := config.Load()
	pgClient, _ := postgres.NewClient(cfg)
	defer pgClient.Close()
	redisClient, _ := redis.NewClient()

	eventRepo := repository.NewEventRepository(pgClient.DB)
	eventService := service.NewEventService(eventRepo, redisClient)

	// Create an event
	event := &models.Event{
		Title:       "Rajesh's Birthday",
		Description: "Celebrating Rajesh's 50th birthday",
		EventDate:   time.Date(2024, 12, 25, 10, 0, 0, 0, time.UTC),
		EventType:   "birthday",
		MemberIDs:   []string{"member-uuid-1", "member-uuid-2"},
		CreatedBy:   1,
	}

	ctx := context.Background()
	if err := eventService.Create(ctx, event); err != nil {
		log.Fatalf("Failed to create event: %v", err)
	}

	log.Printf("Event created successfully: ID=%d, Title=%s", event.ID, event.Title)
}

// ExampleGetUpcomingEvents demonstrates retrieving upcoming events
func ExampleGetUpcomingEvents() {
	// Setup
	cfg, _ := config.Load()
	pgClient, _ := postgres.NewClient(cfg)
	defer pgClient.Close()
	redisClient, _ := redis.NewClient()

	eventRepo := repository.NewEventRepository(pgClient.DB)
	eventService := service.NewEventService(eventRepo, redisClient)

	// Get upcoming events for the next 30 days
	ctx := context.Background()
	events, err := eventService.GetUpcoming(ctx, 30)
	if err != nil {
		log.Fatalf("Failed to get upcoming events: %v", err)
	}

	log.Printf("Found %d upcoming events:", len(events))
	for _, event := range events {
		log.Printf("- %s on %s", event.Title, event.EventDate.Format("2006-01-02"))
	}
}

// ExampleFilterEventsByType demonstrates filtering events by type
func ExampleFilterEventsByType() {
	// Setup
	cfg, _ := config.Load()
	pgClient, _ := postgres.NewClient(cfg)
	defer pgClient.Close()
	redisClient, _ := redis.NewClient()

	eventRepo := repository.NewEventRepository(pgClient.DB)
	eventService := service.NewEventService(eventRepo, redisClient)

	// Get all birthday events
	ctx := context.Background()
	events, err := eventService.GetByType(ctx, "birthday")
	if err != nil {
		log.Fatalf("Failed to get birthday events: %v", err)
	}

	log.Printf("Found %d birthday events:", len(events))
	for _, event := range events {
		log.Printf("- %s on %s", event.Title, event.EventDate.Format("2006-01-02"))
	}
}

// ExampleFilterEventsByMember demonstrates filtering events by member
func ExampleFilterEventsByMember() {
	// Setup
	cfg, _ := config.Load()
	pgClient, _ := postgres.NewClient(cfg)
	defer pgClient.Close()
	redisClient, _ := redis.NewClient()

	eventRepo := repository.NewEventRepository(pgClient.DB)
	eventService := service.NewEventService(eventRepo, redisClient)

	// Get all events for a specific member
	ctx := context.Background()
	memberID := "member-uuid-1"
	events, err := eventService.GetByMember(ctx, memberID)
	if err != nil {
		log.Fatalf("Failed to get events for member: %v", err)
	}

	log.Printf("Found %d events for member %s:", len(events), memberID)
	for _, event := range events {
		log.Printf("- %s on %s", event.Title, event.EventDate.Format("2006-01-02"))
	}
}

// ExampleUpdateEvent demonstrates updating an event
func ExampleUpdateEvent() {
	// Setup
	cfg, _ := config.Load()
	pgClient, _ := postgres.NewClient(cfg)
	defer pgClient.Close()
	redisClient, _ := redis.NewClient()

	eventRepo := repository.NewEventRepository(pgClient.DB)
	eventService := service.NewEventService(eventRepo, redisClient)

	ctx := context.Background()

	// Get an existing event
	event, err := eventService.GetByID(ctx, 1)
	if err != nil {
		log.Fatalf("Failed to get event: %v", err)
	}

	// Update event details
	event.Title = "Updated Event Title"
	event.Description = "Updated description"
	event.EventDate = time.Date(2024, 12, 31, 18, 0, 0, 0, time.UTC)

	if err := eventService.Update(ctx, event); err != nil {
		log.Fatalf("Failed to update event: %v", err)
	}

	log.Printf("Event updated successfully: ID=%d, Title=%s", event.ID, event.Title)
}

// ExampleDeleteEvent demonstrates deleting an event
func ExampleDeleteEvent() {
	// Setup
	cfg, _ := config.Load()
	pgClient, _ := postgres.NewClient(cfg)
	defer pgClient.Close()
	redisClient, _ := redis.NewClient()

	eventRepo := repository.NewEventRepository(pgClient.DB)
	eventService := service.NewEventService(eventRepo, redisClient)

	ctx := context.Background()

	// Delete an event
	eventID := uint(1)
	if err := eventService.Delete(ctx, eventID); err != nil {
		log.Fatalf("Failed to delete event: %v", err)
	}

	log.Printf("Event deleted successfully: ID=%d", eventID)
}

// ExampleCompleteEventAPI demonstrates a complete event management workflow
func ExampleCompleteEventAPI() {
	// Setup
	cfg, _ := config.Load()
	pgClient, _ := postgres.NewClient(cfg)
	defer pgClient.Close()
	redisClient, _ := redis.NewClient()

	// Auto-migrate
	pgClient.DB.AutoMigrate(&models.User{}, &models.Event{})

	// Create test user
	testUser := &models.User{
		Email:        "demo@vamsasetu.com",
		PasswordHash: "hashedpassword",
		Name:         "Demo User",
		Role:         "owner",
	}
	pgClient.DB.Create(testUser)

	eventRepo := repository.NewEventRepository(pgClient.DB)
	eventService := service.NewEventService(eventRepo, redisClient)

	ctx := context.Background()

	// 1. Create multiple events
	events := []*models.Event{
		{
			Title:       "Rajesh's Birthday",
			Description: "50th birthday celebration",
			EventDate:   time.Now().Add(5 * 24 * time.Hour),
			EventType:   "birthday",
			MemberIDs:   []string{"member-1"},
			CreatedBy:   testUser.ID,
		},
		{
			Title:       "Wedding Anniversary",
			Description: "25th wedding anniversary",
			EventDate:   time.Now().Add(10 * 24 * time.Hour),
			EventType:   "anniversary",
			MemberIDs:   []string{"member-1", "member-2"},
			CreatedBy:   testUser.ID,
		},
		{
			Title:       "Diwali Celebration",
			Description: "Family Diwali gathering",
			EventDate:   time.Now().Add(15 * 24 * time.Hour),
			EventType:   "ceremony",
			MemberIDs:   []string{"member-1", "member-2", "member-3"},
			CreatedBy:   testUser.ID,
		},
	}

	for _, event := range events {
		if err := eventService.Create(ctx, event); err != nil {
			log.Printf("Failed to create event: %v", err)
			continue
		}
		log.Printf("Created event: %s (ID: %d)", event.Title, event.ID)
	}

	// 2. Get upcoming events
	upcomingEvents, _ := eventService.GetUpcoming(ctx, 30)
	log.Printf("\nUpcoming events (next 30 days): %d", len(upcomingEvents))

	// 3. Filter by type
	birthdayEvents, _ := eventService.GetByType(ctx, "birthday")
	log.Printf("Birthday events: %d", len(birthdayEvents))

	// 4. Filter by member
	memberEvents, _ := eventService.GetByMember(ctx, "member-1")
	log.Printf("Events for member-1: %d", len(memberEvents))

	// 5. Update an event
	if len(upcomingEvents) > 0 {
		event := upcomingEvents[0]
		event.Description = "Updated description"
		eventService.Update(ctx, event)
		log.Printf("Updated event: %s", event.Title)
	}

	log.Println("\nEvent management workflow completed successfully!")
}
