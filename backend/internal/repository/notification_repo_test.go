package repository

import (
	"context"
	"testing"
	"time"

	"vamsasetu/backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupNotificationRepoTest(t *testing.T) (*NotificationRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewNotificationRepository(gormDB)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestNotificationRepository_Create(t *testing.T) {
	repo, mock, cleanup := setupNotificationRepoTest(t)
	defer cleanup()

	notification := &models.Notification{
		EventID:     1,
		UserID:      1,
		Channel:     "email",
		ScheduledAt: time.Now().Add(24 * time.Hour),
		Status:      "pending",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "notifications"`).
		WithArgs(
			sqlmock.AnyArg(), // event_id
			sqlmock.AnyArg(), // user_id
			sqlmock.AnyArg(), // channel
			sqlmock.AnyArg(), // scheduled_at
			sqlmock.AnyArg(), // sent_at
			sqlmock.AnyArg(), // status
			sqlmock.AnyArg(), // retry_count
			sqlmock.AnyArg(), // error_msg
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), notification)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_GetPending(t *testing.T) {
	repo, mock, cleanup := setupNotificationRepoTest(t)
	defer cleanup()

	now := time.Now()
	scheduledAt := now.Add(-1 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "event_id", "user_id", "channel", "scheduled_at",
		"sent_at", "status", "retry_count", "error_msg", "created_at", "updated_at",
	}).
		AddRow(1, 1, 1, "email", scheduledAt, nil, "pending", 0, "", now, now).
		AddRow(2, 2, 1, "whatsapp", scheduledAt, nil, "pending", 0, "", now, now)

	mock.ExpectQuery(`SELECT \* FROM "notifications"`).
		WithArgs("pending", sqlmock.AnyArg()).
		WillReturnRows(rows)

	notifications, err := repo.GetPending(context.Background())
	assert.NoError(t, err)
	assert.Len(t, notifications, 2)
	assert.Equal(t, uint(1), notifications[0].ID)
	assert.Equal(t, "email", notifications[0].Channel)
	assert.Equal(t, uint(2), notifications[1].ID)
	assert.Equal(t, "whatsapp", notifications[1].Channel)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_UpdateStatus(t *testing.T) {
	repo, mock, cleanup := setupNotificationRepoTest(t)
	defer cleanup()

	sentAt := time.Now()

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

	err := repo.UpdateStatus(context.Background(), 1, "sent", &sentAt, "")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_UpdateStatus_Failed(t *testing.T) {
	repo, mock, cleanup := setupNotificationRepoTest(t)
	defer cleanup()

	errorMsg := "failed to send email"

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "notifications"`).
		WithArgs(
			sqlmock.AnyArg(), // error_msg
			sqlmock.AnyArg(), // status
			sqlmock.AnyArg(), // updated_at
			1,                // id
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateStatus(context.Background(), 1, "failed", nil, errorMsg)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_IncrementRetry(t *testing.T) {
	repo, mock, cleanup := setupNotificationRepoTest(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "notifications" SET "retry_count"=retry_count \+ \$1 WHERE id = \$2`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.IncrementRetry(context.Background(), 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_GetPending_EmptyResult(t *testing.T) {
	repo, mock, cleanup := setupNotificationRepoTest(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "event_id", "user_id", "channel", "scheduled_at",
		"sent_at", "status", "retry_count", "error_msg", "created_at", "updated_at",
	})

	mock.ExpectQuery(`SELECT \* FROM "notifications"`).
		WithArgs("pending", sqlmock.AnyArg()).
		WillReturnRows(rows)

	notifications, err := repo.GetPending(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, notifications)
	assert.NoError(t, mock.ExpectationsWereMet())
}
