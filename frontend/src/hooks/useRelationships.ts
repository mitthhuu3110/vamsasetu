import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import relationshipService from '../services/relationshipService';
import type { Relationship, CreateRelationshipRequest, RelationshipPath } from '../types/relationship';
import type { APIResponse } from '../services/api';

/**
 * Hook for fetching all relationships
 */
export function useRelationships() {
  return useQuery<APIResponse<Relationship[]>, Error>({
    queryKey: ['relationships'],
    queryFn: relationshipService.getAll,
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}

/**
 * Hook for creating a new relationship
 * Invalidates relationships and familyTree cache on success
 */
export function useCreateRelationship() {
  const queryClient = useQueryClient();

  return useMutation<APIResponse<Relationship>, Error, CreateRelationshipRequest>({
    mutationFn: relationshipService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['relationships'] });
      queryClient.invalidateQueries({ queryKey: ['familyTree'] });
    },
  });
}

/**
 * Hook for deleting a relationship
 * Invalidates relationships and familyTree cache on success
 */
export function useDeleteRelationship() {
  const queryClient = useQueryClient();

  return useMutation<
    APIResponse<{ message: string }>,
    Error,
    { id: string; fromId: string; toId: string; type: string }
  >({
    mutationFn: ({ id, fromId, toId, type }) => relationshipService.delete(id, fromId, toId, type),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['relationships'] });
      queryClient.invalidateQueries({ queryKey: ['familyTree'] });
    },
  });
}

/**
 * Hook for finding relationship path between two members
 */
export function useFindPath(fromId: string, toId: string) {
  return useQuery<APIResponse<RelationshipPath>, Error>({
    queryKey: ['relationships', 'path', fromId, toId],
    queryFn: () => relationshipService.findPath(fromId, toId),
    enabled: !!fromId && !!toId && fromId !== toId,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
