package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the VamsaSetu application
type Config struct {
	// Database Configuration
	PostgresURL    string
	Neo4jURI       string
	Neo4jUsername  string
	Neo4jPassword  string
	RedisAddr      string

	// Authentication & Security
	JWTSecret string

	// Notification Services
	SendGridAPIKey       string
	TwilioAccountSID     string
	TwilioAuthToken      string
	TwilioPhoneNumber    string
	TwilioWhatsAppNumber string

	// Application Configuration (optional with defaults)
	Port        string
	Environment string
}

// Load reads configuration from environment variables and validates required fields
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		// Database Configuration
		PostgresURL:   os.Getenv("POSTGRES_URL"),
		Neo4jURI:      os.Getenv("NEO4J_URI"),
		Neo4jUsername: os.Getenv("NEO4J_USERNAME"),
		Neo4jPassword: os.Getenv("NEO4J_PASSWORD"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),

		// Authentication & Security
		JWTSecret: os.Getenv("JWT_SECRET"),

		// Notification Services
		SendGridAPIKey:       os.Getenv("SENDGRID_API_KEY"),
		TwilioAccountSID:     os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:      os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioPhoneNumber:    os.Getenv("TWILIO_PHONE_NUMBER"),
		TwilioWhatsAppNumber: os.Getenv("TWILIO_WHATSAPP_NUMBER"),

		// Application Configuration (with defaults)
		Port:        getEnvOrDefault("PORT", "8080"),
		Environment: getEnvOrDefault("ENV", "development"),
	}

	// Validate required environment variables
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that all required environment variables are present
func (c *Config) Validate() error {
	var missingVars []string

	// Check required database variables
	if c.PostgresURL == "" {
		missingVars = append(missingVars, "POSTGRES_URL")
	}
	if c.Neo4jURI == "" {
		missingVars = append(missingVars, "NEO4J_URI")
	}
	if c.Neo4jUsername == "" {
		missingVars = append(missingVars, "NEO4J_USERNAME")
	}
	if c.Neo4jPassword == "" {
		missingVars = append(missingVars, "NEO4J_PASSWORD")
	}
	if c.RedisAddr == "" {
		missingVars = append(missingVars, "REDIS_ADDR")
	}

	// Check required authentication variables
	if c.JWTSecret == "" {
		missingVars = append(missingVars, "JWT_SECRET")
	}

	// Check required notification service variables
	if c.SendGridAPIKey == "" {
		missingVars = append(missingVars, "SENDGRID_API_KEY")
	}
	if c.TwilioAccountSID == "" {
		missingVars = append(missingVars, "TWILIO_ACCOUNT_SID")
	}
	if c.TwilioAuthToken == "" {
		missingVars = append(missingVars, "TWILIO_AUTH_TOKEN")
	}
	if c.TwilioPhoneNumber == "" {
		missingVars = append(missingVars, "TWILIO_PHONE_NUMBER")
	}
	if c.TwilioWhatsAppNumber == "" {
		missingVars = append(missingVars, "TWILIO_WHATSAPP_NUMBER")
	}

	// Return descriptive error if any variables are missing
	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %v. Please ensure all required variables are set in your .env file or environment", missingVars)
	}

	return nil
}

// getEnvOrDefault returns the environment variable value or a default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
