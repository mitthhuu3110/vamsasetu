package service

import (
	"context"
	"testing"
	"time"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
	"vamsasetu/backend/pkg/neo4j"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock repositories for testing
type MockMemberRepository struct {
	mock.Mock
}

func (m *MockMemberRepository) GetAll(ctx context.Context) ([]*models.Member, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Member), args.Error(1)
}

func (m *MockMemberRepository) Create(ctx context.Context, member *models.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockMemberRepository) GetByID(ctx context.Context, id string) (*models.Member, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Member), args.Error(1)
}

func (m *MockMemberRepository) Update(ctx context.Context, member *models.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockMemberRepository) SoftDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMemberRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockRelationshipRepository struct {
	mock.Mock
}

func (m *MockRelationshipRepository) GetAll(ctx context.Context) ([]*models.Relationship, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Relationship), args.Error(1)
}

func (m *MockRelationshipRepository) Create(ctx context.Context, rel *models.Relationship) error {
	args := m.Called(ctx, rel)
	return args.Error(0)
}

func (m *MockRelationshipRepository) Delete(ctx context.Context, fromID, toID, relType string) error {
	args := m.Called(ctx, fromID, toID, relType)
	return args.Error(0)
}

func (m *MockRelationshipRepository) FindPath(ctx context.Context, fromID, toID string) (*repository.RelationshipPath, error) {
	args := m.Called(ctx, fromID, toID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.RelationshipPath), args.Error(1)
}

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) GetUpcoming(ctx context.Context, startDate, endDate time.Time) ([]*models.Event, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRepository) Create(ctx context.Context, event *models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) GetByID(ctx context.Context, id uint) (*models.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventRepository) GetAll(ctx context.Context) ([]*models.Event, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRepository) Update(ctx context.Context, event *models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventRepository) GetByType(ctx context.Context, eventType string) ([]*models.Event, error) {
	args := m.Called(ctx, eventType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRepository) GetByMember(ctx context.Context, memberID string) ([]*models.Event, error) {
	args := m.Called(ctx, memberID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.Event, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventRepository) GetByCreator(ctx context.Context, userID uint) ([]*models.Event, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

// Test helper functions
func createTestMember(id, name, gender string) *models.Member {
	return &models.Member{
		ID:          id,
		Name:        name,
		Gender:      gender,
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}
}

func createTestRelationship(fromID, toID, relType string) *models.Relationship {
	return &models.Relationship{
		FromID:    fromID,
		ToID:      toID,
		Type:      relType,
		CreatedAt: time.Now(),
	}
}

func createTestEvent(memberIDs []string, eventDate time.Time) *models.Event {
	return &models.Event{
		ID:        1,
		Title:     "Test Event",
		EventDate: eventDate,
		EventType: "birthday",
		MemberIDs: memberIDs,
		CreatedBy: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Test BuildTree with simple family structure
func TestBuildTree_SimpleFamily(t *testing.T) {
	// Setup
	mockMemberRepo := new(MockMemberRepository)
	mockRelationshipRepo := new(MockRelationshipRepository)
	mockEventRepo := new(MockEventRepository)

	// Create test data: Parent -> Child
	parent := createTestMember("parent-1", "Parent", "male")
	child := createTestMember("child-1", "Child", "female")
	members := []*models.Member{parent, child}

	parentChildRel := createTestRelationship("parent-1", "child-1", models.RelationshipTypeParentOf)
	relationships := []*models.Relationship{parentChildRel}

	// Mock expectations
	mockMemberRepo.On("GetAll", mock.Anything).Return(members, nil)
	mockRelationshipRepo.On("GetAll", mock.Anything).Return(relationships, nil)
	mockEventRepo.On("GetUpcoming", mock.Anything, mock.Anything, mock.Anything).Return([]*models.Event{}, nil)

	// Create tree builder
	tb := NewTreeBuilder(mockMemberRepo, mockRelationshipRepo, mockEventRepo)

	// Execute
	tree, err := tb.BuildTree(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 2, len(tree.Nodes))
	assert.Equal(t, 1, len(tree.Edges))

	// Verify nodes have positions
	for _, node := range tree.Nodes {
		assert.NotEqual(t, 0.0, node.Position.X)
		assert.NotEqual(t, 0.0, node.Position.Y)
	}

	// Verify edge color
	assert.Equal(t, ParentEdgeColor, tree.Edges[0].Style["stroke"])
}

// Test BuildTree with multiple generations
func TestBuildTree_MultipleGenerations(t *testing.T) {
	// Setup
	mockMemberRepo := new(MockMemberRepository)
	mockRelationshipRepo := new(MockRelationshipRepository)
	mockEventRepo := new(MockEventRepository)

	// Create test data: Grandparent -> Parent -> Child
	grandparent := createTestMember("gp-1", "Grandparent", "male")
	parent := createTestMember("parent-1", "Parent", "female")
	child := createTestMember("child-1", "Child", "male")
	members := []*models.Member{grandparent, parent, child}

	gpToParent := createTestRelationship("gp-1", "parent-1", models.RelationshipTypeParentOf)
	parentToChild := createTestRelationship("parent-1", "child-1", models.RelationshipTypeParentOf)
	relationships := []*models.Relationship{gpToParent, parentToChild}

	// Mock expectations
	mockMemberRepo.On("GetAll", mock.Anything).Return(members, nil)
	mockRelationshipRepo.On("GetAll", mock.Anything).Return(relationships, nil)
	mockEventRepo.On("GetUpcoming", mock.Anything, mock.Anything, mock.Anything).Return([]*models.Event{}, nil)

	// Create tree builder
	tb := NewTreeBuilder(mockMemberRepo, mockRelationshipRepo, mockEventRepo)

	// Execute
	tree, err := tb.BuildTree(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 3, len(tree.Nodes))
	assert.Equal(t, 2, len(tree.Edges))

	// Verify generation spacing (Y coordinates should increase)
	nodeMap := make(map[string]ReactFlowNode)
	for _, node := range tree.Nodes {
		nodeMap[node.ID] = node
	}

	assert.True(t, nodeMap["parent-1"].Position.Y > nodeMap["gp-1"].Position.Y)
	assert.True(t, nodeMap["child-1"].Position.Y > nodeMap["parent-1"].Position.Y)
}

// Test BuildTree with spouse relationships
func TestBuildTree_WithSpouses(t *testing.T) {
	// Setup
	mockMemberRepo := new(MockMemberRepository)
	mockRelationshipRepo := new(MockRelationshipRepository)
	mockEventRepo := new(MockEventRepository)

	// Create test data: Husband <-> Wife
	husband := createTestMember("husband-1", "Husband", "male")
	wife := createTestMember("wife-1", "Wife", "female")
	members := []*models.Member{husband, wife}

	spouseRel := createTestRelationship("husband-1", "wife-1", models.RelationshipTypeSpouseOf)
	relationships := []*models.Relationship{spouseRel}

	// Mock expectations
	mockMemberRepo.On("GetAll", mock.Anything).Return(members, nil)
	mockRelationshipRepo.On("GetAll", mock.Anything).Return(relationships, nil)
	mockEventRepo.On("GetUpcoming", mock.Anything, mock.Anything, mock.Anything).Return([]*models.Event{}, nil)

	// Create tree builder
	tb := NewTreeBuilder(mockMemberRepo, mockRelationshipRepo, mockEventRepo)

	// Execute
	tree, err := tb.BuildTree(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 2, len(tree.Nodes))
	assert.Equal(t, 1, len(tree.Edges))

	// Verify edge color for spouse
	assert.Equal(t, SpouseEdgeColor, tree.Edges[0].Style["stroke"])

	// Verify spouses are on same Y level (same generation)
	nodeMap := make(map[string]ReactFlowNode)
	for _, node := range tree.Nodes {
		nodeMap[node.ID] = node
	}
	assert.Equal(t, nodeMap["husband-1"].Position.Y, nodeMap["wife-1"].Position.Y)
}

// Test BuildTree with upcoming events
func TestBuildTree_WithUpcomingEvents(t *testing.T) {
	// Setup
	mockMemberRepo := new(MockMemberRepository)
	mockRelationshipRepo := new(MockRelationshipRepository)
	mockEventRepo := new(MockEventRepository)

	// Create test data
	member := createTestMember("member-1", "Member", "male")
	members := []*models.Member{member}

	// Create upcoming event (within 7 days)
	upcomingEvent := createTestEvent([]string{"member-1"}, time.Now().AddDate(0, 0, 3))
	events := []*models.Event{upcomingEvent}

	// Mock expectations
	mockMemberRepo.On("GetAll", mock.Anything).Return(members, nil)
	mockRelationshipRepo.On("GetAll", mock.Anything).Return([]*models.Relationship{}, nil)
	mockEventRepo.On("GetUpcoming", mock.Anything, mock.Anything, mock.Anything).Return(events, nil)

	// Create tree builder
	tb := NewTreeBuilder(mockMemberRepo, mockRelationshipRepo, mockEventRepo)

	// Execute
	tree, err := tb.BuildTree(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 1, len(tree.Nodes))
	assert.True(t, tree.Nodes[0].Data.HasUpcomingEvent)
}

// Test BuildTree with sibling relationships
func TestBuildTree_WithSiblings(t *testing.T) {
	// Setup
	mockMemberRepo := new(MockMemberRepository)
	mockRelationshipRepo := new(MockRelationshipRepository)
	mockEventRepo := new(MockEventRepository)

	// Create test data: Sibling1 <-> Sibling2
	sibling1 := createTestMember("sibling-1", "Sibling1", "male")
	sibling2 := createTestMember("sibling-2", "Sibling2", "female")
	members := []*models.Member{sibling1, sibling2}

	siblingRel := createTestRelationship("sibling-1", "sibling-2", models.RelationshipTypeSiblingOf)
	relationships := []*models.Relationship{siblingRel}

	// Mock expectations
	mockMemberRepo.On("GetAll", mock.Anything).Return(members, nil)
	mockRelationshipRepo.On("GetAll", mock.Anything).Return(relationships, nil)
	mockEventRepo.On("GetUpcoming", mock.Anything, mock.Anything, mock.Anything).Return([]*models.Event{}, nil)

	// Create tree builder
	tb := NewTreeBuilder(mockMemberRepo, mockRelationshipRepo, mockEventRepo)

	// Execute
	tree, err := tb.BuildTree(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 2, len(tree.Nodes))
	assert.Equal(t, 1, len(tree.Edges))

	// Verify edge color for sibling
	assert.Equal(t, SiblingEdgeColor, tree.Edges[0].Style["stroke"])

	// Verify siblings are on same Y level
	nodeMap := make(map[string]ReactFlowNode)
	for _, node := range tree.Nodes {
		nodeMap[node.ID] = node
	}
	assert.Equal(t, nodeMap["sibling-1"].Position.Y, nodeMap["sibling-2"].Position.Y)
}

// Test edge color mapping
func TestGetEdgeColor(t *testing.T) {
	tb := &TreeBuilder{}

	tests := []struct {
		relType  string
		expected string
	}{
		{models.RelationshipTypeSpouseOf, SpouseEdgeColor},
		{models.RelationshipTypeParentOf, ParentEdgeColor},
		{models.RelationshipTypeSiblingOf, SiblingEdgeColor},
		{"UNKNOWN", "#6B7280"},
	}

	for _, tt := range tests {
		t.Run(tt.relType, func(t *testing.T) {
			color := tb.getEdgeColor(tt.relType)
			assert.Equal(t, tt.expected, color)
		})
	}
}

// Test collision detection and spacing
func TestOptimizeSpacing(t *testing.T) {
	tb := &TreeBuilder{}

	// Create overlapping layout
	layout := map[string]NodeLayout{
		"member-1": {X: 0, Y: 0, Generation: 0},
		"member-2": {X: 50, Y: 0, Generation: 0}, // Too close (should be at least 200 apart)
		"member-3": {X: 100, Y: 0, Generation: 0},
	}

	graph := &Graph{
		Members: map[string]*models.Member{
			"member-1": createTestMember("member-1", "Member1", "male"),
			"member-2": createTestMember("member-2", "Member2", "female"),
			"member-3": createTestMember("member-3", "Member3", "male"),
		},
	}

	// Execute
	tb.optimizeSpacing(layout, graph)

	// Assert - nodes should be spaced apart
	assert.True(t, layout["member-2"].X-layout["member-1"].X >= NodeWidth+HorizontalSpace)
	assert.True(t, layout["member-3"].X-layout["member-2"].X >= NodeWidth+HorizontalSpace)
}
