package repository

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"vamsasetu/backend/internal/models"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Validates: Requirements 5.1, 11.2**
// Property 21: Event Creation and Retrieval
// For any valid event data (title, date, type, member IDs), creating an event
// should result in that event being stored in PostgreSQL and retrievable by ID.
func TestProperty21_EventCreationAndRetrieval(t *testing.T) {
	repo, _, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any valid event data, creating an event should result in that event being retrievable by ID",
		prop.ForAll(
			func(title string, description string, eventDate time.Time, eventType string, memberIDs []string, createdBy uint) bool {
				// Create event with generated data
				event := &models.Event{
					Title:       title,
					Description: description,
					EventDate:   eventDate,
					EventType:   eventType,
					MemberIDs:   memberIDs,
					CreatedBy:   createdBy,
				}

				// Create in repository
				err := repo.Create(ctx, event)
				if err != nil {
					t.Logf("Failed to create event: %v", err)
					return false
				}

				// Verify ID was assigned
				if event.ID == 0 {
					t.Logf("Event ID should be assigned after creation")
					return false
				}

				// Retrieve by ID
				retrieved, err := repo.GetByID(ctx, event.ID)
				if err != nil {
					t.Logf("Failed to retrieve event: %v", err)
					return false
				}

				// Verify all fields match
				if retrieved.ID != event.ID {
					t.Logf("ID mismatch: expected %d, got %d", event.ID, retrieved.ID)
					return false
				}
				if retrieved.Title != event.Title {
					t.Logf("Title mismatch: expected %s, got %s", event.Title, retrieved.Title)
					return false
				}
				if retrieved.Description != event.Description {
					t.Logf("Description mismatch: expected %s, got %s", event.Description, retrieved.Description)
					return false
				}
				if !retrieved.EventDate.Equal(event.EventDate) {
					t.Logf("EventDate mismatch: expected %v, got %v", event.EventDate, retrieved.EventDate)
					return false
				}
				if retrieved.EventType != event.EventType {
					t.Logf("EventType mismatch: expected %s, got %s", event.EventType, retrieved.EventType)
					return false
				}
				if len(retrieved.MemberIDs) != len(event.MemberIDs) {
					t.Logf("MemberIDs length mismatch: expected %d, got %d", len(event.MemberIDs), len(retrieved.MemberIDs))
					return false
				}
				for i, id := range event.MemberIDs {
					if retrieved.MemberIDs[i] != id {
						t.Logf("MemberID mismatch at index %d: expected %s, got %s", i, id, retrieved.MemberIDs[i])
						return false
					}
				}
				if retrieved.CreatedBy != event.CreatedBy {
					t.Logf("CreatedBy mismatch: expected %d, got %d", event.CreatedBy, retrieved.CreatedBy)
					return false
				}

				// Cleanup
				repo.Delete(ctx, event.ID)

				return true
			},
			genValidEventTitle(),
			genValidEventDescription(),
			genValidEventDate(),
			genValidEventType(),
			genValidMemberIDs(),
			genValidUserID(),
		),
	)

	properties.TestingRun(t)
}

// **Validates: Requirements 5.2**
// Property 22: Event Update Persistence
// For any existing event and any valid attribute changes, updating the event
// should persist those changes such that subsequent retrieval returns the updated values.
func TestProperty22_EventUpdatePersistence(t *testing.T) {
	repo, _, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any existing event and valid attribute changes, updating should persist changes",
		prop.ForAll(
			func(originalTitle string, originalDescription string, originalDate time.Time, originalType string,
				updatedTitle string, updatedDescription string, updatedDate time.Time) bool {

				// Create initial event
				event := &models.Event{
					Title:       originalTitle,
					Description: originalDescription,
					EventDate:   originalDate,
					EventType:   originalType,
					MemberIDs:   []string{"test-member-1"},
					CreatedBy:   1,
				}

				err := repo.Create(ctx, event)
				if err != nil {
					t.Logf("Failed to create event: %v", err)
					return false
				}

				// Store original ID
				originalID := event.ID

				// Update event attributes
				event.Title = updatedTitle
				event.Description = updatedDescription
				event.EventDate = updatedDate

				// Persist update
				err = repo.Update(ctx, event)
				if err != nil {
					t.Logf("Failed to update event: %v", err)
					return false
				}

				// Retrieve and verify updates
				retrieved, err := repo.GetByID(ctx, event.ID)
				if err != nil {
					t.Logf("Failed to retrieve updated event: %v", err)
					return false
				}

				// Verify updated fields
				if retrieved.Title != updatedTitle {
					t.Logf("Title not updated: expected %s, got %s", updatedTitle, retrieved.Title)
					return false
				}
				if retrieved.Description != updatedDescription {
					t.Logf("Description not updated: expected %s, got %s", updatedDescription, retrieved.Description)
					return false
				}
				if !retrieved.EventDate.Equal(updatedDate) {
					t.Logf("EventDate not updated: expected %v, got %v", updatedDate, retrieved.EventDate)
					return false
				}

				// Verify unchanged fields
				if retrieved.ID != originalID {
					t.Logf("ID should not change: expected %d, got %d", originalID, retrieved.ID)
					return false
				}
				if retrieved.EventType != originalType {
					t.Logf("EventType should not change: expected %s, got %s", originalType, retrieved.EventType)
					return false
				}

				// Cleanup
				repo.Delete(ctx, event.ID)

				return true
			},
			genValidEventTitle(),
			genValidEventDescription(),
			genValidEventDate(),
			genValidEventType(),
			genValidEventTitle(),
			genValidEventDescription(),
			genValidEventDate(),
		),
	)

	properties.TestingRun(t)
}

// **Validates: Requirements 5.3**
// Property 23: Event Deletion
// For any existing event, deleting that event should remove it from PostgreSQL
// such that it is no longer retrievable.
func TestProperty23_EventDeletion(t *testing.T) {
	repo, _, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any existing event, deleting should remove it from PostgreSQL",
		prop.ForAll(
			func(title string, description string, eventDate time.Time, eventType string) bool {
				// Create event
				event := &models.Event{
					Title:       title,
					Description: description,
					EventDate:   eventDate,
					EventType:   eventType,
					MemberIDs:   []string{"test-member-1"},
					CreatedBy:   1,
				}

				err := repo.Create(ctx, event)
				if err != nil {
					t.Logf("Failed to create event: %v", err)
					return false
				}

				eventID := event.ID

				// Verify event exists
				_, err = repo.GetByID(ctx, eventID)
				if err != nil {
					t.Logf("Event should exist before deletion: %v", err)
					return false
				}

				// Delete event
				err = repo.Delete(ctx, eventID)
				if err != nil {
					t.Logf("Failed to delete event: %v", err)
					return false
				}

				// Verify event is no longer retrievable
				_, err = repo.GetByID(ctx, eventID)
				if err == nil {
					t.Logf("Deleted event should not be retrievable")
					return false
				}

				// Verify event is not in GetAll results
				allEvents, err := repo.GetAll(ctx)
				if err != nil {
					t.Logf("Failed to get all events: %v", err)
					return false
				}
				for _, e := range allEvents {
					if e.ID == eventID {
						t.Logf("Deleted event should not appear in GetAll results")
						return false
					}
				}

				return true
			},
			genValidEventTitle(),
			genValidEventDescription(),
			genValidEventDate(),
			genValidEventType(),
		),
	)

	properties.TestingRun(t)
}

// **Validates: Requirements 5.4**
// Property 24: Event Type Validity
// For any event in the system, that event's type must be one of:
// birthday, anniversary, ceremony, or custom.
func TestProperty24_EventTypeValidity(t *testing.T) {
	repo, _, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 100,
		MaxSize:            100,
	})

	properties.Property("For any event in the system, event type must be one of: birthday, anniversary, ceremony, custom",
		prop.ForAll(
			func(title string, description string, eventDate time.Time, eventType string) bool {
				// Create event with valid type
				event := &models.Event{
					Title:       title,
					Description: description,
					EventDate:   eventDate,
					EventType:   eventType,
					MemberIDs:   []string{"test-member-1"},
					CreatedBy:   1,
				}

				err := repo.Create(ctx, event)
				if err != nil {
					t.Logf("Failed to create event: %v", err)
					return false
				}

				// Retrieve and verify type is valid
				retrieved, err := repo.GetByID(ctx, event.ID)
				if err != nil {
					t.Logf("Failed to retrieve event: %v", err)
					return false
				}

				validTypes := map[string]bool{
					"birthday":    true,
					"anniversary": true,
					"ceremony":    true,
					"custom":      true,
				}

				if !validTypes[retrieved.EventType] {
					t.Logf("Invalid event type: %s", retrieved.EventType)
					return false
				}

				// Cleanup
				repo.Delete(ctx, event.ID)

				return true
			},
			genValidEventTitle(),
			genValidEventDescription(),
			genValidEventDate(),
			genValidEventType(),
		),
	)

	properties.TestingRun(t)
}

// **Validates: Requirements 5.7, 8.3, 8.4**
// Property 26: Event Filtering
// For any event filter criteria (type, member ID, date range), the returned events
// should match all specified filter criteria.
func TestProperty26_EventFiltering(t *testing.T) {
	repo, _, ctx, cleanup := setupEventTestRepo(t)
	defer cleanup()

	properties := gopter.NewProperties(&gopter.TestParameters{
		MinSuccessfulTests: 50,
		MaxSize:            50,
	})

	// Test filtering by type
	properties.Property("Events filtered by type should all have the specified type",
		prop.ForAll(
			func(filterType string, events []eventData) bool {
				// Create multiple events with different types
				createdIDs := []uint{}
				for _, ed := range events {
					event := &models.Event{
						Title:       ed.title,
						Description: ed.description,
						EventDate:   ed.eventDate,
						EventType:   ed.eventType,
						MemberIDs:   ed.memberIDs,
						CreatedBy:   1,
					}
					err := repo.Create(ctx, event)
					if err != nil {
						t.Logf("Failed to create event: %v", err)
						return false
					}
					createdIDs = append(createdIDs, event.ID)
				}

				// Filter by type
				filtered, err := repo.GetByType(ctx, filterType)
				if err != nil {
					t.Logf("Failed to filter events by type: %v", err)
					return false
				}

				// Verify all returned events have the specified type
				for _, e := range filtered {
					if e.EventType != filterType {
						t.Logf("Filtered event has wrong type: expected %s, got %s", filterType, e.EventType)
						return false
					}
				}

				// Cleanup
				for _, id := range createdIDs {
					repo.Delete(ctx, id)
				}

				return true
			},
			genValidEventType(),
			genEventDataList(),
		),
	)

	// Test filtering by member
	properties.Property("Events filtered by member should all contain the specified member ID",
		prop.ForAll(
			func(filterMemberID string, events []eventData) bool {
				// Create multiple events with different member IDs
				createdIDs := []uint{}
				for _, ed := range events {
					event := &models.Event{
						Title:       ed.title,
						Description: ed.description,
						EventDate:   ed.eventDate,
						EventType:   ed.eventType,
						MemberIDs:   ed.memberIDs,
						CreatedBy:   1,
					}
					err := repo.Create(ctx, event)
					if err != nil {
						t.Logf("Failed to create event: %v", err)
						return false
					}
					createdIDs = append(createdIDs, event.ID)
				}

				// Filter by member
				filtered, err := repo.GetByMember(ctx, filterMemberID)
				if err != nil {
					t.Logf("Failed to filter events by member: %v", err)
					return false
				}

				// Verify all returned events contain the specified member ID
				for _, e := range filtered {
					found := false
					for _, mid := range e.MemberIDs {
						if mid == filterMemberID {
							found = true
							break
						}
					}
					if !found {
						t.Logf("Filtered event does not contain member ID: %s", filterMemberID)
						return false
					}
				}

				// Cleanup
				for _, id := range createdIDs {
					repo.Delete(ctx, id)
				}

				return true
			},
			genValidMemberID(),
			genEventDataListWithMember(),
		),
	)

	// Test filtering by date range
	properties.Property("Events filtered by date range should all fall within the specified range",
		prop.ForAll(
			func(startDate time.Time, endDate time.Time, events []eventData) bool {
				// Ensure startDate is before endDate
				if startDate.After(endDate) {
					startDate, endDate = endDate, startDate
				}

				// Create multiple events with different dates
				createdIDs := []uint{}
				for _, ed := range events {
					event := &models.Event{
						Title:       ed.title,
						Description: ed.description,
						EventDate:   ed.eventDate,
						EventType:   ed.eventType,
						MemberIDs:   ed.memberIDs,
						CreatedBy:   1,
					}
					err := repo.Create(ctx, event)
					if err != nil {
						t.Logf("Failed to create event: %v", err)
						return false
					}
					createdIDs = append(createdIDs, event.ID)
				}

				// Filter by date range
				filtered, err := repo.GetByDateRange(ctx, startDate, endDate)
				if err != nil {
					t.Logf("Failed to filter events by date range: %v", err)
					return false
				}

				// Verify all returned events fall within the date range
				for _, e := range filtered {
					if e.EventDate.Before(startDate) || e.EventDate.After(endDate) {
						t.Logf("Filtered event outside date range: %v not in [%v, %v]", e.EventDate, startDate, endDate)
						return false
					}
				}

				// Cleanup
				for _, id := range createdIDs {
					repo.Delete(ctx, id)
				}

				return true
			},
			genValidEventDate(),
			genValidEventDate(),
			genEventDataList(),
		),
	)

	properties.TestingRun(t)
}

// Helper type for generating event data
type eventData struct {
	title       string
	description string
	eventDate   time.Time
	eventType   string
	memberIDs   []string
}

// Generator functions for property-based testing

// genValidEventTitle generates valid event titles (3-255 characters)
func genValidEventTitle() gopter.Gen {
	return gen.AlphaString().
		SuchThat(func(s string) bool {
			return len(s) >= 3 && len(s) <= 100
		}).
		Map(func(s string) string {
			return "Test Event " + s
		})
}

// genValidEventDescription generates valid event descriptions
func genValidEventDescription() gopter.Gen {
	return gen.AlphaString().
		SuchThat(func(s string) bool {
			return len(s) <= 500
		}).
		Map(func(s string) string {
			if s == "" {
				return "Test description"
			}
			return "Test description: " + s
		})
}

// genValidEventDate generates valid event dates (between now and 2 years from now)
func genValidEventDate() gopter.Gen {
	now := time.Now()
	maxDate := now.AddDate(2, 0, 0)

	return gen.Int64Range(now.Unix(), maxDate.Unix()).
		Map(func(timestamp int64) time.Time {
			return time.Unix(timestamp, 0).UTC()
		})
}

// genValidEventType generates valid event types
func genValidEventType() gopter.Gen {
	return gen.OneConstOf("birthday", "anniversary", "ceremony", "custom")
}

// genValidMemberIDs generates a list of valid member IDs (1-5 members)
func genValidMemberIDs() gopter.Gen {
	return gen.SliceOfN(
		gen.IntRange(1, 5),
		gen.AlphaString().
			SuchThat(func(s string) bool {
				return len(s) >= 5 && len(s) <= 36
			}).
			Map(func(s string) string {
				return "test-member-" + s
			}),
	).SuchThat(func(slice []string) bool {
		return len(slice) > 0
	})
}

// genValidMemberID generates a single valid member ID
func genValidMemberID() gopter.Gen {
	return gen.AlphaString().
		SuchThat(func(s string) bool {
			return len(s) >= 5 && len(s) <= 36
		}).
		Map(func(s string) string {
			return "test-member-" + s
		})
}

// genValidUserID generates valid user IDs (1-1000)
func genValidUserID() gopter.Gen {
	return gen.UIntRange(1, 1000).Map(func(n uint) uint {
		return n
	})
}

// genEventDataList generates a list of event data for filtering tests
func genEventDataList() gopter.Gen {
	return gen.SliceOfN(
		gen.IntRange(3, 10),
		gen.Struct(reflect.TypeOf(eventData{}), map[string]gopter.Gen{
			"title":       genValidEventTitle(),
			"description": genValidEventDescription(),
			"eventDate":   genValidEventDate(),
			"eventType":   genValidEventType(),
			"memberIDs":   genValidMemberIDs(),
		}),
	)
}

// genEventDataListWithMember generates a list of event data where some events contain a specific member
func genEventDataListWithMember() gopter.Gen {
	return gen.SliceOfN(
		gen.IntRange(3, 10),
		gen.Struct(reflect.TypeOf(eventData{}), map[string]gopter.Gen{
			"title":       genValidEventTitle(),
			"description": genValidEventDescription(),
			"eventDate":   genValidEventDate(),
			"eventType":   genValidEventType(),
			"memberIDs": gen.OneGenOf(
				// Some events with the filter member
				gen.Const([]string{"test-member-filter"}),
				// Some events with the filter member plus others
				gen.Const([]string{"test-member-filter", "test-member-other"}),
				// Some events without the filter member
				genValidMemberIDs(),
			),
		}),
	)
}
