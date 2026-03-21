/**
 * Example usage of EventService
 * 
 * This file demonstrates how to use the EventService for event management operations.
 */

import eventService from './eventService';
import type { CreateEventRequest } from '../types/event';

/**
 * Example: Get all events with pagination
 */
async function getAllEventsExample() {
  const response = await eventService.getAll({
    page: 1,
    limit: 10,
  });

  if (response.success && response.data) {
    console.log('Events:', response.data.events);
    console.log('Total:', response.data.total);
  } else {
    console.error('Error:', response.error);
  }
}

/**
 * Example: Get events filtered by type
 */
async function getEventsByTypeExample() {
  const response = await eventService.getAll({
    type: 'birthday',
    page: 1,
    limit: 20,
  });

  if (response.success && response.data) {
    console.log('Birthday events:', response.data.events);
  }
}

/**
 * Example: Get events filtered by member
 */
async function getEventsByMemberExample() {
  const memberUUID = '123e4567-e89b-12d3-a456-426614174000';
  
  const response = await eventService.getAll({
    member: memberUUID,
  });

  if (response.success && response.data) {
    console.log('Member events:', response.data.events);
  }
}

/**
 * Example: Get events filtered by date range
 */
async function getEventsByDateRangeExample() {
  const response = await eventService.getAll({
    startDate: '2024-01-01T00:00:00Z',
    endDate: '2024-12-31T23:59:59Z',
  });

  if (response.success && response.data) {
    console.log('Events in 2024:', response.data.events);
  }
}

/**
 * Example: Get a single event by ID
 */
async function getEventByIdExample() {
  const eventId = 1;
  
  const response = await eventService.getById(eventId);

  if (response.success && response.data) {
    console.log('Event:', response.data);
  } else {
    console.error('Error:', response.error);
  }
}

/**
 * Example: Create a new event
 */
async function createEventExample() {
  const newEvent: CreateEventRequest = {
    title: 'Birthday Celebration',
    description: 'Celebrating John\'s 30th birthday',
    eventDate: '2024-06-15T18:00:00Z',
    eventType: 'birthday',
    memberIds: ['123e4567-e89b-12d3-a456-426614174000'],
  };

  const response = await eventService.create(newEvent);

  if (response.success && response.data) {
    console.log('Created event:', response.data);
  } else {
    console.error('Error:', response.error);
  }
}

/**
 * Example: Update an existing event
 */
async function updateEventExample() {
  const eventId = 1;
  
  const updates = {
    title: 'Updated Birthday Celebration',
    description: 'Updated description',
  };

  const response = await eventService.update(eventId, updates);

  if (response.success && response.data) {
    console.log('Updated event:', response.data);
  } else {
    console.error('Error:', response.error);
  }
}

/**
 * Example: Delete an event
 */
async function deleteEventExample() {
  const eventId = 1;
  
  const response = await eventService.delete(eventId);

  if (response.success && response.data) {
    console.log('Success:', response.data.message);
  } else {
    console.error('Error:', response.error);
  }
}

/**
 * Example: Get upcoming events
 */
async function getUpcomingEventsExample() {
  // Get events in the next 30 days (default)
  const response = await eventService.getUpcoming();

  if (response.success && response.data) {
    console.log('Upcoming events:', response.data);
  } else {
    console.error('Error:', response.error);
  }
}

/**
 * Example: Get upcoming events with custom days
 */
async function getUpcomingEventsCustomDaysExample() {
  // Get events in the next 7 days
  const response = await eventService.getUpcoming(7);

  if (response.success && response.data) {
    console.log('Events in next 7 days:', response.data);
  }
}

// Export examples for reference
export {
  getAllEventsExample,
  getEventsByTypeExample,
  getEventsByMemberExample,
  getEventsByDateRangeExample,
  getEventByIdExample,
  createEventExample,
  updateEventExample,
  deleteEventExample,
  getUpcomingEventsExample,
  getUpcomingEventsCustomDaysExample,
};
