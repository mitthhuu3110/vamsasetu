import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './hooks/useAuth.ts';
import Layout from './components/Layout/Layout.tsx';
import HomePage from './pages/HomePage.tsx';
import LoginPage from './pages/LoginPage.tsx';
import RegisterPage from './pages/RegisterPage.tsx';
import FamilyTreePage from './pages/FamilyTreePage.tsx';
import EventsPage from './pages/EventsPage.tsx';
import ProfilePage from './pages/ProfilePage.tsx';
import LoadingSpinner from './components/UI/LoadingSpinner.tsx';

function App() {
  const { user, loading } = useAuth();

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="App">
      <Routes>
        {/* Public routes (auth bypass enabled) */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        
        {/* App routes (auth bypass: always render) */}
        <Route path="/" element={<Layout />}>
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
