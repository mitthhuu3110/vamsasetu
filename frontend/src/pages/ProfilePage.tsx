import React, { useState } from 'react';
import { useAuth } from '../hooks/useAuth';
import { User } from '../types';
import { 
  UserIcon, 
  EnvelopeIcon, 
  PhoneIcon, 
  CalendarIcon,
  PencilIcon,
  CheckIcon,
  XMarkIcon
} from '@heroicons/react/24/outline';

const ProfilePage: React.FC = () => {
  const { user } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    firstName: user?.firstName || '',
    lastName: user?.lastName || '',
    email: user?.email || '',
    phone: user?.phone || '',
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSave = () => {
    // TODO: Implement profile update
    console.log('Saving profile:', formData);
    setIsEditing(false);
  };

  const handleCancel = () => {
    setFormData({
      firstName: user?.firstName || '',
      lastName: user?.lastName || '',
      email: user?.email || '',
      phone: user?.phone || '',
    });
    setIsEditing(false);
  };

  if (!user) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-600">User not found</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Profile</h1>
          <p className="text-gray-600">Manage your account information</p>
        </div>
        <button
          onClick={() => setIsEditing(!isEditing)}
          className="btn-primary flex items-center"
        >
          <PencilIcon className="w-5 h-5 mr-2" />
          {isEditing ? 'Cancel' : 'Edit Profile'}
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Profile Card */}
        <div className="lg:col-span-1">
          <div className="card text-center">
            <div className="w-24 h-24 mx-auto mb-4 bg-primary-100 rounded-full flex items-center justify-center">
              {user.profilePicture ? (
                <img
                  src={user.profilePicture}
                  alt={`${user.firstName} ${user.lastName}`}
                  className="w-24 h-24 rounded-full object-cover"
                />
              ) : (
                <UserIcon className="w-12 h-12 text-primary-600" />
              )}
            </div>
            <h2 className="text-xl font-semibold text-gray-900">
              {user.firstName} {user.lastName}
            </h2>
            <p className="text-gray-500 capitalize">{user.role.toLowerCase()}</p>
            <div className="mt-4 pt-4 border-t border-gray-200">
              <p className="text-sm text-gray-500">
                Member since {new Date(user.createdAt).toLocaleDateString()}
              </p>
            </div>
          </div>
        </div>

        {/* Profile Information */}
        <div className="lg:col-span-2">
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-6">Personal Information</h3>
            
            <div className="space-y-6">
              {/* Name Fields */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    First Name
                  </label>
                  {isEditing ? (
                    <input
                      type="text"
                      name="firstName"
                      value={formData.firstName}
                      onChange={handleInputChange}
                      className="input-field"
                    />
                  ) : (
                    <p className="text-gray-900 py-2">{user.firstName}</p>
                  )}
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Last Name
                  </label>
                  {isEditing ? (
                    <input
                      type="text"
                      name="lastName"
                      value={formData.lastName}
                      onChange={handleInputChange}
                      className="input-field"
                    />
                  ) : (
                    <p className="text-gray-900 py-2">{user.lastName}</p>
                  )}
                </div>
              </div>

              {/* Email */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Email Address
                </label>
                {isEditing ? (
                  <input
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleInputChange}
                    className="input-field"
                  />
                ) : (
                  <div className="flex items-center space-x-2 py-2">
                    <EnvelopeIcon className="w-5 h-5 text-gray-400" />
                    <p className="text-gray-900">{user.email}</p>
                  </div>
                )}
              </div>

              {/* Phone */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Phone Number
                </label>
                {isEditing ? (
                  <input
                    type="tel"
                    name="phone"
                    value={formData.phone}
                    onChange={handleInputChange}
                    className="input-field"
                    placeholder="Enter phone number"
                  />
                ) : (
                  <div className="flex items-center space-x-2 py-2">
                    <PhoneIcon className="w-5 h-5 text-gray-400" />
                    <p className="text-gray-900">{user.phone || 'Not provided'}</p>
                  </div>
                )}
              </div>

              {/* Account Information */}
              <div className="pt-6 border-t border-gray-200">
                <h4 className="text-md font-medium text-gray-900 mb-4">Account Information</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Account Created
                    </label>
                    <div className="flex items-center space-x-2 py-2">
                      <CalendarIcon className="w-5 h-5 text-gray-400" />
                      <p className="text-gray-900">
                        {new Date(user.createdAt).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Last Updated
                    </label>
                    <div className="flex items-center space-x-2 py-2">
                      <CalendarIcon className="w-5 h-5 text-gray-400" />
                      <p className="text-gray-900">
                        {new Date(user.updatedAt).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Action Buttons */}
              {isEditing && (
                <div className="flex justify-end space-x-3 pt-6 border-t border-gray-200">
                  <button
                    onClick={handleCancel}
                    className="btn-secondary flex items-center"
                  >
                    <XMarkIcon className="w-4 h-4 mr-2" />
                    Cancel
                  </button>
                  <button
                    onClick={handleSave}
                    className="btn-primary flex items-center"
                  >
                    <CheckIcon className="w-4 h-4 mr-2" />
                    Save Changes
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Additional Settings */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Notification Preferences */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Notification Preferences</h3>
          <div className="space-y-3">
            <label className="flex items-center">
              <input type="checkbox" defaultChecked className="mr-3" />
              <span className="text-sm text-gray-700">Email notifications</span>
            </label>
            <label className="flex items-center">
              <input type="checkbox" defaultChecked className="mr-3" />
              <span className="text-sm text-gray-700">SMS notifications</span>
            </label>
            <label className="flex items-center">
              <input type="checkbox" defaultChecked className="mr-3" />
              <span className="text-sm text-gray-700">WhatsApp notifications</span>
            </label>
          </div>
        </div>

        {/* Privacy Settings */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Privacy Settings</h3>
          <div className="space-y-3">
            <label className="flex items-center">
              <input type="checkbox" defaultChecked className="mr-3" />
              <span className="text-sm text-gray-700">Make profile visible to family members</span>
            </label>
            <label className="flex items-center">
              <input type="checkbox" className="mr-3" />
              <span className="text-sm text-gray-700">Allow family members to add me to events</span>
            </label>
            <label className="flex items-center">
              <input type="checkbox" defaultChecked className="mr-3" />
              <span className="text-sm text-gray-700">Show my birthday in family calendar</span>
            </label>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProfilePage;
