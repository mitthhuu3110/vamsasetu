import React from 'react';
import { useForm } from 'react-hook-form';
import { useCreateEvent } from '../../hooks/useEvents';
import Modal from '../ui/Modal';
import Input from '../ui/Input';
import Button from '../ui/Button';
import type { CreateEventRequest } from '../../types/event';

interface AddEventModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const AddEventModal: React.FC<AddEventModalProps> = ({ isOpen, onClose }) => {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateEventRequest>();

  const { mutate: createEvent, isPending } = useCreateEvent();

  const onSubmit = (data: CreateEventRequest) => {
    createEvent(data, {
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
    <Modal isOpen={isOpen} onClose={handleClose} title="Add Event">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <Input
          label="Event Title"
          type="text"
          fullWidth
          error={errors.title?.message}
          {...register('title', {
            required: 'Title is required',
            minLength: {
              value: 2,
              message: 'Title must be at least 2 characters',
            },
          })}
        />

        <div>
          <label className="block text-sm font-medium text-charcoal mb-2">
            Description (optional)
          </label>
          <textarea
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-saffron focus:border-transparent resize-none"
            rows={3}
            {...register('description')}
          />
        </div>

        <Input
          label="Event Date"
          type="date"
          fullWidth
          error={errors.eventDate?.message}
          {...register('eventDate', {
            required: 'Event date is required',
          })}
        />

        <div>
          <label className="block text-sm font-medium text-charcoal mb-2">
            Event Type
          </label>
          <select
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-saffron focus:border-transparent"
            {...register('eventType', { required: 'Event type is required' })}
          >
            <option value="">Select event type</option>
            <option value="birthday">🎂 Birthday</option>
            <option value="anniversary">💍 Anniversary</option>
            <option value="ceremony">🎉 Ceremony</option>
            <option value="custom">📌 Custom</option>
          </select>
          {errors.eventType && (
            <p className="text-rose text-sm mt-1">{errors.eventType.message}</p>
          )}
        </div>

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
            {isPending ? 'Adding...' : 'Add Event'}
          </Button>
        </div>
      </form>
    </Modal>
  );
};

export default AddEventModal;
