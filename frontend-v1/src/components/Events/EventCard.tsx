import React from 'react';
import { Event, EventType } from '../../types';
import { 
  CalendarDaysIcon, 
  MapPinIcon, 
  ClockIcon,
  PencilIcon,
  TrashIcon
} from '@heroicons/react/24/outline';
import { format, isToday, isTomorrow, isYesterday } from 'date-fns';

interface EventCardProps {
  event: Event;
  onEdit: () => void;
  onDelete: () => void;
}

const EventCard: React.FC<EventCardProps> = ({ event, onEdit, onDelete }) => {
  const eventDate = new Date(event.date);
  const isUpcoming = eventDate >= new Date();

  const getEventTypeColor = (type: EventType): string => {
    const colors = {
      [EventType.BIRTHDAY]: 'bg-green-100 text-green-800',
      [EventType.ANNIVERSARY]: 'bg-pink-100 text-pink-800',
      [EventType.WEDDING]: 'bg-purple-100 text-purple-800',
      [EventType.FUNERAL]: 'bg-gray-100 text-gray-800',
      [EventType.FESTIVAL]: 'bg-orange-100 text-orange-800',
      [EventType.CUSTOM]: 'bg-blue-100 text-blue-800',
    };
    return colors[type] || colors[EventType.CUSTOM];
  };

  const getEventTypeIcon = (type: EventType) => {
    // You can customize icons for different event types
    return CalendarDaysIcon;
  };

  const formatEventDate = (date: Date): string => {
    if (isToday(date)) return 'Today';
    if (isTomorrow(date)) return 'Tomorrow';
    if (isYesterday(date)) return 'Yesterday';
    return format(date, 'MMM dd, yyyy');
  };

  const EventIcon = getEventTypeIcon(event.eventType);

  return (
    <div className={`card hover:shadow-md transition-shadow ${
      isUpcoming ? 'border-l-4 border-l-primary-500' : 'opacity-75'
    }`}>
      {/* Header */}
      <div className="flex justify-between items-start mb-3">
        <div className="flex items-center space-x-2">
          <EventIcon className="w-5 h-5 text-primary-600" />
          <span className={`px-2 py-1 rounded-full text-xs font-medium ${getEventTypeColor(event.eventType)}`}>
            {event.eventType}
          </span>
        </div>
        <div className="flex space-x-1">
          <button
            onClick={onEdit}
            className="p-1 text-gray-400 hover:text-gray-600"
          >
            <PencilIcon className="w-4 h-4" />
          </button>
          <button
            onClick={onDelete}
            className="p-1 text-gray-400 hover:text-red-600"
          >
            <TrashIcon className="w-4 h-4" />
          </button>
        </div>
      </div>

      {/* Title and Description */}
      <h3 className="font-semibold text-gray-900 mb-2">{event.title}</h3>
      {event.description && (
        <p className="text-sm text-gray-600 mb-3 line-clamp-2">{event.description}</p>
      )}

      {/* Date and Time */}
      <div className="flex items-center space-x-4 text-sm text-gray-500 mb-3">
        <div className="flex items-center space-x-1">
          <CalendarDaysIcon className="w-4 h-4" />
          <span>{formatEventDate(eventDate)}</span>
        </div>
        {event.time && (
          <div className="flex items-center space-x-1">
            <ClockIcon className="w-4 h-4" />
            <span>{event.time}</span>
          </div>
        )}
      </div>

      {/* Location */}
      {event.location && (
        <div className="flex items-center space-x-1 text-sm text-gray-500 mb-3">
          <MapPinIcon className="w-4 h-4" />
          <span className="truncate">{event.location}</span>
        </div>
      )}

      {/* Attendees */}
      {event.attendees.length > 0 && (
        <div className="mb-3">
          <p className="text-sm font-medium text-gray-700 mb-1">Attendees ({event.attendees.length})</p>
          <div className="flex flex-wrap gap-1">
            {event.attendees.slice(0, 3).map((attendee, index) => (
              <span
                key={index}
                className="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded-full"
              >
                {attendee.memberName}
              </span>
            ))}
            {event.attendees.length > 3 && (
              <span className="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded-full">
                +{event.attendees.length - 3} more
              </span>
            )}
          </div>
        </div>
      )}

      {/* Recurring Indicator */}
      {event.isRecurring && (
        <div className="flex items-center space-x-1 text-xs text-primary-600">
          <span>🔄</span>
          <span>Recurring</span>
        </div>
      )}

      {/* Reminder Settings */}
      {event.reminderSettings.enabled && (
        <div className="mt-3 pt-3 border-t border-gray-200">
          <div className="flex items-center justify-between text-xs text-gray-500">
            <span>Reminders enabled</span>
            <span>{event.reminderSettings.advanceTime}h before</span>
          </div>
        </div>
      )}
    </div>
  );
};

export default EventCard;
