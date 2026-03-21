# Requirements Document

## Introduction

VamsaSetu ("Vamsa" = lineage, "Setu" = bridge) is a comprehensive family tree and event management system designed specifically for Indian families. The system enables users to visualize their family tree as an interactive graph, understand complex Indian kinship relationships, track family events (birthdays, anniversaries, ceremonies), and receive smart notifications via WhatsApp, Email, and SMS. The system uses a graph database to model complex family relationships and provides an intuitive visual interface with culturally-appropriate design elements.

## Glossary

- **VamsaSetu_System**: The complete family tree and event management application
- **Family_Tree**: A graph structure representing family members as nodes and relationships as edges
- **Member**: An individual person in a family tree with attributes like name, date of birth, gender, contact information
- **Relationship**: A connection between two members (spouse, parent-child, sibling)
- **Relationship_Engine**: The core algorithm that computes relationship paths and generates kinship labels
- **Event**: A calendar item associated with one or more family members (birthday, anniversary, ceremony)
- **Notification_Service**: The subsystem responsible for sending reminders via WhatsApp, SMS, and Email
- **Tree_Canvas**: The interactive React Flow visualization component displaying the family tree
- **User**: An authenticated person who owns or views family trees
- **Owner**: A user role with full edit permissions on a family tree
- **Viewer**: A user role with read-only access to a family tree
- **Admin**: A user role with system-wide administrative privileges
- **Kinship_Term**: Indian relationship labels (Amma, Nanna, Akka, Anna, Tammudu, Chelli, Attha, Mamayya, Pinni, Babai, Menarikam, etc.)
- **Neo4j_Database**: The graph database storing members and relationships
- **PostgreSQL_Database**: The relational database storing users, events, and audit logs
- **Redis_Cache**: The caching layer for performance optimization
- **WebSocket_Hub**: The real-time communication channel for live updates
- **Notification_Scheduler**: The background process that checks and dispatches scheduled notifications
- **JWT_Token**: JSON Web Token used for authentication
- **React_Flow**: The library used for rendering the interactive family tree graph
- **Fiber_Framework**: The Go HTTP framework for the backend API
- **GORM**: The Go ORM library for PostgreSQL interactions
- **Twilio_Service**: The third-party service for SMS and WhatsApp messaging
- **SendGrid_Service**: The third-party service for email delivery

## Requirements

### Requirement 1: User Authentication and Authorization

**User Story:** As a user, I want to securely register and log in to the system, so that I can access my family trees and protect my family data.

#### Acceptance Criteria

1. WHEN a user submits valid registration data (email, password, name), THE VamsaSetu_System SHALL create a new user account with a hashed password
2. WHEN a user submits valid login credentials, THE VamsaSetu_System SHALL return a JWT_Token valid for the configured expiry period
3. WHEN a user provides a valid JWT_Token, THE VamsaSetu_System SHALL authenticate the user and grant access to protected endpoints
4. WHEN a user provides an expired or invalid JWT_Token, THE VamsaSetu_System SHALL return an authentication error
5. THE VamsaSetu_System SHALL assign one of three roles to each user: Owner, Viewer, or Admin
6. WHEN a Viewer attempts to modify family tree data, THE VamsaSetu_System SHALL deny the request with an authorization error
7. WHEN a user requests a token refresh with a valid refresh token, THE VamsaSetu_System SHALL issue a new JWT_Token

### Requirement 2: Family Tree Data Management

**User Story:** As an Owner, I want to add, edit, and remove family members and their relationships, so that I can build an accurate representation of my family structure.

#### Acceptance Criteria

1. WHEN an Owner creates a Member with required attributes (name, date of birth, gender), THE VamsaSetu_System SHALL store the Member in the Neo4j_Database
2. WHEN an Owner updates a Member's attributes, THE VamsaSetu_System SHALL persist the changes to the Neo4j_Database
3. WHEN an Owner deletes a Member, THE VamsaSetu_System SHALL perform a soft delete and preserve the Member's data for audit purposes
4. WHEN an Owner creates a Relationship between two Members, THE VamsaSetu_System SHALL store the Relationship as an edge in the Neo4j_Database with the appropriate type (SPOUSE_OF, PARENT_OF, SIBLING_OF)
5. WHEN an Owner deletes a Relationship, THE VamsaSetu_System SHALL remove the edge from the Neo4j_Database
6. THE VamsaSetu_System SHALL validate that relationship types are semantically correct (e.g., a person cannot be their own parent)
7. WHEN a Member has associated Relationships, THE VamsaSetu_System SHALL prevent hard deletion and require soft delete

### Requirement 3: Interactive Family Tree Visualization

**User Story:** As a user, I want to view my family tree as an interactive graph with custom-styled nodes and edges, so that I can visually understand my family structure.

#### Acceptance Criteria

1. WHEN a user requests the family tree, THE VamsaSetu_System SHALL return nodes and edges in a format compatible with React_Flow
2. THE Tree_Canvas SHALL render each Member as a custom node displaying avatar, name, relation badge, and event indicator
3. THE Tree_Canvas SHALL render Relationships as curved bezier edges color-coded by type (spouse=rose, parent-child=teal, sibling=amber)
4. WHEN a Member has an Event within 7 days, THE Tree_Canvas SHALL display a glowing amber indicator on the Member's node
5. THE Tree_Canvas SHALL provide zoom, pan, fit-view, and minimap controls
6. WHEN a user clicks a Member node, THE Tree_Canvas SHALL display a side panel with full Member details and direct Relationships
7. THE VamsaSetu_System SHALL calculate node positions using a layered layout algorithm (generation-based Y-axis, sibling-order X-axis)

### Requirement 4: Relationship Path Finding and Kinship Intelligence

**User Story:** As a user, I want to select two family members and discover how they are related, so that I can understand complex family connections using Indian kinship terms.

#### Acceptance Criteria

1. WHEN a user requests the relationship between two Members, THE Relationship_Engine SHALL compute the shortest path in the Neo4j_Database
2. THE Relationship_Engine SHALL map the edge sequence to appropriate Kinship_Terms based on Indian family relationship conventions
3. THE Relationship_Engine SHALL return a result containing the path nodes, relation label, description, and natural language explanation
4. WHEN no path exists between two Members, THE Relationship_Engine SHALL return a "not related" result
5. WHEN a user selects two Members on the Tree_Canvas, THE VamsaSetu_System SHALL highlight the connecting path with an animated traveling dot effect
6. THE Relationship_Engine SHALL handle complex multi-hop relationships (e.g., father's sister = Attha, mother's brother = Mamayya)

### Requirement 5: Event Management and Calendar

**User Story:** As an Owner, I want to create and manage family events (birthdays, anniversaries, ceremonies), so that I can track important dates and receive reminders.

#### Acceptance Criteria

1. WHEN an Owner creates an Event with required attributes (title, date, type, associated Members), THE VamsaSetu_System SHALL store the Event in the PostgreSQL_Database
2. WHEN an Owner updates an Event, THE VamsaSetu_System SHALL persist the changes to the PostgreSQL_Database
3. WHEN an Owner deletes an Event, THE VamsaSetu_System SHALL remove the Event from the PostgreSQL_Database
4. THE VamsaSetu_System SHALL support Event types: birthday, anniversary, ceremony, and custom
5. THE VamsaSetu_System SHALL display Events in both calendar view and list view
6. WHEN an Event is within 7 days, THE VamsaSetu_System SHALL display a countdown chip on the Event card
7. THE VamsaSetu_System SHALL filter Events by upcoming, member, and type

### Requirement 6: Notification Scheduling and Delivery

**User Story:** As a user, I want to receive automated reminders for upcoming family events via WhatsApp, Email, and SMS, so that I never miss important occasions.

#### Acceptance Criteria

1. WHEN an Owner creates an Event, THE VamsaSetu_System SHALL schedule notifications based on user preferences (e.g., 7 days before, 1 day before, day of)
2. THE Notification_Scheduler SHALL run every hour and query the PostgreSQL_Database for due notifications
3. WHEN a notification is due, THE Notification_Service SHALL dispatch the notification via the configured channel (WhatsApp, SMS, Email)
4. THE Notification_Service SHALL use a goroutine pool with a maximum of 10 concurrent dispatches
5. WHEN a notification is successfully sent, THE VamsaSetu_System SHALL mark the notification as sent in the PostgreSQL_Database
6. WHEN a notification fails to send, THE VamsaSetu_System SHALL log the error and retry up to 3 times with exponential backoff
7. THE VamsaSetu_System SHALL send WhatsApp messages via Twilio_Service using approved template messages
8. THE VamsaSetu_System SHALL send SMS messages via Twilio_Service
9. THE VamsaSetu_System SHALL send Email messages via SendGrid_Service using HTML templates

### Requirement 7: Real-Time Updates via WebSocket

**User Story:** As a user, I want to see live updates when other users modify the family tree, so that I always have the most current information.

#### Acceptance Criteria

1. WHEN a user connects to the WebSocket endpoint with a valid JWT_Token, THE WebSocket_Hub SHALL establish a connection
2. WHEN an Owner adds, updates, or deletes a Member or Relationship, THE VamsaSetu_System SHALL broadcast the change to all connected clients viewing the same Family_Tree
3. WHEN a client receives a WebSocket message, THE Tree_Canvas SHALL update the visualization without requiring a page refresh
4. WHEN a WebSocket connection is lost, THE VamsaSetu_System SHALL attempt to reconnect automatically
5. THE WebSocket_Hub SHALL run as a goroutine and manage client connections via channels

### Requirement 8: Search and Filtering

**User Story:** As a user, I want to search for family members by name and filter events by type or member, so that I can quickly find specific information.

#### Acceptance Criteria

1. WHEN a user enters a search query, THE VamsaSetu_System SHALL return Members whose names match the query (case-insensitive partial match)
2. THE VamsaSetu_System SHALL support filtering Members by gender, generation, and relationship type
3. THE VamsaSetu_System SHALL support filtering Events by type (birthday, anniversary, ceremony, custom)
4. THE VamsaSetu_System SHALL support filtering Events by associated Member
5. WHEN a user applies filters, THE VamsaSetu_System SHALL return results within 500 milliseconds using Redis_Cache for frequently accessed queries

### Requirement 9: Responsive Design and Mobile Support

**User Story:** As a mobile user, I want to access VamsaSetu on my phone with an optimized interface, so that I can view and manage my family tree on the go.

#### Acceptance Criteria

1. WHEN a user accesses VamsaSetu on a mobile device, THE VamsaSetu_System SHALL display a bottom tab navigation instead of a sidebar
2. THE Tree_Canvas SHALL support pinch-to-zoom and pan gestures on touch devices
3. WHEN a user views Member details on mobile, THE VamsaSetu_System SHALL display swipeable Member cards
4. THE VamsaSetu_System SHALL use responsive breakpoints: mobile (<768px), tablet (768px-1024px), desktop (>1024px)
5. THE VamsaSetu_System SHALL ensure all interactive elements have a minimum touch target size of 44x44 pixels on mobile

### Requirement 10: Visual Design and Cultural Theming

**User Story:** As a user, I want VamsaSetu to have a visually distinctive Indian-inspired design, so that the application feels culturally appropriate and aesthetically pleasing.

#### Acceptance Criteria

1. THE VamsaSetu_System SHALL use the color palette: deep saffron (#E8650A), turmeric gold (#F5A623), ivory (#FBF5E6), deep teal (#0D4A52), warm charcoal (#2C2420)
2. THE VamsaSetu_System SHALL use "Playfair Display" font for headings and "DM Sans" font for body text
3. THE VamsaSetu_System SHALL display subtle rangoli-inspired geometric patterns as background textures using inline SVG
4. THE Tree_Canvas SHALL render the root node with a mandala-like radial layout
5. THE VamsaSetu_System SHALL use Framer Motion for page transitions with fade and slide effects
6. WHEN a user hovers over a Member node, THE Tree_Canvas SHALL display a glow effect
7. WHEN a Relationship path is highlighted, THE Tree_Canvas SHALL animate the path with a drawing effect
8. THE VamsaSetu_System SHALL default to dark mode with the warm charcoal background

### Requirement 11: Data Persistence and Backup

**User Story:** As an Owner, I want my family tree data to be reliably stored and backed up, so that I don't lose important family information.

#### Acceptance Criteria

1. THE VamsaSetu_System SHALL persist all Member and Relationship data in the Neo4j_Database
2. THE VamsaSetu_System SHALL persist all User, Event, and audit log data in the PostgreSQL_Database
3. THE VamsaSetu_System SHALL use Docker named volumes for Neo4j, PostgreSQL, and Redis data persistence
4. WHEN a database operation fails, THE VamsaSetu_System SHALL roll back the transaction and return an error to the client
5. THE VamsaSetu_System SHALL log all data modification operations (create, update, delete) with timestamp and user ID in the PostgreSQL_Database

### Requirement 12: Configuration and Environment Management

**User Story:** As a system administrator, I want to configure VamsaSetu using environment variables, so that I can deploy the system in different environments without code changes.

#### Acceptance Criteria

1. THE VamsaSetu_System SHALL load configuration from environment variables using godotenv
2. THE VamsaSetu_System SHALL require the following environment variables: POSTGRES_URL, NEO4J_URI, NEO4J_USERNAME, NEO4J_PASSWORD, REDIS_ADDR, JWT_SECRET
3. THE VamsaSetu_System SHALL require the following notification service variables: SENDGRID_API_KEY, TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_PHONE_NUMBER, TWILIO_WHATSAPP_NUMBER
4. THE VamsaSetu_System SHALL provide a .env.example file with all required variables and placeholder values
5. WHEN a required environment variable is missing, THE VamsaSetu_System SHALL fail to start and log a descriptive error message

### Requirement 13: API Response Consistency

**User Story:** As a frontend developer, I want all API responses to follow a consistent format, so that I can handle responses uniformly in the client code.

#### Acceptance Criteria

1. THE VamsaSetu_System SHALL return all API responses in the format: { "success": bool, "data": any, "error": string }
2. WHEN an API request succeeds, THE VamsaSetu_System SHALL set success to true, populate data, and set error to empty string
3. WHEN an API request fails, THE VamsaSetu_System SHALL set success to false, set data to null, and populate error with a descriptive message
4. THE VamsaSetu_System SHALL use appropriate HTTP status codes: 200 (success), 201 (created), 400 (bad request), 401 (unauthorized), 403 (forbidden), 404 (not found), 500 (server error)

### Requirement 14: Error Handling and Validation

**User Story:** As a user, I want to receive clear error messages when I submit invalid data, so that I can correct my input and successfully complete my task.

#### Acceptance Criteria

1. WHEN a user submits a form with invalid data, THE VamsaSetu_System SHALL perform client-side validation and display inline error messages
2. WHEN a user submits a form with invalid data, THE VamsaSetu_System SHALL perform server-side validation and return descriptive error messages
3. THE VamsaSetu_System SHALL validate required fields, data types, format constraints, and business rules
4. WHEN a database constraint is violated, THE VamsaSetu_System SHALL return a user-friendly error message (not a raw database error)
5. THE VamsaSetu_System SHALL display empty states with meaningful illustrations when no data is available

### Requirement 15: Seed Data for Development and Demo

**User Story:** As a developer, I want to populate the system with realistic seed data, so that I can test features and demonstrate the application.

#### Acceptance Criteria

1. THE VamsaSetu_System SHALL provide a seed command that creates a demo user with email demo@vamsasetu.com and password Demo@1234
2. THE VamsaSetu_System SHALL create a Family_Tree with 12 Members across 3 generations
3. THE VamsaSetu_System SHALL create 5 Events: 2 upcoming birthdays, 1 anniversary, 1 ceremony, and 1 custom event
4. THE VamsaSetu_System SHALL create Relationships connecting all Members in a realistic family structure
5. THE VamsaSetu_System SHALL store all seed data in both Neo4j_Database and PostgreSQL_Database as appropriate

### Requirement 16: Docker Compose Orchestration

**User Story:** As a developer, I want to start all VamsaSetu services with a single command, so that I can quickly set up a development environment.

#### Acceptance Criteria

1. THE VamsaSetu_System SHALL provide a docker-compose.yml file that defines all services: frontend, backend, neo4j, postgres, redis
2. THE VamsaSetu_System SHALL configure service dependencies so that backend waits for neo4j and postgres health checks before starting
3. THE VamsaSetu_System SHALL expose the following ports: frontend (3000), backend (8080), neo4j (7474, 7687), postgres (5432), redis (6379)
4. THE VamsaSetu_System SHALL use named volumes for data persistence: neo4j_data, postgres_data, redis_data
5. WHEN a developer runs "docker-compose up", THE VamsaSetu_System SHALL start all services and make the application accessible at http://localhost:3000

### Requirement 17: Performance and Caching

**User Story:** As a user, I want VamsaSetu to respond quickly to my requests, so that I have a smooth and efficient experience.

#### Acceptance Criteria

1. WHEN a user requests frequently accessed data (e.g., family tree structure), THE VamsaSetu_System SHALL check Redis_Cache before querying the database
2. WHEN cached data is found, THE VamsaSetu_System SHALL return the cached result within 50 milliseconds
3. WHEN cached data is not found, THE VamsaSetu_System SHALL query the database, cache the result with an appropriate TTL, and return the data
4. WHEN a Member or Relationship is modified, THE VamsaSetu_System SHALL invalidate the relevant cache entries
5. THE VamsaSetu_System SHALL cache family tree data with a TTL of 5 minutes
6. THE VamsaSetu_System SHALL cache Member search results with a TTL of 2 minutes

### Requirement 18: Frontend State Management

**User Story:** As a frontend developer, I want a clear separation between server state and client state, so that the application is maintainable and predictable.

#### Acceptance Criteria

1. THE VamsaSetu_System SHALL use React Query for managing server state (API data, caching, refetching)
2. THE VamsaSetu_System SHALL use Zustand for managing client state (UI state, user preferences, temporary form data)
3. WHEN an API mutation succeeds, THE VamsaSetu_System SHALL invalidate relevant React Query cache entries to trigger refetch
4. THE VamsaSetu_System SHALL persist user preferences (theme, notification settings) to localStorage via Zustand middleware

### Requirement 19: Accessibility and User Experience

**User Story:** As a user with accessibility needs, I want VamsaSetu to be usable with keyboard navigation and screen readers, so that I can access all features.

#### Acceptance Criteria

1. THE VamsaSetu_System SHALL ensure all interactive elements are keyboard accessible with visible focus indicators
2. THE VamsaSetu_System SHALL provide ARIA labels for all icon buttons and interactive elements
3. THE VamsaSetu_System SHALL maintain a logical tab order throughout the application
4. THE VamsaSetu_System SHALL ensure color contrast ratios meet WCAG AA standards (4.5:1 for normal text, 3:1 for large text)
5. THE VamsaSetu_System SHALL provide loading states for all asynchronous operations
6. THE VamsaSetu_System SHALL display error states with clear recovery actions

### Requirement 20: Documentation and Developer Experience

**User Story:** As a developer, I want comprehensive documentation and clear code structure, so that I can understand and contribute to the project.

#### Acceptance Criteria

1. THE VamsaSetu_System SHALL provide a README.md file with setup instructions, architecture overview, and API documentation
2. THE VamsaSetu_System SHALL include inline code comments for complex algorithms (especially Relationship_Engine)
3. THE VamsaSetu_System SHALL define PropTypes for all React components
4. THE VamsaSetu_System SHALL use a constants file for all UI strings (no hardcoded text)
5. THE VamsaSetu_System SHALL provide a .env.example file with all required environment variables documented
