export interface User {
  id: number;
  email: string;
  name: string;
  role: 'owner' | 'viewer' | 'admin';
  createdAt: string;
  updatedAt: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

// Alias for task requirement
export type LoginRequest = LoginCredentials;

export interface RegisterData {
  email: string;
  password: string;
  name: string;
}

// Alias for task requirement
export type RegisterRequest = RegisterData;

export interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
}
