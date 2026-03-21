import React, { useCallback } from 'react';
import ReactFlow, {
  Background,
  Controls,
  MiniMap,
  useNodesState,
  useEdgesState,
  type Node,
  type Edge,
} from 'reactflow';
import 'reactflow/dist/style.css';
import { useFamilyTree } from '../../hooks/useFamilyTree';
import { useResponsive } from '../../hooks/useResponsive';
import MemberNode from './MemberNode';
import RelationshipEdge from './RelationshipEdge';

// Register custom node and edge types
const nodeTypes = {
  memberNode: MemberNode,
};

const edgeTypes = {
  relationshipEdge: RelationshipEdge,
};

export interface TreeCanvasProps {
  onNodeClick?: (nodeId: string) => void;
}

const TreeCanvas: React.FC<TreeCanvasProps> = ({ onNodeClick }) => {
  const { isMobile } = useResponsive();
  const { data, isLoading, error } = useFamilyTree();
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);

  // Update nodes and edges when data is fetched
  React.useEffect(() => {
    if (data?.data) {
      setNodes(data.data.nodes as Node[]);
      setEdges(data.data.edges as Edge[]);
    }
  }, [data, setNodes, setEdges]);

  // Handle node click
  const handleNodeClick = useCallback((nodeId: string) => {
    if (onNodeClick) {
      onNodeClick(nodeId);
    }
  }, [onNodeClick]);

  // Add click handler to node data
  const nodesWithClickHandler = React.useMemo(
    () =>
      nodes.map((node) => ({
        ...node,
        data: {
          ...node.data,
          onNodeClick: handleNodeClick,
        },
      })),
    [nodes, handleNodeClick]
  );

  if (isLoading) {
    return (
      <div className="w-full h-full flex items-center justify-center bg-cream">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-saffron border-t-transparent rounded-full animate-spin mx-auto mb-4" />
          <p className="text-charcoal font-medium">Loading family tree...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full h-full flex items-center justify-center bg-cream">
        <div className="text-center max-w-md p-6">
          <div className="text-rose-500 text-5xl mb-4">⚠️</div>
          <h3 className="text-xl font-semibold text-charcoal mb-2">
            Failed to load family tree
          </h3>
          <p className="text-gray-600">{error.message}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full h-full bg-cream">
      <ReactFlow
        nodes={nodesWithClickHandler}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        fitView
        minZoom={0.1}
        maxZoom={2}
        defaultViewport={{ x: 0, y: 0, zoom: 1 }}
        // Mobile gesture support
        panOnScroll
        panOnDrag
        zoomOnPinch
        zoomOnScroll
        zoomOnDoubleClick
        // Styling
        className="bg-cream"
      >
        {/* Background pattern */}
        <Background color="#D4AF37" gap={16} size={1} />

        {/* Zoom and pan controls */}
        <Controls
          className="bg-white shadow-lg rounded-lg border border-gray-200"
          showInteractive={false}
          position={isMobile ? 'bottom-right' : 'top-left'}
        />

        {/* Mini map for navigation - hide on mobile */}
        {!isMobile && (
          <MiniMap
            className="bg-white shadow-lg rounded-lg border border-gray-200"
            nodeColor={(node) => {
              // Color nodes by gender in minimap
              const gender = node.data?.gender;
              if (gender === 'male') return '#3B82F6';
              if (gender === 'female') return '#EC4899';
              return '#9CA3AF';
            }}
            maskColor="rgba(0, 0, 0, 0.1)"
            zoomable
            pannable
          />
        )}
      </ReactFlow>
    </div>
  );
};

export default TreeCanvas;
