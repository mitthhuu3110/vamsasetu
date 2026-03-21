import api, { type APIResponse } from './api';
import type { Event, CreateEventRequest } from '../types/event';

/**
 * Event service for handling event CRUD operations and filtering
 */
class EventService {
  /**
   * Get all events with optional pagination and filters
   * @param params - Query parameters (page, limit, type, member, startDate, endDate)
   * @returns Promise with paginated event list
   */
  async getAll(params?: {
    page?: number;
    limit?: number;
    type?: string;
    member?: string;
    startDate?: string;
    endDate?: string;
  }): Promise<APIResponse<{ events: Event[]; total: number; page: number; limit: number }>> {
    try {
      const queryParams = new URLSearchParams();
      
      if (params?.page) queryParams.append('page', params.page.toString());
      if (params?.limit) queryParams.append('limit', params.limit.toString());
      if (params?.type) queryParams.append('type', params.type);
      if (params?.member) queryParams.append('member', params.member);
      if (params?.startDate) queryParams.append('startDate', params.startDate);
      if (params?.endDate) queryParams.append('endDate', params.endDate);
      
      const queryString = queryParams.toString();
      const url = queryString ? `/api/events?${queryString}` : '/api/events';
      
      const response = await api.get<APIResponse<{ events: Event[]; total: number; page: number; limit: number }>>(url);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to fetch events',
      };
    }
  }

  /**
   * Get an event by ID
   * @param id - Event ID
   * @returns Promise with event data
   */
  async getById(id: number): Promise<APIResponse<Event>> {
    try {
      const response = await api.get<APIResponse<Event>>(`/api/events/${id}`);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to fetch event',
      };
    }
  }

  /**
   * Create a new event
   * @param data - Event creation data
   * @returns Promise with created event
   */
  async create(data: CreateEventRequest): Promise<APIResponse<Event>> {
    try {
      const response = await api.post<APIResponse<Event>>('/api/events', data);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to create event',
      };
    }
  }

  /**
   * Update an existing event
   * @param id - Event ID
   * @param data - Event update data
   * @returns Promise with updated event
   */
  async update(id: number, data: Partial<CreateEventRequest>): Promise<APIResponse<Event>> {
    try {
      const response = await api.put<APIResponse<Event>>(`/api/events/${id}`, data);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to update event',
      };
    }
  }

  /**
   * Delete an event
   * @param id - Event ID
   * @returns Promise with success message
   */
  async delete(id: number): Promise<APIResponse<{ message: string }>> {
    try {
      const response = await api.delete<APIResponse<{ message: string }>>(`/api/events/${id}`);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to delete event',
      };
    }
  }

  /**
   * Get upcoming events within the next N days
   * @param days - Number of days to look ahead (default: 30)
   * @returns Promise with upcoming events
   */
  async getUpcoming(days?: number): Promise<APIResponse<Event[]>> {
    try {
      const queryParams = new URLSearchParams();
      if (days) queryParams.append('days', days.toString());
      
      const queryString = queryParams.toString();
      const url = queryString ? `/api/events/upcoming?${queryString}` : '/api/events/upcoming';
      
      const response = await api.get<APIResponse<Event[]>>(url);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to fetch upcoming events',
      };
    }
  }
}

// Export singleton instance
export default new EventService();
