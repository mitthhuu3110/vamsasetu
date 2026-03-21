export interface Relationship {
  type: 'SPOUSE_OF' | 'PARENT_OF' | 'SIBLING_OF';
  fromId: string;
  toId: string;
  createdAt: string;
}

export interface CreateRelationshipData {
  type: 'SPOUSE_OF' | 'PARENT_OF' | 'SIBLING_OF';
  fromId: string;
  toId: string;
}

// Alias for task requirement
export type CreateRelationshipRequest = CreateRelationshipData;

export interface RelationshipPath {
  path: Array<{ id: string; name: string }>;
  relationLabel: string;
  kinshipTerm: string;
  description: string;
}
