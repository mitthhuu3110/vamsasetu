# Design Document: VamsaSetu

## Overview

VamsaSetu is a full-stack family tree and event management system that combines graph database technology with modern web frameworks to provide an intuitive, culturally-appropriate platform for Indian families. The system enables users to visualize complex family relationships, track important events, and receive automated notifications through multiple channels.

### Design Goals

- **Scalability**: Support family trees with hundreds of members and complex multi-generational relationships
- **Performance**: Sub-second response times for tree visualization and relationship queries
- **Real-time Collaboration**: Multiple users can view and edit family trees simultaneously with live updates
- **Cultural Authenticity**: Indian kinship terminology and visual design elements
- **Reliability**: Robust error handling, data persistence, and notification delivery with retry logic
- **Developer Experience**: Clear separation of concerns, consistent patterns, comprehensive documentation

### Technology Stack

**Frontend**:
- React 18 with TypeScript for type safety
- React Flow for interactive graph visualization
- Tailwind CSS for utility-first styling
- React Query for server state management
- Zustand for client state management
- Framer Motion for animations
- Axios for HTTP requests

**Backend**:
- Go 1.21+ with Fiber v2 framework
- GORM for PostgreSQL ORM
- Neo4j Go Driver for graph database access
- Go-Redis for caching
- Gorilla WebSocket for real-time updates
- JWT-Go for authentication

**Databases**:
- Neo4j 5.x for graph data (members and relationships)
- PostgreSQL 15 for relational data (users, events, notifications, audit logs)
- Redis 7.x for caching and session management

**Infrastructure**:
- Docker and Docker Compose for containerization
- Twilio for SMS and WhatsApp messaging
- SendGrid for email delivery


## Architecture

### System Architecture Diagram

```mermaid
graph TB
    subgraph "Client Layer"
        Browser[Web Browser]
        Mobile[Mobile Browser]
    end
    
    subgraph "Frontend Layer"
        React[React App<br/>TypeScript + Tailwind]
        ReactFlow[React Flow<br/>Tree Visualization]
        ReactQuery[React Query<br/>Server State]
        Zustand[Zustand<br/>Client State]
    end
    
    subgraph "Backend Layer"
        API[Fiber API Server<br/>Go]
        WS[WebSocket Hub<br/>Goroutines]
        Auth[JWT Auth<br/>Middleware]
        NotifScheduler[Notification Scheduler<br/>Goroutine]
    end
    
    subgraph "Service Layer"
        MemberSvc[Member Service]
        RelationSvc[Relationship Service]
        EventSvc[Event Service]
        NotifSvc[Notification Service]
        CacheSvc[Cache Service]
    end
    
    subgraph "Data Layer"
        Neo4j[(Neo4j<br/>Graph DB)]
        Postgres[(PostgreSQL<br/>Relational DB)]
        Redis[(Redis<br/>Cache)]
    end
    
    subgraph "External Services"
        Twilio[Twilio<br/>SMS/WhatsApp]
        SendGrid[SendGrid<br/>Email]
    end
    
    Browser --> React
    Mobile --> React
    React --> ReactFlow
    React --> ReactQuery
    React --> Zustand
    ReactQuery --> API
    React --> WS
    
    API --> Auth
    Auth --> MemberSvc
    Auth --> RelationSvc
    Auth --> EventSvc
    Auth --> NotifSvc
    
    MemberSvc --> Neo4j
    MemberSvc --> CacheSvc
    RelationSvc --> Neo4j
    RelationSvc --> CacheSvc
    EventSvc --> Postgres
    NotifSvc --> Postgres
    CacheSvc --> Redis
    
    NotifScheduler --> NotifSvc
    NotifSvc --> Twilio
    NotifSvc --> SendGrid
    
    WS --> MemberSvc
    WS --> RelationSvc
```

### Multi-Tier Architecture

VamsaSetu follows a clean multi-tier architecture with clear separation of concerns:

**Presentation Tier (Frontend)**:
- React components organized by feature (auth, family, events, notifications)
- Custom hooks for business logic reuse
- React Flow for tree canvas rendering
- Responsive design with mobile-first approach

**API Tier (Backend)**:
- RESTful endpoints following resource-based routing
- WebSocket endpoint for real-time updates
- JWT middleware for authentication
- CORS middleware for cross-origin requests
- Request validation middleware
- Error handling middleware

**Service Tier**:
- Business logic encapsulation
- Transaction management
- Cache invalidation logic
- Relationship path finding algorithms
- Notification dispatch logic

**Data Tier**:
- Neo4j for graph queries (Cypher)
- PostgreSQL for relational queries (SQL via GORM)
- Redis for caching (key-value operations)


## Components and Interfaces

### Backend Module Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Environment configuration
│   ├── middleware/
│   │   ├── auth.go                # JWT authentication
│   │   ├── cors.go                # CORS configuration
│   │   └── logger.go              # Request logging
│   ├── models/
│   │   ├── user.go                # User entity (PostgreSQL)
│   │   ├── event.go               # Event entity (PostgreSQL)
│   │   ├── notification.go        # Notification entity (PostgreSQL)
│   │   ├── member.go              # Member node (Neo4j)
│   │   └── relationship.go        # Relationship edge (Neo4j)
│   ├── repository/
│   │   ├── user_repo.go           # User data access
│   │   ├── event_repo.go          # Event data access
│   │   ├── notification_repo.go   # Notification data access
│   │   ├── member_repo.go         # Member graph queries
│   │   └── relationship_repo.go   # Relationship graph queries
│   ├── service/
│   │   ├── auth_service.go        # Authentication logic
│   │   ├── member_service.go      # Member business logic
│   │   ├── relationship_service.go # Relationship engine
│   │   ├── event_service.go       # Event management
│   │   ├── notification_service.go # Notification dispatch
│   │   └── cache_service.go       # Redis caching
│   ├── handler/
│   │   ├── auth_handler.go        # Auth endpoints
│   │   ├── member_handler.go      # Member endpoints
│   │   ├── relationship_handler.go # Relationship endpoints
│   │   ├── event_handler.go       # Event endpoints
│   │   └── websocket_handler.go   # WebSocket hub
│   ├── scheduler/
│   │   └── notification_scheduler.go # Background notification job
│   └── utils/
│       ├── jwt.go                 # JWT utilities
│       ├── validator.go           # Input validation
│       └── response.go            # API response helpers
└── pkg/
    ├── neo4j/
    │   └── client.go              # Neo4j connection
    ├── postgres/
    │   └── client.go              # PostgreSQL connection
    └── redis/
        └── client.go              # Redis connection
```

### Frontend Component Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── auth/
│   │   │   ├── LoginForm.tsx
│   │   │   └── RegisterForm.tsx
│   │   ├── family/
│   │   │   ├── TreeCanvas.tsx          # React Flow visualization
│   │   │   ├── MemberNode.tsx          # Custom node component
│   │   │   ├── RelationshipEdge.tsx    # Custom edge component
│   │   │   ├── MemberDetailsPanel.tsx
│   │   │   ├── AddMemberModal.tsx
│   │   │   └── AddRelationshipModal.tsx
│   │   ├── events/
│   │   │   ├── EventCalendar.tsx
│   │   │   ├── EventList.tsx
│   │   │   ├── EventCard.tsx
│   │   │   └── AddEventModal.tsx
│   │   ├── common/
│   │   │   ├── Navbar.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   ├── LoadingSpinner.tsx
│   │   │   ├── ErrorBoundary.tsx
│   │   │   └── EmptyState.tsx
│   │   └── ui/
│   │       ├── Button.tsx
│   │       ├── Input.tsx
│   │       ├── Modal.tsx
│   │       └── Card.tsx
│   ├── pages/
│   │   ├── LoginPage.tsx
│   │   ├── RegisterPage.tsx
│   │   ├── DashboardPage.tsx
│   │   ├── FamilyTreePage.tsx
│   │   └── EventsPage.tsx
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── useMembers.ts
│   │   ├── useRelationships.ts
│   │   ├── useEvents.ts
│   │   └── useWebSocket.ts
│   ├── services/
│   │   ├── api.ts                      # Axios instance
│   │   ├── authService.ts
│   │   ├── memberService.ts
│   │   ├── relationshipService.ts
│   │   └── eventService.ts
│   ├── stores/
│   │   ├── authStore.ts                # Zustand auth state
│   │   └── uiStore.ts                  # Zustand UI state
│   ├── types/
│   │   ├── user.ts
│   │   ├── member.ts
│   │   ├── relationship.ts
│   │   └── event.ts
│   ├── utils/
│   │   ├── constants.ts
│   │   ├── validators.ts
│   │   └── formatters.ts
│   └── App.tsx
└── package.json
```


### Key Interfaces

#### API Response Format

All API responses follow a consistent structure:

```typescript
interface APIResponse<T> {
  success: boolean;
  data: T | null;
  error: string;
}
```

#### Member Node (Neo4j)

```go
type Member struct {
    ID          string    `json:"id"`           // UUID
    Name        string    `json:"name"`
    DateOfBirth time.Time `json:"dateOfBirth"`
    Gender      string    `json:"gender"`       // "male", "female", "other"
    Email       string    `json:"email"`
    Phone       string    `json:"phone"`
    AvatarURL   string    `json:"avatarUrl"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    IsDeleted   bool      `json:"isDeleted"`    // Soft delete flag
}
```

#### Relationship Edge (Neo4j)

```go
type Relationship struct {
    Type      string    `json:"type"`         // "SPOUSE_OF", "PARENT_OF", "SIBLING_OF"
    FromID    string    `json:"fromId"`
    ToID      string    `json:"toId"`
    CreatedAt time.Time `json:"createdAt"`
}
```

#### User Entity (PostgreSQL)

```go
type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Email        string    `gorm:"unique;not null" json:"email"`
    PasswordHash string    `gorm:"not null" json:"-"`
    Name         string    `gorm:"not null" json:"name"`
    Role         string    `gorm:"not null" json:"role"` // "owner", "viewer", "admin"
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}
```

#### Event Entity (PostgreSQL)

```go
type Event struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `gorm:"not null" json:"title"`
    Description string    `json:"description"`
    EventDate   time.Time `gorm:"not null" json:"eventDate"`
    EventType   string    `gorm:"not null" json:"eventType"` // "birthday", "anniversary", "ceremony", "custom"
    MemberIDs   []string  `gorm:"type:text[]" json:"memberIds"` // Array of member UUIDs
    CreatedBy   uint      `gorm:"not null" json:"createdBy"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

#### Notification Entity (PostgreSQL)

```go
type Notification struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    EventID     uint      `gorm:"not null" json:"eventId"`
    UserID      uint      `gorm:"not null" json:"userId"`
    Channel     string    `gorm:"not null" json:"channel"` // "whatsapp", "sms", "email"
    ScheduledAt time.Time `gorm:"not null" json:"scheduledAt"`
    SentAt      *time.Time `json:"sentAt"`
    Status      string    `gorm:"not null" json:"status"` // "pending", "sent", "failed"
    RetryCount  int       `gorm:"default:0" json:"retryCount"`
    ErrorMsg    string    `json:"errorMsg"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

#### React Flow Node Data

```typescript
interface MemberNodeData {
  id: string;
  name: string;
  avatarUrl: string;
  relationBadge: string;      // e.g., "Father", "Mother", "Sibling"
  hasUpcomingEvent: boolean;
  gender: string;
}

interface ReactFlowNode {
  id: string;
  type: 'memberNode';
  position: { x: number; y: number };
  data: MemberNodeData;
}
```

#### React Flow Edge Data

```typescript
interface ReactFlowEdge {
  id: string;
  source: string;
  target: string;
  type: 'bezier';
  animated: boolean;
  style: {
    stroke: string;           // Color based on relationship type
    strokeWidth: number;
  };
  label?: string;
}
```


## Data Models

### Neo4j Graph Schema

The graph database stores the family tree structure with members as nodes and relationships as edges.

**Member Node Labels**: `(:Member)`

**Member Node Properties**:
- `id`: UUID (indexed, unique)
- `name`: String
- `dateOfBirth`: DateTime
- `gender`: String (enum: "male", "female", "other")
- `email`: String
- `phone`: String
- `avatarUrl`: String
- `createdAt`: DateTime
- `updatedAt`: DateTime
- `isDeleted`: Boolean (for soft delete)

**Relationship Types**:
- `SPOUSE_OF`: Bidirectional relationship between married partners
- `PARENT_OF`: Directed relationship from parent to child
- `SIBLING_OF`: Bidirectional relationship between siblings

**Relationship Properties**:
- `createdAt`: DateTime

**Cypher Indexes**:
```cypher
CREATE INDEX member_id_index FOR (m:Member) ON (m.id);
CREATE INDEX member_name_index FOR (m:Member) ON (m.name);
```

**Example Graph Structure**:
```cypher
// Create members
CREATE (father:Member {id: 'uuid-1', name: 'Rajesh', gender: 'male'})
CREATE (mother:Member {id: 'uuid-2', name: 'Lakshmi', gender: 'female'})
CREATE (child:Member {id: 'uuid-3', name: 'Arjun', gender: 'male'})

// Create relationships
CREATE (father)-[:SPOUSE_OF]->(mother)
CREATE (father)-[:PARENT_OF]->(child)
CREATE (mother)-[:PARENT_OF]->(child)
```

### PostgreSQL Schema

**users table**:
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('owner', 'viewer', 'admin')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

**events table**:
```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMP NOT NULL,
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN ('birthday', 'anniversary', 'ceremony', 'custom')),
    member_ids TEXT[] NOT NULL,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_events_created_by ON events(created_by);
```

**notifications table**:
```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    channel VARCHAR(50) NOT NULL CHECK (channel IN ('whatsapp', 'sms', 'email')),
    scheduled_at TIMESTAMP NOT NULL,
    sent_at TIMESTAMP,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'sent', 'failed')),
    retry_count INTEGER DEFAULT 0,
    error_msg TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_scheduled ON notifications(scheduled_at, status);
CREATE INDEX idx_notifications_event ON notifications(event_id);
```

**audit_logs table**:
```sql
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    details JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);
```

### Redis Cache Keys

**Cache Key Patterns**:
- `family_tree:{userId}` - Complete family tree structure (TTL: 5 minutes)
- `member:{memberId}` - Individual member details (TTL: 10 minutes)
- `relationship:{fromId}:{toId}` - Relationship path result (TTL: 10 minutes)
- `search:members:{query}` - Member search results (TTL: 2 minutes)
- `events:upcoming:{userId}` - Upcoming events list (TTL: 5 minutes)

**Cache Invalidation Triggers**:
- Member create/update/delete → Invalidate `family_tree:*`, `member:{id}`
- Relationship create/delete → Invalidate `family_tree:*`, `relationship:*`
- Event create/update/delete → Invalidate `events:upcoming:*`


### API Design

#### RESTful Endpoints

**Authentication Endpoints**:
```
POST   /api/auth/register          # Register new user
POST   /api/auth/login             # Login and get JWT token
POST   /api/auth/refresh           # Refresh JWT token
GET    /api/auth/profile           # Get current user profile
```

**Member Endpoints**:
```
GET    /api/members                # Get all members (with pagination)
POST   /api/members                # Create new member
GET    /api/members/:id            # Get member by ID
PUT    /api/members/:id            # Update member
DELETE /api/members/:id            # Soft delete member
GET    /api/members/search?q=name  # Search members by name
```

**Relationship Endpoints**:
```
GET    /api/relationships          # Get all relationships
POST   /api/relationships          # Create new relationship
DELETE /api/relationships/:id      # Delete relationship
GET    /api/relationships/path     # Find path between two members
       ?from=uuid1&to=uuid2
```

**Event Endpoints**:
```
GET    /api/events                 # Get all events (with filters)
POST   /api/events                 # Create new event
GET    /api/events/:id             # Get event by ID
PUT    /api/events/:id             # Update event
DELETE /api/events/:id             # Delete event
GET    /api/events/upcoming        # Get upcoming events
```

**Family Tree Endpoint**:
```
GET    /api/family/tree            # Get complete family tree
                                   # Returns nodes and edges for React Flow
```

**WebSocket Endpoint**:
```
WS     /ws                         # WebSocket connection for real-time updates
```

#### Request/Response Examples

**POST /api/members** (Create Member):
```json
// Request
{
  "name": "Arjun Kumar",
  "dateOfBirth": "1995-06-15T00:00:00Z",
  "gender": "male",
  "email": "arjun@example.com",
  "phone": "+919876543210"
}

// Response
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Arjun Kumar",
    "dateOfBirth": "1995-06-15T00:00:00Z",
    "gender": "male",
    "email": "arjun@example.com",
    "phone": "+919876543210",
    "avatarUrl": "",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z",
    "isDeleted": false
  },
  "error": ""
}
```

**POST /api/relationships** (Create Relationship):
```json
// Request
{
  "type": "PARENT_OF",
  "fromId": "550e8400-e29b-41d4-a716-446655440000",
  "toId": "660e8400-e29b-41d4-a716-446655440001"
}

// Response
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

**GET /api/relationships/path?from=uuid1&to=uuid2** (Find Relationship):
```json
// Response
{
  "success": true,
  "data": {
    "path": [
      {"id": "uuid1", "name": "Arjun"},
      {"id": "uuid2", "name": "Rajesh"},
      {"id": "uuid3", "name": "Lakshmi"}
    ],
    "relationLabel": "Grandmother (Father's Mother)",
    "kinshipTerm": "Ammamma",
    "description": "Lakshmi is Arjun's grandmother through his father Rajesh"
  },
  "error": ""
}
```

**GET /api/family/tree** (Get Family Tree):
```json
// Response
{
  "success": true,
  "data": {
    "nodes": [
      {
        "id": "uuid1",
        "type": "memberNode",
        "position": {"x": 0, "y": 0},
        "data": {
          "id": "uuid1",
          "name": "Rajesh",
          "avatarUrl": "https://...",
          "relationBadge": "Father",
          "hasUpcomingEvent": true,
          "gender": "male"
        }
      }
    ],
    "edges": [
      {
        "id": "edge1",
        "source": "uuid1",
        "target": "uuid2",
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


### Relationship Engine Algorithm

The Relationship Engine is the core component that computes kinship relationships between family members using graph traversal and rule-based mapping.

#### Algorithm Overview

```mermaid
flowchart TD
    A[Start: Find Relationship] --> B[Query Neo4j for Shortest Path]
    B --> C{Path Found?}
    C -->|No| D[Return 'Not Related']
    C -->|Yes| E[Extract Edge Sequence]
    E --> F[Apply Kinship Mapping Rules]
    F --> G[Generate Natural Language Description]
    G --> H[Return Relationship Result]
```

#### Cypher Query for Path Finding

```cypher
MATCH path = shortestPath(
  (from:Member {id: $fromId})-[*]-(to:Member {id: $toId})
)
WHERE from.isDeleted = false AND to.isDeleted = false
RETURN path
```

#### Kinship Mapping Rules

The engine uses a rule-based system to map edge sequences to Indian kinship terms:

**Direct Relationships**:
- `PARENT_OF` (male) → "Father" / "Nanna"
- `PARENT_OF` (female) → "Mother" / "Amma"
- `SPOUSE_OF` → "Husband" / "Wife" / "Menarikam"
- `SIBLING_OF` (older male) → "Elder Brother" / "Anna"
- `SIBLING_OF` (older female) → "Elder Sister" / "Akka"
- `SIBLING_OF` (younger male) → "Younger Brother" / "Tammudu"
- `SIBLING_OF` (younger female) → "Younger Sister" / "Chelli"

**Two-Hop Relationships**:
- `PARENT_OF` → `PARENT_OF` (male) → "Grandfather" / "Tata"
- `PARENT_OF` → `PARENT_OF` (female) → "Grandmother" / "Ammamma"
- `PARENT_OF` (father) → `SIBLING_OF` (male) → "Uncle (Father's Brother)" / "Babai"
- `PARENT_OF` (father) → `SIBLING_OF` (female) → "Aunt (Father's Sister)" / "Attha"
- `PARENT_OF` (mother) → `SIBLING_OF` (male) → "Uncle (Mother's Brother)" / "Mamayya"
- `PARENT_OF` (mother) → `SIBLING_OF` (female) → "Aunt (Mother's Sister)" / "Pinni"

**Three-Hop and Beyond**:
- For paths longer than 3 hops, the engine generates descriptive labels like "Great Grandfather", "Second Cousin", etc.
- The algorithm uses generation counting and lateral distance to determine cousin relationships

#### Implementation Pseudocode

```go
func FindRelationship(fromID, toID string) (*RelationshipResult, error) {
    // 1. Query Neo4j for shortest path
    path := neo4j.ShortestPath(fromID, toID)
    
    if path == nil {
        return &RelationshipResult{
            Path: nil,
            RelationLabel: "Not Related",
            KinshipTerm: "",
            Description: "No family connection found"
        }, nil
    }
    
    // 2. Extract edge sequence and node genders
    edges := extractEdges(path)
    nodes := extractNodes(path)
    
    // 3. Apply kinship mapping
    relationLabel := mapToKinshipTerm(edges, nodes)
    
    // 4. Generate natural language description
    description := generateDescription(nodes, relationLabel)
    
    return &RelationshipResult{
        Path: nodes,
        RelationLabel: relationLabel,
        KinshipTerm: getTeluguTerm(relationLabel),
        Description: description
    }, nil
}

func mapToKinshipTerm(edges []Edge, nodes []Node) string {
    if len(edges) == 1 {
        return mapDirectRelationship(edges[0], nodes)
    } else if len(edges) == 2 {
        return mapTwoHopRelationship(edges, nodes)
    } else {
        return mapMultiHopRelationship(edges, nodes)
    }
}
```

#### Caching Strategy

Relationship queries are expensive, so results are cached:
- Cache key: `relationship:{fromId}:{toId}`
- TTL: 10 minutes
- Invalidation: When any relationship is created or deleted


### Notification System Design

The notification system uses a background scheduler with a goroutine pool for concurrent dispatch.

#### Architecture

```mermaid
flowchart LR
    A[Notification Scheduler<br/>Goroutine] --> B[Query Pending<br/>Notifications]
    B --> C[Dispatch Pool<br/>10 Workers]
    C --> D[WhatsApp Worker]
    C --> E[SMS Worker]
    C --> F[Email Worker]
    D --> G[Twilio API]
    E --> G
    F --> H[SendGrid API]
    G --> I[Update Status]
    H --> I
    I --> J{Success?}
    J -->|Yes| K[Mark as Sent]
    J -->|No| L[Increment Retry]
    L --> M{Retry < 3?}
    M -->|Yes| N[Schedule Retry<br/>Exponential Backoff]
    M -->|No| O[Mark as Failed]
```

#### Scheduler Implementation

```go
type NotificationScheduler struct {
    db          *gorm.DB
    notifSvc    *NotificationService
    ticker      *time.Ticker
    workerPool  chan struct{}
}

func (s *NotificationScheduler) Start() {
    s.ticker = time.NewTicker(1 * time.Hour)
    s.workerPool = make(chan struct{}, 10) // Max 10 concurrent workers
    
    go func() {
        for range s.ticker.C {
            s.processNotifications()
        }
    }()
}

func (s *NotificationScheduler) processNotifications() {
    // Query pending notifications
    var notifications []Notification
    s.db.Where("status = ? AND scheduled_at <= ?", "pending", time.Now()).
        Find(&notifications)
    
    // Dispatch each notification in a goroutine
    for _, notif := range notifications {
        s.workerPool <- struct{}{} // Acquire worker slot
        
        go func(n Notification) {
            defer func() { <-s.workerPool }() // Release worker slot
            
            err := s.notifSvc.Dispatch(n)
            if err != nil {
                s.handleFailure(n, err)
            } else {
                s.markAsSent(n)
            }
        }(notif)
    }
}

func (s *NotificationScheduler) handleFailure(n Notification, err error) {
    n.RetryCount++
    n.ErrorMsg = err.Error()
    
    if n.RetryCount >= 3 {
        n.Status = "failed"
    } else {
        // Exponential backoff: 5min, 15min, 45min
        backoff := time.Duration(math.Pow(3, float64(n.RetryCount))) * 5 * time.Minute
        n.ScheduledAt = time.Now().Add(backoff)
    }
    
    s.db.Save(&n)
}
```

#### Notification Templates

**WhatsApp Template** (Twilio approved):
```
🎉 Upcoming Event Reminder

{{event_title}} is coming up on {{event_date}}!

Family members: {{member_names}}

Don't forget to wish them! 🎂

- VamsaSetu
```

**Email Template** (HTML):
```html
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: 'DM Sans', sans-serif; background: #FBF5E6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #E8650A; color: white; padding: 20px; }
        .content { background: white; padding: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Upcoming Event Reminder</h1>
        </div>
        <div class="content">
            <h2>{{event_title}}</h2>
            <p><strong>Date:</strong> {{event_date}}</p>
            <p><strong>Family Members:</strong> {{member_names}}</p>
            <p>Don't forget to wish them!</p>
        </div>
    </div>
</body>
</html>
```

**SMS Template**:
```
VamsaSetu Reminder: {{event_title}} on {{event_date}}. Members: {{member_names}}. Don't forget!
```


### Authentication Flow

VamsaSetu uses JWT-based authentication with access and refresh tokens.

#### Authentication Sequence

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant Auth
    participant DB
    participant Redis
    
    Client->>API: POST /api/auth/login
    API->>Auth: Validate credentials
    Auth->>DB: Query user by email
    DB-->>Auth: User record
    Auth->>Auth: Compare password hash
    Auth->>Auth: Generate JWT tokens
    Auth->>Redis: Store refresh token
    Auth-->>API: Access + Refresh tokens
    API-->>Client: {accessToken, refreshToken}
    
    Note over Client: Store tokens in localStorage
    
    Client->>API: GET /api/members (with JWT)
    API->>Auth: Validate JWT
    Auth->>Auth: Verify signature & expiry
    Auth-->>API: User context
    API->>DB: Query members
    DB-->>API: Member data
    API-->>Client: Response
    
    Note over Client: Access token expires
    
    Client->>API: POST /api/auth/refresh
    API->>Auth: Validate refresh token
    Auth->>Redis: Check refresh token
    Redis-->>Auth: Token valid
    Auth->>Auth: Generate new access token
    Auth-->>API: New access token
    API-->>Client: {accessToken}
```

#### JWT Token Structure

**Access Token** (expires in 15 minutes):
```json
{
  "sub": "user-id",
  "email": "user@example.com",
  "role": "owner",
  "exp": 1705329600,
  "iat": 1705328700
}
```

**Refresh Token** (expires in 7 days):
```json
{
  "sub": "user-id",
  "type": "refresh",
  "exp": 1705933500,
  "iat": 1705328700
}
```

#### Middleware Implementation

```go
func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract token from Authorization header
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(APIResponse{
                Success: false,
                Data:    nil,
                Error:   "Missing authorization header",
            })
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        // Validate JWT
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        
        if err != nil || !token.Valid {
            return c.Status(401).JSON(APIResponse{
                Success: false,
                Data:    nil,
                Error:   "Invalid or expired token",
            })
        }
        
        // Extract claims
        claims := token.Claims.(jwt.MapClaims)
        c.Locals("userId", claims["sub"])
        c.Locals("userRole", claims["role"])
        
        return c.Next()
    }
}
```

#### Role-Based Access Control

```go
func RequireRole(allowedRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole := c.Locals("userRole").(string)
        
        for _, role := range allowedRoles {
            if userRole == role {
                return c.Next()
            }
        }
        
        return c.Status(403).JSON(APIResponse{
            Success: false,
            Data:    nil,
            Error:   "Insufficient permissions",
        })
    }
}

// Usage in routes
app.Post("/api/members", AuthMiddleware(), RequireRole("owner", "admin"), memberHandler.Create)
app.Get("/api/members", AuthMiddleware(), RequireRole("owner", "viewer", "admin"), memberHandler.GetAll)
```


### Frontend Architecture

#### Component Hierarchy

```
App
├── AuthProvider (Context)
├── QueryClientProvider (React Query)
└── Router
    ├── PublicRoute
    │   ├── LoginPage
    │   └── RegisterPage
    └── ProtectedRoute
        ├── DashboardLayout
        │   ├── Navbar
        │   ├── Sidebar (Desktop) / BottomNav (Mobile)
        │   └── Outlet
        │       ├── DashboardPage
        │       ├── FamilyTreePage
        │       │   ├── TreeCanvas (React Flow)
        │       │   │   ├── MemberNode (Custom)
        │       │   │   └── RelationshipEdge (Custom)
        │       │   ├── TreeControls
        │       │   ├── MemberDetailsPanel
        │       │   └── AddMemberModal
        │       └── EventsPage
        │           ├── EventCalendar
        │           ├── EventList
        │           └── AddEventModal
```

#### State Management Strategy

**React Query** (Server State):
- API data fetching and caching
- Automatic background refetching
- Optimistic updates
- Cache invalidation

```typescript
// Example: useMembers hook
export function useMembers() {
  return useQuery({
    queryKey: ['members'],
    queryFn: memberService.getAll,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useCreateMember() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: memberService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members'] });
      queryClient.invalidateQueries({ queryKey: ['familyTree'] });
    },
  });
}
```

**Zustand** (Client State):
- UI state (modals, sidebars, selected items)
- User preferences (theme, language)
- Temporary form data

```typescript
// authStore.ts
interface AuthState {
  user: User | null;
  accessToken: string | null;
  setAuth: (user: User, token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      setAuth: (user, token) => set({ user, accessToken: token }),
      logout: () => set({ user: null, accessToken: null }),
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
    }
  )
);

// uiStore.ts
interface UIState {
  isSidebarOpen: boolean;
  selectedMemberId: string | null;
  isAddMemberModalOpen: boolean;
  toggleSidebar: () => void;
  setSelectedMember: (id: string | null) => void;
  openAddMemberModal: () => void;
  closeAddMemberModal: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  isSidebarOpen: true,
  selectedMemberId: null,
  isAddMemberModalOpen: false,
  toggleSidebar: () => set((state) => ({ isSidebarOpen: !state.isSidebarOpen })),
  setSelectedMember: (id) => set({ selectedMemberId: id }),
  openAddMemberModal: () => set({ isAddMemberModalOpen: true }),
  closeAddMemberModal: () => set({ isAddMemberModalOpen: false }),
}));
```

#### React Flow Integration

**Custom Member Node Component**:
```typescript
import { Handle, Position } from 'reactflow';

interface MemberNodeProps {
  data: MemberNodeData;
}

export function MemberNode({ data }: MemberNodeProps) {
  const { setSelectedMember } = useUIStore();
  
  return (
    <div
      className={cn(
        "relative w-32 h-40 rounded-lg border-2 transition-all",
        "hover:shadow-lg hover:scale-105 cursor-pointer",
        data.gender === 'male' ? 'border-blue-400' : 'border-pink-400',
        data.hasUpcomingEvent && 'ring-2 ring-amber-400 animate-pulse'
      )}
      onClick={() => setSelectedMember(data.id)}
    >
      <Handle type="target" position={Position.Top} />
      
      <div className="p-2 flex flex-col items-center">
        <img
          src={data.avatarUrl || '/default-avatar.png'}
          alt={data.name}
          className="w-16 h-16 rounded-full object-cover"
        />
        <p className="mt-2 text-sm font-medium text-center">{data.name}</p>
        {data.relationBadge && (
          <span className="mt-1 px-2 py-0.5 text-xs bg-teal-100 text-teal-800 rounded">
            {data.relationBadge}
          </span>
        )}
        {data.hasUpcomingEvent && (
          <span className="absolute top-1 right-1 w-3 h-3 bg-amber-400 rounded-full" />
        )}
      </div>
      
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
}
```

**Tree Canvas Component**:
```typescript
import ReactFlow, { Background, Controls, MiniMap } from 'reactflow';

const nodeTypes = {
  memberNode: MemberNode,
};

export function TreeCanvas() {
  const { data: treeData, isLoading } = useFamilyTree();
  const [nodes, setNodes] = useState([]);
  const [edges, setEdges] = useState([]);
  
  useEffect(() => {
    if (treeData) {
      setNodes(treeData.nodes);
      setEdges(treeData.edges);
    }
  }, [treeData]);
  
  if (isLoading) return <LoadingSpinner />;
  
  return (
    <div className="w-full h-full">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        fitView
        attributionPosition="bottom-left"
      >
        <Background />
        <Controls />
        <MiniMap />
      </ReactFlow>
    </div>
  );
}
```


### Real-Time Updates via WebSocket

#### WebSocket Hub Architecture

```go
type Client struct {
    ID     string
    UserID uint
    Conn   *websocket.Conn
    Send   chan []byte
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.Send)
            }
            h.mu.Unlock()
            
        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.Send <- message:
                default:
                    close(client.Send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *Hub) BroadcastUpdate(updateType string, data interface{}) {
    message := map[string]interface{}{
        "type": updateType,
        "data": data,
        "timestamp": time.Now(),
    }
    
    jsonMessage, _ := json.Marshal(message)
    h.broadcast <- jsonMessage
}
```

#### WebSocket Handler

```go
func HandleWebSocket(c *fiber.Ctx, hub *Hub) error {
    // Upgrade HTTP connection to WebSocket
    conn, err := websocket.Upgrade(c)
    if err != nil {
        return err
    }
    
    // Extract user ID from JWT
    userID := c.Locals("userId").(uint)
    
    client := &Client{
        ID:     uuid.New().String(),
        UserID: userID,
        Conn:   conn,
        Send:   make(chan []byte, 256),
    }
    
    hub.register <- client
    
    // Start goroutines for reading and writing
    go client.writePump()
    go client.readPump(hub)
    
    return nil
}

func (c *Client) writePump() {
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-c.Send:
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            c.Conn.WriteMessage(websocket.TextMessage, message)
            
        case <-ticker.C:
            // Send ping to keep connection alive
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func (c *Client) readPump(hub *Hub) {
    defer func() {
        hub.unregister <- c
        c.Conn.Close()
    }()
    
    for {
        _, _, err := c.Conn.ReadMessage()
        if err != nil {
            break
        }
    }
}
```

#### Broadcasting Updates

When data changes occur, broadcast to all connected clients:

```go
// In member service after creating a member
func (s *MemberService) Create(member *Member) error {
    // Save to database
    err := s.repo.Create(member)
    if err != nil {
        return err
    }
    
    // Invalidate cache
    s.cache.Delete("family_tree:*")
    
    // Broadcast WebSocket update
    s.hub.BroadcastUpdate("member_created", member)
    
    return nil
}
```

#### Frontend WebSocket Hook

```typescript
export function useWebSocket() {
  const { accessToken } = useAuthStore();
  const queryClient = useQueryClient();
  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  
  useEffect(() => {
    if (!accessToken) return;
    
    const ws = new WebSocket(`ws://localhost:8080/ws?token=${accessToken}`);
    wsRef.current = ws;
    
    ws.onopen = () => {
      setIsConnected(true);
      console.log('WebSocket connected');
    };
    
    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      
      switch (message.type) {
        case 'member_created':
        case 'member_updated':
        case 'member_deleted':
          queryClient.invalidateQueries({ queryKey: ['members'] });
          queryClient.invalidateQueries({ queryKey: ['familyTree'] });
          break;
          
        case 'relationship_created':
        case 'relationship_deleted':
          queryClient.invalidateQueries({ queryKey: ['relationships'] });
          queryClient.invalidateQueries({ queryKey: ['familyTree'] });
          break;
          
        case 'event_created':
        case 'event_updated':
        case 'event_deleted':
          queryClient.invalidateQueries({ queryKey: ['events'] });
          break;
      }
    };
    
    ws.onclose = () => {
      setIsConnected(false);
      console.log('WebSocket disconnected');
      
      // Attempt reconnection after 3 seconds
      setTimeout(() => {
        if (accessToken) {
          // Reconnect logic
        }
      }, 3000);
    };
    
    return () => {
      ws.close();
    };
  }, [accessToken, queryClient]);
  
  return { isConnected };
}
```


### Caching Strategy

#### Cache Patterns

**Read-Through Cache**:
```go
func (s *MemberService) GetByID(id string) (*Member, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("member:%s", id)
    cached, err := s.cache.Get(cacheKey)
    if err == nil {
        var member Member
        json.Unmarshal([]byte(cached), &member)
        return &member, nil
    }
    
    // Cache miss - query database
    member, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    memberJSON, _ := json.Marshal(member)
    s.cache.Set(cacheKey, string(memberJSON), 10*time.Minute)
    
    return member, nil
}
```

**Write-Through Cache**:
```go
func (s *MemberService) Update(member *Member) error {
    // Update database
    err := s.repo.Update(member)
    if err != nil {
        return err
    }
    
    // Update cache
    cacheKey := fmt.Sprintf("member:%s", member.ID)
    memberJSON, _ := json.Marshal(member)
    s.cache.Set(cacheKey, string(memberJSON), 10*time.Minute)
    
    // Invalidate related caches
    s.cache.Delete("family_tree:*")
    
    return nil
}
```

**Cache-Aside Pattern** (for complex queries):
```go
func (s *MemberService) GetFamilyTree(userID uint) (*FamilyTree, error) {
    cacheKey := fmt.Sprintf("family_tree:%d", userID)
    
    // Check cache
    cached, err := s.cache.Get(cacheKey)
    if err == nil {
        var tree FamilyTree
        json.Unmarshal([]byte(cached), &tree)
        return &tree, nil
    }
    
    // Build tree from database
    tree, err := s.buildFamilyTree(userID)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    treeJSON, _ := json.Marshal(tree)
    s.cache.Set(cacheKey, string(treeJSON), 5*time.Minute)
    
    return tree, nil
}
```

#### Cache Invalidation Rules

| Operation | Invalidate Keys |
|-----------|----------------|
| Create Member | `family_tree:*`, `search:members:*` |
| Update Member | `member:{id}`, `family_tree:*` |
| Delete Member | `member:{id}`, `family_tree:*`, `relationship:*:{id}`, `relationship:{id}:*` |
| Create Relationship | `family_tree:*`, `relationship:*` |
| Delete Relationship | `family_tree:*`, `relationship:*` |
| Create Event | `events:upcoming:*` |
| Update Event | `events:upcoming:*` |
| Delete Event | `events:upcoming:*` |

#### TTL Configuration

```go
const (
    MemberCacheTTL       = 10 * time.Minute
    FamilyTreeCacheTTL   = 5 * time.Minute
    RelationshipCacheTTL = 10 * time.Minute
    SearchCacheTTL       = 2 * time.Minute
    EventsCacheTTL       = 5 * time.Minute
)
```


### UI Component Design and Cultural Theming

#### Color Palette

```typescript
// tailwind.config.js
export default {
  theme: {
    extend: {
      colors: {
        saffron: {
          DEFAULT: '#E8650A',
          light: '#FF8C42',
          dark: '#C55000',
        },
        turmeric: {
          DEFAULT: '#F5A623',
          light: '#FFD166',
          dark: '#E89005',
        },
        ivory: {
          DEFAULT: '#FBF5E6',
          light: '#FFFEF9',
          dark: '#F5EDD6',
        },
        teal: {
          DEFAULT: '#0D4A52',
          light: '#0D9488',
          dark: '#083339',
        },
        charcoal: {
          DEFAULT: '#2C2420',
          light: '#3D3530',
          dark: '#1A1512',
        },
      },
      fontFamily: {
        heading: ['Playfair Display', 'serif'],
        body: ['DM Sans', 'sans-serif'],
      },
    },
  },
};
```

#### Rangoli Pattern Background

```typescript
// RangoliPattern.tsx
export function RangoliPattern() {
  return (
    <svg className="absolute inset-0 w-full h-full opacity-5" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <pattern id="rangoli" x="0" y="0" width="100" height="100" patternUnits="userSpaceOnUse">
          <circle cx="50" cy="50" r="2" fill="#E8650A" />
          <path d="M50,30 L60,40 L50,50 L40,40 Z" fill="#F5A623" />
          <path d="M50,50 L60,60 L50,70 L40,60 Z" fill="#0D9488" />
          <circle cx="30" cy="30" r="1" fill="#E8650A" />
          <circle cx="70" cy="30" r="1" fill="#E8650A" />
          <circle cx="30" cy="70" r="1" fill="#E8650A" />
          <circle cx="70" cy="70" r="1" fill="#E8650A" />
        </pattern>
      </defs>
      <rect width="100%" height="100%" fill="url(#rangoli)" />
    </svg>
  );
}
```

#### Relationship Edge Colors

```typescript
export const EDGE_COLORS = {
  SPOUSE_OF: '#E11D48',      // Rose
  PARENT_OF: '#0D9488',      // Teal
  SIBLING_OF: '#F59E0B',     // Amber
};

export function getEdgeStyle(type: string): React.CSSProperties {
  return {
    stroke: EDGE_COLORS[type] || '#6B7280',
    strokeWidth: 2,
  };
}
```

#### Animation Variants (Framer Motion)

```typescript
// animations.ts
export const pageTransition = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -20 },
  transition: { duration: 0.3 },
};

export const modalTransition = {
  initial: { opacity: 0, scale: 0.95 },
  animate: { opacity: 1, scale: 1 },
  exit: { opacity: 0, scale: 0.95 },
  transition: { duration: 0.2 },
};

export const pathDrawing = {
  initial: { pathLength: 0 },
  animate: { pathLength: 1 },
  transition: { duration: 1, ease: 'easeInOut' },
};
```

#### Responsive Breakpoints

```typescript
// useResponsive.ts
export function useResponsive() {
  const [breakpoint, setBreakpoint] = useState<'mobile' | 'tablet' | 'desktop'>('desktop');
  
  useEffect(() => {
    const handleResize = () => {
      const width = window.innerWidth;
      if (width < 768) {
        setBreakpoint('mobile');
      } else if (width < 1024) {
        setBreakpoint('tablet');
      } else {
        setBreakpoint('desktop');
      }
    };
    
    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);
  
  return {
    isMobile: breakpoint === 'mobile',
    isTablet: breakpoint === 'tablet',
    isDesktop: breakpoint === 'desktop',
  };
}
```

#### Mobile Navigation

```typescript
// BottomNav.tsx (Mobile)
export function BottomNav() {
  const location = useLocation();
  
  const navItems = [
    { path: '/dashboard', icon: HomeIcon, label: 'Home' },
    { path: '/family', icon: UsersIcon, label: 'Family' },
    { path: '/events', icon: CalendarIcon, label: 'Events' },
    { path: '/profile', icon: UserIcon, label: 'Profile' },
  ];
  
  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-charcoal border-t border-charcoal-light md:hidden">
      <div className="flex justify-around items-center h-16">
        {navItems.map((item) => (
          <Link
            key={item.path}
            to={item.path}
            className={cn(
              "flex flex-col items-center justify-center flex-1 h-full",
              location.pathname === item.path ? 'text-saffron' : 'text-gray-400'
            )}
          >
            <item.icon className="w-6 h-6" />
            <span className="text-xs mt-1">{item.label}</span>
          </Link>
        ))}
      </div>
    </nav>
  );
}
```


### Docker Orchestration

#### Docker Compose Configuration

```yaml
version: '3.8'

services:
  # Frontend Service
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      - REACT_APP_API_URL=http://localhost:8080
      - REACT_APP_WS_URL=ws://localhost:8080
    depends_on:
      - backend
    volumes:
      - ./frontend:/app
      - /app/node_modules
    networks:
      - vamsasetu-network

  # Backend Service
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - POSTGRES_URL=postgresql://vamsasetu:vamsasetu123@postgres:5432/vamsasetu
      - NEO4J_URI=bolt://neo4j:7687
      - NEO4J_USERNAME=neo4j
      - NEO4J_PASSWORD=vamsasetu123
      - REDIS_ADDR=redis:6379
      - JWT_SECRET=${JWT_SECRET}
      - SENDGRID_API_KEY=${SENDGRID_API_KEY}
      - TWILIO_ACCOUNT_SID=${TWILIO_ACCOUNT_SID}
      - TWILIO_AUTH_TOKEN=${TWILIO_AUTH_TOKEN}
      - TWILIO_PHONE_NUMBER=${TWILIO_PHONE_NUMBER}
      - TWILIO_WHATSAPP_NUMBER=${TWILIO_WHATSAPP_NUMBER}
    depends_on:
      neo4j:
        condition: service_healthy
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./backend:/app
    networks:
      - vamsasetu-network

  # Neo4j Graph Database
  neo4j:
    image: neo4j:5.13
    ports:
      - "7474:7474"  # HTTP
      - "7687:7687"  # Bolt
    environment:
      - NEO4J_AUTH=neo4j/vamsasetu123
      - NEO4J_PLUGINS=["apoc"]
      - NEO4J_dbms_security_procedures_unrestricted=apoc.*
    volumes:
      - neo4j_data:/data
      - neo4j_logs:/logs
    healthcheck:
      test: ["CMD-SHELL", "cypher-shell -u neo4j -p vamsasetu123 'RETURN 1'"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - vamsasetu-network

  # PostgreSQL Database
  postgres:
    image: postgres:15
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_DB=vamsasetu
      - POSTGRES_USER=vamsasetu
      - POSTGRES_PASSWORD=vamsasetu123
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U vamsasetu"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - vamsasetu-network

  # Redis Cache
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
    networks:
      - vamsasetu-network

volumes:
  neo4j_data:
  neo4j_logs:
  postgres_data:
  redis_data:

networks:
  vamsasetu-network:
    driver: bridge
```

#### Backend Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run
CMD ["./main"]
```

#### Frontend Dockerfile

```dockerfile
# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build app
RUN npm run build

# Runtime stage
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/build /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
```

#### Health Check Endpoints

```go
// In main.go
app.Get("/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "healthy",
        "timestamp": time.Now(),
        "services": fiber.Map{
            "neo4j": checkNeo4jHealth(),
            "postgres": checkPostgresHealth(),
            "redis": checkRedisHealth(),
        },
    })
})

func checkNeo4jHealth() string {
    // Ping Neo4j
    _, err := neo4jDriver.VerifyConnectivity()
    if err != nil {
        return "unhealthy"
    }
    return "healthy"
}

func checkPostgresHealth() string {
    // Ping PostgreSQL
    err := db.Exec("SELECT 1").Error
    if err != nil {
        return "unhealthy"
    }
    return "healthy"
}

func checkRedisHealth() string {
    // Ping Redis
    _, err := redisClient.Ping(context.Background()).Result()
    if err != nil {
        return "unhealthy"
    }
    return "healthy"
}
```


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

After analyzing all acceptance criteria, I've identified the following testable properties. Properties that are redundant or can be combined have been consolidated to ensure each property provides unique validation value.

### Property 1: User Registration Round Trip

*For any* valid registration data (email, password, name), creating a user account should result in a user record being stored with a hashed password, and that user should be able to log in with the original credentials.

**Validates: Requirements 1.1, 1.2**

### Property 2: JWT Authentication Round Trip

*For any* valid user credentials, logging in should produce a JWT token, and that token should grant access to protected endpoints until expiry.

**Validates: Requirements 1.3**

### Property 3: Invalid Token Rejection

*For any* expired or malformed JWT token, authentication should fail with an appropriate error.

**Validates: Requirements 1.4**

### Property 4: User Role Invariant

*For any* user in the system, that user must have exactly one role from the set {owner, viewer, admin}.

**Validates: Requirements 1.5**

### Property 5: Role-Based Authorization

*For any* modification operation on family tree data, a user with the "viewer" role should be denied access with an authorization error.

**Validates: Requirements 1.6**

### Property 6: Token Refresh Round Trip

*For any* valid refresh token, requesting a token refresh should produce a new valid access token.

**Validates: Requirements 1.7**

### Property 7: Member Creation and Retrieval

*For any* valid member data (name, date of birth, gender), creating a member should result in that member being stored in Neo4j and retrievable by ID.

**Validates: Requirements 2.1, 11.1**

### Property 8: Member Update Persistence

*For any* existing member and any valid attribute changes, updating the member should persist those changes such that subsequent retrieval returns the updated values.

**Validates: Requirements 2.2**

### Property 9: Soft Delete Preservation

*For any* member, deleting that member should set the isDeleted flag to true while preserving all other member data.

**Validates: Requirements 2.3**

### Property 10: Relationship Creation and Retrieval

*For any* two existing members and a valid relationship type (SPOUSE_OF, PARENT_OF, SIBLING_OF), creating a relationship should result in that relationship being stored in Neo4j and retrievable.

**Validates: Requirements 2.4**

### Property 11: Relationship Deletion

*For any* existing relationship, deleting that relationship should remove it from Neo4j such that it is no longer retrievable.

**Validates: Requirements 2.5**

### Property 12: Relationship Semantic Validation

*For any* member, attempting to create a PARENT_OF relationship from that member to itself should be rejected with a validation error.

**Validates: Requirements 2.6**

### Property 13: Soft Delete Enforcement for Connected Members

*For any* member with at least one relationship, attempting to delete that member should result in a soft delete (isDeleted=true) rather than removal from the database.

**Validates: Requirements 2.7**

### Property 14: Family Tree Format Validity

*For any* family tree request, the returned data structure should contain "nodes" and "edges" arrays where each node has the required React Flow properties (id, type, position, data) and each edge has the required properties (id, source, target, type).

**Validates: Requirements 3.1**

### Property 15: Member Node Rendering Completeness

*For any* member, the rendered node data should include all required fields: id, name, avatarUrl, relationBadge, hasUpcomingEvent, and gender.

**Validates: Requirements 3.2**

### Property 16: Relationship Edge Color Mapping

*For any* relationship edge, the edge style should have a stroke color that matches the relationship type: rose for SPOUSE_OF, teal for PARENT_OF, amber for SIBLING_OF.

**Validates: Requirements 3.3**

### Property 17: Upcoming Event Indicator

*For any* member with an event scheduled within 7 days, the member's node data should have hasUpcomingEvent set to true.

**Validates: Requirements 3.4**

### Property 18: Member Node Click Interaction

*For any* member node, clicking that node should trigger the display of a side panel containing the member's full details.

**Validates: Requirements 3.6**

### Property 19: Relationship Path Finding

*For any* two members in the same family tree, requesting the relationship between them should return either a valid path with nodes and edges, or a "not related" result if no path exists.

**Validates: Requirements 4.1, 4.4**

### Property 20: Relationship Result Completeness

*For any* relationship query that finds a path, the result should contain all required fields: path nodes, relation label, kinship term, and natural language description.

**Validates: Requirements 4.2, 4.3**

### Property 21: Event Creation and Retrieval

*For any* valid event data (title, date, type, member IDs), creating an event should result in that event being stored in PostgreSQL and retrievable by ID.

**Validates: Requirements 5.1, 11.2**

### Property 22: Event Update Persistence

*For any* existing event and any valid attribute changes, updating the event should persist those changes such that subsequent retrieval returns the updated values.

**Validates: Requirements 5.2**

### Property 23: Event Deletion

*For any* existing event, deleting that event should remove it from PostgreSQL such that it is no longer retrievable.

**Validates: Requirements 5.3**

### Property 24: Event Type Validity

*For any* event in the system, that event's type must be one of: birthday, anniversary, ceremony, or custom.

**Validates: Requirements 5.4**

### Property 25: Event Countdown Display

*For any* event scheduled within 7 days from now, the event card should display a countdown indicator.

**Validates: Requirements 5.6**

### Property 26: Event Filtering

*For any* event filter criteria (type, member ID, date range), the returned events should match all specified filter criteria.

**Validates: Requirements 5.7, 8.3, 8.4**

### Property 27: Notification Scheduling

*For any* event creation, the system should create notification records in PostgreSQL based on the configured notification preferences.

**Validates: Requirements 6.1**

### Property 28: Notification Dispatch

*For any* notification with status "pending" and scheduledAt <= now, the notification service should attempt to dispatch it via the configured channel (whatsapp, sms, email).

**Validates: Requirements 6.3**

### Property 29: Notification Status Update on Success

*For any* notification that is successfully dispatched, the notification's status should be updated to "sent" and sentAt should be set to the current timestamp.

**Validates: Requirements 6.5**

### Property 30: Notification Retry Logic

*For any* notification that fails to dispatch, the system should increment retryCount and reschedule the notification with exponential backoff, up to a maximum of 3 retries.

**Validates: Requirements 6.6**

### Property 31: WebSocket Connection Establishment

*For any* valid JWT token, connecting to the WebSocket endpoint with that token should result in a successful connection.

**Validates: Requirements 7.1**

### Property 32: WebSocket Broadcast

*For any* data modification (member, relationship, event), all connected WebSocket clients should receive a broadcast message containing the modification type and data.

**Validates: Requirements 7.2**

### Property 33: Member Search

*For any* search query string, the returned members should have names that contain the query string (case-insensitive partial match).

**Validates: Requirements 8.1**

### Property 34: Member Filtering

*For any* member filter criteria (gender, generation, relationship type), the returned members should match all specified filter criteria.

**Validates: Requirements 8.2**

### Property 35: Database Transaction Rollback

*For any* database operation that fails, the system should roll back the transaction such that no partial changes are persisted.

**Validates: Requirements 11.4**

### Property 36: Audit Logging

*For any* data modification operation (create, update, delete), an audit log entry should be created in PostgreSQL with the operation type, entity type, entity ID, user ID, and timestamp.

**Validates: Requirements 11.5**

### Property 37: API Response Format Consistency

*For any* API response, the response body should be a JSON object with exactly three fields: success (boolean), data (any type or null), and error (string).

**Validates: Requirements 13.1**

### Property 38: Success Response Format

*For any* successful API request, the response should have success=true, data populated with the result, and error as an empty string.

**Validates: Requirements 13.2**

### Property 39: Error Response Format

*For any* failed API request, the response should have success=false, data=null, and error populated with a descriptive message.

**Validates: Requirements 13.3**

### Property 40: HTTP Status Code Mapping

*For any* API response, the HTTP status code should match the response outcome: 2xx for success, 4xx for client errors, 5xx for server errors.

**Validates: Requirements 13.4**

### Property 41: Client-Side Validation

*For any* form submission with invalid data (missing required fields, invalid format), client-side validation should prevent submission and display inline error messages.

**Validates: Requirements 14.1**

### Property 42: Server-Side Validation

*For any* API request with invalid data, server-side validation should reject the request and return a descriptive error message.

**Validates: Requirements 14.2**

### Property 43: Database Constraint Error Handling

*For any* database operation that violates a constraint (unique, foreign key, check), the system should return a user-friendly error message rather than exposing the raw database error.

**Validates: Requirements 14.4**

### Property 44: Cache Read-Through

*For any* request for frequently accessed data (family tree, member details), if the data exists in Redis cache, it should be returned from cache without querying the database.

**Validates: Requirements 17.1**

### Property 45: Cache Write-Through

*For any* cache miss on frequently accessed data, the system should query the database, store the result in Redis with appropriate TTL, and return the data.

**Validates: Requirements 17.3**

### Property 46: Cache Invalidation on Modification

*For any* modification to a member or relationship, all related cache entries (family_tree:*, member:{id}, relationship:*) should be invalidated.

**Validates: Requirements 17.4**

### Property 47: React Query Cache Invalidation

*For any* successful API mutation (create, update, delete), the relevant React Query cache entries should be invalidated to trigger refetch.

**Validates: Requirements 18.3**

### Property 48: User Preference Persistence

*For any* user preference change (theme, notification settings), the preference should be persisted to localStorage and retrievable on subsequent page loads.

**Validates: Requirements 18.4**

### Property 49: ARIA Label Presence

*For any* interactive element (button, link, input), the element should have an appropriate ARIA label or aria-label attribute.

**Validates: Requirements 19.2**

### Property 50: Loading State Display

*For any* asynchronous operation (API call, data fetch), a loading indicator should be displayed while the operation is in progress.

**Validates: Requirements 19.5**

### Property 51: Error State Display

*For any* error condition, an error message should be displayed with clear information about what went wrong and how to recover.

**Validates: Requirements 19.6**


## Error Handling

### Error Categories

VamsaSetu implements comprehensive error handling across all layers:

**Validation Errors** (400 Bad Request):
- Missing required fields
- Invalid data formats (email, date, phone)
- Business rule violations (self-referential relationships)
- Invalid enum values

**Authentication Errors** (401 Unauthorized):
- Missing JWT token
- Expired JWT token
- Invalid JWT signature
- Malformed token

**Authorization Errors** (403 Forbidden):
- Insufficient role permissions
- Attempting to access another user's family tree
- Viewer attempting to modify data

**Not Found Errors** (404 Not Found):
- Member ID not found
- Event ID not found
- Relationship not found
- User not found

**Conflict Errors** (409 Conflict):
- Duplicate email during registration
- Duplicate relationship between same members
- Attempting to delete member with relationships (hard delete)

**Server Errors** (500 Internal Server Error):
- Database connection failures
- Neo4j query errors
- Redis connection failures
- External service failures (Twilio, SendGrid)

### Error Response Structure

All errors follow the consistent API response format:

```json
{
  "success": false,
  "data": null,
  "error": "Descriptive error message for the user"
}
```

### Backend Error Handling Pattern

```go
// Custom error types
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    return e.Message
}

// Error handling middleware
func ErrorHandler() fiber.Handler {
    return func(c *fiber.Ctx) error {
        err := c.Next()
        
        if err != nil {
            var appErr *AppError
            if errors.As(err, &appErr) {
                return c.Status(appErr.Code).JSON(APIResponse{
                    Success: false,
                    Data:    nil,
                    Error:   appErr.Message,
                })
            }
            
            // Log unexpected errors
            log.Printf("Unexpected error: %v", err)
            
            return c.Status(500).JSON(APIResponse{
                Success: false,
                Data:    nil,
                Error:   "An unexpected error occurred",
            })
        }
        
        return nil
    }
}

// Usage in handlers
func (h *MemberHandler) Create(c *fiber.Ctx) error {
    var member Member
    if err := c.BodyParser(&member); err != nil {
        return &AppError{
            Code:    400,
            Message: "Invalid request body",
            Err:     err,
        }
    }
    
    if err := h.validator.Validate(member); err != nil {
        return &AppError{
            Code:    400,
            Message: fmt.Sprintf("Validation failed: %v", err),
            Err:     err,
        }
    }
    
    createdMember, err := h.service.Create(&member)
    if err != nil {
        return &AppError{
            Code:    500,
            Message: "Failed to create member",
            Err:     err,
        }
    }
    
    return c.Status(201).JSON(APIResponse{
        Success: true,
        Data:    createdMember,
        Error:   "",
    })
}
```

### Frontend Error Handling

```typescript
// Error boundary for React components
export class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean; error: Error | null }
> {
  constructor(props: any) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="flex flex-col items-center justify-center h-screen">
          <h1 className="text-2xl font-bold text-red-600">Something went wrong</h1>
          <p className="mt-2 text-gray-600">{this.state.error?.message}</p>
          <button
            onClick={() => window.location.reload()}
            className="mt-4 px-4 py-2 bg-saffron text-white rounded"
          >
            Reload Page
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}

// API error handling with React Query
export function useCreateMember() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: memberService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members'] });
      toast.success('Member created successfully');
    },
    onError: (error: any) => {
      const message = error.response?.data?.error || 'Failed to create member';
      toast.error(message);
    },
  });
}

// Form validation with error display
export function AddMemberForm() {
  const [errors, setErrors] = useState<Record<string, string>>({});
  const createMember = useCreateMember();
  
  const handleSubmit = (data: MemberFormData) => {
    const validationErrors = validateMemberData(data);
    
    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }
    
    createMember.mutate(data);
  };
  
  return (
    <form onSubmit={handleSubmit}>
      <Input
        name="name"
        label="Name"
        error={errors.name}
        required
      />
      {/* More fields */}
    </form>
  );
}
```

### Database Error Handling

```go
// Transaction with rollback
func (s *MemberService) CreateWithRelationship(member *Member, relationship *Relationship) error {
    // Start transaction
    tx := s.neo4jSession.BeginTransaction()
    defer tx.Close()
    
    // Create member
    _, err := tx.Run(
        "CREATE (m:Member {id: $id, name: $name}) RETURN m",
        map[string]interface{}{
            "id":   member.ID,
            "name": member.Name,
        },
    )
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create member: %w", err)
    }
    
    // Create relationship
    _, err = tx.Run(
        "MATCH (a:Member {id: $fromId}), (b:Member {id: $toId}) " +
        "CREATE (a)-[:PARENT_OF]->(b)",
        map[string]interface{}{
            "fromId": relationship.FromID,
            "toId":   relationship.ToID,
        },
    )
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create relationship: %w", err)
    }
    
    // Commit transaction
    return tx.Commit()
}
```

### Retry Logic for External Services

```go
func (s *NotificationService) SendWithRetry(notification *Notification) error {
    maxRetries := 3
    baseDelay := 5 * time.Second
    
    for attempt := 0; attempt <= maxRetries; attempt++ {
        err := s.send(notification)
        if err == nil {
            return nil
        }
        
        if attempt < maxRetries {
            delay := time.Duration(math.Pow(3, float64(attempt))) * baseDelay
            log.Printf("Retry %d/%d after %v: %v", attempt+1, maxRetries, delay, err)
            time.Sleep(delay)
        }
    }
    
    return fmt.Errorf("failed after %d retries", maxRetries)
}
```


## Testing Strategy

VamsaSetu employs a comprehensive dual testing approach combining unit tests for specific examples and edge cases with property-based tests for universal correctness guarantees.

### Testing Philosophy

**Unit Tests**: Focus on specific examples, edge cases, and integration points
**Property-Based Tests**: Verify universal properties across randomized inputs (minimum 100 iterations per test)

Together, these approaches provide comprehensive coverage where unit tests catch concrete bugs and property-based tests verify general correctness.

### Backend Testing

#### Technology Stack

- **Testing Framework**: Go's built-in `testing` package
- **Property-Based Testing**: `gopter` library
- **Mocking**: `gomock` for interface mocking
- **Database Testing**: Testcontainers for Neo4j and PostgreSQL
- **HTTP Testing**: `httptest` package

#### Test Organization

```
backend/
├── internal/
│   ├── service/
│   │   ├── member_service.go
│   │   ├── member_service_test.go          # Unit tests
│   │   └── member_service_property_test.go # Property tests
│   ├── handler/
│   │   ├── member_handler.go
│   │   └── member_handler_test.go
│   └── repository/
│       ├── member_repo.go
│       └── member_repo_test.go
```

#### Property-Based Test Example

```go
// Feature: vamsasetu-full-system, Property 7: Member Creation and Retrieval
func TestProperty_MemberCreationAndRetrieval(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("For any valid member data, creating a member should result in that member being retrievable by ID",
        prop.ForAll(
            func(name string, dob time.Time, gender string) bool {
                // Setup
                service := setupMemberService(t)
                defer teardown(t)
                
                // Create member
                member := &Member{
                    ID:          uuid.New().String(),
                    Name:        name,
                    DateOfBirth: dob,
                    Gender:      gender,
                }
                
                err := service.Create(member)
                if err != nil {
                    return false
                }
                
                // Retrieve member
                retrieved, err := service.GetByID(member.ID)
                if err != nil {
                    return false
                }
                
                // Verify
                return retrieved.ID == member.ID &&
                       retrieved.Name == member.Name &&
                       retrieved.Gender == member.Gender
            },
            genValidName(),
            genValidDateOfBirth(),
            genValidGender(),
        ),
    )
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Generators for property-based testing
func genValidName() gopter.Gen {
    return gen.AlphaString().SuchThat(func(s string) bool {
        return len(s) >= 2 && len(s) <= 100
    })
}

func genValidDateOfBirth() gopter.Gen {
    return gen.TimeRange(
        time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
        time.Now(),
    )
}

func genValidGender() gopter.Gen {
    return gen.OneConstOf("male", "female", "other")
}
```

#### Unit Test Example

```go
func TestMemberService_Create_ValidData(t *testing.T) {
    service := setupMemberService(t)
    defer teardown(t)
    
    member := &Member{
        ID:          uuid.New().String(),
        Name:        "Arjun Kumar",
        DateOfBirth: time.Date(1995, 6, 15, 0, 0, 0, 0, time.UTC),
        Gender:      "male",
        Email:       "arjun@example.com",
    }
    
    err := service.Create(member)
    assert.NoError(t, err)
    
    retrieved, err := service.GetByID(member.ID)
    assert.NoError(t, err)
    assert.Equal(t, member.Name, retrieved.Name)
}

func TestMemberService_Create_SelfReferentialRelationship(t *testing.T) {
    service := setupRelationshipService(t)
    defer teardown(t)
    
    memberID := uuid.New().String()
    
    relationship := &Relationship{
        Type:   "PARENT_OF",
        FromID: memberID,
        ToID:   memberID,
    }
    
    err := service.Create(relationship)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "cannot be their own parent")
}
```

#### Integration Test Example

```go
func TestMemberAPI_CreateAndRetrieve(t *testing.T) {
    app := setupTestApp(t)
    defer teardown(t)
    
    // Create member
    reqBody := `{
        "name": "Arjun Kumar",
        "dateOfBirth": "1995-06-15T00:00:00Z",
        "gender": "male",
        "email": "arjun@example.com"
    }`
    
    req := httptest.NewRequest("POST", "/api/members", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+getTestToken())
    
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
    
    var createResp APIResponse
    json.NewDecoder(resp.Body).Decode(&createResp)
    assert.True(t, createResp.Success)
    
    memberID := createResp.Data.(map[string]interface{})["id"].(string)
    
    // Retrieve member
    req = httptest.NewRequest("GET", "/api/members/"+memberID, nil)
    req.Header.Set("Authorization", "Bearer "+getTestToken())
    
    resp, err = app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

### Frontend Testing

#### Technology Stack

- **Testing Framework**: Jest + React Testing Library
- **Property-Based Testing**: `fast-check` library
- **Component Testing**: React Testing Library
- **E2E Testing**: Playwright (optional)

#### Test Organization

```
frontend/
├── src/
│   ├── components/
│   │   ├── family/
│   │   │   ├── MemberNode.tsx
│   │   │   ├── MemberNode.test.tsx
│   │   │   └── MemberNode.property.test.tsx
│   │   └── events/
│   │       ├── EventCard.tsx
│   │       └── EventCard.test.tsx
│   ├── hooks/
│   │   ├── useMembers.ts
│   │   └── useMembers.test.ts
│   └── services/
│       ├── memberService.ts
│       └── memberService.test.ts
```

#### Property-Based Test Example (Frontend)

```typescript
// Feature: vamsasetu-full-system, Property 15: Member Node Rendering Completeness
import fc from 'fast-check';

describe('MemberNode - Property Tests', () => {
  it('Property 15: For any member, the rendered node should include all required fields', () => {
    fc.assert(
      fc.property(
        fc.record({
          id: fc.uuid(),
          name: fc.string({ minLength: 2, maxLength: 100 }),
          avatarUrl: fc.webUrl(),
          relationBadge: fc.constantFrom('Father', 'Mother', 'Sibling', 'Child'),
          hasUpcomingEvent: fc.boolean(),
          gender: fc.constantFrom('male', 'female', 'other'),
        }),
        (memberData) => {
          const { container } = render(<MemberNode data={memberData} />);
          
          // Verify all required fields are present
          expect(container.querySelector('[data-testid="member-name"]')).toHaveTextContent(memberData.name);
          expect(container.querySelector('[data-testid="member-avatar"]')).toHaveAttribute('src', memberData.avatarUrl);
          expect(container.querySelector('[data-testid="relation-badge"]')).toHaveTextContent(memberData.relationBadge);
          
          if (memberData.hasUpcomingEvent) {
            expect(container.querySelector('[data-testid="event-indicator"]')).toBeInTheDocument();
          }
          
          return true;
        }
      ),
      { numRuns: 100 }
    );
  });
});
```

#### Unit Test Example (Frontend)

```typescript
describe('MemberNode', () => {
  it('should render member name and avatar', () => {
    const memberData = {
      id: '123',
      name: 'Arjun Kumar',
      avatarUrl: 'https://example.com/avatar.jpg',
      relationBadge: 'Father',
      hasUpcomingEvent: false,
      gender: 'male',
    };
    
    const { getByText, getByAltText } = render(<MemberNode data={memberData} />);
    
    expect(getByText('Arjun Kumar')).toBeInTheDocument();
    expect(getByAltText('Arjun Kumar')).toHaveAttribute('src', 'https://example.com/avatar.jpg');
  });
  
  it('should display event indicator when hasUpcomingEvent is true', () => {
    const memberData = {
      id: '123',
      name: 'Arjun Kumar',
      avatarUrl: 'https://example.com/avatar.jpg',
      relationBadge: 'Father',
      hasUpcomingEvent: true,
      gender: 'male',
    };
    
    const { container } = render(<MemberNode data={memberData} />);
    
    expect(container.querySelector('[data-testid="event-indicator"]')).toBeInTheDocument();
  });
});
```

### Test Configuration

#### Backend Test Configuration

Each property-based test must run a minimum of 100 iterations:

```go
// In gopter properties
properties := gopter.NewProperties(&gopter.TestParameters{
    MinSuccessfulTests: 100,
    MaxSize:            100,
})
```

#### Frontend Test Configuration

```typescript
// In fast-check assertions
fc.assert(
  fc.property(/* ... */),
  { numRuns: 100 }
);
```

### Test Coverage Goals

- **Unit Test Coverage**: Minimum 80% line coverage
- **Property Test Coverage**: All 51 correctness properties implemented
- **Integration Test Coverage**: All API endpoints tested
- **E2E Test Coverage**: Critical user flows (registration, login, create family tree, add event)

### Continuous Integration

```yaml
# .github/workflows/test.yml
name: Test Suite

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: cd backend && go test ./... -v
      - name: Run property tests
        run: cd backend && go test ./... -v -tags=property
      - name: Generate coverage
        run: cd backend && go test ./... -coverprofile=coverage.out
  
  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: cd frontend && npm ci
      - name: Run tests
        run: cd frontend && npm test -- --coverage
```

### Test Data Management

**Seed Data for Testing**:
- Use factories to generate test data
- Testcontainers for isolated database instances
- Reset database state between tests
- Use transactions that rollback for unit tests

**Test Fixtures**:
```go
// fixtures/members.go
func CreateTestMember(name string) *Member {
    return &Member{
        ID:          uuid.New().String(),
        Name:        name,
        DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
        Gender:      "male",
        Email:       fmt.Sprintf("%s@test.com", strings.ToLower(name)),
    }
}

func CreateTestFamily() []*Member {
    father := CreateTestMember("Rajesh")
    mother := CreateTestMember("Lakshmi")
    child := CreateTestMember("Arjun")
    return []*Member{father, mother, child}
}
```

