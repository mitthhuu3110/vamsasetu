import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth.ts';
import { useTheme } from '../../contexts/ThemeContext';
import { 
  BellIcon, 
  UserCircleIcon, 
  Cog6ToothIcon,
  ArrowRightOnRectangleIcon,
  SunIcon,
  MoonIcon
} from '@heroicons/react/24/outline';

const Header: React.FC = () => {
  const { user, logout } = useAuth();
  const { isDark, toggleTheme } = useTheme();

  return (
    <header className="bg-white dark:bg-dark-card shadow-lg border-b border-gray-200 dark:border-dark-accent sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex items-center">
            <Link to="/" className="flex items-center space-x-3 group">
              <div className="w-10 h-10 bg-gradient-to-br from-soft-gold to-deep-gold rounded-xl flex items-center justify-center shadow-lg group-hover:shadow-xl transition-all duration-300">
                <span className="text-white font-bold text-lg">🌳</span>
              </div>
              <span className="text-2xl font-display font-bold text-gradient">
                VamsaSetu
              </span>
            </Link>
          </div>

          {/* Navigation */}
          <nav className="hidden md:flex space-x-2">
            <Link 
              to="/" 
              className="text-warm-brown dark:text-dark-text hover:text-soft-gold dark:hover:text-soft-gold px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-warm-beige dark:hover:bg-dark-accent"
            >
              Dashboard
            </Link>
            <Link 
              to="/family-tree" 
              className="text-warm-brown dark:text-dark-text hover:text-soft-gold dark:hover:text-soft-gold px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-warm-beige dark:hover:bg-dark-accent"
            >
              Family Tree
            </Link>
            <Link 
              to="/events" 
              className="text-warm-brown dark:text-dark-text hover:text-soft-gold dark:hover:text-soft-gold px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-warm-beige dark:hover:bg-dark-accent"
            >
              Events
            </Link>
          </nav>

          {/* User Menu */}
          <div className="flex items-center space-x-3">
            {/* Dark Mode Toggle */}
            <button
              onClick={toggleTheme}
              className="p-2 rounded-lg bg-warm-beige dark:bg-dark-accent text-warm-brown dark:text-dark-text hover:bg-soft-gold hover:text-white transition-all duration-200"
              aria-label="Toggle theme"
            >
              {isDark ? (
                <SunIcon className="h-5 w-5" />
              ) : (
                <MoonIcon className="h-5 w-5" />
              )}
            </button>

            {/* Notifications */}
            <button className="p-2 text-warm-brown dark:text-dark-text hover:text-soft-gold dark:hover:text-soft-gold relative transition-all duration-200">
              <BellIcon className="w-6 h-6" />
              <span className="absolute top-0 right-0 block h-2 w-2 rounded-full bg-red-400 ring-2 ring-white"></span>
            </button>

            {/* Profile Dropdown */}
            <div className="relative group">
              <button className="flex items-center space-x-2 text-warm-brown dark:text-dark-text hover:text-soft-gold dark:hover:text-soft-gold transition-all duration-200">
                <div className="w-8 h-8 bg-gradient-to-br from-soft-gold to-deep-gold rounded-full flex items-center justify-center">
                  <UserCircleIcon className="w-5 h-5 text-white" />
                </div>
                <span className="hidden md:block text-sm font-medium">
                  {user?.firstName} {user?.lastName}
                </span>
              </button>
              
              {/* Dropdown Menu */}
              <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-dark-card rounded-xl shadow-xl py-2 z-50 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 border border-gray-200 dark:border-dark-accent">
                <Link
                  to="/profile"
                  className="flex items-center px-4 py-2 text-sm text-warm-brown dark:text-dark-text hover:bg-warm-beige dark:hover:bg-dark-accent transition-all duration-200"
                >
                  <UserCircleIcon className="w-4 h-4 mr-3" />
                  Profile
                </Link>
                <Link
                  to="/settings"
                  className="flex items-center px-4 py-2 text-sm text-warm-brown dark:text-dark-text hover:bg-warm-beige dark:hover:bg-dark-accent transition-all duration-200"
                >
                  <Cog6ToothIcon className="w-4 h-4 mr-3" />
                  Settings
                </Link>
                <hr className="my-2 border-gray-200 dark:border-dark-accent" />
                <button
                  onClick={logout}
                  className="flex items-center w-full px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-all duration-200"
                >
                  <ArrowRightOnRectangleIcon className="w-4 h-4 mr-3" />
                  Sign out
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
