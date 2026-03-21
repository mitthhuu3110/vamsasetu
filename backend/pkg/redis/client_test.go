package redis

import (
	"context"
	"testing"

	"vamsasetu/backend/internal/config"
)

func TestNewClient_MissingAddr(t *testing.T) {
	cfg := &config.Config{
		RedisAddr: "",
	}

	_, err := NewClient(cfg)
	if err == nil {
		t.Error("expected error for missing Redis address, got nil")
	}
}

func TestHealthCheck_NilClient(t *testing.T) {
	client := &Client{Client: nil}
	ctx := context.Background()

	err := client.HealthCheck(ctx)
	if err == nil {
		t.Error("expected error for nil client, got nil")
	}
}
