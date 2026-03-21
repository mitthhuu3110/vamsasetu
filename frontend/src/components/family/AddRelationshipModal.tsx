import React from 'react';
import { useForm } from 'react-hook-form';
import { useCreateRelationship } from '../../hooks/useRelationships';
import { useMembers } from '../../hooks/useMembers';
import Modal from '../ui/Modal';
import Button from '../ui/Button';
import type { CreateRelationshipRequest } from '../../types/relationship';
import type { Member } from '../../types/member';

interface AddRelationshipModalProps {
  isOpen: boolean;
  onClose: () => void;
}

/**
 * AddRelationshipModal Component
 * 
 * Modal for creating a new relationship between two family members.
 * 
 * Features:
 * - Dropdown selection for both members
 * - Relationship type selection (SPOUSE_OF, PARENT_OF, SIBLING_OF)
 * - Form validation
 * - Integration with useCreateRelationship hook
 * 
 * @example
 * ```tsx
 * <AddRelationshipModal
 *   isOpen={showModal}
 *   onClose={() => setShowModal(false)}
 * />
 * ```
 * 
 * **Validates: Requirements 2.4**
 */
const AddRelationshipModal: React.FC<AddRelationshipModalProps> = ({
  isOpen,
  onClose,
}) => {
  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors },
  } = useForm<CreateRelationshipRequest>();

  const { data: membersResponse } = useMembers();
  const { mutate: createRelationship, isPending } = useCreateRelationship();

  const members = membersResponse?.data?.members || [];
  const fromMemberId = watch('fromId');

  const onSubmit = (data: CreateRelationshipRequest) => {
    createRelationship(data, {
      onSuccess: () => {
        reset();
        onClose();
      },
    });
  };

  const handleClose = () => {
    reset();
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title="Add Relationship">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {/* From Member */}
        <div>
          <label className="block text-sm font-medium text-charcoal mb-2">
            From Member
          </label>
          <select
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-saffron focus:border-transparent"
            {...register('fromId', { required: 'Please select a member' })}
          >
            <option value="">Select member</option>
            {members.map((member: Member) => (
              <option key={member.id} value={member.id}>
                {member.name}
              </option>
            ))}
          </select>
          {errors.fromId && (
            <p className="text-rose text-sm mt-1">{errors.fromId.message}</p>
          )}
        </div>

        {/* Relationship Type */}
        <div>
          <label className="block text-sm font-medium text-charcoal mb-2">
            Relationship Type
          </label>
          <select
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-saffron focus:border-transparent"
            {...register('type', { required: 'Please select a relationship type' })}
          >
            <option value="">Select type</option>
            <option value="SPOUSE_OF">Spouse</option>
            <option value="PARENT_OF">Parent (from is parent of to)</option>
            <option value="SIBLING_OF">Sibling</option>
          </select>
          {errors.type && (
            <p className="text-rose text-sm mt-1">{errors.type.message}</p>
          )}
        </div>

        {/* To Member */}
        <div>
          <label className="block text-sm font-medium text-charcoal mb-2">
            To Member
          </label>
          <select
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-saffron focus:border-transparent"
            {...register('toId', {
              required: 'Please select a member',
              validate: (value) =>
                value !== fromMemberId || 'Cannot create relationship with the same member',
            })}
          >
            <option value="">Select member</option>
            {members.map((member: Member) => (
              <option key={member.id} value={member.id}>
                {member.name}
              </option>
            ))}
          </select>
          {errors.toId && (
            <p className="text-rose text-sm mt-1">{errors.toId.message}</p>
          )}
        </div>

        {/* Help Text */}
        <div className="bg-ivory rounded-lg p-3">
          <p className="text-xs text-charcoal/70">
            <strong>Note:</strong> For PARENT_OF relationships, the "From Member" is the parent
            and the "To Member" is the child.
          </p>
        </div>

        {/* Actions */}
        <div className="flex space-x-3 pt-4">
          <Button
            type="button"
            variant="outline"
            onClick={handleClose}
            fullWidth
            disabled={isPending}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            variant="primary"
            fullWidth
            isLoading={isPending}
            disabled={isPending}
          >
            {isPending ? 'Adding...' : 'Add Relationship'}
          </Button>
        </div>
      </form>
    </Modal>
  );
};

export default AddRelationshipModal;
