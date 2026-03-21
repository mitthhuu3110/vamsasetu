# VamsaSetu Backend Integration Guide

## Quick Start

### 1. Environment Setup

Create a `.env` file in the project root:

```bash
# Database Configuration
POSTGRES_URL=postgresql://postgres:password@localhost:5432/vamsasetu?sslmode=disable
NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=password
REDIS_ADDR=localhost:6379

# Authentication
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Notification Services
SENDGRID_API_KEY=your-sendgrid-api-key
TWILIO_ACCOUNT_SID=your-twilio-account-sid
TWILIO_AUTH_TOKEN=your-twilio-auth-token
TWILIO_PHONE_NUMBER=+1234567890
TWILIO_WHATSAPP_NUMBER=whatsapp:+1234567890

# Optional
PORT=8080
ENV=development
FRONTEND_ORIGIN=http://localhost:3000
```

### 2. Start Services with Docker Compose

```bash
docker-compose up -d
```

This starts:
- PostgreSQL on port 5432
- Neo4j on ports 7474 (HTTP) and 7687 (Bolt)
- Redis on port 6379

### 3. Run the Backend Server

```bash
cd backend
go run cmd/server/main.go
```

Expected output:
```
VamsaSetu Backend Server - Starting...
Configuration loaded successfully
Connecting to PostgreSQL...
PostgreSQL connection established
Running GORM auto-migration...
✓ Migrated User table
✓ Migrated Event table
✓ Migrated Notification table
✓ Migrated AuditLog table
Database migrations completed successfully
Connecting to Neo4j...
Neo4j connection established
Connecting to Redis...
Redis connection established
All database connections established
Repositories initialized
Services initialized
Handlers initialized
Notification scheduler started
Server starting on port 8080
```

## API Endpoints

### Authentication

#### Register a New User
```bash
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "name": "John Doe",
  "role": "owner"
}

Response:
{
  "success": true,
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "role": "owner"
    }
  },
  "error": ""
}
```

#### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}

Response: Same as register
```

#### Get Profile
```bash
GET /api/auth/profile
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "role": "owner"
  },
  "error": ""
}
```

### Members

#### Create a Member
```bash
POST /api/members
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Rajesh Kumar",
  "dateOfBirth": "1970-05-15T00:00:00Z",
  "gender": "male",
  "email": "rajesh@example.com",
  "phone": "+919876543210",
  "avatarUrl": "https://example.com/avatar.jpg"
}

Response:
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Rajesh Kumar",
    "dateOfBirth": "1970-05-15T00:00:00Z",
    "gender": "male",
    "email": "rajesh@example.com",
    "phone": "+919876543210",
    "avatarUrl": "https://example.com/avatar.jpg",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z",
    "isDeleted": false
  },
  "error": ""
}
```

#### List Members
```bash
GET /api/members?page=1&limit=10&search=rajesh&gender=male
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": {
    "members": [...],
    "total": 5,
    "page": 1,
    "limit": 10
  },
  "error": ""
}
```

#### Get Member by ID
```bash
GET /api/members/:id
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Rajesh Kumar",
    ...
  },
  "error": ""
}
```

#### Update Member
```bash
PUT /api/members/:id
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Rajesh Kumar Updated",
  "email": "rajesh.new@example.com"
}

Response: Updated member object
```

#### Delete Member (Soft Delete)
```bash
DELETE /api/members/:id
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": {
    "message": "Member deleted successfully"
  },
  "error": ""
}
```

### Relationships

#### Create a Relationship
```bash
POST /api/relationships
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "type": "PARENT_OF",
  "fromId": "550e8400-e29b-41d4-a716-446655440000",
  "toId": "660e8400-e29b-41d4-a716-446655440001"
}

Response:
{
  "success": true,
  "data": {
    "type": "PARENT_OF",
    "fromId": "550e8400-e29b-41d4-a716-446655440000",
    "toId": "660e8400-e29b-41d4-a716-446655440001",
    "createdAt": "2024-01-15T10:35:00Z"
  },
  "error": ""
}
```

Valid relationship types:
- `SPOUSE_OF` - Bidirectional relationship between married partners
- `PARENT_OF` - Directed relationship from parent to child
- `SIBLING_OF` - Bidirectional relationship between siblings

#### List All Relationships
```bash
GET /api/relationships
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": [...],
  "error": ""
}
```

#### Get Member Relationships
```bash
GET /api/members/:id/relationships
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": {
    "relationships": [...]
  },
  "error": ""
}
```

#### Delete Relationship
```bash
DELETE /api/relationships/:id?fromId=uuid1&toId=uuid2&type=PARENT_OF
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": {
    "message": "Relationship deleted successfully"
  },
  "error": ""
}
```

### Events

#### Create an Event
```bash
POST /api/events
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "title": "Rajesh's Birthday",
  "description": "Birthday celebration",
  "eventDate": "2024-05-15T00:00:00Z",
  "eventType": "birthday",
  "memberIds": ["550e8400-e29b-41d4-a716-446655440000"]
}

Response:
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Rajesh's Birthday",
    "description": "Birthday celebration",
    "eventDate": "2024-05-15T00:00:00Z",
    "eventType": "birthday",
    "memberIds": ["550e8400-e29b-41d4-a716-446655440000"],
    "createdBy": 1,
    "createdAt": "2024-01-15T10:40:00Z",
    "updatedAt": "2024-01-15T10:40:00Z"
  },
  "error": ""
}
```

Valid event types:
- `birthday`
- `anniversary`
- `ceremony`
- `custom`

#### List Events
```bash
GET /api/events?page=1&limit=10&type=birthday&member=uuid
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": {
    "events": [...],
    "total": 5,
    "page": 1,
    "limit": 10
  },
  "error": ""
}
```

#### Get Upcoming Events
```bash
GET /api/events/upcoming?days=30
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": [...],
  "error": ""
}
```

#### Update Event
```bash
PUT /api/events/:id
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "title": "Updated Title",
  "description": "Updated description"
}

Response: Updated event object
```

#### Delete Event
```bash
DELETE /api/events/:id
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": {
    "message": "Event deleted successfully"
  },
  "error": ""
}
```

### Family Tree

#### Get Family Tree
```bash
GET /api/family/tree
Authorization: Bearer <access_token>

Response:
{
  "success": true,
  "data": {
    "nodes": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "type": "memberNode",
        "position": { "x": 0, "y": 0 },
        "data": {
          "id": "550e8400-e29b-41d4-a716-446655440000",
          "name": "Rajesh Kumar",
          "avatarUrl": "https://example.com/avatar.jpg",
          "relationBadge": "Father",
          "hasUpcomingEvent": true,
          "gender": "male"
        }
      }
    ],
    "edges": [
      {
        "id": "edge1",
        "source": "550e8400-e29b-41d4-a716-446655440000",
        "target": "660e8400-e29b-41d4-a716-446655440001",
        "type": "bezier",
        "animated": false,
        "style": {
          "stroke": "#0D9488",
          "strokeWidth": 2
        }
      }
    ]
  },
  "error": ""
}
```

### Health Check

#### Check System Health
```bash
GET /health

Response:
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

## Error Handling

All errors follow the consistent format:

```json
{
  "success": false,
  "data": null,
  "error": "Descriptive error message"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation error)
- `401` - Unauthorized (missing or invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

## Authentication Flow

1. **Register** or **Login** to get access and refresh tokens
2. **Store tokens** securely (localStorage or httpOnly cookies)
3. **Include access token** in Authorization header for all protected endpoints:
   ```
   Authorization: Bearer <access_token>
   ```
4. **Refresh token** when access token expires:
   ```bash
   POST /api/auth/refresh
   Content-Type: application/json
   
   {
     "refreshToken": "<refresh_token>"
   }
   ```

## Role-Based Access Control

### Roles
- **owner** - Full access to all operations
- **viewer** - Read-only access
- **admin** - System-wide administrative privileges

### Permissions
- **Public endpoints**: `/api/auth/register`, `/api/auth/login`, `/health`
- **Authenticated endpoints**: All `/api/*` endpoints require valid JWT
- **Owner/Admin only**: POST, PUT, DELETE operations on members, relationships, events

## Notification System

The notification scheduler runs automatically in the background:

- **Frequency**: Every hour
- **Worker Pool**: Max 10 concurrent workers
- **Retry Logic**: Up to 3 retries with exponential backoff (5min, 15min, 45min)
- **Channels**: WhatsApp, SMS, Email

Notifications are automatically created when events are added and scheduled based on user preferences.

## Caching Strategy

Redis caching is used for:

- **Family Tree**: TTL 5 minutes, key `family_tree:{userId}`
- **Member Details**: TTL 10 minutes, key `member:{memberId}`
- **Search Results**: TTL 2 minutes, key `search:members:{query}`
- **Upcoming Events**: TTL 5 minutes, key `events:upcoming:{userId}`

Cache is automatically invalidated on data modifications.

## Graceful Shutdown

The server handles shutdown signals (SIGINT, SIGTERM) gracefully:

1. Stop accepting new requests
2. Stop notification scheduler
3. Complete in-flight requests (10-second timeout)
4. Close database connections
5. Exit

To shutdown:
```bash
# Press Ctrl+C or send SIGTERM
kill -TERM <pid>
```

## Troubleshooting

### Database Connection Issues

**PostgreSQL:**
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check logs
docker logs vamsasetu-postgres
```

**Neo4j:**
```bash
# Check if Neo4j is running
docker ps | grep neo4j

# Access Neo4j browser
open http://localhost:7474
```

**Redis:**
```bash
# Check if Redis is running
docker ps | grep redis

# Test connection
redis-cli ping
```

### Common Errors

**"Missing required environment variables"**
- Ensure all required variables are set in `.env`
- Check `.env.example` for reference

**"Failed to connect to database"**
- Verify database services are running
- Check connection strings in `.env`
- Ensure ports are not blocked by firewall

**"Invalid or expired token"**
- Token may have expired (15 minutes for access token)
- Use refresh token to get new access token
- Re-login if refresh token expired (7 days)

## Development Tips

### Hot Reload
Use `air` for hot reload during development:
```bash
go install github.com/cosmtrek/air@latest
cd backend
air
```

### Database Migrations
Migrations run automatically on startup using GORM AutoMigrate.

### Logging
All requests are logged with:
- HTTP method
- Path
- Status code
- Duration

Example log:
```
[GET] /api/members - 200 - 45.2ms
[POST] /api/auth/login - 200 - 123.5ms
```

## Next Steps

1. **Frontend Integration**: Connect React frontend to these APIs
2. **WebSocket**: Implement real-time updates (Task 9.1-9.3)
3. **Audit Logging**: Add comprehensive audit trails (Task 10.10)
4. **Testing**: Write integration tests for all endpoints
5. **Production**: Deploy with proper security, monitoring, and backups
