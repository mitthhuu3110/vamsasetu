// +build ignore

package pkg

import (
	"context"
	"fmt"
	"log"

	"vamsasetu/backend/internal/config"
	"vamsasetu/backend/pkg/neo4j"
	"vamsasetu/backend/pkg/postgres"
	"vamsasetu/backend/pkg/redis"
)

// ExampleUsage demonstrates how to initialize and use the database clients
func ExampleUsage() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Neo4j client
	neo4jClient, err := neo4j.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Neo4j client: %v", err)
	}
	defer neo4jClient.Close(context.Background())

	// Initialize PostgreSQL client
	postgresClient, err := postgres.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL client: %v", err)
	}
	defer postgresClient.Close()

	// Initialize Redis client
	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}
	defer redisClient.Close()

	// Perform health checks
	ctx := context.Background()

	if err := neo4jClient.HealthCheck(ctx); err != nil {
		log.Printf("Neo4j health check failed: %v", err)
	} else {
		fmt.Println("Neo4j is healthy")
	}

	if err := postgresClient.HealthCheck(ctx); err != nil {
		log.Printf("PostgreSQL health check failed: %v", err)
	} else {
		fmt.Println("PostgreSQL is healthy")
	}

	if err := redisClient.HealthCheck(ctx); err != nil {
		log.Printf("Redis health check failed: %v", err)
	} else {
		fmt.Println("Redis is healthy")
	}

	// Now you can use the clients:
	// - neo4jClient.Driver for Neo4j operations
	// - postgresClient.DB for GORM operations
	// - redisClient.Client for Redis operations
}
