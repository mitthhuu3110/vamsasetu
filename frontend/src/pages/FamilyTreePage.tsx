import React, { useState } from 'react';
import { ReactFlowProvider } from 'reactflow';
import TreeCanvas from '../components/family/TreeCanvas';
import TreeControls from '../components/family/TreeControls';
import MemberDetailsPanel from '../components/family/MemberDetailsPanel';
import AddMemberModal from '../components/family/AddMemberModal';
import AddRelationshipModal from '../components/family/AddRelationshipModal';

/**
 * FamilyTreePage Component
 * 
 * Main page for interactive family tree visualization.
 * Composes TreeCanvas, TreeControls, MemberDetailsPanel, AddMemberModal, and AddRelationshipModal.
 * 
 * Features:
 * - Interactive React Flow tree visualization
 * - Zoom, pan, and fit-view controls
 * - Member details side panel
 * - Add member and relationship modals
 * - Relationship path highlighting with animated traveling dot
 * 
 * **Validates: Requirements 3.1, 4.5, 10.7**
 */
const FamilyTreePage: React.FC = () => {
  const [selectedMemberId, setSelectedMemberId] = useState<string | null>(null);
  const [showAddMemberModal, setShowAddMemberModal] = useState(false);
  const [showAddRelationshipModal, setShowAddRelationshipModal] = useState(false);

  return (
    <div className="h-screen flex flex-col">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-6 py-4">
        <h1 className="font-display text-2xl font-bold text-charcoal">
          Family Tree
        </h1>
        <p className="text-charcoal/70 text-sm mt-1">
          Visualize and explore your family relationships
        </p>
      </div>

      {/* Tree Container */}
      <div className="flex-1 relative">
        <ReactFlowProvider>
          <TreeCanvas onNodeClick={setSelectedMemberId} />
          
          {/* Tree Controls - Positioned absolutely */}
          <div className="absolute top-4 left-4 z-10">
            <TreeControls
              onAddMember={() => setShowAddMemberModal(true)}
              onAddRelationship={() => setShowAddRelationshipModal(true)}
            />
          </div>
        </ReactFlowProvider>

        {/* Member Details Panel */}
        <MemberDetailsPanel
          memberId={selectedMemberId}
          onClose={() => setSelectedMemberId(null)}
        />

        {/* Add Member Modal */}
        <AddMemberModal
          isOpen={showAddMemberModal}
          onClose={() => setShowAddMemberModal(false)}
        />

        {/* Add Relationship Modal */}
        <AddRelationshipModal
          isOpen={showAddRelationshipModal}
          onClose={() => setShowAddRelationshipModal(false)}
        />
      </div>
    </div>
  );
};

export default FamilyTreePage;
