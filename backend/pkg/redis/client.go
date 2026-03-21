package redis

import (
	"context"
	"fmt"

	"vamsasetu/backend/internal/config"

	"github.com/redis/go-redis/v9"
)

// Client wraps the Redis client
type Client struct {
	Client *redis.Client
}

// NewClient initializes a new Redis client with the provided configuration
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.RedisAddr == "" {
		return nil, fmt.Errorf("redis address is required")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	// Verify connectivity
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{
		Client: rdb,
	}, nil
}

// HealthCheck verifies the Redis connection is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	if c.Client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	if err := c.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

// Close closes the Redis client connection
func (c *Client) Close() error {
	if c.Client != nil {
		return c.Client.Close()
	}
	return nil
}
