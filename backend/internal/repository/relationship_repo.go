package repository

import (
	"context"
	"fmt"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/pkg/neo4j"

	neo4jDriver "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// RelationshipRepository handles relationship data operations in Neo4j
type RelationshipRepository struct {
	client *neo4j.Client
}

// PathNode represents a node in a relationship path
type PathNode struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

// RelationshipPath represents the result of a path finding query
type RelationshipPath struct {
	Nodes        []PathNode `json:"nodes"`
	Relationships []string   `json:"relationships"`
	Length       int        `json:"length"`
}

// NewRelationshipRepository creates a new relationship repository instance
func NewRelationshipRepository(client *neo4j.Client) *RelationshipRepository {
	return &RelationshipRepository{
		client: client,
	}
}

// Create creates a new relationship edge in Neo4j
func (r *RelationshipRepository) Create(ctx context.Context, rel *models.Relationship) error {
	session := r.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
		AccessMode: neo4jDriver.AccessModeWrite,
	})
	defer session.Close(ctx)

	// Build the query based on relationship type
	var query string
	if rel.IsBidirectional() {
		// For bidirectional relationships (SPOUSE_OF, SIBLING_OF), create edges in both directions
		query = fmt.Sprintf(`
			MATCH (from:Member {id: $fromId})
			MATCH (to:Member {id: $toId})
			WHERE from.isDeleted = false AND to.isDeleted = false
			CREATE (from)-[r1:%s {createdAt: datetime($createdAt)}]->(to)
			CREATE (to)-[r2:%s {createdAt: datetime($createdAt)}]->(from)
			RETURN r1, r2
		`, rel.Type, rel.Type)
	} else {
		// For directed relationships (PARENT_OF), create edge in one direction
		query = fmt.Sprintf(`
			MATCH (from:Member {id: $fromId})
			MATCH (to:Member {id: $toId})
			WHERE from.isDeleted = false AND to.isDeleted = false
			CREATE (from)-[r:%s {createdAt: datetime($createdAt)}]->(to)
			RETURN r
		`, rel.Type)
	}

	params := map[string]interface{}{
		"fromId":    rel.FromID,
		"toId":      rel.ToID,
		"createdAt": rel.CreatedAt.Format(time.RFC3339),
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create relationship: %w", err)
	}

	if !result.Next(ctx) {
		return fmt.Errorf("failed to create relationship: members not found or already deleted")
	}

	return nil
}

// GetAll retrieves all relationships from Neo4j
func (r *RelationshipRepository) GetAll(ctx context.Context) ([]*models.Relationship, error) {
	session := r.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
		AccessMode: neo4jDriver.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
		MATCH (from:Member)-[r]->(to:Member)
		WHERE from.isDeleted = false AND to.isDeleted = false
		AND type(r) IN ['SPOUSE_OF', 'PARENT_OF', 'SIBLING_OF']
		RETURN from.id AS fromId, to.id AS toId, type(r) AS relType, r.createdAt AS createdAt
		ORDER BY r.createdAt DESC
	`

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships: %w", err)
	}

	var relationships []*models.Relationship
	seen := make(map[string]bool) // To deduplicate bidirectional relationships

	for result.Next(ctx) {
		record := result.Record()
		rel, err := r.recordToRelationship(record)
		if err != nil {
			return nil, err
		}

		// For bidirectional relationships, only include one direction
		if rel.IsBidirectional() {
			key := r.makeRelationshipKey(rel)
			if seen[key] {
				continue
			}
			seen[key] = true
		}

		relationships = append(relationships, rel)
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	return relationships, nil
}

// Delete removes a relationship from Neo4j
func (r *RelationshipRepository) Delete(ctx context.Context, fromID, toID, relType string) error {
	session := r.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
		AccessMode: neo4jDriver.AccessModeWrite,
	})
	defer session.Close(ctx)

	// Check if relationship is bidirectional
	rel := &models.Relationship{Type: relType}
	var query string

	if rel.IsBidirectional() {
		// Delete both directions for bidirectional relationships
		query = fmt.Sprintf(`
			MATCH (from:Member {id: $fromId})-[r1:%s]->(to:Member {id: $toId})
			MATCH (to)-[r2:%s]->(from)
			DELETE r1, r2
			RETURN count(r1) + count(r2) AS deletedCount
		`, relType, relType)
	} else {
		// Delete single direction for directed relationships
		query = fmt.Sprintf(`
			MATCH (from:Member {id: $fromId})-[r:%s]->(to:Member {id: $toId})
			DELETE r
			RETURN count(r) AS deletedCount
		`, relType)
	}

	params := map[string]interface{}{
		"fromId": fromID,
		"toId":   toID,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete relationship: %w", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		deletedCount, _ := record.Get("deletedCount")
		if deletedCount.(int64) == 0 {
			return fmt.Errorf("relationship not found")
		}
	}

	return nil
}

// FindPath finds the shortest path between two members using Cypher
func (r *RelationshipRepository) FindPath(ctx context.Context, fromID, toID string) (*RelationshipPath, error) {
	session := r.client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{
		AccessMode: neo4jDriver.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
		MATCH (from:Member {id: $fromId})
		MATCH (to:Member {id: $toId})
		WHERE from.isDeleted = false AND to.isDeleted = false
		MATCH path = shortestPath((from)-[*]-(to))
		WITH path, [node IN nodes(path) | {id: node.id, name: node.name, gender: node.gender}] AS nodeList,
		     [rel IN relationships(path) | type(rel)] AS relList
		RETURN nodeList, relList, length(path) AS pathLength
		LIMIT 1
	`

	params := map[string]interface{}{
		"fromId": fromID,
		"toId":   toID,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to find path: %w", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		return r.recordToPath(record)
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	// No path found
	return nil, nil
}

// recordToRelationship converts a Neo4j record to a Relationship model
func (r *RelationshipRepository) recordToRelationship(record *neo4jDriver.Record) (*models.Relationship, error) {
	fromID, _ := record.Get("fromId")
	toID, _ := record.Get("toId")
	relType, _ := record.Get("relType")
	createdAtRaw, _ := record.Get("createdAt")

	createdAt, err := r.parseNeo4jTime(createdAtRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse createdAt: %w", err)
	}

	rel := &models.Relationship{
		Type:      relType.(string),
		FromID:    fromID.(string),
		ToID:      toID.(string),
		CreatedAt: createdAt,
	}

	return rel, nil
}

// recordToPath converts a Neo4j record to a RelationshipPath
func (r *RelationshipRepository) recordToPath(record *neo4jDriver.Record) (*RelationshipPath, error) {
	nodeListRaw, _ := record.Get("nodeList")
	relListRaw, _ := record.Get("relList")
	pathLength, _ := record.Get("pathLength")

	// Parse nodes
	nodeList := nodeListRaw.([]interface{})
	nodes := make([]PathNode, len(nodeList))
	for i, nodeRaw := range nodeList {
		nodeMap := nodeRaw.(map[string]interface{})
		nodes[i] = PathNode{
			ID:     nodeMap["id"].(string),
			Name:   nodeMap["name"].(string),
			Gender: nodeMap["gender"].(string),
		}
	}

	// Parse relationships
	relList := relListRaw.([]interface{})
	relationships := make([]string, len(relList))
	for i, relRaw := range relList {
		relationships[i] = relRaw.(string)
	}

	path := &RelationshipPath{
		Nodes:        nodes,
		Relationships: relationships,
		Length:       int(pathLength.(int64)),
	}

	return path, nil
}

// parseNeo4jTime converts Neo4j datetime to Go time.Time
func (r *RelationshipRepository) parseNeo4jTime(value interface{}) (time.Time, error) {
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

// makeRelationshipKey creates a unique key for deduplicating bidirectional relationships
func (r *RelationshipRepository) makeRelationshipKey(rel *models.Relationship) string {
	// Always use lexicographically smaller ID first to ensure consistency
	if rel.FromID < rel.ToID {
		return fmt.Sprintf("%s-%s-%s", rel.FromID, rel.ToID, rel.Type)
	}
	return fmt.Sprintf("%s-%s-%s", rel.ToID, rel.FromID, rel.Type)
}
