package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Member represents a family member node in the Neo4j graph database
type Member struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Gender      string    `json:"gender"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	AvatarURL   string    `json:"avatarUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsDeleted   bool      `json:"isDeleted"`
}

// NewMember creates a new Member with a generated UUID
func NewMember(name string, dateOfBirth time.Time, gender string) *Member {
	now := time.Now()
	return &Member{
		ID:          uuid.New().String(),
		Name:        name,
		DateOfBirth: dateOfBirth,
		Gender:      gender,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsDeleted:   false,
	}
}

// Validate checks if the member data is valid
func (m *Member) Validate() error {
	if m.Name == "" {
		return errors.New("name is required")
	}
	if m.DateOfBirth.IsZero() {
		return errors.New("date of birth is required")
	}
	if m.Gender == "" {
		return errors.New("gender is required")
	}
	if m.Gender != "male" && m.Gender != "female" && m.Gender != "other" {
		return errors.New("gender must be one of: male, female, other")
	}
	if m.DateOfBirth.After(time.Now()) {
		return errors.New("date of birth cannot be in the future")
	}
	return nil
}

// SoftDelete marks the member as deleted without removing from database
func (m *Member) SoftDelete() {
	m.IsDeleted = true
	m.UpdatedAt = time.Now()
}

// Update updates the member's timestamp
func (m *Member) Update() {
	m.UpdatedAt = time.Now()
}
