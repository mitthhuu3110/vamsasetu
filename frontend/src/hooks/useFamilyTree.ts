import { useQuery } from '@tanstack/react-query';
import familyTreeService from '../services/familyTreeService';
import type { FamilyTree } from '../types/familyTree';
import type { APIResponse } from '../services/api';

/**
 * Hook for fetching the complete family tree structure
 * Returns nodes and edges in React Flow format
 */
export function useFamilyTree() {
  return useQuery<APIResponse<FamilyTree>, Error>({
    queryKey: ['familyTree'],
    queryFn: familyTreeService.getTree,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
