import api, { type APIResponse } from './api';
import type { Relationship, CreateRelationshipRequest, RelationshipPath } from '../types/relationship';

/**
 * Relationship service for handling relationship CRUD operations and path finding
 */
class RelationshipService {
  /**
   * Get all relationships
   * @returns Promise with list of all relationships
   */
  async getAll(): Promise<APIResponse<Relationship[]>> {
    try {
      const response = await api.get<APIResponse<Relationship[]>>('/api/relationships');
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to fetch relationships',
      };
    }
  }

  /**
   * Create a new relationship
   * @param data - Relationship creation data
   * @returns Promise with created relationship
   */
  async create(data: CreateRelationshipRequest): Promise<APIResponse<Relationship>> {
    try {
      const response = await api.post<APIResponse<Relationship>>('/api/relationships', data);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to create relationship',
      };
    }
  }

  /**
   * Delete a relationship
   * @param id - Relationship ID
   * @param fromId - Source member ID
   * @param toId - Target member ID
   * @param type - Relationship type
   * @returns Promise with success message
   */
  async delete(
    id: string,
    fromId: string,
    toId: string,
    type: string
  ): Promise<APIResponse<{ message: string }>> {
    try {
      const response = await api.delete<APIResponse<{ message: string }>>(
        `/api/relationships/${id}?fromId=${fromId}&toId=${toId}&type=${type}`
      );
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to delete relationship',
      };
    }
  }

  /**
   * Find relationship path between two members
   * @param fromId - Source member ID
   * @param toId - Target member ID
   * @returns Promise with relationship path and kinship information
   */
  async findPath(fromId: string, toId: string): Promise<APIResponse<RelationshipPath>> {
    try {
      const response = await api.get<APIResponse<RelationshipPath>>(
        `/api/relationships/path?from=${fromId}&to=${toId}`
      );
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to find relationship path',
      };
    }
  }
}

// Export singleton instance
export default new RelationshipService();
