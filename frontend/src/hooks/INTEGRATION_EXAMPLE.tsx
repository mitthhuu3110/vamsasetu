/**
 * Integration Example: Complete Component Using Multiple Hooks
 * 
 * This example demonstrates how to use multiple React Query hooks
 * together in a real component scenario.
 */

import { useState } from 'react';
import {
  useMembers,
  useCreateMember,
  useDeleteMember,
  useRelationships,
  useFamilyTree,
  useUpcomingEvents,
} from './index';

/**
 * Example: Family Dashboard Component
 * 
 * This component demonstrates:
 * - Fetching multiple data sources
 * - Creating and deleting members
 * - Managing relationships
 * - Displaying family tree
 * - Showing upcoming events
 */
export function FamilyDashboard() {
  const [_selectedMemberId, setSelectedMemberId] = useState<string | null>(null);

  // Fetch data using query hooks
  const { data: membersData, isLoading: membersLoading } = useMembers();
  const { data: relationshipsData } = useRelationships();
  const { data: treeData } = useFamilyTree();
  const { data: upcomingEventsData } = useUpcomingEvents(7);

  // Mutation hooks
  const createMember = useCreateMember();
  const deleteMember = useDeleteMember();
  // const _createRelationship = useCreateRelationship(); // Example - not used in this component

  // Handle member creation
  const handleAddMember = async () => {
    const result = await createMember.mutateAsync({
      name: 'New Member',
      dateOfBirth: '2000-01-01',
      gender: 'male',
      email: 'new@example.com',
    });

    if (result.success) {
      console.log('Member created:', result.data);
      // Cache is automatically invalidated, UI will update
    } else {
      console.error('Failed to create member:', result.error);
    }
  };

  // Handle member deletion
  const handleDeleteMember = (memberId: string) => {
    if (confirm('Are you sure you want to delete this member?')) {
      deleteMember.mutate(memberId, {
        onSuccess: () => {
          console.log('Member deleted successfully');
          setSelectedMemberId(null);
        },
        onError: (error) => {
          console.error('Failed to delete member:', error);
        },
      });
    }
  };

  // Handle relationship creation (example - not used in this component)
  // const handleAddRelationship = async (fromId: string, toId: string) => {
  //   await _createRelationship.mutateAsync({
  //     type: 'PARENT_OF',
  //     fromId,
  //     toId,
  //   });
  // };

  // Loading state
  if (membersLoading) {
    return <div className="p-4">Loading family data...</div>;
  }

  const members = membersData?.data?.members || [];
  const relationships = relationshipsData?.data || [];
  const upcomingEvents = upcomingEventsData?.data || [];

  return (
    <div className="p-4 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Family Dashboard</h1>
        <button
          onClick={handleAddMember}
          disabled={createMember.isPending}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
        >
          {createMember.isPending ? 'Adding...' : 'Add Member'}
        </button>
      </div>

      {/* Error Display */}
      {createMember.isError && (
        <div className="p-4 bg-red-100 text-red-700 rounded">
          Error: {createMember.error.message}
        </div>
      )}

      {/* Stats */}
      <div className="grid grid-cols-3 gap-4">
        <div className="p-4 bg-white rounded shadow">
          <h3 className="text-lg font-semibold">Total Members</h3>
          <p className="text-3xl">{members.length}</p>
        </div>
        <div className="p-4 bg-white rounded shadow">
          <h3 className="text-lg font-semibold">Relationships</h3>
          <p className="text-3xl">{relationships.length}</p>
        </div>
        <div className="p-4 bg-white rounded shadow">
          <h3 className="text-lg font-semibold">Upcoming Events</h3>
          <p className="text-3xl">{upcomingEvents.length}</p>
        </div>
      </div>

      {/* Members List */}
      <div className="bg-white rounded shadow">
        <h2 className="text-xl font-semibold p-4 border-b">Family Members</h2>
        <div className="divide-y">
          {members.map((member) => (
            <div
              key={member.id}
              className="p-4 flex justify-between items-center hover:bg-gray-50"
            >
              <div>
                <h3 className="font-semibold">{member.name}</h3>
                <p className="text-sm text-gray-600">
                  {member.gender} • Born: {member.dateOfBirth}
                </p>
              </div>
              <div className="space-x-2">
                <button
                  onClick={() => setSelectedMemberId(member.id)}
                  className="px-3 py-1 bg-blue-500 text-white rounded text-sm"
                >
                  View
                </button>
                <button
                  onClick={() => handleDeleteMember(member.id)}
                  disabled={deleteMember.isPending}
                  className="px-3 py-1 bg-red-500 text-white rounded text-sm disabled:opacity-50"
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Upcoming Events */}
      {upcomingEvents.length > 0 && (
        <div className="bg-white rounded shadow">
          <h2 className="text-xl font-semibold p-4 border-b">Upcoming Events</h2>
          <div className="divide-y">
            {upcomingEvents.map((event) => (
              <div key={event.id} className="p-4">
                <h3 className="font-semibold">{event.title}</h3>
                <p className="text-sm text-gray-600">
                  {event.eventType} • {event.eventDate}
                </p>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Family Tree Visualization */}
      {treeData?.success && (
        <div className="bg-white rounded shadow">
          <h2 className="text-xl font-semibold p-4 border-b">Family Tree</h2>
          <div className="p-4">
            <p className="text-gray-600">
              {treeData.data?.nodes.length} members in tree
            </p>
            {/* React Flow component would go here */}
          </div>
        </div>
      )}
    </div>
  );
}

/**
 * Example: Member Detail Component with Relationships
 */
export function MemberDetail({ memberId }: { memberId: string }) {
  const { data: memberData, isLoading } = useMembers();
  const { data: relationshipsData } = useRelationships();

  if (isLoading) return <div>Loading...</div>;

  const member = memberData?.data?.members.find((m) => m.id === memberId);
  const memberRelationships = relationshipsData?.data?.filter(
    (r) => r.fromId === memberId || r.toId === memberId
  );

  if (!member) return <div>Member not found</div>;

  return (
    <div className="p-4 space-y-4">
      <h2 className="text-2xl font-bold">{member.name}</h2>
      
      <div className="space-y-2">
        <p><strong>Date of Birth:</strong> {member.dateOfBirth}</p>
        <p><strong>Gender:</strong> {member.gender}</p>
        <p><strong>Email:</strong> {member.email}</p>
        <p><strong>Phone:</strong> {member.phone}</p>
      </div>

      <div>
        <h3 className="text-xl font-semibold mb-2">Relationships</h3>
        {memberRelationships?.map((rel, idx) => (
          <div key={idx} className="p-2 bg-gray-100 rounded mb-2">
            {rel.type}: {rel.fromId} → {rel.toId}
          </div>
        ))}
      </div>
    </div>
  );
}

/**
 * Example: Authentication Flow
 */
export function AuthExample() {
  const { useLogin, useRegister, useProfile } = require('./index');
  
  const login = useLogin();
  const register = useRegister();
  const { data: profileData } = useProfile();

  const handleLogin = async () => {
    const result = await login.mutateAsync({
      email: 'user@example.com',
      password: 'password123',
    });

    if (result.success) {
      console.log('Logged in as:', result.data?.user.name);
      // User is automatically stored in authStore
    }
  };

  const handleRegister = async () => {
    const result = await register.mutateAsync({
      email: 'newuser@example.com',
      password: 'password123',
      name: 'New User',
    });

    if (result.success) {
      console.log('Registered and logged in as:', result.data?.user.name);
    }
  };

  return (
    <div className="p-4 space-y-4">
      <button
        onClick={handleLogin}
        disabled={login.isPending}
        className="px-4 py-2 bg-blue-500 text-white rounded"
      >
        {login.isPending ? 'Logging in...' : 'Login'}
      </button>

      <button
        onClick={handleRegister}
        disabled={register.isPending}
        className="px-4 py-2 bg-green-500 text-white rounded"
      >
        {register.isPending ? 'Registering...' : 'Register'}
      </button>

      {profileData?.success && (
        <div className="p-4 bg-gray-100 rounded">
          <h3 className="font-semibold">Current User</h3>
          <p>Name: {profileData.data?.name}</p>
          <p>Email: {profileData.data?.email}</p>
          <p>Role: {profileData.data?.role}</p>
        </div>
      )}
    </div>
  );
}

/**
 * Key Takeaways from These Examples:
 * 
 * 1. Multiple hooks can be used together seamlessly
 * 2. Mutations automatically invalidate related queries
 * 3. Loading and error states are easily accessible
 * 4. Auth state is automatically managed via Zustand
 * 5. Type safety is maintained throughout
 * 6. Cache invalidation ensures UI stays in sync
 */
