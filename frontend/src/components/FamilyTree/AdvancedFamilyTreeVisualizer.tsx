import React, { useEffect, useRef, useState, useCallback } from 'react';
import * as d3 from 'd3';
import { motion, AnimatePresence } from 'framer-motion';
import { FamilyMember, Relationship } from '../../types';
import { 
  ChevronDownIcon, 
  ChevronRightIcon, 
  UserIcon,
  HeartIcon,
  SparklesIcon,
  PlusIcon,
  MinusIcon
} from '@heroicons/react/24/outline';

interface FamilyTreeVisualizerProps {
  familyMembers: FamilyMember[];
  relationships: Relationship[];
  onMemberClick: (member: FamilyMember) => void;
}

interface TreeNode {
  id: string;
  name: string;
  data: FamilyMember;
  children: TreeNode[];
  parent?: TreeNode;
  x?: number;
  y?: number;
  depth?: number;
  expanded?: boolean;
}

const AdvancedFamilyTreeVisualizer: React.FC<FamilyTreeVisualizerProps> = ({
  familyMembers,
  relationships,
  onMemberClick,
}) => {
  const svgRef = useRef<SVGSVGElement>(null);
  const [treeData, setTreeData] = useState<TreeNode | null>(null);
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set());
  const [hoveredNode, setHoveredNode] = useState<string | null>(null);
  const [selectedNode, setSelectedNode] = useState<string | null>(null);

  // Build tree structure from flat data
  const buildTree = useCallback((members: FamilyMember[], rels: Relationship[]): TreeNode | null => {
    if (members.length === 0) return null;

    // Find root node (person with no parents or the first person)
    const rootMember = members[0];
    
    const buildNode = (member: FamilyMember, parent?: TreeNode): TreeNode => {
      const node: TreeNode = {
        id: member.id,
        name: `${member.firstName} ${member.lastName}`,
        data: member,
        children: [],
        parent,
        expanded: expandedNodes.has(member.id),
      };

      // Find children
      const children = rels
        .filter(rel => rel.fromMemberId === member.id && rel.relationshipType === 'PARENT')
        .map(rel => members.find(m => m.id === rel.toMemberId))
        .filter(Boolean) as FamilyMember[];

      node.children = children.map(child => buildNode(child, node));
      return node;
    };

    return buildNode(rootMember);
  }, [expandedNodes]);

  useEffect(() => {
    const tree = buildTree(familyMembers, relationships);
    setTreeData(tree);
  }, [familyMembers, relationships, buildTree]);

  useEffect(() => {
    if (!treeData || !svgRef.current) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll('*').remove();

    const width = 1200;
    const height = 800;
    const margin = { top: 20, right: 20, bottom: 20, left: 20 };

    svg.attr('width', width).attr('height', height);

    const g = svg.append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    // Create tree layout
    const treeLayout = d3.tree<TreeNode>()
      .size([width - margin.left - margin.right, height - margin.top - margin.bottom])
      .separation((a, b) => (a.parent === b.parent ? 1 : 2) / a.depth!);

    // Process tree data for D3
    const processNode = (node: TreeNode, depth = 0): any => {
      const d3Node: any = {
        id: node.id,
        name: node.name,
        data: node.data,
        depth,
        children: node.children.length > 0 && node.expanded ? 
          node.children.map(child => processNode(child, depth + 1)) : 
          undefined,
        _children: node.children.length > 0 && !node.expanded ? 
          node.children.map(child => processNode(child, depth + 1)) : 
          undefined,
      };
      return d3Node;
    };

    const root = d3.hierarchy(processNode(treeData));
    treeLayout(root);

    // Draw connections with beautiful curves
    const links = g.selectAll('.link')
      .data(root.links())
      .enter().append('path')
      .attr('class', 'link')
      .attr('d', d3.linkVertical<any, any>()
        .x(d => d.x!)
        .y(d => d.y!))
      .style('fill', 'none')
      .style('stroke', '#D4AF37')
      .style('stroke-width', 4)
      .style('stroke-opacity', 0.8)
      .style('filter', 'drop-shadow(0 2px 4px rgba(212, 175, 55, 0.3))')
      .style('stroke-linecap', 'round');

    // Draw nodes
    const nodes = g.selectAll('.node')
      .data(root.descendants())
      .enter().append('g')
      .attr('class', 'node')
      .attr('transform', d => `translate(${d.x},${d.y})`)
      .style('cursor', 'pointer');

    // Add gradient definitions
    const defs = svg.append('defs');
    
    // Gold gradient
    const goldGradient = defs.append('radialGradient')
      .attr('id', 'goldGradient')
      .attr('cx', '30%')
      .attr('cy', '30%')
      .attr('r', '70%');
    
    goldGradient.append('stop')
      .attr('offset', '0%')
      .attr('stop-color', '#FFD700');
    
    goldGradient.append('stop')
      .attr('offset', '100%')
      .attr('stop-color', '#B8860B');

    // Green gradient for expand buttons
    const greenGradient = defs.append('radialGradient')
      .attr('id', 'greenGradient')
      .attr('cx', '30%')
      .attr('cy', '30%')
      .attr('r', '70%');
    
    greenGradient.append('stop')
      .attr('offset', '0%')
      .attr('stop-color', '#90EE90');
    
    greenGradient.append('stop')
      .attr('offset', '100%')
      .attr('stop-color', '#228B22');

    // Add node circles with beautiful styling
    nodes.append('circle')
      .attr('r', 35)
      .attr('fill', 'url(#goldGradient)')
      .attr('stroke', '#B8860B')
      .attr('stroke-width', 3)
      .style('filter', 'drop-shadow(0 4px 8px rgba(0, 0, 0, 0.2))')
      .on('mouseover', function(event, d) {
        setHoveredNode(d.data.id);
        d3.select(this)
          .transition()
          .duration(200)
          .attr('r', 40)
          .style('filter', 'drop-shadow(0 6px 12px rgba(212, 175, 55, 0.4))');
      })
      .on('mouseout', function(event, d) {
        setHoveredNode(null);
        d3.select(this)
          .transition()
          .duration(200)
          .attr('r', 35)
          .style('filter', 'drop-shadow(0 4px 8px rgba(0, 0, 0, 0.2))');
      })
      .on('click', function(event, d) {
        setSelectedNode(d.data.id);
        onMemberClick(d.data.data);
      });

    // Add expand/collapse buttons for nodes with children
    const expandButtons = nodes.filter(d => d.children || d._children);
    
    expandButtons.append('circle')
      .attr('r', 15)
      .attr('cx', 25)
      .attr('cy', -25)
      .attr('fill', 'url(#greenGradient)')
      .attr('stroke', '#228B22')
      .attr('stroke-width', 2)
      .style('cursor', 'pointer')
      .on('click', function(event, d) {
        event.stopPropagation();
        toggleNode(d.data.id);
      });

    // Add expand/collapse icons
    expandButtons.append('text')
      .attr('x', 25)
      .attr('y', -20)
      .attr('text-anchor', 'middle')
      .attr('font-size', '14px')
      .attr('fill', '#228B22')
      .attr('font-weight', 'bold')
      .text(d => d.children ? '−' : '+');

    // Add member names with beautiful typography
    nodes.append('text')
      .attr('dy', 50)
      .attr('text-anchor', 'middle')
      .attr('font-size', '14px')
      .attr('font-weight', 'bold')
      .attr('fill', '#8B4513')
      .text(d => d.data.name.split(' ')[0]) // First name only
      .style('pointer-events', 'none')
      .style('text-shadow', '0 1px 2px rgba(255, 255, 255, 0.8)');

    // Add relationship indicators with emojis
    nodes.append('text')
      .attr('dy', 65)
      .attr('text-anchor', 'middle')
      .attr('font-size', '12px')
      .text(d => {
        if (d.data.data.gender === 'MALE') return '♂';
        if (d.data.data.gender === 'FEMALE') return '♀';
        return '⚥';
      })
      .style('pointer-events', 'none');

    // Add selection indicator
    if (selectedNode) {
      nodes.filter(d => d.data.id === selectedNode)
        .append('circle')
        .attr('r', 45)
        .attr('fill', 'none')
        .attr('stroke', '#FF6B6B')
        .attr('stroke-width', 3)
        .attr('stroke-dasharray', '5,5')
        .style('animation', 'dash 1s linear infinite');
    }

  }, [treeData, expandedNodes, selectedNode, onMemberClick]);

  const toggleNode = (nodeId: string) => {
    setExpandedNodes(prev => {
      const newSet = new Set(prev);
      if (newSet.has(nodeId)) {
        newSet.delete(nodeId);
      } else {
        newSet.add(nodeId);
      }
      return newSet;
    });
  };

  if (!treeData) {
    return (
      <div className="flex items-center justify-center h-full bg-gradient-to-br from-warm-beige to-cream dark:from-dark-bg dark:to-dark-card">
        <motion.div 
          className="text-center"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          <SparklesIcon className="w-16 h-16 text-soft-gold mx-auto mb-4 animate-pulse" />
          <p className="text-warm-brown dark:text-dark-text text-lg font-display">
            Building your family tree...
          </p>
          <p className="text-gray-600 dark:text-gray-400 text-sm mt-2">
            Connecting generations with love
          </p>
        </motion.div>
      </div>
    );
  }

  return (
    <div className="w-full h-full bg-gradient-to-br from-warm-beige to-cream dark:from-dark-bg dark:to-dark-card indian-pattern relative overflow-hidden">
      <div className="absolute inset-0 opacity-10">
        <div className="absolute top-10 left-10 w-32 h-32 bg-soft-gold rounded-full blur-3xl"></div>
        <div className="absolute bottom-10 right-10 w-24 h-24 bg-soft-green rounded-full blur-2xl"></div>
        <div className="absolute top-1/2 left-1/2 w-16 h-16 bg-deep-gold rounded-full blur-xl"></div>
      </div>
      
      <div className="relative w-full h-full">
        <svg
          ref={svgRef}
          className="w-full h-full"
          style={{ background: 'transparent' }}
        />
        
        {/* Hovered member info */}
        <AnimatePresence>
          {hoveredNode && (
            <motion.div
              initial={{ opacity: 0, scale: 0.8, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.8, y: 20 }}
              className="absolute top-4 right-4 bg-white dark:bg-dark-card rounded-xl shadow-xl p-4 border border-gray-200 dark:border-dark-accent max-w-xs backdrop-blur-sm"
            >
              <div className="flex items-center space-x-3">
                <div className="w-12 h-12 bg-gradient-to-br from-soft-gold to-deep-gold rounded-full flex items-center justify-center shadow-lg">
                  <UserIcon className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h3 className="font-display font-bold text-warm-brown dark:text-dark-text">
                    {familyMembers.find(m => m.id === hoveredNode)?.firstName} {familyMembers.find(m => m.id === hoveredNode)?.lastName}
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {familyMembers.find(m => m.id === hoveredNode)?.relationship || 'Family Member'}
                  </p>
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Tree controls */}
        <div className="absolute bottom-4 left-4 flex space-x-2">
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={() => setExpandedNodes(new Set())}
            className="btn-secondary text-sm flex items-center space-x-1"
          >
            <MinusIcon className="w-4 h-4" />
            <span>Collapse All</span>
          </motion.button>
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={() => {
              const allNodeIds = new Set(familyMembers.map(m => m.id));
              setExpandedNodes(allNodeIds);
            }}
            className="btn-primary text-sm flex items-center space-x-1"
          >
            <PlusIcon className="w-4 h-4" />
            <span>Expand All</span>
          </motion.button>
        </div>

        {/* Family tree title */}
        <div className="absolute top-4 left-4">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.2 }}
            className="bg-white dark:bg-dark-card rounded-xl shadow-lg p-3 border border-gray-200 dark:border-dark-accent"
          >
            <h2 className="text-lg font-display font-bold text-gradient">
              🌳 Vamsa Tree
            </h2>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {familyMembers.length} family members
            </p>
          </motion.div>
        </div>
      </div>
    </div>
  );
};

export default AdvancedFamilyTreeVisualizer;
