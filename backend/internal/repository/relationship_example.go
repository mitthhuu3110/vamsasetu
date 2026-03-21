// +build ignore

package repository

import (
	"context"
	"fmt"
	"log"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/pkg/neo4j"
)

// ExampleRelationshipRepository demonstrates how to use the RelationshipRepository
func ExampleRelationshipRepository(client *neo4j.Client) {
	ctx := context.Background()
	repo := NewRelationshipRepository(client)

	// Example 1: Create a parent-child relationship
	fmt.Println("=== Example 1: Create Parent-Child Relationship ===")
	parentChildRel := models.NewRelationship(
		models.RelationshipTypeParentOf,
		"parent-uuid-123",
		"child-uuid-456",
	)
	if err := repo.Create(ctx, parentChildRel); err != nil {
		log.Printf("Error creating parent-child relationship: %v", err)
	} else {
		fmt.Println("✓ Parent-child relationship created successfully")
	}

	// Example 2: Create a spouse relationship (bidirectional)
	fmt.Println("\n=== Example 2: Create Spouse Relationship ===")
	spouseRel := models.NewRelationship(
		models.RelationshipTypeSpouseOf,
		"person1-uuid-789",
		"person2-uuid-012",
	)
	if err := repo.Create(ctx, spouseRel); err != nil {
		log.Printf("Error creating spouse relationship: %v", err)
	} else {
		fmt.Println("✓ Spouse relationship created (bidirectional edges)")
	}

	// Example 3: Create a sibling relationship (bidirectional)
	fmt.Println("\n=== Example 3: Create Sibling Relationship ===")
	siblingRel := models.NewRelationship(
		models.RelationshipTypeSiblingOf,
		"sibling1-uuid-345",
		"sibling2-uuid-678",
	)
	if err := repo.Create(ctx, siblingRel); err != nil {
		log.Printf("Error creating sibling relationship: %v", err)
	} else {
		fmt.Println("✓ Sibling relationship created (bidirectional edges)")
	}

	// Example 4: Get all relationships
	fmt.Println("\n=== Example 4: Get All Relationships ===")
	relationships, err := repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error getting relationships: %v", err)
	} else {
		fmt.Printf("✓ Found %d relationships:\n", len(relationships))
		for i, rel := range relationships {
			fmt.Printf("  %d. %s: %s -> %s (created: %s)\n",
				i+1, rel.Type, rel.FromID, rel.ToID, rel.CreatedAt.Format("2006-01-02"))
		}
	}

	// Example 5: Find path between two members
	fmt.Println("\n=== Example 5: Find Shortest Path ===")
	path, err := repo.FindPath(ctx, "member1-uuid", "member2-uuid")
	if err != nil {
		log.Printf("Error finding path: %v", err)
	} else if path == nil {
		fmt.Println("✗ No path found - members are not related")
	} else {
		fmt.Printf("✓ Path found with %d hops:\n", path.Length)
		fmt.Println("  Nodes:")
		for i, node := range path.Nodes {
			fmt.Printf("    %d. %s (ID: %s, Gender: %s)\n",
				i+1, node.Name, node.ID, node.Gender)
		}
		fmt.Println("  Relationships:")
		for i, relType := range path.Relationships {
			fmt.Printf("    %d. %s\n", i+1, relType)
		}
	}

	// Example 6: Find path that doesn't exist
	fmt.Println("\n=== Example 6: Find Non-Existent Path ===")
	noPath, err := repo.FindPath(ctx, "unrelated1-uuid", "unrelated2-uuid")
	if err != nil {
		log.Printf("Error finding path: %v", err)
	} else if noPath == nil {
		fmt.Println("✓ Correctly returned nil for unrelated members")
	}

	// Example 7: Delete a relationship
	fmt.Println("\n=== Example 7: Delete Relationship ===")
	if err := repo.Delete(ctx, "parent-uuid-123", "child-uuid-456", models.RelationshipTypeParentOf); err != nil {
		log.Printf("Error deleting relationship: %v", err)
	} else {
		fmt.Println("✓ Relationship deleted successfully")
	}

	// Example 8: Delete a bidirectional relationship
	fmt.Println("\n=== Example 8: Delete Bidirectional Relationship ===")
	if err := repo.Delete(ctx, "person1-uuid-789", "person2-uuid-012", models.RelationshipTypeSpouseOf); err != nil {
		log.Printf("Error deleting spouse relationship: %v", err)
	} else {
		fmt.Println("✓ Spouse relationship deleted (both directions removed)")
	}

	// Example 9: Complex family path
	fmt.Println("\n=== Example 9: Complex Family Path (Grandparent to Grandchild) ===")
	// Assuming we have: Grandparent -> Parent -> Child
	complexPath, err := repo.FindPath(ctx, "grandparent-uuid", "grandchild-uuid")
	if err != nil {
		log.Printf("Error finding complex path: %v", err)
	} else if complexPath != nil {
		fmt.Printf("✓ Found path with %d hops (2 generations):\n", complexPath.Length)
		for i := 0; i < len(complexPath.Nodes)-1; i++ {
			fmt.Printf("  %s -[%s]-> %s\n",
				complexPath.Nodes[i].Name,
				complexPath.Relationships[i],
				complexPath.Nodes[i+1].Name)
		}
	}

	// Example 10: Error handling - try to delete non-existent relationship
	fmt.Println("\n=== Example 10: Error Handling ===")
	if err := repo.Delete(ctx, "fake-uuid-1", "fake-uuid-2", models.RelationshipTypeParentOf); err != nil {
		fmt.Printf("✓ Correctly returned error: %v\n", err)
	}
}

// ExamplePathAnalysis demonstrates how to analyze a relationship path
func ExamplePathAnalysis(path *RelationshipPath) {
	if path == nil {
		fmt.Println("No relationship exists between these members")
		return
	}

	fmt.Printf("Relationship Analysis:\n")
	fmt.Printf("  Path Length: %d hops\n", path.Length)
	fmt.Printf("  Number of Nodes: %d\n", len(path.Nodes))
	fmt.Printf("  Number of Relationships: %d\n", len(path.Relationships))

	// Analyze relationship types
	fmt.Println("\n  Relationship Breakdown:")
	for i, relType := range path.Relationships {
		fromNode := path.Nodes[i]
		toNode := path.Nodes[i+1]
		fmt.Printf("    %s (%s) -[%s]-> %s (%s)\n",
			fromNode.Name, fromNode.Gender,
			relType,
			toNode.Name, toNode.Gender)
	}

	// Determine if path is direct or indirect
	if path.Length == 1 {
		fmt.Println("\n  Type: Direct relationship")
	} else {
		fmt.Printf("\n  Type: Indirect relationship (%d degrees of separation)\n", path.Length)
	}

	// Count relationship types
	relTypeCounts := make(map[string]int)
	for _, relType := range path.Relationships {
		relTypeCounts[relType]++
	}
	fmt.Println("\n  Relationship Type Counts:")
	for relType, count := range relTypeCounts {
		fmt.Printf("    %s: %d\n", relType, count)
	}
}

// ExampleBatchRelationshipCreation demonstrates creating multiple relationships efficiently
func ExampleBatchRelationshipCreation(client *neo4j.Client) {
	ctx := context.Background()
	repo := NewRelationshipRepository(client)

	// Define a family structure
	relationships := []*models.Relationship{
		// Parents
		models.NewRelationship(models.RelationshipTypeSpouseOf, "father-uuid", "mother-uuid"),
		// Children
		models.NewRelationship(models.RelationshipTypeParentOf, "father-uuid", "child1-uuid"),
		models.NewRelationship(models.RelationshipTypeParentOf, "mother-uuid", "child1-uuid"),
		models.NewRelationship(models.RelationshipTypeParentOf, "father-uuid", "child2-uuid"),
		models.NewRelationship(models.RelationshipTypeParentOf, "mother-uuid", "child2-uuid"),
		// Siblings
		models.NewRelationship(models.RelationshipTypeSiblingOf, "child1-uuid", "child2-uuid"),
	}

	fmt.Println("=== Batch Relationship Creation ===")
	successCount := 0
	for i, rel := range relationships {
		if err := repo.Create(ctx, rel); err != nil {
			log.Printf("Failed to create relationship %d: %v", i+1, err)
		} else {
			successCount++
		}
	}
	fmt.Printf("✓ Successfully created %d out of %d relationships\n", successCount, len(relationships))
}
