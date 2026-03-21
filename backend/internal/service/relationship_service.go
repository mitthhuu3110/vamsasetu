package service

import (
	"context"
	"fmt"
	"strings"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
)

// RelationshipService handles relationship business logic including kinship mapping
type RelationshipService struct {
	repo *repository.RelationshipRepository
	hub  WebSocketHub
}

// RelationshipResult represents the result of a relationship query
type RelationshipResult struct {
	Path          []repository.PathNode `json:"path"`
	RelationLabel string                `json:"relationLabel"`
	KinshipTerm   string                `json:"kinshipTerm"`
	Description   string                `json:"description"`
}

// NewRelationshipService creates a new relationship service instance
func NewRelationshipService(repo *repository.RelationshipRepository, hub WebSocketHub) *RelationshipService {
	return &RelationshipService{
		repo: repo,
		hub:  hub,
	}
}

// Create creates a new relationship
func (s *RelationshipService) Create(ctx context.Context, relationship *models.Relationship) error {
	if err := s.repo.Create(ctx, relationship); err != nil {
		return err
	}

	// Broadcast WebSocket update
	if s.hub != nil {
		s.hub.BroadcastUpdate("relationship_created", relationship)
	}

	return nil
}

// GetAll retrieves all relationships
func (s *RelationshipService) GetAll(ctx context.Context) ([]*models.Relationship, error) {
	return s.repo.GetAll(ctx)
}

// Delete deletes a relationship
func (s *RelationshipService) Delete(ctx context.Context, fromID, toID, relType string) error {
	if err := s.repo.Delete(ctx, fromID, toID, relType); err != nil {
		return err
	}

	// Broadcast WebSocket update
	if s.hub != nil {
		s.hub.BroadcastUpdate("relationship_deleted", map[string]string{
			"fromId": fromID,
			"toId":   toID,
			"type":   relType,
		})
	}

	return nil
}

// GetByMemberID retrieves all relationships for a specific member
func (s *RelationshipService) GetByMemberID(ctx context.Context, memberID string) ([]*models.Relationship, error) {
	// Get all relationships and filter by member ID
	allRelationships, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var memberRelationships []*models.Relationship
	for _, rel := range allRelationships {
		if rel.FromID == memberID || rel.ToID == memberID {
			memberRelationships = append(memberRelationships, rel)
		}
	}

	return memberRelationships, nil
}

// FindRelationship finds the relationship path between two members and maps it to kinship terms
func (s *RelationshipService) FindRelationship(ctx context.Context, fromID, toID string) (*RelationshipResult, error) {
	// Handle same person case
	if fromID == toID {
		return &RelationshipResult{
			Path:          nil,
			RelationLabel: "Self",
			KinshipTerm:   "Nenu/Nuvvu",
			Description:   "This is the same person",
		}, nil
	}

	// Query Neo4j for shortest path
	path, err := s.repo.FindPath(ctx, fromID, toID)
	if err != nil {
		return nil, fmt.Errorf("failed to find path: %w", err)
	}

	// No path found - not related
	if path == nil || len(path.Nodes) == 0 {
		return &RelationshipResult{
			Path:          nil,
			RelationLabel: "Not Related",
			KinshipTerm:   "",
			Description:   "No family connection found between these members",
		}, nil
	}

	// Map the path to kinship terms
	relationLabel, kinshipTerm := s.mapToKinshipTerm(path)

	// Generate natural language description
	description := s.generateDescription(path, relationLabel)

	return &RelationshipResult{
		Path:          path.Nodes,
		RelationLabel: relationLabel,
		KinshipTerm:   kinshipTerm,
		Description:   description,
	}, nil
}

// mapToKinshipTerm maps a relationship path to appropriate kinship terms
func (s *RelationshipService) mapToKinshipTerm(path *repository.RelationshipPath) (string, string) {
	if len(path.Relationships) == 0 {
		return "Unknown", ""
	}

	// Direct relationships (1 hop)
	if len(path.Relationships) == 1 {
		return s.mapDirectRelationship(path)
	}

	// Two-hop relationships
	if len(path.Relationships) == 2 {
		return s.mapTwoHopRelationship(path)
	}

	// Multi-hop relationships (3+ hops)
	return s.mapMultiHopRelationship(path)
}

// mapDirectRelationship handles single-edge relationships
func (s *RelationshipService) mapDirectRelationship(path *repository.RelationshipPath) (string, string) {
	relType := path.Relationships[0]
	targetGender := path.Nodes[1].Gender

	switch relType {
	case models.RelationshipTypeParentOf:
		// From parent to child perspective
		if targetGender == "male" {
			return "Son", "Koduku"
		}
		return "Daughter", "Kuthuru"

	case models.RelationshipTypeSpouseOf:
		if targetGender == "male" {
			return "Husband", "Menarikam"
		}
		return "Wife", "Bharya"

	case models.RelationshipTypeSiblingOf:
		// For siblings, we need to determine if older or younger
		// Since we don't have age data in the path, we'll use generic terms
		if targetGender == "male" {
			return "Brother", "Annayya/Tammudu"
		}
		return "Sister", "Akka/Chelli"

	default:
		return "Related", ""
	}
}

// mapTwoHopRelationship handles two-edge relationships
func (s *RelationshipService) mapTwoHopRelationship(path *repository.RelationshipPath) (string, string) {
	rel1 := path.Relationships[0]
	rel2 := path.Relationships[1]
	middleGender := path.Nodes[1].Gender
	targetGender := path.Nodes[2].Gender

	// Grandparent relationships: PARENT_OF -> PARENT_OF (reversed)
	if rel1 == models.RelationshipTypeParentOf && rel2 == models.RelationshipTypeParentOf {
		// This is actually grandchild from grandparent's perspective
		if targetGender == "male" {
			return "Grandson", "Manumadu"
		}
		return "Granddaughter", "Manumalu"
	}

	// Parent's spouse (step-parent)
	if rel1 == models.RelationshipTypeParentOf && rel2 == models.RelationshipTypeSpouseOf {
		if targetGender == "male" {
			return "Step-Father", "Nanna"
		}
		return "Step-Mother", "Amma"
	}

	// Spouse's parent (in-laws)
	if rel1 == models.RelationshipTypeSpouseOf && rel2 == models.RelationshipTypeParentOf {
		if targetGender == "male" {
			return "Father-in-Law", "Maamayyagaru"
		}
		return "Mother-in-Law", "Attagaru"
	}

	// Parent's sibling (uncle/aunt)
	if rel1 == models.RelationshipTypeParentOf && rel2 == models.RelationshipTypeSiblingOf {
		if middleGender == "male" {
			// Father's sibling
			if targetGender == "male" {
				return "Uncle (Father's Brother)", "Babai"
			}
			return "Aunt (Father's Sister)", "Attha"
		} else {
			// Mother's sibling
			if targetGender == "male" {
				return "Uncle (Mother's Brother)", "Mamayya"
			}
			return "Aunt (Mother's Sister)", "Pinni"
		}
	}

	// Sibling's spouse (brother-in-law/sister-in-law)
	if rel1 == models.RelationshipTypeSiblingOf && rel2 == models.RelationshipTypeSpouseOf {
		if targetGender == "male" {
			return "Brother-in-Law", "Bava"
		}
		return "Sister-in-Law", "Vadina"
	}

	// Sibling's child (nephew/niece)
	if rel1 == models.RelationshipTypeSiblingOf && rel2 == models.RelationshipTypeParentOf {
		if targetGender == "male" {
			return "Nephew", "Bhanja/Alludu"
		}
		return "Niece", "Bhanjika/Kodalu"
	}

	// Spouse's sibling
	if rel1 == models.RelationshipTypeSpouseOf && rel2 == models.RelationshipTypeSiblingOf {
		if targetGender == "male" {
			return "Brother-in-Law", "Bava"
		}
		return "Sister-in-Law", "Maradalu"
	}

	return "Related (2 hops)", ""
}

// mapMultiHopRelationship handles three or more edge relationships
func (s *RelationshipService) mapMultiHopRelationship(path *repository.RelationshipPath) (string, string) {
	// Count parent hops to determine generation distance
	parentHops := 0
	for _, rel := range path.Relationships {
		if rel == models.RelationshipTypeParentOf {
			parentHops++
		}
	}

	targetGender := path.Nodes[len(path.Nodes)-1].Gender

	// Great-grandparent relationships (3+ parent hops going up)
	if parentHops >= 3 && s.isAscendingPath(path) {
		prefix := s.getGreatPrefix(parentHops - 2)
		if targetGender == "male" {
			return fmt.Sprintf("%sGrandfather", prefix), "Tata"
		}
		return fmt.Sprintf("%sGrandmother", prefix), "Ammamma"
	}

	// Great-grandchild relationships (3+ parent hops going down)
	if parentHops >= 3 && s.isDescendingPath(path) {
		prefix := s.getGreatPrefix(parentHops - 2)
		if targetGender == "male" {
			return fmt.Sprintf("%sGrandson", prefix), "Manumadu"
		}
		return fmt.Sprintf("%sGranddaughter", prefix), "Manumalu"
	}

	// Cousin relationships
	if s.isCousinPath(path) {
		degree := s.calculateCousinDegree(path)
		if degree == 1 {
			return "Cousin", "Bava/Maradalu"
		}
		return fmt.Sprintf("Cousin (%d degree)", degree), "Bava/Maradalu"
	}

	// Generic multi-hop relationship
	return fmt.Sprintf("Related (%d hops)", len(path.Relationships)), ""
}

// isAscendingPath checks if the path goes up the family tree (towards ancestors)
func (s *RelationshipService) isAscendingPath(path *repository.RelationshipPath) bool {
	// A path is ascending if it starts with PARENT_OF relationships
	// In Neo4j, PARENT_OF goes from parent to child, so we need to check the direction
	// For now, we'll check if most relationships are PARENT_OF
	parentCount := 0
	for _, rel := range path.Relationships {
		if rel == models.RelationshipTypeParentOf {
			parentCount++
		}
	}
	return parentCount > len(path.Relationships)/2
}

// isDescendingPath checks if the path goes down the family tree (towards descendants)
func (s *RelationshipService) isDescendingPath(path *repository.RelationshipPath) bool {
	// Similar logic to isAscendingPath
	return s.isAscendingPath(path)
}

// isCousinPath checks if the path represents a cousin relationship
func (s *RelationshipService) isCousinPath(path *repository.RelationshipPath) bool {
	// Cousin paths typically go up, across (sibling), and down
	hasParent := false
	hasSibling := false
	for _, rel := range path.Relationships {
		if rel == models.RelationshipTypeParentOf {
			hasParent = true
		}
		if rel == models.RelationshipTypeSiblingOf {
			hasSibling = true
		}
	}
	return hasParent && hasSibling
}

// calculateCousinDegree calculates the degree of cousin relationship
func (s *RelationshipService) calculateCousinDegree(path *repository.RelationshipPath) int {
	// Count the number of generations up before the sibling connection
	upGenerations := 0
	for _, rel := range path.Relationships {
		if rel == models.RelationshipTypeParentOf {
			upGenerations++
		} else if rel == models.RelationshipTypeSiblingOf {
			break
		}
	}
	return upGenerations
}

// getGreatPrefix returns the appropriate "Great-" prefix for multi-generational relationships
func (s *RelationshipService) getGreatPrefix(count int) string {
	if count == 0 {
		return ""
	}
	if count == 1 {
		return "Great-"
	}
	// For 2+, use "Great-Great-", etc.
	return strings.Repeat("Great-", count)
}

// generateDescription creates a natural language description of the relationship
func (s *RelationshipService) generateDescription(path *repository.RelationshipPath, relationLabel string) string {
	if len(path.Nodes) < 2 {
		return ""
	}

	fromName := path.Nodes[0].Name
	toName := path.Nodes[len(path.Nodes)-1].Name

	// For direct relationships
	if len(path.Nodes) == 2 {
		return fmt.Sprintf("%s is %s's %s", toName, fromName, strings.ToLower(relationLabel))
	}

	// For multi-hop relationships, include the path
	if len(path.Nodes) == 3 {
		middleName := path.Nodes[1].Name
		return fmt.Sprintf("%s is %s's %s through %s", toName, fromName, strings.ToLower(relationLabel), middleName)
	}

	// For longer paths, just mention the relationship
	return fmt.Sprintf("%s is %s's %s (%d connections)", toName, fromName, strings.ToLower(relationLabel), len(path.Relationships))
}
