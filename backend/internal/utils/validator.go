package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ValidationError represents a validation error with field and message
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface for ValidationError
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

// Error implements the error interface for ValidationErrors
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	if len(e) == 1 {
		return e[0].Error()
	}
	messages := make([]string, len(e))
	for i, err := range e {
		messages[i] = err.Error()
	}
	return strings.Join(messages, "; ")
}

// Add appends a validation error to the list
func (e *ValidationErrors) Add(field, message string) {
	*e = append(*e, ValidationError{Field: field, Message: message})
}

// HasErrors returns true if there are any validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Email validation regex (RFC 5322 simplified)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Phone validation regex (international format: +[country code][number])
// Supports formats like: +919876543210, +1234567890, +44-1234-567890
var phoneRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

// ValidateEmail validates an email address according to RFC 5322 format
// Returns a ValidationError if the email is invalid
func ValidateEmail(field, email string) *ValidationError {
	if email == "" {
		return &ValidationError{Field: field, Message: "email is required"}
	}
	
	email = strings.TrimSpace(email)
	
	if len(email) > 254 {
		return &ValidationError{Field: field, Message: "email must not exceed 254 characters"}
	}
	
	if !emailRegex.MatchString(email) {
		return &ValidationError{Field: field, Message: "email format is invalid (expected: user@example.com)"}
	}
	
	return nil
}

// ValidateDate validates a date string in ISO 8601 format (YYYY-MM-DD or RFC3339)
// Returns the parsed time.Time and a ValidationError if the date is invalid
func ValidateDate(field, dateStr string) (time.Time, *ValidationError) {
	if dateStr == "" {
		return time.Time{}, &ValidationError{Field: field, Message: "date is required"}
	}
	
	dateStr = strings.TrimSpace(dateStr)
	
	// Try parsing as RFC3339 (ISO 8601 with time)
	t, err := time.Parse(time.RFC3339, dateStr)
	if err == nil {
		return t, nil
	}
	
	// Try parsing as date only (YYYY-MM-DD)
	t, err = time.Parse("2006-01-02", dateStr)
	if err == nil {
		return t, nil
	}
	
	return time.Time{}, &ValidationError{
		Field:   field,
		Message: "date format is invalid (expected: YYYY-MM-DD or RFC3339)",
	}
}

// ValidateDateOptional validates an optional date string
// Returns the parsed time.Time and a ValidationError if the date is provided but invalid
func ValidateDateOptional(field, dateStr string) (time.Time, *ValidationError) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	return ValidateDate(field, dateStr)
}

// ValidatePhone validates a phone number in international format
// Expected format: +[country code][number] (e.g., +919876543210)
func ValidatePhone(field, phone string) *ValidationError {
	if phone == "" {
		return &ValidationError{Field: field, Message: "phone number is required"}
	}
	
	phone = strings.TrimSpace(phone)
	
	// Remove common separators for validation
	cleanPhone := strings.ReplaceAll(phone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")
	
	if !phoneRegex.MatchString(cleanPhone) {
		return &ValidationError{
			Field:   field,
			Message: "phone number format is invalid (expected: +[country code][number], e.g., +919876543210)",
		}
	}
	
	return nil
}

// ValidatePhoneOptional validates an optional phone number
// Returns a ValidationError if the phone is provided but invalid
func ValidatePhoneOptional(field, phone string) *ValidationError {
	if phone == "" {
		return nil
	}
	return ValidatePhone(field, phone)
}

// ValidateRequired validates that a string field is not empty
func ValidateRequired(field, value string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{Field: field, Message: fmt.Sprintf("%s is required", field)}
	}
	return nil
}

// ValidateRequiredInt validates that an integer field is not zero
func ValidateRequiredInt(field string, value int) *ValidationError {
	if value == 0 {
		return &ValidationError{Field: field, Message: fmt.Sprintf("%s is required", field)}
	}
	return nil
}

// ValidateRequiredUint validates that an unsigned integer field is not zero
func ValidateRequiredUint(field string, value uint) *ValidationError {
	if value == 0 {
		return &ValidationError{Field: field, Message: fmt.Sprintf("%s is required", field)}
	}
	return nil
}

// ValidateMinLength validates that a string meets minimum length requirement
func ValidateMinLength(field, value string, minLength int) *ValidationError {
	if len(strings.TrimSpace(value)) < minLength {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be at least %d characters", field, minLength),
		}
	}
	return nil
}

// ValidateMaxLength validates that a string does not exceed maximum length
func ValidateMaxLength(field, value string, maxLength int) *ValidationError {
	if len(value) > maxLength {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must not exceed %d characters", field, maxLength),
		}
	}
	return nil
}

// ValidateEnum validates that a value is one of the allowed options
func ValidateEnum(field, value string, allowedValues []string) *ValidationError {
	value = strings.TrimSpace(value)
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}
	return &ValidationError{
		Field:   field,
		Message: fmt.Sprintf("%s must be one of: %s", field, strings.Join(allowedValues, ", ")),
	}
}

// ValidateDateRange validates that a date falls within a specified range
func ValidateDateRange(field string, date, minDate, maxDate time.Time) *ValidationError {
	if !minDate.IsZero() && date.Before(minDate) {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be on or after %s", field, minDate.Format("2006-01-02")),
		}
	}
	if !maxDate.IsZero() && date.After(maxDate) {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be on or before %s", field, maxDate.Format("2006-01-02")),
		}
	}
	return nil
}

// ValidateFutureDate validates that a date is in the future
func ValidateFutureDate(field string, date time.Time) *ValidationError {
	if date.Before(time.Now()) {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be a future date", field),
		}
	}
	return nil
}

// ValidatePastDate validates that a date is in the past
func ValidatePastDate(field string, date time.Time) *ValidationError {
	if date.After(time.Now()) {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be a past date", field),
		}
	}
	return nil
}
