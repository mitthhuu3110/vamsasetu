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

// **Validates: Requirements 2.1, 11.1**
// Property 7: Member Creation and Retrieval
// For any valid member data (name, date of birth, gender), creating a member
// should result in that member being stored in Neo4j and retrievable by ID.
func TestProperty7_MemberCreationAndRetrieval(t *testing.T) {
	repo, ctx, cleanup := setupTestRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any valid member data, creating a member should result in that member being retrievable by ID",
		prop.ForAll(
			func(name string, dob time.Time, gender string) bool {
				// Create member with generated data
				member := models.NewMember(name, dob, gender)

				// Create in repository
				err := repo.Create(ctx, member)
				if err != nil {
					t.Logf("Failed to create member: %v", err)
					return false
				}

				// Retrieve by ID
				retrieved, err := repo.GetByID(ctx, member.ID)
				if err != nil {
					t.Logf("Failed to retrieve member: %v", err)
					return false
				}

				// Verify all fields match
				if retrieved.ID != member.ID {
					t.Logf("ID mismatch: expected %s, got %s", member.ID, retrieved.ID)
					return false
				}
				if retrieved.Name != member.Name {
					t.Logf("Name mismatch: expected %s, got %s", member.Name, retrieved.Name)
					return false
				}
				if !retrieved.DateOfBirth.Equal(member.DateOfBirth) {
					t.Logf("DateOfBirth mismatch: expected %v, got %v", member.DateOfBirth, retrieved.DateOfBirth)
					return false
				}
				if retrieved.Gender != member.Gender {
					t.Logf("Gender mismatch: expected %s, got %s", member.Gender, retrieved.Gender)
					return false
				}
				if retrieved.IsDeleted != false {
					t.Logf("IsDeleted should be false, got %v", retrieved.IsDeleted)
					return false
				}

				// Cleanup: soft delete the test member
				repo.SoftDelete(ctx, member.ID)

				return true
			},
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
		),
	)

	properties.TestingRun(t)
}

// **Validates: Requirements 2.2**
// Property 8: Member Update Persistence
// For any existing member and any valid attribute changes, updating the member
// should persist those changes such that subsequent retrieval returns the updated values.
func TestProperty8_MemberUpdatePersistence(t *testing.T) {
	repo, ctx, cleanup := setupTestRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any existing member and valid attribute changes, updating should persist changes",
		prop.ForAll(
			func(originalName string, originalDob time.Time, originalGender string,
				updatedName string, updatedEmail string, updatedPhone string) bool {

				// Create initial member
				member := models.NewMember(originalName, originalDob, originalGender)
				err := repo.Create(ctx, member)
				if err != nil {
					t.Logf("Failed to create member: %v", err)
					return false
				}

				// Update member attributes
				member.Name = updatedName
				member.Email = updatedEmail
				member.Phone = updatedPhone
				member.Update()

				// Persist update
				err = repo.Update(ctx, member)
				if err != nil {
					t.Logf("Failed to update member: %v", err)
					return false
				}

				// Retrieve and verify updates
				retrieved, err := repo.GetByID(ctx, member.ID)
				if err != nil {
					t.Logf("Failed to retrieve updated member: %v", err)
					return false
				}

				// Verify updated fields
				if retrieved.Name != updatedName {
					t.Logf("Name not updated: expected %s, got %s", updatedName, retrieved.Name)
					return false
				}
				if retrieved.Email != updatedEmail {
					t.Logf("Email not updated: expected %s, got %s", updatedEmail, retrieved.Email)
					return false
				}
				if retrieved.Phone != updatedPhone {
					t.Logf("Phone not updated: expected %s, got %s", updatedPhone, retrieved.Phone)
					return false
				}

				// Verify unchanged fields
				if retrieved.ID != member.ID {
					t.Logf("ID should not change: expected %s, got %s", member.ID, retrieved.ID)
					return false
				}
				if retrieved.Gender != originalGender {
					t.Logf("Gender should not change: expected %s, got %s", originalGender, retrieved.Gender)
					return false
				}

				// Cleanup
				repo.SoftDelete(ctx, member.ID)

				return true
			},
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
			genValidName(),
			genValidEmail(),
			genValidPhone(),
		),
	)

	properties.TestingRun(t)
}

// **Validates: Requirements 2.3**
// Property 9: Soft Delete Preservation
// For any member, deleting that member should set the isDeleted flag to true
// while preserving all other member data.
func TestProperty9_SoftDeletePreservation(t *testing.T) {
	repo, ctx, cleanup := setupTestRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any member, soft delete should set isDeleted=true while preserving all data",
		prop.ForAll(
			func(name string, dob time.Time, gender string, email string, phone string) bool {
				// Create member with all fields populated
				member := models.NewMember(name, dob, gender)
				member.Email = email
				member.Phone = phone

				err := repo.Create(ctx, member)
				if err != nil {
					t.Logf("Failed to create member: %v", err)
					return false
				}

				// Store original values for comparison
				originalID := member.ID
				originalName := member.Name
				originalDob := member.DateOfBirth
				originalGender := member.Gender
				originalEmail := member.Email
				originalPhone := member.Phone

				// Perform soft delete
				err = repo.SoftDelete(ctx, member.ID)
				if err != nil {
					t.Logf("Failed to soft delete member: %v", err)
					return false
				}

				// Verify member is not retrievable via GetByID (filters out deleted)
				_, err = repo.GetByID(ctx, member.ID)
				if err == nil {
					t.Logf("Soft deleted member should not be retrievable via GetByID")
					return false
				}

				// Verify member is not in GetAll results (filters out deleted)
				allMembers, err := repo.GetAll(ctx)
				if err != nil {
					t.Logf("Failed to get all members: %v", err)
					return false
				}
				for _, m := range allMembers {
					if m.ID == originalID {
						t.Logf("Soft deleted member should not appear in GetAll results")
						return false
					}
				}

				// Verify data is preserved by querying directly (bypassing isDeleted filter)
				session := repo.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
					AccessMode: neo4jDriver.AccessModeRead,
				})
				defer session.Close(ctx)

				query := `
					MATCH (m:Member {id: $id})
					RETURN m.id, m.name, m.dateOfBirth, m.gender, m.email, m.phone, m.isDeleted
				`
				result, err := session.Run(ctx, query, map[string]interface{}{"id": originalID})
				if err != nil {
					t.Logf("Failed to query deleted member: %v", err)
					return false
				}

				if !result.Next(ctx) {
					t.Logf("Soft deleted member should still exist in database")
					return false
				}

				record := result.Record()
				id, _ := record.Get("m.id")
				name, _ := record.Get("m.name")
				gender, _ := record.Get("m.gender")
				email, _ := record.Get("m.email")
				phone, _ := record.Get("m.phone")
				isDeleted, _ := record.Get("m.isDeleted")

				// Verify isDeleted is true
				if isDeleted.(bool) != true {
					t.Logf("isDeleted should be true, got %v", isDeleted)
					return false
				}

				// Verify all other data is preserved
				if id.(string) != originalID {
					t.Logf("ID not preserved: expected %s, got %s", originalID, id)
					return false
				}
				if name.(string) != originalName {
					t.Logf("Name not preserved: expected %s, got %s", originalName, name)
					return false
				}
				if gender.(string) != originalGender {
					t.Logf("Gender not preserved: expected %s, got %s", originalGender, gender)
					return false
				}

				// Handle optional fields (may be nil)
				emailStr := ""
				if email != nil {
					emailStr = email.(string)
				}
				if emailStr != originalEmail {
					t.Logf("Email not preserved: expected %s, got %s", originalEmail, emailStr)
					return false
				}

				phoneStr := ""
				if phone != nil {
					phoneStr = phone.(string)
				}
				if phoneStr != originalPhone {
					t.Logf("Phone not preserved: expected %s, got %s", originalPhone, phoneStr)
					return false
				}

				// Final cleanup: hard delete the test member
				session.Run(ctx, "MATCH (m:Member {id: $id}) DELETE m", map[string]interface{}{"id": originalID})

				return true
			},
			genValidName(),
			genValidDateOfBirth(),
			genValidGender(),
			genValidEmail(),
			genValidPhone(),
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

// genValidEmail generates valid email addresses
func genValidEmail() gopter.Gen {
	return gen.AlphaString().
		SuchThat(func(s string) bool {
			return len(s) >= 3 && len(s) <= 50
		}).
		Map(func(s string) string {
			return s + "@test.com"
		})
}

// genValidPhone generates valid phone numbers
func genValidPhone() gopter.Gen {
	return gen.Int64Range(1000000000, 9999999999).
		Map(func(n int64) string {
			return "+91" + fmt.Sprintf("%d", n)
		})
}
