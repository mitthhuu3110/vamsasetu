# VamsaSetu Backend

## Overview
Spring Boot backend for VamsaSetu – family tree visualization, events, and notification platform. Uses Neo4j (relations), PostgreSQL (users/events), connects to common notification APIs.

## Setup
1. Ensure Java 17+ is installed.
2. Setup PostgreSQL and Neo4j servers (local/docker recommended, see configs below).
3. Copy `.env.example` to `.env` and fill in values.
4. `mvnw spring-boot:run` (or build and run via IDE).

## Configuration
- `src/main/resources/application.yml` for active profiles & DB connection setup.
- `.env` (not under version control) for API secrets & sensitive config.

## Main Modules
- **Auth**: Registration, login, social login, JWT handling
- **FamilyTree**: Graph APIs (CRUD on nodes/edges with Neo4j)
- **RelationsEngine**: Relationship inference endpoints
- **EventManager**: Add/edit/delete/notify events (Postgres main store)
- **NotificationService**: Webhook/APIs for SMS, WhatsApp, Email

## Testing
`mvn test` runs unit and integration tests (mock DBs required)
