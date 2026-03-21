// +build ignore

package utils

import (
	"fmt"
	"time"
)

// ExampleMemberValidation demonstrates how to validate member creation data
func ExampleMemberValidation() {
	var errors ValidationErrors

	// Sample member data
	name := "Arjun Kumar"
	email := "arjun@example.com"
	phone := "+919876543210"
	dateOfBirth := "1995-06-15"
	gender := "male"

	// Validate required fields
	if err := ValidateRequired("name", name); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Validate email
	if err := ValidateEmail("email", email); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Validate phone (optional)
	if err := ValidatePhoneOptional("phone", phone); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Validate date of birth
	dob, err := ValidateDate("dateOfBirth", dateOfBirth)
	if err != nil {
		errors.Add(err.Field, err.Message)
	} else {
		// Additional validation: date of birth should be in the past
		if err := ValidatePastDate("dateOfBirth", dob); err != nil {
			errors.Add(err.Field, err.Message)
		}
	}

	// Validate gender enum
	allowedGenders := []string{"male", "female", "other"}
	if err := ValidateEnum("gender", gender, allowedGenders); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Check if there are any validation errors
	if errors.HasErrors() {
		fmt.Println("Validation failed:")
		for _, err := range errors {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
		return
	}

	fmt.Println("Member validation passed!")
}

// ExampleEventValidation demonstrates how to validate event creation data
func ExampleEventValidation() {
	var errors ValidationErrors

	// Sample event data
	title := "Birthday Celebration"
	description := "Arjun's 30th birthday party"
	eventDate := "2024-06-15T18:00:00Z"
	eventType := "birthday"

	// Validate required fields
	if err := ValidateRequired("title", title); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Validate title length
	if err := ValidateMaxLength("title", title, 255); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Validate description length (optional)
	if description != "" {
		if err := ValidateMaxLength("description", description, 1000); err != nil {
			errors.Add(err.Field, err.Message)
		}
	}

	// Validate event date
	evtDate, err := ValidateDate("eventDate", eventDate)
	if err != nil {
		errors.Add(err.Field, err.Message)
	} else {
		// Additional validation: event date should be in the future
		if err := ValidateFutureDate("eventDate", evtDate); err != nil {
			errors.Add(err.Field, err.Message)
		}
	}

	// Validate event type enum
	allowedEventTypes := []string{"birthday", "anniversary", "ceremony", "custom"}
	if err := ValidateEnum("eventType", eventType, allowedEventTypes); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Check if there are any validation errors
	if errors.HasErrors() {
		fmt.Println("Validation failed:")
		for _, err := range errors {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
		return
	}

	fmt.Println("Event validation passed!")
}

// ExampleUserRegistrationValidation demonstrates how to validate user registration data
func ExampleUserRegistrationValidation() {
	var errors ValidationErrors

	// Sample registration data
	email := "user@example.com"
	password := "SecurePass123!"
	name := "John Doe"
	role := "owner"

	// Validate email
	if err := ValidateEmail("email", email); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Validate password
	if err := ValidateRequired("password", password); err != nil {
		errors.Add(err.Field, err.Message)
	} else {
		// Password must be at least 8 characters
		if err := ValidateMinLength("password", password, 8); err != nil {
			errors.Add(err.Field, err.Message)
		}
	}

	// Validate name
	if err := ValidateRequired("name", name); err != nil {
		errors.Add(err.Field, err.Message)
	} else {
		if err := ValidateMaxLength("name", name, 255); err != nil {
			errors.Add(err.Field, err.Message)
		}
	}

	// Validate role enum
	allowedRoles := []string{"owner", "viewer", "admin"}
	if err := ValidateEnum("role", role, allowedRoles); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Check if there are any validation errors
	if errors.HasErrors() {
		fmt.Println("Validation failed:")
		for _, err := range errors {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
		return
	}

	fmt.Println("User registration validation passed!")
}

// ExampleNotificationValidation demonstrates how to validate notification data
func ExampleNotificationValidation() {
	var errors ValidationErrors

	// Sample notification data
	channel := "whatsapp"
	scheduledAt := "2024-06-15T10:00:00Z"
	eventID := uint(123)
	userID := uint(456)

	// Validate required IDs
	if err := ValidateRequiredUint("eventId", eventID); err != nil {
		errors.Add(err.Field, err.Message)
	}

	if err := ValidateRequiredUint("userId", userID); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Validate channel enum
	allowedChannels := []string{"whatsapp", "sms", "email"}
	if err := ValidateEnum("channel", channel, allowedChannels); err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Validate scheduled date
	schedDate, err := ValidateDate("scheduledAt", scheduledAt)
	if err != nil {
		errors.Add(err.Field, err.Message)
	} else {
		// Scheduled date should be in the future
		if err := ValidateFutureDate("scheduledAt", schedDate); err != nil {
			errors.Add(err.Field, err.Message)
		}
	}

	// Check if there are any validation errors
	if errors.HasErrors() {
		fmt.Println("Validation failed:")
		for _, err := range errors {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
		return
	}

	fmt.Println("Notification validation passed!")
}

// ExampleDateRangeValidation demonstrates how to validate a date within a range
func ExampleDateRangeValidation() {
	var errors ValidationErrors

	// Sample date range for event filtering
	startDate := "2024-01-01"
	endDate := "2024-12-31"

	// Parse dates
	start, err := ValidateDate("startDate", startDate)
	if err != nil {
		errors.Add(err.Field, err.Message)
	}

	end, err := ValidateDate("endDate", endDate)
	if err != nil {
		errors.Add(err.Field, err.Message)
	}

	// Validate that end date is after start date
	if !start.IsZero() && !end.IsZero() {
		if end.Before(start) {
			errors.Add("endDate", "end date must be after start date")
		}
	}

	// Validate date range (e.g., within current year)
	minDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	if !start.IsZero() {
		if err := ValidateDateRange("startDate", start, minDate, maxDate); err != nil {
			errors.Add(err.Field, err.Message)
		}
	}

	if !end.IsZero() {
		if err := ValidateDateRange("endDate", end, minDate, maxDate); err != nil {
			errors.Add(err.Field, err.Message)
		}
	}

	// Check if there are any validation errors
	if errors.HasErrors() {
		fmt.Println("Validation failed:")
		for _, err := range errors {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
		return
	}

	fmt.Println("Date range validation passed!")
}
