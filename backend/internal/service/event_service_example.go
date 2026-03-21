// +build ignore

package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"vamsasetu/backend/internal/config"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/pkg/postgres"
	"vamsasetu/backend/pkg/redis"
)

// ExampleEventService demonstrates how to use the EventService
func ExampleEventService() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize PostgreSQL client
	pgClient, err := postgres.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgClient.Close()

	// Initialize Redis client
	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Create repository and service
	eventRepo := repository.NewEventRepository(pgClient.DB)
	eventService := NewEventService(eventRepo, redisClient)

	ctx := context.Background()

	// Example 1: Create a new event
	fmt.Println("=== Creating a new event ===")
	newEvent := &models.Event{
		Title:       "Rajesh's Birthday",
		Description: "Birthday celebration for Rajesh",
		EventDate:   time.Now().AddDate(0, 0, 7), // 7 days from now
		EventType:   "birthday",
		MemberIDs:   []string{"member-uuid-1"},
		CreatedBy:   1, // User ID
	}

	if err := eventService.Create(ctx, newEvent); err != nil {
		log.Printf("Failed to create event: %v", err)
	} else {
		fmt.Printf("Created event with ID: %d\n", newEvent.ID)
	}

	// Example 2: Get event by ID (uses cache)
	fmt.Println("\n=== Getting event by ID ===")
	event, err := eventService.GetByID(ctx, newEvent.ID)
	if err != nil {
		log.Printf("Failed to get event: %v", err)
	} else {
		fmt.Printf("Event: %s on %s\n", event.Title, event.EventDate.Format("2006-01-02"))
	}

	// Example 3: Get upcoming events within 7 days
	fmt.Println("\n=== Getting upcoming events (7 days) ===")
	upcomingEvents, err := eventService.GetUpcoming(ctx, 7)
	if err != nil {
		log.Printf("Failed to get upcoming events: %v", err)
	} else {
		fmt.Printf("Found %d upcoming events:\n", len(upcomingEvents))
		for _, e := range upcomingEvents {
			daysUntil := int(time.Until(e.EventDate).Hours() / 24)
			fmt.Printf("  - %s (%s) in %d days\n", e.Title, e.EventType, daysUntil)
		}
	}

	// Example 4: Get events by type
	fmt.Println("\n=== Getting events by type (birthday) ===")
	birthdayEvents, err := eventService.GetByType(ctx, "birthday")
	if err != nil {
		log.Printf("Failed to get events by type: %v", err)
	} else {
		fmt.Printf("Found %d birthday events\n", len(birthdayEvents))
	}

	// Example 5: Get events by member
	fmt.Println("\n=== Getting events by member ===")
	memberEvents, err := eventService.GetByMember(ctx, "member-uuid-1")
	if err != nil {
		log.Printf("Failed to get events by member: %v", err)
	} else {
		fmt.Printf("Found %d events for member\n", len(memberEvents))
	}

	// Example 6: Get events by date range
	fmt.Println("\n=== Getting events by date range ===")
	startDate := time.Now()
	endDate := startDate.AddDate(0, 1, 0) // Next month
	rangeEvents, err := eventService.GetByDateRange(ctx, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get events by date range: %v", err)
	} else {
		fmt.Printf("Found %d events in the next month\n", len(rangeEvents))
	}

	// Example 7: Update an event
	fmt.Println("\n=== Updating an event ===")
	event.Title = "Rajesh's Birthday Party (Updated)"
	event.Description = "Updated description with more details"
	if err := eventService.Update(ctx, event); err != nil {
		log.Printf("Failed to update event: %v", err)
	} else {
		fmt.Println("Event updated successfully")
		fmt.Println("Note: Cache was automatically invalidated")
	}

	// Example 8: Delete an event
	fmt.Println("\n=== Deleting an event ===")
	if err := eventService.Delete(ctx, newEvent.ID); err != nil {
		log.Printf("Failed to delete event: %v", err)
	} else {
		fmt.Println("Event deleted successfully")
		fmt.Println("Note: All related caches were automatically invalidated")
	}

	// Example 9: Get all events
	fmt.Println("\n=== Getting all events ===")
	allEvents, err := eventService.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to get all events: %v", err)
	} else {
		fmt.Printf("Total events in system: %d\n", len(allEvents))
	}
}

// ExampleCachingBehavior demonstrates the caching behavior of EventService
func ExampleCachingBehavior() {
	fmt.Println("=== Event Service Caching Behavior ===\n")

	fmt.Println("Cache Keys Used:")
	fmt.Println("  - event:{id}              : Individual event (TTL: 5 minutes)")
	fmt.Println("  - events:upcoming:{days}  : Upcoming events list (TTL: 5 minutes)")
	fmt.Println("  - events:type:{type}      : Events by type (TTL: 5 minutes)")
	fmt.Println("  - events:member:{memberID}: Events by member (TTL: 5 minutes)")

	fmt.Println("\nCache Invalidation Triggers:")
	fmt.Println("  - Create Event: Invalidates events:upcoming:*, events:type:*, events:member:*")
	fmt.Println("  - Update Event: Invalidates event:{id}, events:upcoming:*, events:type:*, events:member:*")
	fmt.Println("  - Delete Event: Invalidates event:{id}, events:upcoming:*, events:type:*, events:member:*")

	fmt.Println("\nPerformance Benefits:")
	fmt.Println("  - First GetByID call: Queries database (~10-50ms)")
	fmt.Println("  - Subsequent GetByID calls: Returns from cache (~1-5ms)")
	fmt.Println("  - GetUpcoming with cache: ~1-5ms vs ~20-100ms without cache")
	fmt.Println("  - Filter operations benefit from caching frequently accessed queries")
}

// ExampleEventTypes demonstrates the valid event types
func ExampleEventTypes() {
	fmt.Println("=== Valid Event Types ===\n")

	eventTypes := []struct {
		Type        string
		Description string
		Example     string
	}{
		{
			Type:        "birthday",
			Description: "Birthday celebrations",
			Example:     "Rajesh's 50th Birthday",
		},
		{
			Type:        "anniversary",
			Description: "Wedding anniversaries",
			Example:     "Rajesh & Lakshmi's 25th Anniversary",
		},
		{
			Type:        "ceremony",
			Description: "Religious or cultural ceremonies",
			Example:     "Upanayanam Ceremony for Arjun",
		},
		{
			Type:        "custom",
			Description: "Custom family events",
			Example:     "Annual Family Reunion",
		},
	}

	for _, et := range eventTypes {
		fmt.Printf("Type: %s\n", et.Type)
		fmt.Printf("  Description: %s\n", et.Description)
		fmt.Printf("  Example: %s\n\n", et.Example)
	}
}
