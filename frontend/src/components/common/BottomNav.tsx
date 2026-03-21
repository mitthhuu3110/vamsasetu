import React from 'react';
import { NavLink } from 'react-router-dom';

interface NavItem {
  path: string;
  label: string;
  icon: string;
}

const navItems: NavItem[] = [
  { path: '/dashboard', label: 'Home', icon: '🏠' },
  { path: '/family-tree', label: 'Tree', icon: '🌳' },
  { path: '/members', label: 'Members', icon: '👥' },
  { path: '/events', label: 'Events', icon: '📅' },
];

const BottomNav: React.FC = () => {
  return (
    <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 shadow-lg z-50" aria-label="Bottom navigation">
      <div className="flex justify-around items-center h-16">
        {navItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            aria-label={item.label}
            className={({ isActive }) =>
              `flex flex-col items-center justify-center flex-1 h-full min-w-[44px] min-h-[44px] transition-colors ${
                isActive
                  ? 'text-saffron'
                  : 'text-charcoal/60 hover:text-charcoal'
              }`
            }
          >
            <span className="text-2xl mb-1" aria-hidden="true">{item.icon}</span>
            <span className="text-xs font-medium">{item.label}</span>
          </NavLink>
        ))}
      </div>
    </nav>
  );
};

export default BottomNav;
