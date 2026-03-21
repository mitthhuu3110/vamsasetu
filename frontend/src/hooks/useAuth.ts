import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import authService from '../services/authService';
import { useAuthStore } from '../stores/authStore';
import type { LoginRequest, RegisterRequest, User, AuthResponse } from '../types/user';
import type { APIResponse } from '../services/api';

/**
 * Hook for user login
 * Automatically stores auth data in authStore on success
 */
export function useLogin() {
  const { setAuth } = useAuthStore();
  const queryClient = useQueryClient();

  return useMutation<APIResponse<AuthResponse>, Error, LoginRequest>({
    mutationFn: authService.login,
    onSuccess: (response) => {
      if (response.success && response.data) {
        const { user, accessToken, refreshToken } = response.data;
        setAuth(user, accessToken, refreshToken);
        queryClient.invalidateQueries({ queryKey: ['profile'] });
      }
    },
  });
}

/**
 * Hook for user registration
 * Automatically stores auth data in authStore on success
 */
export function useRegister() {
  const { setAuth } = useAuthStore();
  const queryClient = useQueryClient();

  return useMutation<APIResponse<AuthResponse>, Error, RegisterRequest>({
    mutationFn: authService.register,
    onSuccess: (response) => {
      if (response.success && response.data) {
        const { user, accessToken, refreshToken } = response.data;
        setAuth(user, accessToken, refreshToken);
        queryClient.invalidateQueries({ queryKey: ['profile'] });
      }
    },
  });
}

/**
 * Hook for fetching user profile
 * Only fetches when user is authenticated
 */
export function useProfile() {
  const { isAuthenticated } = useAuthStore();

  return useQuery<APIResponse<User>, Error>({
    queryKey: ['profile'],
    queryFn: authService.getProfile,
    enabled: isAuthenticated,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
