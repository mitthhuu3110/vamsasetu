# VamsaSetu Models

This directory contains the GORM models for the VamsaSetu PostgreSQL database.

## Models Overview

### User Model (`user.go`)

Represents authenticated users in the system.

**Fields:**
- `ID` (uint, primary key): Auto-incrementing user ID
- `Email` (string, unique, indexed): User's email address
- `PasswordHash` (string): Bcrypt hashed password (excluded from JSON)
- `Name` (string): User's full name
- `Role` (string): User role - must be one of: `owner`, `viewer`, `admin`
- `CreatedAt` (time.Time): Timestamp of user creation
- `UpdatedAt` (time.Time): Timestamp of last update

**Constraints:**
- Email must be unique and not null
- Role must be one of the valid values (enforced by CHECK constraint)
- Email is indexed for fast lookups

**Table Name:** `users`

### Event Model (`event.go`)

Represents family events (birthdays, anniversaries, ceremonies, etc.).

**Fields:**
- `ID` (uint, primary key): Auto-incrementing event ID
- `Title` (string): Event title
- `Description` (string): Event description (optional)
- `EventDate` (time.Time, indexed): Date and time of the event
- `EventType` (string): Event type - must be one of: `birthday`, `anniversary`, `ceremony`, `custom`
- `MemberIDs` ([]string): Array of member UUIDs associated with the event
- `CreatedBy` (uint, foreign key, indexed): User ID who created the event
- `CreatedAt` (time.Time): Timestamp of event creation
- `UpdatedAt` (time.Time): Timestamp of last update
- `User` (User): Foreign key relationship to User model

**Constraints:**
- EventType must be one of the valid values (enforced by CHECK constraint)
- EventDate is indexed for efficient date-based queries
- CreatedBy is indexed for user-based queries
- Foreign key constraint with CASCADE delete on User

**Table Name:** `events`

### Notification Model (`notification.go`)

Represents scheduled notifications for events.

**Fields:**
- `ID` (uint, primary key): Auto-incrementing notification ID
- `EventID` (uint, foreign key, indexed): Associated event ID
- `UserID` (uint, foreign key): User to notify
- `Channel` (string): Notification channel - must be one of: `whatsapp`, `sms`, `email`
- `ScheduledAt` (time.Time, indexed): When the notification should be sent
- `SentAt` (*time.Time): When the notification was actually sent (nullable)
- `Status` (string, indexed): Notification status - must be one of: `pending`, `sent`, `failed`
- `RetryCount` (int): Number of retry attempts (default: 0)
- `ErrorMsg` (string): Error message if notification failed
- `CreatedAt` (time.Time): Timestamp of notification creation
- `UpdatedAt` (time.Time): Timestamp of last update
- `Event` (Event): Foreign key relationship to Event model
- `User` (User): Foreign key relationship to User model

**Constraints:**
- Channel must be one of the valid values (enforced by CHECK constraint)
- Status must be one of the valid values (enforced by CHECK constraint)
- Composite index on (ScheduledAt, Status) for efficient scheduler queries
- Foreign key constraints with CASCADE delete on Event and User

**Hooks:**
- `BeforeCreate`: Sets default status to "pending" if not specified

**Table Name:** `notifications`

## Database Indexes

The following indexes are automatically created by GORM:

1. **users table:**
   - Primary key index on `id`
   - Unique index on `email`

2. **events table:**
   - Primary key index on `id`
   - Index on `event_date` (for date-based queries)
   - Index on `created_by` (for user-based queries)

3. **notifications table:**
   - Primary key index on `id`
   - Index on `event_id` (for event-based queries)
   - Composite index `idx_notifications_scheduled` on `(scheduled_at, status)` (for scheduler queries)

## Usage Example

```go
package main

import (
    "vamsasetu/backend/internal/models"
    "gorm.io/gorm"
    "time"
)

func CreateUser(db *gorm.DB) error {
    user := models.User{
        Email:        "user@example.com",
        PasswordHash: "hashed_password_here",
        Name:         "John Doe",
        Role:         "owner",
    }
    
    return db.Create(&user).Error
}

func CreateEvent(db *gorm.DB, userID uint) error {
    event := models.Event{
        Title:       "Birthday Party",
        Description: "John's 30th birthday",
        EventDate:   time.Now().AddDate(0, 0, 7),
        EventType:   "birthday",
        MemberIDs:   []string{"uuid-1", "uuid-2"},
        CreatedBy:   userID,
    }
    
    return db.Create(&event).Error
}

func CreateNotification(db *gorm.DB, eventID, userID uint) error {
    notification := models.Notification{
        EventID:     eventID,
        UserID:      userID,
        Channel:     "email",
        ScheduledAt: time.Now().AddDate(0, 0, 6), // 6 days before event
        Status:      "pending",
    }
    
    return db.Create(&notification).Error
}
```

## Migration

To create the database tables, use GORM's AutoMigrate:

```go
db.AutoMigrate(&models.User{}, &models.Event{}, &models.Notification{})
```

This will create all tables with the appropriate constraints, indexes, and foreign keys.

## Testing

Run the model tests with:

```bash
go test ./internal/models/...
```

The test suite includes:
- Model structure validation
- Table name verification
- Field constraint validation
- Hook behavior verification
- Valid enum value testing

## Requirements Mapping

These models satisfy the following requirements from the design document:

- **Requirement 2.1**: User entity with email, password_hash, name, role
- **Requirement 5.1**: Event entity with title, description, event_date, event_type, member_ids
- **Requirement 6.1**: Notification entity with event_id, user_id, channel, scheduled_at, status

## Notes

- All timestamps use GORM's automatic timestamp management (`autoCreateTime`, `autoUpdateTime`)
- The `PasswordHash` field is excluded from JSON serialization using the `json:"-"` tag
- Foreign key relationships use CASCADE delete to maintain referential integrity
- Array fields (like `MemberIDs`) use PostgreSQL's native array type (`text[]`)
- All CHECK constraints are defined at the database level for data integrity
