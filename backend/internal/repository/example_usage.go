// +build ignore

package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"vamsasetu/backend/internal/config"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/pkg/neo4j"
)

// ExampleMemberRepositoryUsage demonstrates how to use the MemberRepository
func ExampleMemberRepositoryUsage() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Neo4j client
	client, err := neo4j.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Neo4j client: %v", err)
	}
	defer client.Close(context.Background())

	// Create repository
	repo := NewMemberRepository(client)
	ctx := context.Background()

	// Ensure indexes are created
	if err := repo.EnsureIndexes(ctx); err != nil {
		log.Fatalf("Failed to ensure indexes: %v", err)
	}
	fmt.Println("✓ Indexes created successfully")

	// Example 1: Create a new member
	member := models.NewMember(
		"Rajesh Kumar",
		time.Date(1980, 5, 15, 0, 0, 0, 0, time.UTC),
		"male",
	)
	member.Email = "rajesh@example.com"
	member.Phone = "+919876543210"
	member.AvatarURL = "https://example.com/avatars/rajesh.jpg"

	if err := repo.Create(ctx, member); err != nil {
		log.Fatalf("Failed to create member: %v", err)
	}
	fmt.Printf("✓ Created member: %s (ID: %s)\n", member.Name, member.ID)

	// Example 2: Retrieve member by ID
	retrieved, err := repo.GetByID(ctx, member.ID)
	if err != nil {
		log.Fatalf("Failed to get member: %v", err)
	}
	fmt.Printf("✓ Retrieved member: %s, Gender: %s, Email: %s\n",
		retrieved.Name, retrieved.Gender, retrieved.Email)

	// Example 3: Update member
	retrieved.Email = "rajesh.kumar@example.com"
	retrieved.Phone = "+919876543211"
	retrieved.Update() // Updates the UpdatedAt timestamp

	if err := repo.Update(ctx, retrieved); err != nil {
		log.Fatalf("Failed to update member: %v", err)
	}
	fmt.Printf("✓ Updated member email to: %s\n", retrieved.Email)

	// Example 4: Get all members
	allMembers, err := repo.GetAll(ctx)
	if err != nil {
		log.Fatalf("Failed to get all members: %v", err)
	}
	fmt.Printf("✓ Retrieved %d total members\n", len(allMembers))

	// Example 5: Soft delete member
	if err := repo.SoftDelete(ctx, member.ID); err != nil {
		log.Fatalf("Failed to soft delete member: %v", err)
	}
	fmt.Printf("✓ Soft deleted member: %s\n", member.Name)

	// Example 6: Verify soft delete (should not be found)
	_, err = repo.GetByID(ctx, member.ID)
	if err != nil {
		fmt.Printf("✓ Confirmed: Soft deleted member is not retrievable\n")
	} else {
		fmt.Println("✗ Warning: Soft deleted member is still retrievable")
	}

	// Example 7: Create multiple members
	members := []*models.Member{
		models.NewMember("Lakshmi Devi", time.Date(1982, 8, 20, 0, 0, 0, 0, time.UTC), "female"),
		models.NewMember("Arjun Kumar", time.Date(2005, 3, 10, 0, 0, 0, 0, time.UTC), "male"),
		models.NewMember("Priya Sharma", time.Date(2008, 11, 5, 0, 0, 0, 0, time.UTC), "female"),
	}

	for _, m := range members {
		if err := repo.Create(ctx, m); err != nil {
			log.Printf("Failed to create member %s: %v", m.Name, err)
			continue
		}
		fmt.Printf("✓ Created member: %s\n", m.Name)
	}

	fmt.Println("\n✓ All examples completed successfully!")
}

// ExampleErrorHandling demonstrates error handling patterns
func ExampleErrorHandling() {
	cfg, _ := config.Load()
	client, _ := neo4j.NewClient(cfg)
	defer client.Close(context.Background())

	repo := NewMemberRepository(client)
	ctx := context.Background()

	// Example 1: Handle non-existent member
	_, err := repo.GetByID(ctx, "non-existent-id")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Example 2: Handle update of non-existent member
	nonExistent := models.NewMember("Ghost", time.Now(), "male")
	err = repo.Update(ctx, nonExistent)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Example 3: Handle soft delete of non-existent member
	err = repo.SoftDelete(ctx, "non-existent-id")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Example 4: Validate member before creating
	invalidMember := &models.Member{
		Name:   "", // Invalid: empty name
		Gender: "invalid",
	}
	if err := invalidMember.Validate(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	}
}

// ExampleBatchOperations demonstrates batch operations
func ExampleBatchOperations() {
	cfg, _ := config.Load()
	client, _ := neo4j.NewClient(cfg)
	defer client.Close(context.Background())

	repo := NewMemberRepository(client)
	ctx := context.Background()

	// Create a family tree
	family := []struct {
		name   string
		dob    time.Time
		gender string
		email  string
	}{
		{"Grandfather", time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC), "male", "grandfather@example.com"},
		{"Grandmother", time.Date(1952, 2, 2, 0, 0, 0, 0, time.UTC), "female", "grandmother@example.com"},
		{"Father", time.Date(1975, 3, 3, 0, 0, 0, 0, time.UTC), "male", "father@example.com"},
		{"Mother", time.Date(1977, 4, 4, 0, 0, 0, 0, time.UTC), "female", "mother@example.com"},
		{"Son", time.Date(2000, 5, 5, 0, 0, 0, 0, time.UTC), "male", "son@example.com"},
		{"Daughter", time.Date(2002, 6, 6, 0, 0, 0, 0, time.UTC), "female", "daughter@example.com"},
	}

	fmt.Println("Creating family tree...")
	for _, f := range family {
		member := models.NewMember(f.name, f.dob, f.gender)
		member.Email = f.email

		if err := repo.Create(ctx, member); err != nil {
			log.Printf("Failed to create %s: %v", f.name, err)
			continue
		}
		fmt.Printf("✓ Created: %s\n", f.name)
	}

	// Retrieve all and display
	allMembers, _ := repo.GetAll(ctx)
	fmt.Printf("\nTotal family members: %d\n", len(allMembers))
	for _, m := range allMembers {
		age := time.Now().Year() - m.DateOfBirth.Year()
		fmt.Printf("  - %s (%s, age %d)\n", m.Name, m.Gender, age)
	}
}
