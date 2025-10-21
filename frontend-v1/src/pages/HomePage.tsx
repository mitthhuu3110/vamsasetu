import React from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { 
  UserGroupIcon, 
  CalendarDaysIcon, 
  ChartBarIcon,
  PlusIcon,
  BellIcon,
  CakeIcon,
  GiftIcon,
  HeartIcon,
  SparklesIcon,
  TreeIcon,
  ArrowRightIcon
} from '@heroicons/react/24/outline';

const HomePage: React.FC = () => {
  const recentEvents = [
    { 
      id: 1, 
      title: "Ravi's Birthday", 
      date: "2024-01-15", 
      type: "Birthday",
      icon: CakeIcon,
      color: "bg-pink-100 text-pink-600"
    },
    { 
      id: 2, 
      title: "Wedding Anniversary", 
      date: "2024-01-20", 
      type: "Anniversary",
      icon: HeartIcon,
      color: "bg-red-100 text-red-600"
    },
    { 
      id: 3, 
      title: "Diwali Celebration", 
      date: "2024-01-25", 
      type: "Festival",
      icon: SparklesIcon,
      color: "bg-yellow-100 text-yellow-600"
    },
    { 
      id: 4, 
      title: "Aarav's Birthday", 
      date: "2024-02-02", 
      type: "Birthday",
      icon: GiftIcon,
      color: "bg-blue-100 text-blue-600"
    },
  ];

  const familyStats = [
    { 
      label: "Total Members", 
      value: "24", 
      icon: UserGroupIcon,
      color: "from-soft-gold to-deep-gold",
      bgColor: "bg-gradient-to-br from-soft-gold/10 to-deep-gold/10"
    },
    { 
      label: "Upcoming Events", 
      value: "4", 
      icon: CalendarDaysIcon,
      color: "from-soft-green to-forest-green",
      bgColor: "bg-gradient-to-br from-soft-green/10 to-forest-green/10"
    },
    { 
      label: "Generations", 
      value: "4", 
      icon: ChartBarIcon,
      color: "from-warm-brown to-maroon",
      bgColor: "bg-gradient-to-br from-warm-brown/10 to-maroon/10"
    },
    { 
      label: "Active Relationships", 
      value: "18", 
      icon: HeartIcon,
      color: "from-saffron to-orange-500",
      bgColor: "bg-gradient-to-br from-saffron/10 to-orange-500/10"
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-warm-beige to-cream dark:from-dark-bg dark:to-dark-card indian-pattern">
      <div className="space-y-8 p-6">
        {/* Welcome Section */}
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-soft-gold via-deep-gold to-soft-gold p-8 text-white shadow-2xl"
        >
          <div className="absolute inset-0 bg-black/10"></div>
          <div className="relative z-10">
            <div className="flex items-center space-x-3 mb-4">
              <TreeIcon className="w-8 h-8" />
              <h1 className="text-4xl font-display font-bold">Welcome to VamsaSetu</h1>
            </div>
            <p className="text-xl text-white/90 mb-6">
              Your intelligent family relationship visualization and event management system
            </p>
            <div className="flex items-center space-x-4">
              <motion.div
                whileHover={{ scale: 1.05 }}
                className="bg-white/20 backdrop-blur-sm rounded-lg px-4 py-2"
              >
                <span className="text-sm font-medium">🌳 Connect Generations</span>
              </motion.div>
              <motion.div
                whileHover={{ scale: 1.05 }}
                className="bg-white/20 backdrop-blur-sm rounded-lg px-4 py-2"
              >
                <span className="text-sm font-medium">💝 Celebrate Together</span>
              </motion.div>
            </div>
          </div>
          <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full blur-3xl"></div>
          <div className="absolute bottom-0 left-0 w-24 h-24 bg-white/10 rounded-full blur-2xl"></div>
        </motion.div>

        {/* Quick Actions */}
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.2 }}
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"
        >
          <motion.div whileHover={{ scale: 1.02, y: -5 }} whileTap={{ scale: 0.98 }}>
            <Link
              to="/family-tree"
              className="card card-hover group cursor-pointer"
            >
              <div className="flex items-center">
                <div className="w-12 h-12 bg-gradient-to-br from-soft-gold to-deep-gold rounded-xl flex items-center justify-center mr-4 group-hover:shadow-lg transition-all duration-300">
                  <TreeIcon className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h3 className="font-display font-bold text-warm-brown dark:text-dark-text">Family Tree</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">View & manage family</p>
                </div>
              </div>
            </Link>
          </motion.div>

          <motion.div whileHover={{ scale: 1.02, y: -5 }} whileTap={{ scale: 0.98 }}>
            <Link
              to="/events"
              className="card card-hover group cursor-pointer"
            >
              <div className="flex items-center">
                <div className="w-12 h-12 bg-gradient-to-br from-soft-green to-forest-green rounded-xl flex items-center justify-center mr-4 group-hover:shadow-lg transition-all duration-300">
                  <CalendarDaysIcon className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h3 className="font-display font-bold text-warm-brown dark:text-dark-text">Events</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">Manage celebrations</p>
                </div>
              </div>
            </Link>
          </motion.div>

          <motion.div whileHover={{ scale: 1.02, y: -5 }} whileTap={{ scale: 0.98 }}>
            <button className="card card-hover group w-full">
              <div className="flex items-center">
                <div className="w-12 h-12 bg-gradient-to-br from-saffron to-orange-500 rounded-xl flex items-center justify-center mr-4 group-hover:shadow-lg transition-all duration-300">
                  <PlusIcon className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h3 className="font-display font-bold text-warm-brown dark:text-dark-text">Add Member</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">Expand your family</p>
                </div>
              </div>
            </button>
          </motion.div>

          <motion.div whileHover={{ scale: 1.02, y: -5 }} whileTap={{ scale: 0.98 }}>
            <button className="card card-hover group w-full">
              <div className="flex items-center">
                <div className="w-12 h-12 bg-gradient-to-br from-warm-brown to-maroon rounded-xl flex items-center justify-center mr-4 group-hover:shadow-lg transition-all duration-300">
                  <BellIcon className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h3 className="font-display font-bold text-warm-brown dark:text-dark-text">Notifications</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">View reminders</p>
                </div>
              </div>
            </button>
          </motion.div>
        </motion.div>

        {/* Stats and Recent Activity */}
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.4 }}
          className="grid grid-cols-1 lg:grid-cols-3 gap-6"
        >
          {/* Family Stats */}
          <div className="lg:col-span-2">
            <div className="card card-hover">
              <div className="flex items-center space-x-2 mb-6">
                <ChartBarIcon className="w-6 h-6 text-soft-gold" />
                <h2 className="text-2xl font-display font-bold text-warm-brown dark:text-dark-text">Family Overview</h2>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {familyStats.map((stat, index) => (
                  <motion.div 
                    key={stat.label}
                    initial={{ opacity: 0, scale: 0.8 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ duration: 0.5, delay: 0.6 + index * 0.1 }}
                    className={`text-center p-6 rounded-xl ${stat.bgColor} border border-gray-200 dark:border-dark-accent`}
                  >
                    <div className={`w-16 h-16 bg-gradient-to-br ${stat.color} rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-lg`}>
                      <stat.icon className="w-8 h-8 text-white" />
                    </div>
                    <div className="text-3xl font-display font-bold text-warm-brown dark:text-dark-text mb-2">{stat.value}</div>
                    <div className="text-sm font-medium text-gray-600 dark:text-gray-400">{stat.label}</div>
                  </motion.div>
                ))}
              </div>
            </div>
          </div>

          {/* Recent Events */}
          <div className="card card-hover">
            <div className="flex items-center space-x-2 mb-6">
              <CalendarDaysIcon className="w-6 h-6 text-soft-green" />
              <h2 className="text-2xl font-display font-bold text-warm-brown dark:text-dark-text">Upcoming Events</h2>
            </div>
            <div className="space-y-4">
              {recentEvents.map((event, index) => (
                <motion.div 
                  key={event.id}
                  initial={{ opacity: 0, x: 20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ duration: 0.5, delay: 0.8 + index * 0.1 }}
                  className="flex items-center space-x-4 p-4 bg-warm-beige dark:bg-dark-accent rounded-xl hover:shadow-md transition-all duration-200 border border-gray-200 dark:border-dark-accent"
                >
                  <div className={`w-12 h-12 ${event.color} rounded-xl flex items-center justify-center`}>
                    <event.icon className="w-6 h-6" />
                  </div>
                  <div className="flex-1">
                    <h4 className="font-display font-bold text-warm-brown dark:text-dark-text">{event.title}</h4>
                    <p className="text-sm text-gray-600 dark:text-gray-400">{event.type}</p>
                  </div>
                  <div className="text-sm font-medium text-warm-brown dark:text-dark-text">{event.date}</div>
                </motion.div>
              ))}
            </div>
            <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
              <Link
                to="/events"
                className="block mt-6 text-center btn-primary text-sm"
              >
                <span className="flex items-center justify-center space-x-2">
                  <span>View all events</span>
                  <ArrowRightIcon className="w-4 h-4" />
                </span>
              </Link>
            </motion.div>
          </div>
        </motion.div>

        {/* Quick Tips */}
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 1.0 }}
          className="card card-hover"
        >
          <div className="flex items-center space-x-2 mb-6">
            <SparklesIcon className="w-6 h-6 text-soft-gold" />
            <h2 className="text-2xl font-display font-bold text-warm-brown dark:text-dark-text">Getting Started</h2>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <motion.div 
              whileHover={{ scale: 1.02, y: -2 }}
              className="p-6 bg-gradient-to-br from-soft-gold/10 to-deep-gold/10 rounded-xl border border-soft-gold/20"
            >
              <div className="flex items-center space-x-3 mb-3">
                <TreeIcon className="w-8 h-8 text-soft-gold" />
                <h3 className="font-display font-bold text-warm-brown dark:text-dark-text">Build Your Family Tree</h3>
              </div>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Start by adding family members and their relationships to create your beautiful family tree visualization.
              </p>
            </motion.div>
            <motion.div 
              whileHover={{ scale: 1.02, y: -2 }}
              className="p-6 bg-gradient-to-br from-soft-green/10 to-forest-green/10 rounded-xl border border-soft-green/20"
            >
              <div className="flex items-center space-x-3 mb-3">
                <HeartIcon className="w-8 h-8 text-soft-green" />
                <h3 className="font-display font-bold text-warm-brown dark:text-dark-text">Set Up Events</h3>
              </div>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Add birthdays, anniversaries, and other important dates to get smart reminders and celebrate together.
              </p>
            </motion.div>
          </div>
        </motion.div>
      </div>
    </div>
  );
};

export default HomePage;
