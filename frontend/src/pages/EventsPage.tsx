import React, { useState } from 'react';
import { useQuery } from 'react-query';
import { eventsApi } from '../services/api';
import { Event, EventType } from '../types';
import EventCard from '../components/Events/EventCard';
import AddEventModal from '../components/Events/AddEventModal';
import EventFilters from '../components/Events/EventFilters';
import { PlusIcon, CalendarDaysIcon } from '@heroicons/react/24/outline';

const EventsPage: React.FC = () => {
  const [showAddEvent, setShowAddEvent] = useState(false);
  const [selectedEventType, setSelectedEventType] = useState<EventType | 'ALL'>('ALL');
  const [searchQuery, setSearchQuery] = useState('');

  const { data: events, isLoading, error, refetch } = useQuery(
    'events',
    () => eventsApi.getEvents(),
    {
      refetchOnWindowFocus: false,
    }
  );

  const filteredEvents = events?.filter(event => {
    const matchesType = selectedEventType === 'ALL' || event.eventType === selectedEventType;
    const matchesSearch = event.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         event.description?.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesType && matchesSearch;
  }) || [];

  const upcomingEvents = filteredEvents.filter(event => 
    new Date(event.date) >= new Date()
  ).sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());

  const pastEvents = filteredEvents.filter(event => 
    new Date(event.date) < new Date()
  ).sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading events...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600">Error loading events. Please try again.</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Family Events</h1>
          <p className="text-gray-600">Manage birthdays, anniversaries, and special occasions</p>
        </div>
        <button
          onClick={() => setShowAddEvent(true)}
          className="btn-primary flex items-center"
        >
          <PlusIcon className="w-5 h-5 mr-2" />
          Add Event
        </button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="card text-center">
          <CalendarDaysIcon className="w-8 h-8 text-primary-600 mx-auto mb-2" />
          <div className="text-2xl font-bold text-gray-900">{upcomingEvents.length}</div>
          <div className="text-sm text-gray-500">Upcoming Events</div>
        </div>
        <div className="card text-center">
          <CalendarDaysIcon className="w-8 h-8 text-green-600 mx-auto mb-2" />
          <div className="text-2xl font-bold text-gray-900">
            {filteredEvents.filter(e => e.eventType === EventType.BIRTHDAY).length}
          </div>
          <div className="text-sm text-gray-500">Birthdays</div>
        </div>
        <div className="card text-center">
          <CalendarDaysIcon className="w-8 h-8 text-pink-600 mx-auto mb-2" />
          <div className="text-2xl font-bold text-gray-900">
            {filteredEvents.filter(e => e.eventType === EventType.ANNIVERSARY).length}
          </div>
          <div className="text-sm text-gray-500">Anniversaries</div>
        </div>
        <div className="card text-center">
          <CalendarDaysIcon className="w-8 h-8 text-purple-600 mx-auto mb-2" />
          <div className="text-2xl font-bold text-gray-900">
            {filteredEvents.filter(e => e.eventType === EventType.FESTIVAL).length}
          </div>
          <div className="text-sm text-gray-500">Festivals</div>
        </div>
      </div>

      {/* Filters */}
      <EventFilters
        selectedType={selectedEventType}
        onTypeChange={setSelectedEventType}
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
      />

      {/* Events List */}
      <div className="space-y-8">
        {/* Upcoming Events */}
        {upcomingEvents.length > 0 && (
          <div>
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Upcoming Events</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {upcomingEvents.map((event) => (
                <EventCard
                  key={event.id}
                  event={event}
                  onEdit={() => {/* TODO: Implement edit */}}
                  onDelete={() => {/* TODO: Implement delete */}}
                />
              ))}
            </div>
          </div>
        )}

        {/* Past Events */}
        {pastEvents.length > 0 && (
          <div>
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Past Events</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {pastEvents.map((event) => (
                <EventCard
                  key={event.id}
                  event={event}
                  onEdit={() => {/* TODO: Implement edit */}}
                  onDelete={() => {/* TODO: Implement delete */}}
                />
              ))}
            </div>
          </div>
        )}

        {/* Empty State */}
        {filteredEvents.length === 0 && (
          <div className="text-center py-12">
            <CalendarDaysIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No events found</h3>
            <p className="text-gray-500 mb-4">
              {searchQuery || selectedEventType !== 'ALL' 
                ? 'Try adjusting your filters or search terms.'
                : 'Get started by adding your first family event.'
              }
            </p>
            <button
              onClick={() => setShowAddEvent(true)}
              className="btn-primary"
            >
              Add Event
            </button>
          </div>
        )}
      </div>

      {/* Add Event Modal */}
      {showAddEvent && (
        <AddEventModal
          onClose={() => setShowAddEvent(false)}
          onSuccess={() => {
            setShowAddEvent(false);
            refetch();
          }}
        />
      )}
    </div>
  );
};

export default EventsPage;
