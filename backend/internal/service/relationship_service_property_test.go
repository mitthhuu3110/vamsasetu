package service

import (
	"context"
	"testing"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/pkg/neo4j"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	neo4jDriver "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// setupTestService creates a test relationship service with Neo4j connection
func setupTestService(t *testing.T) (*RelationshipService, *repository.MemberRepository, context.Context, func()) {
	// Create Neo4j client
	client, err := neo4j.NewClient(
		"bolt://localhost:7687",
		"neo4j",
		"vamsasetu123",
	)
	if err != nil {
		t.Fatalf("Failed to create Neo4j client: %v", err)
	}

	// Create repositories
	memberRepo := repository.NewMemberRepository(client)
	relRepo := repository.NewRelationshipRepository(client)

	// Create service
	service := NewRelationshipService(relRepo, nil)

	ctx := context.Background()

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		session := client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
			AccessMode: neo4jDriver.AccessModeWrite,
		})
		defer session.Close(ctx)

		// Delete all test members (those with names starting with "Test")
		session.Run(ctx, "MATCH (m:Member) WHERE m.name STARTS WITH 'Test' DETACH DELETE m", nil)

		client.Close()
	}

	return service, memberRepo, ctx, cleanup
}

// **Validates: Requirements 4.1, 4.4**
// Property 19: Relationship Path Finding
// For any two members in the same family tree, requesting the relationship between them
// should return either a valid path with nodes and edges, or a "not related" result if no path exists.
func TestProperty19_RelationshipPathFinding(t *testing.T) {
	service, memberRepo, ctx, cleanup := setupTestService(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 50,
		MaxSize:            50,
	})

	properties.Property("For any two members, FindRelationship should return either a valid path or 'not related'",
		prop.ForAll(
			func(name1 string, dob1 time.Time, gender1 string,
				name2 string, dob2 time.Time, gender2 string,
				createConnection bool) bool {

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

				// Optionally create a connection between them
				if createConnection {
					rel := models.NewRelationship(models.RelationshipTypeParentOf, member1.ID, member2.ID)
					err = service.repo.Create(ctx, rel)
					if err != nil {
						t.Logf("Failed to create relationship: %v", err)
						memberRepo.SoftDelete(ctx, member1.ID)
						memberRepo.SoftDelete(ctx, member2.ID)
						return false
					}
				}

				// Find relationship
				result, err := service.FindRelationship(ctx, member1.ID, member2.ID)
				if err != nil {
					t.Logf("FindRelationship returned error: %v", err)
					if createConnection {
						service.repo.Delete(ctx, member1.ID, member2.ID, models.RelationshipTypeParentOf)
					}
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Verify result is not nil
				if result == nil {
					t.Logf("FindRelationship returned nil result")
					if createConnection {
						service.repo.Delete(ctx, member1.ID, member2.ID, models.RelationshipTypeParentOf)
					}
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// If connection was created, verify we got a valid path
				if createConnection {
					if result.RelationLabel == "Not Related" {
						t.Logf("Expected a valid path but got 'Not Related'")
						service.repo.Delete(ctx, member1.ID, member2.ID, models.RelationshipTypeParentOf)
						memberRepo.SoftDelete(ctx, member1.ID)
						memberRepo.SoftDelete(ctx, member2.ID)
						return false
					}

					// Verify path contains both members
					if len(result.Path) < 2 {
						t.Logf("Path should contain at least 2 nodes, got %d", len(result.Path))
						service.repo.Delete(ctx, member1.ID, member2.ID, models.RelationshipTypeParentOf)
						memberRepo.SoftDelete(ctx, member1.ID)
						memberRepo.SoftDelete(ctx, member2.ID)
						return false
					}

					// Verify path starts with member1 and ends with member2
					if result.Path[0].ID != member1.ID {
						t.Logf("Path should start with member1 ID %s, got %s", member1.ID, result.Path[0].ID)
						service.repo.Delete(ctx, member1.ID, member2.ID, models.RelationshipTypeParentOf)
						memberRepo.SoftDelete(ctx, member1.ID)
						memberRepo.SoftDelete(ctx, member2.ID)
						return false
					}

					if result.Path[len(result.Path)-1].ID != member2.ID {
						t.Logf("Path should end with member2 ID %s, got %s", member2.ID, result.Path[len(result.Path)-1].ID)
						service.repo.Delete(ctx, member1.ID, member2.ID, models.RelationshipTypeParentOf)
						memberRepo.SoftDelete(ctx, member1.ID)
						memberRepo.SoftDelete(ctx, member2.ID)
						return false
					}
				} else {
					// If no connection was created, verify we got "Not Related"
					if result.RelationLabel != "Not Related" {
						t.Logf("Expected 'Not Related' but got '%s'", result.RelationLabel)
						memberRepo.SoftDelete(ctx, member1.ID)
						memberRepo.SoftDelete(ctx, member2.ID)
						return false
					}

					// Verify path is nil or empty for unrelated members
					if result.Path != nil && len(result.Path) > 0 {
						t.Logf("Path should be nil or empty for unrelated members, got %d nodes", len(result.Path))
						memberRepo.SoftDelete(ctx, member1.ID)
						memberRepo.SoftDelete(ctx, member2.ID)
						return false
					}
				}

				// Cleanup
				if createConnection {
					service.repo.Delete(ctx, member1.ID, member2.ID, models.RelationshipTypeParentOf)
				}
				memberRepo.SoftDelete(ctx, member1.ID)
				memberRepo.SoftDelete(ctx, member2.ID)

				return true
			},
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
			gen.Bool(), // createConnection
		),
	)

	properties.TestingRun(t)
}

// **Validates: Requirements 4.2, 4.3**
// Property 20: Relationship Result Completeness
// For any relationship query that finds a path, the result should contain all required fields:
// path nodes, relation label, kinship term, and natural language description.
func TestProperty20_RelationshipResultCompleteness(t *testing.T) {
	service, memberRepo, ctx, cleanup := setupTestService(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 50,
		MaxSize:            50,
	})

	properties.Property("For any relationship query that finds a path, result should contain all required fields",
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

				// Create a relationship between them
				rel := models.NewRelationship(relType, member1.ID, member2.ID)
				err = service.repo.Create(ctx, rel)
				if err != nil {
					t.Logf("Failed to create relationship: %v", err)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Find relationship
				result, err := service.FindRelationship(ctx, member1.ID, member2.ID)
				if err != nil {
					t.Logf("FindRelationship returned error: %v", err)
					service.repo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Verify result is not nil
				if result == nil {
					t.Logf("FindRelationship returned nil result")
					service.repo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Since we created a connection, this should not be "Not Related"
				if result.RelationLabel == "Not Related" {
					t.Logf("Expected a valid relationship but got 'Not Related'")
					service.repo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Verify all required fields are present and non-empty

				// 1. Path nodes should be present
				if result.Path == nil {
					t.Logf("Path should not be nil for connected members")
					service.repo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				if len(result.Path) < 2 {
					t.Logf("Path should contain at least 2 nodes, got %d", len(result.Path))
					service.repo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Verify each path node has required fields
				for i, node := range result.Path {
					if node.ID == "" {
						t.Logf("Path node %d has empty ID", i)
						service.repo.Delete(ctx, member1.ID, member2.ID, relType)
						memberRepo.SoftDelete(ctx, member1.ID)
						memberRepo.SoftDelete(ctx, member2.ID)
						return false
					}
					if node.Name == "" {
						t.Logf("Path node %d has empty Name", i)
						service.repo.Delete(ctx, member1.ID, member2.ID, relType)
						memberRepo.SoftDelete(ctx, member1.ID)
						memberRepo.SoftDelete(ctx, member2.ID)
						return false
					}
					if node.Gender == "" {
						t.Logf("Path node %d has empty Gender", i)
						service.repo.Delete(ctx, member1.ID, member2.ID, relType)
						memberRepo.SoftDelete(ctx, member1.ID)
						memberRepo.SoftDelete(ctx, member2.ID)
						return false
					}
				}

				// 2. Relation label should be present and non-empty
				if result.RelationLabel == "" {
					t.Logf("RelationLabel should not be empty")
					service.repo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// 3. Kinship term should be present (can be empty for some complex relationships)
				// We just verify the field exists, not that it's non-empty
				// since some relationships may not have Telugu terms

				// 4. Description should be present and non-empty
				if result.Description == "" {
					t.Logf("Description should not be empty")
					service.repo.Delete(ctx, member1.ID, member2.ID, relType)
					memberRepo.SoftDelete(ctx, member1.ID)
					memberRepo.SoftDelete(ctx, member2.ID)
					return false
				}

				// Cleanup
				service.repo.Delete(ctx, member1.ID, member2.ID, relType)
				memberRepo.SoftDelete(ctx, member1.ID)
				memberRepo.SoftDelete(ctx, member2.ID)

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

// genValidName generates valid member names (2-100 characters, alphanumeric with spaces)
func genValidName() gopter.Gen {
	return gen.AlphaString().
		SuchThat(func(s string) bool {
			return len(s) >= 2 && len(s) <= 100
		}).
		Map(func(s string) string {
			// Ensure it starts with "Test" for easy cleanup
			return "Test " + s
		})
}

// genValidDateOfBirth generates valid dates of birth (between 1900 and today)
func genValidDateOfBirth() gopter.Gen {
	minDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Now()

	return gen.Int64Range(minDate.Unix(), maxDate.Unix()).
		Map(func(timestamp int64) time.Time {
			return time.Unix(timestamp, 0).UTC()
		})
}

// genValidGender generates valid gender values
func genValidGender() gopter.Gen {
	return gen.OneConstOf("male", "female", "other")
}

// genValidRelationshipType generates valid relationship types
func genValidRelationshipType() gopter.Gen {
	return gen.OneConstOf(
		models.RelationshipTypeSpouseOf,
		models.RelationshipTypeParentOf,
		models.RelationshipTypeSiblingOf,
	)
}
