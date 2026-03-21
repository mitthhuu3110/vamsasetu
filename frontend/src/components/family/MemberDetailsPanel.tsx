import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useMembers } from '../../hooks/useMembers';
import { useRelationships } from '../../hooks/useRelationships';
import { useResponsive } from '../../hooks/useResponsive';
import Button from '../ui/Button';
import type { Member } from '../../types/member';

export interface MemberDetailsPanelProps {
  memberId: string | null;
  onClose: () => void;
  onEdit?: (memberId: string) => void;
}

/**
 * MemberDetailsPanel Component
 * 
 * Displays detailed information about a selected family member
 * including their personal details and direct relationships.
 * 
 * Features:
 * - Full member details (name, DOB, gender, email, phone, avatar)
 * - Direct relationships list
 * - Edit and close actions
 * - Slide-in animation
 * 
 * @example
 * ```tsx
 * <MemberDetailsPanel
 *   memberId={selectedMemberId}
 *   onClose={() => setSelectedMemberId(null)}
 *   onEdit={(id) => handleEdit(id)}
 * />
 * ```
 * 
 * **Validates: Requirements 2.1**
 */
const MemberDetailsPanel: React.FC<MemberDetailsPanelProps> = ({
  memberId,
  onClose,
  onEdit,
}) => {
  const { isMobile } = useResponsive();
  const { data: membersResponse } = useMembers();
  const { data: relationshipsResponse } = useRelationships();

  const member = React.useMemo(() => {
    if (!memberId || !membersResponse?.data) return null;
    return membersResponse.data.members.find((m: Member) => m.id === memberId);
  }, [memberId, membersResponse]);

  const directRelationships = React.useMemo(() => {
    if (!memberId || !relationshipsResponse?.data) return [];
    return relationshipsResponse.data.filter(
      (rel) => rel.fromId === memberId || rel.toId === memberId
    );
  }, [memberId, relationshipsResponse]);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-IN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const getRelationshipLabel = (type: string) => {
    switch (type) {
      case 'SPOUSE_OF':
        return 'Spouse';
      case 'PARENT_OF':
        return 'Parent/Child';
      case 'SIBLING_OF':
        return 'Sibling';
      default:
        return type;
    }
  };

  const getRelatedMemberName = (rel: any) => {
    if (!membersResponse?.data) return 'Unknown';
    const relatedId = rel.fromId === memberId ? rel.toId : rel.fromId;
    const relatedMember = membersResponse.data.members.find((m: Member) => m.id === relatedId);
    return relatedMember?.name || 'Unknown';
  };

  return (
    <AnimatePresence>
      {memberId && member && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/30 z-40"
          />

          {/* Panel */}
          <motion.div
            initial={isMobile ? { y: '100%' } : { x: '100%' }}
            animate={isMobile ? { y: 0 } : { x: 0 }}
            exit={isMobile ? { y: '100%' } : { x: '100%' }}
            transition={{ type: 'spring', damping: 25, stiffness: 200 }}
            className={`fixed ${
              isMobile ? 'inset-0' : 'right-0 top-0 h-full w-96'
            } bg-white shadow-2xl z-50 overflow-y-auto`}
          >
            {/* Header */}
            <div className="sticky top-0 bg-gradient-to-r from-saffron to-turmeric p-4 flex items-center justify-between">
              <h2 className="text-xl font-display font-bold text-white">
                Member Details
              </h2>
              <button
                onClick={onClose}
                className="text-white hover:bg-white/20 rounded-full p-2 transition-colors"
                aria-label="Close panel"
              >
                <span className="text-2xl leading-none">×</span>
              </button>
            </div>

            {/* Content */}
            <div className="p-6 space-y-6">
              {/* Avatar and Name */}
              <div className="flex flex-col items-center text-center">
                <div className="w-24 h-24 rounded-full overflow-hidden bg-gray-200 mb-4">
                  {member.avatarUrl ? (
                    <img
                      src={member.avatarUrl}
                      alt={member.name}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-4xl text-gray-500 font-semibold">
                      {member.name.charAt(0).toUpperCase()}
                    </div>
                  )}
                </div>
                <h3 className="text-2xl font-display font-bold text-charcoal">
                  {member.name}
                </h3>
                <span
                  className={`inline-block px-3 py-1 rounded-full text-sm font-medium mt-2 ${
                    member.gender === 'male'
                      ? 'bg-blue-100 text-blue-700'
                      : member.gender === 'female'
                      ? 'bg-pink-100 text-pink-700'
                      : 'bg-gray-100 text-gray-700'
                  }`}
                >
                  {member.gender.charAt(0).toUpperCase() + member.gender.slice(1)}
                </span>
              </div>

              {/* Personal Details */}
              <div className="space-y-3">
                <h4 className="text-sm font-semibold text-charcoal/70 uppercase tracking-wide">
                  Personal Information
                </h4>
                
                <div className="bg-ivory rounded-lg p-4 space-y-3">
                  <div>
                    <p className="text-xs text-charcoal/60 mb-1">Date of Birth</p>
                    <p className="text-sm font-medium text-charcoal">
                      {formatDate(member.dateOfBirth)}
                    </p>
                  </div>

                  {member.email && (
                    <div>
                      <p className="text-xs text-charcoal/60 mb-1">Email</p>
                      <p className="text-sm font-medium text-charcoal break-all">
                        {member.email}
                      </p>
                    </div>
                  )}

                  {member.phone && (
                    <div>
                      <p className="text-xs text-charcoal/60 mb-1">Phone</p>
                      <p className="text-sm font-medium text-charcoal">
                        {member.phone}
                      </p>
                    </div>
                  )}
                </div>
              </div>

              {/* Direct Relationships */}
              <div className="space-y-3">
                <h4 className="text-sm font-semibold text-charcoal/70 uppercase tracking-wide">
                  Direct Relationships
                </h4>
                
                {directRelationships.length > 0 ? (
                  <div className="space-y-2">
                    {directRelationships.map((rel, index) => (
                      <div
                        key={index}
                        className="bg-ivory rounded-lg p-3 flex items-center justify-between"
                      >
                        <div>
                          <p className="text-sm font-medium text-charcoal">
                            {getRelatedMemberName(rel)}
                          </p>
                          <p className="text-xs text-charcoal/60">
                            {getRelationshipLabel(rel.type)}
                          </p>
                        </div>
                        <span className="text-teal text-xl">🔗</span>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="bg-ivory rounded-lg p-4 text-center">
                    <p className="text-sm text-charcoal/60">
                      No direct relationships found
                    </p>
                  </div>
                )}
              </div>

              {/* Actions */}
              {onEdit && (
                <div className="pt-4">
                  <Button
                    variant="primary"
                    fullWidth
                    onClick={() => onEdit(member.id)}
                  >
                    Edit Member
                  </Button>
                </div>
              )}
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
};

export default MemberDetailsPanel;
