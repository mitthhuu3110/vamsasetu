import React, { useState } from 'react';
import { useMembers, useSearchMembers } from '../../hooks/useMembers';
import Card, { CardContent } from '../ui/Card';
import Input from '../ui/Input';
import LoadingSpinner from '../common/LoadingSpinner';
import EmptyState from '../common/EmptyState';
import type { Member } from '../../types/member';

interface MemberListProps {
  onMemberClick?: (member: Member) => void;
}

const MemberList: React.FC<MemberListProps> = ({ onMemberClick }) => {
  const [searchQuery, setSearchQuery] = useState('');
  const { data: membersResponse, isLoading: isLoadingAll } = useMembers();
  const { data: searchResponse, isLoading: isSearching } = useSearchMembers(searchQuery);

  const isLoading = searchQuery ? isSearching : isLoadingAll;
  const membersData = searchQuery
    ? searchResponse?.data
    : membersResponse?.data;
  
  const members = membersData?.members || [];

  return (
    <div className="space-y-4">
      {/* Search Bar */}
      <Input
        type="text"
        placeholder="Search members by name..."
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        fullWidth
      />

      {/* Members Grid */}
      {isLoading ? (
        <LoadingSpinner />
      ) : members.length > 0 ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {members.map((member) => (
            <Card
              key={member.id}
              variant="elevated"
              hoverable
              onClick={() => onMemberClick?.(member)}
              className="cursor-pointer"
            >
              <CardContent className="flex flex-col items-center text-center p-4">
                {/* Avatar */}
                <div className="w-20 h-20 rounded-full bg-gradient-to-br from-saffron to-turmeric flex items-center justify-center text-white text-2xl font-bold mb-3">
                  {member.avatarUrl ? (
                    <img
                      src={member.avatarUrl}
                      alt={member.name}
                      className="w-full h-full rounded-full object-cover"
                    />
                  ) : (
                    member.name.charAt(0).toUpperCase()
                  )}
                </div>

                {/* Name */}
                <h3 className="font-semibold text-charcoal mb-1">{member.name}</h3>

                {/* Gender Badge */}
                <span
                  className={`text-xs px-2 py-1 rounded-full ${
                    member.gender === 'male'
                      ? 'bg-blue-100 text-blue-700'
                      : member.gender === 'female'
                      ? 'bg-pink-100 text-pink-700'
                      : 'bg-gray-100 text-gray-700'
                  }`}
                >
                  {member.gender}
                </span>

                {/* Date of Birth */}
                {member.dateOfBirth && (
                  <p className="text-sm text-charcoal/60 mt-2">
                    {new Date(member.dateOfBirth).toLocaleDateString()}
                  </p>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <EmptyState
          icon="👥"
          title={searchQuery ? 'No members found' : 'No members yet'}
          description={
            searchQuery
              ? `No members match "${searchQuery}"`
              : 'Add family members to start building your tree'
          }
        />
      )}
    </div>
  );
};

export default MemberList;
