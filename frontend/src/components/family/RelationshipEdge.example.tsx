import React from 'react';
import ReactFlow, { Background, Controls } from 'reactflow';
import type { Node, Edge } from 'reactflow';
import 'reactflow/dist/style.css';
import MemberNode from './MemberNode';
import RelationshipEdge from './RelationshipEdge';
import type { MemberNodeData } from './MemberNode';
import type { RelationshipEdgeData } from './RelationshipEdge';

/**
 * Example demonstrating RelationshipEdge component usage
 * 
 * This example shows:
 * - Color-coded edges for different relationship types
 * - SPOUSE_OF: rose color (#E11D48)
 * - PARENT_OF: teal color (#14b8a6)
 * - SIBLING_OF: amber color (#F59E0B)
 * - Bezier curve styling for smooth connections
 */

const nodeTypes = {
  memberNode: MemberNode,
};

const edgeTypes = {
  relationshipEdge: RelationshipEdge,
};

const RelationshipEdgeExample: React.FC = () => {
  // Sample nodes
  const nodes: Node<MemberNodeData>[] = [
    {
      id: '1',
      type: 'memberNode',
      position: { x: 100, y: 100 },
      data: {
        id: '1',
        name: 'John Doe',
        gender: 'male',
      },
    },
    {
      id: '2',
      type: 'memberNode',
      position: { x: 300, y: 100 },
      data: {
        id: '2',
        name: 'Jane Doe',
        gender: 'female',
      },
    },
    {
      id: '3',
      type: 'memberNode',
      position: { x: 200, y: 250 },
      data: {
        id: '3',
        name: 'Child Doe',
        gender: 'other',
      },
    },
    {
      id: '4',
      type: 'memberNode',
      position: { x: 500, y: 100 },
      data: {
        id: '4',
        name: 'Sibling Doe',
        gender: 'male',
      },
    },
  ];

  // Sample edges with different relationship types
  const edges: Edge<RelationshipEdgeData>[] = [
    {
      id: 'e1-2',
      source: '1',
      target: '2',
      type: 'relationshipEdge',
      data: { type: 'SPOUSE_OF' }, // Rose color
    },
    {
      id: 'e1-3',
      source: '1',
      target: '3',
      type: 'relationshipEdge',
      data: { type: 'PARENT_OF' }, // Teal color
    },
    {
      id: 'e2-3',
      source: '2',
      target: '3',
      type: 'relationshipEdge',
      data: { type: 'PARENT_OF' }, // Teal color
    },
    {
      id: 'e1-4',
      source: '1',
      target: '4',
      type: 'relationshipEdge',
      data: { type: 'SIBLING_OF' }, // Amber color
    },
  ];

  return (
    <div style={{ width: '100%', height: '600px' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        fitView
      >
        <Background />
        <Controls />
      </ReactFlow>
    </div>
  );
};

export default RelationshipEdgeExample;
