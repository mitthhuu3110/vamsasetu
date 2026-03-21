// Event type enum
export type EventType = 'birthday' | 'anniversary' | 'ceremony' | 'custom';

export interface Event {
  id: number;
  title: string;
  description: string;
  eventDate: string;
  eventType: EventType;
  memberIds: string[];
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateEventData {
  title: string;
  description?: string;
  eventDate: string;
  eventType: EventType;
  memberIds: string[];
}

// Alias for task requirement
export type CreateEventRequest = CreateEventData;

export interface UpdateEventData extends Partial<CreateEventData> {
  id: number;
}
