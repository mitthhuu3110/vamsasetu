# Service Layer

This directory contains the business logic layer for VamsaSetu. Services encapsulate complex operations, orchestrate repository calls, and implement domain-specific algorithms.

## Relationship Service

The `RelationshipService` is the core component that implements the **Relationship Engine** - the algorithm that computes kinship relationships between family members using graph traversal and rule-based mapping.

### Features

- **Path Finding**: Uses Neo4j's shortest path algorithm to find connections between family members
- **Kinship Mapping**: Maps relationship paths to culturally appropriate Indian kinship terms
- **Multi-hop Support**: Handles complex relationships spanning multiple generations
- **Natural Language**: Generates human-readable descriptions of relationships

### Architecture

```
FindRelationship(fromID, toID)
    ↓
Query Neo4j for shortest path
    ↓
Map path to kinship terms
    ↓
Generate natural language description
    ↓
Return RelationshipResult
```

### Kinship Mapping Rules

#### Direct Relationships (1 hop)

| Relationship Type | Target Gender | English Term | Telugu Term |
|------------------|---------------|--------------|-------------|
| PARENT_OF | Male | Son | Koduku |
| PARENT_OF | Female | Daughter | Kuthuru |
| SPOUSE_OF | Male | Husband | Menarikam |
| SPOUSE_OF | Female | Wife | Bharya |
| SIBLING_OF | Male | Brother | Annayya/Tammudu |
| SIBLING_OF | Female | Sister | Akka/Chelli |

#### Two-Hop Relationships

| Path Pattern | Middle Gender | Target Gender | English Term | Telugu Term |
|-------------|---------------|---------------|--------------|-------------|
| PARENT_OF → PARENT_OF | Any | Male | Grandson | Manumadu |
| PARENT_OF → PARENT_OF | Any | Female | Granddaughter | Manumalu |
| PARENT_OF → SIBLING_OF | Male | Male | Uncle (Father's Brother) | Babai |
| PARENT_OF → SIBLING_OF | Male | Female | Aunt (Father's Sister) | Attha |
| PARENT_OF → SIBLING_OF | Female | Male | Uncle (Mother's Brother) | Mamayya |
| PARENT_OF → SIBLING_OF | Female | Female | Aunt (Mother's Sister) | Pinni |
| SPOUSE_OF → PARENT_OF | Any | Male | Father-in-Law | Maamayyagaru |
| SPOUSE_OF → PARENT_OF | Any | Female | Mother-in-Law | Attagaru |
| SIBLING_OF → SPOUSE_OF | Any | Male | Brother-in-Law | Bava |
| SIBLING_OF → SPOUSE_OF | Any | Female | Sister-in-Law | Vadina |
| SIBLING_OF → PARENT_OF | Any | Male | Nephew | Bhanja/Alludu |
| SIBLING_OF → PARENT_OF | Any | Female | Niece | Bhanjika/Kodalu |

#### Multi-Hop Relationships (3+ hops)

- **Great-Grandparents**: 3+ PARENT_OF edges going up the tree
- **Great-Grandchildren**: 3+ PARENT_OF edges going down the tree
- **Cousins**: Paths that go up (PARENT_OF), across (SIBLING_OF), and down (PARENT_OF)
  - First cousin: 1 generation up, sibling connection, 1 generation down
  - Second cousin: 2 generations up, sibling connection, 2 generations down

### Usage Example

```go
// Initialize service
relationshipRepo := repository.NewRelationshipRepository(neo4jClient)
relationshipService := service.NewRelationshipService(relationshipRepo)

// Find relationship between two members
result, err := relationshipService.FindRelationship(ctx, "member-id-1", "member-id-2")
if err != nil {
    log.Fatal(err)
}

// Result contains:
// - Path: Array of nodes in the relationship path
// - RelationLabel: English kinship term (e.g., "Uncle (Father's Brother)")
// - KinshipTerm: Telugu kinship term (e.g., "Babai")
// - Description: Natural language description (e.g., "Rajesh is Arjun's uncle through his father")
```

### API Response Format

```json
{
  "path": [
    {"id": "uuid-1", "name": "Arjun", "gender": "male"},
    {"id": "uuid-2", "name": "Rajesh", "gender": "male"},
    {"id": "uuid-3", "name": "Kumar", "gender": "male"}
  ],
  "relationLabel": "Uncle (Father's Brother)",
  "kinshipTerm": "Babai",
  "description": "Kumar is Arjun's uncle (father's brother) through Rajesh"
}
```

### Special Cases

1. **Same Person**: Returns "Self" with Telugu term "Nenu/Nuvvu"
2. **No Path Found**: Returns "Not Related" with empty kinship term
3. **Complex Paths**: For paths longer than 3 hops, provides generic descriptions with hop count

### Testing

The service includes comprehensive unit tests covering:
- Direct relationships (parent, spouse, sibling)
- Two-hop relationships (grandparent, uncle/aunt, in-laws, nephew/niece)
- Multi-hop relationships (great-grandparents, cousins)
- Edge cases (same person, no path)
- Helper functions (prefix generation, cousin degree calculation)

Run tests:
```bash
cd backend
go test -v ./internal/service/...
```

### Future Enhancements

1. **Age-based Sibling Terms**: Distinguish between older (Anna/Akka) and younger (Tammudu/Chelli) siblings
2. **Regional Variations**: Support for different Indian language kinship terms
3. **Relationship Strength**: Calculate relationship "closeness" based on path length
4. **Multiple Paths**: Handle cases where multiple relationship paths exist
5. **Caching**: Cache frequently queried relationships for performance

### Requirements Validation

This implementation validates the following requirements:
- **Requirement 4.2**: Relationship path finding using Neo4j shortest path algorithm
- **Requirement 4.3**: Kinship mapping with Indian family relationship conventions
- **Requirement 4.6**: Natural language description generation

### Related Files

- `internal/repository/relationship_repo.go`: Repository layer for Neo4j queries
- `internal/models/relationship.go`: Relationship data model
- `internal/handler/relationship_handler.go`: HTTP handlers (to be implemented)
