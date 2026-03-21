package neo4j

import (
	"context"
	"testing"

	"vamsasetu/backend/internal/config"
)

func TestNewClient_MissingURI(t *testing.T) {
	cfg := &config.Config{
		Neo4jURI:      "",
		Neo4jUsername: "neo4j",
		Neo4jPassword: "password",
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Error("expected error for missing Neo4j URI, got nil")
	}
}

func TestNewClient_MissingUsername(t *testing.T) {
	cfg := &config.Config{
		Neo4jURI:      "bolt://localhost:7687",
		Neo4jUsername: "",
		Neo4jPassword: "password",
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Error("expected error for missing Neo4j username, got nil")
	}
}

func TestNewClient_MissingPassword(t *testing.T) {
	cfg := &config.Config{
		Neo4jURI:      "bolt://localhost:7687",
		Neo4jUsername: "neo4j",
		Neo4jPassword: "",
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Error("expected error for missing Neo4j password, got nil")
	}
}

func TestHealthCheck_NilDriver(t *testing.T) {
	client := &Client{Driver: nil}
	ctx := context.Background()

	err := client.HealthCheck(ctx)
	if err == nil {
		t.Error("expected error for nil driver, got nil")
	}
}
