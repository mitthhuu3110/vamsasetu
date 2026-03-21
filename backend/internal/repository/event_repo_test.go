package repository

import (
	"context"
	"testing"
	"time"

	"vamsasetu/backend/internal/config"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/pkg/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupEventTestRepo(t *testing.T) (*EventRepository, *UserRepository, context.Context, func()) {
	cfg := &config.Config{
		PostgresURL: "postgres://vamsasetu:vamsasetu123@localhost:5432/vamsasetu?sslmode=disable",
	}

	client, err := postgres.NewClient(cfg)
	require.NoError(t, err, "Failed to create PostgreSQL client")

	// Auto-migrate the models
	err = client.DB.AutoMigrate(&models.User{}, &models.Event{})
	require.NoError(t, err, "Failed to migrate models")

	ctx := context.Background()
	eventRepo := NewEventRepository(client.DB)
	userRepo := NewUserRepository(client.DB)

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		client.DB.Where("title LIKE ?", "Test%").Delete(&models.Event{})
		client.DB.Where("email LIKE ?", "test%@example.com").Delete(&models.User{})
		client.Close()
	}

	return eventRepo, userRepo, ctx, cleanup
}

func createTestUser(t *testing.T, userRepo *UserRepository, ctx context.Context) *models.User {
	user := &models.User{
		Email:        "test_event_user@example.com",
		PasswordHash: "hashed_password",
		Name:         "Test Event User",
		Role:         "owner",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	return user
}

func TestEventRepository_Create(t *testing.T) {
	repo, userRepo, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	user := createTestUser(t, userRepo, ctx)

	event := &models.Event{
		Title:       "Test Birthday",
		Description: "Test birthday event",
		EventDate:   time.Now().Add(7 * 24 * time.Hour),
		EventType:   "birthday",
		MemberIDs:   []string{"member-uuid-1", "member-uuid-2"},
		CreatedBy:   user.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, event)
	assert.NoError(t, err, "Create should not return error")
	assert.NotZero(t, event.ID, "Event ID should be set after creation")

	// Verify event was created
	retrieved, err := repo.GetByID(ctx, event.ID)
	assert.NoError(t, err)
	assert.Equal(t, event.Title, retrieved.Title)
	assert.Equal(t, event.EventType, retrieved.EventType)
	assert.Equal(t, len(event.MemberIDs), len(retrieved.MemberIDs))
}

func TestEventRepository_GetByID(t *testing.T) {
	repo, userRepo, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	user := createTestUser(t, userRepo, ctx)

	event := &models.Event{
		Title:       "Test Anniversary",
		Description: "Test anniversary event",
		EventDate:   time.Now().Add(14 * 24 * time.Hour),
		EventType:   "anniversary",
		MemberIDs:   []string{"member-uuid-3"},
		CreatedBy:   user.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, event)
	require.NoError(t, err)

	// Test retrieval by ID
	retrieved, err := repo.GetByID(ctx, event.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, event.ID, retrieved.ID)
	assert.Equal(t, event.Title, retrieved.Title)
	assert.Equal(t, event.EventType, retrieved.EventType)

	// Test non-existent ID
	_, err = repo.GetByID(ctx, 99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "event not found")
}

func TestEventRepository_GetAll(t *testing.T) {
	repo, userRepo, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	user := createTestUser(t, userRepo, ctx)

	// Create multiple events
	events := []*models.Event{
		{
			Title:       "Test Event 1",
			Description: "First test event",
			EventDate:   time.Now().Add(1 * 24 * time.Hour),
			EventType:   "birthday",
			MemberIDs:   []string{"member-1"},
			CreatedBy:   user.ID,
		},
		{
			Title:       "Test Event 2",
			Description: "Second test event",
			EventDate:   time.Now().Add(2 * 24 * time.Hour),
			EventType:   "ceremony",
			MemberIDs:   []string{"member-2"},
			CreatedBy:   user.ID,
		},
	}

	for _, event := range events {
		err := repo.Create(ctx, event)
		require.NoError(t, err)
	}

	// Retrieve all events
	retrieved, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(retrieved), 2, "Should retrieve at least 2 events")
}

func TestEventRepository_Update(t *testing.T) {
	repo, userRepo, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	user := createTestUser(t, userRepo, ctx)

	event := &models.Event{
		Title:       "Test Original Title",
		Description: "Original description",
		EventDate:   time.Now().Add(5 * 24 * time.Hour),
		EventType:   "custom",
		MemberIDs:   []string{"member-uuid-4"},
		CreatedBy:   user.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, event)
	require.NoError(t, err)

	// Update the event
	event.Title = "Test Updated Title"
	event.Description = "Updated description"
	event.UpdatedAt = time.Now()

	err = repo.Update(ctx, event)
	assert.NoError(t, err)

	// Verify update
	retrieved, err := repo.GetByID(ctx, event.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Test Updated Title", retrieved.Title)
	assert.Equal(t, "Updated description", retrieved.Description)
}

func TestEventRepository_Delete(t *testing.T) {
	repo, userRepo, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	user := createTestUser(t, userRepo, ctx)

	event := &models.Event{
		Title:       "Test Delete Event",
		Description: "Event to be deleted",
		EventDate:   time.Now().Add(3 * 24 * time.Hour),
		EventType:   "birthday",
		MemberIDs:   []string{"member-uuid-5"},
		CreatedBy:   user.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, event)
	require.NoError(t, err)

	// Delete the event
	err = repo.Delete(ctx, event.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(ctx, event.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "event not found")
}

func TestEventRepository_GetUpcoming(t *testing.T) {
	repo, userRepo, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	user := createTestUser(t, userRepo, ctx)

	now := time.Now()

	// Create events at different dates
	events := []*models.Event{
		{
			Title:       "Test Upcoming 1",
			Description: "Event in 2 days",
			EventDate:   now.Add(2 * 24 * time.Hour),
			EventType:   "birthday",
			MemberIDs:   []string{"member-1"},
			CreatedBy:   user.ID,
		},
		{
			Title:       "Test Upcoming 2",
			Description: "Event in 5 days",
			EventDate:   now.Add(5 * 24 * time.Hour),
			EventType:   "anniversary",
			MemberIDs:   []string{"member-2"},
			CreatedBy:   user.ID,
		},
		{
			Title:       "Test Past Event",
			Description: "Event in the past",
			EventDate:   now.Add(-2 * 24 * time.Hour),
			EventType:   "ceremony",
			MemberIDs:   []string{"member-3"},
			CreatedBy:   user.ID,
		},
	}

	for _, event := range events {
		err := repo.Create(ctx, event)
		require.NoError(t, err)
	}

	// Get upcoming events (next 7 days)
	startDate := now
	endDate := now.Add(7 * 24 * time.Hour)
	upcoming, err := repo.GetUpcoming(ctx, startDate, endDate)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(upcoming), 2, "Should retrieve at least 2 upcoming events")

	// Verify all retrieved events are within range
	for _, event := range upcoming {
		assert.True(t, event.EventDate.After(startDate) || event.EventDate.Equal(startDate))
		assert.True(t, event.EventDate.Before(endDate) || event.EventDate.Equal(endDate))
	}
}

func TestEventRepository_GetByType(t *testing.T) {
	repo, userRepo, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	user := createTestUser(t, userRepo, ctx)

	// Create events of different types
	events := []*models.Event{
		{
			Title:       "Test Birthday 1",
			Description: "Birthday event",
			EventDate:   time.Now().Add(1 * 24 * time.Hour),
			EventType:   "birthday",
			MemberIDs:   []string{"member-1"},
			CreatedBy:   user.ID,
		},
		{
			Title:       "Test Birthday 2",
			Description: "Another birthday event",
			EventDate:   time.Now().Add(2 * 24 * time.Hour),
			EventType:   "birthday",
			MemberIDs:   []string{"member-2"},
			CreatedBy:   user.ID,
		},
		{
			Title:       "Test Ceremony",
			Description: "Ceremony event",
			EventDate:   time.Now().Add(3 * 24 * time.Hour),
			EventType:   "ceremony",
			MemberIDs:   []string{"member-3"},
			CreatedBy:   user.ID,
		},
	}

	for _, event := range events {
		err := repo.Create(ctx, event)
		require.NoError(t, err)
	}

	// Get birthday events
	birthdays, err := repo.GetByType(ctx, "birthday")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(birthdays), 2, "Should retrieve at least 2 birthday events")

	// Verify all retrieved events are birthdays
	for _, event := range birthdays {
		assert.Equal(t, "birthday", event.EventType)
	}
}

func TestEventRepository_GetByMember(t *testing.T) {
	repo, userRepo, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	user := createTestUser(t, userRepo, ctx)

	memberID := "test-member-uuid-123"

	// Create events with different members
	events := []*models.Event{
		{
			Title:       "Test Member Event 1",
			Description: "Event for specific member",
			EventDate:   time.Now().Add(1 * 24 * time.Hour),
			EventType:   "birthday",
			MemberIDs:   []string{memberID, "other-member"},
			CreatedBy:   user.ID,
		},
		{
			Title:       "Test Member Event 2",
			Description: "Another event for specific member",
			EventDate:   time.Now().Add(2 * 24 * time.Hour),
			EventType:   "anniversary",
			MemberIDs:   []string{memberID},
			CreatedBy:   user.ID,
		},
		{
			Title:       "Test Other Event",
			Description: "Event for different member",
			EventDate:   time.Now().Add(3 * 24 * time.Hour),
			EventType:   "ceremony",
			MemberIDs:   []string{"different-member"},
			CreatedBy:   user.ID,
		},
	}

	for _, event := range events {
		err := repo.Create(ctx, event)
		require.NoError(t, err)
	}

	// Get events for specific member
	memberEvents, err := repo.GetByMember(ctx, memberID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(memberEvents), 2, "Should retrieve at least 2 events for the member")

	// Verify all retrieved events contain the member ID
	for _, event := range memberEvents {
		found := false
		for _, id := range event.MemberIDs {
			if id == memberID {
				found = true
				break
			}
		}
		assert.True(t, found, "Event should contain the member ID")
	}
}

func TestEventRepository_GetByDateRange(t *testing.T) {
	repo, userRepo, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	user := createTestUser(t, userRepo, ctx)

	now := time.Now()

	// Create events at different dates
	events := []*models.Event{
		{
			Title:       "Test Range Event 1",
			Description: "Event in range",
			EventDate:   now.Add(3 * 24 * time.Hour),
			EventType:   "birthday",
			MemberIDs:   []string{"member-1"},
			CreatedBy:   user.ID,
		},
		{
			Title:       "Test Range Event 2",
			Description: "Another event in range",
			EventDate:   now.Add(5 * 24 * time.Hour),
			EventType:   "anniversary",
			MemberIDs:   []string{"member-2"},
			CreatedBy:   user.ID,
		},
		{
			Title:       "Test Out of Range",
			Description: "Event outside range",
			EventDate:   now.Add(15 * 24 * time.Hour),
			EventType:   "ceremony",
			MemberIDs:   []string{"member-3"},
			CreatedBy:   user.ID,
		},
	}

	for _, event := range events {
		err := repo.Create(ctx, event)
		require.NoError(t, err)
	}

	// Get events in date range (2-10 days from now)
	startDate := now.Add(2 * 24 * time.Hour)
	endDate := now.Add(10 * 24 * time.Hour)
	rangeEvents, err := repo.GetByDateRange(ctx, startDate, endDate)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(rangeEvents), 2, "Should retrieve at least 2 events in range")

	// Verify all retrieved events are within range
	for _, event := range rangeEvents {
		assert.True(t, event.EventDate.After(startDate) || event.EventDate.Equal(startDate))
		assert.True(t, event.EventDate.Before(endDate) || event.EventDate.Equal(endDate))
	}
}

func TestEventRepository_GetByCreator(t *testing.T) {
	repo, userRepo, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	user1 := createTestUser(t, userRepo, ctx)

	// Create second user
	user2 := &models.User{
		Email:        "test_event_user2@example.com",
		PasswordHash: "hashed_password",
		Name:         "Test Event User 2",
		Role:         "owner",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, user2)
	require.NoError(t, err)

	// Create events for different users
	events := []*models.Event{
		{
			Title:       "Test User1 Event 1",
			Description: "Event by user 1",
			EventDate:   time.Now().Add(1 * 24 * time.Hour),
			EventType:   "birthday",
			MemberIDs:   []string{"member-1"},
			CreatedBy:   user1.ID,
		},
		{
			Title:       "Test User1 Event 2",
			Description: "Another event by user 1",
			EventDate:   time.Now().Add(2 * 24 * time.Hour),
			EventType:   "anniversary",
			MemberIDs:   []string{"member-2"},
			CreatedBy:   user1.ID,
		},
		{
			Title:       "Test User2 Event",
			Description: "Event by user 2",
			EventDate:   time.Now().Add(3 * 24 * time.Hour),
			EventType:   "ceremony",
			MemberIDs:   []string{"member-3"},
			CreatedBy:   user2.ID,
		},
	}

	for _, event := range events {
		err := repo.Create(ctx, event)
		require.NoError(t, err)
	}

	// Get events by user1
	user1Events, err := repo.GetByCreator(ctx, user1.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(user1Events), 2, "Should retrieve at least 2 events for user1")

	// Verify all retrieved events are created by user1
	for _, event := range user1Events {
		assert.Equal(t, user1.ID, event.CreatedBy)
	}
}
