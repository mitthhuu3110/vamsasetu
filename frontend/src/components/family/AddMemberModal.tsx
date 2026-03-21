import React from 'react';
import { useForm } from 'react-hook-form';
import { useCreateMember } from '../../hooks/useMembers';
import Modal from '../ui/Modal';
import Input from '../ui/Input';
import Button from '../ui/Button';
import type { CreateMemberRequest } from '../../types/member';

interface AddMemberModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const AddMemberModal: React.FC<AddMemberModalProps> = ({ isOpen, onClose }) => {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateMemberRequest>();

  const { mutate: createMember, isPending } = useCreateMember();

  const onSubmit = (data: CreateMemberRequest) => {
    createMember(data, {
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
    <Modal isOpen={isOpen} onClose={handleClose} title="Add Family Member">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <Input
          label="Full Name"
          type="text"
          fullWidth
          error={errors.name?.message}
          {...register('name', {
            required: 'Name is required',
            minLength: {
              value: 2,
              message: 'Name must be at least 2 characters',
            },
          })}
        />

        <Input
          label="Date of Birth"
          type="date"
          fullWidth
          error={errors.dateOfBirth?.message}
          {...register('dateOfBirth', {
            required: 'Date of birth is required',
          })}
        />

        <div>
          <label className="block text-sm font-medium text-charcoal mb-2">
            Gender
          </label>
          <select
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-saffron focus:border-transparent"
            {...register('gender', { required: 'Gender is required' })}
          >
            <option value="">Select gender</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </select>
          {errors.gender && (
            <p className="text-rose text-sm mt-1">{errors.gender.message}</p>
          )}
        </div>

        <Input
          label="Email (optional)"
          type="email"
          fullWidth
          error={errors.email?.message}
          {...register('email', {
            pattern: {
              value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
              message: 'Invalid email address',
            },
          })}
        />

        <Input
          label="Phone (optional)"
          type="tel"
          fullWidth
          error={errors.phone?.message}
          {...register('phone')}
        />

        <Input
          label="Avatar URL (optional)"
          type="url"
          fullWidth
          error={errors.avatarUrl?.message}
          {...register('avatarUrl')}
        />

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
            {isPending ? 'Adding...' : 'Add Member'}
          </Button>
        </div>
      </form>
    </Modal>
  );
};

export default AddMemberModal;
