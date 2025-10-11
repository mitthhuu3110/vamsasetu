import React from 'react';
import { Link } from 'react-router-dom';
import { 
  UserGroupIcon, 
  CalendarDaysIcon, 
  ChartBarIcon,
  PlusIcon,
  BellIcon
} from '@heroicons/react/24/outline';

const HomePage: React.FC = () => {
  const recentEvents = [
    { id: 1, title: "Ravi's Birthday", date: "2024-01-15", type: "Birthday" },
    { id: 2, title: "Wedding Anniversary", date: "2024-01-20", type: "Anniversary" },
    { id: 3, title: "Diwali Celebration", date: "2024-01-25", type: "Festival" },
  ];

  const familyStats = [
    { label: "Total Members", value: "24", icon: UserGroupIcon },
    { label: "Upcoming Events", value: "3", icon: CalendarDaysIcon },
    { label: "Generations", value: "4", icon: ChartBarIcon },
  ];

  return (
    <div className="space-y-6">
      {/* Welcome Section */}
      <div className="bg-gradient-to-r from-primary-600 to-primary-700 rounded-lg p-6 text-white">
        <h1 className="text-2xl font-bold mb-2">Welcome to VamsaSetu</h1>
        <p className="text-primary-100">
          Your intelligent family relationship visualization and event management system
        </p>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Link
          to="/family-tree"
          className="card hover:shadow-md transition-shadow cursor-pointer"
        >
          <div className="flex items-center">
            <UserGroupIcon className="w-8 h-8 text-primary-600 mr-3" />
            <div>
              <h3 className="font-medium text-gray-900">Family Tree</h3>
              <p className="text-sm text-gray-500">View & manage family</p>
            </div>
          </div>
        </Link>

        <Link
          to="/events"
          className="card hover:shadow-md transition-shadow cursor-pointer"
        >
          <div className="flex items-center">
            <CalendarDaysIcon className="w-8 h-8 text-secondary-600 mr-3" />
            <div>
              <h3 className="font-medium text-gray-900">Events</h3>
              <p className="text-sm text-gray-500">Manage celebrations</p>
            </div>
          </div>
        </Link>

        <button className="card hover:shadow-md transition-shadow">
          <div className="flex items-center">
            <PlusIcon className="w-8 h-8 text-green-600 mr-3" />
            <div>
              <h3 className="font-medium text-gray-900">Add Member</h3>
              <p className="text-sm text-gray-500">Expand your family</p>
            </div>
          </div>
        </button>

        <button className="card hover:shadow-md transition-shadow">
          <div className="flex items-center">
            <BellIcon className="w-8 h-8 text-orange-600 mr-3" />
            <div>
              <h3 className="font-medium text-gray-900">Notifications</h3>
              <p className="text-sm text-gray-500">View reminders</p>
            </div>
          </div>
        </button>
      </div>

      {/* Stats and Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Family Stats */}
        <div className="lg:col-span-2">
          <div className="card">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Family Overview</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {familyStats.map((stat) => (
                <div key={stat.label} className="text-center p-4 bg-gray-50 rounded-lg">
                  <stat.icon className="w-8 h-8 text-primary-600 mx-auto mb-2" />
                  <div className="text-2xl font-bold text-gray-900">{stat.value}</div>
                  <div className="text-sm text-gray-500">{stat.label}</div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Recent Events */}
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Upcoming Events</h2>
          <div className="space-y-3">
            {recentEvents.map((event) => (
              <div key={event.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <h4 className="font-medium text-gray-900">{event.title}</h4>
                  <p className="text-sm text-gray-500">{event.type}</p>
                </div>
                <div className="text-sm text-gray-500">{event.date}</div>
              </div>
            ))}
          </div>
          <Link
            to="/events"
            className="block mt-4 text-center text-primary-600 hover:text-primary-700 text-sm font-medium"
          >
            View all events →
          </Link>
        </div>
      </div>

      {/* Quick Tips */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Getting Started</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="p-4 bg-blue-50 rounded-lg">
            <h3 className="font-medium text-blue-900 mb-2">Build Your Family Tree</h3>
            <p className="text-sm text-blue-700">
              Start by adding family members and their relationships to create your family tree.
            </p>
          </div>
          <div className="p-4 bg-green-50 rounded-lg">
            <h3 className="font-medium text-green-900 mb-2">Set Up Events</h3>
            <p className="text-sm text-green-700">
              Add birthdays, anniversaries, and other important dates to get smart reminders.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default HomePage;
