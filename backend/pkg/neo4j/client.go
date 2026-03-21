package neo4j

import (
	"context"
	"fmt"

	"vamsasetu/backend/internal/config"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Client wraps the Neo4j driver
type Client struct {
	Driver neo4j.DriverWithContext
}

// NewClient initializes a new Neo4j client with the provided configuration
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.Neo4jURI == "" {
		return nil, fmt.Errorf("neo4j URI is required")
	}
	if cfg.Neo4jUsername == "" {
		return nil, fmt.Errorf("neo4j username is required")
	}
	if cfg.Neo4jPassword == "" {
		return nil, fmt.Errorf("neo4j password is required")
	}

	driver, err := neo4j.NewDriverWithContext(
		cfg.Neo4jURI,
		neo4j.BasicAuth(cfg.Neo4jUsername, cfg.Neo4jPassword, ""),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	// Verify connectivity
	ctx := context.Background()
	if err := driver.VerifyConnectivity(ctx); err != nil {
		driver.Close(ctx)
		return nil, fmt.Errorf("failed to verify neo4j connectivity: %w", err)
	}

	return &Client{
		Driver: driver,
	}, nil
}

// HealthCheck verifies the Neo4j connection is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	if c.Driver == nil {
		return fmt.Errorf("neo4j driver is not initialized")
	}

	if err := c.Driver.VerifyConnectivity(ctx); err != nil {
		return fmt.Errorf("neo4j health check failed: %w", err)
	}

	return nil
}

// Close closes the Neo4j driver connection
func (c *Client) Close(ctx context.Context) error {
	if c.Driver != nil {
		return c.Driver.Close(ctx)
	}
	return nil
}
