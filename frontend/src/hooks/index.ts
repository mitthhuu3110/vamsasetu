/**
 * React Query hooks for VamsaSetu
 * 
 * This module exports all custom hooks for data fetching and mutations.
 * All hooks use @tanstack/react-query for server state management.
 */

// Auth hooks
export { useLogin, useRegister, useProfile } from './useAuth';

// Member hooks
export {
  useMembers,
  useMember,
  useCreateMember,
  useUpdateMember,
  useDeleteMember,
  useSearchMembers,
} from './useMembers';

// Relationship hooks
export {
  useRelationships,
  useCreateRelationship,
  useDeleteRelationship,
  useFindPath,
} from './useRelationships';

// Event hooks
export {
  useEvents,
  useEvent,
  useCreateEvent,
  useUpdateEvent,
  useDeleteEvent,
  useUpcomingEvents,
} from './useEvents';

// Family tree hooks
export { useFamilyTree } from './useFamilyTree';

// WebSocket hooks
export { useWebSocket } from './useWebSocket';
