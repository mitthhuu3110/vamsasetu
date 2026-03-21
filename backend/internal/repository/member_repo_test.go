package repository

import (
	"context"
	"testing"
	"time"

	"vamsasetu/backend/internal/config"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/pkg/neo4j"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) (*MemberRepository, context.Context, func()) {
	cfg := &config.Config{
		Neo4jURI:      "bolt://localhost:7687",
		Neo4jUsername: "neo4j",
		Neo4jPassword: "vamsasetu123",
	}

	client, err := neo4j.NewClient(cfg)
	require.NoError(t, err, "Failed to create Neo4j client")

	ctx := context.Background()

	// Ensure indexes exist
	repo := NewMemberRepository(client)
	err = repo.EnsureIndexes(ctx)
	require.NoError(t, err, "Failed to ensure indexes")

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		session := client.Driver.NewSession(ctx, neo4jDriver.SessionConfig{AccessMode: neo4jDriver.AccessModeWrite})
		defer session.Close(ctx)
		session.Run(ctx, "MATCH (m:Member) WHERE m.name STARTS WITH 'Test' DELETE m", nil)
		client.Close(ctx)
	}

	return repo, ctx, cleanup
}

func TestMemberRepository_Create(t *testing.T) {
	repo, ctx, cleanup := setupTestRepo(t)
	defer cleanup()

	member := models.NewMember("Test User", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "male")
	member.Email = "test@example.com"
	member.Phone = "+919876543210"

	err := repo.Create(ctx, member)
	assert.NoError(t, err, "Create should not return error")

	// Verify member was created
	retrieved, err := repo.GetByID(ctx, member.ID)
	assert.NoError(t, err)
	assert.Equal(t, member.Name, retrieved.Name)
	assert.Equal(t, member.Gender, retrieved.Gender)
	assert.Equal(t, member.Email, retrieved.Email)
}

func TestMemberRepository_GetByID(t *testing.T) {
	repo, ctx, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test member
	member := models.NewMember("Test GetByID", time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC), "female")
	err := repo.Create(ctx, member)
	require.NoError(t, err)

	// Test retrieval
	retrieved, err := repo.GetByID(ctx, member.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, member.ID, retrieved.ID)
	assert.Equal(t, member.Name, retrieved.Name)
	assert.Equal(t, member.Gender, retrieved.Gender)

	// Test non-existent ID
	_, err = repo.GetByID(ctx, "non-existent-id")
	assert.Error(t, err)
}

func TestMemberRepository_GetAll(t *testing.T) {
	repo, ctx, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create multiple test members
	member1 := models.NewMember("Test GetAll 1", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "male")
	member2 := models.NewMember("Test GetAll 2", time.Date(1992, 2, 2, 0, 0, 0, 0, time.UTC), "female")

	err := repo.Create(ctx, member1)
	require.NoError(t, err)
	err = repo.Create(ctx, member2)
	require.NoError(t, err)

	// Test retrieval
	members, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(members), 2, "Should retrieve at least 2 members")

	// Verify our test members are in the results
	found := 0
	for _, m := range members {
		if m.ID == member1.ID || m.ID == member2.ID {
			found++
		}
	}
	assert.Equal(t, 2, found, "Should find both test members")
}

func TestMemberRepository_Update(t *testing.T) {
	repo, ctx, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test member
	member := models.NewMember("Test Update Original", time.Date(1988, 3, 10, 0, 0, 0, 0, time.UTC), "male")
	err := repo.Create(ctx, member)
	require.NoError(t, err)

	// Update the member
	member.Name = "Test Update Modified"
	member.Email = "updated@example.com"
	member.Update()

	err = repo.Update(ctx, member)
	assert.NoError(t, err)

	// Verify update
	retrieved, err := repo.GetByID(ctx, member.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Test Update Modified", retrieved.Name)
	assert.Equal(t, "updated@example.com", retrieved.Email)

	// Test updating non-existent member
	nonExistent := models.NewMember("Non Existent", time.Now(), "male")
	err = repo.Update(ctx, nonExistent)
	assert.Error(t, err)
}

func TestMemberRepository_SoftDelete(t *testing.T) {
	repo, ctx, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test member
	member := models.NewMember("Test SoftDelete", time.Date(1995, 7, 20, 0, 0, 0, 0, time.UTC), "female")
	err := repo.Create(ctx, member)
	require.NoError(t, err)

	// Soft delete the member
	err = repo.SoftDelete(ctx, member.ID)
	assert.NoError(t, err)

	// Verify member is not returned by GetByID (because isDeleted = true)
	_, err = repo.GetByID(ctx, member.ID)
	assert.Error(t, err, "Soft deleted member should not be retrievable")

	// Verify member is not in GetAll results
	members, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	for _, m := range members {
		assert.NotEqual(t, member.ID, m.ID, "Soft deleted member should not appear in GetAll")
	}

	// Test soft deleting non-existent member
	err = repo.SoftDelete(ctx, "non-existent-id")
	assert.Error(t, err)
}

func TestMemberRepository_EnsureIndexes(t *testing.T) {
	repo, ctx, cleanup := setupTestRepo(t)
	defer cleanup()

	// Test that EnsureIndexes can be called multiple times without error
	err := repo.EnsureIndexes(ctx)
	assert.NoError(t, err)

	err = repo.EnsureIndexes(ctx)
	assert.NoError(t, err, "EnsureIndexes should be idempotent")
}

func TestMemberRepository_EmptyFields(t *testing.T) {
	repo, ctx, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create member with minimal fields
	member := models.NewMember("Test Minimal", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "other")
	// Leave email, phone, avatarUrl empty

	err := repo.Create(ctx, member)
	assert.NoError(t, err)

	// Verify retrieval handles empty fields
	retrieved, err := repo.GetByID(ctx, member.ID)
	assert.NoError(t, err)
	assert.Equal(t, "", retrieved.Email)
	assert.Equal(t, "", retrieved.Phone)
	assert.Equal(t, "", retrieved.AvatarURL)
}
