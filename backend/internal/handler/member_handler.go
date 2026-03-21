package handler

import (
	"time"
	"vamsasetu/backend/internal/middleware"
	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

// MemberHandler handles member-related HTTP requests
type MemberHandler struct {
	memberService *service.MemberService
}

// NewMemberHandler creates a new member handler instance
func NewMemberHandler(memberService *service.MemberService) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
	}
}

// RegisterRoutes registers all member routes with the Fiber app
func (h *MemberHandler) RegisterRoutes(app *fiber.App) {
	members := app.Group("/api/members")

	// All member endpoints require authentication
	members.Use(middleware.AuthMiddleware())

	// Member CRUD endpoints
	members.Get("/", h.ListMembers)
	members.Post("/", middleware.RequireRole("owner", "admin"), h.CreateMember)
	members.Get("/:id", h.GetMember)
	members.Put("/:id", middleware.RequireRole("owner", "admin"), h.UpdateMember)
	members.Delete("/:id", middleware.RequireRole("owner", "admin"), h.DeleteMember)
}

// CreateMember handles member creation
// POST /api/members
// Request body: { "name": string, "dateOfBirth": string, "gender": string, "email": string, "phone": string, "avatarUrl": string }
// Response: { "success": bool, "data": Member, "error": string }
func (h *MemberHandler) CreateMember(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		DateOfBirth string `json:"dateOfBirth"`
		Gender      string `json:"gender"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		AvatarURL   string `json:"avatarUrl"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	// Parse date of birth
	dateOfBirth, err := time.Parse(time.RFC3339, req.DateOfBirth)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid date format. Expected RFC3339 format (e.g., 2006-01-02T15:04:05Z)",
		})
	}

	// Create member model
	member := models.NewMember(req.Name, dateOfBirth, req.Gender)
	member.Email = req.Email
	member.Phone = req.Phone
	member.AvatarURL = req.AvatarURL

	// Call member service
	if err := h.memberService.Create(c.Context(), member); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    member,
		"error":   "",
	})
}

// GetMember retrieves a member by ID with relationships
// GET /api/members/:id
// Response: { "success": bool, "data": Member, "error": string }
func (h *MemberHandler) GetMember(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Member ID is required",
		})
	}

	// Get member from service
	member, err := h.memberService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Member not found",
		})
	}

	// Return member
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    member,
		"error":   "",
	})
}

// UpdateMember updates an existing member
// PUT /api/members/:id
// Request body: { "name": string, "dateOfBirth": string, "gender": string, "email": string, "phone": string, "avatarUrl": string }
// Response: { "success": bool, "data": Member, "error": string }
func (h *MemberHandler) UpdateMember(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Member ID is required",
		})
	}

	var req struct {
		Name        string `json:"name"`
		DateOfBirth string `json:"dateOfBirth"`
		Gender      string `json:"gender"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		AvatarURL   string `json:"avatarUrl"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	// Get existing member
	member, err := h.memberService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Member not found",
		})
	}

	// Update member fields
	member.Name = req.Name
	member.Gender = req.Gender
	member.Email = req.Email
	member.Phone = req.Phone
	member.AvatarURL = req.AvatarURL

	// Parse date of birth if provided
	if req.DateOfBirth != "" {
		dateOfBirth, err := time.Parse(time.RFC3339, req.DateOfBirth)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Invalid date format. Expected RFC3339 format (e.g., 2006-01-02T15:04:05Z)",
			})
		}
		member.DateOfBirth = dateOfBirth
	}

	// Call member service
	if err := h.memberService.Update(c.Context(), member); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return updated member
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    member,
		"error":   "",
	})
}

// DeleteMember soft deletes a member
// DELETE /api/members/:id
// Response: { "success": bool, "data": { "message": string }, "error": string }
func (h *MemberHandler) DeleteMember(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Member ID is required",
		})
	}

	// Call member service
	if err := h.memberService.SoftDelete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "Member deleted successfully",
		},
		"error": "",
	})
}

// ListMembers retrieves all members with pagination and filters
// GET /api/members?page=1&limit=10&search=name&gender=male
// Response: { "success": bool, "data": { "members": []Member, "total": int, "page": int, "limit": int }, "error": string }
func (h *MemberHandler) ListMembers(c *fiber.Ctx) error {
	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	search := c.Query("search", "")
	gender := c.Query("gender", "")

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// Get members from service
	var members []*models.Member
	var err error

	if search != "" {
		// Search by name
		members, err = h.memberService.Search(c.Context(), search)
	} else {
		// Get all members
		members, err = h.memberService.GetAll(c.Context())
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Filter by gender if specified
	if gender != "" {
		var filtered []*models.Member
		for _, member := range members {
			if member.Gender == gender {
				filtered = append(filtered, member)
			}
		}
		members = filtered
	}

	// Calculate pagination
	total := len(members)
	start := (page - 1) * limit
	end := start + limit

	// Handle out of bounds
	if start >= total {
		members = []*models.Member{}
	} else {
		if end > total {
			end = total
		}
		members = members[start:end]
	}

	// Return paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"members": members,
			"total":   total,
			"page":    page,
			"limit":   limit,
		},
		"error": "",
	})
}
