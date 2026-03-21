/**
 * Example usage of the MemberService
 * This file demonstrates how to use the member service methods
 */

import memberService from './memberService';

// Example 1: Get all members with pagination
async function getAllMembersExample() {
  const response = await memberService.getAll({
    page: 1,
    limit: 10,
  });

  if (response.success && response.data) {
    console.log('Members retrieved successfully!');
    console.log('Total members:', response.data.total);
    console.log('Members:', response.data.members);
  } else {
    console.error('Failed to get members:', response.error);
  }
}

// Example 2: Search members by name
async function searchMembersExample() {
  const response = await memberService.search('John');

  if (response.success && response.data) {
    console.log('Search results:', response.data.members);
  } else {
    console.error('Search failed:', response.error);
  }
}

// Example 3: Get a specific member by ID
async function getMemberByIdExample() {
  const response = await memberService.getById('member-123');

  if (response.success && response.data) {
    console.log('Member details:', response.data);
    console.log('Name:', response.data.name);
    console.log('Email:', response.data.email);
  } else {
    console.error('Failed to get member:', response.error);
  }
}

// Example 4: Create a new member
async function createMemberExample() {
  const response = await memberService.create({
    name: 'John Doe',
    dateOfBirth: '1990-01-15T00:00:00Z',
    gender: 'male',
    email: 'john.doe@example.com',
    phone: '+91-9876543210',
    avatarUrl: 'https://example.com/avatar.jpg',
  });

  if (response.success && response.data) {
    console.log('Member created successfully!');
    console.log('New member ID:', response.data.id);
    console.log('Member:', response.data);
  } else {
    console.error('Failed to create member:', response.error);
  }
}

// Example 5: Update an existing member
async function updateMemberExample() {
  const response = await memberService.update('member-123', {
    name: 'John Smith',
    email: 'john.smith@example.com',
  });

  if (response.success && response.data) {
    console.log('Member updated successfully!');
    console.log('Updated member:', response.data);
  } else {
    console.error('Failed to update member:', response.error);
  }
}

// Example 6: Delete a member
async function deleteMemberExample() {
  const response = await memberService.delete('member-123');

  if (response.success && response.data) {
    console.log('Member deleted successfully!');
    console.log('Message:', response.data.message);
  } else {
    console.error('Failed to delete member:', response.error);
  }
}

// Example 7: Filter members by gender
async function filterByGenderExample() {
  const response = await memberService.getAll({
    page: 1,
    limit: 20,
    gender: 'male',
  });

  if (response.success && response.data) {
    console.log('Male members:', response.data.members);
  } else {
    console.error('Failed to filter members:', response.error);
  }
}

// Example 8: Complete member management flow
async function completeMemberFlowExample() {
  // 1. Create a new member
  const createResponse = await memberService.create({
    name: 'Jane Doe',
    dateOfBirth: '1995-05-20T00:00:00Z',
    gender: 'female',
    email: 'jane.doe@example.com',
    phone: '+91-9876543211',
  });

  if (!createResponse.success) {
    console.error('Failed to create member:', createResponse.error);
    return;
  }

  const memberId = createResponse.data!.id;
  console.log('Created member with ID:', memberId);

  // 2. Get the member details
  const getResponse = await memberService.getById(memberId);
  console.log('Member details:', getResponse.data);

  // 3. Update the member
  const updateResponse = await memberService.update(memberId, {
    phone: '+91-9876543299',
  });
  console.log('Updated member:', updateResponse.data);

  // 4. Search for the member
  const searchResponse = await memberService.search('Jane');
  console.log('Search results:', searchResponse.data?.members);

  // 5. Delete the member
  const deleteResponse = await memberService.delete(memberId);
  console.log('Delete result:', deleteResponse.data?.message);
}

export {
  getAllMembersExample,
  searchMembersExample,
  getMemberByIdExample,
  createMemberExample,
  updateMemberExample,
  deleteMemberExample,
  filterByGenderExample,
  completeMemberFlowExample,
};
