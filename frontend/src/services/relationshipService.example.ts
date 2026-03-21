/**
 * Example usage of the RelationshipService
 * This file demonstrates how to use the relationship service methods
 */

import relationshipService from './relationshipService';

// Example 1: Get all relationships
async function getAllRelationshipsExample() {
  const response = await relationshipService.getAll();

  if (response.success && response.data) {
    console.log('Relationships retrieved successfully!');
    console.log('Total relationships:', response.data.length);
    console.log('Relationships:', response.data);
  } else {
    console.error('Failed to get relationships:', response.error);
  }
}

// Example 2: Create a new relationship
async function createRelationshipExample() {
  const response = await relationshipService.create({
    type: 'PARENT_OF',
    fromId: 'member-123',
    toId: 'member-456',
  });

  if (response.success && response.data) {
    console.log('Relationship created successfully!');
    console.log('Relationship:', response.data);
    console.log('Type:', response.data.type);
    console.log('From:', response.data.fromId);
    console.log('To:', response.data.toId);
  } else {
    console.error('Failed to create relationship:', response.error);
  }
}

// Example 3: Delete a relationship
async function deleteRelationshipExample() {
  const response = await relationshipService.delete(
    'relationship-789',
    'member-123',
    'member-456',
    'PARENT_OF'
  );

  if (response.success && response.data) {
    console.log('Relationship deleted successfully!');
    console.log('Message:', response.data.message);
  } else {
    console.error('Failed to delete relationship:', response.error);
  }
}

// Example 4: Find relationship path between two members
async function findPathExample() {
  const response = await relationshipService.findPath('member-123', 'member-789');

  if (response.success && response.data) {
    console.log('Relationship path found!');
    console.log('Path:', response.data.path);
    console.log('Relation Label:', response.data.relationLabel);
    console.log('Kinship Term:', response.data.kinshipTerm);
    console.log('Description:', response.data.description);
  } else {
    console.error('Failed to find path:', response.error);
  }
}

// Example 5: Create multiple relationship types
async function createMultipleRelationshipsExample() {
  // Create a spouse relationship
  const spouseResponse = await relationshipService.create({
    type: 'SPOUSE_OF',
    fromId: 'member-100',
    toId: 'member-101',
  });
  console.log('Spouse relationship:', spouseResponse.data);

  // Create a parent-child relationship
  const parentResponse = await relationshipService.create({
    type: 'PARENT_OF',
    fromId: 'member-100',
    toId: 'member-102',
  });
  console.log('Parent-child relationship:', parentResponse.data);

  // Create a sibling relationship
  const siblingResponse = await relationshipService.create({
    type: 'SIBLING_OF',
    fromId: 'member-102',
    toId: 'member-103',
  });
  console.log('Sibling relationship:', siblingResponse.data);
}

// Example 6: Complete relationship management flow
async function completeRelationshipFlowExample() {
  // 1. Get all existing relationships
  const getAllResponse = await relationshipService.getAll();
  console.log('Existing relationships:', getAllResponse.data?.length);

  // 2. Create a new relationship
  const createResponse = await relationshipService.create({
    type: 'PARENT_OF',
    fromId: 'parent-id',
    toId: 'child-id',
  });

  if (!createResponse.success) {
    console.error('Failed to create relationship:', createResponse.error);
    return;
  }

  console.log('Created relationship:', createResponse.data);

  // 3. Find the path between two members
  const pathResponse = await relationshipService.findPath('parent-id', 'grandchild-id');
  
  if (pathResponse.success && pathResponse.data) {
    console.log('Relationship path:');
    pathResponse.data.path.forEach((node, index) => {
      console.log(`  ${index + 1}. ${node.name} (${node.id})`);
    });
    console.log('Relation:', pathResponse.data.relationLabel);
    console.log('Telugu term:', pathResponse.data.kinshipTerm);
  }

  // 4. Delete the relationship
  const deleteResponse = await relationshipService.delete(
    'relationship-id',
    'parent-id',
    'child-id',
    'PARENT_OF'
  );
  console.log('Delete result:', deleteResponse.data?.message);
}

// Example 7: Find complex family relationships
async function findComplexRelationshipsExample() {
  // Find relationship between grandparent and grandchild
  const grandparentPath = await relationshipService.findPath('grandparent-id', 'grandchild-id');
  console.log('Grandparent relationship:', grandparentPath.data?.relationLabel);

  // Find relationship between uncle and nephew
  const unclePath = await relationshipService.findPath('uncle-id', 'nephew-id');
  console.log('Uncle relationship:', unclePath.data?.relationLabel);

  // Find relationship between cousins
  const cousinPath = await relationshipService.findPath('cousin1-id', 'cousin2-id');
  console.log('Cousin relationship:', cousinPath.data?.relationLabel);
}

export {
  getAllRelationshipsExample,
  createRelationshipExample,
  deleteRelationshipExample,
  findPathExample,
  createMultipleRelationshipsExample,
  completeRelationshipFlowExample,
  findComplexRelationshipsExample,
};
