package repository

import (
	"context"
	"fmt"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/pkg/neo4j"

	neo4jDriver "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// MemberRepository handles member data operations in Neo4j
type MemberRepository struct {
	client *neo4j.Client
}

// NewMemberRepository creates a new member repository instance
func NewMemberRepository(client *neo4j.Client) *MemberRepository {
	return &MemberRepository{
		client: client,
	}
}

// Create creates a new member node in Neo4j
func (r *MemberRepository) Create(ctx context.Context, member *models.Member) error {
	session := r.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
		AccessMode: neo4jDriver.AccessModeWrite,
	})
	defer session.Close(ctx)

	query := `
		CREATE (m:Member {
			id: $id,
			name: $name,
			dateOfBirth: datetime($dateOfBirth),
			gender: $gender,
			email: $email,
			phone: $phone,
			avatarUrl: $avatarUrl,
			createdAt: datetime($createdAt),
			updatedAt: datetime($updatedAt),
			isDeleted: $isDeleted
		})
		RETURN m
	`

	params := map[string]interface{}{
		"id":          member.ID,
		"name":        member.Name,
		"dateOfBirth": member.DateOfBirth.Format(time.RFC3339),
		"gender":      member.Gender,
		"email":       member.Email,
		"phone":       member.Phone,
		"avatarUrl":   member.AvatarURL,
		"createdAt":   member.CreatedAt.Format(time.RFC3339),
		"updatedAt":   member.UpdatedAt.Format(time.RFC3339),
		"isDeleted":   member.IsDeleted,
	}

	_, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}

	return nil
}

// GetByID retrieves a member by ID from Neo4j
func (r *MemberRepository) GetByID(ctx context.Context, id string) (*models.Member, error) {
	session := r.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
		AccessMode: neo4jDriver.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
		MATCH (m:Member {id: $id})
		WHERE m.isDeleted = false
		RETURN m.id, m.name, m.dateOfBirth, m.gender, m.email, m.phone, 
		       m.avatarUrl, m.createdAt, m.updatedAt, m.isDeleted
	`

	result, err := session.Run(ctx, query, map[string]interface{}{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to query member: %w", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		return r.recordToMember(record)
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	return nil, fmt.Errorf("member not found")
}

// GetAll retrieves all non-deleted members from Neo4j
func (r *MemberRepository) GetAll(ctx context.Context) ([]*models.Member, error) {
	session := r.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
		AccessMode: neo4jDriver.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
		MATCH (m:Member)
		WHERE m.isDeleted = false
		RETURN m.id, m.name, m.dateOfBirth, m.gender, m.email, m.phone, 
		       m.avatarUrl, m.createdAt, m.updatedAt, m.isDeleted
		ORDER BY m.name
	`

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query members: %w", err)
	}

	var members []*models.Member
	for result.Next(ctx) {
		record := result.Record()
		member, err := r.recordToMember(record)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	return members, nil
}

// Update updates an existing member in Neo4j
func (r *MemberRepository) Update(ctx context.Context, member *models.Member) error {
	session := r.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
		AccessMode: neo4jDriver.AccessModeWrite,
	})
	defer session.Close(ctx)

	query := `
		MATCH (m:Member {id: $id})
		SET m.name = $name,
		    m.dateOfBirth = datetime($dateOfBirth),
		    m.gender = $gender,
		    m.email = $email,
		    m.phone = $phone,
		    m.avatarUrl = $avatarUrl,
		    m.updatedAt = datetime($updatedAt)
		RETURN m
	`

	params := map[string]interface{}{
		"id":          member.ID,
		"name":        member.Name,
		"dateOfBirth": member.DateOfBirth.Format(time.RFC3339),
		"gender":      member.Gender,
		"email":       member.Email,
		"phone":       member.Phone,
		"avatarUrl":   member.AvatarURL,
		"updatedAt":   member.UpdatedAt.Format(time.RFC3339),
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}

	if !result.Next(ctx) {
		return fmt.Errorf("member not found")
	}

	return nil
}

// SoftDelete marks a member as deleted without removing from database
func (r *MemberRepository) SoftDelete(ctx context.Context, id string) error {
	session := r.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
		AccessMode: neo4jDriver.AccessModeWrite,
	})
	defer session.Close(ctx)

	query := `
		MATCH (m:Member {id: $id})
		SET m.isDeleted = true,
		    m.updatedAt = datetime($updatedAt)
		RETURN m
	`

	params := map[string]interface{}{
		"id":        id,
		"updatedAt": time.Now().Format(time.RFC3339),
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to soft delete member: %w", err)
	}

	if !result.Next(ctx) {
		return fmt.Errorf("member not found")
	}

	return nil
}

// EnsureIndexes creates indexes for member.id and member.name
func (r *MemberRepository) EnsureIndexes(ctx context.Context) error {
	session := r.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
		AccessMode: neo4jDriver.AccessModeWrite,
	})
	defer session.Close(ctx)

	// Create index on member.id
	idIndexQuery := `
		CREATE INDEX member_id_index IF NOT EXISTS
		FOR (m:Member) ON (m.id)
	`

	_, err := session.Run(ctx, idIndexQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to create id index: %w", err)
	}

	// Create index on member.name
	nameIndexQuery := `
		CREATE INDEX member_name_index IF NOT EXISTS
		FOR (m:Member) ON (m.name)
	`

	_, err = session.Run(ctx, nameIndexQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to create name index: %w", err)
	}

	return nil
}

// recordToMember converts a Neo4j record to a Member model
func (r *MemberRepository) recordToMember(record *neo4jDriver.Record) (*models.Member, error) {
	id, _ := record.Get("m.id")
	name, _ := record.Get("m.name")
	dateOfBirthRaw, _ := record.Get("m.dateOfBirth")
	gender, _ := record.Get("m.gender")
	email, _ := record.Get("m.email")
	phone, _ := record.Get("m.phone")
	avatarUrl, _ := record.Get("m.avatarUrl")
	createdAtRaw, _ := record.Get("m.createdAt")
	updatedAtRaw, _ := record.Get("m.updatedAt")
	isDeleted, _ := record.Get("m.isDeleted")

	// Parse datetime fields
	dateOfBirth, err := r.parseNeo4jTime(dateOfBirthRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dateOfBirth: %w", err)
	}

	createdAt, err := r.parseNeo4jTime(createdAtRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse createdAt: %w", err)
	}

	updatedAt, err := r.parseNeo4jTime(updatedAtRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updatedAt: %w", err)
	}

	member := &models.Member{
		ID:          id.(string),
		Name:        name.(string),
		DateOfBirth: dateOfBirth,
		Gender:      gender.(string),
		Email:       r.stringOrEmpty(email),
		Phone:       r.stringOrEmpty(phone),
		AvatarURL:   r.stringOrEmpty(avatarUrl),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		IsDeleted:   isDeleted.(bool),
	}

	return member, nil
}

// parseNeo4jTime converts Neo4j datetime to Go time.Time
func (r *MemberRepository) parseNeo4jTime(value interface{}) (time.Time, error) {
	switch v := value.(type) {
	case neo4jDriver.LocalDateTime:
		return v.Time(), nil
	case time.Time:
		return v, nil
	case string:
		return time.Parse(time.RFC3339, v)
	default:
		return time.Time{}, fmt.Errorf("unsupported time type: %T", value)
	}
}

// stringOrEmpty safely converts interface{} to string, returning empty string if nil
func (r *MemberRepository) stringOrEmpty(value interface{}) string {
	if value == nil {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
