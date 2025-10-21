import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { familyApi } from '../../services/api';
import { FamilyMember, Gender } from '../../types';
import { XMarkIcon } from '@heroicons/react/24/outline';
import LoadingSpinner from '../UI/LoadingSpinner';
import toast from 'react-hot-toast';

interface AddMemberModalProps {
  onClose: () => void;
  onSuccess: () => void;
}

interface FormData {
  firstName: string;
  lastName: string;
  middleName?: string;
  dateOfBirth?: string;
  dateOfDeath?: string;
  gender: Gender;
  phone?: string;
  email?: string;
  occupation?: string;
  notes?: string;
  isAlive: boolean;
}

const AddMemberModal: React.FC<AddMemberModalProps> = ({ onClose, onSuccess }) => {
  const [loading, setLoading] = useState(false);
  const { register, handleSubmit, watch, setValue, formState: { errors } } = useForm<FormData>({
    defaultValues: {
      gender: Gender.MALE,
      isAlive: true,
    }
  });

  const isAlive = watch('isAlive');

  const onSubmit = async (data: FormData) => {
    setLoading(true);
    try {
      await familyApi.createMember({
        ...data,
        dateOfDeath: data.isAlive ? undefined : data.dateOfDeath,
      });
      toast.success('Family member added successfully!');
      onSuccess();
    } catch (error: any) {
      toast.error(error.message || 'Failed to add family member');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="flex justify-between items-center p-6 border-b border-gray-200">
          <h2 className="text-xl font-semibold text-gray-900">Add Family Member</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="p-6 space-y-6">
          {/* Basic Information */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                First Name *
              </label>
              <input
                {...register('firstName', { required: 'First name is required' })}
                className="input-field"
                placeholder="Enter first name"
              />
              {errors.firstName && (
                <p className="text-red-500 text-sm mt-1">{errors.firstName.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Middle Name
              </label>
              <input
                {...register('middleName')}
                className="input-field"
                placeholder="Enter middle name"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Last Name *
              </label>
              <input
                {...register('lastName', { required: 'Last name is required' })}
                className="input-field"
                placeholder="Enter last name"
              />
              {errors.lastName && (
                <p className="text-red-500 text-sm mt-1">{errors.lastName.message}</p>
              )}
            </div>
          </div>

          {/* Gender and Status */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Gender *
              </label>
              <select {...register('gender')} className="input-field">
                <option value={Gender.MALE}>Male</option>
                <option value={Gender.FEMALE}>Female</option>
                <option value={Gender.OTHER}>Other</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Status
              </label>
              <div className="flex space-x-4">
                <label className="flex items-center">
                  <input
                    type="radio"
                    value="true"
                    checked={isAlive}
                    onChange={() => setValue('isAlive', true)}
                    className="mr-2"
                  />
                  Alive
                </label>
                <label className="flex items-center">
                  <input
                    type="radio"
                    value="false"
                    checked={!isAlive}
                    onChange={() => setValue('isAlive', false)}
                    className="mr-2"
                  />
                  Deceased
                </label>
              </div>
            </div>
          </div>

          {/* Dates */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Date of Birth
              </label>
              <input
                {...register('dateOfBirth')}
                type="date"
                className="input-field"
              />
            </div>

            {!isAlive && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Date of Death
                </label>
                <input
                  {...register('dateOfDeath')}
                  type="date"
                  className="input-field"
                />
              </div>
            )}
          </div>

          {/* Contact Information */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Phone
              </label>
              <input
                {...register('phone')}
                type="tel"
                className="input-field"
                placeholder="Enter phone number"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Email
              </label>
              <input
                {...register('email')}
                type="email"
                className="input-field"
                placeholder="Enter email address"
              />
            </div>
          </div>

          {/* Additional Information */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Occupation
            </label>
            <input
              {...register('occupation')}
              className="input-field"
              placeholder="Enter occupation"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Notes
            </label>
            <textarea
              {...register('notes')}
              rows={3}
              className="input-field"
              placeholder="Enter any additional notes"
            />
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
              {loading ? <LoadingSpinner size="sm" /> : 'Add Member'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddMemberModal;
