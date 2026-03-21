# Relationship Repository - Quick Reference

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  RelationshipRepository                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Create(ctx, rel)                                           │
│    ├─ Bidirectional: SPOUSE_OF, SIBLING_OF                 │
│    │   └─ Creates edges in BOTH directions                 │
│    └─ Directed: PARENT_OF                                   │
│        └─ Creates edge in ONE direction                     │
│                                                              │
│  GetAll(ctx)                                                │
│    ├─ Fetches all relationships                            │
│    └─ Deduplicates bidirectional relationships             │
│                                                              │
│  Delete(ctx, fromID, toID, relType)                        │
│    ├─ Bidirectional: Deletes BOTH edges                    │
│    └─ Directed: Deletes ONE edge                           │
│                                                              │
│  FindPath(ctx, fromID, toID)                               │
│    ├─ Uses Cypher shortestPath()                           │
│    ├─ Returns PathNode[] with ID, Name, Gender             │
│    ├─ Returns relationship types in path                    │
│    └─ Returns nil if no path exists                        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Relationship Types

| Type | Direction | Example |
|------|-----------|---------|
| `SPOUSE_OF` | ↔️ Bidirectional | Husband ↔️ Wife |
| `SIBLING_OF` | ↔️ Bidirectional | Brother ↔️ Sister |
| `PARENT_OF` | → Directed | Parent → Child |

## Key Features

### 1. Bidirectional Relationship Handling
```
Create SPOUSE_OF(A, B):
  A ─[SPOUSE_OF]→ B
  B ─[SPOUSE_OF]→ A

GetAll():
  Returns only: A ─[SPOUSE_OF]→ B
  (Deduplicates using lexicographic key)

Delete SPOUSE_OF(A, B):
  Removes both edges
```

### 2. Shortest Path Finding
```cypher
MATCH path = shortestPath((from)-[*]-(to))
```

Returns:
```json
{
  "nodes": [
    {"id": "1", "name": "Arjun", "gender": "male"},
    {"id": "2", "name": "Rajesh", "gender": "male"}
  ],
  "relationships": ["PARENT_OF"],
  "length": 1
}
```

### 3. Soft Delete Awareness
All queries filter out deleted members:
```cypher
WHERE from.isDeleted = false AND to.isDeleted = false
```

## Usage Patterns

### Creating Relationships
```go
// Parent → Child (directed)
rel := models.NewRelationship(
    models.RelationshipTypeParentOf,
    "parent-id",
    "child-id",
)
repo.Create(ctx, rel)

// Spouse (bidirectional)
rel := models.NewRelationship(
    models.RelationshipTypeSpouseOf,
    "person1-id",
    "person2-id",
)
repo.Create(ctx, rel)
```

### Finding Paths
```go
path, err := repo.FindPath(ctx, "member1-id", "member2-id")
if path == nil {
    // Not related
} else {
    // path.Nodes = connecting members
    // path.Relationships = edge types
    // path.Length = number of hops
}
```

### Analyzing Paths
```go
// Direct relationship (1 hop)
if path.Length == 1 {
    relType := path.Relationships[0]
    // PARENT_OF, SPOUSE_OF, or SIBLING_OF
}

// Multi-hop relationship
for i := 0; i < len(path.Nodes)-1; i++ {
    from := path.Nodes[i]
    to := path.Nodes[i+1]
    rel := path.Relationships[i]
    // Analyze each hop for kinship term mapping
}
```

## Data Structures

### PathNode
```go
type PathNode struct {
    ID     string  // Member UUID
    Name   string  // Member name
    Gender string  // "male", "female", "other"
}
```

### RelationshipPath
```go
type RelationshipPath struct {
    Nodes         []PathNode  // Ordered list of members in path
    Relationships []string    // Ordered list of edge types
    Length        int         // Number of edges (hops)
}
```

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Members not found | Returns error |
| Members deleted | Returns error |
| No path exists | Returns `nil, nil` (not error) |
| Invalid relationship type | Caught by model validation |
| Database error | Returns wrapped error |

## Performance Considerations

1. **Indexes**: Relies on `member_id_index` for fast lookups
2. **Shortest Path**: Uses Neo4j's optimized `shortestPath()` algorithm
3. **Deduplication**: O(n) with hash map for bidirectional relationships
4. **Caching**: Results should be cached at service layer (TTL: 10 minutes)

## Integration Points

### Service Layer
```go
// RelationshipService will use this repository
type RelationshipService struct {
    repo  *RelationshipRepository
    cache *CacheService
}

// Add caching layer
func (s *RelationshipService) FindPath(ctx, from, to) {
    // Check cache first
    // Call repo.FindPath()
    // Map to kinship terms
    // Cache result
}
```

### Handler Layer
```go
// API endpoint: GET /api/relationships/path?from=uuid1&to=uuid2
func (h *RelationshipHandler) FindPath(c *fiber.Ctx) {
    path := h.service.FindPath(ctx, fromID, toID)
    // Return path with kinship term mapping
}
```

## Testing Strategy

### Unit Tests ✅
- Bidirectional detection
- Key generation
- Time parsing
- Structure validation

### Integration Tests (Future)
- Create relationships in Neo4j
- Query relationships
- Find paths in real graph
- Delete relationships

### Property-Based Tests (Future)
- Path symmetry (if A→B exists, B→A should exist for bidirectional)
- Path transitivity (if A→B and B→C, then path A→C should exist)
- Deduplication correctness

## Files Created

1. ✅ `relationship_repo.go` - Main implementation
2. ✅ `relationship_repo_test.go` - Unit tests
3. ✅ `relationship_example.go` - Usage examples
4. ✅ `TASK_4.3_SUMMARY.md` - Task completion summary
5. ✅ `RELATIONSHIP_REPO_OVERVIEW.md` - This file
6. ✅ Updated `README.md` - Documentation

## Requirements Satisfied

- ✅ **Requirement 2.4**: Create and delete relationships
- ✅ **Requirement 2.5**: Store as Neo4j edges with proper types
- ✅ **Requirement 4.1**: Find shortest path between members
- ✅ Accept Neo4j client as dependency
- ✅ Implement Create, GetAll, Delete, FindPath methods
- ✅ Handle bidirectional relationships
- ✅ Return proper error messages
- ✅ Use context for all operations

## Next Steps

1. **Service Layer** (Task 4.4): Implement RelationshipService with:
   - Kinship term mapping
   - Cache integration
   - Business logic validation

2. **Handler Layer**: Create API endpoints:
   - POST /api/relationships
   - GET /api/relationships
   - DELETE /api/relationships/:id
   - GET /api/relationships/path

3. **WebSocket Integration**: Broadcast relationship changes to connected clients

4. **Cache Layer**: Add Redis caching for path queries
