# Configuration Management

This package provides configuration management for the VamsaSetu backend application.

## Features

- Loads configuration from environment variables using `godotenv`
- Validates all required environment variables on startup
- Provides descriptive error messages for missing variables
- Supports optional variables with sensible defaults

## Usage

### Basic Usage

```go
package main

import (
    "log"
    "vamsasetu/backend/internal/config"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Use configuration
    log.Printf("Starting server on port %s", cfg.Port)
    log.Printf("Environment: %s", cfg.Environment)
    
    // Access database credentials
    log.Printf("PostgreSQL URL: %s", cfg.PostgresURL)
    log.Printf("Neo4j URI: %s", cfg.Neo4jURI)
}
```

## Required Environment Variables

The following environment variables **must** be set:

### Database Configuration
- `POSTGRES_URL` - PostgreSQL connection string (e.g., `postgresql://user:pass@localhost:5432/dbname`)
- `NEO4J_URI` - Neo4j connection URI (e.g., `bolt://localhost:7687`)
- `NEO4J_USERNAME` - Neo4j username
- `NEO4J_PASSWORD` - Neo4j password
- `REDIS_ADDR` - Redis server address (e.g., `localhost:6379`)

### Authentication & Security
- `JWT_SECRET` - Secret key for JWT token signing

### Notification Services
- `SENDGRID_API_KEY` - SendGrid API key for email notifications
- `TWILIO_ACCOUNT_SID` - Twilio account SID
- `TWILIO_AUTH_TOKEN` - Twilio authentication token
- `TWILIO_PHONE_NUMBER` - Twilio phone number for SMS (E.164 format)
- `TWILIO_WHATSAPP_NUMBER` - Twilio WhatsApp number (e.g., `whatsapp:+1234567890`)

## Optional Environment Variables

These variables have default values:

- `PORT` - Server port (default: `8080`)
- `ENV` - Environment mode (default: `development`)

## Configuration Struct

```go
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

    // Application Configuration
    Port        string
    Environment string
}
```

## Error Handling

If any required environment variables are missing, `Load()` will return a descriptive error:

```
missing required environment variables: [JWT_SECRET SENDGRID_API_KEY]. 
Please ensure all required variables are set in your .env file or environment
```

## Testing

Run tests with:

```bash
go test -v ./internal/config/...
```

## Example .env File

See `.env.example` in the project root for a complete example of all required variables.
