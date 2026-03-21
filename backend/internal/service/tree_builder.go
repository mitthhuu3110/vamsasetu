package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
)

// TreeBuilder handles family tree structure generation and layout
type TreeBuilder struct {
	memberRepo       *repository.MemberRepository
	relationshipRepo *repository.RelationshipRepository
	eventRepo        *repository.EventRepository
}

// NewTreeBuilder creates a new tree builder instance
func NewTreeBuilder(
	memberRepo *repository.MemberRepository,
	relationshipRepo *repository.RelationshipRepository,
	eventRepo *repository.EventRepository,
) *TreeBuilder {
	return &TreeBuilder{
		memberRepo:       memberRepo,
		relationshipRepo: relationshipRepo,
		eventRepo:        eventRepo,
	}
}

// ReactFlowNode represents a node in React Flow format
type ReactFlowNode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Position Position               `json:"position"`
	Data     MemberNodeData         `json:"data"`
}

// Position represents x, y coordinates
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// MemberNodeData represents the data payload for a member node
type MemberNodeData struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	AvatarURL        string `json:"avatarUrl"`
	RelationBadge    string `json:"relationBadge"`
	HasUpcomingEvent bool   `json:"hasUpcomingEvent"`
	Gender           string `json:"gender"`
}

// ReactFlowEdge represents an edge in React Flow format
type ReactFlowEdge struct {
	ID       string            `json:"id"`
	Source   string            `json:"source"`
	Target   string            `json:"target"`
	Type     string            `json:"type"`
	Animated bool              `json:"animated"`
	Style    map[string]string `json:"style"`
	Label    string            `json:"label,omitempty"`
}

// FamilyTree represents the complete tree structure
type FamilyTree struct {
	Nodes []ReactFlowNode `json:"nodes"`
	Edges []ReactFlowEdge `json:"edges"`
}

// Edge colors based on relationship type
const (
	SpouseEdgeColor  = "#E11D48" // Rose
	ParentEdgeColor  = "#0D9488" // Teal
	SiblingEdgeColor = "#F59E0B" // Amber
)

// Layout constants
const (
	NodeWidth       = 150.0
	NodeHeight      = 180.0
	HorizontalSpace = 50.0
	VerticalSpace   = 100.0
	GenerationGap   = 250.0
)

// BuildTree generates the complete family tree structure with layout
func (tb *TreeBuilder) BuildTree(ctx context.Context) (*FamilyTree, error) {
	// 1. Query all members and relationships
	members, err := tb.memberRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}

	relationships, err := tb.relationshipRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationships: %w", err)
	}

	// 2. Query upcoming events (within 7 days)
	now := time.Now()
	upcomingEnd := now.AddDate(0, 0, 7)
	upcomingEvents, err := tb.eventRepo.GetUpcoming(ctx, now, upcomingEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}

	// 3. Build member ID to upcoming event map
	upcomingEventMap := tb.buildUpcomingEventMap(upcomingEvents)

	// 4. Build graph structure
	graph := tb.buildGraph(members, relationships)

	// 5. Calculate layout (assign x, y coordinates)
	layout := tb.calculateLayout(graph)

	// 6. Transform to React Flow format
	nodes := tb.buildReactFlowNodes(members, layout, upcomingEventMap)
	edges := tb.buildReactFlowEdges(relationships)

	return &FamilyTree{
		Nodes: nodes,
		Edges: edges,
	}, nil
}

// Graph represents the family tree graph structure
type Graph struct {
	Members       map[string]*models.Member
	Adjacency     map[string][]string // member ID -> connected member IDs
	ParentOf      map[string][]string // parent ID -> child IDs
	ChildOf       map[string][]string // child ID -> parent IDs
	SpouseOf      map[string][]string // member ID -> spouse IDs
	SiblingOf     map[string][]string // member ID -> sibling IDs
	Generations   map[string]int      // member ID -> generation level
	RootNodes     []string            // members with no parents
}

// buildGraph constructs the graph structure from members and relationships
func (tb *TreeBuilder) buildGraph(members []*models.Member, relationships []*models.Relationship) *Graph {
	graph := &Graph{
		Members:     make(map[string]*models.Member),
		Adjacency:   make(map[string][]string),
		ParentOf:    make(map[string][]string),
		ChildOf:     make(map[string][]string),
		SpouseOf:    make(map[string][]string),
		SiblingOf:   make(map[string][]string),
		Generations: make(map[string]int),
		RootNodes:   []string{},
	}

	// Build member map
	for _, member := range members {
		graph.Members[member.ID] = member
		graph.Adjacency[member.ID] = []string{}
	}

	// Build relationship maps
	for _, rel := range relationships {
		switch rel.Type {
		case models.RelationshipTypeParentOf:
			graph.ParentOf[rel.FromID] = append(graph.ParentOf[rel.FromID], rel.ToID)
			graph.ChildOf[rel.ToID] = append(graph.ChildOf[rel.ToID], rel.FromID)
			graph.Adjacency[rel.FromID] = append(graph.Adjacency[rel.FromID], rel.ToID)
			graph.Adjacency[rel.ToID] = append(graph.Adjacency[rel.ToID], rel.FromID)

		case models.RelationshipTypeSpouseOf:
			graph.SpouseOf[rel.FromID] = append(graph.SpouseOf[rel.FromID], rel.ToID)
			graph.Adjacency[rel.FromID] = append(graph.Adjacency[rel.FromID], rel.ToID)

		case models.RelationshipTypeSiblingOf:
			graph.SiblingOf[rel.FromID] = append(graph.SiblingOf[rel.FromID], rel.ToID)
			graph.Adjacency[rel.FromID] = append(graph.Adjacency[rel.FromID], rel.ToID)
		}
	}

	// Identify root nodes (members with no parents)
	for memberID := range graph.Members {
		if len(graph.ChildOf[memberID]) == 0 {
			graph.RootNodes = append(graph.RootNodes, memberID)
		}
	}

	// Calculate generations
	tb.calculateGenerations(graph)

	return graph
}

// calculateGenerations assigns generation levels to all members
func (tb *TreeBuilder) calculateGenerations(graph *Graph) {
	// Start with root nodes at generation 0
	for _, rootID := range graph.RootNodes {
		tb.assignGeneration(graph, rootID, 0, make(map[string]bool))
	}

	// Handle disconnected components (members not connected to any root)
	for memberID := range graph.Members {
		if _, exists := graph.Generations[memberID]; !exists {
			// Assign to generation 0 if not yet assigned
			graph.Generations[memberID] = 0
		}
	}
}

// assignGeneration recursively assigns generation levels
func (tb *TreeBuilder) assignGeneration(graph *Graph, memberID string, generation int, visited map[string]bool) {
	if visited[memberID] {
		return
	}
	visited[memberID] = true

	// Set generation (use minimum if already set)
	if existingGen, exists := graph.Generations[memberID]; exists {
		if generation < existingGen {
			graph.Generations[memberID] = generation
		}
	} else {
		graph.Generations[memberID] = generation
	}

	// Assign same generation to spouses
	for _, spouseID := range graph.SpouseOf[memberID] {
		if !visited[spouseID] {
			tb.assignGeneration(graph, spouseID, generation, visited)
		}
	}

	// Assign same generation to siblings
	for _, siblingID := range graph.SiblingOf[memberID] {
		if !visited[siblingID] {
			tb.assignGeneration(graph, siblingID, generation, visited)
		}
	}

	// Assign next generation to children
	for _, childID := range graph.ParentOf[memberID] {
		if !visited[childID] {
			tb.assignGeneration(graph, childID, generation+1, visited)
		}
	}
}

// NodeLayout represents the calculated position for a node
type NodeLayout struct {
	X          float64
	Y          float64
	Generation int
}

// calculateLayout computes x, y coordinates for all members using hierarchical layout
func (tb *TreeBuilder) calculateLayout(graph *Graph) map[string]NodeLayout {
	layout := make(map[string]NodeLayout)

	// Group members by generation
	generationGroups := make(map[int][]string)
	for memberID, generation := range graph.Generations {
		generationGroups[generation] = append(generationGroups[generation], memberID)
	}

	// Sort generations
	var generations []int
	for gen := range generationGroups {
		generations = append(generations, gen)
	}
	sort.Ints(generations)

	// Layout each generation
	for _, generation := range generations {
		members := generationGroups[generation]
		
		// Sort members within generation for consistent ordering
		sort.Slice(members, func(i, j int) bool {
			// Sort by name for now (could be enhanced with sibling order)
			return graph.Members[members[i]].Name < graph.Members[members[j]].Name
		})

		// Calculate Y position (generation-based)
		y := float64(generation) * GenerationGap

		// Calculate X positions with spacing
		tb.layoutGeneration(graph, members, generation, y, layout)
	}

	// Apply collision detection and spacing optimization
	tb.optimizeSpacing(layout, graph)

	return layout
}

// layoutGeneration positions members within a single generation
func (tb *TreeBuilder) layoutGeneration(graph *Graph, members []string, generation int, y float64, layout map[string]NodeLayout) {
	// Group members by family clusters (spouses and siblings together)
	clusters := tb.identifyClusters(graph, members)

	currentX := 0.0

	for _, cluster := range clusters {
		// Calculate cluster width
		clusterWidth := float64(len(cluster)) * (NodeWidth + HorizontalSpace)

		// Position each member in the cluster
		for i, memberID := range cluster {
			x := currentX + float64(i)*(NodeWidth+HorizontalSpace)
			layout[memberID] = NodeLayout{
				X:          x,
				Y:          y,
				Generation: generation,
			}
		}

		// Move to next cluster position
		currentX += clusterWidth + HorizontalSpace*2
	}
}

// identifyClusters groups members who should be positioned together
func (tb *TreeBuilder) identifyClusters(graph *Graph, members []string) [][]string {
	visited := make(map[string]bool)
	var clusters [][]string

	for _, memberID := range members {
		if visited[memberID] {
			continue
		}

		// Start a new cluster
		cluster := []string{memberID}
		visited[memberID] = true

		// Add spouses to cluster
		for _, spouseID := range graph.SpouseOf[memberID] {
			if !visited[spouseID] && tb.contains(members, spouseID) {
				cluster = append(cluster, spouseID)
				visited[spouseID] = true
			}
		}

		// Add siblings to cluster (if they're in the same generation group)
		for _, siblingID := range graph.SiblingOf[memberID] {
			if !visited[siblingID] && tb.contains(members, siblingID) {
				cluster = append(cluster, siblingID)
				visited[siblingID] = true
			}
		}

		clusters = append(clusters, cluster)
	}

	return clusters
}

// contains checks if a slice contains a string
func (tb *TreeBuilder) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// optimizeSpacing adjusts positions to prevent overlaps and improve visual balance
func (tb *TreeBuilder) optimizeSpacing(layout map[string]NodeLayout, graph *Graph) {
	// Group by generation for collision detection
	generationNodes := make(map[int][]string)
	for memberID, nodeLayout := range layout {
		generationNodes[nodeLayout.Generation] = append(generationNodes[nodeLayout.Generation], memberID)
	}

	// Check each generation for overlaps
	for generation, members := range generationNodes {
		// Sort by X position
		sort.Slice(members, func(i, j int) bool {
			return layout[members[i]].X < layout[members[j]].X
		})

		// Detect and resolve overlaps
		for i := 0; i < len(members)-1; i++ {
			currentID := members[i]
			nextID := members[i+1]

			currentLayout := layout[currentID]
			nextLayout := layout[nextID]

			// Calculate minimum required spacing
			minSpacing := NodeWidth + HorizontalSpace

			// Check for overlap
			if nextLayout.X-currentLayout.X < minSpacing {
				// Shift all subsequent nodes to the right
				shift := minSpacing - (nextLayout.X - currentLayout.X)
				for j := i + 1; j < len(members); j++ {
					shiftID := members[j]
					shiftLayout := layout[shiftID]
					layout[shiftID] = NodeLayout{
						X:          shiftLayout.X + shift,
						Y:          shiftLayout.Y,
						Generation: generation,
					}
				}
			}
		}
	}

	// Center the tree horizontally
	tb.centerTree(layout)
}

// centerTree shifts all nodes to center the tree around x=0
func (tb *TreeBuilder) centerTree(layout map[string]NodeLayout) {
	if len(layout) == 0 {
		return
	}

	// Find min and max X
	minX := layout[tb.getFirstKey(layout)].X
	maxX := minX

	for _, nodeLayout := range layout {
		if nodeLayout.X < minX {
			minX = nodeLayout.X
		}
		if nodeLayout.X > maxX {
			maxX = nodeLayout.X
		}
	}

	// Calculate center offset
	centerOffset := -(minX + maxX) / 2

	// Apply offset to all nodes
	for memberID, nodeLayout := range layout {
		layout[memberID] = NodeLayout{
			X:          nodeLayout.X + centerOffset,
			Y:          nodeLayout.Y,
			Generation: nodeLayout.Generation,
		}
	}
}

// getFirstKey returns the first key from a map (helper for finding initial values)
func (tb *TreeBuilder) getFirstKey(layout map[string]NodeLayout) string {
	for key := range layout {
		return key
	}
	return ""
}

// buildUpcomingEventMap creates a map of member IDs to whether they have upcoming events
func (tb *TreeBuilder) buildUpcomingEventMap(events []*models.Event) map[string]bool {
	eventMap := make(map[string]bool)

	for _, event := range events {
		for _, memberID := range event.GetMemberIDs() {
			eventMap[memberID] = true
		}
	}

	return eventMap
}

// buildReactFlowNodes transforms members and layout into React Flow node format
func (tb *TreeBuilder) buildReactFlowNodes(
	members []*models.Member,
	layout map[string]NodeLayout,
	upcomingEventMap map[string]bool,
) []ReactFlowNode {
	var nodes []ReactFlowNode

	for _, member := range members {
		nodeLayout, exists := layout[member.ID]
		if !exists {
			// Skip members without layout (shouldn't happen)
			continue
		}

		node := ReactFlowNode{
			ID:   member.ID,
			Type: "memberNode",
			Position: Position{
				X: nodeLayout.X,
				Y: nodeLayout.Y,
			},
			Data: MemberNodeData{
				ID:               member.ID,
				Name:             member.Name,
				AvatarURL:        member.AvatarURL,
				RelationBadge:    "", // Can be enhanced with relationship context
				HasUpcomingEvent: upcomingEventMap[member.ID],
				Gender:           member.Gender,
			},
		}

		nodes = append(nodes, node)
	}

	return nodes
}

// buildReactFlowEdges transforms relationships into React Flow edge format
func (tb *TreeBuilder) buildReactFlowEdges(relationships []*models.Relationship) []ReactFlowEdge {
	var edges []ReactFlowEdge
	edgeIDMap := make(map[string]bool) // To avoid duplicate edges

	for _, rel := range relationships {
		// Create edge ID
		edgeID := fmt.Sprintf("%s-%s-%s", rel.FromID, rel.ToID, rel.Type)

		// Skip if already added (for bidirectional relationships)
		if edgeIDMap[edgeID] {
			continue
		}

		// Determine edge color based on relationship type
		color := tb.getEdgeColor(rel.Type)

		edge := ReactFlowEdge{
			ID:       edgeID,
			Source:   rel.FromID,
			Target:   rel.ToID,
			Type:     "bezier",
			Animated: false,
			Style: map[string]string{
				"stroke":      color,
				"strokeWidth": "2",
			},
		}

		edges = append(edges, edge)
		edgeIDMap[edgeID] = true

		// For bidirectional relationships, add reverse edge ID to map
		if rel.Type == models.RelationshipTypeSpouseOf || rel.Type == models.RelationshipTypeSiblingOf {
			reverseEdgeID := fmt.Sprintf("%s-%s-%s", rel.ToID, rel.FromID, rel.Type)
			edgeIDMap[reverseEdgeID] = true
		}
	}

	return edges
}

// getEdgeColor returns the color for a relationship type
func (tb *TreeBuilder) getEdgeColor(relType string) string {
	switch relType {
	case models.RelationshipTypeSpouseOf:
		return SpouseEdgeColor
	case models.RelationshipTypeParentOf:
		return ParentEdgeColor
	case models.RelationshipTypeSiblingOf:
		return SiblingEdgeColor
	default:
		return "#6B7280" // Gray fallback
	}
}
