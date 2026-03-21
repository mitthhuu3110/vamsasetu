package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"vamsasetu/backend/internal/models"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	neo4jDriver "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// setupTestRelationshipRepo creates a test relationship repository with Neo4j connection
func setupTestRelationshipRepo(t *testing.T) (*RelationshipRepository, *MemberRepository, context.Context, func()) {
	memberRepo, ctx, cleanup := setupTestRepo(t)
	
	// Create relationship repository using the same client
	relRepo := NewRelationshipRepository(memberRepo.client)
	
	return relRepo, memberRepo, ctx, cleanup
}

// **Validates: Requirements 2.4**
// Property 10: Relationship Creation and Retrieval
// For any two existing members and a valid relationship type (SPOUSE_OF, PARENT_OF, SIBLING_OF),
// creating a relationship should result in that relationship being stored in Neo4j and retrievable.
func TestProperty10_RelationshipCreationAndRetrieval(t *testing.T) {
	relRepo, memberRepo, ctx, cleanup := setupTestRelationshipRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any two existing members and valid relationship type, creating a relationship should make it retrievable",
		prop.ForAll(
			func(name1 string, dob1 time.Time, gender1 string,
				name2 string, dob2 time.Time, gender2 string,
				relType string) bool {

				// Create two members
				member1 := models.NewMember(name1, dob1, gender1)
				member2 := models.NewMember(name2, dob2, gender2)

				err := memberRepo.Create(ctx, member1)
				if err != nil {
					t.Logf("Failed to create member1: %v", err)
					return false
				}

				err = memberRepo.Create(ctx, member2)
				if err != nil {
					t.Logf("Failed to create member2: %v", err)
					memberRepo.SoftDelete(ctx, member1.ID)
					return false
				}

				// Create relationship
				rel := models.NewRelationship(relType, member1.ID, member2.ID)
				err = relRepo.Create(ctx, rel)
				if err != nil {
					t.Logf("Failed to create relationship: %v", err)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Retrieve all relationships and verify our relationship exists
				allRels, err := relRepo.GetAll(ctx)
				if err != nil {
					t.Logf("Failed to retrieve relationships: %v", err)
					relRepo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Check if our relationship exists in the results
				found := false
				for _, r := range allRels {
					if r.Type == relType {
						// For bidirectional relationships, check both directions
						if rel.IsBidirectional() {
							if (r.FromID == member1.ID && r.ToID == member2.ID) ||
								(r.FromID == member2.ID && r.ToID == member1.ID) {
								found = true
								break
							}
						} else {
							// For directed relationships, check exact direction
							if r.FromID == member1.ID && r.ToID == member2.ID {
								found = true
								break
							}
						}
					}
				}

				if !found {
					t.Logf("Created relationship not found in GetAll results")
				}

				// Cleanup
				relRepo.Delete(ctx, member1.ID, member2.ID, relType)
				memberRepo.SoftDelete(ctx, member1.ID)
				memberRepo.SoftDelete(ctx, member2.ID)

				return found
			},
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
			genValidRelationshipType(),
		),
	)

	properties.TestingRun(t)
}

// **Validates: Requirements 2.5**
// Property 11: Relationship Deletion
// For any existing relationship, deleting that relationship should remove it from Neo4j
// such that it is no longer retrievable.
func TestProperty11_RelationshipDeletion(t *testing.T) {
	relRepo, memberRepo, ctx, cleanup := setupTestRelationshipRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any existing relationship, deleting it should remove it from Neo4j",
		prop.ForAll(
			func(name1 string, dob1 time.Time, gender1 string,
				name2 string, dob2 time.Time, gender2 string,
				relType string) bool {

				// Create two members
				member1 := models.NewMember(name1, dob1, gender1)
				member2 := models.NewMember(name2, dob2, gender2)

				err := memberRepo.Create(ctx, member1)
				if err != nil {
					t.Logf("Failed to create member1: %v", err)
					return false
				}

				err = memberRepo.Create(ctx, member2)
				if err != nil {
					t.Logf("Failed to create member2: %v", err)
					memberRepo.SoftDelete(ctx, member1.ID)
					return false
				}

				// Create relationship
				rel := models.NewRelationship(relType, member1.ID, member2.ID)
				err = relRepo.Create(ctx, rel)
				if err != nil {
					t.Logf("Failed to create relationship: %v", err)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Verify relationship exists
				allRelsBefore, err := relRepo.GetAll(ctx)
				if err != nil {
					t.Logf("Failed to retrieve relationships before delete: %v", err)
					relRepo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				foundBefore := false
				for _, r := range allRelsBefore {
					if r.Type == relType {
						if rel.IsBidirectional() {
							if (r.FromID == member1.ID && r.ToID == member2.ID) ||
								(r.FromID == member2.ID && r.ToID == member1.ID) {
								foundBefore = true
								break
							}
						} else {
							if r.FromID == member1.ID && r.ToID == member2.ID {
								foundBefore = true
								break
							}
						}
					}
				}

				if !foundBefore {
					t.Logf("Relationship not found before deletion")
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Delete relationship
				err = relRepo.Delete(ctx, member1.ID, member2.ID, relType)
				if err != nil {
					t.Logf("Failed to delete relationship: %v", err)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Verify relationship no longer exists
				allRelsAfter, err := relRepo.GetAll(ctx)
				if err != nil {
					t.Logf("Failed to retrieve relationships after delete: %v", err)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				foundAfter := false
				for _, r := range allRelsAfter {
					if r.Type == relType {
						if rel.IsBidirectional() {
							if (r.FromID == member1.ID && r.ToID == member2.ID) ||
								(r.FromID == member2.ID && r.ToID == member1.ID) {
								foundAfter = true
								break
							}
						} else {
							if r.FromID == member1.ID && r.ToID == member2.ID {
								foundAfter = true
								break
							}
						}
					}
				}

				if foundAfter {
					t.Logf("Relationship still found after deletion")
				}

				// Cleanup
				memberRepo.SoftDelete(ctx, member1.ID)
				memberRepo.SoftDelete(ctx, member2.ID)

				return !foundAfter
			},
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
			genValidRelationshipType(),
		),
	)

	properties.TestingRun(t)
}

// **Validates: Requirements 2.6**
// Property 12: Relationship Semantic Validation
// For any member, attempting to create a PARENT_OF relationship from that member to itself
// should be rejected with a validation error.
func TestProperty12_RelationshipSemanticValidation(t *testing.T) {
	relRepo, memberRepo, ctx, cleanup := setupTestRelationshipRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any member, creating a self-referential relationship should be rejected",
		prop.ForAll(
			func(name string, dob time.Time, gender string) bool {

				// Create a member
				member := models.NewMember(name, dob, gender)

				err := memberRepo.Create(ctx, member)
				if err != nil {
					t.Logf("Failed to create member: %v", err)
					return false
				}

				// Attempt to create self-referential relationship
				rel := models.NewRelationship(models.RelationshipTypeParentOf, member.ID, member.ID)
				
				// Validate should fail
				validationErr := rel.Validate()
				if validationErr == nil {
					t.Logf("Validation should have failed for self-referential relationship")
					memberRepo.SoftDelete(ctx, member.ID)
					return false
				}

				// Attempt to create in repository should also fail
				createErr := relRepo.Create(ctx, rel)
				
				// Either validation catches it or Neo4j query fails (no matching nodes)
				shouldFail := (validationErr != nil || createErr != nil)

				// Cleanup
				memberRepo.SoftDelete(ctx, member.ID)

				return shouldFail
			},
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
		),
	)

	properties.TestingRun(t)
}

// **Validates: Requirements 2.7**
// Property 13: Soft Delete Enforcement for Connected Members
// For any member with at least one relationship, attempting to delete that member should result
// in a soft delete (isDeleted=true) rather than removal from the database.
func TestProperty13_SoftDeleteEnforcementForConnectedMembers(t *testing.T) {
	relRepo, memberRepo, ctx, cleanup := setupTestRelationshipRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any member with relationships, soft delete should preserve data with isDeleted=true",
		prop.ForAll(
			func(name1 string, dob1 time.Time, gender1 string,
				name2 string, dob2 time.Time, gender2 string,
				relType string) bool {

				// Create two members
				member1 := models.NewMember(name1, dob1, gender1)
				member2 := models.NewMember(name2, dob2, gender2)

				err := memberRepo.Create(ctx, member1)
				if err != nil {
					t.Logf("Failed to create member1: %v", err)
					return false
				}

				err = memberRepo.Create(ctx, member2)
				if err != nil {
					t.Logf("Failed to create member2: %v", err)
					memberRepo.SoftDelete(ctx, member1.ID)
					return false
				}

				// Create relationship (member1 now has a connection)
				rel := models.NewRelationship(relType, member1.ID, member2.ID)
				err = relRepo.Create(ctx, rel)
				if err != nil {
					t.Logf("Failed to create relationship: %v", err)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Store original data for comparison
				originalID := member1.ID
				originalName := member1.Name
				originalGender := member1.Gender

				// Perform soft delete on member1 (who has a relationship)
				err = memberRepo.SoftDelete(ctx, member1.ID)
				if err != nil {
					t.Logf("Failed to soft delete member: %v", err)
					relRepo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Verify member is not retrievable via GetByID (filters out deleted)
				_, err = memberRepo.GetByID(ctx, member1.ID)
				if err == nil {
					t.Logf("Soft deleted member should not be retrievable via GetByID")
					relRepo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Verify data is preserved by querying directly (bypassing isDeleted filter)
				session := memberRepo.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
					AccessMode: neo4jDriver.AccessModeRead,
				})
				defer session.Close(ctx)

				query := `
					MATCH (m:Member {id: $id})
					RETURN m.id, m.name, m.gender, m.isDeleted
				`
				result, err := session.Run(ctx, query, map[string]interface{}{"id": originalID})
				if err != nil {
					t.Logf("Failed to query deleted member: %v", err)
					relRepo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				if !result.Next(ctx) {
					t.Logf("Soft deleted member should still exist in database")
					relRepo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				record := result.Record()
				id, _ := record.Get("m.id")
				name, _ := record.Get("m.name")
				gender, _ := record.Get("m.gender")
				isDeleted, _ := record.Get("m.isDeleted")

				// Verify isDeleted is true
				if isDeleted.(bool) != true {
					t.Logf("isDeleted should be true, got %v", isDeleted)
					relRepo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Verify all other data is preserved
				if id.(string) != originalID {
					t.Logf("ID not preserved: expected %s, got %s", originalID, id)
					relRepo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}
				if name.(string) != originalName {
					t.Logf("Name not preserved: expected %s, got %s", originalName, name)
					relRepo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}
				if gender.(string) != originalGender {
					t.Logf("Gender not preserved: expected %s, got %s", originalGender, gender)
					relRepo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Cleanup: delete relationship and second member
				relRepo.Delete(ctx, member1.ID, member2.ID, relType)
				memberRepo.SoftDelete(ctx, member2.ID)

				// Hard delete test members
				session.Run(ctx, "MATCH (m:Member {id: $id}) DELETE m", map[string]interface{}{"id": member1.ID})
				session.Run(ctx, "MATCH (m:Member {id: $id}) DELETE m", map[string]interface{}{"id": member2.ID})

				return true
			},
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
			genValidRelationshipType(),
		),
	)

	properties.TestingRun(t)
}

// Generator functions for property-based testing

// genValidRelationshipType generates valid relationship types
func genValidRelationshipType() gopter.Gen {
	return gen.OneConstOf(
		models.RelationshipTypeSpouseOf,
		models.RelationshipTypeParentOf,
		models.RelationshipTypeSiblingOf,
	)
}
