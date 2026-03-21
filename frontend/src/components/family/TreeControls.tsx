import React from 'react';
import { useReactFlow } from 'reactflow';
import Button from '../ui/Button';

export interface TreeControlsProps {
  onAddMember?: () => void;
  onAddRelationship?: () => void;
  className?: string;
}

/**
 * TreeControls Component
 * 
 * Provides action buttons for the family tree visualization:
 * - Zoom in/out controls
 * - Fit view to see entire tree
 * - Add member button
 * - Add relationship button
 * 
 * Uses ReactFlow's useReactFlow hook for zoom and fit view functionality.
 * 
 * @example
 * ```tsx
 * <TreeControls
 *   onAddMember={() => setShowAddMemberModal(true)}
 *   onAddRelationship={() => setShowAddRelationshipModal(true)}
 * />
 * ```
 * 
 * **Validates: Requirements 3.5**
 */
const TreeControls: React.FC<TreeControlsProps> = ({
  onAddMember,
  onAddRelationship,
  className = '',
}) => {
  const { zoomIn, zoomOut, fitView } = useReactFlow();

  const handleZoomIn = () => {
    zoomIn({ duration: 300 });
  };

  const handleZoomOut = () => {
    zoomOut({ duration: 300 });
  };

  const handleFitView = () => {
    fitView({ duration: 300, padding: 0.2 });
  };

  return (
    <div
      className={`flex flex-col gap-2 bg-white rounded-lg shadow-lg border border-gray-200 p-2 ${className}`}
      role="toolbar"
      aria-label="Family tree controls"
    >
      {/* Zoom Controls */}
      <div className="flex flex-col gap-1 pb-2 border-b border-gray-200">
        <Button
          variant="outline"
          size="sm"
          onClick={handleZoomIn}
          aria-label="Zoom in"
          title="Zoom in"
          className="w-10 h-10 p-0 flex items-center justify-center"
        >
          <span className="text-xl leading-none">+</span>
        </Button>
        
        <Button
          variant="outline"
          size="sm"
          onClick={handleZoomOut}
          aria-label="Zoom out"
          title="Zoom out"
          className="w-10 h-10 p-0 flex items-center justify-center"
        >
          <span className="text-xl leading-none">−</span>
        </Button>
        
        <Button
          variant="outline"
          size="sm"
          onClick={handleFitView}
          aria-label="Fit view"
          title="Fit view to see entire tree"
          className="w-10 h-10 p-0 flex items-center justify-center"
        >
          <span className="text-lg leading-none">⊡</span>
        </Button>
      </div>

      {/* Action Controls */}
      <div className="flex flex-col gap-1 pt-1">
        {onAddMember && (
          <Button
            variant="primary"
            size="sm"
            onClick={onAddMember}
            aria-label="Add family member"
            title="Add a new family member"
            className="w-10 h-10 p-0 flex items-center justify-center"
          >
            <span className="text-xl leading-none">👤</span>
          </Button>
        )}
        
        {onAddRelationship && (
          <Button
            variant="secondary"
            size="sm"
            onClick={onAddRelationship}
            aria-label="Add relationship"
            title="Add a relationship between members"
            className="w-10 h-10 p-0 flex items-center justify-center"
          >
            <span className="text-xl leading-none">🔗</span>
          </Button>
        )}
      </div>
    </div>
  );
};

export default TreeControls;
