package models

import (
	"time"

	"gorm.io/gorm"
)

// Notification represents a notification entity in the PostgreSQL database
type Notification struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	EventID     uint       `gorm:"not null;index" json:"eventId"`
	UserID      uint       `gorm:"not null" json:"userId"`
	Channel     string     `gorm:"not null" json:"channel"`
	ScheduledAt time.Time  `gorm:"not null;index:idx_notifications_scheduled" json:"scheduledAt"`
	SentAt      *time.Time `json:"sentAt"`
	Status      string     `gorm:"not null;index:idx_notifications_scheduled" json:"status"`
	RetryCount  int        `gorm:"default:0" json:"retryCount"`
	ErrorMsg    string     `json:"errorMsg"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName specifies the table name for the Notification model
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate hook to set default values
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.Status == "" {
		n.Status = "pending"
	}
	return nil
}
