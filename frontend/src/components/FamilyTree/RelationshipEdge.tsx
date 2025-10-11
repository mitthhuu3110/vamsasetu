import React from 'react';
import { EdgeProps, getBezierPath, EdgeLabelRenderer, BaseEdge } from 'react-flow-renderer';
import { RelationshipType } from '../../types';

interface RelationshipEdgeData {
  relationshipType: RelationshipType;
  isActive: boolean;
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
}) => {
  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  const getRelationshipLabel = (type: RelationshipType): string => {
    const labels: { [key in RelationshipType]: string } = {
      [RelationshipType.PARENT]: 'Parent',
      [RelationshipType.CHILD]: 'Child',
      [RelationshipType.SPOUSE]: 'Spouse',
      [RelationshipType.SIBLING]: 'Sibling',
      [RelationshipType.GRANDPARENT]: 'Grandparent',
      [RelationshipType.GRANDCHILD]: 'Grandchild',
      [RelationshipType.UNCLE]: 'Uncle',
      [RelationshipType.AUNT]: 'Aunt',
      [RelationshipType.NEPHEW]: 'Nephew',
      [RelationshipType.NIECE]: 'Niece',
      [RelationshipType.COUSIN]: 'Cousin',
      [RelationshipType.FATHER_IN_LAW]: 'Father-in-law',
      [RelationshipType.MOTHER_IN_LAW]: 'Mother-in-law',
      [RelationshipType.SON_IN_LAW]: 'Son-in-law',
      [RelationshipType.DAUGHTER_IN_LAW]: 'Daughter-in-law',
      [RelationshipType.BROTHER_IN_LAW]: 'Brother-in-law',
      [RelationshipType.SISTER_IN_LAW]: 'Sister-in-law',
      [RelationshipType.MATERNAL_UNCLE]: 'Maternal Uncle',
      [RelationshipType.PATERNAL_UNCLE]: 'Paternal Uncle',
      [RelationshipType.MATERNAL_AUNT]: 'Maternal Aunt',
      [RelationshipType.PATERNAL_AUNT]: 'Paternal Aunt',
      [RelationshipType.MATERNAL_GRANDFATHER]: 'Maternal Grandfather',
      [RelationshipType.PATERNAL_GRANDFATHER]: 'Paternal Grandfather',
      [RelationshipType.MATERNAL_GRANDMOTHER]: 'Maternal Grandmother',
      [RelationshipType.PATERNAL_GRANDMOTHER]: 'Paternal Grandmother',
    };
    
    return labels[type] || type;
  };

  return (
    <>
      <BaseEdge
        id={id}
        path={edgePath}
        style={{
          stroke: data?.isActive ? '#0ea5e9' : '#d1d5db',
          strokeWidth: data?.isActive ? 2 : 1,
        }}
      />
      <EdgeLabelRenderer>
        <div
          style={{
            position: 'absolute',
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            background: 'white',
            padding: '2px 6px',
            borderRadius: '4px',
            fontSize: '10px',
            fontWeight: 500,
            color: data?.isActive ? '#0ea5e9' : '#6b7280',
            border: `1px solid ${data?.isActive ? '#0ea5e9' : '#d1d5db'}`,
            pointerEvents: 'all',
          }}
        >
          {data?.relationshipType && getRelationshipLabel(data.relationshipType)}
        </div>
      </EdgeLabelRenderer>
    </>
  );
};

export default RelationshipEdge;
