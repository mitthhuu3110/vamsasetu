# Task 4.2: Property Tests for Member Repository

## Overview

This document describes the property-based tests implemented for the member repository, validating Requirements 2.1, 2.2, 2.3, and 11.1.

## Implemented Properties

### Property 7: Member Creation and Retrieval
**Validates: Requirements 2.1, 11.1**

**Property Statement:**
For any valid member data (name, date of birth, gender), creating a member should result in that member being stored in Neo4j and retrievable by ID.

**Test Function:** `TestProperty7_MemberCreationAndRetrieval`

**Test Strategy:**
- Generates 100+ random valid member data combinations
- Creates each member in Neo4j
- Retrieves the member by ID
- Verifies all fields match (ID, Name, DateOfBirth, Gender, IsDeleted)
- Cleans up by soft deleting the test member

**Generators Used:**
- `genValidName()`: Generates names 2-100 characters, prefixed with "Test" for cleanup
- `genValidDateOfBirth()`: Generates dates between 1900 and today
- `genValidGender()`: Generates "male", "female", or "other"

### Property 8: Member Update Persistence
**Validates: Requirements 2.2**

**Property Statement:**
For any existing member and any valid attribute changes, updating the member should persist those changes such that subsequent retrieval returns the updated values.

**Test Function:** `TestProperty8_MemberUpdatePersistence`

**Test Strategy:**
- Generates 100+ combinations of original and updated member data
- Creates a member with original data
- Updates the member's name, email, and phone
- Persists the update to Neo4j
- Retrieves the member and verifies updated fields match
- Verifies unchanged fields (ID, Gender) remain the same
- Cleans up by soft deleting the test member

**Generators Used:**
- `genValidName()`: For original and updated names
- `genValidDateOfBirth()`: For date of birth
- `genValidGender()`: For gender
- `genValidEmail()`: Generates email addresses (alphanumeric@test.com)
- `genValidPhone()`: Generates phone numbers (+91XXXXXXXXXX)

### Property 9: Soft Delete Preservation
**Validates: Requirements 2.3**

**Property Statement:**
For any member, deleting that member should set the isDeleted flag to true while preserving all other member data.

**Test Function:** `TestProperty9_SoftDeletePreservation`

**Test Strategy:**
- Generates 100+ member data combinations with all fields populated
- Creates a member in Neo4j
- Stores original values for comparison
- Performs soft delete
- Verifies member is NOT retrievable via GetByID (filters deleted members)
- Verifies member is NOT in GetAll results (filters deleted members)
- Directly queries Neo4j (bypassing isDeleted filter) to verify:
  - isDeleted flag is set to true
  - All other data is preserved (ID, Name, Gender, Email, Phone)
- Cleans up by hard deleting the test member

**Generators Used:**
- `genValidName()`: For member name
- `genValidDateOfBirth()`: For date of birth
- `genValidGender()`: For gender
- `genValidEmail()`: For email address
- `genValidPhone()`: For phone number

## Generator Functions

### genValidName()
Generates valid member names:
- Length: 2-100 characters
- Format: Alphanumeric strings
- Prefix: "Test " for easy cleanup
- Example: "Test AbCdEf"

### genValidDateOfBirth()
Generates valid dates of birth:
- Range: January 1, 1900 to current date
- Format: time.Time in UTC
- Implementation: Generates Unix timestamps and converts to time.Time

### genValidGender()
Generates valid gender values:
- Values: "male", "female", "other"
- Implementation: Uses gen.OneConstOf for uniform distribution

### genValidEmail()
Generates valid email addresses:
- Format: [alphanumeric]@test.com
- Length: 3-50 characters before @
- Example: "AbCdEf@test.com"

### genValidPhone()
Generates valid phone numbers:
- Format: +91[10-digit number]
- Range: +911000000000 to +919999999999
- Example: "+915432109876"

## Test Configuration

All property tests are configured with:
- **MinSuccessfulTests**: 100 (minimum iterations)
- **MaxSize**: 100 (maximum size for generated values)

This ensures comprehensive coverage across the input space.

## Running the Tests

### Prerequisites
1. Neo4j database running on `bolt://localhost:7687`
2. Neo4j credentials: username=`neo4j`, password=`vamsasetu123`
3. Go 1.21+ installed

### Using Docker Compose
```bash
# Start Neo4j database
docker-compose up -d neo4j

# Wait for Neo4j to be healthy
docker-compose ps

# Run property tests
cd backend
go test -v ./internal/repository -run "TestProperty"
```

### Using Makefile
```bash
cd backend

# Run all tests (includes property tests)
make test

# Run only property tests
go test -v ./internal/repository -run "TestProperty"
```

### Direct Go Command
```bash
cd backend
go test -v ./internal/repository -run "TestProperty7"
go test -v ./internal/repository -run "TestProperty8"
go test -v ./internal/repository -run "TestProperty9"
```

## Expected Output

### Successful Test Run
```
=== RUN   TestProperty7_MemberCreationAndRetrieval
+ For any valid member data, creating a member should result in that member being retrievable by ID: OK, passed 100 tests.
--- PASS: TestProperty7_MemberCreationAndRetrieval (X.XXs)

=== RUN   TestProperty8_MemberUpdatePersistence
+ For any existing member and valid attribute changes, updating should persist changes: OK, passed 100 tests.
--- PASS: TestProperty8_MemberUpdatePersistence (X.XXs)

=== RUN   TestProperty9_SoftDeletePreservation
+ For any member, soft delete should set isDeleted=true while preserving all data: OK, passed 100 tests.
--- PASS: TestProperty9_SoftDeletePreservation (X.XXs)

PASS
ok      vamsasetu/backend/internal/repository   X.XXXs
```

### Failed Test Example
If a property fails, gopter will provide a counterexample:
```
! For any valid member data, creating a member should result in that member being retrievable by ID: Falsified after 42 passed tests.
ARG_0: "Test XyZ"
ARG_1: 1990-05-15 00:00:00 +0000 UTC
ARG_2: "male"
```

## Test Cleanup

All property tests include cleanup logic:
- **Property 7 & 8**: Soft delete test members after verification
- **Property 9**: Hard delete test members after verification (since they're already soft deleted)

The `setupTestRepo` function also includes a cleanup function that removes all members with names starting with "Test".

## Integration with CI/CD

These tests can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run Property Tests
  run: |
    docker-compose up -d neo4j
    docker-compose exec -T neo4j cypher-shell -u neo4j -p vamsasetu123 'RETURN 1'
    cd backend && go test -v ./internal/repository -run "TestProperty"
```

## Troubleshooting

### Test Failures
1. **Database Connection Issues**: Ensure Neo4j is running and accessible
2. **Authentication Errors**: Verify Neo4j credentials match configuration
3. **Timeout Issues**: Increase test timeout or reduce MinSuccessfulTests
4. **Data Conflicts**: Ensure test cleanup is working properly

### Performance
- Each property test runs 100+ iterations
- Total execution time: ~10-30 seconds depending on database performance
- Can be parallelized using `go test -parallel N`

## File Location

**Test File**: `backend/internal/repository/member_repo_property_test.go`

**Related Files**:
- `backend/internal/repository/member_repo.go` (implementation)
- `backend/internal/repository/member_repo_test.go` (unit tests)
- `backend/internal/models/member.go` (model definition)

## Compliance

These property tests fulfill the requirements specified in:
- **Design Document**: Section "Correctness Properties" - Properties 7, 8, 9
- **Requirements Document**: Requirements 2.1, 2.2, 2.3, 11.1
- **Testing Strategy**: Minimum 100 iterations per property test

## Next Steps

After these property tests pass:
1. Proceed to Task 4.3: Implement relationship repository
2. Write property tests for relationship repository
3. Integrate with service layer tests
4. Add integration tests for complete workflows
