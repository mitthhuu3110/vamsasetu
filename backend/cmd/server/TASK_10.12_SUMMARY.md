# Task 10.12: Wire Everything Together in main.go

## Summary

Successfully implemented the complete application wiring in `main.go`, integrating all components of the VamsaSetu backend system.

## Implementation Details

### 1. Configuration and Database Initialization

**Configuration Loading:**
- Loads environment variables using `config.Load()`
- Validates all required configuration parameters
- Provides default values for optional parameters (PORT, ENV)

**Database Clients:**
- **PostgreSQL**: Initialized with GORM for relational data (users, events, notifications, audit logs)
- **Neo4j**: Initialized for graph data (members and relationships)
- **Redis**: Initialized for caching and session management

**Health Checks:**
- Each database client is verified with a health check before proceeding
- Proper error handling with cleanup on failure
- All connections are deferred for cleanup on shutdown

### 2. Repository Layer Initialization

Created instances of all repositories:
- `UserRepository` - PostgreSQL-based user data access
- `MemberRepository` - Neo4j-based member graph queries
- `RelationshipRepository` - Neo4j-based relationship graph queries
- `EventRepository` - PostgreSQL-based event data access
- `NotificationRepository` - PostgreSQL-based notification data access

### 3. Service Layer Initialization

Created instances of all services with proper dependencies:

**AuthService:**
- Dependencies: UserRepository, Redis Client
- Handles user registration, login, token refresh

**MemberService:**
- Dependencies: MemberRepository, Redis Client
- Handles member CRUD operations with caching

**RelationshipService:**
- Dependencies: RelationshipRepository, MemberRepository
- Handles relationship CRUD and path finding

**EventService:**
- Dependencies: EventRepository, Redis Client
- Handles event CRUD operations with caching

**NotificationService:**
- Dependencies: NotificationRepository, SendGrid API Key, Twilio credentials
- Handles notification dispatch via WhatsApp, SMS, and Email

**TreeBuilder:**
- Dependencies: MemberRepository, RelationshipRepository, EventRepository
- Builds family tree structure in React Flow format

### 4. Handler Layer Initialization

Created instances of all HTTP handlers:
- `AuthHandler` - Authentication endpoints
- `MemberHandler` - Member CRUD endpoints
- `RelationshipHandler` - Relationship CRUD endpoints
- `EventHandler` - Event CRUD endpoints
- `FamilyHandler` - Family tree visualization endpoint

### 5. Background Services

**Notification Scheduler:**
- Initialized with NotificationRepository and NotificationService
- Started as a goroutine
- Runs hourly to process pending notifications
- Uses worker pool with max 10 concurrent workers
- Implements retry logic with exponential backoff

### 6. Fiber Application Setup

**Middleware Order (as specified in requirements):**
1. **Logger Middleware** - Logs all HTTP requests with method, path, status, duration
2. **CORS Middleware** - Allows frontend origin with credentials
3. **Error Handler** - Global error handling with consistent response format

**Route Registration:**
- Auth routes: `/api/auth/*` (public and protected)
- Member routes: `/api/members/*` (protected, owner/admin for mutations)
- Relationship routes: `/api/relationships/*` (protected, owner/admin for mutations)
- Event routes: `/api/events/*` (protected, owner/admin for mutations)
- Family tree routes: `/api/family/*` (protected)

### 7. Health Check Endpoint

**Endpoint:** `GET /health`

**Response Format:**
```json
{
  "status": "healthy" | "degraded",
  "services": {
    "postgres": "healthy" | "unhealthy: error",
    "neo4j": "healthy" | "unhealthy: error",
    "redis": "healthy" | "unhealthy: error"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Health Check Logic:**
- Checks PostgreSQL, Neo4j, and Redis connections
- Returns "healthy" if all services are up
- Returns "degraded" if any service is down
- Provides detailed status for each service

### 8. Graceful Shutdown

**Signal Handling:**
- Listens for SIGINT and SIGTERM signals
- Initiates graceful shutdown on signal receipt

**Shutdown Sequence:**
1. Stop notification scheduler
2. Shutdown Fiber app with 10-second timeout
3. Close database connections (deferred)
4. Log shutdown completion

### 9. Server Configuration

**Port Configuration:**
- Reads from `PORT` environment variable
- Defaults to "8080" if not specified
- Server listens on `0.0.0.0:<port>`

**Frontend Origin:**
- Reads from `FRONTEND_ORIGIN` environment variable
- Defaults to "http://localhost:3000"
- Used for CORS configuration

## Architecture Flow

```
main()
  ├─ Load Configuration
  ├─ Initialize Databases (PostgreSQL, Neo4j, Redis)
  ├─ Run Migrations (GORM AutoMigrate)
  ├─ Initialize Repositories
  ├─ Initialize Services
  ├─ Initialize Handlers
  ├─ Start Notification Scheduler (goroutine)
  ├─ Create Fiber App
  ├─ Apply Middleware (Logger → CORS)
  ├─ Register Routes
  ├─ Add Health Check Endpoint
  ├─ Start Server (goroutine)
  └─ Wait for Shutdown Signal
      ├─ Stop Scheduler
      ├─ Shutdown Fiber
      └─ Close Database Connections
```

## Middleware Order

As specified in the design document:

1. **Logger** - Logs all incoming requests
2. **CORS** - Handles cross-origin requests
3. **Routes** - Route-specific handlers
4. **Auth** - Applied per-route in handlers (not globally)
5. **Detailed Logger** - Could be added after auth for user-specific logging

## Key Features

### Dependency Injection
- All components use constructor injection
- Clear dependency graph
- Easy to test and mock

### Error Handling
- Comprehensive error handling at each initialization step
- Proper cleanup on failure
- Descriptive error messages

### Logging
- Startup progress logging
- Database connection status
- Service initialization status
- Graceful shutdown logging

### Resource Management
- Deferred cleanup for database connections
- Graceful shutdown with timeout
- Proper signal handling

## Environment Variables Required

The application requires the following environment variables:

**Database Configuration:**
- `POSTGRES_URL` - PostgreSQL connection string
- `NEO4J_URI` - Neo4j connection URI
- `NEO4J_USERNAME` - Neo4j username
- `NEO4J_PASSWORD` - Neo4j password
- `REDIS_ADDR` - Redis address (host:port)

**Authentication:**
- `JWT_SECRET` - Secret key for JWT signing

**Notification Services:**
- `SENDGRID_API_KEY` - SendGrid API key for email
- `TWILIO_ACCOUNT_SID` - Twilio account SID
- `TWILIO_AUTH_TOKEN` - Twilio auth token
- `TWILIO_PHONE_NUMBER` - Twilio phone number for SMS
- `TWILIO_WHATSAPP_NUMBER` - Twilio WhatsApp number

**Optional:**
- `PORT` - Server port (default: 8080)
- `ENV` - Environment (default: development)
- `FRONTEND_ORIGIN` - Frontend URL for CORS (default: http://localhost:3000)

## Testing the Implementation

### 1. Start the Server

```bash
cd backend
go run cmd/server/main.go
```

### 2. Check Health Endpoint

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "services": {
    "postgres": "healthy",
    "neo4j": "healthy",
    "redis": "healthy"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 3. Test Authentication

```bash
# Register a new user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test@1234",
    "name": "Test User",
    "role": "owner"
  }'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test@1234"
  }'
```

### 4. Test Protected Endpoints

```bash
# Get profile (requires JWT token)
curl http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer <access_token>"

# List members (requires JWT token)
curl http://localhost:8080/api/members \
  -H "Authorization: Bearer <access_token>"

# Get family tree (requires JWT token)
curl http://localhost:8080/api/family/tree \
  -H "Authorization: Bearer <access_token>"
```

## Notes on WebSocket Implementation

**Current Status:**
- WebSocket hub (Task 9.1-9.3) is not yet implemented
- The main.go is ready to integrate WebSocket when implemented
- Placeholder for future WebSocket integration:

```go
// TODO: Initialize WebSocket hub when Task 9.1-9.3 is completed
// wsHub := handler.NewWebSocketHub()
// wsHub.Start()
// app.Get("/ws", wsHub.HandleWebSocket)
```

**Integration Steps (when WebSocket is ready):**
1. Create WebSocket hub instance
2. Start hub as goroutine
3. Register WebSocket endpoint
4. Pass hub to services for broadcasting updates
5. Stop hub during graceful shutdown

## Compliance with Requirements

This implementation satisfies the following requirements:

- **Requirement 12.1**: Configuration loaded from environment variables
- **Requirement 16.3**: All services properly orchestrated
- **Requirement 11.1**: Neo4j data persistence
- **Requirement 11.2**: PostgreSQL data persistence
- **Requirement 11.4**: Transaction management in services
- **Requirement 13.1-13.4**: Consistent API response format
- **Requirement 14.4**: Error handling middleware
- **Requirement 6.2**: Notification scheduler running
- **Requirement 1.3**: JWT authentication middleware
- **Requirement 1.5-1.6**: Role-based authorization

## Next Steps

1. **Test the complete system** with Docker Compose
2. **Implement WebSocket hub** (Task 9.1-9.3)
3. **Add audit logging** (Task 10.10)
4. **Implement transaction rollback** (Task 10.8)
5. **Frontend integration** (Tasks 12-20)

## Conclusion

Task 10.12 is complete. The VamsaSetu backend is fully wired and ready to serve API requests. All components are properly initialized, middleware is applied in the correct order, and graceful shutdown is implemented.
