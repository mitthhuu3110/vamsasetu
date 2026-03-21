import React from 'react';
import { NavLink } from 'react-router-dom';

interface NavItem {
  path: string;
  label: string;
  icon: string;
}

const navItems: NavItem[] = [
  { path: '/dashboard', label: 'Dashboard', icon: '🏠' },
  { path: '/family-tree', label: 'Family Tree', icon: '🌳' },
  { path: '/members', label: 'Members', icon: '👥' },
  { path: '/events', label: 'Events', icon: '📅' },
  { path: '/relationship-finder', label: 'Relationships', icon: '🔍' },
  { path: '/settings', label: 'Settings', icon: '⚙️' },
];

const Sidebar: React.FC = () => {
  return (
    <aside className="hidden md:flex md:flex-col w-64 bg-white border-r border-gray-200 h-[calc(100vh-4rem)]" aria-label="Sidebar navigation">
      <nav className="flex-1 px-4 py-6 space-y-2" aria-label="Main menu">
        {navItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            aria-label={item.label}
            className={({ isActive }) =>
              `flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors ${
                isActive
                  ? 'bg-saffron text-white font-medium'
                  : 'text-charcoal/70 hover:bg-gray-100 hover:text-charcoal'
              }`
            }
          >
            <span className="text-xl" aria-hidden="true">{item.icon}</span>
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>
    </aside>
  );
};

export default Sidebar;
