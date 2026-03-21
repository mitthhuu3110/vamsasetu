package main

import (
	"testing"

	"vamsasetu/backend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestMigrations verifies that all models can be migrated successfully
func TestMigrations(t *testing.T) {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Notification{},
		&models.AuditLog{},
	)
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Verify tables exist
	tables := []string{"users", "events", "notifications", "audit_logs"}
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			t.Errorf("Table %s was not created", table)
		}
	}
}

// TestAuditLogModel verifies the AuditLog model structure
func TestAuditLogModel(t *testing.T) {
	auditLog := models.AuditLog{
		UserID:     1,
		Action:     "CREATE",
		EntityType: "member",
		EntityID:   "uuid-123",
	}

	if auditLog.TableName() != "audit_logs" {
		t.Errorf("Expected table name 'audit_logs', got '%s'", auditLog.TableName())
	}
}
