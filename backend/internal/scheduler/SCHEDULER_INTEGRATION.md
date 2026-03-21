# Notification Scheduler Integration Guide

## Overview

The notification scheduler is a background service that runs every hour to check for pending notifications and dispatch them via WhatsApp, SMS, or Email. It uses a goroutine pool with a maximum of 10 concurrent workers to handle notification dispatch efficiently.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   Notification Scheduler                     │
│                                                              │
│  ┌──────────────┐      ┌──────────────────────────────┐   │
│  │   Ticker     │─────▶│   processNotifications()     │   │
│  │  (1 hour)    │      │                              │   │
│  └──────────────┘      │  1. Query pending notifs     │   │
│                        │  2. Dispatch via worker pool │   │
│                        │  3. Update status            │   │
│                        └──────────────────────────────┘   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │           Worker Pool (10 concurrent)                 │  │
│  │                                                       │  │
│  │  Worker 1 ─▶ Dispatch ─▶ Update Status              │  │
│  │  Worker 2 ─▶ Dispatch ─▶ Update Status              │  │
│  │  ...                                                  │  │
│  │  Worker 10 ─▶ Dispatch ─▶ Update Status             │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Usage in main.go

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/internal/repository"
    "vamsasetu/backend/internal/scheduler"
    "vamsasetu/backend/internal/service"
    "vamsasetu/backend/pkg/postgres"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Initialize PostgreSQL connection
    db, err := postgres.NewClient(cfg.PostgresURL)
    if err != nil {
        log.Fatalf("Failed to connect to PostgreSQL: %v", err)
    }
    
    // Initialize repositories
    notifRepo := repository.NewNotificationRepository(db)
    
    // Initialize notification service
    notifSvc := service.NewNotificationService(
        notifRepo,
        cfg.SendGridAPIKey,
        cfg.TwilioAccountSID,
        cfg.TwilioAuthToken,
        cfg.TwilioPhoneNumber,
        cfg.TwilioWhatsApp,
    )
    
    // Initialize and start notification scheduler
    notifScheduler := scheduler.NewNotificationScheduler(notifRepo, notifSvc)
    notifScheduler.Start()
    log.Println("Notification scheduler started")
    
    // Set up graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    // Wait for shutdown signal
    <-sigChan
    log.Println("Shutting down...")
    
    // Stop scheduler
    notifScheduler.Stop()
    log.Println("Notification scheduler stopped")
}
```

## Features

### 1. Hourly Execution

The scheduler runs every hour using `time.Ticker`:

```go
s.ticker = time.NewTicker(1 * time.Hour)
```

It also runs immediately on start to process any pending notifications.

### 2. Worker Pool (10 Concurrent Workers)

The scheduler uses a buffered channel to limit concurrent dispatches:

```go
s.workerPool = make(chan struct{}, 10) // Max 10 concurrent workers
```

This prevents overwhelming external services (Twilio, SendGrid) with too many concurrent requests.

### 3. Automatic Retry with Exponential Backoff

When a notification fails to dispatch:
- Retry count is incremented
- If retry count < 3: notification remains pending for next scheduler run
- If retry count >= 3: notification is marked as "failed"

Backoff schedule:
- 1st retry: 5 minutes (3^1 * 5)
- 2nd retry: 15 minutes (3^2 * 5)
- 3rd retry: 45 minutes (3^3 * 5)

### 4. Graceful Shutdown

The scheduler supports graceful shutdown via context cancellation:

```go
scheduler.Stop() // Cancels context and stops ticker
```

## Notification Status Flow

```
┌─────────┐
│ pending │ ◀─── Initial state when notification is created
└────┬────┘
     │
     ├─── Dispatch Success ───▶ ┌──────┐
     │                          │ sent │
     │                          └──────┘
     │
     └─── Dispatch Failure ───▶ ┌─────────┐
                                 │ pending │ (retry_count++)
                                 └────┬────┘
                                      │
                                      ├─── retry_count < 3 ───▶ Try again
                                      │
                                      └─── retry_count >= 3 ───▶ ┌────────┐
                                                                  │ failed │
                                                                  └────────┘
```

## Database Queries

### Query Pending Notifications

```sql
SELECT * FROM notifications
WHERE status = 'pending' AND scheduled_at <= NOW()
ORDER BY scheduled_at ASC
```

### Update Status to Sent

```sql
UPDATE notifications
SET status = 'sent', sent_at = NOW(), error_msg = '', updated_at = NOW()
WHERE id = ?
```

### Update Status to Failed

```sql
UPDATE notifications
SET status = 'failed', error_msg = ?, updated_at = NOW()
WHERE id = ?
```

### Increment Retry Count

```sql
UPDATE notifications
SET retry_count = retry_count + 1
WHERE id = ?
```

## Testing

Run the scheduler tests:

```bash
cd backend
go test -v ./internal/scheduler/...
```

Test coverage includes:
- Scheduler start/stop
- Processing pending notifications
- Worker pool concurrency (15 notifications with 10 workers)
- Graceful shutdown

## Monitoring

The scheduler logs important events:

```
Notification scheduler started
Processing 5 pending notifications
Dispatching notification ID=1, channel=email
Successfully dispatched notification ID=1
Failed to dispatch notification ID=2: connection timeout
Scheduling retry for notification ID=2 in 5 minutes
Notification ID=3 failed after 3 retries
Notification scheduler stopped
```

## Performance Considerations

1. **Hourly Execution**: Reduces database load while ensuring timely delivery
2. **Worker Pool**: Prevents overwhelming external APIs
3. **Exponential Backoff**: Gives temporary failures time to resolve
4. **Goroutines**: Non-blocking dispatch allows processing multiple notifications simultaneously

## Requirements Satisfied

- ✅ 6.2: Scheduler runs every hour
- ✅ 6.4: Uses goroutine pool with max 10 concurrent workers
- ✅ 6.5: Updates notification status (sent/failed)
- ✅ 6.6: Implements retry logic with exponential backoff (max 3 retries)
- ✅ Graceful shutdown with context cancellation
- ✅ Comprehensive logging for monitoring

## Next Steps

1. Implement actual Twilio and SendGrid integrations in `notification_service.go`
2. Add metrics/monitoring (e.g., Prometheus counters for sent/failed notifications)
3. Consider adding a manual trigger endpoint for immediate notification processing
4. Add alerting for high failure rates
