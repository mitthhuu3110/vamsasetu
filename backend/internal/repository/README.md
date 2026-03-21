# Repository Layer

This directory contains data access layer implementations for VamsaSetu.

## Overview

The repository layer provides data access abstractions for:
- **Member Repository**: CRUD operations for family members (Neo4j nodes)
- **Relationship Repository**: CRUD operations for family relationships (Neo4j edges) and path finding

## Member Repository

The `MemberRepository` provides CRUD operations for family members stored in Neo4j graph database.

### Features

- **Create**: Add new member nodes to the graph
- **GetByID**: Retrieve a specific member by UUID
- **GetAll**: Fetch all non-deleted members
- **Update**: Modify member attributes
- **SoftDelete**: Mark members as deleted without removing data
- **EnsureIndexes**: Create database indexes for performance

### Usage Example

```go
import (
    "context"
    "vamsasetu/backend/internal/repository"
    "vamsasetu/backend/pkg/neo4j"
)

// Initialize repository
client, _ := neo4j.NewClient(cfg)
repo := repository.NewMemberRepository(client)

// Ensure indexes are created
ctx := context.Background()
repo.EnsureIndexes(ctx)

// Create a member
member := models.NewMember("John Doe", dateOfBirth, "male")
member.Email = "john@example.com"
err := repo.Create(ctx, member)

// Retrieve a member
member, err := repo.GetByID(ctx, memberID)

// Update a member
member.Name = "John Smith"
member.Update()
err = repo.Update(ctx, member)

// Soft delete a member
err = repo.SoftDelete(ctx, memberID)

// Get all members
members, err := repo.GetAll(ctx)
```

### Database Schema

**Node Label**: `Member`

**Properties**:
- `id` (string): UUID, indexed
- `name` (string): Full name, indexed
- `dateOfBirth` (datetime): Date of birth
- `gender` (string): "male", "female", or "other"
- `email` (string): Email address (optional)
- `phone` (string): Phone number (optional)
- `avatarUrl` (string): Avatar image URL (optional)
- `createdAt` (datetime): Creation timestamp
- `updatedAt` (datetime): Last update timestamp
- `isDeleted` (boolean): Soft delete flag

**Indexes**:
- `member_id_index`: Index on `id` property
- `member_name_index`: Index on `name` property

### Cypher Queries

**Create Member**:
```cypher
CREATE (m:Member {
    id: $id,
    name: $name,
    dateOfBirth: datetime($dateOfBirth),
    gender: $gender,
    email: $email,
    phone: $phone,
    avatarUrl: $avatarUrl,
    createdAt: datetime($createdAt),
    updatedAt: datetime($updatedAt),
    isDeleted: $isDeleted
})
RETURN m
```

**Get Member by ID**:
```cypher
MATCH (m:Member {id: $id})
WHERE m.isDeleted = false
RETURN m.id, m.name, m.dateOfBirth, m.gender, m.email, m.phone, 
       m.avatarUrl, m.createdAt, m.updatedAt, m.isDeleted
```

**Get All Members**:
```cypher
MATCH (m:Member)
WHERE m.isDeleted = false
RETURN m.id, m.name, m.dateOfBirth, m.gender, m.email, m.phone, 
       m.avatarUrl, m.createdAt, m.updatedAt, m.isDeleted
ORDER BY m.name
```

**Update Member**:
```cypher
MATCH (m:Member {id: $id})
SET m.name = $name,
    m.dateOfBirth = datetime($dateOfBirth),
    m.gender = $gender,
    m.email = $email,
    m.phone = $phone,
    m.avatarUrl = $avatarUrl,
    m.updatedAt = datetime($updatedAt)
RETURN m
```

**Soft Delete Member**:
```cypher
MATCH (m:Member {id: $id})
SET m.isDeleted = true,
    m.updatedAt = datetime($updatedAt)
RETURN m
```

### Error Handling

The repository returns descriptive errors:
- `"failed to create member"`: Error during member creation
- `"member not found"`: Member ID doesn't exist or is soft deleted
- `"failed to query member"`: Database query error
- `"failed to update member"`: Error during update operation
- `"failed to soft delete member"`: Error during soft delete

### Testing

Run tests with Neo4j running locally:

```bash
# Start Neo4j (via Docker Compose)
docker-compose up -d neo4j

# Run tests
cd backend
go test ./internal/repository -v

# Run specific test
go test ./internal/repository -v -run TestMemberRepository_Create
```

Test coverage includes:
- Creating members with all fields
- Creating members with minimal fields
- Retrieving members by ID
- Retrieving all members
- Updating member attributes
- Soft deleting members
- Index creation (idempotent)
- Error cases (non-existent IDs)

### Requirements Satisfied

This implementation satisfies the following requirements from the spec:

- **Requirement 2.1**: Store Member in Neo4j with required attributes
- **Requirement 2.2**: Persist Member updates to Neo4j
- **Requirement 2.3**: Perform soft delete and preserve data for audit

### Design Alignment

Follows the design document specifications:
- Uses Neo4j Go Driver v5
- Implements context-based operations
- Creates indexes on `member.id` and `member.name`
- Returns proper error messages
- Handles datetime conversions between Go and Neo4j
- Supports soft delete pattern


## Relationship Repository

The `RelationshipRepository` provides CRUD operations for family relationships stored as edges in Neo4j graph database, plus shortest path finding capabilities.

### Features

- **Create**: Add relationship edges between members (handles bidirectional relationships)
- **GetAll**: Retrieve all relationships (with deduplication for bidirectional types)
- **Delete**: Remove relationships (handles bidirectional deletion)
- **FindPath**: Find shortest path between two members using Cypher

### Usage Example

```go
import (
    "context"
    "vamsasetu/backend/internal/repository"
    "vamsasetu/backend/internal/models"
    "vamsasetu/backend/pkg/neo4j"
)

// Initialize repository
client, _ := neo4j.NewClient(cfg)
repo := repository.NewRelationshipRepository(client)

ctx := context.Background()

// Create a parent-child relationship (directed)
parentChildRel := models.NewRelationship(
    models.RelationshipTypeParentOf,
    "parent-uuid",
    "child-uuid",
)
err := repo.Create(ctx, parentChildRel)

// Create a spouse relationship (bidirectional)
spouseRel := models.NewRelationship(
    models.RelationshipTypeSpouseOf,
    "person1-uuid",
    "person2-uuid",
)
err = repo.Create(ctx, spouseRel)

// Get all relationships
relationships, err := repo.GetAll(ctx)

// Find shortest path between two members
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

### Relationship Types

**Bidirectional** (creates edges in both directions):
- `SPOUSE_OF`: Marriage relationship
- `SIBLING_OF`: Brother/sister relationship

**Directed** (creates edge in one direction):
- `PARENT_OF`: Parent to child relationship

### Database Schema

**Relationship Types**: `SPOUSE_OF`, `PARENT_OF`, `SIBLING_OF`

**Properties**:
- `createdAt` (datetime): Creation timestamp

### Cypher Queries

**Create Bidirectional Relationship**:
```cypher
MATCH (from:Member {id: $fromId})
MATCH (to:Member {id: $toId})
WHERE from.isDeleted = false AND to.isDeleted = false
CREATE (from)-[r1:SPOUSE_OF {createdAt: datetime($createdAt)}]->(to)
CREATE (to)-[r2:SPOUSE_OF {createdAt: datetime($createdAt)}]->(from)
RETURN r1, r2
```

**Create Directed Relationship**:
```cypher
MATCH (from:Member {id: $fromId})
MATCH (to:Member {id: $toId})
WHERE from.isDeleted = false AND to.isDeleted = false
CREATE (from)-[r:PARENT_OF {createdAt: datetime($createdAt)}]->(to)
RETURN r
```

**Get All Relationships**:
```cypher
MATCH (from:Member)-[r]->(to:Member)
WHERE from.isDeleted = false AND to.isDeleted = false
AND type(r) IN ['SPOUSE_OF', 'PARENT_OF', 'SIBLING_OF']
RETURN from.id AS fromId, to.id AS toId, type(r) AS relType, r.createdAt AS createdAt
ORDER BY r.createdAt DESC
```

**Find Shortest Path**:
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

**Delete Bidirectional Relationship**:
```cypher
MATCH (from:Member {id: $fromId})-[r1:SPOUSE_OF]->(to:Member {id: $toId})
MATCH (to)-[r2:SPOUSE_OF]->(from)
DELETE r1, r2
RETURN count(r1) + count(r2) AS deletedCount
```

**Delete Directed Relationship**:
```cypher
MATCH (from:Member {id: $fromId})-[r:PARENT_OF]->(to:Member {id: $toId})
DELETE r
RETURN count(r) AS deletedCount
```

### Path Finding

The `FindPath` method returns a `RelationshipPath` structure:

```go
type PathNode struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Gender string `json:"gender"`
}

type RelationshipPath struct {
    Nodes         []PathNode `json:"nodes"`
    Relationships []string   `json:"relationships"`
    Length        int        `json:"length"`
}
```

**Example Path Result**:
```json
{
  "nodes": [
    {"id": "uuid1", "name": "Arjun", "gender": "male"},
    {"id": "uuid2", "name": "Rajesh", "gender": "male"},
    {"id": "uuid3", "name": "Lakshmi", "gender": "female"}
  ],
  "relationships": ["PARENT_OF", "SPOUSE_OF"],
  "length": 2
}
```

This represents: Arjun → (PARENT_OF) → Rajesh → (SPOUSE_OF) → Lakshmi

### Bidirectional Relationship Deduplication

For bidirectional relationships, the repository automatically:
1. Creates edges in both directions during `Create()`
2. Deduplicates results during `GetAll()` (returns only one direction)
3. Deletes both edges during `Delete()`

Deduplication uses lexicographically ordered IDs to ensure consistency:
```go
// Always returns the same key regardless of direction
key := makeRelationshipKey(rel)
// "aaa-bbb-SPOUSE_OF" for both (aaa→bbb) and (bbb→aaa)
```

### Error Handling

The repository returns descriptive errors:
- `"failed to create relationship"`: Error during creation or members not found
- `"relationship not found"`: Relationship doesn't exist
- `"failed to query relationships"`: Database query error
- `"failed to delete relationship"`: Error during deletion
- `"failed to find path"`: Error during path query

Returns `nil` (not error) when no path exists between members.

### Testing

Run tests with Neo4j running locally:

```bash
# Start Neo4j (via Docker Compose)
docker-compose up -d neo4j

# Run tests
cd backend
go test ./internal/repository -v -run TestRelationship

# Run all repository tests
go test ./internal/repository -v
```

Test coverage includes:
- Bidirectional relationship detection
- Relationship key generation
- Neo4j time parsing
- PathNode structure validation
- RelationshipPath structure validation

### Requirements Satisfied

This implementation satisfies the following requirements from the spec:

- **Requirement 2.4**: Create and delete relationships between members
- **Requirement 2.5**: Store relationships as edges in Neo4j with appropriate types
- **Requirement 4.1**: Find shortest path between two members

### Design Alignment

Follows the design document specifications:
- Uses Neo4j Go Driver v5
- Implements context-based operations
- Uses Cypher's `shortestPath()` function for optimal performance
- Returns proper error messages
- Handles bidirectional relationships correctly
- Supports all three relationship types (SPOUSE_OF, PARENT_OF, SIBLING_OF)
- Returns path with nodes and relationship types for kinship term mapping

### Example Files

See `relationship_example.go` for comprehensive usage examples including:
- Creating different relationship types
- Batch relationship creation
- Path finding and analysis
- Error handling patterns
