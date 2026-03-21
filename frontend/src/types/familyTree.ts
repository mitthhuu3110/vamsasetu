/**
 * Family tree types for React Flow visualization
 */

export interface Position {
  x: number;
  y: number;
}

export interface MemberNodeData {
  id: string;
  name: string;
  avatarUrl: string;
  relationBadge: string;
  hasUpcomingEvent: boolean;
  gender: string;
}

export interface ReactFlowNode {
  id: string;
  type: string;
  position: Position;
  data: MemberNodeData;
}

export interface ReactFlowEdge {
  id: string;
  source: string;
  target: string;
  type: string;
  animated: boolean;
  style: {
    stroke: string;
    strokeWidth: string;
  };
  label?: string;
}

export interface FamilyTree {
  nodes: ReactFlowNode[];
  edges: ReactFlowEdge[];
}
