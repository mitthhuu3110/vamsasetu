import React, { useState } from 'react';
import { useEvents } from '../../hooks/useEvents';
import EventCard from './EventCard';
import LoadingSpinner from '../common/LoadingSpinner';
import EmptyState from '../common/EmptyState';
import type { Event } from '../../types/event';

interface EventListProps {
  onEdit?: (event: Event) => void;
  onDelete?: (event: Event) => void;
}

const EventList: React.FC<EventListProps> = ({ onEdit, onDelete }) => {
  const [typeFilter, setTypeFilter] = useState<string>('');
  const { data: eventsResponse, isLoading } = useEvents({
    type: typeFilter || undefined,
  });

  const events = eventsResponse?.data?.events || [];

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-wrap gap-2">
        <button
          onClick={() => setTypeFilter('')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            typeFilter === ''
              ? 'bg-saffron text-white'
              : 'bg-white text-charcoal border border-gray-300 hover:bg-gray-50'
          }`}
        >
          All Events
        </button>
        <button
          onClick={() => setTypeFilter('birthday')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            typeFilter === 'birthday'
              ? 'bg-saffron text-white'
              : 'bg-white text-charcoal border border-gray-300 hover:bg-gray-50'
          }`}
        >
          🎂 Birthdays
        </button>
        <button
          onClick={() => setTypeFilter('anniversary')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            typeFilter === 'anniversary'
              ? 'bg-saffron text-white'
              : 'bg-white text-charcoal border border-gray-300 hover:bg-gray-50'
          }`}
        >
          💍 Anniversaries
        </button>
        <button
          onClick={() => setTypeFilter('ceremony')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            typeFilter === 'ceremony'
              ? 'bg-saffron text-white'
              : 'bg-white text-charcoal border border-gray-300 hover:bg-gray-50'
          }`}
        >
          🎉 Ceremonies
        </button>
        <button
          onClick={() => setTypeFilter('custom')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            typeFilter === 'custom'
              ? 'bg-saffron text-white'
              : 'bg-white text-charcoal border border-gray-300 hover:bg-gray-50'
          }`}
        >
          📌 Custom
        </button>
      </div>

      {/* Event List */}
      {isLoading ? (
        <LoadingSpinner />
      ) : events.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {events.map((event) => (
            <EventCard
              key={event.id}
              event={event}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          ))}
        </div>
      ) : (
        <EmptyState
          icon="📅"
          title={typeFilter ? `No ${typeFilter} events` : 'No events yet'}
          description={
            typeFilter
              ? `No ${typeFilter} events found. Try a different filter.`
              : 'Add events to keep track of important family dates'
          }
        />
      )}
    </div>
  );
};

export default EventList;
