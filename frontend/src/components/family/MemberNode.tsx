import React from 'react';
import { Handle, Position } from 'reactflow';
import type { NodeProps } from 'reactflow';
import { motion } from 'framer-motion';

export interface MemberNodeData {
  id: string;
  name: string;
  gender: 'male' | 'female' | 'other';
  avatarUrl?: string;
  relationBadge?: string;
  hasUpcomingEvent?: boolean;
  onNodeClick?: (id: string) => void;
}

const MemberNode: React.FC<NodeProps<MemberNodeData>> = ({ data }) => {
  const {
    id,
    name,
    gender,
    avatarUrl,
    relationBadge,
    hasUpcomingEvent = false,
    onNodeClick,
  } = data;

  // Gender-based border colors
  const getBorderColor = () => {
    switch (gender) {
      case 'male':
        return 'border-blue-500';
      case 'female':
        return 'border-pink-500';
      default:
        return 'border-gray-400';
    }
  };

  const handleClick = () => {
    if (onNodeClick) {
      onNodeClick(id);
    }
  };

  return (
    <>
      {/* Connection handles */}
      <Handle type="target" position={Position.Top} className="!bg-teal opacity-0" />
      <Handle type="source" position={Position.Bottom} className="!bg-teal opacity-0" />

      <motion.div
        className={`relative bg-white rounded-lg shadow-md border-4 ${getBorderColor()} 
          cursor-pointer transition-all duration-200 min-w-[160px]`}
        whileHover={{
          scale: 1.05,
          boxShadow: '0 10px 30px rgba(0, 0, 0, 0.15)',
        }}
        whileTap={{ scale: 0.98 }}
        onClick={handleClick}
      >
        {/* Upcoming event indicator */}
        {hasUpcomingEvent && (
          <motion.div
            className="absolute -top-2 -right-2 w-4 h-4 bg-amber-500 rounded-full"
            animate={{
              boxShadow: [
                '0 0 0 0 rgba(245, 158, 11, 0.7)',
                '0 0 0 8px rgba(245, 158, 11, 0)',
              ],
            }}
            transition={{
              duration: 1.5,
              repeat: Infinity,
              ease: 'easeInOut',
            }}
          />
        )}

        <div className="p-3 flex flex-col items-center gap-2">
          {/* Avatar */}
          <div className="w-16 h-16 rounded-full overflow-hidden bg-gray-200 flex items-center justify-center">
            {avatarUrl ? (
              <img
                src={avatarUrl}
                alt={name}
                className="w-full h-full object-cover"
              />
            ) : (
              <span className="text-2xl text-gray-500 font-semibold">
                {name.charAt(0).toUpperCase()}
              </span>
            )}
          </div>

          {/* Name */}
          <div className="text-center">
            <p className="text-sm font-semibold text-charcoal truncate max-w-[140px]">
              {name}
            </p>
          </div>

          {/* Relation badge */}
          {relationBadge && (
            <div className="px-2 py-1 bg-turmeric/20 text-turmeric text-xs font-medium rounded-full">
              {relationBadge}
            </div>
          )}
        </div>

        {/* Hover glow effect */}
        <div className="absolute inset-0 rounded-lg opacity-0 hover:opacity-100 transition-opacity duration-200 pointer-events-none">
          <div className="absolute inset-0 rounded-lg bg-gradient-to-br from-saffron/10 to-teal/10" />
        </div>
      </motion.div>
    </>
  );
};

export default MemberNode;
