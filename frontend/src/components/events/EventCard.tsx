import React from 'react';
import Card, { CardContent } from '../ui/Card';
import Button from '../ui/Button';
import type { Event } from '../../types/event';

interface EventCardProps {
  event: Event;
  onEdit?: (event: Event) => void;
  onDelete?: (event: Event) => void;
  showActions?: boolean;
}

const EventCard: React.FC<EventCardProps> = ({
  event,
  onEdit,
  onDelete,
  showActions = true,
}) => {
  const eventDate = new Date(event.eventDate);
  const today = new Date();
  const daysUntil = Math.ceil((eventDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
  const isUpcoming = daysUntil >= 0 && daysUntil <= 7;

  const eventTypeIcons: Record<string, string> = {
    birthday: '🎂',
    anniversary: '💍',
    ceremony: '🎉',
    custom: '📌',
  };

  const eventTypeColors: Record<string, string> = {
    birthday: 'bg-rose/10 text-rose border-rose/20',
    anniversary: 'bg-pink-100 text-pink-700 border-pink-200',
    ceremony: 'bg-turmeric/10 text-turmeric border-turmeric/20',
    custom: 'bg-teal/10 text-teal border-teal/20',
  };

  return (
    <Card variant="elevated" className="hover:shadow-xl transition-shadow">
      <CardContent className="p-4">
        <div className="flex items-start justify-between mb-3">
          <div className="flex items-center space-x-3">
            <span className="text-3xl">{eventTypeIcons[event.eventType] || '📌'}</span>
            <div>
              <h3 className="font-semibold text-charcoal text-lg">{event.title}</h3>
              <span
                className={`inline-block text-xs px-2 py-1 rounded-full border ${
                  eventTypeColors[event.eventType] || eventTypeColors.custom
                }`}
              >
                {event.eventType}
              </span>
            </div>
          </div>
          {isUpcoming && (
            <span className="bg-amber text-white text-xs px-2 py-1 rounded-full font-medium animate-pulse">
              {daysUntil === 0 ? 'Today!' : `${daysUntil}d`}
            </span>
          )}
        </div>

        {event.description && (
          <p className="text-charcoal/70 text-sm mb-3">{event.description}</p>
        )}

        <div className="flex items-center justify-between text-sm">
          <div className="text-charcoal/60">
            📅 {eventDate.toLocaleDateString('en-US', {
              weekday: 'short',
              year: 'numeric',
              month: 'short',
              day: 'numeric',
            })}
          </div>
        </div>

        {showActions && (onEdit || onDelete) && (
          <div className="flex space-x-2 mt-4 pt-4 border-t border-gray-200">
            {onEdit && (
              <Button
                variant="outline"
                onClick={() => onEdit(event)}
                className="flex-1 text-sm"
              >
                Edit
              </Button>
            )}
            {onDelete && (
              <Button
                variant="outline"
                onClick={() => onDelete(event)}
                className="flex-1 text-sm text-rose hover:bg-rose/10"
              >
                Delete
              </Button>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default EventCard;
