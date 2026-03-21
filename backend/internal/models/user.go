package models

import (
	"time"
)

// User represents a user entity in the PostgreSQL database
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"unique;not null;index" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Name         string    `gorm:"not null" json:"name"`
	Role         string    `gorm:"not null" json:"role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
