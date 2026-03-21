# Database Migration Implementation

## Overview

This document describes the GORM auto-migration implementation for VamsaSetu's PostgreSQL database.

## Implementation Details

### Migration Script Location
- **File**: `backend/cmd/server/main.go`
- **Function**: `runMigrations()`

### Migrated Models

The migration script auto-migrates the following models:

1. **User** (`users` table)
   - User accounts with authentication credentials
   - Roles: owner, viewer, admin
   - Foreign key relationships with events and audit logs

2. **Event** (`events` table)
   - Family events (birthdays, anniversaries, ceremonies)
   - Associated with multiple members via member_ids array
   - Foreign key to users table (created_by)

3. **Notification** (`notifications` table)
   - Scheduled notifications for events
   - Channels: whatsapp, sms, email
   - Status tracking with retry logic
   - Foreign keys to events and users tables

4. **AuditLog** (`audit_logs` table)
   - Audit trail for all data modifications
   - Tracks user actions, entity types, and changes
   - JSONB details field for flexible metadata storage
   - Foreign key to users table

### Database Schema Features

#### Indexes
- User email (unique index)
- Event date and created_by (indexes)
- Notification scheduled_at and status (composite index)
- Notification event_id (index)
- AuditLog user_id and created_at (indexes)

#### Constraints
- Check constraints for enum fields (role, event_type, channel, status)
- Foreign key constraints with CASCADE delete
- NOT NULL constraints on required fields

#### Data Types
- Text arrays for member_ids in events
- JSONB for flexible details in audit_logs
- Timestamps with auto-creation and auto-update

### Migration Execution Flow

```
1. Load configuration from environment variables
2. Initialize PostgreSQL client
3. Verify database connectivity
4. Run GORM AutoMigrate for all models
5. Log migration status for each table
6. Initialize Neo4j client (for graph data)
7. Verify Neo4j connectivity
```

### Error Handling

The migration script includes comprehensive error handling:
- Configuration validation errors
- Database connection errors
- Migration execution errors
- Health check failures

All errors are logged with descriptive messages and cause the application to exit gracefully.

### Running Migrations

#### Prerequisites
1. PostgreSQL database running and accessible
2. Neo4j database running and accessible
3. Environment variables configured (see .env.example)

#### Execution
```bash
# Using Go directly
cd backend
go run cmd/server/main.go

# Using Docker Compose
docker-compose up backend
```

#### Expected Output
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
VamsaSetu Backend Server - Ready
```

### Testing

A test suite is provided in `main_test.go` to verify:
- All models can be migrated successfully
- Tables are created with correct names
- Model structures are valid

Run tests with:
```bash
cd backend/cmd/server
go test -v
```

## Requirements Validation

This implementation satisfies **Requirement 11.2** from the requirements document:

> "THE VamsaSetu_System SHALL persist all User, Event, and audit log data in the PostgreSQL_Database"

The migration script ensures:
- ✅ User table is created and migrated
- ✅ Event table is created and migrated
- ✅ Notification table is created and migrated
- ✅ AuditLog table is created and migrated
- ✅ All foreign key relationships are established
- ✅ All indexes are created for performance
- ✅ All constraints are enforced for data integrity

## Next Steps

After successful migration, the following tasks can proceed:
- Implement repository layer for data access
- Implement service layer for business logic
- Implement API handlers for HTTP endpoints
- Implement audit logging service
