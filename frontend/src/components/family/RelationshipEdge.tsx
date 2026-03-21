import React from 'react';
import { BaseEdge, getBezierPath } from 'reactflow';
import type { EdgeProps } from 'reactflow';

export interface RelationshipEdgeData {
  type: 'SPOUSE_OF' | 'PARENT_OF' | 'SIBLING_OF';
}

const RelationshipEdge: React.FC<EdgeProps<RelationshipEdgeData>> = ({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  data,
  markerEnd,
}) => {
  // Get color based on relationship type
  const getEdgeColor = () => {
    switch (data?.type) {
      case 'SPOUSE_OF':
        return '#E11D48'; // rose
      case 'PARENT_OF':
        return '#14b8a6'; // teal
      case 'SIBLING_OF':
        return '#F59E0B'; // amber
      default:
        return '#9CA3AF'; // gray fallback
    }
  };

  const [edgePath] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  const edgeColor = getEdgeColor();

  return (
    <BaseEdge
      id={id}
      path={edgePath}
      markerEnd={markerEnd}
      style={{
        stroke: edgeColor,
        strokeWidth: 2,
      }}
    />
  );
};

export default RelationshipEdge;
