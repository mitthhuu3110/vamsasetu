package service

import (
	"context"
	"fmt"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/internal/utils"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    *repository.UserRepository
	redisClient *redis.Client
}

// NewAuthService creates a new authentication service instance
func NewAuthService(userRepo *repository.UserRepository, redisClient *redis.Client) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"` // defaults to "owner" if not provided
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response with tokens
type AuthResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         *models.User `json:"user"`
}

// Register creates a new user account with hashed password
// Validates email uniqueness and hashes password with bcrypt cost factor 10
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, fmt.Errorf("email, password, and name are required")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password with bcrypt cost factor 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "owner"
	}

	// Validate role
	if role != "owner" && role != "viewer" && role != "admin" {
		return nil, fmt.Errorf("invalid role: must be owner, viewer, or admin")
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Role:         role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in Redis with 7-day TTL
	refreshTokenKey := fmt.Sprintf("refresh_token:%d", user.ID)
	if err := s.redisClient.Set(ctx, refreshTokenKey, refreshToken, 7*24*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// Login validates user credentials and returns JWT tokens
// Compares provided password with stored bcrypt hash
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Compare password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in Redis with 7-day TTL
	refreshTokenKey := fmt.Sprintf("refresh_token:%d", user.ID)
	if err := s.redisClient.Set(ctx, refreshTokenKey, refreshToken, 7*24*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken validates a refresh token and issues a new access token
// Verifies the refresh token exists in Redis before issuing new access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	// Parse and validate refresh token
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Verify token type
	if claims.Type != "refresh" {
		return nil, fmt.Errorf("invalid token type: expected refresh token")
	}

	// Check if refresh token exists in Redis
	refreshTokenKey := fmt.Sprintf("refresh_token:%d", claims.UserID)
	storedToken, err := s.redisClient.Get(ctx, refreshTokenKey).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("refresh token not found or expired")
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve refresh token: %w", err)
	}

	// Verify stored token matches provided token
	if storedToken != refreshToken {
		return nil, fmt.Errorf("refresh token mismatch")
	}

	// Get user details
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // Return the same refresh token
		User:         user,
	}, nil
}

// GetUserByID retrieves a user by their ID
// Used by the profile endpoint to fetch authenticated user details
func (s *AuthService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}
