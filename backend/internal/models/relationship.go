package models

import (
	"errors"
	"time"
)

// Relationship type constants
const (
	RelationshipTypeSpouseOf  = "SPOUSE_OF"
	RelationshipTypeParentOf  = "PARENT_OF"
	RelationshipTypeSiblingOf = "SIBLING_OF"
)

// Relationship represents a relationship edge between two members in the Neo4j graph database
type Relationship struct {
	Type      string    `json:"type"`
	FromID    string    `json:"fromId"`
	ToID      string    `json:"toId"`
	CreatedAt time.Time `json:"createdAt"`
}

// NewRelationship creates a new Relationship
func NewRelationship(relType, fromID, toID string) *Relationship {
	return &Relationship{
		Type:      relType,
		FromID:    fromID,
		ToID:      toID,
		CreatedAt: time.Now(),
	}
}

// Validate checks if the relationship data is valid
func (r *Relationship) Validate() error {
	if r.Type == "" {
		return errors.New("relationship type is required")
	}
	if !IsValidRelationshipType(r.Type) {
		return errors.New("invalid relationship type: must be one of SPOUSE_OF, PARENT_OF, SIBLING_OF")
	}
	if r.FromID == "" {
		return errors.New("fromId is required")
	}
	if r.ToID == "" {
		return errors.New("toId is required")
	}
	if r.FromID == r.ToID {
		return errors.New("a member cannot have a relationship with themselves")
	}
	return nil
}

// IsValidRelationshipType checks if the given type is a valid relationship type
func IsValidRelationshipType(relType string) bool {
	return relType == RelationshipTypeSpouseOf ||
		relType == RelationshipTypeParentOf ||
		relType == RelationshipTypeSiblingOf
}

// IsBidirectional returns true if the relationship type is bidirectional
func (r *Relationship) IsBidirectional() bool {
	return r.Type == RelationshipTypeSpouseOf || r.Type == RelationshipTypeSiblingOf
}
