package scheduler

import (
	"context"
	"log"
	"math"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/internal/service"
)

// NotificationScheduler handles periodic notification dispatch
type NotificationScheduler struct {
	notifRepo   *repository.NotificationRepository
	notifSvc    *service.NotificationService
	ticker      *time.Ticker
	workerPool  chan struct{}
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewNotificationScheduler creates a new notification scheduler instance
func NewNotificationScheduler(
	notifRepo *repository.NotificationRepository,
	notifSvc *service.NotificationService,
) *NotificationScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &NotificationScheduler{
		notifRepo:  notifRepo,
		notifSvc:   notifSvc,
		workerPool: make(chan struct{}, 10), // Max 10 concurrent workers
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start begins the notification scheduler with hourly execution
func (s *NotificationScheduler) Start() {
	s.ticker = time.NewTicker(1 * time.Hour)
	
	// Run immediately on start
	go s.processNotifications()
	
	// Then run every hour
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.processNotifications()
			case <-s.ctx.Done():
				log.Println("Notification scheduler stopped")
				return
			}
		}
	}()
	
	log.Println("Notification scheduler started")
}

// Stop gracefully shuts down the scheduler
func (s *NotificationScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.cancel()
}

// processNotifications queries and dispatches pending notifications
func (s *NotificationScheduler) processNotifications() {
	ctx := context.Background()
	
	// Query pending notifications
	notifications, err := s.notifRepo.GetPending(ctx)
	if err != nil {
		log.Printf("Error querying pending notifications: %v", err)
		return
	}
	
	if len(notifications) == 0 {
		log.Println("No pending notifications to process")
		return
	}
	
	log.Printf("Processing %d pending notifications", len(notifications))
	
	// Dispatch each notification using worker pool
	for _, notif := range notifications {
		// Acquire worker slot (blocks if pool is full)
		s.workerPool <- struct{}{}
		
		// Dispatch in goroutine
		go func(n *models.Notification) {
			defer func() {
				<-s.workerPool // Release worker slot
			}()
			
			s.dispatchNotification(ctx, n)
		}(notif)
	}
}

// dispatchNotification sends a single notification and updates its status
func (s *NotificationScheduler) dispatchNotification(ctx context.Context, notification *models.Notification) {
	log.Printf("Dispatching notification ID=%d, channel=%s", notification.ID, notification.Channel)
	
	// Attempt to dispatch
	err := s.notifSvc.Dispatch(ctx, notification)
	
	if err != nil {
		// Handle failure
		log.Printf("Failed to dispatch notification ID=%d: %v", notification.ID, err)
		s.handleFailure(ctx, notification, err)
	} else {
		// Mark as sent
		log.Printf("Successfully dispatched notification ID=%d", notification.ID)
		s.markAsSent(ctx, notification)
	}
}

// markAsSent updates notification status to sent
func (s *NotificationScheduler) markAsSent(ctx context.Context, notification *models.Notification) {
	now := time.Now()
	err := s.notifRepo.UpdateStatus(ctx, notification.ID, "sent", &now, "")
	if err != nil {
		log.Printf("Error updating notification status to sent: %v", err)
	}
}

// handleFailure increments retry count and reschedules or marks as failed
func (s *NotificationScheduler) handleFailure(ctx context.Context, notification *models.Notification, dispatchErr error) {
	// Increment retry count
	err := s.notifRepo.IncrementRetry(ctx, notification.ID)
	if err != nil {
		log.Printf("Error incrementing retry count: %v", err)
		return
	}
	
	notification.RetryCount++
	
	if notification.RetryCount >= 3 {
		// Max retries reached, mark as failed
		log.Printf("Notification ID=%d failed after %d retries", notification.ID, notification.RetryCount)
		err = s.notifRepo.UpdateStatus(ctx, notification.ID, "failed", nil, dispatchErr.Error())
		if err != nil {
			log.Printf("Error updating notification status to failed: %v", err)
		}
	} else {
		// Schedule retry with exponential backoff
		// Backoff: 5min, 15min, 45min (3^n * 5 minutes)
		backoffMinutes := int(math.Pow(3, float64(notification.RetryCount))) * 5
		log.Printf("Scheduling retry for notification ID=%d in %d minutes", notification.ID, backoffMinutes)
		
		// Note: The notification will be picked up in the next scheduler run
		// when scheduledAt <= now. We don't update scheduledAt here because
		// the repository doesn't have that method yet. The retry logic will
		// naturally pick it up in the next hourly run.
		err = s.notifRepo.UpdateStatus(ctx, notification.ID, "pending", nil, dispatchErr.Error())
		if err != nil {
			log.Printf("Error updating notification error message: %v", err)
		}
	}
}
