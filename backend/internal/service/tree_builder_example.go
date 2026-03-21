// +build ignore

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/pkg/neo4j"

	"gorm.io/gorm"
)

// ExampleTreeBuilderUsage demonstrates how to use the TreeBuilder service
func ExampleTreeBuilderUsage(neo4jClient *neo4j.Client, db *gorm.DB) {
	// 1. Initialize repositories
	memberRepo := repository.NewMemberRepository(neo4jClient)
	relationshipRepo := repository.NewRelationshipRepository(neo4jClient)
	eventRepo := repository.NewEventRepository(db)

	// 2. Create tree builder
	treeBuilder := NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)

	// 3. Build the family tree
	ctx := context.Background()
	familyTree, err := treeBuilder.BuildTree(ctx)
	if err != nil {
		log.Fatalf("Failed to build family tree: %v", err)
	}

	// 4. Use the tree data
	fmt.Printf("Family tree generated with %d nodes and %d edges\n", 
		len(familyTree.Nodes), len(familyTree.Edges))

	// 5. Convert to JSON for API response
	treeJSON, err := json.MarshalIndent(familyTree, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal tree to JSON: %v", err)
	}

	fmt.Println("Family Tree JSON:")
	fmt.Println(string(treeJSON))

	// 6. Example: Find nodes with upcoming events
	var upcomingEventMembers []string
	for _, node := range familyTree.Nodes {
		if node.Data.HasUpcomingEvent {
			upcomingEventMembers = append(upcomingEventMembers, node.Data.Name)
		}
	}
	fmt.Printf("Members with upcoming events: %v\n", upcomingEventMembers)

	// 7. Example: Count relationships by type
	edgeTypeCounts := make(map[string]int)
	for _, edge := range familyTree.Edges {
		color := edge.Style["stroke"]
		switch color {
		case SpouseEdgeColor:
			edgeTypeCounts["spouse"]++
		case ParentEdgeColor:
			edgeTypeCounts["parent-child"]++
		case SiblingEdgeColor:
			edgeTypeCounts["sibling"]++
		}
	}
	fmt.Printf("Relationship counts: %v\n", edgeTypeCounts)
}

// ExampleTreeBuilderWithCaching demonstrates integration with caching
func ExampleTreeBuilderWithCaching(
	neo4jClient *neo4j.Client, 
	db *gorm.DB,
	cacheService *CacheService,
	userID uint,
) (*FamilyTree, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("family_tree:%d", userID)

	// 1. Try to get from cache
	var cachedTree FamilyTree
	err := cacheService.Get(ctx, cacheKey, &cachedTree)
	if err == nil {
		fmt.Println("Cache hit: returning cached family tree")
		return &cachedTree, nil
	}

	// 2. Cache miss - build tree
	fmt.Println("Cache miss: building family tree")
	memberRepo := repository.NewMemberRepository(neo4jClient)
	relationshipRepo := repository.NewRelationshipRepository(neo4jClient)
	eventRepo := repository.NewEventRepository(db)

	treeBuilder := NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)
	familyTree, err := treeBuilder.BuildTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build tree: %w", err)
	}

	// 3. Store in cache (5 minutes TTL as per requirements)
	err = cacheService.Set(ctx, cacheKey, familyTree, 5*60) // 5 minutes
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to cache tree: %v\n", err)
	}

	return familyTree, nil
}

// ExampleReactFlowIntegration shows the expected JSON structure for React Flow
func ExampleReactFlowIntegration() {
	// This is the structure that will be sent to the frontend
	exampleTree := FamilyTree{
		Nodes: []ReactFlowNode{
			{
				ID:   "member-1",
				Type: "memberNode",
				Position: Position{
					X: 0,
					Y: 0,
				},
				Data: MemberNodeData{
					ID:               "member-1",
					Name:             "Rajesh Kumar",
					AvatarURL:        "https://example.com/avatar1.jpg",
					RelationBadge:    "Father",
					HasUpcomingEvent: true,
					Gender:           "male",
				},
			},
			{
				ID:   "member-2",
				Type: "memberNode",
				Position: Position{
					X: 200,
					Y: 0,
				},
				Data: MemberNodeData{
					ID:               "member-2",
					Name:             "Lakshmi Kumar",
					AvatarURL:        "https://example.com/avatar2.jpg",
					RelationBadge:    "Mother",
					HasUpcomingEvent: false,
					Gender:           "female",
				},
			},
			{
				ID:   "member-3",
				Type: "memberNode",
				Position: Position{
					X: 100,
					Y: 250,
				},
				Data: MemberNodeData{
					ID:               "member-3",
					Name:             "Arjun Kumar",
					AvatarURL:        "https://example.com/avatar3.jpg",
					RelationBadge:    "Son",
					HasUpcomingEvent: false,
					Gender:           "male",
				},
			},
		},
		Edges: []ReactFlowEdge{
			{
				ID:       "member-1-member-2-SPOUSE_OF",
				Source:   "member-1",
				Target:   "member-2",
				Type:     "bezier",
				Animated: false,
				Style: map[string]string{
					"stroke":      SpouseEdgeColor,
					"strokeWidth": "2",
				},
			},
			{
				ID:       "member-1-member-3-PARENT_OF",
				Source:   "member-1",
				Target:   "member-3",
				Type:     "bezier",
				Animated: false,
				Style: map[string]string{
					"stroke":      ParentEdgeColor,
					"strokeWidth": "2",
				},
			},
			{
				ID:       "member-2-member-3-PARENT_OF",
				Source:   "member-2",
				Target:   "member-3",
				Type:     "bezier",
				Animated: false,
				Style: map[string]string{
					"stroke":      ParentEdgeColor,
					"strokeWidth": "2",
				},
			},
		},
	}

	// Convert to JSON
	treeJSON, _ := json.MarshalIndent(exampleTree, "", "  ")
	fmt.Println("React Flow JSON Structure:")
	fmt.Println(string(treeJSON))
}
