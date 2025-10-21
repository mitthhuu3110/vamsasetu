import { useState } from 'react';
import { motion } from 'motion/react';
import { ZoomIn, ZoomOut, Maximize2, Users, ChevronDown, ChevronUp } from 'lucide-react';
import { Button } from './ui/button';
import { Card } from './ui/card';

interface TreeNode {
  id: string;
  name: string;
  relation: string;
  avatar?: string;
  children?: TreeNode[];
  spouse?: string;
}

export function FamilyTree() {
  const [zoom, setZoom] = useState(1);
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set(['1', '2', '3', '4', '5']));
  const [hoveredNode, setHoveredNode] = useState<string | null>(null);

  const familyData: TreeNode = {
    id: '1',
    name: 'Venkata Rao',
    relation: 'Grandfather',
    spouse: 'Saraswati',
    children: [
      {
        id: '2',
        name: 'Krishna',
        relation: 'Father',
        spouse: 'Lakshmi',
        children: [
          { id: '4', name: 'Ravi', relation: 'Brother', spouse: 'Priya' },
          { id: '5', name: 'You', relation: 'Self' },
          { id: '6', name: 'Sita', relation: 'Sister' },
        ],
      },
      {
        id: '3',
        name: 'Ramesh',
        relation: 'Uncle',
        spouse: 'Radha',
        children: [
          { id: '7', name: 'Arun', relation: 'Cousin' },
          { id: '8', name: 'Meera', relation: 'Cousin' },
        ],
      },
    ],
  };

  const toggleNode = (nodeId: string) => {
    const newExpanded = new Set(expandedNodes);
    if (newExpanded.has(nodeId)) {
      newExpanded.delete(nodeId);
    } else {
      newExpanded.add(nodeId);
    }
    setExpandedNodes(newExpanded);
  };

  const renderNode = (node: TreeNode, generation: number = 0): JSX.Element => {
    const isExpanded = expandedNodes.has(node.id);
    const hasChildren = node.children && node.children.length > 0;
    const isHovered = hoveredNode === node.id;

    return (
      <div key={node.id} className="flex flex-col items-center">
        <motion.div
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.4, delay: generation * 0.1 }}
          className="relative mb-8"
        >
          {/* Connection line to parent */}
          {generation > 0 && (
            <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 w-0.5 h-8 bg-gradient-to-b from-primary to-secondary opacity-40" />
          )}

          {/* Node Card */}
          <motion.div
            whileHover={{ scale: 1.05, y: -4 }}
            onHoverStart={() => setHoveredNode(node.id)}
            onHoverEnd={() => setHoveredNode(null)}
            className="relative"
          >
            <Card
              className={`relative p-4 bg-card border-2 transition-all duration-300 cursor-pointer min-w-[180px] ${
                isHovered
                  ? 'border-primary shadow-xl shadow-primary/20'
                  : 'border-border shadow-lg'
              }`}
            >
              {/* Glow effect on hover */}
              {isHovered && (
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  className="absolute -inset-1 bg-gradient-to-r from-primary to-secondary opacity-20 blur-xl rounded-xl -z-10"
                />
              )}

              <div className="flex flex-col items-center space-y-3">
                {/* Avatar */}
                <div className="relative">
                  <div
                    className={`w-16 h-16 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center overflow-hidden shadow-md border-2 transition-all ${
                      isHovered ? 'border-primary scale-110' : 'border-background'
                    }`}
                  >
                    {node.avatar ? (
                      <img src={node.avatar} alt={node.name} className="w-full h-full object-cover" />
                    ) : (
                      <span className="text-xl text-primary-foreground">
                        {node.name.charAt(0)}
                      </span>
                    )}
                  </div>
                  {node.id === '5' && (
                    <div className="absolute -top-1 -right-1 w-5 h-5 bg-primary rounded-full flex items-center justify-center border-2 border-background">
                      <Users className="w-3 h-3 text-primary-foreground" />
                    </div>
                  )}
                </div>

                {/* Name and Relation */}
                <div className="text-center">
                  <p className="text-foreground mb-1">{node.name}</p>
                  <p className="text-xs text-muted-foreground">{node.relation}</p>
                  {node.spouse && (
                    <p className="text-xs text-primary mt-1">& {node.spouse}</p>
                  )}
                </div>

                {/* Expand/Collapse Button */}
                {hasChildren && (
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => toggleNode(node.id)}
                    className="mt-2 h-7 px-2 border-border hover:bg-muted"
                  >
                    {isExpanded ? (
                      <>
                        <ChevronUp className="w-3 h-3 mr-1" />
                        <span className="text-xs">Hide</span>
                      </>
                    ) : (
                      <>
                        <ChevronDown className="w-3 h-3 mr-1" />
                        <span className="text-xs">Show {node.children?.length}</span>
                      </>
                    )}
                  </Button>
                )}
              </div>
            </Card>
          </motion.div>
        </motion.div>

        {/* Children Nodes */}
        {hasChildren && isExpanded && (
          <div className="relative">
            {/* Horizontal connection line */}
            {node.children && node.children.length > 1 && (
              <div
                className="absolute top-0 left-0 right-0 h-0.5 bg-gradient-to-r from-primary via-secondary to-primary opacity-40"
                style={{
                  left: '50%',
                  right: '50%',
                  transform: 'translateX(-50%)',
                  width: `${(node.children.length - 1) * 220}px`,
                }}
              />
            )}

            <div className="flex justify-center gap-8 pt-8">
              {node.children?.map((child) => renderNode(child, generation + 1))}
            </div>
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-background relative">
      {/* Background Pattern */}
      <div className="absolute inset-0 overflow-hidden opacity-5 pointer-events-none">
        <div className="absolute top-0 left-0 w-full h-full">
          <svg className="w-full h-full" style={{ backgroundImage: 'radial-gradient(circle, #C9A961 1px, transparent 1px)', backgroundSize: '40px 40px' }} />
        </div>
      </div>

      {/* Parchment Texture Overlay */}
      <div className="absolute inset-0 bg-gradient-to-br from-[#FAF7F2] via-transparent to-[#F5EFE6] opacity-50 pointer-events-none" />

      {/* Header */}
      <div className="relative bg-card/80 backdrop-blur-sm border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1>Family Tree</h1>
              <p className="text-sm text-muted-foreground mt-1">
                Explore your family connections across generations
              </p>
            </div>

            {/* Zoom Controls */}
            <div className="flex items-center space-x-2">
              <Button
                size="icon"
                variant="outline"
                onClick={() => setZoom(Math.max(0.5, zoom - 0.1))}
                className="border-border hover:bg-muted"
              >
                <ZoomOut className="w-4 h-4" />
              </Button>
              <span className="text-sm text-muted-foreground min-w-[60px] text-center">
                {Math.round(zoom * 100)}%
              </span>
              <Button
                size="icon"
                variant="outline"
                onClick={() => setZoom(Math.min(1.5, zoom + 0.1))}
                className="border-border hover:bg-muted"
              >
                <ZoomIn className="w-4 h-4" />
              </Button>
              <Button
                size="icon"
                variant="outline"
                onClick={() => setZoom(1)}
                className="border-border hover:bg-muted"
              >
                <Maximize2 className="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Tree Container */}
      <div className="relative overflow-auto p-8">
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.6 }}
          style={{ transform: `scale(${zoom})`, transformOrigin: 'top center' }}
          className="inline-block min-w-full"
        >
          <div className="flex justify-center py-12">
            {renderNode(familyData)}
          </div>
        </motion.div>
      </div>

      {/* Legend */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, delay: 0.3 }}
        className="fixed bottom-6 left-1/2 transform -translate-x-1/2"
      >
        <Card className="px-6 py-3 bg-card/95 backdrop-blur-sm border-border shadow-xl">
          <div className="flex items-center space-x-6 text-sm text-muted-foreground">
            <div className="flex items-center space-x-2">
              <div className="w-3 h-3 rounded-full bg-gradient-to-br from-primary to-secondary" />
              <span>Hover to highlight path</span>
            </div>
            <div className="flex items-center space-x-2">
              <Users className="w-3 h-3" />
              <span>You</span>
            </div>
            <div className="flex items-center space-x-2">
              <ChevronDown className="w-3 h-3" />
              <span>Expand/Collapse</span>
            </div>
          </div>
        </Card>
      </motion.div>
    </div>
  );
}
