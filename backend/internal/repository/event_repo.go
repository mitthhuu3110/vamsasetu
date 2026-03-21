package repository

import (
	"context"
	"fmt"
	"time"

	"vamsasetu/backend/internal/models"

	"gorm.io/gorm"
)

// EventRepository handles event data operations in PostgreSQL
type EventRepository struct {
	db *gorm.DB
}

// NewEventRepository creates a new event repository instance
func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

// Create creates a new event in PostgreSQL
func (r *EventRepository) Create(ctx context.Context, event *models.Event) error {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

// GetByID retrieves an event by ID from PostgreSQL
func (r *EventRepository) GetByID(ctx context.Context, id uint) (*models.Event, error) {
	var event models.Event
	if err := r.db.WithContext(ctx).First(&event, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to query event: %w", err)
	}
	return &event, nil
}

// GetAll retrieves all events from PostgreSQL
func (r *EventRepository) GetAll(ctx context.Context) ([]*models.Event, error) {
	var events []*models.Event
	if err := r.db.WithContext(ctx).Order("event_date ASC").Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	return events, nil
}

// Update updates an existing event in PostgreSQL
func (r *EventRepository) Update(ctx context.Context, event *models.Event) error {
	if err := r.db.WithContext(ctx).Save(event).Error; err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	return nil
}

// Delete deletes an event from PostgreSQL
func (r *EventRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Event{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}

// GetUpcoming retrieves events within a date range
func (r *EventRepository) GetUpcoming(ctx context.Context, startDate, endDate time.Time) ([]*models.Event, error) {
	var events []*models.Event
	if err := r.db.WithContext(ctx).
		Where("event_date >= ? AND event_date <= ?", startDate, endDate).
		Order("event_date ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query upcoming events: %w", err)
	}
	return events, nil
}

// GetByType retrieves events filtered by event type
func (r *EventRepository) GetByType(ctx context.Context, eventType string) ([]*models.Event, error) {
	var events []*models.Event
	if err := r.db.WithContext(ctx).
		Where("event_type = ?", eventType).
		Order("event_date ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query events by type: %w", err)
	}
	return events, nil
}

// GetByMember retrieves events associated with a specific member
func (r *EventRepository) GetByMember(ctx context.Context, memberID string) ([]*models.Event, error) {
	var events []*models.Event
	// PostgreSQL array contains operator
	if err := r.db.WithContext(ctx).
		Where("? = ANY(member_ids)", memberID).
		Order("event_date ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query events by member: %w", err)
	}
	return events, nil
}

// GetByDateRange retrieves events within a specific date range
func (r *EventRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.Event, error) {
	var events []*models.Event
	if err := r.db.WithContext(ctx).
		Where("event_date >= ? AND event_date <= ?", startDate, endDate).
		Order("event_date ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query events by date range: %w", err)
	}
	return events, nil
}

// GetByCreator retrieves events created by a specific user
func (r *EventRepository) GetByCreator(ctx context.Context, userID uint) ([]*models.Event, error) {
	var events []*models.Event
	if err := r.db.WithContext(ctx).
		Where("created_by = ?", userID).
		Order("event_date ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query events by creator: %w", err)
	}
	return events, nil
}
