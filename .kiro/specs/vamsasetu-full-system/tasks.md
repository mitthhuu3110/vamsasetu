# Implementation Plan: VamsaSetu Full System

## Overview

This implementation plan breaks down the VamsaSetu family tree and event management system into discrete, actionable coding tasks. The system uses Go + Fiber for the backend, React + TypeScript for the frontend, Neo4j for graph data, PostgreSQL for relational data, and Redis for caching. Tasks follow the specified build order to ensure incremental progress with early validation.

## Tasks

- [x] 1. Infrastructure Setup
  - Create docker-compose.yml with all services (neo4j, postgres, redis, backend, frontend)
  - Create .env.example with all required environment variables documented
  - Configure service dependencies and health checks
  - Set up named volumes for data persistence
  - _Requirements: 12.1, 12.2, 12.3, 12.4, 16.1, 16.2, 16.3, 16.4, 16.5_

- [ ] 2. Backend Foundation
  - [x] 2.1 Initialize Go module and project structure
    - Run `go mod init vamsasetu/backend`
    - Create directory structure: cmd/server, internal/{config,middleware,models,repository,service,handler,scheduler,utils}, pkg/{neo4j,postgres,redis}
    - Install dependencies: fiber, gorm, neo4j-go-driver, go-redis, jwt-go, godotenv
    - _Requirements: 20.1_

  - [x] 2.2 Implement configuration management
    - Create internal/config/config.go to load environment variables
    - Validate required environment variables on startup
    - Implement config struct with all database and service credentials
    - _Requirements: 12.1, 12.2, 12.5_

  - [x] 2.3 Implement database connection clients
    - Create pkg/neo4j/client.go with Neo4j driver initialization
    - Create pkg/postgres/client.go with GORM PostgreSQL connection
    - Create pkg/redis/client.go with Redis client initialization
    - Add connection health check functions
    - _Requirements: 11.1, 11.2_

  - [x] 2.4 Write property test for configuration validation
    - **Property 4: User Role Invariant** (partial - validates config loading)
    - **Validates: Requirements 12.5**

- [ ] 3. Backend Data Models
  - [x] 3.1 Define PostgreSQL models with GORM
    - Create internal/models/user.go (User entity with email, password_hash, name, role)
    - Create internal/models/event.go (Event entity with title, description, event_date, event_type, member_ids)
    - Create internal/models/notification.go (Notification entity with event_id, user_id, channel, scheduled_at, status)
    - Add GORM tags for constraints and indexes
    - _Requirements: 2.1, 5.1, 6.1_

  - [x] 3.2 Define Neo4j models
    - Create internal/models/member.go (Member node with id, name, dateOfBirth, gender, email, phone, avatarUrl, isDeleted)
    - Create internal/models/relationship.go (Relationship edge with type, fromId, toId, createdAt)
    - Define relationship type constants (SPOUSE_OF, PARENT_OF, SIBLING_OF)
    - _Requirements: 2.1, 2.4_

  - [x] 3.3 Run GORM auto-migration
    - Create migration script in cmd/server/main.go
    - Auto-migrate User, Event, Notification, and audit_logs tables
    - _Requirements: 11.2_

- [ ] 4. Neo4j Repository and Relationship Engine
  - [x] 4.1 Implement Neo4j member repository
    - Create internal/repository/member_repo.go
    - Implement Create, GetByID, GetAll, Update, SoftDelete methods
    - Write Cypher queries for member CRUD operations
    - Create indexes for member.id and member.name
    - _Requirements: 2.1, 2.2, 2.3_

  - [x] 4.2 Write property tests for member repository
    - **Property 7: Member Creation and Retrieval**
    - **Property 8: Member Update Persistence**
    - **Property 9: Soft Delete Preservation**
    - **Validates: Requirements 2.1, 2.2, 2.3, 11.1**

  - [x] 4.3 Implement Neo4j relationship repository
    - Create internal/repository/relationship_repo.go
    - Implement Create, GetAll, Delete, FindPath methods
    - Write Cypher query for shortest path finding between two members
    - _Requirements: 2.4, 2.5, 4.1_

  - [x] 4.4 Write property tests for relationship repository
    - **Property 10: Relationship Creation and Retrieval**
    - **Property 11: Relationship Deletion**
    - **Property 12: Relationship Semantic Validation**
    - **Property 13: Soft Delete Enforcement for Connected Members**
    - **Validates: Requirements 2.4, 2.5, 2.6, 2.7**

  - [x] 4.5 Implement Relationship Engine with kinship mapping
    - Create internal/service/relationship_service.go
    - Implement FindRelationship(fromID, toID) method
    - Implement kinship mapping rules for Indian family relationships
    - Map edge sequences to kinship terms (Father/Nanna, Mother/Amma, Uncle/Babai, Aunt/Attha, etc.)
    - Generate natural language descriptions for relationships
    - _Requirements: 4.2, 4.3, 4.6_

  - [x] 4.6 Write property tests for Relationship Engine
    - **Property 19: Relationship Path Finding**
    - **Property 20: Relationship Result Completeness**
    - **Validates: Requirements 4.1, 4.2, 4.3, 4.4**

- [ ] 5. PostgreSQL Repositories
  - [x] 5.1 Implement user repository
    - Create internal/repository/user_repo.go
    - Implement Create, GetByEmail, GetByID, Update methods
    - _Requirements: 1.1, 1.2_

  - [x] 5.2 Implement event repository
    - Create internal/repository/event_repo.go
    - Implement Create, GetByID, GetAll, Update, Delete methods
    - Implement GetUpcoming and filter methods (by type, member, date range)
    - _Requirements: 5.1, 5.2, 5.3, 5.7_

  - [x] 5.3 Write property tests for event repository
    - **Property 21: Event Creation and Retrieval**
    - **Property 22: Event Update Persistence**
    - **Property 23: Event Deletion**
    - **Property 24: Event Type Validity**
    - **Property 26: Event Filtering**
    - **Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.7, 8.3, 8.4, 11.2**

  - [x] 5.4 Implement notification repository
    - Create internal/repository/notification_repo.go
    - Implement Create, GetPending, UpdateStatus, IncrementRetry methods
    - _Requirements: 6.1, 6.5, 6.6_

  - [ ] 5.5 Write property tests for notification repository
    - **Property 27: Notification Scheduling**
    - **Property 29: Notification Status Update on Success**
    - **Property 30: Notification Retry Logic**
    - **Validates: Requirements 6.1, 6.5, 6.6**

- [ ] 6. Authentication Service and Handlers
  - [x] 6.1 Implement JWT utilities
    - Create internal/utils/jwt.go
    - Implement GenerateAccessToken, GenerateRefreshToken, ValidateToken functions
    - Use JWT_SECRET from environment
    - Set access token expiry to 15 minutes, refresh token to 7 days
    - _Requirements: 1.2, 1.3, 1.7_

  - [x] 6.2 Implement authentication service
    - Create internal/service/auth_service.go
    - Implement Register (hash password with bcrypt)
    - Implement Login (validate credentials, generate tokens)
    - Implement RefreshToken
    - Store refresh tokens in Redis
    - _Requirements: 1.1, 1.2, 1.7_

  - [ ] 6.3 Write property tests for authentication
    - **Property 1: User Registration Round Trip**
    - **Property 2: JWT Authentication Round Trip**
    - **Property 3: Invalid Token Rejection**
    - **Property 6: Token Refresh Round Trip**
    - **Validates: Requirements 1.1, 1.2, 1.3, 1.4, 1.7**

  - [x] 6.4 Implement JWT authentication middleware
    - Create internal/middleware/auth.go
    - Extract and validate JWT from Authorization header
    - Store user ID and role in Fiber context locals
    - Return 401 for missing/invalid tokens
    - _Requirements: 1.3, 1.4_

  - [x] 6.5 Implement role-based authorization middleware
    - Create RequireRole middleware in internal/middleware/auth.go
    - Check user role against allowed roles
    - Return 403 for insufficient permissions
    - _Requirements: 1.5, 1.6_

  - [ ] 6.6 Write property tests for authorization
    - **Property 4: User Role Invariant**
    - **Property 5: Role-Based Authorization**
    - **Validates: Requirements 1.5, 1.6**

  - [x] 6.7 Implement authentication handlers
    - Create internal/handler/auth_handler.go
    - Implement POST /api/auth/register endpoint
    - Implement POST /api/auth/login endpoint
    - Implement POST /api/auth/refresh endpoint
    - Implement GET /api/auth/profile endpoint
    - Return consistent APIResponse format
    - _Requirements: 1.1, 1.2, 1.7, 13.1, 13.2_

- [ ] 7. Family Tree Handlers (Members and Relationships)
  - [x] 7.1 Implement member service with caching
    - Create internal/service/member_service.go
    - Implement Create, GetByID, GetAll, Update, SoftDelete with Redis caching
    - Implement Search by name (case-insensitive partial match)
    - Cache keys: member:{id}, search:members:{query}
    - Invalidate family_tree:* cache on modifications
    - _Requirements: 2.1, 2.2, 2.3, 8.1, 17.1, 17.3, 17.4_

  - [ ] 7.2 Write property tests for member service caching
    - **Property 44: Cache Read-Through**
    - **Property 45: Cache Write-Through**
    - **Property 46: Cache Invalidation on Modification**
    - **Validates: Requirements 17.1, 17.3, 17.4**

  - [x] 7.3 Implement member handlers
    - Create internal/handler/member_handler.go
    - Implement GET /api/members (with pagination)
    - Implement POST /api/members (owner/admin only)
    - Implement GET /api/members/:id
    - Implement PUT /api/members/:id (owner/admin only)
    - Implement DELETE /api/members/:id (owner/admin only)
    - Implement GET /api/members/search?q=name
    - _Requirements: 2.1, 2.2, 2.3, 8.1_

  - [ ] 7.4 Write property tests for member search
    - **Property 33: Member Search**
    - **Property 34: Member Filtering**
    - **Validates: Requirements 8.1, 8.2**

  - [x] 7.5 Implement relationship handlers
    - Create internal/handler/relationship_handler.go
    - Implement GET /api/relationships
    - Implement POST /api/relationships (owner/admin only)
    - Implement DELETE /api/relationships/:id (owner/admin only)
    - Implement GET /api/relationships/path?from=uuid1&to=uuid2
    - _Requirements: 2.4, 2.5, 4.1_

  - [x] 7.6 Implement family tree builder with layout algorithm
    - Create internal/service/family_tree_service.go
    - Implement GetFamilyTree method that queries all members and relationships
    - Calculate node positions using generation-based Y-axis and sibling-order X-axis
    - Transform to React Flow format (nodes with id, type, position, data; edges with id, source, target, type, style)
    - Apply edge color mapping (SPOUSE_OF=rose, PARENT_OF=teal, SIBLING_OF=amber)
    - Check for upcoming events (within 7 days) and set hasUpcomingEvent flag
    - Cache result with key family_tree:{userId}, TTL 5 minutes
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.7, 17.5_

  - [ ] 7.7 Write property tests for family tree builder
    - **Property 14: Family Tree Format Validity**
    - **Property 15: Member Node Rendering Completeness**
    - **Property 16: Relationship Edge Color Mapping**
    - **Property 17: Upcoming Event Indicator**
    - **Validates: Requirements 3.1, 3.2, 3.3, 3.4**

  - [x] 7.8 Implement family tree handler
    - Create internal/handler/family_handler.go
    - Implement GET /api/family/tree endpoint
    - Return nodes and edges in React Flow format
    - _Requirements: 3.1_

- [ ] 8. Event Handlers and Notification Scheduler
  - [x] 8.1 Implement event service
    - Create internal/service/event_service.go
    - Implement Create, GetByID, GetAll, Update, Delete methods
    - Implement GetUpcoming (events within configurable days)
    - Implement filter methods (by type, member, date range)
    - Invalidate events:upcoming:* cache on modifications
    - _Requirements: 5.1, 5.2, 5.3, 5.7_

  - [ ] 8.2 Write property tests for event service
    - **Property 25: Event Countdown Display**
    - **Validates: Requirements 5.6**

  - [x] 8.3 Implement event handlers
    - Create internal/handler/event_handler.go
    - Implement GET /api/events (with filters)
    - Implement POST /api/events (owner/admin only)
    - Implement GET /api/events/:id
    - Implement PUT /api/events/:id (owner/admin only)
    - Implement DELETE /api/events/:id (owner/admin only)
    - Implement GET /api/events/upcoming
    - _Requirements: 5.1, 5.2, 5.3, 5.7_

  - [x] 8.4 Implement notification service
    - Create internal/service/notification_service.go
    - Implement CreateNotifications (called when event is created)
    - Implement Dispatch method for WhatsApp, SMS, Email
    - Integrate Twilio for SMS and WhatsApp
    - Integrate SendGrid for Email
    - Use notification templates from design document
    - Implement retry logic with exponential backoff (max 3 retries)
    - _Requirements: 6.1, 6.3, 6.6, 6.7, 6.8, 6.9_

  - [ ] 8.5 Write property tests for notification service
    - **Property 28: Notification Dispatch**
    - **Validates: Requirements 6.3**

  - [x] 8.6 Implement notification scheduler
    - Create internal/scheduler/notification_scheduler.go
    - Implement Start method that runs every hour
    - Query pending notifications with scheduledAt <= now
    - Use goroutine pool with max 10 concurrent workers
    - Call notification service Dispatch for each notification
    - Update notification status (sent/failed) and retry count
    - _Requirements: 6.2, 6.4, 6.5, 6.6_

- [ ] 9. WebSocket Hub for Real-Time Updates
  - [x] 9.1 Implement WebSocket hub
    - Create internal/handler/websocket_handler.go
    - Implement Hub struct with clients map, broadcast/register/unregister channels
    - Implement Run method as goroutine
    - Implement BroadcastUpdate method
    - _Requirements: 7.1, 7.2, 7.5_

  - [x] 9.2 Implement WebSocket connection handler
    - Implement HandleWebSocket function
    - Upgrade HTTP connection to WebSocket
    - Extract user ID from JWT token
    - Register client with hub
    - Implement readPump and writePump goroutines
    - Send ping messages every 54 seconds to keep connection alive
    - _Requirements: 7.1, 7.4_

  - [ ] 9.3 Write property tests for WebSocket
    - **Property 31: WebSocket Connection Establishment**
    - **Property 32: WebSocket Broadcast**
    - **Validates: Requirements 7.1, 7.2**

  - [x] 9.3 Integrate WebSocket broadcasts in services
    - Update member service to broadcast member_created, member_updated, member_deleted
    - Update relationship service to broadcast relationship_created, relationship_deleted
    - Update event service to broadcast event_created, event_updated, event_deleted
    - _Requirements: 7.2, 7.3_

- [ ] 10. Backend Wiring and Error Handling
  - [x] 10.1 Implement API response utilities
    - Create internal/utils/response.go
    - Define APIResponse struct with success, data, error fields
    - Implement helper functions for success and error responses
    - _Requirements: 13.1, 13.2, 13.3_

  - [ ] 10.2 Write property tests for API response format
    - **Property 37: API Response Format Consistency**
    - **Property 38: Success Response Format**
    - **Property 39: Error Response Format**
    - **Property 40: HTTP Status Code Mapping**
    - **Validates: Requirements 13.1, 13.2, 13.3, 13.4**

  - [x] 10.3 Implement validation utilities
    - Create internal/utils/validator.go
    - Implement validation functions for email, date, phone, required fields
    - _Requirements: 14.2, 14.3_

  - [ ] 10.4 Write property tests for validation
    - **Property 41: Client-Side Validation** (backend validation part)
    - **Property 42: Server-Side Validation**
    - **Validates: Requirements 14.1, 14.2**

  - [x] 10.5 Implement error handling middleware
    - Create internal/middleware/error.go
    - Define AppError struct with code, message, error
    - Implement ErrorHandler middleware
    - Map errors to appropriate HTTP status codes
    - Return user-friendly error messages
    - _Requirements: 14.4_

  - [ ] 10.6 Write property tests for error handling
    - **Property 43: Database Constraint Error Handling**
    - **Validates: Requirements 14.4**

  - [x] 10.7 Implement CORS and logging middleware
    - Create internal/middleware/cors.go
    - Create internal/middleware/logger.go
    - Configure CORS for frontend origin
    - Log all requests with method, path, status, duration
    - _Requirements: 20.1_

  - [ ] 10.8 Implement transaction rollback in services
    - Update member service to use Neo4j transactions
    - Update event service to use GORM transactions
    - Ensure rollback on errors
    - _Requirements: 11.4_

  - [ ] 10.9 Write property tests for transactions
    - **Property 35: Database Transaction Rollback**
    - **Validates: Requirements 11.4**

  - [ ] 10.10 Implement audit logging
    - Create audit_logs table migration
    - Create internal/service/audit_service.go
    - Log all create, update, delete operations with user ID, entity type, entity ID, timestamp
    - _Requirements: 11.5_

  - [ ] 10.11 Write property tests for audit logging
    - **Property 36: Audit Logging**
    - **Validates: Requirements 11.5**

  - [x] 10.12 Wire everything together in main.go
    - Create cmd/server/main.go
    - Initialize config, database clients, Redis client
    - Initialize repositories, services, handlers
    - Set up Fiber app with middleware (CORS, logger, error handler, auth)
    - Register all routes (auth, members, relationships, family, events, websocket)
    - Start notification scheduler as goroutine
    - Start WebSocket hub as goroutine
    - Implement /health endpoint with database health checks
    - Start server on port 8080
    - _Requirements: 12.1, 16.3_

- [ ] 11. Checkpoint - Backend Complete
  - Ensure all backend tests pass
  - Verify all API endpoints are accessible
  - Test WebSocket connection
  - Ask the user if questions arise

- [ ] 12. Frontend Scaffold and Configuration
  - [x] 12.1 Initialize React + TypeScript + Vite project
    - Run `npm create vite@latest frontend -- --template react-ts`
    - Install dependencies: react-router-dom, axios, zustand, @tanstack/react-query, reactflow, framer-motion, react-hook-form
    - Install Tailwind CSS and configure with custom theme
    - _Requirements: 10.1, 10.2_

  - [x] 12.2 Configure Tailwind with VamsaSetu theme
    - Update tailwind.config.js with custom colors (saffron, turmeric, ivory, teal, charcoal)
    - Add custom fonts (Playfair Display for headings, DM Sans for body)
    - Configure responsive breakpoints (mobile <768px, tablet 768-1024px, desktop >1024px)
    - _Requirements: 9.4, 10.1, 10.2_

  - [x] 12.3 Set up project structure
    - Create directory structure: src/{components/{auth,family,events,common,ui},pages,hooks,services,stores,types,utils}
    - Create src/App.tsx with router setup
    - Create src/main.tsx with React Query and Zustand providers
    - _Requirements: 20.1_

- [ ] 13. Frontend Type Definitions and Services
  - [x] 13.1 Define TypeScript types
    - Create src/types/user.ts (User, LoginRequest, RegisterRequest)
    - Create src/types/member.ts (Member, MemberNodeData, CreateMemberRequest)
    - Create src/types/relationship.ts (Relationship, RelationshipPath, CreateRelationshipRequest)
    - Create src/types/event.ts (Event, CreateEventRequest, EventType)
    - Create src/types/api.ts (APIResponse<T>)
    - _Requirements: 20.3_

  - [x] 13.2 Implement API service layer
    - Create src/services/api.ts with Axios instance
    - Configure base URL from environment variable
    - Add request interceptor to attach JWT token
    - Add response interceptor for error handling
    - _Requirements: 13.1_

  - [x] 13.3 Implement auth service
    - Create src/services/authService.ts
    - Implement register, login, refresh, getProfile methods
    - _Requirements: 1.1, 1.2, 1.7_

  - [x] 13.4 Implement member service
    - Create src/services/memberService.ts
    - Implement getAll, getById, create, update, delete, search methods
    - _Requirements: 2.1, 2.2, 2.3, 8.1_

  - [x] 13.5 Implement relationship service
    - Create src/services/relationshipService.ts
    - Implement getAll, create, delete, findPath methods
    - _Requirements: 2.4, 2.5, 4.1_

  - [x] 13.6 Implement event service
    - Create src/services/eventService.ts
    - Implement getAll, getById, create, update, delete, getUpcoming methods
    - _Requirements: 5.1, 5.2, 5.3, 5.7_

  - [x] 13.7 Implement family tree service
    - Create src/services/familyTreeService.ts
    - Implement getTree method
    - _Requirements: 3.1_

- [ ] 14. Frontend State Management
  - [x] 14.1 Implement Zustand auth store
    - Create src/stores/authStore.ts
    - Define AuthState with user, accessToken, setAuth, logout
    - Persist to localStorage using Zustand persist middleware
    - _Requirements: 18.2, 18.4_

  - [ ] 14.2 Write property tests for auth store persistence
    - **Property 48: User Preference Persistence**
    - **Validates: Requirements 18.4**

  - [x] 14.3 Implement Zustand UI store
    - Create src/stores/uiStore.ts
    - Define UIState with sidebar, modals, selected items
    - _Requirements: 18.2_

  - [x] 14.4 Implement React Query hooks
    - Create src/hooks/useAuth.ts (useLogin, useRegister, useProfile)
    - Create src/hooks/useMembers.ts (useMembers, useCreateMember, useUpdateMember, useDeleteMember, useSearchMembers)
    - Create src/hooks/useRelationships.ts (useRelationships, useCreateRelationship, useDeleteRelationship, useFindPath)
    - Create src/hooks/useEvents.ts (useEvents, useCreateEvent, useUpdateEvent, useDeleteEvent, useUpcomingEvents)
    - Create src/hooks/useFamilyTree.ts (useFamilyTree)
    - Configure cache invalidation on mutations
    - _Requirements: 18.1, 18.3_

  - [ ] 14.5 Write property tests for React Query cache invalidation
    - **Property 47: React Query Cache Invalidation**
    - **Validates: Requirements 18.3**

  - [x] 14.6 Implement WebSocket hook
    - Create src/hooks/useWebSocket.ts
    - Connect to WebSocket endpoint with JWT token
    - Listen for messages and invalidate React Query cache
    - Implement auto-reconnect on disconnect
    - _Requirements: 7.1, 7.3, 7.4_

- [ ] 15. Authentication Pages
  - [x] 15.1 Create UI components
    - Create src/components/ui/Button.tsx with variants (primary, secondary, outline)
    - Create src/components/ui/Input.tsx with label, error display, validation
    - Create src/components/ui/Modal.tsx with Framer Motion animations
    - Create src/components/ui/Card.tsx
    - Apply Tailwind styling with VamsaSetu theme colors
    - _Requirements: 10.1, 10.5, 10.6_

  - [x] 15.2 Create LoginForm component
    - Create src/components/auth/LoginForm.tsx
    - Use react-hook-form for form handling
    - Implement client-side validation (required fields, email format)
    - Display inline error messages
    - Call useLogin hook on submit
    - _Requirements: 1.2, 14.1_

  - [ ] 15.3 Write property tests for login form validation
    - **Property 41: Client-Side Validation** (frontend part)
    - **Validates: Requirements 14.1**

  - [x] 15.4 Create RegisterForm component
    - Create src/components/auth/RegisterForm.tsx
    - Use react-hook-form for form handling
    - Implement client-side validation (required fields, email format, password strength)
    - Display inline error messages
    - Call useRegister hook on submit
    - _Requirements: 1.1, 14.1_

  - [x] 15.5 Create LoginPage and RegisterPage
    - Create src/pages/LoginPage.tsx with LoginForm
    - Create src/pages/RegisterPage.tsx with RegisterForm
    - Apply VamsaSetu visual design (rangoli background, saffron accents)
    - Add Framer Motion page transitions
    - _Requirements: 10.1, 10.5, 10.6_

- [ ] 16. Common Layout Components
  - [x] 16.1 Create Navbar component
    - Create src/components/common/Navbar.tsx
    - Display logo, user name, logout button
    - Apply VamsaSetu theme styling
    - _Requirements: 10.1_

  - [x] 16.2 Create Sidebar component (desktop)
    - Create src/components/common/Sidebar.tsx
    - Display navigation links (Dashboard, Family Tree, Events, Settings)
    - Highlight active route
    - Apply VamsaSetu theme styling
    - _Requirements: 10.1_

  - [x] 16.3 Create BottomNav component (mobile)
    - Create src/components/common/BottomNav.tsx
    - Display bottom tab navigation with icons
    - Show only on mobile (<768px)
    - Ensure minimum touch target size of 44x44 pixels
    - _Requirements: 9.1, 9.5_

  - [x] 16.4 Create LoadingSpinner component
    - Create src/components/common/LoadingSpinner.tsx
    - Use Framer Motion for animation
    - _Requirements: 19.5_

  - [ ] 16.5 Write property tests for loading states
    - **Property 50: Loading State Display**
    - **Validates: Requirements 19.5**

  - [x] 16.6 Create ErrorBoundary component
    - Create src/components/common/ErrorBoundary.tsx
    - Catch React errors and display user-friendly message
    - Provide reload button
    - _Requirements: 19.6_

  - [ ] 16.7 Write property tests for error states
    - **Property 51: Error State Display**
    - **Validates: Requirements 19.6**

  - [x] 16.8 Create EmptyState component
    - Create src/components/common/EmptyState.tsx
    - Display meaningful illustrations when no data
    - _Requirements: 14.5_

- [ ] 17. Custom React Flow Nodes and Edges
  - [x] 17.1 Create MemberNode component
    - Create src/components/family/MemberNode.tsx
    - Display avatar, name, relation badge, event indicator
    - Apply gender-based border colors (blue for male, pink for female)
    - Show glowing amber indicator for upcoming events
    - Add hover glow effect
    - Implement click handler to open member details panel
    - _Requirements: 3.2, 3.4, 3.6, 10.6_

  - [ ] 17.2 Write property tests for MemberNode
    - **Property 18: Member Node Click Interaction**
    - **Validates: Requirements 3.6**

  - [x] 17.3 Create RelationshipEdge component
    - Create src/components/family/RelationshipEdge.tsx
    - Apply color-coded bezier edges (rose for SPOUSE_OF, teal for PARENT_OF, amber for SIBLING_OF)
    - _Requirements: 3.3_

  - [x] 17.4 Create RangoliPattern component
    - Create src/components/common/RangoliPattern.tsx
    - Implement SVG rangoli-inspired geometric pattern
    - Use as background texture with low opacity
    - _Requirements: 10.3_

- [ ] 18. Family Tree Page
  - [x] 18.1 Create TreeCanvas component
    - Create src/components/family/TreeCanvas.tsx
    - Use React Flow with custom node types (memberNode)
    - Fetch family tree data with useFamilyTree hook
    - Render nodes and edges from API response
    - Add Background, Controls, MiniMap components
    - Enable zoom, pan, fit-view controls
    - Support pinch-to-zoom and pan gestures on mobile
    - _Requirements: 3.1, 3.5, 9.2_

  - [x] 18.2 Create TreeControls component
    - Create src/components/family/TreeControls.tsx
    - Add buttons for zoom in, zoom out, fit view, add member, add relationship
    - _Requirements: 3.5_

  - [x] 18.3 Create MemberDetailsPanel component
    - Create src/components/family/MemberDetailsPanel.tsx
    - Display full member details (name, DOB, gender, email, phone, avatar)
    - Show direct relationships
    - Add edit and delete buttons (owner/admin only)
    - Slide in from right on member node click
    - _Requirements: 3.6_

  - [x] 18.4 Create AddMemberModal component
    - Create src/components/family/AddMemberModal.tsx
    - Form with fields: name, dateOfBirth, gender, email, phone, avatarUrl
    - Use react-hook-form with validation
    - Call useCreateMember hook on submit
    - _Requirements: 2.1_

  - [x] 18.5 Create AddRelationshipModal component
    - Create src/components/family/AddRelationshipModal.tsx
    - Form with fields: fromMember (dropdown), toMember (dropdown), type (SPOUSE_OF, PARENT_OF, SIBLING_OF)
    - Call useCreateRelationship hook on submit
    - _Requirements: 2.4_

  - [x] 18.6 Create FamilyTreePage
    - Create src/pages/FamilyTreePage.tsx
    - Compose TreeCanvas, TreeControls, MemberDetailsPanel, AddMemberModal, AddRelationshipModal
    - Implement relationship path highlighting with animated traveling dot
    - _Requirements: 3.1, 4.5, 10.7_

- [ ] 19. Dashboard Page
  - [x] 19.1 Create DashboardPage
    - Create src/pages/DashboardPage.tsx
    - Display welcome message with user name
    - Show upcoming events widget (next 3 events)
    - Show family tree summary (total members, recent additions)
    - Apply VamsaSetu theme styling with card layout
    - _Requirements: 10.1_

- [ ] 20. Members Page
  - [x] 20.1 Create MemberList component
    - Create src/components/family/MemberList.tsx
    - Display all members in a grid or list view
    - Show avatar, name, relation badge
    - Add search bar for filtering by name
    - _Requirements: 8.1_

  - [x] 20.2 Create MembersPage
    - Create src/pages/MembersPage.tsx
    - Compose MemberList, AddMemberModal
    - Add "Add Member" button
    - _Requirements: 2.1, 8.1_

- [ ] 21. Events Page
  - [x] 21.1 Create EventCard component
    - Create src/components/events/EventCard.tsx
    - Display event title, date, type, associated members
    - Show countdown chip for events within 7 days
    - Add edit and delete buttons (owner/admin only)
    - _Requirements: 5.6_

  - [x] 21.2 Create EventList component
    - Create src/components/events/EventList.tsx
    - Display events in list view with EventCard components
    - Add filters for type (birthday, anniversary, ceremony, custom) and member
    - _Requirements: 5.5, 5.7_

  - [ ] 21.3 Create EventCalendar component
    - Create src/components/events/EventCalendar.tsx
    - Display events in calendar view (use a calendar library like react-big-calendar)
    - Highlight dates with events
    - _Requirements: 5.5_

  - [x] 21.4 Create AddEventModal component
    - Create src/components/events/AddEventModal.tsx
    - Form with fields: title, description, eventDate, eventType, memberIds (multi-select)
    - Use react-hook-form with validation
    - Call useCreateEvent hook on submit
    - _Requirements: 5.1_

  - [x] 21.5 Create EventsPage
    - Create src/pages/EventsPage.tsx
    - Compose EventCalendar, EventList, AddEventModal
    - Add toggle between calendar and list view
    - Add "Add Event" button
    - _Requirements: 5.1, 5.5_

- [ ] 22. Relationship Finder Page
  - [x] 22.1 Create RelationshipFinderPage
    - Create src/pages/RelationshipFinderPage.tsx
    - Add two dropdowns to select members (from and to)
    - Add "Find Relationship" button
    - Call useFindPath hook on submit
    - Display relationship result (path, relation label, kinship term, description)
    - Highlight path on family tree canvas
    - _Requirements: 4.1, 4.2, 4.3, 4.5_

- [ ] 23. Settings Page
  - [x] 23.1 Create SettingsPage
    - Create src/pages/SettingsPage.tsx
    - Add sections for profile settings, notification preferences, theme settings
    - Allow user to update email, name, password
    - Allow user to configure notification channels (WhatsApp, SMS, Email)
    - _Requirements: 18.4_

- [ ] 24. Mobile Responsive Pass
  - [x] 24.1 Implement responsive breakpoints
    - Create src/hooks/useResponsive.ts hook
    - Detect breakpoint (mobile, tablet, desktop)
    - _Requirements: 9.4_

  - [x] 24.2 Apply mobile-specific layouts
    - Show BottomNav on mobile, Sidebar on desktop
    - Adjust TreeCanvas for mobile (smaller nodes, simplified controls)
    - Make MemberDetailsPanel full-screen on mobile with swipeable cards
    - Ensure all modals are mobile-friendly
    - _Requirements: 9.1, 9.3_

  - [x] 24.3 Test touch interactions
    - Verify pinch-to-zoom and pan gestures on TreeCanvas
    - Ensure minimum touch target size of 44x44 pixels for all interactive elements
    - _Requirements: 9.2, 9.5_

- [ ] 25. Accessibility Pass
  - [x] 25.1 Add ARIA labels to all interactive elements
    - Add aria-label to icon buttons
    - Add aria-label to form inputs
    - Add aria-label to navigation links
    - _Requirements: 19.2_

  - [ ] 25.2 Write property tests for ARIA labels
    - **Property 49: ARIA Label Presence**
    - **Validates: Requirements 19.2**

  - [x] 25.3 Ensure keyboard accessibility
    - Test tab order throughout application
    - Add visible focus indicators
    - Ensure all interactive elements are keyboard accessible
    - _Requirements: 19.1, 19.3_

  - [x] 25.4 Verify color contrast
    - Check color contrast ratios meet WCAG AA standards (4.5:1 for normal text, 3:1 for large text)
    - Adjust colors if needed
    - _Requirements: 19.4_

- [ ] 26. Seed Data and Demo Setup
  - [x] 26.1 Create seed data script
    - Create backend/cmd/seed/main.go
    - Create demo user (demo@vamsasetu.com / Demo@1234)
    - Create 12 members across 3 generations
    - Create relationships connecting all members
    - Create 5 events (2 birthdays, 1 anniversary, 1 ceremony, 1 custom)
    - _Requirements: 15.1, 15.2, 15.3, 15.4_

- [ ] 27. Documentation
  - [x] 27.1 Write comprehensive README.md
    - Document system architecture with diagram
    - Document tech stack
    - Document quick start instructions (Docker Compose)
    - Document local development setup
    - Document API endpoints with examples
    - Document environment variables
    - Document testing instructions
    - Document deployment instructions
    - Document project structure
    - Document roadmap and future features
    - _Requirements: 20.1, 20.2, 20.5_

  - [x] 27.2 Add inline code comments
    - Add comments to complex algorithms (Relationship Engine, layout algorithm)
    - Add comments to all public functions and methods
    - _Requirements: 20.2_

  - [x] 27.3 Define PropTypes for React components
    - Add TypeScript interfaces for all component props
    - _Requirements: 20.3_

  - [x] 27.4 Create constants file
    - Create src/utils/constants.ts
    - Move all hardcoded strings to constants
    - _Requirements: 20.4_

- [ ] 28. Final Integration and Testing
  - [ ] 28.1 Run all backend tests
    - Run unit tests: `go test ./... -v`
    - Run property tests: `go test ./... -v -tags=property`
    - Ensure minimum 100 iterations for each property test
    - Verify 80% line coverage
    - _Requirements: All_

  - [ ] 28.2 Run all frontend tests
    - Run unit tests: `npm test`
    - Run property tests with fast-check (minimum 100 iterations)
    - Verify coverage
    - _Requirements: All_

  - [ ] 28.3 End-to-end testing
    - Test user registration and login flow
    - Test creating family tree with members and relationships
    - Test adding events and verifying notifications
    - Test WebSocket real-time updates
    - Test mobile responsive design
    - _Requirements: All_

  - [ ] 28.4 Performance testing
    - Test family tree rendering with 100+ members
    - Test cache hit rates
    - Verify API response times (<500ms for cached queries)
    - _Requirements: 17.1, 17.2, 17.5_

  - [ ] 28.5 Security testing
    - Test JWT authentication and authorization
    - Test role-based access control
    - Test input validation and SQL injection prevention
    - Test CORS configuration
    - _Requirements: 1.3, 1.4, 1.5, 1.6, 14.2_

- [ ] 29. Final Checkpoint - System Complete
  - Ensure all tests pass
  - Verify all features are working end-to-end
  - Verify Docker Compose brings up all services successfully
  - Verify seed data populates correctly
  - Ask the user if questions arise

## Notes

- Tasks marked with `*` are optional property-based testing tasks and can be skipped for faster MVP delivery
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation at key milestones
- Property tests validate universal correctness properties with minimum 100 iterations
- Unit tests validate specific examples and edge cases
- All 51 correctness properties from the design document are covered in property test tasks
- The implementation follows the specified build order for incremental progress
- Backend uses Go + Fiber, Frontend uses React + TypeScript
- Databases: Neo4j (graph), PostgreSQL (relational), Redis (cache)
- External services: Twilio (SMS/WhatsApp), SendGrid (Email)
