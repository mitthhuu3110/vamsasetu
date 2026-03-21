package config

import (
	"os"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
)

// Feature: vamsasetu-full-system, Property 4: User Role Invariant (partial - validates config loading)
// **Validates: Requirements 12.5**
//
// This property test validates that configuration loading correctly identifies missing environment variables
// and handles various combinations of present/absent variables. While Property 4 is about user roles,
// this partial implementation focuses on the configuration validation aspect that ensures the system
// fails fast when required variables are missing.
func TestProperty_ConfigurationValidation(t *testing.T) {
	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	// Property: For any subset of required environment variables that is missing at least one variable,
	// Load() should return an error mentioning the missing variables
	properties.Property("Missing required environment variables should cause Load() to fail with descriptive error",
		gen.SliceOfN(11, gen.Bool()).
			SuchThat(func(v interface{}) bool {
				// Ensure at least one variable is missing (not all true)
				bools := v.([]bool)
				allPresent := true
				for _, b := range bools {
					if !b {
						allPresent = false
						break
					}
				}
				return !allPresent
			}).
			Map(func(bools []bool) map[string]bool {
				// Map boolean array to variable presence map
				varNames := []string{
					"POSTGRES_URL",
					"NEO4J_URI",
					"NEO4J_USERNAME",
					"NEO4J_PASSWORD",
					"REDIS_ADDR",
					"JWT_SECRET",
					"SENDGRID_API_KEY",
					"TWILIO_ACCOUNT_SID",
					"TWILIO_AUTH_TOKEN",
					"TWILIO_PHONE_NUMBER",
					"TWILIO_WHATSAPP_NUMBER",
				}
				
				presence := make(map[string]bool)
				for i, varName := range varNames {
					presence[varName] = bools[i]
				}
				return presence
			}),
		func(presence map[string]bool) bool {
			// Clear all environment variables first
			clearAllConfigEnvVars()
			defer clearAllConfigEnvVars()

			// Set only the variables that should be present
			for varName, shouldBePresent := range presence {
				if shouldBePresent {
					os.Setenv(varName, getTestValueForVar(varName))
				}
			}

			// Attempt to load configuration
			_, err := Load()

			// Should always return an error since at least one variable is missing
			if err == nil {
				t.Logf("Expected error but got none. Presence map: %v", presence)
				return false
			}

			// Verify error message mentions all missing variables
			errMsg := err.Error()
			for varName, shouldBePresent := range presence {
				if !shouldBePresent {
					// This variable should be mentioned in the error
					if !containsString(errMsg, varName) {
						t.Logf("Error message should mention missing variable %s. Error: %s", varName, errMsg)
						return false
					}
				}
			}

			return true
		},
	)

	// Property: For any complete set of required environment variables, Load() should succeed
	properties.Property("All required environment variables present should cause Load() to succeed",
		gen.Const(true), // Dummy generator since we're testing a fixed scenario
		func(_ bool) bool {
			// Clear all environment variables first
			clearAllConfigEnvVars()
			defer clearAllConfigEnvVars()

			// Set all required variables
			setAllRequiredEnvVars()

			// Attempt to load configuration
			cfg, err := Load()

			// Should succeed
			if err != nil {
				t.Logf("Expected no error but got: %v", err)
				return false
			}

			// Verify all fields are populated correctly
			if cfg.PostgresURL != "postgresql://test:test@localhost:5432/test" {
				t.Logf("PostgresURL not set correctly: %s", cfg.PostgresURL)
				return false
			}
			if cfg.Neo4jURI != "bolt://localhost:7687" {
				t.Logf("Neo4jURI not set correctly: %s", cfg.Neo4jURI)
				return false
			}
			if cfg.Neo4jUsername != "neo4j" {
				t.Logf("Neo4jUsername not set correctly: %s", cfg.Neo4jUsername)
				return false
			}
			if cfg.Neo4jPassword != "password" {
				t.Logf("Neo4jPassword not set correctly: %s", cfg.Neo4jPassword)
				return false
			}
			if cfg.RedisAddr != "localhost:6379" {
				t.Logf("RedisAddr not set correctly: %s", cfg.RedisAddr)
				return false
			}
			if cfg.JWTSecret != "test-secret" {
				t.Logf("JWTSecret not set correctly: %s", cfg.JWTSecret)
				return false
			}
			if cfg.SendGridAPIKey != "SG.test" {
				t.Logf("SendGridAPIKey not set correctly: %s", cfg.SendGridAPIKey)
				return false
			}
			if cfg.TwilioAccountSID != "AC123" {
				t.Logf("TwilioAccountSID not set correctly: %s", cfg.TwilioAccountSID)
				return false
			}
			if cfg.TwilioAuthToken != "token" {
				t.Logf("TwilioAuthToken not set correctly: %s", cfg.TwilioAuthToken)
				return false
			}
			if cfg.TwilioPhoneNumber != "+1234567890" {
				t.Logf("TwilioPhoneNumber not set correctly: %s", cfg.TwilioPhoneNumber)
				return false
			}
			if cfg.TwilioWhatsAppNumber != "whatsapp:+1234567890" {
				t.Logf("TwilioWhatsAppNumber not set correctly: %s", cfg.TwilioWhatsAppNumber)
				return false
			}

			return true
		},
	)

	// Property: Optional variables should have defaults when not provided
	properties.Property("Optional environment variables should use defaults when not provided",
		gen.Const(true), // Dummy generator
		func(_ bool) bool {
			// Clear all environment variables first
			clearAllConfigEnvVars()
			defer clearAllConfigEnvVars()

			// Set only required variables (not optional ones)
			setAllRequiredEnvVars()
			os.Unsetenv("PORT")
			os.Unsetenv("ENV")

			// Attempt to load configuration
			cfg, err := Load()

			// Should succeed
			if err != nil {
				t.Logf("Expected no error but got: %v", err)
				return false
			}

			// Verify defaults are applied
			if cfg.Port != "8080" {
				t.Logf("Expected default Port to be 8080, got: %s", cfg.Port)
				return false
			}
			if cfg.Environment != "development" {
				t.Logf("Expected default Environment to be development, got: %s", cfg.Environment)
				return false
			}

			return true
		},
	)

	// Property: Optional variables can be overridden
	properties.Property("Optional environment variables can be overridden",
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) < 20 }).
			Map(func(s string) string { return s }),
		func(portValue string) bool {
			// Clear all environment variables first
			clearAllConfigEnvVars()
			defer clearAllConfigEnvVars()

			// Set all required variables and override PORT
			setAllRequiredEnvVars()
			os.Setenv("PORT", portValue)

			// Attempt to load configuration
			cfg, err := Load()

			// Should succeed
			if err != nil {
				t.Logf("Expected no error but got: %v", err)
				return false
			}

			// Verify PORT is overridden
			if cfg.Port != portValue {
				t.Logf("Expected Port to be %s, got: %s", portValue, cfg.Port)
				return false
			}

			return true
		},
	)

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Helper functions

func clearAllConfigEnvVars() {
	varNames := []string{
		"POSTGRES_URL",
		"NEO4J_URI",
		"NEO4J_USERNAME",
		"NEO4J_PASSWORD",
		"REDIS_ADDR",
		"JWT_SECRET",
		"SENDGRID_API_KEY",
		"TWILIO_ACCOUNT_SID",
		"TWILIO_AUTH_TOKEN",
		"TWILIO_PHONE_NUMBER",
		"TWILIO_WHATSAPP_NUMBER",
		"PORT",
		"ENV",
	}
	for _, varName := range varNames {
		os.Unsetenv(varName)
	}
}

func setAllRequiredEnvVars() {
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

func getTestValueForVar(varName string) string {
	testValues := map[string]string{
		"POSTGRES_URL":           "postgresql://test:test@localhost:5432/test",
		"NEO4J_URI":              "bolt://localhost:7687",
		"NEO4J_USERNAME":         "neo4j",
		"NEO4J_PASSWORD":         "password",
		"REDIS_ADDR":             "localhost:6379",
		"JWT_SECRET":             "test-secret",
		"SENDGRID_API_KEY":       "SG.test",
		"TWILIO_ACCOUNT_SID":     "AC123",
		"TWILIO_AUTH_TOKEN":      "token",
		"TWILIO_PHONE_NUMBER":    "+1234567890",
		"TWILIO_WHATSAPP_NUMBER": "whatsapp:+1234567890",
	}
	return testValues[varName]
}

func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
