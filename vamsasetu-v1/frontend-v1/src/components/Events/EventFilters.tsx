import React from 'react';
import { EventType } from '../../types';
import { MagnifyingGlassIcon } from '@heroicons/react/24/outline';

interface EventFiltersProps {
  selectedType: EventType | 'ALL';
  onTypeChange: (type: EventType | 'ALL') => void;
  searchQuery: string;
  onSearchChange: (query: string) => void;
}

const EventFilters: React.FC<EventFiltersProps> = ({
  selectedType,
  onTypeChange,
  searchQuery,
  onSearchChange,
}) => {
  const eventTypes = [
    { value: 'ALL', label: 'All Events' },
    { value: EventType.BIRTHDAY, label: 'Birthdays' },
    { value: EventType.ANNIVERSARY, label: 'Anniversaries' },
    { value: EventType.WEDDING, label: 'Weddings' },
    { value: EventType.FESTIVAL, label: 'Festivals' },
    { value: EventType.FUNERAL, label: 'Funerals' },
    { value: EventType.CUSTOM, label: 'Custom' },
  ];

  return (
    <div className="card">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between space-y-4 md:space-y-0">
        {/* Search */}
        <div className="flex-1 max-w-md">
          <div className="relative">
            <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search events..."
              value={searchQuery}
              onChange={(e) => onSearchChange(e.target.value)}
              className="input-field pl-10"
            />
          </div>
        </div>

        {/* Event Type Filter */}
        <div className="flex flex-wrap gap-2">
          {eventTypes.map((type) => (
            <button
              key={type.value}
              onClick={() => onTypeChange(type.value as EventType | 'ALL')}
              className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                selectedType === type.value
                  ? 'bg-primary-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {type.label}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
};

export default EventFilters;
