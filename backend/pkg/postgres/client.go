package postgres

import (
	"context"
	"fmt"

	"vamsasetu/backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Client wraps the GORM database connection
type Client struct {
	DB *gorm.DB
}

// NewClient initializes a new PostgreSQL client with GORM
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.PostgresURL == "" {
		return nil, fmt.Errorf("postgres URL is required")
	}

	db, err := gorm.Open(postgres.Open(cfg.PostgresURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Get underlying SQL DB to verify connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Verify connectivity
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &Client{
		DB: db,
	}, nil
}

// HealthCheck verifies the PostgreSQL connection is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	if c.DB == nil {
		return fmt.Errorf("postgres DB is not initialized")
	}

	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("postgres health check failed: %w", err)
	}

	return nil
}

// Close closes the PostgreSQL database connection
func (c *Client) Close() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}
