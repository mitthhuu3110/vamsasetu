# Task 4.4: Relationship Repository Property Tests

## Overview

This document summarizes the property-based tests implemented for the relationship repository, covering Properties 10-13 from the design document.

## Implemented Property Tests

### Property 10: Relationship Creation and Retrieval
**Validates: Requirements 2.4**

**Property Statement**: For any two existing members and a valid relationship type (SPOUSE_OF, PARENT_OF, SIBLING_OF), creating a relationship should result in that relationship being stored in Neo4j and retrievable.

**Test Implementation**:
- Generates random pairs of members with valid attributes
- Generates random valid relationship types
- Creates both members in Neo4j
- Creates a relationship between them
- Retrieves all relationships and verifies the created relationship exists
- Handles bidirectional relationships (SPOUSE_OF, SIBLING_OF) by checking both directions
- Handles directed relationships (PARENT_OF) by checking exact direction
- Cleans up test data after verification

**Test Configuration**:
- Minimum successful tests: 100
- Uses gopter property-based testing framework

### Property 11: Relationship Deletion
**Validates: Requirements 2.5**

**Property Statement**: For any existing relationship, deleting that relationship should remove it from Neo4j such that it is no longer retrievable.

**Test Implementation**:
- Generates random pairs of members and relationship types
- Creates both members and a relationship between them
- Verifies the relationship exists before deletion
- Deletes the relationship
- Verifies the relationship no longer exists after deletion
- Handles both bidirectional and directed relationships correctly
- Cleans up test data

**Test Configuration**:
- Minimum successful tests: 100
- Verifies both presence before deletion and absence after deletion

### Property 12: Relationship Semantic Validation
**Validates: Requirements 2.6**

**Property Statement**: For any member, attempting to create a PARENT_OF relationship from that member to itself should be rejected with a validation error.

**Test Implementation**:
- Generates random members
- Attempts to create a self-referential PARENT_OF relationship
- Verifies that validation fails (either at model level or repository level)
- Tests the semantic constraint that prevents illogical relationships
- Cleans up test data

**Test Configuration**:
- Minimum successful tests: 100
- Validates that self-referential relationships are rejected

### Property 13: Soft Delete Enforcement for Connected Members
**Validates: Requirements 2.7**

**Property Statement**: For any member with at least one relationship, attempting to delete that member should result in a soft delete (isDeleted=true) rather than removal from the database.

**Test Implementation**:
- Generates random pairs of members and creates a relationship
- Performs soft delete on a member that has a relationship
- Verifies the member is not retrievable via GetByID (which filters deleted members)
- Directly queries Neo4j to verify the member still exists in the database
- Verifies isDeleted flag is set to true
- Verifies all other member data is preserved (ID, name, gender)
- Cleans up test data including hard deletion of test members

**Test Configuration**:
- Minimum successful tests: 100
- Uses direct Neo4j queries to bypass soft delete filters and verify data preservation

## Generator Functions

### genValidRelationshipType()
Generates valid relationship types from the three allowed values:
- SPOUSE_OF
- PARENT_OF
- SIBLING_OF

## Test Setup

The tests use a shared setup function `setupTestRelationshipRepo()` that:
1. Creates a member repository with Neo4j connection
2. Creates a relationship repository using the same client
3. Returns both repositories, context, and cleanup function
4. Reuses the existing `setupTestRepo()` function from member tests

## Running the Tests

### Prerequisites
- Neo4j database running on bolt://localhost:7687
- Neo4j credentials: neo4j/vamsasetu123
- Go 1.21+ installed
- Dependencies installed (`go mod download`)

### Commands

Run all property tests:
```bash
cd backend
go test -v -run "TestProperty10_RelationshipCreationAndRetrieval|TestProperty11_RelationshipDeletion|TestProperty12_RelationshipSemanticValidation|TestProperty13_SoftDeleteEnforcementForConnectedMembers" ./internal/repository/
```

Run a specific property test:
```bash
cd backend
go test -v -run TestProperty10_RelationshipCreationAndRetrieval ./internal/repository/
```

Run with coverage:
```bash
cd backend
go test -v -coverprofile=coverage.out -run "TestProperty1[0-3]" ./internal/repository/
go tool cover -html=coverage.out
```

### Using Docker

If Go is not installed locally, you can run tests in the Docker container:
```bash
docker-compose up -d neo4j postgres redis
docker-compose run --rm backend go test -v ./internal/repository/
```

## Test Patterns

### Bidirectional Relationship Handling
For SPOUSE_OF and SIBLING_OF relationships, the tests check both directions:
```go
if rel.IsBidirectional() {
    if (r.FromID == member1.ID && r.ToID == member2.ID) ||
       (r.FromID == member2.ID && r.ToID == member1.ID) {
        found = true
    }
}
```

### Directed Relationship Handling
For PARENT_OF relationships, the tests check exact direction:
```go
if r.FromID == member1.ID && r.ToID == member2.ID {
    found = true
}
```

### Soft Delete Verification
Tests verify soft delete by:
1. Confirming member is not in GetByID results (filtered)
2. Directly querying Neo4j to confirm member still exists
3. Verifying isDeleted=true
4. Verifying all other data is preserved

## Integration with Existing Tests

These property tests complement the existing unit tests in `relationship_repo_test.go`:
- Unit tests verify specific behaviors and edge cases
- Property tests verify universal properties across many random inputs
- Both test types are valuable and provide comprehensive coverage

## Notes

- All tests include proper cleanup to avoid test data pollution
- Tests use the "Test" prefix in member names for easy identification
- Tests handle both success and failure cases appropriately
- Tests log detailed error messages for debugging failures
- Property-based testing provides high confidence in correctness across diverse inputs
