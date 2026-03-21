import api, { type APIResponse } from './api';
import type { Member, CreateMemberRequest } from '../types/member';

/**
 * Member service for handling member CRUD operations and search
 */
class MemberService {
  /**
   * Get all members with optional pagination and filters
   * @param params - Query parameters (page, limit, search, gender)
   * @returns Promise with paginated member list
   */
  async getAll(params?: {
    page?: number;
    limit?: number;
    search?: string;
    gender?: string;
  }): Promise<APIResponse<{ members: Member[]; total: number; page: number; limit: number }>> {
    try {
      const queryParams = new URLSearchParams();
      
      if (params?.page) queryParams.append('page', params.page.toString());
      if (params?.limit) queryParams.append('limit', params.limit.toString());
      if (params?.search) queryParams.append('search', params.search);
      if (params?.gender) queryParams.append('gender', params.gender);
      
      const queryString = queryParams.toString();
      const url = queryString ? `/api/members?${queryString}` : '/api/members';
      
      const response = await api.get<APIResponse<{ members: Member[]; total: number; page: number; limit: number }>>(url);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to fetch members',
      };
    }
  }

  /**
   * Get a member by ID
   * @param id - Member ID
   * @returns Promise with member data
   */
  async getById(id: string): Promise<APIResponse<Member>> {
    try {
      const response = await api.get<APIResponse<Member>>(`/api/members/${id}`);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to fetch member',
      };
    }
  }

  /**
   * Create a new member
   * @param data - Member creation data
   * @returns Promise with created member
   */
  async create(data: CreateMemberRequest): Promise<APIResponse<Member>> {
    try {
      const response = await api.post<APIResponse<Member>>('/api/members', data);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to create member',
      };
    }
  }

  /**
   * Update an existing member
   * @param id - Member ID
   * @param data - Member update data
   * @returns Promise with updated member
   */
  async update(id: string, data: Partial<CreateMemberRequest>): Promise<APIResponse<Member>> {
    try {
      const response = await api.put<APIResponse<Member>>(`/api/members/${id}`, data);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to update member',
      };
    }
  }

  /**
   * Delete a member (soft delete)
   * @param id - Member ID
   * @returns Promise with success message
   */
  async delete(id: string): Promise<APIResponse<{ message: string }>> {
    try {
      const response = await api.delete<APIResponse<{ message: string }>>(`/api/members/${id}`);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to delete member',
      };
    }
  }

  /**
   * Search members by name
   * @param query - Search query string
   * @returns Promise with matching members
   */
  async search(query: string): Promise<APIResponse<{ members: Member[]; total: number; page: number; limit: number }>> {
    try {
      const response = await api.get<APIResponse<{ members: Member[]; total: number; page: number; limit: number }>>(
        `/api/members?search=${encodeURIComponent(query)}`
      );
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to search members',
      };
    }
  }
}

// Export singleton instance
export default new MemberService();
