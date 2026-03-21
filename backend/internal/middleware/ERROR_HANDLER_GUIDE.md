# Error Handling Middleware Guide

## Overview

The error handling middleware provides a centralized, consistent way to handle errors across the VamsaSetu backend. It catches all errors returned from handlers and formats them using the standard `APIResponse` format.

## Features

- **Structured Errors**: `AppError` type with HTTP status code, user-friendly message, and underlying error
- **Automatic Error Mapping**: Converts common errors (GORM, validation, Fiber) to appropriate HTTP responses
- **Consistent Format**: All errors return the standard `{ success, data, error }` format
- **Error Logging**: Logs underlying errors for debugging while keeping user-facing messages clean
- **Type Safety**: Uses Go's `errors.Is` and `errors.As` for proper error handling

## Setup

### 1. Configure Error Handler in main.go

```go
package main

import (
    "vamsasetu/backend/internal/middleware"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New(fiber.Config{
        // Set the custom error handler
        ErrorHandler: middleware.ErrorHandler(),
    })

    // Register routes...
    
    app.Listen(":8080")
}
```

### 2. Use in Handlers

Once configured, simply return errors from your handlers:

```go
func (h *MemberHandler) GetMember(c *fiber.Ctx) error {
    id := c.Params("id")
    
    member, err := h.memberService.GetByID(c.Context(), id)
    if err != nil {
        // Return error - ErrorHandler will catch and format it
        return middleware.MapServiceError(err)
    }
    
    return utils.SuccessResponse(c, member)
}
```

## AppError Type

### Structure

```go
type AppError struct {
    Code    int    // HTTP status code
    Message string // User-friendly error message
    Err     error  // Original error (for logging)
}
```

### Constructor Functions

```go
// 400 Bad Request
middleware.BadRequest("Invalid input", err)

// 401 Unauthorized
middleware.Unauthorized("Authentication required", err)

// 403 Forbidden
middleware.Forbidden("Access denied", err)

// 404 Not Found
middleware.NotFound("Resource not found", err)

// 409 Conflict
middleware.Conflict("Resource already exists", err)

// 500 Internal Server Error
middleware.InternalServerError("Server error", err)
```

## Error Mapping

### Automatic Error Type Detection

The `ErrorHandler` automatically detects and maps common error types:

| Error Type | HTTP Status | Example |
|------------|-------------|---------|
| `AppError` | Custom code | `BadRequest("Invalid input", nil)` |
| `gorm.ErrRecordNotFound` | 404 | Database record not found |
| `utils.ValidationErrors` | 400 | Multiple validation errors |
| `utils.ValidationError` | 400 | Single validation error |
| `fiber.Error` | Custom code | Fiber's built-in errors |
| Unknown errors | 500 | Any other error |

### MapServiceError Helper

Use `MapServiceError` to automatically convert service layer errors:

```go
func (h *Handler) Create(c *fiber.Ctx) error {
    result, err := h.service.Create(data)
    if err != nil {
        // Automatically maps to appropriate HTTP error
        return middleware.MapServiceError(err)
    }
    return utils.CreatedResponse(c, result)
}
```

## Usage Patterns

### Pattern 1: Input Validation

```go
func (h *Handler) CreateMember(c *fiber.Ctx) error {
    var req CreateMemberRequest
    if err := c.BodyParser(&req); err != nil {
        return middleware.BadRequest("Invalid request body", err)
    }
    
    // Validate fields
    var errs utils.ValidationErrors
    if err := utils.ValidateRequired("name", req.Name); err != nil {
        errs = append(errs, *err)
    }
    if err := utils.ValidateEmail("email", req.Email); err != nil {
        errs = append(errs, *err)
    }
    
    if errs.HasErrors() {
        return errs // Returns 400 with validation details
    }
    
    // Process request...
}
```

### Pattern 2: Resource Not Found

```go
func (h *Handler) GetMember(c *fiber.Ctx) error {
    id := c.Params("id")
    
    member, err := h.service.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return middleware.NotFound("Member not found", err)
        }
        return middleware.InternalServerError("Failed to fetch member", err)
    }
    
    return utils.SuccessResponse(c, member)
}
```

### Pattern 3: Authorization Checks

```go
func (h *Handler) UpdateMember(c *fiber.Ctx) error {
    // Check authentication
    userID := c.Locals("userId")
    if userID == nil {
        return middleware.Unauthorized("Authentication required", nil)
    }
    
    // Check permissions
    if !h.service.HasPermission(userID, resourceID) {
        return middleware.Forbidden("Insufficient permissions", nil)
    }
    
    // Process update...
}
```

### Pattern 4: Conflict Detection

```go
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return middleware.BadRequest("Invalid request body", err)
    }
    
    // Check if user already exists
    exists, err := h.service.UserExists(req.Email)
    if err != nil {
        return middleware.InternalServerError("Failed to check user", err)
    }
    if exists {
        return middleware.Conflict("User with this email already exists", nil)
    }
    
    // Create user...
}
```

### Pattern 5: Service Layer Errors

```go
// In service layer
func (s *MemberService) Create(member *Member) error {
    // Validate business rules
    if err := s.validateMember(member); err != nil {
        // Return validation errors directly
        return err
    }
    
    // Database operation
    if err := s.repo.Create(member); err != nil {
        // Wrap with context
        return fmt.Errorf("failed to create member: %w", err)
    }
    
    return nil
}

// In handler
func (h *Handler) CreateMember(c *fiber.Ctx) error {
    member := parseRequest(c)
    
    if err := h.service.Create(member); err != nil {
        // MapServiceError handles all error types
        return middleware.MapServiceError(err)
    }
    
    return utils.CreatedResponse(c, member)
}
```

## Response Format

All errors return the standard APIResponse format:

```json
{
  "success": false,
  "data": null,
  "error": "User-friendly error message"
}
```

### Examples

**400 Bad Request:**
```json
{
  "success": false,
  "data": null,
  "error": "email: email is required; name: name must be at least 3 characters"
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "data": null,
  "error": "Authentication required"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "data": null,
  "error": "Member not found"
}
```

**500 Internal Server Error:**
```json
{
  "success": false,
  "data": null,
  "error": "An unexpected error occurred"
}
```

## Error Logging

The middleware automatically logs errors with appropriate context:

- **AppError with underlying error**: Logs both user message and underlying error
- **GORM errors**: Logs "Record not found" with error details
- **Validation errors**: Logs validation details
- **Unexpected errors**: Logs full error details

Example log output:
```
[ERROR] Failed to fetch member: connection timeout
[VALIDATION] email: email is required; name: name must be at least 3 characters
[ERROR] Unexpected error: some unexpected error
```

## Best Practices

### 1. Use Specific Error Types

```go
// Good: Specific error type
return middleware.NotFound("Member not found", err)

// Avoid: Generic error
return middleware.InternalServerError("Error", err)
```

### 2. Provide User-Friendly Messages

```go
// Good: Clear, actionable message
return middleware.BadRequest("Email format is invalid. Expected: user@example.com", nil)

// Avoid: Technical jargon
return middleware.BadRequest("regex match failed", nil)
```

### 3. Include Underlying Errors for Logging

```go
// Good: Includes underlying error for debugging
return middleware.InternalServerError("Failed to fetch data", dbErr)

// Avoid: Loses error context
return middleware.InternalServerError("Failed to fetch data", nil)
```

### 4. Use MapServiceError for Service Errors

```go
// Good: Automatic error mapping
if err := h.service.Create(data); err != nil {
    return middleware.MapServiceError(err)
}

// Avoid: Manual error checking
if err := h.service.Create(data); err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return middleware.NotFound("Not found", err)
    }
    // ... more checks
}
```

### 5. Validate Early, Fail Fast

```go
// Good: Validate input before processing
if id == "" {
    return middleware.BadRequest("ID is required", nil)
}

// Avoid: Processing invalid input
member, err := h.service.GetByID(id) // id might be empty
```

## Integration with Existing Code

### Updating Existing Handlers

**Before:**
```go
func (h *Handler) GetMember(c *fiber.Ctx) error {
    member, err := h.service.GetByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "data":    nil,
            "error":   "Member not found",
        })
    }
    return c.JSON(fiber.Map{
        "success": true,
        "data":    member,
        "error":   "",
    })
}
```

**After:**
```go
func (h *Handler) GetMember(c *fiber.Ctx) error {
    member, err := h.service.GetByID(id)
    if err != nil {
        return middleware.MapServiceError(err)
    }
    return utils.SuccessResponse(c, member)
}
```

### Benefits

- **Less Code**: Reduced boilerplate in handlers
- **Consistency**: All errors follow the same format
- **Maintainability**: Centralized error handling logic
- **Debugging**: Better error logging with context
- **Type Safety**: Proper error type checking with `errors.Is` and `errors.As`

## Testing

The error handler is fully tested with comprehensive test coverage:

```bash
cd backend
go test ./internal/middleware/error_test.go ./internal/middleware/error.go -v
```

Test coverage includes:
- AppError construction and methods
- All error constructor functions
- ErrorHandler with different error types
- MapServiceError with various inputs
- Integration tests with full request/response cycle
- Response format consistency

## Migration Checklist

- [ ] Add `ErrorHandler: middleware.ErrorHandler()` to Fiber config in main.go
- [ ] Update handlers to return errors instead of manual JSON responses
- [ ] Use `middleware.MapServiceError()` for service layer errors
- [ ] Replace manual error responses with error constructors
- [ ] Use `utils.SuccessResponse()` and `utils.CreatedResponse()` for success cases
- [ ] Test error responses to ensure consistent format
- [ ] Update API documentation with error response examples

