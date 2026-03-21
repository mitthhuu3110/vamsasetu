import React, { useCallback } from 'react';
import ReactFlow, { 
  Controls, 
  Background,
  useNodesState,
  useEdgesState,
} from 'reactflow';
import type { Node, Edge } from 'reactflow';
import 'reactflow/dist/style.css';
import MemberNode from './MemberNode';
import type { MemberNodeData } from './MemberNode';

// Register custom node type
const nodeTypes = {
  memberNode: MemberNode,
};

// Example usage of MemberNode component
const MemberNodeExample: React.FC = () => {
  const handleNodeClick = useCallback((id: string) => {
    console.log('Member node clicked:', id);
    // In real implementation, this would open a member details panel
  }, []);

  const initialNodes: Node<MemberNodeData>[] = [
    {
      id: '1',
      type: 'memberNode',
      position: { x: 250, y: 0 },
      data: {
        id: '1',
        name: 'Rajesh Kumar',
        gender: 'male',
        avatarUrl: 'https://i.pravatar.cc/150?img=12',
        relationBadge: 'Father',
        hasUpcomingEvent: true,
        onNodeClick: handleNodeClick,
      },
    },
    {
      id: '2',
      type: 'memberNode',
      position: { x: 100, y: 150 },
      data: {
        id: '2',
        name: 'Priya Kumar',
        gender: 'female',
        avatarUrl: 'https://i.pravatar.cc/150?img=5',
        relationBadge: 'Daughter',
        hasUpcomingEvent: false,
        onNodeClick: handleNodeClick,
      },
    },
    {
      id: '3',
      type: 'memberNode',
      position: { x: 400, y: 150 },
      data: {
        id: '3',
        name: 'Amit Kumar',
        gender: 'male',
        avatarUrl: '',
        relationBadge: 'Son',
        hasUpcomingEvent: true,
        onNodeClick: handleNodeClick,
      },
    },
  ];

  const initialEdges: Edge[] = [
    { id: 'e1-2', source: '1', target: '2', type: 'smoothstep' },
    { id: 'e1-3', source: '1', target: '3', type: 'smoothstep' },
  ];

  const [nodes] = useNodesState(initialNodes);
  const [edges] = useEdgesState(initialEdges);

  return (
    <div className="w-full h-screen">
      <div className="p-4 bg-white border-b border-gray-200">
        <h2 className="text-2xl font-heading font-semibold text-charcoal">
          MemberNode Component Example
        </h2>
        <p className="text-sm text-gray-600 mt-1">
          Interactive family tree nodes with gender-based colors and event indicators
        </p>
      </div>
      
      <div className="w-full h-[calc(100vh-100px)]">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          nodeTypes={nodeTypes}
          fitView
          minZoom={0.5}
          maxZoom={2}
        >
          <Background />
          <Controls />
        </ReactFlow>
      </div>

      <div className="absolute bottom-4 left-4 bg-white p-4 rounded-lg shadow-lg max-w-md">
        <h3 className="font-semibold text-charcoal mb-2">Features:</h3>
        <ul className="text-sm text-gray-700 space-y-1">
          <li>• <span className="text-blue-500">Blue border</span> for male members</li>
          <li>• <span className="text-pink-500">Pink border</span> for female members</li>
          <li>• <span className="text-amber-500">Amber glowing indicator</span> for upcoming events</li>
          <li>• Hover effect with scale and shadow</li>
          <li>• Click to open member details (check console)</li>
          <li>• Relation badge display</li>
          <li>• Avatar with fallback to initials</li>
        </ul>
      </div>
    </div>
  );
};

export default MemberNodeExample;
