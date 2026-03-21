package utils

import (
	"testing"
	"time"
)

// TestValidateEmail tests email validation
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		email     string
		wantError bool
		errMsg    string
	}{
		{
			name:      "valid email",
			field:     "email",
			email:     "user@example.com",
			wantError: false,
		},
		{
			name:      "valid email with subdomain",
			field:     "email",
			email:     "user@mail.example.com",
			wantError: false,
		},
		{
			name:      "valid email with plus",
			field:     "email",
			email:     "user+tag@example.com",
			wantError: false,
		},
		{
			name:      "valid email with numbers",
			field:     "email",
			email:     "user123@example.com",
			wantError: false,
		},
		{
			name:      "empty email",
			field:     "email",
			email:     "",
			wantError: true,
			errMsg:    "email is required",
		},
		{
			name:      "email without @",
			field:     "email",
			email:     "userexample.com",
			wantError: true,
			errMsg:    "email format is invalid",
		},
		{
			name:      "email without domain",
			field:     "email",
			email:     "user@",
			wantError: true,
			errMsg:    "email format is invalid",
		},
		{
			name:      "email without local part",
			field:     "email",
			email:     "@example.com",
			wantError: true,
			errMsg:    "email format is invalid",
		},
		{
			name:      "email with spaces",
			field:     "email",
			email:     "user @example.com",
			wantError: true,
			errMsg:    "email format is invalid",
		},
		{
			name:      "email too long",
			field:     "email",
			email:     "verylongemailaddressthatexceedsthemaximumlengthof254characters" + string(make([]byte, 200)) + "@example.com",
			wantError: true,
			errMsg:    "email must not exceed 254 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.field, tt.email)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateEmail() expected error but got nil")
					return
				}
				if err.Field != tt.field {
					t.Errorf("ValidateEmail() field = %v, want %v", err.Field, tt.field)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateEmail() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateDate tests date validation
func TestValidateDate(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		dateStr   string
		wantError bool
		errMsg    string
	}{
		{
			name:      "valid RFC3339 date",
			field:     "dateOfBirth",
			dateStr:   "1995-06-15T00:00:00Z",
			wantError: false,
		},
		{
			name:      "valid ISO 8601 date",
			field:     "dateOfBirth",
			dateStr:   "1995-06-15",
			wantError: false,
		},
		{
			name:      "valid RFC3339 with timezone",
			field:     "eventDate",
			dateStr:   "2024-01-15T10:30:00+05:30",
			wantError: false,
		},
		{
			name:      "empty date",
			field:     "dateOfBirth",
			dateStr:   "",
			wantError: true,
			errMsg:    "date is required",
		},
		{
			name:      "invalid date format",
			field:     "dateOfBirth",
			dateStr:   "15-06-1995",
			wantError: true,
			errMsg:    "date format is invalid",
		},
		{
			name:      "invalid date format with slashes",
			field:     "dateOfBirth",
			dateStr:   "06/15/1995",
			wantError: true,
			errMsg:    "date format is invalid",
		},
		{
			name:      "invalid date",
			field:     "dateOfBirth",
			dateStr:   "1995-13-45",
			wantError: true,
			errMsg:    "date format is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedDate, err := ValidateDate(tt.field, tt.dateStr)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateDate() expected error but got nil")
					return
				}
				if err.Field != tt.field {
					t.Errorf("ValidateDate() field = %v, want %v", err.Field, tt.field)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateDate() unexpected error: %v", err)
				}
				if parsedDate.IsZero() {
					t.Errorf("ValidateDate() returned zero time for valid date")
				}
			}
		})
	}
}

// TestValidateDateOptional tests optional date validation
func TestValidateDateOptional(t *testing.T) {
	// Empty date should not return error
	parsedDate, err := ValidateDateOptional("dateOfBirth", "")
	if err != nil {
		t.Errorf("ValidateDateOptional() with empty date should not return error, got: %v", err)
	}
	if !parsedDate.IsZero() {
		t.Errorf("ValidateDateOptional() with empty date should return zero time")
	}

	// Valid date should parse correctly
	parsedDate, err = ValidateDateOptional("dateOfBirth", "1995-06-15")
	if err != nil {
		t.Errorf("ValidateDateOptional() with valid date returned error: %v", err)
	}
	if parsedDate.IsZero() {
		t.Errorf("ValidateDateOptional() with valid date returned zero time")
	}

	// Invalid date should return error
	_, err = ValidateDateOptional("dateOfBirth", "invalid-date")
	if err == nil {
		t.Errorf("ValidateDateOptional() with invalid date should return error")
	}
}

// TestValidatePhone tests phone number validation
func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		phone     string
		wantError bool
		errMsg    string
	}{
		{
			name:      "valid Indian phone",
			field:     "phone",
			phone:     "+919876543210",
			wantError: false,
		},
		{
			name:      "valid US phone",
			field:     "phone",
			phone:     "+11234567890",
			wantError: false,
		},
		{
			name:      "valid UK phone",
			field:     "phone",
			phone:     "+441234567890",
			wantError: false,
		},
		{
			name:      "valid phone with hyphens",
			field:     "phone",
			phone:     "+91-987-654-3210",
			wantError: false,
		},
		{
			name:      "valid phone with spaces",
			field:     "phone",
			phone:     "+91 9876 543 210",
			wantError: false,
		},
		{
			name:      "empty phone",
			field:     "phone",
			phone:     "",
			wantError: true,
			errMsg:    "phone number is required",
		},
		{
			name:      "phone without plus",
			field:     "phone",
			phone:     "919876543210",
			wantError: true,
			errMsg:    "phone number format is invalid",
		},
		{
			name:      "phone without country code",
			field:     "phone",
			phone:     "+9876543210",
			wantError: true,
			errMsg:    "phone number format is invalid",
		},
		{
			name:      "phone too short",
			field:     "phone",
			phone:     "+91987",
			wantError: true,
			errMsg:    "phone number format is invalid",
		},
		{
			name:      "phone with letters",
			field:     "phone",
			phone:     "+91ABC7654321",
			wantError: true,
			errMsg:    "phone number format is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhone(tt.field, tt.phone)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidatePhone() expected error but got nil")
					return
				}
				if err.Field != tt.field {
					t.Errorf("ValidatePhone() field = %v, want %v", err.Field, tt.field)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePhone() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidatePhoneOptional tests optional phone validation
func TestValidatePhoneOptional(t *testing.T) {
	// Empty phone should not return error
	err := ValidatePhoneOptional("phone", "")
	if err != nil {
		t.Errorf("ValidatePhoneOptional() with empty phone should not return error, got: %v", err)
	}

	// Valid phone should pass
	err = ValidatePhoneOptional("phone", "+919876543210")
	if err != nil {
		t.Errorf("ValidatePhoneOptional() with valid phone returned error: %v", err)
	}

	// Invalid phone should return error
	err = ValidatePhoneOptional("phone", "invalid-phone")
	if err == nil {
		t.Errorf("ValidatePhoneOptional() with invalid phone should return error")
	}
}

// TestValidateRequired tests required field validation
func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		wantError bool
	}{
		{
			name:      "non-empty value",
			field:     "name",
			value:     "John Doe",
			wantError: false,
		},
		{
			name:      "empty value",
			field:     "name",
			value:     "",
			wantError: true,
		},
		{
			name:      "whitespace only",
			field:     "name",
			value:     "   ",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequired(tt.field, tt.value)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateRequired() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateRequired() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateRequiredUint tests required uint validation
func TestValidateRequiredUint(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     uint
		wantError bool
	}{
		{
			name:      "non-zero value",
			field:     "userId",
			value:     123,
			wantError: false,
		},
		{
			name:      "zero value",
			field:     "userId",
			value:     0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequiredUint(tt.field, tt.value)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateRequiredUint() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateRequiredUint() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateMinLength tests minimum length validation
func TestValidateMinLength(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		minLength int
		wantError bool
	}{
		{
			name:      "meets minimum length",
			field:     "password",
			value:     "password123",
			minLength: 8,
			wantError: false,
		},
		{
			name:      "below minimum length",
			field:     "password",
			value:     "pass",
			minLength: 8,
			wantError: true,
		},
		{
			name:      "exact minimum length",
			field:     "password",
			value:     "password",
			minLength: 8,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMinLength(tt.field, tt.value, tt.minLength)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateMinLength() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateMinLength() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateMaxLength tests maximum length validation
func TestValidateMaxLength(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		maxLength int
		wantError bool
	}{
		{
			name:      "within maximum length",
			field:     "name",
			value:     "John Doe",
			maxLength: 50,
			wantError: false,
		},
		{
			name:      "exceeds maximum length",
			field:     "name",
			value:     "This is a very long name that exceeds the maximum allowed length",
			maxLength: 50,
			wantError: true,
		},
		{
			name:      "exact maximum length",
			field:     "name",
			value:     "Exactly fifty characters in this name right here!",
			maxLength: 50,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMaxLength(tt.field, tt.value, tt.maxLength)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateMaxLength() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateMaxLength() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateEnum tests enum validation
func TestValidateEnum(t *testing.T) {
	allowedGenders := []string{"male", "female", "other"}

	tests := []struct {
		name      string
		field     string
		value     string
		wantError bool
	}{
		{
			name:      "valid enum value",
			field:     "gender",
			value:     "male",
			wantError: false,
		},
		{
			name:      "another valid enum value",
			field:     "gender",
			value:     "female",
			wantError: false,
		},
		{
			name:      "invalid enum value",
			field:     "gender",
			value:     "unknown",
			wantError: true,
		},
		{
			name:      "case sensitive mismatch",
			field:     "gender",
			value:     "Male",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnum(tt.field, tt.value, allowedGenders)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateEnum() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateEnum() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateDateRange tests date range validation
func TestValidateDateRange(t *testing.T) {
	minDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		field     string
		date      time.Time
		wantError bool
	}{
		{
			name:      "date within range",
			field:     "eventDate",
			date:      time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			wantError: false,
		},
		{
			name:      "date before minimum",
			field:     "eventDate",
			date:      time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC),
			wantError: true,
		},
		{
			name:      "date after maximum",
			field:     "eventDate",
			date:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			wantError: true,
		},
		{
			name:      "date at minimum boundary",
			field:     "eventDate",
			date:      minDate,
			wantError: false,
		},
		{
			name:      "date at maximum boundary",
			field:     "eventDate",
			date:      maxDate,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDateRange(tt.field, tt.date, minDate, maxDate)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateDateRange() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateDateRange() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateFutureDate tests future date validation
func TestValidateFutureDate(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)
	pastDate := time.Now().Add(-24 * time.Hour)

	err := ValidateFutureDate("eventDate", futureDate)
	if err != nil {
		t.Errorf("ValidateFutureDate() with future date returned error: %v", err)
	}

	err = ValidateFutureDate("eventDate", pastDate)
	if err == nil {
		t.Errorf("ValidateFutureDate() with past date should return error")
	}
}

// TestValidatePastDate tests past date validation
func TestValidatePastDate(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)
	pastDate := time.Now().Add(-24 * time.Hour)

	err := ValidatePastDate("dateOfBirth", pastDate)
	if err != nil {
		t.Errorf("ValidatePastDate() with past date returned error: %v", err)
	}

	err = ValidatePastDate("dateOfBirth", futureDate)
	if err == nil {
		t.Errorf("ValidatePastDate() with future date should return error")
	}
}

// TestValidationErrors tests ValidationErrors collection
func TestValidationErrors(t *testing.T) {
	var errors ValidationErrors

	// Test empty errors
	if errors.HasErrors() {
		t.Errorf("ValidationErrors.HasErrors() should return false for empty errors")
	}

	// Add errors
	errors.Add("email", "email is required")
	errors.Add("phone", "phone format is invalid")

	if !errors.HasErrors() {
		t.Errorf("ValidationErrors.HasErrors() should return true after adding errors")
	}

	if len(errors) != 2 {
		t.Errorf("ValidationErrors length = %d, want 2", len(errors))
	}

	// Test error message
	errMsg := errors.Error()
	if errMsg == "" {
		t.Errorf("ValidationErrors.Error() should return non-empty string")
	}
}

// TestValidationError tests ValidationError
func TestValidationError(t *testing.T) {
	err := ValidationError{Field: "email", Message: "email is required"}
	
	expectedMsg := "email: email is required"
	if err.Error() != expectedMsg {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), expectedMsg)
	}
}
