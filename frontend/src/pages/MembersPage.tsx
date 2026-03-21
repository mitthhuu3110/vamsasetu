import React, { useState } from 'react';
import MemberList from '../components/family/MemberList';
import AddMemberModal from '../components/family/AddMemberModal';
import Button from '../components/ui/Button';
import type { Member } from '../types/member';

const MembersPage: React.FC = () => {
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);

  const handleMemberClick = (member: Member) => {
    // TODO: Open member details panel or modal
    console.log('Selected member:', member);
  };

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-3xl font-bold text-charcoal">
            Family Members
          </h1>
          <p className="text-charcoal/70 mt-1">
            Manage your family tree members
          </p>
        </div>
        <Button
          variant="primary"
          onClick={() => setIsAddModalOpen(true)}
        >
          + Add Member
        </Button>
      </div>

      {/* Member List */}
      <MemberList onMemberClick={handleMemberClick} />

      {/* Add Member Modal */}
      <AddMemberModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
      />
    </div>
  );
};

export default MembersPage;
