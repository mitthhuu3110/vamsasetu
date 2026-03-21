# Task 9.3: WebSocket Broadcast Integration

## Overview

This task integrates WebSocket broadcasting into the member, relationship, and event services to enable real-time updates for all connected clients.

## Changes Made

### 1. Service Layer Updates

#### Member Service (`backend/internal/service/member_service.go`)
- Added `WebSocketHub` interface definition for broadcasting
- Added `hub` field to `MemberService` struct
- Updated `NewMemberService` constructor to accept `hub` parameter
- Added broadcast calls in:
  - `Create()` - broadcasts `member_created` event
  - `Update()` - broadcasts `member_updated` event
  - `SoftDelete()` - broadcasts `member_deleted` event with member ID

#### Relationship Service (`backend/internal/service/relationship_service.go`)
- Added `hub` field to `RelationshipService` struct
- Updated `NewRelationshipService` constructor to accept `hub` parameter
- Added broadcast calls in:
  - `Create()` - broadcasts `relationship_created` event
  - `Delete()` - broadcasts `relationship_deleted` event with fromId, toId, and type

#### Event Service (`backend/internal/service/event_service.go`)
- Added `hub` field to `EventService` struct
- Updated `NewEventService` constructor to accept `hub` parameter
- Added broadcast calls in:
  - `Create()` - broadcasts `event_created` event
  - `Update()` - broadcasts `event_updated` event
  - `Delete()` - broadcasts `event_deleted` event with event ID

### 2. Main Application (`backend/cmd/server/main.go`)
- Initialize WebSocket hub using `handler.NewHub()`
- Start hub as goroutine with `go hub.Run()`
- Pass hub to service constructors:
  - `memberService := service.NewMemberService(memberRepo, redisClientInstance.Client, hub)`
  - `relationshipService := service.NewRelationshipService(relationshipRepo, hub)`
  - `eventService := service.NewEventService(eventRepo, redisClientInstance.Client, hub)`
- Register WebSocket endpoint: `app.Get("/ws", handler.HandleWebSocket(hub))`

### 3. Test Updates

Updated test files to pass `nil` for hub parameter (tests run in isolation):
- `backend/internal/handler/event_handler_test.go`
- `backend/internal/handler/relationship_handler_test.go`
- `backend/internal/service/relationship_service_property_test.go`

## WebSocket Message Format

All broadcast messages follow this format:
```json
{
  "type": "event_type",
  "data": { /* event-specific data */ },
  "timestamp": 1234567890
}
```

### Event Types

**Member Events:**
- `member_created` - Full member object
- `member_updated` - Full member object
- `member_deleted` - `{"id": "member-uuid"}`

**Relationship Events:**
- `relationship_created` - Full relationship object
- `relationship_deleted` - `{"fromId": "uuid1", "toId": "uuid2", "type": "PARENT_OF"}`

**Event Events:**
- `event_created` - Full event object
- `event_updated` - Full event object
- `event_deleted` - `{"id": 123}`

## Implementation Details

### Null Safety
All broadcast calls check if hub is nil before calling `BroadcastUpdate()`:
```go
if s.hub != nil {
    s.hub.BroadcastUpdate("event_type", data)
}
```

This allows services to work without WebSocket support (e.g., in tests or when hub is not initialized).

### Interface Design
The `WebSocketHub` interface is defined in the service package to avoid circular dependencies:
```go
type WebSocketHub interface {
    BroadcastUpdate(eventType string, data interface{})
}
```

The actual `Hub` struct in `handler` package implements this interface.

## Testing

All modified files compile without errors. Tests pass with nil hub parameter.

## Requirements Validated

- **Requirement 7.2**: Real-time updates via WebSocket
- **Requirement 7.3**: Broadcast member, relationship, and event changes

## Next Steps

Frontend integration:
1. Connect to WebSocket endpoint at `/ws`
2. Listen for broadcast messages
3. Update UI state based on event type
4. Refresh family tree visualization on member/relationship changes
5. Update event lists on event changes
