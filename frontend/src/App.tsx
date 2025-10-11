import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './hooks/useAuth';
import Layout from './components/Layout/Layout';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import FamilyTreePage from './pages/FamilyTreePage';
import EventsPage from './pages/EventsPage';
import ProfilePage from './pages/ProfilePage';
import LoadingSpinner from './components/UI/LoadingSpinner';

function App() {
  const { user, loading } = useAuth();

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="App">
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={!user ? <LoginPage /> : <Navigate to="/" />} />
        <Route path="/register" element={!user ? <RegisterPage /> : <Navigate to="/" />} />
        
        {/* Protected routes */}
        <Route path="/" element={user ? <Layout /> : <Navigate to="/login" />}>
          <Route index element={<HomePage />} />
          <Route path="family-tree" element={<FamilyTreePage />} />
          <Route path="events" element={<EventsPage />} />
          <Route path="profile" element={<ProfilePage />} />
        </Route>
        
        {/* Catch all route */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </div>
  );
}

export default App;
