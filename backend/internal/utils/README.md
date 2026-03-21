# API Response Utilities

This package provides standardized API response utilities for the VamsaSetu backend.

## Overview

All API responses in VamsaSetu follow a consistent format to ensure predictable client-side handling:

```json
{
  "success": boolean,
  "data": any | null,
  "error": string
}
```

## APIResponse Structure

```go
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Error   string      `json:"error"`
}
```

### Fields

- **success**: Boolean indicating whether the request was successful
- **data**: Contains the response payload when successful, `null` when failed
- **error**: Contains error message when failed, empty string when successful

## Helper Functions

### Success Responses

#### SuccessResponse
Returns HTTP 200 OK with success data.

```go
func SuccessResponse(c *fiber.Ctx, data interface{}) error
```

**Example:**
```go
return utils.SuccessResponse(c, fiber.Map{
    "message": "Operation successful",
    "userId": 123,
})
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Operation successful",
    "userId": 123
  },
  "error": ""
}
```

#### CreatedResponse
Returns HTTP 201 Created for resource creation.

```go
func CreatedResponse(c *fiber.Ctx, data interface{}) error
```

**Example:**
```go
return utils.CreatedResponse(c, newMember)
```

### Error Responses

#### ErrorResponse
Generic error response with custom status code.

```go
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error
```

**Example:**
```go
return utils.ErrorResponse(c, fiber.StatusConflict, "Resource already exists")
```

#### BadRequestResponse (400)
For validation errors and malformed requests.

```go
func BadRequestResponse(c *fiber.Ctx, message string) error
```

**Example:**
```go
return utils.BadRequestResponse(c, "Email is required")
```

**Response:**
```json
{
  "success": false,
  "data": null,
  "error": "Email is required"
}
```

#### UnauthorizedResponse (401)
For authentication failures.

```go
func UnauthorizedResponse(c *fiber.Ctx, message string) error
```

**Example:**
```go
return utils.UnauthorizedResponse(c, "Invalid or expired token")
```

#### ForbiddenResponse (403)
For authorization failures.

```go
func ForbiddenResponse(c *fiber.Ctx, message string) error
```

**Example:**
```go
return utils.ForbiddenResponse(c, "Insufficient permissions")
```

#### NotFoundResponse (404)
For missing resources.

```go
func NotFoundResponse(c *fiber.Ctx, message string) error
```

**Example:**
```go
return utils.NotFoundResponse(c, "Member not found")
```

#### InternalServerErrorResponse (500)
For server-side errors.

```go
func InternalServerErrorResponse(c *fiber.Ctx, message string) error
```

**Example:**
```go
return utils.InternalServerErrorResponse(c, "Database connection failed")
```

## HTTP Status Code Mapping

| Function | Status Code | Use Case |
|----------|-------------|----------|
| SuccessResponse | 200 OK | Successful GET, PUT, DELETE |
| CreatedResponse | 201 Created | Successful POST (resource created) |
| BadRequestResponse | 400 Bad Request | Validation errors, malformed input |
| UnauthorizedResponse | 401 Unauthorized | Missing or invalid authentication |
| ForbiddenResponse | 403 Forbidden | Valid auth but insufficient permissions |
| NotFoundResponse | 404 Not Found | Resource doesn't exist |
| InternalServerErrorResponse | 500 Internal Server Error | Server-side failures |

## Usage in Handlers

### Basic Handler Example

```go
func (h *MemberHandler) GetMember(c *fiber.Ctx) error {
    id := c.Params("id")
    
    // Validate input
    if id == "" {
        return utils.BadRequestResponse(c, "Member ID is required")
    }
    
    // Fetch member
    member, err := h.memberService.GetByID(c.Context(), id)
    if err != nil {
        return utils.InternalServerErrorResponse(c, "Failed to fetch member")
    }
    
    if member == nil {
        return utils.NotFoundResponse(c, "Member not found")
    }
    
    return utils.SuccessResponse(c, member)
}
```

### Create Handler Example

```go
func (h *MemberHandler) CreateMember(c *fiber.Ctx) error {
    var req CreateMemberRequest
    
    // Parse request body
    if err := c.BodyParser(&req); err != nil {
        return utils.BadRequestResponse(c, "Invalid request body")
    }
    
    // Validate required fields
    if req.Name == "" {
        return utils.BadRequestResponse(c, "Name is required")
    }
    
    // Create member
    member, err := h.memberService.Create(c.Context(), &req)
    if err != nil {
        return utils.InternalServerErrorResponse(c, err.Error())
    }
    
    return utils.CreatedResponse(c, member)
}
```

### Complex Data Response Example

```go
func (h *EventHandler) ListEvents(c *fiber.Ctx) error {
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    
    events, total, err := h.eventService.GetPaginated(c.Context(), page, limit)
    if err != nil {
        return utils.InternalServerErrorResponse(c, "Failed to fetch events")
    }
    
    return utils.SuccessResponse(c, fiber.Map{
        "events": events,
        "pagination": fiber.Map{
            "total": total,
            "page": page,
            "limit": limit,
            "totalPages": (total + limit - 1) / limit,
        },
    })
}
```

## Design Principles

1. **Consistency**: All endpoints use the same response format
2. **Predictability**: Clients can always expect `success`, `data`, and `error` fields
3. **Type Safety**: TypeScript clients can define a single response interface
4. **Error Handling**: Clear error messages help with debugging
5. **HTTP Standards**: Proper status codes follow REST conventions

## Testing

The response utilities include comprehensive unit tests. Run tests with:

```bash
go test -v ./internal/utils/
```

## Requirements Satisfied

This implementation satisfies **Requirement 13: API Response Consistency** from the requirements document:

- ✅ All API responses follow the format: `{ "success": bool, "data": any, "error": string }`
- ✅ Success responses set `success=true`, populate `data`, and set `error=""` 
- ✅ Error responses set `success=false`, set `data=null`, and populate `error`
- ✅ Appropriate HTTP status codes are used (200, 201, 400, 401, 403, 404, 500)

## Related Files

- `response.go` - Main implementation
- `response_test.go` - Unit tests
- `response_example.go` - Usage examples
