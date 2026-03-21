# Task 5.3: Event Repository Property-Based Tests

## Overview

This document summarizes the implementation of property-based tests for the Event Repository, validating the correctness properties defined in the design document.

## Implemented Properties

### Property 21: Event Creation and Retrieval
**Validates: Requirements 5.1, 11.2**

**Property Statement**: For any valid event data (title, date, type, member IDs), creating an event should result in that event being stored in PostgreSQL and retrievable by ID.

**Test Implementation**:
- Generates random valid event data (title, description, date, type, member IDs, creator)
- Creates event in repository
- Retrieves event by ID
- Verifies all fields match (ID, title, description, date, type, member IDs, creator)
- Runs 100 test cases with different random inputs

**Key Validations**:
- Event ID is assigned after creation
- All event attributes are persisted correctly
- Retrieved event matches created event exactly

### Property 22: Event Update Persistence
**Validates: Requirements 5.2**

**Property Statement**: For any existing event and any valid attribute changes, updating the event should persist those changes such that subsequent retrieval returns the updated values.

**Test Implementation**:
- Creates an initial event with original data
- Updates event attributes (title, description, date)
- Persists the update
- Retrieves and verifies updated fields
- Verifies unchanged fields remain the same (ID, type)
- Runs 100 test cases

**Key Validations**:
- Updated fields are persisted correctly
- Unchanged fields remain intact
- Event ID does not change after update

### Property 23: Event Deletion
**Validates: Requirements 5.3**

**Property Statement**: For any existing event, deleting that event should remove it from PostgreSQL such that it is no longer retrievable.

**Test Implementation**:
- Creates an event
- Verifies event exists before deletion
- Deletes the event
- Verifies event is no longer retrievable by ID
- Verifies event does not appear in GetAll results
- Runs 100 test cases

**Key Validations**:
- Deleted events cannot be retrieved by ID
- Deleted events do not appear in GetAll queries
- Delete operation completes without errors

### Property 24: Event Type Validity
**Validates: Requirements 5.4**

**Property Statement**: For any event in the system, that event's type must be one of: birthday, anniversary, ceremony, or custom.

**Test Implementation**:
- Generates events with valid types only
- Creates and retrieves events
- Verifies event type is one of the four valid types
- Runs 100 test cases

**Key Validations**:
- All events have valid types (birthday, anniversary, ceremony, custom)
- Invalid types are rejected by the generator
- Type validation is enforced at the database level

### Property 26: Event Filtering
**Validates: Requirements 5.7, 8.3, 8.4**

**Property Statement**: For any event filter criteria (type, member ID, date range), the returned events should match all specified filter criteria.

**Test Implementation**: Three sub-properties tested:

#### 26a: Filter by Type
- Creates multiple events with different types
- Filters by a specific type
- Verifies all returned events have the specified type
- Runs 50 test cases

#### 26b: Filter by Member
- Creates multiple events with different member IDs
- Filters by a specific member ID
- Verifies all returned events contain the specified member ID
- Runs 50 test cases

#### 26c: Filter by Date Range
- Creates multiple events with different dates
- Filters by a date range (start and end date)
- Verifies all returned events fall within the specified range
- Runs 50 test cases

**Key Validations**:
- Type filtering returns only events of the specified type
- Member filtering returns only events containing the specified member
- Date range filtering returns only events within the specified range

## Generator Functions

The following generator functions were implemented to create valid test data:

1. **genValidEventTitle()**: Generates event titles (3-100 characters, prefixed with "Test Event")
2. **genValidEventDescription()**: Generates event descriptions (up to 500 characters)
3. **genValidEventDate()**: Generates event dates (between now and 2 years from now)
4. **genValidEventType()**: Generates valid event types (birthday, anniversary, ceremony, custom)
5. **genValidMemberIDs()**: Generates lists of 1-5 member IDs
6. **genValidMemberID()**: Generates a single member ID
7. **genValidUserID()**: Generates user IDs (1-1000)
8. **genEventDataList()**: Generates lists of event data for filtering tests
9. **genEventDataListWithMember()**: Generates event lists with specific member IDs for member filtering tests

## Test Configuration

- **Framework**: gopter (Go property-based testing library)
- **Test Parameters**:
  - MinSuccessfulTests: 100 (50 for filtering tests)
  - MaxSize: 100 (50 for filtering tests)
- **Test Isolation**: Each test creates and cleans up its own data
- **Database**: PostgreSQL (via GORM)

## Test Execution

To run the property tests:

```bash
# Run all event repository property tests
go test -v -run TestProperty2[1-6] ./internal/repository/

# Run specific property test
go test -v -run TestProperty21 ./internal/repository/
go test -v -run TestProperty22 ./internal/repository/
go test -v -run TestProperty23 ./internal/repository/
go test -v -run TestProperty24 ./internal/repository/
go test -v -run TestProperty26 ./internal/repository/
```

## Requirements Coverage

The property tests validate the following requirements:

- **Requirement 5.1**: Event creation with required attributes
- **Requirement 5.2**: Event update persistence
- **Requirement 5.3**: Event deletion
- **Requirement 5.4**: Event type validity (birthday, anniversary, ceremony, custom)
- **Requirement 5.7**: Event filtering by type, member, and date range
- **Requirement 8.3**: Search and filtering by type
- **Requirement 8.4**: Search and filtering by member
- **Requirement 11.2**: Data persistence in PostgreSQL

## Notes

- All property tests follow the same pattern as member repository property tests
- Tests use the existing `setupEventTestRepo` function from `event_repo_test.go`
- Each test includes proper cleanup to avoid test data pollution
- Property tests complement the existing unit tests by validating universal properties across many random inputs
- The tests are designed to run independently and can be executed in any order

## Next Steps

After running the tests:
1. Verify all property tests pass
2. Review any failing test cases to identify edge cases
3. Update implementation if property violations are found
4. Consider adding additional properties if new edge cases are discovered

