import { useState } from 'react';
import { ReactFlowProvider } from 'reactflow';
import TreeControls from './TreeControls';
import TreeCanvas from './TreeCanvas';

/**
 * TreeControls Component Examples
 * 
 * This file demonstrates various usage patterns for the TreeControls component.
 */

// Example 1: Basic TreeControls with zoom and fit view only
export function BasicTreeControls() {
  return (
    <ReactFlowProvider>
      <div className="relative w-full h-screen bg-ivory">
        <TreeCanvas />
        
        {/* Basic controls without add buttons */}
        <div className="absolute top-4 right-4 z-10">
          <TreeControls />
        </div>
      </div>
    </ReactFlowProvider>
  );
}

// Example 2: TreeControls with Add Member and Add Relationship buttons
export function TreeControlsWithActions() {
  const [showAddMember, setShowAddMember] = useState(false);
  const [showAddRelationship, setShowAddRelationship] = useState(false);

  const handleAddMember = () => {
    console.log('Add member clicked');
    setShowAddMember(true);
  };

  const handleAddRelationship = () => {
    console.log('Add relationship clicked');
    setShowAddRelationship(true);
  };

  return (
    <ReactFlowProvider>
      <div className="relative w-full h-screen bg-ivory">
        <TreeCanvas />
        
        {/* Controls with all action buttons */}
        <div className="absolute top-4 right-4 z-10">
          <TreeControls
            onAddMember={handleAddMember}
            onAddRelationship={handleAddRelationship}
          />
        </div>

        {/* Mock modals */}
        {showAddMember && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white p-6 rounded-lg shadow-xl">
              <h2 className="text-xl font-bold mb-4">Add Member Modal</h2>
              <button
                onClick={() => setShowAddMember(false)}
                className="px-4 py-2 bg-saffron text-white rounded"
              >
                Close
              </button>
            </div>
          </div>
        )}

        {showAddRelationship && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white p-6 rounded-lg shadow-xl">
              <h2 className="text-xl font-bold mb-4">Add Relationship Modal</h2>
              <button
                onClick={() => setShowAddRelationship(false)}
                className="px-4 py-2 bg-teal text-white rounded"
              >
                Close
              </button>
            </div>
          </div>
        )}
      </div>
    </ReactFlowProvider>
  );
}

// Example 3: Mobile-optimized positioning
export function MobileTreeControls() {
  return (
    <ReactFlowProvider>
      <div className="relative w-full h-screen bg-ivory">
        <TreeCanvas />
        
        {/* Bottom-right on mobile, top-right on desktop */}
        <div className="absolute bottom-4 right-4 md:top-4 md:bottom-auto z-10">
          <TreeControls
            onAddMember={() => console.log('Add member')}
            onAddRelationship={() => console.log('Add relationship')}
          />
        </div>
      </div>
    </ReactFlowProvider>
  );
}

// Example 4: TreeControls with only Add Member button
export function TreeControlsWithAddMemberOnly() {
  const handleAddMember = () => {
    alert('Opening Add Member modal...');
  };

  return (
    <ReactFlowProvider>
      <div className="relative w-full h-screen bg-ivory">
        <TreeCanvas />
        
        {/* Only show add member button */}
        <div className="absolute top-4 right-4 z-10">
          <TreeControls onAddMember={handleAddMember} />
        </div>
      </div>
    </ReactFlowProvider>
  );
}

// Example 5: TreeControls with custom className
export function TreeControlsWithCustomStyling() {
  return (
    <ReactFlowProvider>
      <div className="relative w-full h-screen bg-ivory">
        <TreeCanvas />
        
        {/* Custom styling with additional classes */}
        <div className="absolute top-4 right-4 z-10">
          <TreeControls
            className="shadow-2xl border-2 border-saffron"
            onAddMember={() => console.log('Add member')}
            onAddRelationship={() => console.log('Add relationship')}
          />
        </div>
      </div>
    </ReactFlowProvider>
  );
}

// Example 6: Multiple TreeControls instances (left and right)
export function DualTreeControls() {
  return (
    <ReactFlowProvider>
      <div className="relative w-full h-screen bg-ivory">
        <TreeCanvas />
        
        {/* Zoom controls on the left */}
        <div className="absolute top-4 left-4 z-10">
          <TreeControls />
        </div>

        {/* Action controls on the right */}
        <div className="absolute top-4 right-4 z-10">
          <div className="flex flex-col gap-2 bg-white rounded-lg shadow-lg border border-gray-200 p-2">
            <button
              onClick={() => console.log('Add member')}
              className="w-10 h-10 bg-saffron text-white rounded flex items-center justify-center hover:bg-saffron-light"
              aria-label="Add member"
            >
              👤
            </button>
            <button
              onClick={() => console.log('Add relationship')}
              className="w-10 h-10 bg-teal text-white rounded flex items-center justify-center hover:bg-teal-light"
              aria-label="Add relationship"
            >
              🔗
            </button>
          </div>
        </div>
      </div>
    </ReactFlowProvider>
  );
}

// Example 7: TreeControls with state management
export function TreeControlsWithState() {
  const [memberCount, setMemberCount] = useState(0);
  const [relationshipCount, setRelationshipCount] = useState(0);

  const handleAddMember = () => {
    setMemberCount((prev) => prev + 1);
    console.log(`Total members: ${memberCount + 1}`);
  };

  const handleAddRelationship = () => {
    setRelationshipCount((prev) => prev + 1);
    console.log(`Total relationships: ${relationshipCount + 1}`);
  };

  return (
    <ReactFlowProvider>
      <div className="relative w-full h-screen bg-ivory">
        <TreeCanvas />
        
        {/* Stats panel */}
        <div className="absolute top-4 left-4 z-10 bg-white p-4 rounded-lg shadow-lg">
          <h3 className="font-bold mb-2">Tree Stats</h3>
          <p className="text-sm">Members: {memberCount}</p>
          <p className="text-sm">Relationships: {relationshipCount}</p>
        </div>

        {/* Controls */}
        <div className="absolute top-4 right-4 z-10">
          <TreeControls
            onAddMember={handleAddMember}
            onAddRelationship={handleAddRelationship}
          />
        </div>
      </div>
    </ReactFlowProvider>
  );
}

// Example 8: TreeControls with permission-based rendering
export function TreeControlsWithPermissions() {
  const [userRole, setUserRole] = useState<'owner' | 'viewer'>('viewer');

  return (
    <ReactFlowProvider>
      <div className="relative w-full h-screen bg-ivory">
        <TreeCanvas />
        
        {/* Role selector */}
        <div className="absolute top-4 left-4 z-10 bg-white p-4 rounded-lg shadow-lg">
          <h3 className="font-bold mb-2">User Role</h3>
          <select
            value={userRole}
            onChange={(e) => setUserRole(e.target.value as 'owner' | 'viewer')}
            className="border rounded px-2 py-1"
          >
            <option value="owner">Owner</option>
            <option value="viewer">Viewer</option>
          </select>
        </div>

        {/* Controls - only show add buttons for owners */}
        <div className="absolute top-4 right-4 z-10">
          <TreeControls
            onAddMember={
              userRole === 'owner'
                ? () => console.log('Add member')
                : undefined
            }
            onAddRelationship={
              userRole === 'owner'
                ? () => console.log('Add relationship')
                : undefined
            }
          />
        </div>
      </div>
    </ReactFlowProvider>
  );
}

// Default export with all examples
export default function TreeControlsExamples() {
  const [activeExample, setActiveExample] = useState<string>('basic');

  const examples = {
    basic: <BasicTreeControls />,
    withActions: <TreeControlsWithActions />,
    mobile: <MobileTreeControls />,
    addMemberOnly: <TreeControlsWithAddMemberOnly />,
    customStyling: <TreeControlsWithCustomStyling />,
    dual: <DualTreeControls />,
    withState: <TreeControlsWithState />,
    withPermissions: <TreeControlsWithPermissions />,
  };

  return (
    <div className="w-full h-screen flex flex-col">
      {/* Example selector */}
      <div className="bg-charcoal text-white p-4 flex gap-2 overflow-x-auto">
        <button
          onClick={() => setActiveExample('basic')}
          className={`px-4 py-2 rounded ${
            activeExample === 'basic' ? 'bg-saffron' : 'bg-teal'
          }`}
        >
          Basic
        </button>
        <button
          onClick={() => setActiveExample('withActions')}
          className={`px-4 py-2 rounded ${
            activeExample === 'withActions' ? 'bg-saffron' : 'bg-teal'
          }`}
        >
          With Actions
        </button>
        <button
          onClick={() => setActiveExample('mobile')}
          className={`px-4 py-2 rounded ${
            activeExample === 'mobile' ? 'bg-saffron' : 'bg-teal'
          }`}
        >
          Mobile
        </button>
        <button
          onClick={() => setActiveExample('addMemberOnly')}
          className={`px-4 py-2 rounded ${
            activeExample === 'addMemberOnly' ? 'bg-saffron' : 'bg-teal'
          }`}
        >
          Add Member Only
        </button>
        <button
          onClick={() => setActiveExample('customStyling')}
          className={`px-4 py-2 rounded ${
            activeExample === 'customStyling' ? 'bg-saffron' : 'bg-teal'
          }`}
        >
          Custom Styling
        </button>
        <button
          onClick={() => setActiveExample('dual')}
          className={`px-4 py-2 rounded ${
            activeExample === 'dual' ? 'bg-saffron' : 'bg-teal'
          }`}
        >
          Dual Controls
        </button>
        <button
          onClick={() => setActiveExample('withState')}
          className={`px-4 py-2 rounded ${
            activeExample === 'withState' ? 'bg-saffron' : 'bg-teal'
          }`}
        >
          With State
        </button>
        <button
          onClick={() => setActiveExample('withPermissions')}
          className={`px-4 py-2 rounded ${
            activeExample === 'withPermissions' ? 'bg-saffron' : 'bg-teal'
          }`}
        >
          With Permissions
        </button>
      </div>

      {/* Active example */}
      <div className="flex-1">
        {examples[activeExample as keyof typeof examples]}
      </div>
    </div>
  );
}
