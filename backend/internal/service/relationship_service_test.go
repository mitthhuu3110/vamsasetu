package service

import (
	"testing"

	"vamsasetu/backend/internal/models"
	"vamsasetu/backend/internal/repository"
)

func TestMapDirectRelationship(t *testing.T) {
	service := &RelationshipService{}

	tests := []struct {
		name              string
		path              *repository.RelationshipPath
		expectedLabel     string
		expectedKinship   string
	}{
		{
			name: "Parent to Son",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Father", Gender: "male"},
					{ID: "2", Name: "Son", Gender: "male"},
				},
				Relationships: []string{models.RelationshipTypeParentOf},
			},
			expectedLabel:   "Son",
			expectedKinship: "Koduku",
		},
		{
			name: "Parent to Daughter",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Mother", Gender: "female"},
					{ID: "2", Name: "Daughter", Gender: "female"},
				},
				Relationships: []string{models.RelationshipTypeParentOf},
			},
			expectedLabel:   "Daughter",
			expectedKinship: "Kuthuru",
		},
		{
			name: "Spouse - Husband",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Wife", Gender: "female"},
					{ID: "2", Name: "Husband", Gender: "male"},
				},
				Relationships: []string{models.RelationshipTypeSpouseOf},
			},
			expectedLabel:   "Husband",
			expectedKinship: "Menarikam",
		},
		{
			name: "Spouse - Wife",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Husband", Gender: "male"},
					{ID: "2", Name: "Wife", Gender: "female"},
				},
				Relationships: []string{models.RelationshipTypeSpouseOf},
			},
			expectedLabel:   "Wife",
			expectedKinship: "Bharya",
		},
		{
			name: "Sibling - Brother",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Sister", Gender: "female"},
					{ID: "2", Name: "Brother", Gender: "male"},
				},
				Relationships: []string{models.RelationshipTypeSiblingOf},
			},
			expectedLabel:   "Brother",
			expectedKinship: "Annayya/Tammudu",
		},
		{
			name: "Sibling - Sister",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Brother", Gender: "male"},
					{ID: "2", Name: "Sister", Gender: "female"},
				},
				Relationships: []string{models.RelationshipTypeSiblingOf},
			},
			expectedLabel:   "Sister",
			expectedKinship: "Akka/Chelli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			label, kinship := service.mapDirectRelationship(tt.path)
			if label != tt.expectedLabel {
				t.Errorf("mapDirectRelationship() label = %v, want %v", label, tt.expectedLabel)
			}
			if kinship != tt.expectedKinship {
				t.Errorf("mapDirectRelationship() kinship = %v, want %v", kinship, tt.expectedKinship)
			}
		})
	}
}

func TestMapTwoHopRelationship(t *testing.T) {
	service := &RelationshipService{}

	tests := []struct {
		name              string
		path              *repository.RelationshipPath
		expectedLabel     string
		expectedKinship   string
	}{
		{
			name: "Grandson",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Grandparent", Gender: "male"},
					{ID: "2", Name: "Parent", Gender: "male"},
					{ID: "3", Name: "Grandson", Gender: "male"},
				},
				Relationships: []string{models.RelationshipTypeParentOf, models.RelationshipTypeParentOf},
			},
			expectedLabel:   "Grandson",
			expectedKinship: "Manumadu",
		},
		{
			name: "Granddaughter",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Grandparent", Gender: "female"},
					{ID: "2", Name: "Parent", Gender: "female"},
					{ID: "3", Name: "Granddaughter", Gender: "female"},
				},
				Relationships: []string{models.RelationshipTypeParentOf, models.RelationshipTypeParentOf},
			},
			expectedLabel:   "Granddaughter",
			expectedKinship: "Manumalu",
		},
		{
			name: "Uncle - Father's Brother",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Child", Gender: "male"},
					{ID: "2", Name: "Father", Gender: "male"},
					{ID: "3", Name: "Uncle", Gender: "male"},
				},
				Relationships: []string{models.RelationshipTypeParentOf, models.RelationshipTypeSiblingOf},
			},
			expectedLabel:   "Uncle (Father's Brother)",
			expectedKinship: "Babai",
		},
		{
			name: "Aunt - Father's Sister",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Child", Gender: "female"},
					{ID: "2", Name: "Father", Gender: "male"},
					{ID: "3", Name: "Aunt", Gender: "female"},
				},
				Relationships: []string{models.RelationshipTypeParentOf, models.RelationshipTypeSiblingOf},
			},
			expectedLabel:   "Aunt (Father's Sister)",
			expectedKinship: "Attha",
		},
		{
			name: "Uncle - Mother's Brother",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Child", Gender: "male"},
					{ID: "2", Name: "Mother", Gender: "female"},
					{ID: "3", Name: "Uncle", Gender: "male"},
				},
				Relationships: []string{models.RelationshipTypeParentOf, models.RelationshipTypeSiblingOf},
			},
			expectedLabel:   "Uncle (Mother's Brother)",
			expectedKinship: "Mamayya",
		},
		{
			name: "Aunt - Mother's Sister",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Child", Gender: "female"},
					{ID: "2", Name: "Mother", Gender: "female"},
					{ID: "3", Name: "Aunt", Gender: "female"},
				},
				Relationships: []string{models.RelationshipTypeParentOf, models.RelationshipTypeSiblingOf},
			},
			expectedLabel:   "Aunt (Mother's Sister)",
			expectedKinship: "Pinni",
		},
		{
			name: "Father-in-Law",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Person", Gender: "female"},
					{ID: "2", Name: "Spouse", Gender: "male"},
					{ID: "3", Name: "Father-in-Law", Gender: "male"},
				},
				Relationships: []string{models.RelationshipTypeSpouseOf, models.RelationshipTypeParentOf},
			},
			expectedLabel:   "Father-in-Law",
			expectedKinship: "Maamayyagaru",
		},
		{
			name: "Mother-in-Law",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Person", Gender: "male"},
					{ID: "2", Name: "Spouse", Gender: "female"},
					{ID: "3", Name: "Mother-in-Law", Gender: "female"},
				},
				Relationships: []string{models.RelationshipTypeSpouseOf, models.RelationshipTypeParentOf},
			},
			expectedLabel:   "Mother-in-Law",
			expectedKinship: "Attagaru",
		},
		{
			name: "Brother-in-Law (Sibling's Spouse)",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Person", Gender: "female"},
					{ID: "2", Name: "Sibling", Gender: "female"},
					{ID: "3", Name: "Brother-in-Law", Gender: "male"},
				},
				Relationships: []string{models.RelationshipTypeSiblingOf, models.RelationshipTypeSpouseOf},
			},
			expectedLabel:   "Brother-in-Law",
			expectedKinship: "Bava",
		},
		{
			name: "Sister-in-Law (Sibling's Spouse)",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Person", Gender: "male"},
					{ID: "2", Name: "Sibling", Gender: "male"},
					{ID: "3", Name: "Sister-in-Law", Gender: "female"},
				},
				Relationships: []string{models.RelationshipTypeSiblingOf, models.RelationshipTypeSpouseOf},
			},
			expectedLabel:   "Sister-in-Law",
			expectedKinship: "Vadina",
		},
		{
			name: "Nephew",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Person", Gender: "female"},
					{ID: "2", Name: "Sibling", Gender: "male"},
					{ID: "3", Name: "Nephew", Gender: "male"},
				},
				Relationships: []string{models.RelationshipTypeSiblingOf, models.RelationshipTypeParentOf},
			},
			expectedLabel:   "Nephew",
			expectedKinship: "Bhanja/Alludu",
		},
		{
			name: "Niece",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Person", Gender: "male"},
					{ID: "2", Name: "Sibling", Gender: "female"},
					{ID: "3", Name: "Niece", Gender: "female"},
				},
				Relationships: []string{models.RelationshipTypeSiblingOf, models.RelationshipTypeParentOf},
			},
			expectedLabel:   "Niece",
			expectedKinship: "Bhanjika/Kodalu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			label, kinship := service.mapTwoHopRelationship(tt.path)
			if label != tt.expectedLabel {
				t.Errorf("mapTwoHopRelationship() label = %v, want %v", label, tt.expectedLabel)
			}
			if kinship != tt.expectedKinship {
				t.Errorf("mapTwoHopRelationship() kinship = %v, want %v", kinship, tt.expectedKinship)
			}
		})
	}
}

func TestGenerateDescription(t *testing.T) {
	service := &RelationshipService{}

	tests := []struct {
		name        string
		path        *repository.RelationshipPath
		label       string
		expected    string
	}{
		{
			name: "Direct relationship",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Rajesh", Gender: "male"},
					{ID: "2", Name: "Arjun", Gender: "male"},
				},
				Relationships: []string{models.RelationshipTypeParentOf},
			},
			label:    "Son",
			expected: "Arjun is Rajesh's son",
		},
		{
			name: "Two-hop relationship",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Arjun", Gender: "male"},
					{ID: "2", Name: "Rajesh", Gender: "male"},
					{ID: "3", Name: "Lakshmi", Gender: "female"},
				},
				Relationships: []string{models.RelationshipTypeParentOf, models.RelationshipTypeSpouseOf},
			},
			label:    "Mother",
			expected: "Lakshmi is Arjun's mother through Rajesh",
		},
		{
			name: "Multi-hop relationship",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Person1", Gender: "male"},
					{ID: "2", Name: "Person2", Gender: "male"},
					{ID: "3", Name: "Person3", Gender: "female"},
					{ID: "4", Name: "Person4", Gender: "male"},
				},
				Relationships: []string{
					models.RelationshipTypeParentOf,
					models.RelationshipTypeSiblingOf,
					models.RelationshipTypeParentOf,
				},
			},
			label:    "Cousin",
			expected: "Person4 is Person1's cousin (3 connections)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.generateDescription(tt.path, tt.label)
			if result != tt.expected {
				t.Errorf("generateDescription() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetGreatPrefix(t *testing.T) {
	service := &RelationshipService{}

	tests := []struct {
		name     string
		count    int
		expected string
	}{
		{
			name:     "No prefix",
			count:    0,
			expected: "",
		},
		{
			name:     "One Great",
			count:    1,
			expected: "Great-",
		},
		{
			name:     "Two Greats",
			count:    2,
			expected: "Great-Great-",
		},
		{
			name:     "Three Greats",
			count:    3,
			expected: "Great-Great-Great-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.getGreatPrefix(tt.count)
			if result != tt.expected {
				t.Errorf("getGreatPrefix() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsCousinPath(t *testing.T) {
	service := &RelationshipService{}

	tests := []struct {
		name     string
		path     *repository.RelationshipPath
		expected bool
	}{
		{
			name: "Cousin path with parent and sibling",
			path: &repository.RelationshipPath{
				Relationships: []string{
					models.RelationshipTypeParentOf,
					models.RelationshipTypeSiblingOf,
					models.RelationshipTypeParentOf,
				},
			},
			expected: true,
		},
		{
			name: "Direct parent path",
			path: &repository.RelationshipPath{
				Relationships: []string{
					models.RelationshipTypeParentOf,
				},
			},
			expected: false,
		},
		{
			name: "Sibling only path",
			path: &repository.RelationshipPath{
				Relationships: []string{
					models.RelationshipTypeSiblingOf,
				},
			},
			expected: false,
		},
		{
			name: "Parent only path",
			path: &repository.RelationshipPath{
				Relationships: []string{
					models.RelationshipTypeParentOf,
					models.RelationshipTypeParentOf,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.isCousinPath(tt.path)
			if result != tt.expected {
				t.Errorf("isCousinPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateCousinDegree(t *testing.T) {
	service := &RelationshipService{}

	tests := []struct {
		name     string
		path     *repository.RelationshipPath
		expected int
	}{
		{
			name: "First cousin (1 generation up)",
			path: &repository.RelationshipPath{
				Relationships: []string{
					models.RelationshipTypeParentOf,
					models.RelationshipTypeSiblingOf,
					models.RelationshipTypeParentOf,
				},
			},
			expected: 1,
		},
		{
			name: "Second cousin (2 generations up)",
			path: &repository.RelationshipPath{
				Relationships: []string{
					models.RelationshipTypeParentOf,
					models.RelationshipTypeParentOf,
					models.RelationshipTypeSiblingOf,
					models.RelationshipTypeParentOf,
					models.RelationshipTypeParentOf,
				},
			},
			expected: 2,
		},
		{
			name: "Direct sibling (0 generations up)",
			path: &repository.RelationshipPath{
				Relationships: []string{
					models.RelationshipTypeSiblingOf,
				},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.calculateCousinDegree(tt.path)
			if result != tt.expected {
				t.Errorf("calculateCousinDegree() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMapMultiHopRelationship(t *testing.T) {
	service := &RelationshipService{}

	tests := []struct {
		name            string
		path            *repository.RelationshipPath
		expectedLabel   string
		expectedKinship string
	}{
		{
			name: "Great-Grandfather",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Person", Gender: "male"},
					{ID: "2", Name: "Parent", Gender: "male"},
					{ID: "3", Name: "Grandparent", Gender: "male"},
					{ID: "4", Name: "Great-Grandparent", Gender: "male"},
				},
				Relationships: []string{
					models.RelationshipTypeParentOf,
					models.RelationshipTypeParentOf,
					models.RelationshipTypeParentOf,
				},
			},
			expectedLabel:   "Great-Grandfather",
			expectedKinship: "Tata",
		},
		{
			name: "First Cousin",
			path: &repository.RelationshipPath{
				Nodes: []repository.PathNode{
					{ID: "1", Name: "Person", Gender: "male"},
					{ID: "2", Name: "Parent", Gender: "male"},
					{ID: "3", Name: "Uncle", Gender: "male"},
					{ID: "4", Name: "Cousin", Gender: "male"},
				},
				Relationships: []string{
					models.RelationshipTypeParentOf,
					models.RelationshipTypeSiblingOf,
					models.RelationshipTypeParentOf,
				},
			},
			expectedLabel:   "Cousin",
			expectedKinship: "Bava/Maradalu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			label, kinship := service.mapMultiHopRelationship(tt.path)
			if label != tt.expectedLabel {
				t.Errorf("mapMultiHopRelationship() label = %v, want %v", label, tt.expectedLabel)
			}
			if kinship != tt.expectedKinship {
				t.Errorf("mapMultiHopRelationship() kinship = %v, want %v", kinship, tt.expectedKinship)
			}
		})
	}
}
