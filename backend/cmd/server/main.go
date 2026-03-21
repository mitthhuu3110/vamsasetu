package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vamsasetu/backend/internal/config"
	"vamsasetu/backend/internal/handler"
	"vamsasetu/backend/internal/middleware"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/internal/scheduler"
	"vamsasetu/backend/internal/service"
	"vamsasetu/backend/pkg/neo4j"
	"vamsasetu/backend/pkg/postgres"
	redisClient "vamsasetu/backend/pkg/redis"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("VamsaSetu Backend Server - Starting...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// Initialize database clients
	pgClient, neo4jClient, redisClientInstance, err := initializeDatabases(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer pgClient.Close()
	defer neo4jClient.Close(context.Background())
	defer redisClientInstance.Close()

	log.Println("All database connections established")

	// Initialize WebSocket hub
	hub := handler.NewHub()
	go hub.Run()
	log.Println("WebSocket hub started")

	// Initialize repositories
	userRepo := repository.NewUserRepository(pgClient.DB)
	memberRepo := repository.NewMemberRepository(neo4jClient)
	relationshipRepo := repository.NewRelationshipRepository(neo4jClient)
	eventRepo := repository.NewEventRepository(pgClient.DB)
	notificationRepo := repository.NewNotificationRepository(pgClient.DB)

	log.Println("Repositories initialized")

	// Initialize services
	authService := service.NewAuthService(userRepo, redisClientInstance.Client)
	memberService := service.NewMemberService(memberRepo, redisClientInstance, hub)
	relationshipService := service.NewRelationshipService(relationshipRepo, hub)
	eventService := service.NewEventService(eventRepo, redisClientInstance, hub)
	notificationService := service.NewNotificationService(
		notificationRepo,
		cfg.SendGridAPIKey,
		cfg.TwilioAccountSID,
		cfg.TwilioAuthToken,
		cfg.TwilioPhoneNumber,
		cfg.TwilioWhatsAppNumber,
	)
	treeBuilder := service.NewTreeBuilder(memberRepo, relationshipRepo, eventRepo)

	log.Println("Services initialized")

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	memberHandler := handler.NewMemberHandler(memberService)
	relationshipHandler := handler.NewRelationshipHandler(relationshipService)
	eventHandler := handler.NewEventHandler(eventService)
	familyHandler := handler.NewFamilyHandler(treeBuilder, redisClientInstance.Client)

	log.Println("Handlers initialized")

	// Initialize notification scheduler
	notifScheduler := scheduler.NewNotificationScheduler(notificationRepo, notificationService)
	notifScheduler.Start()
	log.Println("Notification scheduler started")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		AppName:      "VamsaSetu API",
	})

	// Apply middleware in correct order
	// 1. Logger middleware (logs all requests)
	app.Use(middleware.LoggerMiddleware())

	// 2. CORS middleware (allows frontend origin)
	frontendOrigin := getEnv("FRONTEND_ORIGIN", "http://localhost:3000")
	app.Use(middleware.CORSMiddleware(frontendOrigin))

	// Register routes
	// Public routes (no auth required)
	authHandler.RegisterRoutes(app)

	// Protected routes (auth required, applied in handlers)
	memberHandler.RegisterRoutes(app)
	relationshipHandler.RegisterRoutes(app)
	eventHandler.RegisterRoutes(app)
	familyHandler.RegisterRoutes(app)

	// WebSocket endpoint
	app.Get("/ws", handler.HandleWebSocket(hub))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		ctx := context.Background()

		// Check PostgreSQL health
		pgHealth := "healthy"
		if err := pgClient.HealthCheck(ctx); err != nil {
			pgHealth = fmt.Sprintf("unhealthy: %v", err)
		}

		// Check Neo4j health
		neo4jHealth := "healthy"
		if err := neo4jClient.HealthCheck(ctx); err != nil {
			neo4jHealth = fmt.Sprintf("unhealthy: %v", err)
		}

		// Check Redis health
		redisHealth := "healthy"
		if err := redisClientInstance.HealthCheck(ctx); err != nil {
			redisHealth = fmt.Sprintf("unhealthy: %v", err)
		}

		// Determine overall health
		overallHealth := "healthy"
		if pgHealth != "healthy" || neo4jHealth != "healthy" || redisHealth != "healthy" {
			overallHealth = "degraded"
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": overallHealth,
			"services": fiber.Map{
				"postgres": pgHealth,
				"neo4j":    neo4jHealth,
				"redis":    redisHealth,
			},
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Start server in a goroutine
	port := cfg.Port
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop notification scheduler
	notifScheduler.Stop()
	log.Println("Notification scheduler stopped")

	// Shutdown Fiber app with timeout
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("VamsaSetu Backend Server - Stopped")
}

// initializeDatabases initializes PostgreSQL, Neo4j, and Redis clients and runs migrations
func initializeDatabases(cfg *config.Config) (*postgres.Client, *neo4j.Client, *redisClient.Client, error) {
	ctx := context.Background()

	// Initialize PostgreSQL client
	log.Println("Connecting to PostgreSQL...")
	pgClient, err := postgres.NewClient(cfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Verify PostgreSQL connection
	if err := pgClient.HealthCheck(ctx); err != nil {
		pgClient.Close()
		return nil, nil, nil, fmt.Errorf("PostgreSQL health check failed: %w", err)
	}
	log.Println("PostgreSQL connection established")

	// Run GORM auto-migration
	log.Println("Running GORM auto-migration...")
	if err := runMigrations(pgClient); err != nil {
		pgClient.Close()
		return nil, nil, nil, fmt.Errorf("migration failed: %w", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize Neo4j client
	log.Println("Connecting to Neo4j...")
	neo4jClient, err := neo4j.NewClient(cfg)
	if err != nil {
		pgClient.Close()
		return nil, nil, nil, fmt.Errorf("failed to connect to Neo4j: %w", err)
	}

	// Verify Neo4j connection
	if err := neo4jClient.HealthCheck(ctx); err != nil {
		pgClient.Close()
		neo4jClient.Close(ctx)
		return nil, nil, nil, fmt.Errorf("Neo4j health check failed: %w", err)
	}
	log.Println("Neo4j connection established")

	// Initialize Redis client
	log.Println("Connecting to Redis...")
	redisClientInstance, err := redisClient.NewClient(cfg)
	if err != nil {
		pgClient.Close()
		neo4jClient.Close(ctx)
		return nil, nil, nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Verify Redis connection
	if err := redisClientInstance.HealthCheck(ctx); err != nil {
		pgClient.Close()
		neo4jClient.Close(ctx)
		redisClientInstance.Close()
		return nil, nil, nil, fmt.Errorf("Redis health check failed: %w", err)
	}
	log.Println("Redis connection established")

	return pgClient, neo4jClient, redisClientInstance, nil
}

// runMigrations executes GORM AutoMigrate for all models
func runMigrations(pgClient *postgres.Client) error {
	// Auto-migrate User first
	log.Println("Migrating User table...")
	if err := pgClient.DB.AutoMigrate(&models.User{}); err != nil {
		log.Printf("User migration failed: %v", err)
		return err
	}
	log.Println("✓ Migrated User table")

	// Auto-migrate Event
	log.Println("Migrating Event table...")
	if err := pgClient.DB.AutoMigrate(&models.Event{}); err != nil {
		log.Printf("Event migration failed: %v", err)
		return err
	}
	log.Println("✓ Migrated Event table")

	// Auto-migrate Notification
	log.Println("Migrating Notification table...")
	if err := pgClient.DB.AutoMigrate(&models.Notification{}); err != nil {
		log.Printf("Notification migration failed: %v", err)
		return err
	}
	log.Println("✓ Migrated Notification table")

	// Auto-migrate AuditLog
	log.Println("Migrating AuditLog table...")
	if err := pgClient.DB.AutoMigrate(&models.AuditLog{}); err != nil {
		log.Printf("AuditLog migration failed: %v", err)
		return err
	}
	log.Println("✓ Migrated AuditLog table")

	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
