import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User, AuthContextType, RegisterData } from '../types/index.ts';
// import { authApi } from '../services/api.ts';
import toast from 'react-hot-toast';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Temporarily skip auth bootstrap
    setLoading(false);
  }, []);

  // const fetchUserProfile = async () => {
  //   try {
  //     const userData = await authApi.getProfile();
  //     setUser(userData);
  //   } catch (error) {
  //     console.error('Failed to fetch user profile:', error);
  //     localStorage.removeItem('authToken');
  //   } finally {
  //     setLoading(false);
  //   }
  // };

  const login = async (_email: string, _password: string) => {
    // Skip API and set a mock user
    setUser({
      id: '1',
      email: 'demo@example.com',
      firstName: 'Demo',
      lastName: 'User',
      role: 'MEMBER' as any,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    });
    toast.success('Signed in (demo)');
  };

  const register = async (_userData: RegisterData) => {
    // Skip API and set a mock user
    setUser({
      id: '2',
      email: 'newuser@example.com',
      firstName: 'New',
      lastName: 'User',
      role: 'MEMBER' as any,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    });
    toast.success('Registration successful (demo)');
  };

  const logout = () => {
    localStorage.removeItem('authToken');
    setUser(null);
    toast.success('Logged out successfully');
  };

  const value: AuthContextType = {
    user,
    login,
    register,
    logout,
    loading,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
