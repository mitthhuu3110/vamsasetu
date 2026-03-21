package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/pkg/redis"
)

// WebSocketHub interface for broadcasting updates
type WebSocketHub interface {
	BroadcastUpdate(eventType string, data interface{})
}

// MemberService handles member business logic with caching
type MemberService struct {
	repo  *repository.MemberRepository
	cache *redis.Client
	hub   WebSocketHub
}

// Cache TTL constants
const (
	MemberCacheTTL     = 10 * time.Minute
	SearchCacheTTL     = 2 * time.Minute
	FamilyTreeCacheTTL = 5 * time.Minute
)

// NewMemberService creates a new member service instance
func NewMemberService(repo *repository.MemberRepository, cache *redis.Client, hub WebSocketHub) *MemberService {
	return &MemberService{
		repo:  repo,
		cache: cache,
		hub:   hub,
	}
}

// Create creates a new member and invalidates relevant caches
func (s *MemberService) Create(ctx context.Context, member *models.Member) error {
	// Validate member data
	if err := member.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Create member in database
	if err := s.repo.Create(ctx, member); err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}

	// Invalidate family tree and search caches
	s.invalidateFamilyTreeCache(ctx)
	s.invalidateSearchCache(ctx)

	// Broadcast WebSocket update
	if s.hub != nil {
		s.hub.BroadcastUpdate("member_created", member)
	}

	return nil
}

// GetByID retrieves a member by ID with caching
func (s *MemberService) GetByID(ctx context.Context, id string) (*models.Member, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("member:%s", id)
	cached, err := s.cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var member models.Member
		if err := json.Unmarshal([]byte(cached), &member); err == nil {
			return &member, nil
		}
	}

	// Cache miss - query database
	member, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	// Store in cache
	memberJSON, err := json.Marshal(member)
	if err == nil {
		s.cache.Client.Set(ctx, cacheKey, string(memberJSON), MemberCacheTTL)
	}

	return member, nil
}

// GetAll retrieves all non-deleted members
func (s *MemberService) GetAll(ctx context.Context) ([]*models.Member, error) {
	members, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all members: %w", err)
	}

	return members, nil
}

// Update updates an existing member and invalidates relevant caches
func (s *MemberService) Update(ctx context.Context, member *models.Member) error {
	// Validate member data
	if err := member.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update timestamp
	member.Update()

	// Update member in database
	if err := s.repo.Update(ctx, member); err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}

	// Invalidate member cache
	cacheKey := fmt.Sprintf("member:%s", member.ID)
	s.cache.Client.Del(ctx, cacheKey)

	// Invalidate family tree and search caches
	s.invalidateFamilyTreeCache(ctx)
	s.invalidateSearchCache(ctx)

	// Broadcast WebSocket update
	if s.hub != nil {
		s.hub.BroadcastUpdate("member_updated", member)
	}

	return nil
}

// SoftDelete marks a member as deleted and invalidates relevant caches
func (s *MemberService) SoftDelete(ctx context.Context, id string) error {
	// Soft delete member in database
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return fmt.Errorf("failed to soft delete member: %w", err)
	}

	// Invalidate member cache
	cacheKey := fmt.Sprintf("member:%s", id)
	s.cache.Client.Del(ctx, cacheKey)

	// Invalidate family tree and search caches
	s.invalidateFamilyTreeCache(ctx)
	s.invalidateSearchCache(ctx)

	// Broadcast WebSocket update
	if s.hub != nil {
		s.hub.BroadcastUpdate("member_deleted", map[string]string{"id": id})
	}

	return nil
}

// Search searches for members by name (case-insensitive partial match) with caching
func (s *MemberService) Search(ctx context.Context, query string) ([]*models.Member, error) {
	if query == "" {
		return s.GetAll(ctx)
	}

	// Normalize query for cache key
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))

	// Try cache first
	cacheKey := fmt.Sprintf("search:members:%s", normalizedQuery)
	cached, err := s.cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var members []*models.Member
		if err := json.Unmarshal([]byte(cached), &members); err == nil {
			return members, nil
		}
	}

	// Cache miss - query database and filter
	allMembers, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search members: %w", err)
	}

	// Filter members by name (case-insensitive partial match)
	var results []*models.Member
	for _, member := range allMembers {
		if strings.Contains(strings.ToLower(member.Name), normalizedQuery) {
			results = append(results, member)
		}
	}

	// Store in cache
	resultsJSON, err := json.Marshal(results)
	if err == nil {
		s.cache.Client.Set(ctx, cacheKey, string(resultsJSON), SearchCacheTTL)
	}

	return results, nil
}

// invalidateFamilyTreeCache invalidates all family tree cache entries
func (s *MemberService) invalidateFamilyTreeCache(ctx context.Context) {
	// Use SCAN to find all family_tree:* keys
	iter := s.cache.Client.Scan(ctx, 0, "family_tree:*", 0).Iterator()
	for iter.Next(ctx) {
		s.cache.Client.Del(ctx, iter.Val())
	}
}

// invalidateSearchCache invalidates all search cache entries
func (s *MemberService) invalidateSearchCache(ctx context.Context) {
	// Use SCAN to find all search:members:* keys
	iter := s.cache.Client.Scan(ctx, 0, "search:members:*", 0).Iterator()
	for iter.Next(ctx) {
		s.cache.Client.Del(ctx, iter.Val())
	}
}
