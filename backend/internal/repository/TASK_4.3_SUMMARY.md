# Task 4.3: Neo4j Relationship Repository Implementation

## Overview
Implemented the Neo4j relationship repository with full CRUD operations and shortest path finding capabilities.

## Files Created

### 1. `relationship_repo.go`
Complete relationship repository implementation with:

#### Core Methods
- **Create(ctx, rel)**: Creates relationship edges in Neo4j
  - Handles bidirectional relationships (SPOUSE_OF, SIBLING_OF) by creating edges in both directions
  - Handles directed relationships (PARENT_OF) with single direction
  - Validates that both members exist and are not deleted
  - Returns error if members not found

- **GetAll(ctx)**: Retrieves all relationships
  - Filters out deleted members
  - Deduplicates bidirectional relationships (only returns one direction)
  - Orders by creation date (newest first)
  - Returns empty slice if no relationships exist

- **Delete(ctx, fromID, toID, relType)**: Removes relationships
  - Deletes both directions for bidirectional relationships
  - Deletes single direction for directed relationships
  - Returns error if relationship not found
  - Returns count of deleted relationships

- **FindPath(ctx, fromID, toID)**: Finds shortest path between members
  - Uses Cypher's `shortestPath()` function for optimal performance
  - Returns path nodes with ID, name, and gender
  - Returns relationship types in the path
  - Returns path length
  - Returns nil if no path exists (members not related)

#### Supporting Types
- **PathNode**: Represents a node in the relationship path
  - ID: Member UUID
  - Name: Member name
  - Gender: Member gender (needed for kinship term mapping)

- **RelationshipPath**: Complete path result
  - Nodes: Array of PathNode objects
  - Relationships: Array of relationship type strings
  - Length: Number of edges in the path

#### Helper Methods
- **recordToRelationship()**: Converts Neo4j record to Relationship model
- **recordToPath()**: Converts Neo4j record to RelationshipPath
- **parseNeo4jTime()**: Handles Neo4j datetime conversion
- **makeRelationshipKey()**: Creates consistent keys for deduplication

### 2. `relationship_repo_test.go`
Comprehensive unit tests covering:
- Bidirectional relationship detection
- Relationship key generation for deduplication
- Neo4j time parsing
- PathNode structure validation
- RelationshipPath structure validation

## Key Implementation Details

### Bidirectional Relationship Handling
```go
if rel.IsBidirectional() {
    // Create edges in both directions
    CREATE (from)-[r1:SPOUSE_OF]->(to)
    CREATE (to)-[r2:SPOUSE_OF]->(from)
}
```

### Shortest Path Cypher Query
```cypher
MATCH (from:Member {id: $fromId})
MATCH (to:Member {id: $toId})
WHERE from.isDeleted = false AND to.isDeleted = false
MATCH path = shortestPath((from)-[*]-(to))
WITH path, 
     [node IN nodes(path) | {id: node.id, name: node.name, gender: node.gender}] AS nodeList,
     [rel IN relationships(path) | type(rel)] AS relList
RETURN nodeList, relList, length(path) AS pathLength
LIMIT 1
```

### Deduplication Strategy
For bidirectional relationships, we create a consistent key using lexicographically ordered IDs:
```go
func makeRelationshipKey(rel *Relationship) string {
    if rel.FromID < rel.ToID {
        return fmt.Sprintf("%s-%s-%s", rel.FromID, rel.ToID, rel.Type)
    }
    return fmt.Sprintf("%s-%s-%s", rel.ToID, rel.FromID, rel.Type)
}
```

## Error Handling
- Returns descriptive errors for all failure cases
- Validates member existence before creating relationships
- Handles missing paths gracefully (returns nil, not error)
- Proper context propagation for cancellation

## Requirements Satisfied

### Requirement 2.4: Family Tree Data Management
✅ Create relationships between members
✅ Delete relationships
✅ Store relationships as edges in Neo4j
✅ Support SPOUSE_OF, PARENT_OF, SIBLING_OF types

### Requirement 2.5: Relationship Validation
✅ Validate relationship types
✅ Prevent relationships with deleted members
✅ Handle bidirectional relationships correctly

### Requirement 4.1: Relationship Path Finding
✅ Find shortest path between two members
✅ Return path nodes with member details
✅ Return relationship types in the path
✅ Handle "not related" case (no path exists)

## Usage Example

```go
// Initialize repository
repo := NewRelationshipRepository(neo4jClient)

// Create a parent-child relationship
rel := models.NewRelationship(
    models.RelationshipTypeParentOf,
    "parent-uuid",
    "child-uuid",
)
err := repo.Create(ctx, rel)

// Create a spouse relationship (bidirectional)
spouseRel := models.NewRelationship(
    models.RelationshipTypeSpouseOf,
    "person1-uuid",
    "person2-uuid",
)
err = repo.Create(ctx, spouseRel)

// Get all relationships
relationships, err := repo.GetAll(ctx)

// Find path between two members
path, err := repo.FindPath(ctx, "member1-uuid", "member2-uuid")
if path == nil {
    // Members are not related
} else {
    // path.Nodes contains the connecting members
    // path.Relationships contains the edge types
    // path.Length is the number of hops
}

// Delete a relationship
err = repo.Delete(ctx, "from-uuid", "to-uuid", models.RelationshipTypeParentOf)
```

## Testing
All unit tests pass with no syntax errors or diagnostics issues.

## Next Steps
This repository is ready for integration with:
- Relationship service layer (Task 4.4)
- Kinship term mapping engine
- Cache service for path query optimization
- WebSocket updates for real-time relationship changes
