import React, { useState } from 'react';
import EventList from '../components/events/EventList';
import AddEventModal from '../components/events/AddEventModal';
import Button from '../components/ui/Button';
import { useDeleteEvent } from '../hooks/useEvents';
import type { Event } from '../types/event';

const EventsPage: React.FC = () => {
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [viewMode, setViewMode] = useState<'list' | 'calendar'>('list');
  const { mutate: deleteEvent } = useDeleteEvent();

  const handleEdit = (event: Event) => {
    // TODO: Open edit modal
    console.log('Edit event:', event);
  };

  const handleDelete = (event: Event) => {
    if (window.confirm(`Are you sure you want to delete "${event.title}"?`)) {
      deleteEvent(event.id);
    }
  };

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-3xl font-bold text-charcoal">
            Family Events
          </h1>
          <p className="text-charcoal/70 mt-1">
            Track birthdays, anniversaries, and special occasions
          </p>
        </div>
        <div className="flex items-center space-x-3">
          {/* View Toggle */}
          <div className="flex bg-white border border-gray-300 rounded-lg overflow-hidden">
            <button
              onClick={() => setViewMode('list')}
              className={`px-4 py-2 text-sm font-medium transition-colors ${
                viewMode === 'list'
                  ? 'bg-saffron text-white'
                  : 'text-charcoal hover:bg-gray-50'
              }`}
            >
              List
            </button>
            <button
              onClick={() => setViewMode('calendar')}
              className={`px-4 py-2 text-sm font-medium transition-colors ${
                viewMode === 'calendar'
                  ? 'bg-saffron text-white'
                  : 'text-charcoal hover:bg-gray-50'
              }`}
            >
              Calendar
            </button>
          </div>
          <Button
            variant="primary"
            onClick={() => setIsAddModalOpen(true)}
          >
            + Add Event
          </Button>
        </div>
      </div>

      {/* Content */}
      {viewMode === 'list' ? (
        <EventList onEdit={handleEdit} onDelete={handleDelete} />
      ) : (
        <div className="bg-white rounded-lg p-8 text-center border border-gray-200">
          <p className="text-charcoal/60">Calendar view coming soon...</p>
        </div>
      )}

      {/* Add Event Modal */}
      <AddEventModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
      />
    </div>
  );
};

export default EventsPage;
