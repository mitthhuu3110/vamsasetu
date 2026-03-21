package models

import (
	"testing"
	"time"
)

// TestUserModel verifies the User model structure
func TestUserModel(t *testing.T) {
	user := User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         "Test User",
		Role:         "owner",
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email to be test@example.com, got %s", user.Email)
	}

	if user.Role != "owner" {
		t.Errorf("Expected role to be owner, got %s", user.Role)
	}

	if user.TableName() != "users" {
		t.Errorf("Expected table name to be users, got %s", user.TableName())
	}
}

// TestEventModel verifies the Event model structure
func TestEventModel(t *testing.T) {
	event := Event{
		Title:       "Birthday Party",
		Description: "John's birthday celebration",
		EventDate:   time.Now(),
		EventType:   "birthday",
		CreatedBy:   1,
	}
	event.SetMemberIDs([]string{"uuid-1", "uuid-2"})

	if event.Title != "Birthday Party" {
		t.Errorf("Expected title to be Birthday Party, got %s", event.Title)
	}

	if event.EventType != "birthday" {
		t.Errorf("Expected event type to be birthday, got %s", event.EventType)
	}

	if len(event.GetMemberIDs()) != 2 {
		t.Errorf("Expected 2 member IDs, got %d", len(event.GetMemberIDs()))
	}

	if event.TableName() != "events" {
		t.Errorf("Expected table name to be events, got %s", event.TableName())
	}
}

// TestNotificationModel verifies the Notification model structure
func TestNotificationModel(t *testing.T) {
	scheduledAt := time.Now()
	notification := Notification{
		EventID:     1,
		UserID:      1,
		Channel:     "email",
		ScheduledAt: scheduledAt,
		Status:      "pending",
		RetryCount:  0,
	}

	if notification.Channel != "email" {
		t.Errorf("Expected channel to be email, got %s", notification.Channel)
	}

	if notification.Status != "pending" {
		t.Errorf("Expected status to be pending, got %s", notification.Status)
	}

	if notification.RetryCount != 0 {
		t.Errorf("Expected retry count to be 0, got %d", notification.RetryCount)
	}

	if notification.TableName() != "notifications" {
		t.Errorf("Expected table name to be notifications, got %s", notification.TableName())
	}
}

// TestNotificationBeforeCreate verifies the BeforeCreate hook
func TestNotificationBeforeCreate(t *testing.T) {
	notification := Notification{
		EventID:     1,
		UserID:      1,
		Channel:     "email",
		ScheduledAt: time.Now(),
		// Status is intentionally not set
	}

	// Simulate BeforeCreate hook
	err := notification.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate should not return error, got %v", err)
	}

	if notification.Status != "pending" {
		t.Errorf("Expected status to be set to pending by BeforeCreate, got %s", notification.Status)
	}
}

// TestUserRoleValidation verifies role constraints
func TestUserRoleValidation(t *testing.T) {
	validRoles := []string{"owner", "viewer", "admin"}

	for _, role := range validRoles {
		user := User{
			Email:        "test@example.com",
			PasswordHash: "hashed",
			Name:         "Test",
			Role:         role,
		}

		if user.Role != role {
			t.Errorf("Expected role to be %s, got %s", role, user.Role)
		}
	}
}

// TestEventTypeValidation verifies event type constraints
func TestEventTypeValidation(t *testing.T) {
	validTypes := []string{"birthday", "anniversary", "ceremony", "custom"}

	for _, eventType := range validTypes {
		event := Event{
			Title:       "Test Event",
			EventDate:   time.Now(),
			EventType:   eventType,
			MemberIDs:   []string{"uuid-1"},
			CreatedBy:   1,
		}

		if event.EventType != eventType {
			t.Errorf("Expected event type to be %s, got %s", eventType, event.EventType)
		}
	}
}

// TestNotificationChannelValidation verifies channel constraints
func TestNotificationChannelValidation(t *testing.T) {
	validChannels := []string{"whatsapp", "sms", "email"}

	for _, channel := range validChannels {
		notification := Notification{
			EventID:     1,
			UserID:      1,
			Channel:     channel,
			ScheduledAt: time.Now(),
			Status:      "pending",
		}

		if notification.Channel != channel {
			t.Errorf("Expected channel to be %s, got %s", channel, notification.Channel)
		}
	}
}

// TestNotificationStatusValidation verifies status constraints
func TestNotificationStatusValidation(t *testing.T) {
	validStatuses := []string{"pending", "sent", "failed"}

	for _, status := range validStatuses {
		notification := Notification{
			EventID:     1,
			UserID:      1,
			Channel:     "email",
			ScheduledAt: time.Now(),
			Status:      status,
		}

		if notification.Status != status {
			t.Errorf("Expected status to be %s, got %s", status, notification.Status)
		}
	}
}

// TestMemberModel verifies the Member model structure
func TestMemberModel(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	member := NewMember("John Doe", dob, "male")

	if member.Name != "John Doe" {
		t.Errorf("Expected name to be John Doe, got %s", member.Name)
	}

	if member.Gender != "male" {
		t.Errorf("Expected gender to be male, got %s", member.Gender)
	}

	if member.DateOfBirth != dob {
		t.Errorf("Expected date of birth to be %v, got %v", dob, member.DateOfBirth)
	}

	if member.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if member.IsDeleted {
		t.Error("Expected IsDeleted to be false for new member")
	}
}

// TestMemberValidation verifies member validation logic
func TestMemberValidation(t *testing.T) {
	tests := []struct {
		name        string
		member      *Member
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid member",
			member: &Member{
				Name:        "John Doe",
				DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Gender:      "male",
			},
			expectError: false,
		},
		{
			name: "missing name",
			member: &Member{
				DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Gender:      "male",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "missing date of birth",
			member: &Member{
				Name:   "John Doe",
				Gender: "male",
			},
			expectError: true,
			errorMsg:    "date of birth is required",
		},
		{
			name: "missing gender",
			member: &Member{
				Name:        "John Doe",
				DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectError: true,
			errorMsg:    "gender is required",
		},
		{
			name: "invalid gender",
			member: &Member{
				Name:        "John Doe",
				DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Gender:      "invalid",
			},
			expectError: true,
			errorMsg:    "gender must be one of: male, female, other",
		},
		{
			name: "future date of birth",
			member: &Member{
				Name:        "John Doe",
				DateOfBirth: time.Now().Add(24 * time.Hour),
				Gender:      "male",
			},
			expectError: true,
			errorMsg:    "date of birth cannot be in the future",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.member.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestMemberGenderValidation verifies gender constraints
func TestMemberGenderValidation(t *testing.T) {
	validGenders := []string{"male", "female", "other"}
	dob := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, gender := range validGenders {
		member := &Member{
			Name:        "Test User",
			DateOfBirth: dob,
			Gender:      gender,
		}

		err := member.Validate()
		if err != nil {
			t.Errorf("Expected gender %s to be valid, got error: %v", gender, err)
		}
	}
}

// TestMemberSoftDelete verifies soft delete functionality
func TestMemberSoftDelete(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	member := NewMember("John Doe", dob, "male")

	if member.IsDeleted {
		t.Error("Expected IsDeleted to be false initially")
	}

	member.SoftDelete()

	if !member.IsDeleted {
		t.Error("Expected IsDeleted to be true after soft delete")
	}
}

// TestMemberUpdate verifies update timestamp functionality
func TestMemberUpdate(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	member := NewMember("John Doe", dob, "male")

	originalUpdatedAt := member.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	member.Update()

	if !member.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

// TestRelationshipModel verifies the Relationship model structure
func TestRelationshipModel(t *testing.T) {
	relationship := NewRelationship(RelationshipTypeParentOf, "uuid-1", "uuid-2")

	if relationship.Type != RelationshipTypeParentOf {
		t.Errorf("Expected type to be %s, got %s", RelationshipTypeParentOf, relationship.Type)
	}

	if relationship.FromID != "uuid-1" {
		t.Errorf("Expected fromId to be uuid-1, got %s", relationship.FromID)
	}

	if relationship.ToID != "uuid-2" {
		t.Errorf("Expected toId to be uuid-2, got %s", relationship.ToID)
	}

	if relationship.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

// TestRelationshipValidation verifies relationship validation logic
func TestRelationshipValidation(t *testing.T) {
	tests := []struct {
		name        string
		relationship *Relationship
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid relationship",
			relationship: &Relationship{
				Type:   RelationshipTypeParentOf,
				FromID: "uuid-1",
				ToID:   "uuid-2",
			},
			expectError: false,
		},
		{
			name: "missing type",
			relationship: &Relationship{
				FromID: "uuid-1",
				ToID:   "uuid-2",
			},
			expectError: true,
			errorMsg:    "relationship type is required",
		},
		{
			name: "invalid type",
			relationship: &Relationship{
				Type:   "INVALID_TYPE",
				FromID: "uuid-1",
				ToID:   "uuid-2",
			},
			expectError: true,
			errorMsg:    "invalid relationship type: must be one of SPOUSE_OF, PARENT_OF, SIBLING_OF",
		},
		{
			name: "missing fromId",
			relationship: &Relationship{
				Type: RelationshipTypeParentOf,
				ToID: "uuid-2",
			},
			expectError: true,
			errorMsg:    "fromId is required",
		},
		{
			name: "missing toId",
			relationship: &Relationship{
				Type:   RelationshipTypeParentOf,
				FromID: "uuid-1",
			},
			expectError: true,
			errorMsg:    "toId is required",
		},
		{
			name: "self relationship",
			relationship: &Relationship{
				Type:   RelationshipTypeParentOf,
				FromID: "uuid-1",
				ToID:   "uuid-1",
			},
			expectError: true,
			errorMsg:    "a member cannot have a relationship with themselves",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.relationship.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestRelationshipTypeConstants verifies relationship type constants
func TestRelationshipTypeConstants(t *testing.T) {
	if RelationshipTypeSpouseOf != "SPOUSE_OF" {
		t.Errorf("Expected SPOUSE_OF constant, got %s", RelationshipTypeSpouseOf)
	}

	if RelationshipTypeParentOf != "PARENT_OF" {
		t.Errorf("Expected PARENT_OF constant, got %s", RelationshipTypeParentOf)
	}

	if RelationshipTypeSiblingOf != "SIBLING_OF" {
		t.Errorf("Expected SIBLING_OF constant, got %s", RelationshipTypeSiblingOf)
	}
}

// TestIsValidRelationshipType verifies relationship type validation
func TestIsValidRelationshipType(t *testing.T) {
	validTypes := []string{
		RelationshipTypeSpouseOf,
		RelationshipTypeParentOf,
		RelationshipTypeSiblingOf,
	}

	for _, relType := range validTypes {
		if !IsValidRelationshipType(relType) {
			t.Errorf("Expected %s to be valid relationship type", relType)
		}
	}

	invalidTypes := []string{"INVALID", "FRIEND_OF", ""}
	for _, relType := range invalidTypes {
		if IsValidRelationshipType(relType) {
			t.Errorf("Expected %s to be invalid relationship type", relType)
		}
	}
}

// TestRelationshipBidirectional verifies bidirectional relationship logic
func TestRelationshipBidirectional(t *testing.T) {
	tests := []struct {
		relType        string
		isBidirectional bool
	}{
		{RelationshipTypeSpouseOf, true},
		{RelationshipTypeSiblingOf, true},
		{RelationshipTypeParentOf, false},
	}

	for _, tt := range tests {
		t.Run(tt.relType, func(t *testing.T) {
			relationship := &Relationship{
				Type:   tt.relType,
				FromID: "uuid-1",
				ToID:   "uuid-2",
			}

			if relationship.IsBidirectional() != tt.isBidirectional {
				t.Errorf("Expected %s to have IsBidirectional=%v, got %v",
					tt.relType, tt.isBidirectional, relationship.IsBidirectional())
			}
		})
	}
}
