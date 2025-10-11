import React, { useCallback, useMemo } from 'react';
import ReactFlow, {
  Node,
  Edge,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
  Connection,
  NodeTypes,
  EdgeTypes,
} from 'react-flow-renderer';
import { FamilyMember, Relationship } from '../../types';
import FamilyMemberNode from './FamilyMemberNode';
import RelationshipEdge from './RelationshipEdge';

interface FamilyTreeVisualizerProps {
  members: FamilyMember[];
  relationships: Relationship[];
  onMemberSelect: (member: FamilyMember | null) => void;
  selectedMember: FamilyMember | null;
}

const nodeTypes: NodeTypes = {
  familyMember: FamilyMemberNode,
};

const edgeTypes: EdgeTypes = {
  relationship: RelationshipEdge,
};

const FamilyTreeVisualizer: React.FC<FamilyTreeVisualizerProps> = ({
  members,
  relationships,
  onMemberSelect,
  selectedMember,
}) => {
  // Convert family members to React Flow nodes
  const initialNodes: Node[] = useMemo(() => {
    return members.map((member, index) => ({
      id: member.id,
      type: 'familyMember',
      position: {
        x: (index % 4) * 200,
        y: Math.floor(index / 4) * 150,
      },
      data: {
        member,
        isSelected: selectedMember?.id === member.id,
        onSelect: () => onMemberSelect(member),
      },
    }));
  }, [members, selectedMember, onMemberSelect]);

  // Convert relationships to React Flow edges
  const initialEdges: Edge[] = useMemo(() => {
    return relationships.map((relationship) => ({
      id: relationship.id,
      source: relationship.fromMemberId,
      target: relationship.toMemberId,
      type: 'relationship',
      data: {
        relationshipType: relationship.relationshipType,
        isActive: relationship.isActive,
      },
      style: {
        stroke: relationship.isActive ? '#0ea5e9' : '#d1d5db',
        strokeWidth: relationship.isActive ? 2 : 1,
      },
    }));
  }, [relationships]);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  // Update nodes when selectedMember changes
  React.useEffect(() => {
    setNodes((nds) =>
      nds.map((node) => ({
        ...node,
        data: {
          ...node.data,
          isSelected: selectedMember?.id === node.id,
        },
      }))
    );
  }, [selectedMember, setNodes]);

  return (
    <div className="w-full h-full">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        fitView
        attributionPosition="bottom-left"
      >
        <Controls />
        <Background color="#f3f4f6" gap={20} />
      </ReactFlow>
    </div>
  );
};

export default FamilyTreeVisualizer;
