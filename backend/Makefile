.PHONY: help install run build test clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install Go dependencies
	go mod download
	go mod tidy

run: ## Run the server
	go run cmd/server/main.go

build: ## Build the server binary
	go build -o bin/server cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy
