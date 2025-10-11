import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { 
  User, 
  RegisterData, 
  FamilyMember, 
  Relationship, 
  Event, 
  SearchFilters,
  ApiResponse,
  PaginatedResponse 
} from '../types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor to add auth token
    this.client.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('authToken');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor to handle errors
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('authToken');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  private async request<T>(config: any): Promise<T> {
    try {
      const response: AxiosResponse<ApiResponse<T>> = await this.client(config);
      return response.data.data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || error.message || 'An error occurred');
    }
  }

  // Auth API
  async login(email: string, password: string) {
    return this.request<{ user: User; token: string }>({
      method: 'POST',
      url: '/auth/login',
      data: { email, password },
    });
  }

  async register(userData: RegisterData) {
    return this.request<{ user: User; token: string }>({
      method: 'POST',
      url: '/auth/register',
      data: userData,
    });
  }

  async getProfile() {
    return this.request<User>({
      method: 'GET',
      url: '/auth/profile',
    });
  }

  // Family Members API
  async getFamilyMembers(filters?: SearchFilters) {
    return this.request<FamilyMember[]>({
      method: 'GET',
      url: '/family/members',
      params: filters,
    });
  }

  async getFamilyMember(id: string) {
    return this.request<FamilyMember>({
      method: 'GET',
      url: `/family/members/${id}`,
    });
  }

  async createFamilyMember(member: Omit<FamilyMember, 'id' | 'createdAt' | 'updatedAt'>) {
    return this.request<FamilyMember>({
      method: 'POST',
      url: '/family/members',
      data: member,
    });
  }

  async updateFamilyMember(id: string, updates: Partial<FamilyMember>) {
    return this.request<FamilyMember>({
      method: 'PUT',
      url: `/family/members/${id}`,
      data: updates,
    });
  }

  async deleteFamilyMember(id: string) {
    return this.request<void>({
      method: 'DELETE',
      url: `/family/members/${id}`,
    });
  }

  // Relationships API
  async getRelationships() {
    return this.request<Relationship[]>({
      method: 'GET',
      url: '/family/relationships',
    });
  }

  async createRelationship(relationship: Omit<Relationship, 'id' | 'createdAt' | 'updatedAt'>) {
    return this.request<Relationship>({
      method: 'POST',
      url: '/family/relationships',
      data: relationship,
    });
  }

  async updateRelationship(id: string, updates: Partial<Relationship>) {
    return this.request<Relationship>({
      method: 'PUT',
      url: `/family/relationships/${id}`,
      data: updates,
    });
  }

  async deleteRelationship(id: string) {
    return this.request<void>({
      method: 'DELETE',
      url: `/family/relationships/${id}`,
    });
  }

  // Events API
  async getEvents(filters?: SearchFilters) {
    return this.request<Event[]>({
      method: 'GET',
      url: '/events',
      params: filters,
    });
  }

  async getEvent(id: string) {
    return this.request<Event>({
      method: 'GET',
      url: `/events/${id}`,
    });
  }

  async createEvent(event: Omit<Event, 'id' | 'createdAt' | 'updatedAt'>) {
    return this.request<Event>({
      method: 'POST',
      url: '/events',
      data: event,
    });
  }

  async updateEvent(id: string, updates: Partial<Event>) {
    return this.request<Event>({
      method: 'PUT',
      url: `/events/${id}`,
      data: updates,
    });
  }

  async deleteEvent(id: string) {
    return this.request<void>({
      method: 'DELETE',
      url: `/events/${id}`,
    });
  }

  // Tree Visualization API
  async getFamilyTree() {
    return this.request<{ nodes: FamilyMember[]; edges: Relationship[] }>({
      method: 'GET',
      url: '/family/tree',
    });
  }

  async getRelationshipPath(fromMemberId: string, toMemberId: string) {
    return this.request<{ path: string[]; description: string }>({
      method: 'GET',
      url: `/family/relationship-path/${fromMemberId}/${toMemberId}`,
    });
  }

  // Search API
  async searchFamily(query: string, filters?: SearchFilters) {
    return this.request<{ members: FamilyMember[]; events: Event[] }>({
      method: 'GET',
      url: '/search',
      params: { query, ...filters },
    });
  }
}

export const apiClient = new ApiClient();

// Export specific API modules for easier imports
export const authApi = {
  login: (email: string, password: string) => apiClient.login(email, password),
  register: (userData: RegisterData) => apiClient.register(userData),
  getProfile: () => apiClient.getProfile(),
};

export const familyApi = {
  getMembers: (filters?: SearchFilters) => apiClient.getFamilyMembers(filters),
  getMember: (id: string) => apiClient.getFamilyMember(id),
  createMember: (member: Omit<FamilyMember, 'id' | 'createdAt' | 'updatedAt'>) => 
    apiClient.createFamilyMember(member),
  updateMember: (id: string, updates: Partial<FamilyMember>) => 
    apiClient.updateFamilyMember(id, updates),
  deleteMember: (id: string) => apiClient.deleteFamilyMember(id),
  getRelationships: () => apiClient.getRelationships(),
  createRelationship: (relationship: Omit<Relationship, 'id' | 'createdAt' | 'updatedAt'>) => 
    apiClient.createRelationship(relationship),
  updateRelationship: (id: string, updates: Partial<Relationship>) => 
    apiClient.updateRelationship(id, updates),
  deleteRelationship: (id: string) => apiClient.deleteRelationship(id),
  getTree: () => apiClient.getFamilyTree(),
  getRelationshipPath: (fromMemberId: string, toMemberId: string) => 
    apiClient.getRelationshipPath(fromMemberId, toMemberId),
};

export const eventsApi = {
  getEvents: (filters?: SearchFilters) => apiClient.getEvents(filters),
  getEvent: (id: string) => apiClient.getEvent(id),
  createEvent: (event: Omit<Event, 'id' | 'createdAt' | 'updatedAt'>) => 
    apiClient.createEvent(event),
  updateEvent: (id: string, updates: Partial<Event>) => 
    apiClient.updateEvent(id, updates),
  deleteEvent: (id: string) => apiClient.deleteEvent(id),
};

export const searchApi = {
  search: (query: string, filters?: SearchFilters) => apiClient.searchFamily(query, filters),
};
