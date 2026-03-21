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

func setupUserTestRepo(t *testing.T) (*UserRepository, context.Context, func()) {
	cfg := &config.Config{
		PostgresURL: "postgres://vamsasetu:vamsasetu123@localhost:5432/vamsasetu?sslmode=disable",
	}

	client, err := postgres.NewClient(cfg)
	require.NoError(t, err, "Failed to create PostgreSQL client")

	// Auto-migrate the User model
	err = client.DB.AutoMigrate(&models.User{})
	require.NoError(t, err, "Failed to migrate User model")

	ctx := context.Background()
	repo := NewUserRepository(client.DB)

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		client.DB.Where("email LIKE ?", "test%@example.com").Delete(&models.User{})
		client.Close()
	}

	return repo, ctx, cleanup
}

func TestUserRepository_Create(t *testing.T) {
	repo, ctx, cleanup := setupUserTestRepo(t)
	defer cleanup()

	user := &models.User{
		Email:        "test_create@example.com",
		PasswordHash: "hashed_password_123",
		Name:         "Test User",
		Role:         "owner",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err, "Create should not return error")
	assert.NotZero(t, user.ID, "User ID should be set after creation")

	// Verify user was created
	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, retrieved.Email)
	assert.Equal(t, user.Name, retrieved.Name)
	assert.Equal(t, user.Role, retrieved.Role)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo, ctx, cleanup := setupUserTestRepo(t)
	defer cleanup()

	// Create a test user
	user := &models.User{
		Email:        "test_getbyemail@example.com",
		PasswordHash: "hashed_password_456",
		Name:         "Test GetByEmail",
		Role:         "viewer",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test retrieval by email
	retrieved, err := repo.GetByEmail(ctx, "test_getbyemail@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, user.Email, retrieved.Email)
	assert.Equal(t, user.Name, retrieved.Name)
	assert.Equal(t, user.Role, retrieved.Role)

	// Test non-existent email
	_, err = repo.GetByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_GetByID(t *testing.T) {
	repo, ctx, cleanup := setupUserTestRepo(t)
	defer cleanup()

	// Create a test user
	user := &models.User{
		Email:        "test_getbyid@example.com",
		PasswordHash: "hashed_password_789",
		Name:         "Test GetByID",
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test retrieval by ID
	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, user.Email, retrieved.Email)
	assert.Equal(t, user.Name, retrieved.Name)

	// Test non-existent ID
	_, err = repo.GetByID(ctx, 99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_Update(t *testing.T) {
	repo, ctx, cleanup := setupUserTestRepo(t)
	defer cleanup()

	// Create a test user
	user := &models.User{
		Email:        "test_update@example.com",
		PasswordHash: "original_password",
		Name:         "Original Name",
		Role:         "owner",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Update the user
	user.Name = "Updated Name"
	user.PasswordHash = "new_password_hash"
	user.UpdatedAt = time.Now()

	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	// Verify update
	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", retrieved.Name)
	assert.Equal(t, "new_password_hash", retrieved.PasswordHash)
	assert.Equal(t, user.Email, retrieved.Email) // Email should remain unchanged
}

func TestUserRepository_UniqueEmail(t *testing.T) {
	repo, ctx, cleanup := setupUserTestRepo(t)
	defer cleanup()

	// Create first user
	user1 := &models.User{
		Email:        "test_unique@example.com",
		PasswordHash: "password1",
		Name:         "User 1",
		Role:         "owner",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := repo.Create(ctx, user1)
	require.NoError(t, err)

	// Try to create second user with same email
	user2 := &models.User{
		Email:        "test_unique@example.com",
		PasswordHash: "password2",
		Name:         "User 2",
		Role:         "viewer",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = repo.Create(ctx, user2)
	assert.Error(t, err, "Should not allow duplicate email")
}

func TestUserRepository_RoleValidation(t *testing.T) {
	repo, ctx, cleanup := setupUserTestRepo(t)
	defer cleanup()

	// Test valid roles
	validRoles := []string{"owner", "viewer", "admin"}
	for _, role := range validRoles {
		user := &models.User{
			Email:        "test_role_" + role + "@example.com",
			PasswordHash: "password",
			Name:         "Test " + role,
			Role:         role,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		err := repo.Create(ctx, user)
		assert.NoError(t, err, "Should accept valid role: "+role)
	}
}
