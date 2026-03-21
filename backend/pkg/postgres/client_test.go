package postgres

import (
	"context"
	"testing"

	"vamsasetu/backend/internal/config"
)

func TestNewClient_MissingURL(t *testing.T) {
	cfg := &config.Config{
		PostgresURL: "",
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Error("expected error for missing Postgres URL, got nil")
	}
}

func TestHealthCheck_NilDB(t *testing.T) {
	client := &Client{DB: nil}
	ctx := context.Background()

	err := client.HealthCheck(ctx)
	if err == nil {
		t.Error("expected error for nil DB, got nil")
	}
}
