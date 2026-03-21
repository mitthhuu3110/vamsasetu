import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../stores/authStore';
import Button from '../ui/Button';

const Navbar: React.FC = () => {
  const navigate = useNavigate();
  const { user, clearAuth } = useAuthStore();

  const handleLogout = () => {
    clearAuth();
    navigate('/login');
  };

  return (
    <nav className="bg-white border-b border-gray-200 shadow-sm" aria-label="Main navigation">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex items-center">
            <h1 className="font-display text-2xl font-bold text-saffron">
              VamsaSetu
            </h1>
          </div>

          {/* User info and logout */}
          <div className="flex items-center space-x-4">
            {user && (
              <div className="text-sm text-charcoal/70" aria-label={`Logged in as ${user.name}`}>
                Welcome, <span className="font-medium text-charcoal">{user.name}</span>
              </div>
            )}
            <Button
              variant="outline"
              onClick={handleLogout}
              className="text-sm"
              aria-label="Logout from VamsaSetu"
            >
              Logout
            </Button>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
