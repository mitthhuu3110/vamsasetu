# VamsaSetu Backend

Go backend service for VamsaSetu family tree and event management system.

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for running databases)

## Installation

1. Install Go dependencies:
```bash
go mod download
go mod tidy
```

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/                     # Environment configuration
│   ├── middleware/                 # HTTP middleware (auth, CORS, logging)
│   ├── models/                     # Data models (User, Event, Member, Relationship)
│   ├── repository/                 # Data access layer
│   ├── service/                    # Business logic layer
│   ├── handler/                    # HTTP handlers
│   ├── scheduler/                  # Background jobs (notification scheduler)
│   └── utils/                      # Utility functions
└── pkg/
    ├── neo4j/                      # Neo4j client
    ├── postgres/                   # PostgreSQL client
    └── redis/                      # Redis client
```

## Dependencies

- **Fiber v2**: Fast HTTP web framework
- **GORM**: ORM for PostgreSQL
- **Neo4j Go Driver**: Graph database client
- **Go-Redis**: Redis client
- **JWT-Go**: JWT authentication
- **Godotenv**: Environment variable management
- **Bcrypt**: Password hashing

## Running the Server

```bash
# Start databases with Docker Compose (from project root)
docker-compose up -d neo4j postgres redis

# Run the server
go run cmd/server/main.go
```

The server will start on port 8080.

## Environment Variables

See `.env.example` in the project root for required environment variables.

## Development

The backend follows a clean architecture pattern with clear separation of concerns:

- **Handlers**: HTTP request/response handling
- **Services**: Business logic and orchestration
- **Repositories**: Data access and persistence
- **Models**: Data structures and entities
- **Middleware**: Cross-cutting concerns (auth, logging, CORS)

## API Documentation

API endpoints will be documented as they are implemented. The API follows RESTful conventions and returns responses in the format:

```json
{
  "success": boolean,
  "data": any,
  "error": string
}
```
