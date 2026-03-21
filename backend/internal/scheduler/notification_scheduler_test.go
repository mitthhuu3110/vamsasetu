package scheduler

import (
	"context"
	"testing"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/internal/service"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB creates a mock database for testing
func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm database: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

func TestNotificationScheduler_Start(t *testing.T) {
	db, _, cleanup := setupTestDB(t)
	defer cleanup()
	
	notifRepo := repository.NewNotificationRepository(db)
	notifSvc := service.NewNotificationService(notifRepo, "", "", "", "", "")
	
	scheduler := NewNotificationScheduler(notifRepo, notifSvc)
	
	// Test that scheduler can be started
	scheduler.Start()
	
	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Stop the scheduler
	scheduler.Stop()
	
	// Test passes if no panic occurs
}

func TestNotificationScheduler_ProcessNotifications(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	
	now := time.Now()
	scheduledAt := now.Add(-1 * time.Hour)
	
	// Mock GetPending query
	rows := sqlmock.NewRows([]string{
		"id", "event_id", "user_id", "channel", "scheduled_at",
		"sent_at", "status", "retry_count", "error_msg", "created_at", "updated_at",
	}).AddRow(1, 1, 1, "email", scheduledAt, nil, "pending", 0, "", now, now)
	
	mock.ExpectQuery(`SELECT \* FROM "notifications"`).
		WithArgs("pending", sqlmock.AnyArg()).
		WillReturnRows(rows)
	
	// Mock UpdateStatus for marking as sent
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "notifications"`).
		WithArgs(
			sqlmock.AnyArg(), // error_msg
			sqlmock.AnyArg(), // sent_at
			sqlmock.AnyArg(), // status
			sqlmock.AnyArg(), // updated_at
			1,                // id
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	
	// Create scheduler and process
	notifRepo := repository.NewNotificationRepository(db)
	notifSvc := service.NewNotificationService(notifRepo, "", "", "", "", "")
	scheduler := NewNotificationScheduler(notifRepo, notifSvc)
	
	// Process notifications
	scheduler.processNotifications()
	
	// Wait for goroutines to complete
	time.Sleep(500 * time.Millisecond)
	
	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestNotificationScheduler_WorkerPool(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	
	now := time.Now()
	scheduledAt := now.Add(-1 * time.Hour)
	
	// Mock GetPending query - return 15 notifications
	rows := sqlmock.NewRows([]string{
		"id", "event_id", "user_id", "channel", "scheduled_at",
		"sent_at", "status", "retry_count", "error_msg", "created_at", "updated_at",
	})
	
	for i := 1; i <= 15; i++ {
		rows.AddRow(i, 1, 1, "email", scheduledAt, nil, "pending", 0, "", now, now)
	}
	
	mock.ExpectQuery(`SELECT \* FROM "notifications"`).
		WithArgs("pending", sqlmock.AnyArg()).
		WillReturnRows(rows)
	
	// Mock UpdateStatus for each notification
	for i := 0; i < 15; i++ {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "notifications"`).
			WithArgs(
				sqlmock.AnyArg(), // error_msg
				sqlmock.AnyArg(), // sent_at
				sqlmock.AnyArg(), // status
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // id
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}
	
	// Create scheduler and process
	notifRepo := repository.NewNotificationRepository(db)
	notifSvc := service.NewNotificationService(notifRepo, "", "", "", "", "")
	scheduler := NewNotificationScheduler(notifRepo, notifSvc)
	
	// Process notifications
	scheduler.processNotifications()
	
	// Wait for all goroutines to complete
	time.Sleep(1 * time.Second)
	
	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestNotificationScheduler_GracefulShutdown(t *testing.T) {
	db, _, cleanup := setupTestDB(t)
	defer cleanup()
	
	notifRepo := repository.NewNotificationRepository(db)
	notifSvc := service.NewNotificationService(notifRepo, "", "", "", "", "")
	
	scheduler := NewNotificationScheduler(notifRepo, notifSvc)
	
	// Start scheduler
	scheduler.Start()
	time.Sleep(100 * time.Millisecond)
	
	// Stop scheduler
	scheduler.Stop()
	
	// Verify context is cancelled
	select {
	case <-scheduler.ctx.Done():
		// Context cancelled as expected
	case <-time.After(1 * time.Second):
		t.Error("Expected context to be cancelled after Stop()")
	}
}
