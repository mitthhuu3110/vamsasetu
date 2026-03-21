# Error Handler Integration Example

## Quick Start

### Step 1: Update main.go

```go
package main

import (
    "log"
    "vamsasetu/backend/internal/config"
    "vamsasetu/backend/internal/handler"
    "vamsasetu/backend/internal/middleware"
    "vamsasetu/backend/internal/repository"
    "vamsasetu/backend/internal/service"
    "vamsasetu/backend/pkg/neo4j"
    "vamsasetu/backend/pkg/postgres"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize databases
    db, err := postgres.NewClient(cfg.PostgresURL)
    if err != nil {
        log.Fatalf("Failed to connect to PostgreSQL: %v", err)
    }

    neo4jDriver, err := neo4j.NewClient(cfg.Neo4jURI, cfg.Neo4jUsername, cfg.Neo4jPassword)
    if err != nil {
        log.Fatalf("Failed to connect to Neo4j: %v", err)
    }
    defer neo4jDriver.Close()

    // Initialize repositories
    memberRepo := repository.NewMemberRepository(neo4jDriver)

    // Initialize services
    memberService := service.NewMemberService(memberRepo, nil)

    // Initialize handlers
    memberHandler := handler.NewMemberHandler(memberService)

    // Create Fiber app with error handler
    app := fiber.New(fiber.Config{
        // ✅ Add the error handler here
        ErrorHandler: middleware.ErrorHandler(),
    })

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New())

    // Health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "healthy",
        })
    })

    // Register routes
    memberHandler.RegisterRoutes(app)

    // Start server
    log.Printf("Server starting on port %s", cfg.Port)
    if err := app.Listen(":" + cfg.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

### Step 2: Update Handler (Example)

**Before:**
```go
func (h *MemberHandler) GetMember(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "data":    nil,
            "error":   "Member ID is required",
        })
    }

    member, err := h.memberService.GetByID(c.Context(), id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "data":    nil,
            "error":   "Member not found",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "data":    member,
        "error":   "",
    })
}
```

**After:**
```go
func (h *MemberHandler) GetMember(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return middleware.BadRequest("Member ID is required", nil)
    }

    member, err := h.memberService.GetByID(c.Context(), id)
    if err != nil {
        return middleware.MapServiceError(err)
    }

    return utils.SuccessResponse(c, member)
}
```

### Step 3: Update Service Layer (Optional)

Services can return standard errors, and the middleware will handle them:

```go
func (s *MemberService) GetByID(ctx context.Context, id string) (*models.Member, error) {
    member, err := s.repo.GetByID(ctx, id)
    if err != nil {
        // Return GORM error - middleware will map to 404
        return nil, err
    }
    
    if member.IsDeleted {
        // Return GORM not found - middleware will map to 404
        return nil, gorm.ErrRecordNotFound
    }
    
    return member, nil
}
```

## Common Patterns

### Pattern 1: Input Validation

```go
func (h *Handler) CreateMember(c *fiber.Ctx) error {
    var req CreateMemberRequest
    if err := c.BodyParser(&req); err != nil {
        return middleware.BadRequest("Invalid request body", err)
    }

    // Validate using utils
    var errs utils.ValidationErrors
    if err := utils.ValidateRequired("name", req.Name); err != nil {
        errs = append(errs, *err)
    }
    if err := utils.ValidateEmail("email", req.Email); err != nil {
        errs = append(errs, *err)
    }

    if errs.HasErrors() {
        return errs // Middleware handles this automatically
    }

    // Process...
    member, err := h.service.Create(c.Context(), req)
    if err != nil {
        return middleware.MapServiceError(err)
    }

    return utils.CreatedResponse(c, member)
}
```

### Pattern 2: Authentication Check

```go
func (h *Handler) UpdateMember(c *fiber.Ctx) error {
    // Check if user is authenticated (set by AuthMiddleware)
    userID := c.Locals("userId")
    if userID == nil {
        return middleware.Unauthorized("Authentication required", nil)
    }

    // Process update...
}
```

### Pattern 3: Authorization Check

```go
func (h *Handler) DeleteMember(c *fiber.Ctx) error {
    userRole := c.Locals("userRole").(string)
    
    if userRole != "owner" && userRole != "admin" {
        return middleware.Forbidden("Insufficient permissions to delete members", nil)
    }

    // Process deletion...
}
```

### Pattern 4: Resource Existence Check

```go
func (h *Handler) UpdateMember(c *fiber.Ctx) error {
    id := c.Params("id")
    
    member, err := h.service.GetByID(c.Context(), id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return middleware.NotFound("Member not found", err)
        }
        return middleware.InternalServerError("Failed to fetch member", err)
    }

    // Update member...
}
```

### Pattern 5: Conflict Detection

```go
func (h *Handler) CreateRelationship(c *fiber.Ctx) error {
    var req CreateRelationshipRequest
    if err := c.BodyParser(&req); err != nil {
        return middleware.BadRequest("Invalid request body", err)
    }

    // Check if relationship already exists
    exists, err := h.service.RelationshipExists(req.FromID, req.ToID)
    if err != nil {
        return middleware.InternalServerError("Failed to check relationship", err)
    }
    if exists {
        return middleware.Conflict("Relationship already exists", nil)
    }

    // Create relationship...
}
```

## Testing Error Responses

### Using curl

```bash
# Test 400 Bad Request
curl -X POST http://localhost:8080/api/members \
  -H "Content-Type: application/json" \
  -d '{"name": ""}'

# Response:
# {
#   "success": false,
#   "data": null,
#   "error": "name: name is required"
# }

# Test 401 Unauthorized
curl http://localhost:8080/api/members

# Response:
# {
#   "success": false,
#   "data": null,
#   "error": "Missing authorization header"
# }

# Test 404 Not Found
curl http://localhost:8080/api/members/invalid-id \
  -H "Authorization: Bearer <token>"

# Response:
# {
#   "success": false,
#   "data": null,
#   "error": "Member not found"
# }
```

### Using httptest (in tests)

```go
func TestHandler_GetMember_NotFound(t *testing.T) {
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.ErrorHandler(),
    })
    
    handler := NewMemberHandler(mockService)
    handler.RegisterRoutes(app)
    
    req := httptest.NewRequest("GET", "/api/members/invalid-id", nil)
    resp, err := app.Test(req)
    
    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
    
    var response utils.APIResponse
    json.NewDecoder(resp.Body).Decode(&response)
    
    assert.False(t, response.Success)
    assert.Nil(t, response.Data)
    assert.Contains(t, response.Error, "not found")
}
```

## Error Response Examples

### 400 Bad Request - Validation Error
```json
{
  "success": false,
  "data": null,
  "error": "email: email is required; name: name must be at least 3 characters"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "data": null,
  "error": "Authentication required"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "data": null,
  "error": "Insufficient permissions to delete members"
}
```

### 404 Not Found
```json
{
  "success": false,
  "data": null,
  "error": "Member not found"
}
```

### 409 Conflict
```json
{
  "success": false,
  "data": null,
  "error": "User with this email already exists"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "data": null,
  "error": "An unexpected error occurred"
}
```

## Migration Checklist

- [ ] Add `ErrorHandler: middleware.ErrorHandler()` to Fiber config in main.go
- [ ] Update all handlers to return errors instead of manual JSON responses
- [ ] Replace `c.Status().JSON()` error responses with error constructors
- [ ] Use `middleware.MapServiceError()` for service layer errors
- [ ] Use `utils.SuccessResponse()` and `utils.CreatedResponse()` for success cases
- [ ] Update handler tests to expect new error format
- [ ] Test all error scenarios (400, 401, 403, 404, 409, 500)
- [ ] Update API documentation with error response examples

## Benefits After Migration

✅ **Less Code**: ~50% reduction in error handling boilerplate  
✅ **Consistency**: All errors follow the same format  
✅ **Maintainability**: Single source of truth for error handling  
✅ **Debugging**: Better error logging with context  
✅ **Type Safety**: Proper error type checking  
✅ **User Experience**: Clean, user-friendly error messages

