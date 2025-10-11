import React from 'react';
import { Handle, Position, NodeProps } from 'react-flow-renderer';
import { FamilyMember } from '../../types';

interface FamilyMemberNodeData {
  member: FamilyMember;
  isSelected: boolean;
  onSelect: () => void;
}

const FamilyMemberNode: React.FC<NodeProps<FamilyMemberNodeData>> = ({ data }) => {
  const { member, isSelected, onSelect } = data;

  return (
    <div
      className={`
        bg-white rounded-lg shadow-md border-2 p-3 min-w-[120px] cursor-pointer transition-all duration-200
        ${isSelected ? 'border-primary-500 shadow-lg' : 'border-gray-200 hover:border-gray-300'}
      `}
      onClick={onSelect}
    >
      <Handle type="target" position={Position.Top} className="w-3 h-3" />
      
      <div className="text-center">
        {/* Profile Picture or Initials */}
        <div className="w-12 h-12 mx-auto mb-2 bg-primary-100 rounded-full flex items-center justify-center">
          {member.profilePicture ? (
            <img
              src={member.profilePicture}
              alt={`${member.firstName} ${member.lastName}`}
              className="w-12 h-12 rounded-full object-cover"
            />
          ) : (
            <span className="text-primary-600 font-semibold text-sm">
              {member.firstName[0]}{member.lastName[0]}
            </span>
          )}
        </div>
        
        {/* Name */}
        <div className="text-sm font-medium text-gray-900 truncate">
          {member.firstName} {member.lastName}
        </div>
        
        {/* Additional Info */}
        <div className="text-xs text-gray-500 mt-1">
          {member.dateOfBirth && (
            <div>
              Born: {new Date(member.dateOfBirth).getFullYear()}
            </div>
          )}
          {!member.isAlive && (
            <div className="text-red-500">
              Deceased
            </div>
          )}
        </div>
      </div>
      
      <Handle type="source" position={Position.Bottom} className="w-3 h-3" />
    </div>
  );
};

export default FamilyMemberNode;
