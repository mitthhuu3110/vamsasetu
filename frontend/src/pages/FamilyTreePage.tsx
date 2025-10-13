import React, { useState } from 'react';
import { FamilyMember, Relationship, Gender, RelationshipType } from '../types/index.ts';
import FamilyTreeVisualizer from '../components/FamilyTree/FamilyTreeVisualizer.tsx';
import MemberDetails from '../components/FamilyTree/MemberDetails.tsx';
import AddMemberModal from '../components/FamilyTree/AddMemberModal.tsx';
import { PlusIcon, MagnifyingGlassIcon } from '@heroicons/react/24/outline';

const FamilyTreePage: React.FC = () => {
  const [selectedMember, setSelectedMember] = useState<FamilyMember | null>(null);
  const [showAddMember, setShowAddMember] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  // Mock data fallback (used instead of API for now)
  const members: FamilyMember[] = [
    {
      id: '1',
      firstName: 'Ravi',
      lastName: 'Kumar',
      gender: Gender.MALE,
      isAlive: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      dateOfBirth: '1990-05-10',
      email: 'ravi@example.com',
    },
    {
      id: '2',
      firstName: 'Anita',
      lastName: 'Kumar',
      gender: Gender.FEMALE,
      isAlive: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      dateOfBirth: '1992-08-12',
      email: 'anita@example.com',
    },
    {
      id: '3',
      firstName: 'Aarav',
      lastName: 'Kumar',
      gender: Gender.MALE,
      isAlive: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      dateOfBirth: '2018-03-02',
    },
  ];

  const familyData = {
    nodes: members,
    edges: [
      {
        id: 'e1-3',
        fromMemberId: '1',
        toMemberId: '3',
        relationshipType: RelationshipType.PARENT,
        isActive: true,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      {
        id: 'e2-3',
        fromMemberId: '2',
        toMemberId: '3',
        relationshipType: RelationshipType.PARENT,
        isActive: true,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      {
        id: 'e1-2',
        fromMemberId: '1',
        toMemberId: '2',
        relationshipType: RelationshipType.SPOUSE,
        isActive: true,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    ] as Relationship[],
  };

  const filteredMembers = members.filter(member =>
    member.firstName.toLowerCase().includes(searchQuery.toLowerCase()) ||
    member.lastName.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Family Tree</h1>
          <p className="text-gray-600">Visualize and manage your family relationships</p>
        </div>
        <button
          onClick={() => setShowAddMember(true)}
          className="btn-primary flex items-center"
        >
          <PlusIcon className="w-5 h-5 mr-2" />
          Add Member
        </button>
      </div>

      {/* Search and Filters */}
      <div className="card">
        <div className="flex items-center space-x-4">
          <div className="flex-1 relative">
            <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search family members..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input-field pl-10"
            />
          </div>
          <div className="flex space-x-2">
            <button className="btn-secondary">All</button>
            <button className="btn-secondary">Alive</button>
            <button className="btn-secondary">Deceased</button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Family Tree Visualizer */}
        <div className="lg:col-span-3">
          <div className="card h-96">
            <FamilyTreeVisualizer
              members={familyData?.nodes || []}
              relationships={familyData?.edges || []}
              onMemberSelect={setSelectedMember}
              selectedMember={selectedMember}
            />
          </div>
        </div>

        {/* Member Details Sidebar */}
        <div className="lg:col-span-1">
          {selectedMember ? (
            <MemberDetails
              member={selectedMember}
              onEdit={() => {/* TODO: Implement edit */}}
              onClose={() => setSelectedMember(null)}
            />
          ) : (
            <div className="card">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Family Members</h3>
              <div className="space-y-2 max-h-80 overflow-y-auto">
                {filteredMembers.map((member) => (
                  <div
                    key={member.id}
                    onClick={() => setSelectedMember(member)}
                    className="p-3 hover:bg-gray-50 rounded-lg cursor-pointer transition-colors"
                  >
                    <div className="flex items-center space-x-3">
                      <div className="w-10 h-10 bg-primary-100 rounded-full flex items-center justify-center">
                        <span className="text-primary-600 font-medium">
                          {member.firstName[0]}{member.lastName[0]}
                        </span>
                      </div>
                      <div>
                        <p className="font-medium text-gray-900">
                          {member.firstName} {member.lastName}
                        </p>
                        <p className="text-sm text-gray-500">
                          {member.isAlive ? 'Alive' : 'Deceased'}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Add Member Modal */}
      {showAddMember && (
        <AddMemberModal
          onClose={() => setShowAddMember(false)}
          onSuccess={() => {
            setShowAddMember(false);
            // Refetch data
            window.location.reload();
          }}
        />
      )}
    </div>
  );
};

export default FamilyTreePage;
