import React, { useState } from 'react';
import { FamilyMember, Relationship, Gender, RelationshipType } from '../types/index.ts';
import AdvancedFamilyTreeVisualizer from '../components/FamilyTree/AdvancedFamilyTreeVisualizer.tsx';
import MemberDetails from '../components/FamilyTree/MemberDetails.tsx';
import AddMemberModal from '../components/FamilyTree/AddMemberModal.tsx';
import { PlusIcon, MagnifyingGlassIcon, SparklesIcon, HeartIcon } from '@heroicons/react/24/outline';

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
    <div className="space-y-6 min-h-screen bg-gradient-to-br from-warm-beige to-cream dark:from-dark-bg dark:to-dark-card">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-display font-bold text-gradient">
            🌳 Family Tree
          </h1>
          <p className="text-warm-brown dark:text-dark-text mt-2">
            Visualize and manage your family relationships with love
          </p>
        </div>
        <button
          onClick={() => setShowAddMember(true)}
          className="btn-primary flex items-center space-x-2 shadow-lg hover:shadow-xl transition-all duration-300"
        >
          <PlusIcon className="w-5 h-5" />
          <span>Add Member</span>
        </button>
      </div>

      {/* Search and Filters */}
      <div className="card card-hover">
        <div className="flex items-center space-x-4">
          <div className="flex-1 relative">
            <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-warm-brown dark:text-dark-text" />
            <input
              type="text"
              placeholder="Search family members..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-3 bg-warm-beige dark:bg-dark-accent border border-gray-200 dark:border-dark-accent rounded-lg text-warm-brown dark:text-dark-text placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-soft-gold focus:border-transparent transition-all duration-200"
            />
          </div>
          <div className="flex space-x-2">
            <button className="btn-secondary text-sm">All</button>
            <button className="btn-secondary text-sm">Alive</button>
            <button className="btn-secondary text-sm">Deceased</button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Family Tree Visualizer */}
        <div className="lg:col-span-3">
          <div className="card h-[600px] overflow-hidden">
            <AdvancedFamilyTreeVisualizer
              familyMembers={familyData?.nodes || []}
              relationships={familyData?.edges || []}
              onMemberClick={setSelectedMember}
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
            <div className="card card-hover">
              <div className="flex items-center space-x-2 mb-4">
                <HeartIcon className="w-5 h-5 text-soft-gold" />
                <h3 className="text-lg font-display font-bold text-warm-brown dark:text-dark-text">
                  Family Members
                </h3>
              </div>
              <div className="space-y-3 max-h-80 overflow-y-auto scrollbar-hide">
                {filteredMembers.map((member) => (
                  <div
                    key={member.id}
                    onClick={() => setSelectedMember(member)}
                    className="p-3 hover:bg-warm-beige dark:hover:bg-dark-accent rounded-lg cursor-pointer transition-all duration-200 hover:shadow-md border border-transparent hover:border-soft-gold"
                  >
                    <div className="flex items-center space-x-3">
                      <div className="w-10 h-10 bg-gradient-to-br from-soft-gold to-deep-gold rounded-full flex items-center justify-center shadow-md">
                        <span className="text-white font-bold text-sm">
                          {member.firstName[0]}{member.lastName[0]}
                        </span>
                      </div>
                      <div>
                        <p className="font-medium text-warm-brown dark:text-dark-text">
                          {member.firstName} {member.lastName}
                        </p>
                        <p className="text-sm text-gray-600 dark:text-gray-400 flex items-center space-x-1">
                          <span className={`w-2 h-2 rounded-full ${member.isAlive ? 'bg-soft-green' : 'bg-gray-400'}`}></span>
                          <span>{member.isAlive ? 'Alive' : 'Deceased'}</span>
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
