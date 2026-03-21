import api, { type APIResponse } from './api';
import type { FamilyTree } from '../types/familyTree';

/**
 * Family tree service for retrieving family tree visualization data
 */
class FamilyTreeService {
  /**
   * Get the complete family tree structure with nodes and edges
   * @returns Promise with family tree data in React Flow format
   */
  async getTree(): Promise<APIResponse<FamilyTree>> {
    try {
      const response = await api.get<APIResponse<FamilyTree>>('/api/family/tree');
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        data: null,
        error: error.response?.data?.error || 'Failed to fetch family tree',
      };
    }
  }
}

// Export singleton instance
export default new FamilyTreeService();
