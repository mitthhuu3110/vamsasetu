// +build ignore

package service

import (
	"context"
	"fmt"
	"log"

	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/pkg/neo4j"
)

// ExampleRelationshipService demonstrates how to use the RelationshipService
func ExampleRelationshipService() {
	// Initialize Neo4j client
	neo4jClient, err := neo4j.NewClient("bolt://localhost:7687", "neo4j", "vamsasetu123")
	if err != nil {
		log.Fatalf("Failed to create Neo4j client: %v", err)
	}
	defer neo4jClient.Close()

	// Initialize repository and service
	relationshipRepo := repository.NewRelationshipRepository(neo4jClient)
	relationshipService := NewRelationshipService(relationshipRepo)

	ctx := context.Background()

	// Example 1: Find direct parent-child relationship
	fmt.Println("=== Example 1: Parent-Child Relationship ===")
	result1, err := relationshipService.FindRelationship(ctx, "parent-id", "child-id")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		printRelationshipResult(result1)
	}

	// Example 2: Find grandparent relationship
	fmt.Println("\n=== Example 2: Grandparent Relationship ===")
	result2, err := relationshipService.FindRelationship(ctx, "grandparent-id", "grandchild-id")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		printRelationshipResult(result2)
	}

	// Example 3: Find uncle/aunt relationship
	fmt.Println("\n=== Example 3: Uncle/Aunt Relationship ===")
	result3, err := relationshipService.FindRelationship(ctx, "child-id", "uncle-id")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		printRelationshipResult(result3)
	}

	// Example 4: Find cousin relationship
	fmt.Println("\n=== Example 4: Cousin Relationship ===")
	result4, err := relationshipService.FindRelationship(ctx, "person-id", "cousin-id")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		printRelationshipResult(result4)
	}

	// Example 5: Same person
	fmt.Println("\n=== Example 5: Same Person ===")
	result5, err := relationshipService.FindRelationship(ctx, "person-id", "person-id")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		printRelationshipResult(result5)
	}

	// Example 6: No relationship
	fmt.Println("\n=== Example 6: No Relationship ===")
	result6, err := relationshipService.FindRelationship(ctx, "person-a-id", "unrelated-person-id")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		printRelationshipResult(result6)
	}
}

// printRelationshipResult prints a formatted relationship result
func printRelationshipResult(result *RelationshipResult) {
	fmt.Printf("Relationship: %s\n", result.RelationLabel)
	fmt.Printf("Telugu Term: %s\n", result.KinshipTerm)
	fmt.Printf("Description: %s\n", result.Description)
	
	if len(result.Path) > 0 {
		fmt.Printf("Path: ")
		for i, node := range result.Path {
			if i > 0 {
				fmt.Printf(" → ")
			}
			fmt.Printf("%s", node.Name)
		}
		fmt.Println()
	}
}

// ExampleBuildFamilyTreeWithRelationships demonstrates building a family tree and querying relationships
func ExampleBuildFamilyTreeWithRelationships() {
	// This example shows how to build a complete family tree and query relationships
	
	// Initialize clients
	neo4jClient, err := neo4j.NewClient("bolt://localhost:7687", "neo4j", "vamsasetu123")
	if err != nil {
		log.Fatalf("Failed to create Neo4j client: %v", err)
	}
	defer neo4jClient.Close()

	ctx := context.Background()

	// Initialize repositories
	memberRepo := repository.NewMemberRepository(neo4jClient)
	relationshipRepo := repository.NewRelationshipRepository(neo4jClient)
	relationshipService := NewRelationshipService(relationshipRepo)

	// Create family members
	fmt.Println("=== Creating Family Members ===")
	
	// Grandparents
	grandfather := createExampleMember("Venkatesh", "male")
	grandmother := createExampleMember("Lakshmi", "female")
	
	// Parents
	father := createExampleMember("Rajesh", "male")
	mother := createExampleMember("Priya", "female")
	uncle := createExampleMember("Kumar", "male")
	
	// Children
	son := createExampleMember("Arjun", "male")
	daughter := createExampleMember("Ananya", "female")
	cousin := createExampleMember("Rohan", "male")

	// Save members to Neo4j
	members := []*repository.PathNode{
		{ID: "gf-id", Name: grandfather, Gender: "male"},
		{ID: "gm-id", Name: grandmother, Gender: "female"},
		{ID: "f-id", Name: father, Gender: "male"},
		{ID: "m-id", Name: mother, Gender: "female"},
		{ID: "u-id", Name: uncle, Gender: "male"},
		{ID: "s-id", Name: son, Gender: "male"},
		{ID: "d-id", Name: daughter, Gender: "female"},
		{ID: "c-id", Name: cousin, Gender: "male"},
	}

	for _, member := range members {
		fmt.Printf("Created member: %s (%s)\n", member.Name, member.Gender)
	}

	// Create relationships
	fmt.Println("\n=== Creating Relationships ===")
	
	// Grandparents are spouses
	fmt.Println("Venkatesh ←→ Lakshmi (SPOUSE_OF)")
	
	// Grandparents to parents
	fmt.Println("Venkatesh → Rajesh (PARENT_OF)")
	fmt.Println("Lakshmi → Rajesh (PARENT_OF)")
	fmt.Println("Venkatesh → Kumar (PARENT_OF)")
	fmt.Println("Lakshmi → Kumar (PARENT_OF)")
	
	// Parents are spouses
	fmt.Println("Rajesh ←→ Priya (SPOUSE_OF)")
	
	// Parents to children
	fmt.Println("Rajesh → Arjun (PARENT_OF)")
	fmt.Println("Priya → Arjun (PARENT_OF)")
	fmt.Println("Rajesh → Ananya (PARENT_OF)")
	fmt.Println("Priya → Ananya (PARENT_OF)")
	
	// Uncle to cousin
	fmt.Println("Kumar → Rohan (PARENT_OF)")
	
	// Siblings
	fmt.Println("Rajesh ←→ Kumar (SIBLING_OF)")

	// Query relationships
	fmt.Println("\n=== Querying Relationships ===")
	
	// Arjun to Venkatesh (grandfather)
	fmt.Println("\n1. Arjun to Venkatesh:")
	result1, _ := relationshipService.FindRelationship(ctx, "s-id", "gf-id")
	printRelationshipResult(result1)
	
	// Arjun to Kumar (uncle)
	fmt.Println("\n2. Arjun to Kumar:")
	result2, _ := relationshipService.FindRelationship(ctx, "s-id", "u-id")
	printRelationshipResult(result2)
	
	// Arjun to Rohan (cousin)
	fmt.Println("\n3. Arjun to Rohan:")
	result3, _ := relationshipService.FindRelationship(ctx, "s-id", "c-id")
	printRelationshipResult(result3)
	
	// Arjun to Priya (mother)
	fmt.Println("\n4. Arjun to Priya:")
	result4, _ := relationshipService.FindRelationship(ctx, "s-id", "m-id")
	printRelationshipResult(result4)
	
	// Arjun to Ananya (sister)
	fmt.Println("\n5. Arjun to Ananya:")
	result5, _ := relationshipService.FindRelationship(ctx, "s-id", "d-id")
	printRelationshipResult(result5)

	// Note: In a real implementation, you would actually save these to Neo4j
	// using memberRepo.Create() and relationshipRepo.Create()
	_ = memberRepo
	_ = relationshipRepo
}

// createExampleMember is a helper function to create example member names
func createExampleMember(name, gender string) string {
	return name
}
