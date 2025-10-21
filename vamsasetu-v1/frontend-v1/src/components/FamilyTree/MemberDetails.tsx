import React from 'react';
import { FamilyMember } from '../../types';
import { 
  UserIcon, 
  CalendarIcon, 
  PhoneIcon, 
  EnvelopeIcon,
  MapPinIcon,
  BriefcaseIcon,
  XMarkIcon
} from '@heroicons/react/24/outline';

interface MemberDetailsProps {
  member: FamilyMember;
  onEdit: () => void;
  onClose: () => void;
}

const MemberDetails: React.FC<MemberDetailsProps> = ({ member, onEdit, onClose }) => {
  const formatDate = (dateString?: string) => {
    if (!dateString) return 'Not specified';
    return new Date(dateString).toLocaleDateString();
  };

  return (
    <div className="card">
      <div className="flex justify-between items-start mb-4">
        <h3 className="text-lg font-semibold text-gray-900">Member Details</h3>
        <button
          onClick={onClose}
          className="text-gray-400 hover:text-gray-600"
        >
          <XMarkIcon className="w-5 h-5" />
        </button>
      </div>

      {/* Profile Section */}
      <div className="text-center mb-6">
        <div className="w-20 h-20 mx-auto mb-3 bg-primary-100 rounded-full flex items-center justify-center">
          {member.profilePicture ? (
            <img
              src={member.profilePicture}
              alt={`${member.firstName} ${member.lastName}`}
              className="w-20 h-20 rounded-full object-cover"
            />
          ) : (
            <span className="text-primary-600 font-semibold text-xl">
              {member.firstName[0]}{member.lastName[0]}
            </span>
          )}
        </div>
        <h4 className="text-xl font-semibold text-gray-900">
          {member.firstName} {member.middleName} {member.lastName}
        </h4>
        <p className="text-gray-500">
          {member.gender} • {member.isAlive ? 'Alive' : 'Deceased'}
        </p>
      </div>

      {/* Details */}
      <div className="space-y-4">
        {member.dateOfBirth && (
          <div className="flex items-center space-x-3">
            <CalendarIcon className="w-5 h-5 text-gray-400" />
            <div>
              <p className="text-sm font-medium text-gray-900">Date of Birth</p>
              <p className="text-sm text-gray-500">{formatDate(member.dateOfBirth)}</p>
            </div>
          </div>
        )}

        {member.dateOfDeath && (
          <div className="flex items-center space-x-3">
            <CalendarIcon className="w-5 h-5 text-gray-400" />
            <div>
              <p className="text-sm font-medium text-gray-900">Date of Death</p>
              <p className="text-sm text-gray-500">{formatDate(member.dateOfDeath)}</p>
            </div>
          </div>
        )}

        {member.phone && (
          <div className="flex items-center space-x-3">
            <PhoneIcon className="w-5 h-5 text-gray-400" />
            <div>
              <p className="text-sm font-medium text-gray-900">Phone</p>
              <p className="text-sm text-gray-500">{member.phone}</p>
            </div>
          </div>
        )}

        {member.email && (
          <div className="flex items-center space-x-3">
            <EnvelopeIcon className="w-5 h-5 text-gray-400" />
            <div>
              <p className="text-sm font-medium text-gray-900">Email</p>
              <p className="text-sm text-gray-500">{member.email}</p>
            </div>
          </div>
        )}

        {member.occupation && (
          <div className="flex items-center space-x-3">
            <BriefcaseIcon className="w-5 h-5 text-gray-400" />
            <div>
              <p className="text-sm font-medium text-gray-900">Occupation</p>
              <p className="text-sm text-gray-500">{member.occupation}</p>
            </div>
          </div>
        )}

        {member.address && (
          <div className="flex items-start space-x-3">
            <MapPinIcon className="w-5 h-5 text-gray-400 mt-0.5" />
            <div>
              <p className="text-sm font-medium text-gray-900">Address</p>
              <p className="text-sm text-gray-500">
                {member.address.street && `${member.address.street}, `}
                {member.address.city && `${member.address.city}, `}
                {member.address.state && `${member.address.state}, `}
                {member.address.country}
              </p>
            </div>
          </div>
        )}

        {member.notes && (
          <div className="pt-4 border-t border-gray-200">
            <p className="text-sm font-medium text-gray-900 mb-2">Notes</p>
            <p className="text-sm text-gray-500">{member.notes}</p>
          </div>
        )}
      </div>

      {/* Actions */}
      <div className="mt-6 pt-4 border-t border-gray-200">
        <button
          onClick={onEdit}
          className="w-full btn-primary"
        >
          Edit Member
        </button>
      </div>
    </div>
  );
};

export default MemberDetails;
