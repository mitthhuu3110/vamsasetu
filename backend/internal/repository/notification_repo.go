package repository

import (
	"context"
	"fmt"
	"time"

	"vamsasetu/backend/internal/models"

	"gorm.io/gorm"
)

// NotificationRepository handles notification data operations in PostgreSQL
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository instance
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

// Create creates a new notification in PostgreSQL
func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

// GetPending retrieves all pending notifications that are due to be sent
func (r *NotificationRepository) GetPending(ctx context.Context) ([]*models.Notification, error) {
	var notifications []*models.Notification
	if err := r.db.WithContext(ctx).
		Where("status = ? AND scheduled_at <= ?", "pending", time.Now()).
		Order("scheduled_at ASC").
		Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to query pending notifications: %w", err)
	}
	return notifications, nil
}

// UpdateStatus updates the status of a notification
func (r *NotificationRepository) UpdateStatus(ctx context.Context, id uint, status string, sentAt *time.Time, errorMsg string) error {
	updates := map[string]interface{}{
		"status":     status,
		"error_msg":  errorMsg,
		"updated_at": time.Now(),
	}
	
	if sentAt != nil {
		updates["sent_at"] = sentAt
	}
	
	if err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}
	return nil
}

// IncrementRetry increments the retry count for a failed notification
func (r *NotificationRepository) IncrementRetry(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment retry count: %w", err)
	}
	return nil
}
