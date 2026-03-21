package repository

import (
	"testing"
	"time"

	"vamsasetu/backend/internal/models"
)

func TestRelationship_IsBidirectional(t *testing.T) {
	tests := []struct {
		name     string
		relType  string
		expected bool
	}{
		{
			name:     "SPOUSE_OF is bidirectional",
			relType:  models.RelationshipTypeSpouseOf,
			expected: true,
		},
		{
			name:     "SIBLING_OF is bidirectional",
			relType:  models.RelationshipTypeSiblingOf,
			expected: true,
		},
		{
			name:     "PARENT_OF is not bidirectional",
			relType:  models.RelationshipTypeParentOf,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rel := &models.Relationship{Type: tt.relType}
			if got := rel.IsBidirectional(); got != tt.expected {
				t.Errorf("IsBidirectional() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRelationshipRepository_makeRelationshipKey(t *testing.T) {
	repo := &RelationshipRepository{}

	tests := []struct {
		name     string
		rel      *models.Relationship
		expected string
	}{
		{
			name: "FromID < ToID",
			rel: &models.Relationship{
				Type:   models.RelationshipTypeSpouseOf,
				FromID: "aaa",
				ToID:   "bbb",
			},
			expected: "aaa-bbb-SPOUSE_OF",
		},
		{
			name: "FromID > ToID",
			rel: &models.Relationship{
				Type:   models.RelationshipTypeSpouseOf,
				FromID: "zzz",
				ToID:   "aaa",
			},
			expected: "aaa-zzz-SPOUSE_OF",
		},
		{
			name: "Same IDs produce consistent key",
			rel: &models.Relationship{
				Type:   models.RelationshipTypeSiblingOf,
				FromID: "xxx",
				ToID:   "yyy",
			},
			expected: "xxx-yyy-SIBLING_OF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := repo.makeRelationshipKey(tt.rel); got != tt.expected {
				t.Errorf("makeRelationshipKey() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRelationshipRepository_parseNeo4jTime(t *testing.T) {
	repo := &RelationshipRepository{}

	now := time.Now()
	rfc3339Str := now.Format(time.RFC3339)

	tests := []struct {
		name      string
		value     interface{}
		wantError bool
	}{
		{
			name:      "time.Time value",
			value:     now,
			wantError: false,
		},
		{
			name:      "RFC3339 string",
			value:     rfc3339Str,
			wantError: false,
		},
		{
			name:      "invalid type",
			value:     123,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.parseNeo4jTime(tt.value)
			if (err != nil) != tt.wantError {
				t.Errorf("parseNeo4jTime() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestPathNode_Structure(t *testing.T) {
	node := PathNode{
		ID:     "test-id",
		Name:   "Test Name",
		Gender: "male",
	}

	if node.ID != "test-id" {
		t.Errorf("PathNode.ID = %v, want %v", node.ID, "test-id")
	}
	if node.Name != "Test Name" {
		t.Errorf("PathNode.Name = %v, want %v", node.Name, "Test Name")
	}
	if node.Gender != "male" {
		t.Errorf("PathNode.Gender = %v, want %v", node.Gender, "male")
	}
}

func TestRelationshipPath_Structure(t *testing.T) {
	path := RelationshipPath{
		Nodes: []PathNode{
			{ID: "1", Name: "Node1", Gender: "male"},
			{ID: "2", Name: "Node2", Gender: "female"},
		},
		Relationships: []string{"PARENT_OF"},
		Length:       1,
	}

	if len(path.Nodes) != 2 {
		t.Errorf("RelationshipPath.Nodes length = %v, want %v", len(path.Nodes), 2)
	}
	if len(path.Relationships) != 1 {
		t.Errorf("RelationshipPath.Relationships length = %v, want %v", len(path.Relationships), 1)
	}
	if path.Length != 1 {
		t.Errorf("RelationshipPath.Length = %v, want %v", path.Length, 1)
	}
}
