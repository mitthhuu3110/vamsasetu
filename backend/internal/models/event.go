package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// Event represents an event entity in the PostgreSQL database
type Event struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `gorm:"not null;index" json:"eventDate"`
	EventType   string    `gorm:"not null" json:"eventType"`
	MemberIDs   string    `gorm:"type:text;not null;column:member_ids" json:"-"`
	CreatedBy   uint      `gorm:"not null;index" json:"createdBy"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// GetMemberIDs returns the member IDs as a slice
func (e *Event) GetMemberIDs() []string {
	if e.MemberIDs == "" {
		return []string{}
	}
	return strings.Split(e.MemberIDs, ",")
}

// SetMemberIDs sets the member IDs from a slice
func (e *Event) SetMemberIDs(ids []string) {
	e.MemberIDs = strings.Join(ids, ",")
}

// TableName specifies the table name for the Event model
func (Event) TableName() string {
	return "events"
}

// BeforeCreate hook to validate event data before creation
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	// Additional validation can be added here if needed
	return nil
}
