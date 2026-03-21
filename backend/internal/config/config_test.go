package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoad_AllRequiredVariablesPresent(t *testing.T) {
	// Set all required environment variables
	setTestEnvVars()
	defer clearTestEnvVars()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify all fields are populated
	if cfg.PostgresURL != "postgresql://test:test@localhost:5432/test" {
		t.Errorf("Expected PostgresURL to be set, got: %s", cfg.PostgresURL)
	}
	if cfg.Neo4jURI != "bolt://localhost:7687" {
		t.Errorf("Expected Neo4jURI to be set, got: %s", cfg.Neo4jURI)
	}
	if cfg.JWTSecret != "test-secret" {
		t.Errorf("Expected JWTSecret to be set, got: %s", cfg.JWTSecret)
	}
}

func TestLoad_MissingRequiredVariables(t *testing.T) {
	// Clear all environment variables
	clearTestEnvVars()

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error for missing variables, got nil")
	}

	// Verify error message contains information about missing variables
	errMsg := err.Error()
	if !strings.Contains(errMsg, "missing required environment variables") {
		t.Errorf("Expected error message to mention missing variables, got: %s", errMsg)
	}
}

func TestLoad_MissingSingleVariable(t *testing.T) {
	// Set all required variables except one
	setTestEnvVars()
	os.Unsetenv("JWT_SECRET")
	defer clearTestEnvVars()

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error for missing JWT_SECRET, got nil")
	}

	// Verify error message mentions JWT_SECRET
	errMsg := err.Error()
	if !strings.Contains(errMsg, "JWT_SECRET") {
		t.Errorf("Expected error message to mention JWT_SECRET, got: %s", errMsg)
	}
}

func TestLoad_OptionalVariablesHaveDefaults(t *testing.T) {
	// Set only required variables
	setTestEnvVars()
	defer clearTestEnvVars()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify optional variables have defaults
	if cfg.Port != "8080" {
		t.Errorf("Expected default Port to be 8080, got: %s", cfg.Port)
	}
	if cfg.Environment != "development" {
		t.Errorf("Expected default Environment to be development, got: %s", cfg.Environment)
	}
}

func TestLoad_OptionalVariablesCanBeOverridden(t *testing.T) {
	// Set all required variables and override optional ones
	setTestEnvVars()
	os.Setenv("PORT", "3000")
	os.Setenv("ENV", "production")
	defer clearTestEnvVars()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify optional variables are overridden
	if cfg.Port != "3000" {
		t.Errorf("Expected Port to be 3000, got: %s", cfg.Port)
	}
	if cfg.Environment != "production" {
		t.Errorf("Expected Environment to be production, got: %s", cfg.Environment)
	}
}

func TestValidate_AllFieldsPresent(t *testing.T) {
	cfg := &Config{
		PostgresURL:          "postgresql://test:test@localhost:5432/test",
		Neo4jURI:             "bolt://localhost:7687",
		Neo4jUsername:        "neo4j",
		Neo4jPassword:        "password",
		RedisAddr:            "localhost:6379",
		JWTSecret:            "secret",
		SendGridAPIKey:       "SG.test",
		TwilioAccountSID:     "AC123",
		TwilioAuthToken:      "token",
		TwilioPhoneNumber:    "+1234567890",
		TwilioWhatsAppNumber: "whatsapp:+1234567890",
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestValidate_MissingDatabaseFields(t *testing.T) {
	testCases := []struct {
		name      string
		cfg       *Config
		missingVar string
	}{
		{
			name: "Missing PostgresURL",
			cfg: &Config{
				Neo4jURI:             "bolt://localhost:7687",
				Neo4jUsername:        "neo4j",
				Neo4jPassword:        "password",
				RedisAddr:            "localhost:6379",
				JWTSecret:            "secret",
				SendGridAPIKey:       "SG.test",
				TwilioAccountSID:     "AC123",
				TwilioAuthToken:      "token",
				TwilioPhoneNumber:    "+1234567890",
				TwilioWhatsAppNumber: "whatsapp:+1234567890",
			},
			missingVar: "POSTGRES_URL",
		},
		{
			name: "Missing Neo4jURI",
			cfg: &Config{
				PostgresURL:          "postgresql://test:test@localhost:5432/test",
				Neo4jUsername:        "neo4j",
				Neo4jPassword:        "password",
				RedisAddr:            "localhost:6379",
				JWTSecret:            "secret",
				SendGridAPIKey:       "SG.test",
				TwilioAccountSID:     "AC123",
				TwilioAuthToken:      "token",
				TwilioPhoneNumber:    "+1234567890",
				TwilioWhatsAppNumber: "whatsapp:+1234567890",
			},
			missingVar: "NEO4J_URI",
		},
		{
			name: "Missing RedisAddr",
			cfg: &Config{
				PostgresURL:          "postgresql://test:test@localhost:5432/test",
				Neo4jURI:             "bolt://localhost:7687",
				Neo4jUsername:        "neo4j",
				Neo4jPassword:        "password",
				JWTSecret:            "secret",
				SendGridAPIKey:       "SG.test",
				TwilioAccountSID:     "AC123",
				TwilioAuthToken:      "token",
				TwilioPhoneNumber:    "+1234567890",
				TwilioWhatsAppNumber: "whatsapp:+1234567890",
			},
			missingVar: "REDIS_ADDR",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.missingVar) {
				t.Errorf("Expected error to mention %s, got: %s", tc.missingVar, err.Error())
			}
		})
	}
}

func TestValidate_MissingNotificationFields(t *testing.T) {
	cfg := &Config{
		PostgresURL:   "postgresql://test:test@localhost:5432/test",
		Neo4jURI:      "bolt://localhost:7687",
		Neo4jUsername: "neo4j",
		Neo4jPassword: "password",
		RedisAddr:     "localhost:6379",
		JWTSecret:     "secret",
		// Missing all notification service fields
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Verify all notification fields are mentioned
	errMsg := err.Error()
	requiredNotificationVars := []string{
		"SENDGRID_API_KEY",
		"TWILIO_ACCOUNT_SID",
		"TWILIO_AUTH_TOKEN",
		"TWILIO_PHONE_NUMBER",
		"TWILIO_WHATSAPP_NUMBER",
	}

	for _, varName := range requiredNotificationVars {
		if !strings.Contains(errMsg, varName) {
			t.Errorf("Expected error to mention %s, got: %s", varName, errMsg)
		}
	}
}

// Helper functions

func setTestEnvVars() {
	os.Setenv("POSTGRES_URL", "postgresql://test:test@localhost:5432/test")
	os.Setenv("NEO4J_URI", "bolt://localhost:7687")
	os.Setenv("NEO4J_USERNAME", "neo4j")
	os.Setenv("NEO4J_PASSWORD", "password")
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("SENDGRID_API_KEY", "SG.test")
	os.Setenv("TWILIO_ACCOUNT_SID", "AC123")
	os.Setenv("TWILIO_AUTH_TOKEN", "token")
	os.Setenv("TWILIO_PHONE_NUMBER", "+1234567890")
	os.Setenv("TWILIO_WHATSAPP_NUMBER", "whatsapp:+1234567890")
}

func clearTestEnvVars() {
	os.Unsetenv("POSTGRES_URL")
	os.Unsetenv("NEO4J_URI")
	os.Unsetenv("NEO4J_USERNAME")
	os.Unsetenv("NEO4J_PASSWORD")
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("SENDGRID_API_KEY")
	os.Unsetenv("TWILIO_ACCOUNT_SID")
	os.Unsetenv("TWILIO_AUTH_TOKEN")
	os.Unsetenv("TWILIO_PHONE_NUMBER")
	os.Unsetenv("TWILIO_WHATSAPP_NUMBER")
	os.Unsetenv("PORT")
	os.Unsetenv("ENV")
}
