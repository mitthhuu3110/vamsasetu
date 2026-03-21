package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/pkg/redis"
)

// EventService handles event business logic with caching
type EventService struct {
	repo  *repository.EventRepository
	cache *redis.Client
	hub   WebSocketHub
}

// Cache TTL constants
const (
	EventCacheTTL         = 5 * time.Minute
	UpcomingEventsCacheTTL = 5 * time.Minute
)

// NewEventService creates a new event service instance
func NewEventService(repo *repository.EventRepository, cache *redis.Client, hub WebSocketHub) *EventService {
	return &EventService{
		repo:  repo,
		cache: cache,
		hub:   hub,
	}
}

// Create creates a new event and invalidates relevant caches
func (s *EventService) Create(ctx context.Context, event *models.Event) error {
	// Create event in database
	if err := s.repo.Create(ctx, event); err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	// Invalidate upcoming events cache
	s.invalidateUpcomingEventsCache(ctx)

	// Broadcast WebSocket update
	if s.hub != nil {
		s.hub.BroadcastUpdate("event_created", event)
	}

	return nil
}

// GetByID retrieves an event by ID with caching
func (s *EventService) GetByID(ctx context.Context, id uint) (*models.Event, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("event:%d", id)
	cached, err := s.cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var event models.Event
		if err := json.Unmarshal([]byte(cached), &event); err == nil {
			return &event, nil
		}
	}

	// Cache miss - query database
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Store in cache
	eventJSON, err := json.Marshal(event)
	if err == nil {
		s.cache.Client.Set(ctx, cacheKey, string(eventJSON), EventCacheTTL)
	}

	return event, nil
}

// GetAll retrieves all events
func (s *EventService) GetAll(ctx context.Context) ([]*models.Event, error) {
	events, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all events: %w", err)
	}

	return events, nil
}

// Update updates an existing event and invalidates relevant caches
func (s *EventService) Update(ctx context.Context, event *models.Event) error {
	// Update event in database
	if err := s.repo.Update(ctx, event); err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	// Invalidate event cache
	cacheKey := fmt.Sprintf("event:%d", event.ID)
	s.cache.Client.Del(ctx, cacheKey)

	// Invalidate upcoming events cache
	s.invalidateUpcomingEventsCache(ctx)

	// Broadcast WebSocket update
	if s.hub != nil {
		s.hub.BroadcastUpdate("event_updated", event)
	}

	return nil
}

// Delete deletes an event and invalidates relevant caches
func (s *EventService) Delete(ctx context.Context, id uint) error {
	// Delete event from database
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	// Invalidate event cache
	cacheKey := fmt.Sprintf("event:%d", id)
	s.cache.Client.Del(ctx, cacheKey)

	// Invalidate upcoming events cache
	s.invalidateUpcomingEventsCache(ctx)

	// Broadcast WebSocket update
	if s.hub != nil {
		s.hub.BroadcastUpdate("event_deleted", map[string]uint{"id": id})
	}

	return nil
}

// GetUpcoming retrieves events within the specified number of days with caching
func (s *EventService) GetUpcoming(ctx context.Context, days int) ([]*models.Event, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("events:upcoming:%d", days)
	cached, err := s.cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var events []*models.Event
		if err := json.Unmarshal([]byte(cached), &events); err == nil {
			return events, nil
		}
	}

	// Cache miss - query database
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, days)

	events, err := s.repo.GetUpcoming(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}

	// Store in cache
	eventsJSON, err := json.Marshal(events)
	if err == nil {
		s.cache.Client.Set(ctx, cacheKey, string(eventsJSON), UpcomingEventsCacheTTL)
	}

	return events, nil
}

// GetByType retrieves events filtered by event type with caching
func (s *EventService) GetByType(ctx context.Context, eventType string) ([]*models.Event, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("events:type:%s", eventType)
	cached, err := s.cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var events []*models.Event
		if err := json.Unmarshal([]byte(cached), &events); err == nil {
			return events, nil
		}
	}

	// Cache miss - query database
	events, err := s.repo.GetByType(ctx, eventType)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by type: %w", err)
	}

	// Store in cache
	eventsJSON, err := json.Marshal(events)
	if err == nil {
		s.cache.Client.Set(ctx, cacheKey, string(eventsJSON), EventCacheTTL)
	}

	return events, nil
}

// GetByMember retrieves events associated with a specific member with caching
func (s *EventService) GetByMember(ctx context.Context, memberID string) ([]*models.Event, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("events:member:%s", memberID)
	cached, err := s.cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var events []*models.Event
		if err := json.Unmarshal([]byte(cached), &events); err == nil {
			return events, nil
		}
	}

	// Cache miss - query database
	events, err := s.repo.GetByMember(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by member: %w", err)
	}

	// Store in cache
	eventsJSON, err := json.Marshal(events)
	if err == nil {
		s.cache.Client.Set(ctx, cacheKey, string(eventsJSON), EventCacheTTL)
	}

	return events, nil
}

// GetByDateRange retrieves events within a specific date range
func (s *EventService) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.Event, error) {
	events, err := s.repo.GetByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by date range: %w", err)
	}

	return events, nil
}

// invalidateUpcomingEventsCache invalidates all upcoming events cache entries
func (s *EventService) invalidateUpcomingEventsCache(ctx context.Context) {
	// Use SCAN to find all events:upcoming:* keys
	iter := s.cache.Client.Scan(ctx, 0, "events:upcoming:*", 0).Iterator()
	for iter.Next(ctx) {
		s.cache.Client.Del(ctx, iter.Val())
	}

	// Also invalidate type and member caches
	iter = s.cache.Client.Scan(ctx, 0, "events:type:*", 0).Iterator()
	for iter.Next(ctx) {
		s.cache.Client.Del(ctx, iter.Val())
	}

	iter = s.cache.Client.Scan(ctx, 0, "events:member:*", 0).Iterator()
	for iter.Next(ctx) {
		s.cache.Client.Del(ctx, iter.Val())
	}
}
