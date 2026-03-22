import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from './stores/authStore';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import DashboardPage from './pages/DashboardPage';
import MembersPage from './pages/MembersPage';
import EventsPage from './pages/EventsPage';
import FamilyTreePage from './pages/FamilyTreePage';
import RelationshipFinderPage from './pages/RelationshipFinderPage';
import SettingsPage from './pages/SettingsPage';
import Layout from './components/common/Layout';
import ErrorBoundary from './components/common/ErrorBoundary';

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000, // 5 minutes
    },
  },
});

// Protected Route wrapper - DISABLED FOR TESTING
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  // Skip authentication check - allow all access
  return <>{children}</>;
};

function App() {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            {/* Public routes */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            
            {/* Protected routes with Layout */}
            <Route path="/dashboard" element={
              <ProtectedRoute>
                <Layout>
                  <DashboardPage />
                </Layout>
              </ProtectedRoute>
            } />
            <Route path="/family-tree" element={
              <ProtectedRoute>
                <Layout>
                  <FamilyTreePage />
                </Layout>
              </ProtectedRoute>
            } />
            <Route path="/members" element={
              <ProtectedRoute>
                <Layout>
                  <MembersPage />
                </Layout>
              </ProtectedRoute>
            } />
            <Route path="/events" element={
              <ProtectedRoute>
                <Layout>
                  <EventsPage />
                </Layout>
              </ProtectedRoute>
            } />
            <Route path="/relationship-finder" element={
              <ProtectedRoute>
                <Layout>
                  <RelationshipFinderPage />
                </Layout>
              </ProtectedRoute>
            } />
            <Route path="/settings" element={
              <ProtectedRoute>
                <Layout>
                  <SettingsPage />
                </Layout>
              </ProtectedRoute>
            } />
            
            {/* Default redirect - skip login, go straight to dashboard */}
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/login" element={<Navigate to="/dashboard" replace />} />
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    </ErrorBoundary>
  );
}

export default App;
