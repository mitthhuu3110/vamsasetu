# Validator Utilities Guide

## Overview

The `validator.go` module provides comprehensive input validation utilities for the VamsaSetu backend. These utilities ensure data integrity and provide clear, user-friendly error messages for invalid input.

## Features

- **Email Validation**: RFC 5322 format validation
- **Date Validation**: ISO 8601 format support (YYYY-MM-DD and RFC3339)
- **Phone Validation**: International format validation (+[country code][number])
- **Required Field Validation**: String, int, and uint validation
- **Length Validation**: Minimum and maximum length constraints
- **Enum Validation**: Validate against allowed values
- **Date Range Validation**: Ensure dates fall within specified ranges
- **Future/Past Date Validation**: Validate temporal constraints

## Core Types

### ValidationError

Represents a single validation error with field name and message.

```go
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
```

### ValidationErrors

A collection of validation errors that can be accumulated and returned together.

```go
type ValidationErrors []ValidationError

// Methods
func (e *ValidationErrors) Add(field, message string)
func (e ValidationErrors) HasErrors() bool
func (e ValidationErrors) Error() string
```

## Validation Functions

### Email Validation

```go
func ValidateEmail(field, email string) *ValidationError
```

Validates email addresses according to RFC 5322 format.

**Rules:**
- Must not be empty
- Must not exceed 254 characters
- Must match pattern: `user@domain.tld`
- Supports subdomains, plus addressing, numbers, and hyphens

**Examples:**
```go
// Valid
ValidateEmail("email", "user@example.com")           // ✓
ValidateEmail("email", "user+tag@mail.example.com")  // ✓
ValidateEmail("email", "user123@example.co.uk")      // ✓

// Invalid
ValidateEmail("email", "")                           // ✗ email is required
ValidateEmail("email", "userexample.com")            // ✗ missing @
ValidateEmail("email", "user@")                      // ✗ missing domain
ValidateEmail("email", "@example.com")               // ✗ missing local part
```

### Date Validation

```go
func ValidateDate(field, dateStr string) (time.Time, *ValidationError)
func ValidateDateOptional(field, dateStr string) (time.Time, *ValidationError)
```

Validates date strings in ISO 8601 format.

**Supported Formats:**
- `YYYY-MM-DD` (e.g., "2024-01-15")
- RFC3339 (e.g., "2024-01-15T10:30:00Z")
- RFC3339 with timezone (e.g., "2024-01-15T10:30:00+05:30")

**Examples:**
```go
// Valid
ValidateDate("dateOfBirth", "1995-06-15")                    // ✓
ValidateDate("eventDate", "2024-01-15T10:30:00Z")           // ✓
ValidateDate("eventDate", "2024-01-15T10:30:00+05:30")      // ✓

// Invalid
ValidateDate("dateOfBirth", "")                              // ✗ date is required
ValidateDate("dateOfBirth", "15-06-1995")                    // ✗ invalid format
ValidateDate("dateOfBirth", "06/15/1995")                    // ✗ invalid format

// Optional (empty is allowed)
ValidateDateOptional("dateOfBirth", "")                      // ✓ returns zero time
ValidateDateOptional("dateOfBirth", "1995-06-15")           // ✓ returns parsed time
```

### Phone Validation

```go
func ValidatePhone(field, phone string) *ValidationError
func ValidatePhoneOptional(field, phone string) *ValidationError
```

Validates phone numbers in international format.

**Format:** `+[country code][number]`

**Rules:**
- Must start with `+`
- Country code: 1-3 digits (cannot start with 0)
- Total length: 8-15 digits (excluding `+`)
- Supports separators: hyphens, spaces, parentheses (removed during validation)

**Examples:**
```go
// Valid
ValidatePhone("phone", "+919876543210")              // ✓ India
ValidatePhone("phone", "+11234567890")               // ✓ USA
ValidatePhone("phone", "+441234567890")              // ✓ UK
ValidatePhone("phone", "+91-987-654-3210")           // ✓ with hyphens
ValidatePhone("phone", "+91 9876 543 210")           // ✓ with spaces

// Invalid
ValidatePhone("phone", "")                           // ✗ phone is required
ValidatePhone("phone", "919876543210")               // ✗ missing +
ValidatePhone("phone", "+9876543210")                // ✗ invalid country code
ValidatePhone("phone", "+91ABC7654321")              // ✗ contains letters

// Optional (empty is allowed)
ValidatePhoneOptional("phone", "")                   // ✓
ValidatePhoneOptional("phone", "+919876543210")      // ✓
```

### Required Field Validation

```go
func ValidateRequired(field, value string) *ValidationError
func ValidateRequiredInt(field string, value int) *ValidationError
func ValidateRequiredUint(field string, value uint) *ValidationError
```

Validates that required fields are not empty or zero.

**Examples:**
```go
// String
ValidateRequired("name", "John Doe")                 // ✓
ValidateRequired("name", "")                         // ✗ name is required
ValidateRequired("name", "   ")                      // ✗ whitespace only

// Integer
ValidateRequiredInt("age", 25)                       // ✓
ValidateRequiredInt("age", 0)                        // ✗ age is required

// Unsigned Integer
ValidateRequiredUint("userId", 123)                  // ✓
ValidateRequiredUint("userId", 0)                    // ✗ userId is required
```

### Length Validation

```go
func ValidateMinLength(field, value string, minLength int) *ValidationError
func ValidateMaxLength(field, value string, maxLength int) *ValidationError
```

Validates string length constraints.

**Examples:**
```go
// Minimum length
ValidateMinLength("password", "password123", 8)      // ✓
ValidateMinLength("password", "pass", 8)             // ✗ must be at least 8 characters

// Maximum length
ValidateMaxLength("name", "John Doe", 50)            // ✓
ValidateMaxLength("name", "Very long name...", 10)   // ✗ must not exceed 10 characters
```

### Enum Validation

```go
func ValidateEnum(field, value string, allowedValues []string) *ValidationError
```

Validates that a value is one of the allowed options.

**Examples:**
```go
allowedGenders := []string{"male", "female", "other"}

ValidateEnum("gender", "male", allowedGenders)       // ✓
ValidateEnum("gender", "female", allowedGenders)     // ✓
ValidateEnum("gender", "unknown", allowedGenders)    // ✗ must be one of: male, female, other
ValidateEnum("gender", "Male", allowedGenders)       // ✗ case sensitive
```

### Date Range Validation

```go
func ValidateDateRange(field string, date, minDate, maxDate time.Time) *ValidationError
func ValidateFutureDate(field string, date time.Time) *ValidationError
func ValidatePastDate(field string, date time.Time) *ValidationError
```

Validates temporal constraints on dates.

**Examples:**
```go
minDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
maxDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
date := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)

// Date range
ValidateDateRange("eventDate", date, minDate, maxDate)       // ✓
ValidateDateRange("eventDate", tooEarly, minDate, maxDate)   // ✗ must be on or after 2020-01-01
ValidateDateRange("eventDate", tooLate, minDate, maxDate)    // ✗ must be on or before 2025-12-31

// Future date
futureDate := time.Now().Add(24 * time.Hour)
ValidateFutureDate("eventDate", futureDate)                  // ✓
ValidateFutureDate("eventDate", time.Now().Add(-24 * time.Hour)) // ✗ must be a future date

// Past date
pastDate := time.Now().Add(-24 * time.Hour)
ValidatePastDate("dateOfBirth", pastDate)                    // ✓
ValidatePastDate("dateOfBirth", time.Now().Add(24 * time.Hour))  // ✗ must be a past date
```

## Usage Patterns

### Single Field Validation

```go
func CreateMember(name, email string) error {
    if err := ValidateRequired("name", name); err != nil {
        return BadRequestResponse(c, err.Message)
    }
    
    if err := ValidateEmail("email", email); err != nil {
        return BadRequestResponse(c, err.Message)
    }
    
    // Proceed with member creation
    return nil
}
```

### Multiple Field Validation

```go
func CreateMember(name, email, phone, dateOfBirth, gender string) error {
    var errors ValidationErrors
    
    // Validate all fields
    if err := ValidateRequired("name", name); err != nil {
        errors.Add(err.Field, err.Message)
    }
    
    if err := ValidateEmail("email", email); err != nil {
        errors.Add(err.Field, err.Message)
    }
    
    if err := ValidatePhoneOptional("phone", phone); err != nil {
        errors.Add(err.Field, err.Message)
    }
    
    dob, err := ValidateDate("dateOfBirth", dateOfBirth)
    if err != nil {
        errors.Add(err.Field, err.Message)
    } else {
        if err := ValidatePastDate("dateOfBirth", dob); err != nil {
            errors.Add(err.Field, err.Message)
        }
    }
    
    allowedGenders := []string{"male", "female", "other"}
    if err := ValidateEnum("gender", gender, allowedGenders); err != nil {
        errors.Add(err.Field, err.Message)
    }
    
    // Return all errors at once
    if errors.HasErrors() {
        return BadRequestResponse(c, errors.Error())
    }
    
    // Proceed with member creation
    return nil
}
```

### Handler Integration

```go
func (h *MemberHandler) Create(c *fiber.Ctx) error {
    var req CreateMemberRequest
    if err := c.BodyParser(&req); err != nil {
        return BadRequestResponse(c, "Invalid request body")
    }
    
    // Validate request
    var errors ValidationErrors
    
    if err := ValidateRequired("name", req.Name); err != nil {
        errors.Add(err.Field, err.Message)
    }
    
    if err := ValidateEmail("email", req.Email); err != nil {
        errors.Add(err.Field, err.Message)
    }
    
    if err := ValidatePhoneOptional("phone", req.Phone); err != nil {
        errors.Add(err.Field, err.Message)
    }
    
    dob, err := ValidateDate("dateOfBirth", req.DateOfBirth)
    if err != nil {
        errors.Add(err.Field, err.Message)
    } else {
        if err := ValidatePastDate("dateOfBirth", dob); err != nil {
            errors.Add(err.Field, err.Message)
        }
    }
    
    allowedGenders := []string{"male", "female", "other"}
    if err := ValidateEnum("gender", req.Gender, allowedGenders); err != nil {
        errors.Add(err.Field, err.Message)
    }
    
    if errors.HasErrors() {
        return BadRequestResponse(c, errors.Error())
    }
    
    // Create member
    member, err := h.service.Create(req)
    if err != nil {
        return InternalServerErrorResponse(c, "Failed to create member")
    }
    
    return CreatedResponse(c, member)
}
```

## Common Validation Scenarios

### User Registration

```go
var errors ValidationErrors

// Email
if err := ValidateEmail("email", req.Email); err != nil {
    errors.Add(err.Field, err.Message)
}

// Password
if err := ValidateRequired("password", req.Password); err != nil {
    errors.Add(err.Field, err.Message)
} else if err := ValidateMinLength("password", req.Password, 8); err != nil {
    errors.Add(err.Field, err.Message)
}

// Name
if err := ValidateRequired("name", req.Name); err != nil {
    errors.Add(err.Field, err.Message)
} else if err := ValidateMaxLength("name", req.Name, 255); err != nil {
    errors.Add(err.Field, err.Message)
}

// Role
allowedRoles := []string{"owner", "viewer", "admin"}
if err := ValidateEnum("role", req.Role, allowedRoles); err != nil {
    errors.Add(err.Field, err.Message)
}
```

### Event Creation

```go
var errors ValidationErrors

// Title
if err := ValidateRequired("title", req.Title); err != nil {
    errors.Add(err.Field, err.Message)
} else if err := ValidateMaxLength("title", req.Title, 255); err != nil {
    errors.Add(err.Field, err.Message)
}

// Event Date
eventDate, err := ValidateDate("eventDate", req.EventDate)
if err != nil {
    errors.Add(err.Field, err.Message)
} else if err := ValidateFutureDate("eventDate", eventDate); err != nil {
    errors.Add(err.Field, err.Message)
}

// Event Type
allowedTypes := []string{"birthday", "anniversary", "ceremony", "custom"}
if err := ValidateEnum("eventType", req.EventType, allowedTypes); err != nil {
    errors.Add(err.Field, err.Message)
}
```

### Member Creation

```go
var errors ValidationErrors

// Name
if err := ValidateRequired("name", req.Name); err != nil {
    errors.Add(err.Field, err.Message)
}

// Email
if err := ValidateEmail("email", req.Email); err != nil {
    errors.Add(err.Field, err.Message)
}

// Phone (optional)
if err := ValidatePhoneOptional("phone", req.Phone); err != nil {
    errors.Add(err.Field, err.Message)
}

// Date of Birth
dob, err := ValidateDate("dateOfBirth", req.DateOfBirth)
if err != nil {
    errors.Add(err.Field, err.Message)
} else if err := ValidatePastDate("dateOfBirth", dob); err != nil {
    errors.Add(err.Field, err.Message)
}

// Gender
allowedGenders := []string{"male", "female", "other"}
if err := ValidateEnum("gender", req.Gender, allowedGenders); err != nil {
    errors.Add(err.Field, err.Message)
}
```

## Error Messages

All validation functions return clear, user-friendly error messages:

- **Email**: "email format is invalid (expected: user@example.com)"
- **Date**: "date format is invalid (expected: YYYY-MM-DD or RFC3339)"
- **Phone**: "phone number format is invalid (expected: +[country code][number], e.g., +919876543210)"
- **Required**: "[field] is required"
- **Min Length**: "[field] must be at least [n] characters"
- **Max Length**: "[field] must not exceed [n] characters"
- **Enum**: "[field] must be one of: [value1, value2, ...]"
- **Date Range**: "[field] must be on or after [date]" / "[field] must be on or before [date]"
- **Future Date**: "[field] must be a future date"
- **Past Date**: "[field] must be a past date"

## Testing

Comprehensive tests are provided in `validator_test.go`:

```bash
# Run all validator tests
go test -v ./internal/utils/validator_test.go ./internal/utils/validator.go

# Run with coverage
go test -v -coverprofile=coverage.out ./internal/utils/validator_test.go ./internal/utils/validator.go
go tool cover -html=coverage.out
```

## Best Practices

1. **Accumulate Errors**: Use `ValidationErrors` to collect all validation errors and return them together, providing better user experience.

2. **Validate Early**: Perform validation at the handler level before passing data to services.

3. **Use Optional Validators**: For optional fields, use the `*Optional` variants to allow empty values.

4. **Chain Validations**: For complex validations, chain multiple validators (e.g., required + min length + max length).

5. **Consistent Error Handling**: Always use `BadRequestResponse` for validation errors to maintain consistent API responses.

6. **Document Constraints**: Document validation rules in API documentation and request models.

## Integration with Handlers

The validator utilities integrate seamlessly with the existing response utilities:

```go
import (
    "github.com/gofiber/fiber/v2"
    "vamsasetu/internal/utils"
)

func (h *Handler) Create(c *fiber.Ctx) error {
    var req CreateRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.BadRequestResponse(c, "Invalid request body")
    }
    
    var errors utils.ValidationErrors
    
    // Perform validations...
    
    if errors.HasErrors() {
        return utils.BadRequestResponse(c, errors.Error())
    }
    
    // Proceed with business logic...
    
    return utils.CreatedResponse(c, result)
}
```

## See Also

- `response.go` - API response utilities
- `jwt.go` - JWT token utilities
- `validator_example.go` - Usage examples
- `validator_test.go` - Comprehensive test suite
