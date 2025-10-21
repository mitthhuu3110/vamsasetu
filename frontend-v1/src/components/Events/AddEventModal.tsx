import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { eventsApi, familyApi } from '../../services/api';
import { Event, EventType, NotificationMethod, EventAttendee } from '../../types';
import { XMarkIcon } from '@heroicons/react/24/outline';
import LoadingSpinner from '../UI/LoadingSpinner';
import toast from 'react-hot-toast';

interface AddEventModalProps {
  onClose: () => void;
  onSuccess: () => void;
}

interface FormData {
  title: string;
  description?: string;
  eventType: EventType;
  date: string;
  time?: string;
  location?: string;
  isRecurring: boolean;
  reminderEnabled: boolean;
  reminderAdvanceTime: number;
  reminderMethods: NotificationMethod[];
}

const AddEventModal: React.FC<AddEventModalProps> = ({ onClose, onSuccess }) => {
  const [loading, setLoading] = useState(false);
  const [attendees, setAttendees] = useState<EventAttendee[]>([]);
  
  const { register, handleSubmit, watch, setValue, formState: { errors } } = useForm<FormData>({
    defaultValues: {
      eventType: EventType.CUSTOM,
      isRecurring: false,
      reminderEnabled: true,
      reminderAdvanceTime: 24,
      reminderMethods: [NotificationMethod.EMAIL],
    }
  });

  const isRecurring = watch('isRecurring');
  const reminderEnabled = watch('reminderEnabled');

  const onSubmit = async (data: FormData) => {
    setLoading(true);
    try {
      await eventsApi.createEvent({
        title: data.title,
        description: data.description,
        eventType: data.eventType,
        date: data.date,
        time: data.time,
        location: data.location,
        isRecurring: data.isRecurring,
        recurrencePattern: data.isRecurring ? {
          frequency: 'YEARLY',
          interval: 1,
        } : undefined,
        reminderSettings: {
          enabled: data.reminderEnabled,
          methods: data.reminderMethods,
          advanceTime: data.reminderAdvanceTime,
        },
        attendees: attendees,
        createdBy: 'current-user-id', // This should come from auth context
      });
      toast.success('Event created successfully!');
      onSuccess();
    } catch (error: any) {
      toast.error(error.message || 'Failed to create event');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="flex justify-between items-center p-6 border-b border-gray-200">
          <h2 className="text-xl font-semibold text-gray-900">Add Event</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="p-6 space-y-6">
          {/* Basic Information */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Event Title *
              </label>
              <input
                {...register('title', { required: 'Event title is required' })}
                className="input-field"
                placeholder="Enter event title"
              />
              {errors.title && (
                <p className="text-red-500 text-sm mt-1">{errors.title.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Event Type *
              </label>
              <select {...register('eventType')} className="input-field">
                <option value={EventType.BIRTHDAY}>Birthday</option>
                <option value={EventType.ANNIVERSARY}>Anniversary</option>
                <option value={EventType.WEDDING}>Wedding</option>
                <option value={EventType.FESTIVAL}>Festival</option>
                <option value={EventType.FUNERAL}>Funeral</option>
                <option value={EventType.CUSTOM}>Custom</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Date *
              </label>
              <input
                {...register('date', { required: 'Date is required' })}
                type="date"
                className="input-field"
              />
              {errors.date && (
                <p className="text-red-500 text-sm mt-1">{errors.date.message}</p>
              )}
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Description
            </label>
            <textarea
              {...register('description')}
              rows={3}
              className="input-field"
              placeholder="Enter event description"
            />
          </div>

          {/* Time and Location */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Time
              </label>
              <input
                {...register('time')}
                type="time"
                className="input-field"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Location
              </label>
              <input
                {...register('location')}
                className="input-field"
                placeholder="Enter event location"
              />
            </div>
          </div>

          {/* Recurring */}
          <div>
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={isRecurring}
                onChange={(e) => setValue('isRecurring', e.target.checked)}
                className="mr-2"
              />
              <span className="text-sm font-medium text-gray-700">This is a recurring event</span>
            </label>
          </div>

          {/* Reminder Settings */}
          <div className="border-t border-gray-200 pt-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Reminder Settings</h3>
            
            <div className="mb-4">
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={reminderEnabled}
                  onChange={(e) => setValue('reminderEnabled', e.target.checked)}
                  className="mr-2"
                />
                <span className="text-sm font-medium text-gray-700">Enable reminders</span>
              </label>
            </div>

            {reminderEnabled && (
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Remind me
                  </label>
                  <div className="flex items-center space-x-2">
                    <input
                      {...register('reminderAdvanceTime', { valueAsNumber: true })}
                      type="number"
                      min="1"
                      max="168"
                      className="input-field w-20"
                    />
                    <span className="text-sm text-gray-500">hours before</span>
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Notification Methods
                  </label>
                  <div className="space-y-2">
                    {Object.values(NotificationMethod).map((method) => (
                      <label key={method} className="flex items-center">
                        <input
                          type="checkbox"
                          value={method}
                          {...register('reminderMethods')}
                          className="mr-2"
                        />
                        <span className="text-sm text-gray-700">{method}</span>
                      </label>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Actions */}
          <div className="flex justify-end space-x-3 pt-4 border-t border-gray-200">
            <button
              type="button"
              onClick={onClose}
              className="btn-secondary"
              disabled={loading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="btn-primary"
              disabled={loading}
            >
              {loading ? <LoadingSpinner size="sm" /> : 'Create Event'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddEventModal;
