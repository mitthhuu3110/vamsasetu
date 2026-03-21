import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import eventService from '../services/eventService';
import type { Event, CreateEventRequest } from '../types/event';
import type { APIResponse } from '../services/api';

/**
 * Hook for fetching all events with optional filters
 */
export function useEvents(params?: {
  page?: number;
  limit?: number;
  type?: string;
  member?: string;
  startDate?: string;
  endDate?: string;
}) {
  return useQuery<APIResponse<{ events: Event[]; total: number; page: number; limit: number }>, Error>({
    queryKey: ['events', params],
    queryFn: () => eventService.getAll(params),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}

/**
 * Hook for fetching a single event by ID
 */
export function useEvent(id: number) {
  return useQuery<APIResponse<Event>, Error>({
    queryKey: ['events', id],
    queryFn: () => eventService.getById(id),
    enabled: !!id,
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}

/**
 * Hook for creating a new event
 * Invalidates events cache on success
 */
export function useCreateEvent() {
  const queryClient = useQueryClient();

  return useMutation<APIResponse<Event>, Error, CreateEventRequest>({
    mutationFn: eventService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['events'] });
    },
  });
}

/**
 * Hook for updating an existing event
 * Invalidates events cache on success
 */
export function useUpdateEvent() {
  const queryClient = useQueryClient();

  return useMutation<APIResponse<Event>, Error, { id: number; data: Partial<CreateEventRequest> }>({
    mutationFn: ({ id, data }) => eventService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['events'] });
      queryClient.invalidateQueries({ queryKey: ['events', variables.id] });
    },
  });
}

/**
 * Hook for deleting an event
 * Invalidates events cache on success
 */
export function useDeleteEvent() {
  const queryClient = useQueryClient();

  return useMutation<APIResponse<{ message: string }>, Error, number>({
    mutationFn: eventService.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['events'] });
    },
  });
}

/**
 * Hook for fetching upcoming events
 */
export function useUpcomingEvents(days?: number) {
  return useQuery<APIResponse<Event[]>, Error>({
    queryKey: ['events', 'upcoming', days],
    queryFn: () => eventService.getUpcoming(days),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}
