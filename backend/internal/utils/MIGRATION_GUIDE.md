# Migration Guide: Using Response Utilities

This guide shows how to refactor existing handlers to use the new response utilities.

## Benefits of Migration

1. **Less Boilerplate**: Reduce repetitive response code
2. **Consistency**: Ensure all responses follow the same format
3. **Maintainability**: Centralized response logic
4. **Type Safety**: Single source of truth for response structure

## Before and After Examples

### Example 1: Success Response

**Before:**
```go
return c.Status(fiber.StatusOK).JSON(fiber.Map{
    "success": true,
    "data":    member,
    "error":   "",
})
```

**After:**
```go
return utils.SuccessResponse(c, member)
```

### Example 2: Created Response

**Before:**
```go
return c.Status(fiber.StatusCreated).JSON(fiber.Map{
    "success": true,
    "data":    authResponse,
    "error":   "",
})
```

**After:**
```go
return utils.CreatedResponse(c, authResponse)
```

### Example 3: Bad Request Error

**Before:**
```go
return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    "success": false,
    "data":    nil,
    "error":   "Invalid request body",
})
```

**After:**
```go
return utils.BadRequestResponse(c, "Invalid request body")
```

### Example 4: Not Found Error

**Before:**
```go
return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
    "success": false,
    "data":    nil,
    "error":   "User not found",
})
```

**After:**
```go
return utils.NotFoundResponse(c, "User not found")
```

### Example 5: Unauthorized Error

**Before:**
```go
return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
    "success": false,
    "data":    nil,
    "error":   "User ID not found in context",
})
```

**After:**
```go
return utils.UnauthorizedResponse(c, "User ID not found in context")
```

## Complete Handler Refactoring Example

### Before: auth_handler.go

```go
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req service.LoginRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "data":    nil,
            "error":   "Invalid request body",
        })
    }

    authResponse, err := h.authService.Login(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "success": false,
            "data":    nil,
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "data":    authResponse,
        "error":   "",
    })
}
```

### After: auth_handler.go (with utilities)

```go
import "vamsasetu/backend/internal/utils"

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req service.LoginRequest

    if err := c.BodyParser(&req); err != nil {
        return utils.BadRequestResponse(c, "Invalid request body")
    }

    authResponse, err := h.authService.Login(c.Context(), &req)
    if err != nil {
        return utils.UnauthorizedResponse(c, err.Error())
    }

    return utils.SuccessResponse(c, authResponse)
}
```

**Lines of code reduced: 18 → 11 (39% reduction)**

## Migration Checklist

### Step 1: Import the utils package

Add to the top of your handler file:
```go
import "vamsasetu/backend/internal/utils"
```

### Step 2: Replace success responses

Find all instances of:
```go
c.Status(fiber.StatusOK).JSON(fiber.Map{
    "success": true,
    "data":    ...,
    "error":   "",
})
```

Replace with:
```go
utils.SuccessResponse(c, ...)
```

### Step 3: Replace created responses

Find all instances of:
```go
c.Status(fiber.StatusCreated).JSON(fiber.Map{
    "success": true,
    "data":    ...,
    "error":   "",
})
```

Replace with:
```go
utils.CreatedResponse(c, ...)
```

### Step 4: Replace error responses

| Status Code | Replace With |
|-------------|--------------|
| 400 Bad Request | `utils.BadRequestResponse(c, message)` |
| 401 Unauthorized | `utils.UnauthorizedResponse(c, message)` |
| 403 Forbidden | `utils.ForbiddenResponse(c, message)` |
| 404 Not Found | `utils.NotFoundResponse(c, message)` |
| 500 Internal Server Error | `utils.InternalServerErrorResponse(c, message)` |

### Step 5: Test your changes

Run tests to ensure behavior hasn't changed:
```bash
go test -v ./internal/handler/...
```

## Files to Migrate

Priority order for migration:

1. ✅ **New handlers** - Use utilities from the start
2. 🔄 **auth_handler.go** - Authentication endpoints
3. 🔄 **member_handler.go** - Member CRUD operations
4. 🔄 **relationship_handler.go** - Relationship operations
5. 🔄 **event_handler.go** - Event management
6. 🔄 **family_handler.go** - Family tree operations

## Common Patterns

### Pattern 1: Validation Error

```go
// Before
if req.Name == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "success": false,
        "data":    nil,
        "error":   "Name is required",
    })
}

// After
if req.Name == "" {
    return utils.BadRequestResponse(c, "Name is required")
}
```

### Pattern 2: Database Error

```go
// Before
member, err := h.memberService.GetByID(c.Context(), id)
if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "success": false,
        "data":    nil,
        "error":   err.Error(),
    })
}

// After
member, err := h.memberService.GetByID(c.Context(), id)
if err != nil {
    return utils.InternalServerErrorResponse(c, err.Error())
}
```

### Pattern 3: Resource Not Found

```go
// Before
if member == nil {
    return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
        "success": false,
        "data":    nil,
        "error":   "Member not found",
    })
}

// After
if member == nil {
    return utils.NotFoundResponse(c, "Member not found")
}
```

### Pattern 4: Complex Data Response

```go
// Before
return c.Status(fiber.StatusOK).JSON(fiber.Map{
    "success": true,
    "data": fiber.Map{
        "events": events,
        "total":  total,
        "page":   page,
    },
    "error": "",
})

// After
return utils.SuccessResponse(c, fiber.Map{
    "events": events,
    "total":  total,
    "page":   page,
})
```

## Testing After Migration

Ensure your tests still pass:

```bash
# Run all tests
go test -v ./...

# Run specific handler tests
go test -v ./internal/handler/

# Run with coverage
go test -v -coverprofile=coverage.out ./...
```

## Backward Compatibility

The response utilities produce **identical JSON output** to the manual approach, ensuring:

- ✅ No breaking changes for frontend clients
- ✅ Same HTTP status codes
- ✅ Same response structure
- ✅ Same field names and types

## Questions?

If you encounter any issues during migration, refer to:
- `response.go` - Implementation
- `response_test.go` - Test examples
- `response_example.go` - Usage examples
- `README.md` - Full documentation
