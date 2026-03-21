import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import memberService from '../services/memberService';
import type { Member, CreateMemberRequest } from '../types/member';
import type { APIResponse } from '../services/api';

/**
 * Hook for fetching all members with optional filters
 */
export function useMembers(params?: {
  page?: number;
  limit?: number;
  search?: string;
  gender?: string;
}) {
  return useQuery<APIResponse<{ members: Member[]; total: number; page: number; limit: number }>, Error>({
    queryKey: ['members', params],
    queryFn: () => memberService.getAll(params),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}

/**
 * Hook for fetching a single member by ID
 */
export function useMember(id: string) {
  return useQuery<APIResponse<Member>, Error>({
    queryKey: ['members', id],
    queryFn: () => memberService.getById(id),
    enabled: !!id,
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}

/**
 * Hook for creating a new member
 * Invalidates members and familyTree cache on success
 */
export function useCreateMember() {
  const queryClient = useQueryClient();

  return useMutation<APIResponse<Member>, Error, CreateMemberRequest>({
    mutationFn: memberService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members'] });
      queryClient.invalidateQueries({ queryKey: ['familyTree'] });
    },
  });
}

/**
 * Hook for updating an existing member
 * Invalidates members and familyTree cache on success
 */
export function useUpdateMember() {
  const queryClient = useQueryClient();

  return useMutation<APIResponse<Member>, Error, { id: string; data: Partial<CreateMemberRequest> }>({
    mutationFn: ({ id, data }) => memberService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['members'] });
      queryClient.invalidateQueries({ queryKey: ['members', variables.id] });
      queryClient.invalidateQueries({ queryKey: ['familyTree'] });
    },
  });
}

/**
 * Hook for deleting a member (soft delete)
 * Invalidates members and familyTree cache on success
 */
export function useDeleteMember() {
  const queryClient = useQueryClient();

  return useMutation<APIResponse<{ message: string }>, Error, string>({
    mutationFn: memberService.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['members'] });
      queryClient.invalidateQueries({ queryKey: ['familyTree'] });
    },
  });
}

/**
 * Hook for searching members by name
 */
export function useSearchMembers(query: string) {
  return useQuery<APIResponse<{ members: Member[]; total: number; page: number; limit: number }>, Error>({
    queryKey: ['members', 'search', query],
    queryFn: () => memberService.search(query),
    enabled: query.length > 0,
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}
