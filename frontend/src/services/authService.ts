import api, { type APIResponse } from './api';
import type { User, RegisterRequest, LoginRequest, AuthResponse } from '../types/user';

/**
 * Authentication service for handling user registration, login, token refresh, and profile retrieval
 */
class AuthService {
  /**
   * Register a new user
   * @param data - Registration data (email, password, name)
   * @returns Promise with auth response containing user, accessToken, and refreshToken
   */
  async register(data: RegisterRequest): Promise<APIResponse<AuthResponse>> {
    try {
      const response = await api.post<APIResponse<AuthResponse>>('/auth/register', data);
      
      // Store tokens if registration successful
      if (response.data.success && response.data.data) {
        localStorage.setItem('accessToken', response.data.data.accessToken);
        localStorage.setItem('refreshToken', response.data.data.refreshToken);
      }
      
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Registration failed',
      };
    }
  }

  /**
   * Login an existing user
   * @param credentials - Login credentials (email, password)
   * @returns Promise with auth response containing user, accessToken, and refreshToken
   */
  async login(credentials: LoginRequest): Promise<APIResponse<AuthResponse>> {
    try {
      const response = await api.post<APIResponse<AuthResponse>>('/auth/login', credentials);
      
      // Store tokens if login successful
      if (response.data.success && response.data.data) {
        localStorage.setItem('accessToken', response.data.data.accessToken);
        localStorage.setItem('refreshToken', response.data.data.refreshToken);
      }
      
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Login failed',
      };
    }
  }

  /**
   * Refresh access token using refresh token
   * @param refreshToken - The refresh token
   * @returns Promise with auth response containing new accessToken and refreshToken
   */
  async refresh(refreshToken: string): Promise<APIResponse<AuthResponse>> {
    try {
      const response = await api.post<APIResponse<AuthResponse>>('/auth/refresh', {
        refreshToken,
      });
      
      // Update tokens if refresh successful
      if (response.data.success && response.data.data) {
        localStorage.setItem('accessToken', response.data.data.accessToken);
        localStorage.setItem('refreshToken', response.data.data.refreshToken);
      }
      
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Token refresh failed',
      };
    }
  }

  /**
   * Get current user profile
   * @returns Promise with user profile data
   */
  async getProfile(): Promise<APIResponse<User>> {
    try {
      const response = await api.get<APIResponse<User>>('/auth/profile');
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to fetch profile',
      };
    }
  }

  /**
   * Logout user by clearing tokens
   */
  logout(): void {
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
  }

  /**
   * Check if user is authenticated
   * @returns true if access token exists
   */
  isAuthenticated(): boolean {
    return !!localStorage.getItem('accessToken');
  }

  /**
   * Get stored access token
   * @returns access token or null
   */
  getAccessToken(): string | null {
    return localStorage.getItem('accessToken');
  }

  /**
   * Get stored refresh token
   * @returns refresh token or null
   */
  getRefreshToken(): string | null {
    return localStorage.getItem('refreshToken');
  }
}

// Export singleton instance
export default new AuthService();
