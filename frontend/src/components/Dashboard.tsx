import { motion } from 'motion/react';
import { Users, TrendingUp, Heart, Calendar } from 'lucide-react';
import { StatsCard } from './StatsCard';
import { EventCard } from './EventCard';
import { Card } from './ui/card';
import { Button } from './ui/button';

interface DashboardProps {
  onNavigate: (page: string) => void;
}

export function Dashboard({ onNavigate }: DashboardProps) {
  const stats = [
    { title: 'Total Members', value: '48', icon: Users, gradient: 'bg-gradient-to-br from-primary to-amber-300' },
    { title: 'Generations', value: '4', icon: TrendingUp, gradient: 'bg-gradient-to-br from-secondary to-emerald-400' },
    { title: 'Family Branches', value: '7', icon: Heart, gradient: 'bg-gradient-to-br from-rose-400 to-pink-400' },
    { title: 'Upcoming Events', value: '12', icon: Calendar, gradient: 'bg-gradient-to-br from-violet-400 to-purple-400' },
  ];

  const upcomingEvents = [
    {
      name: 'Amma (Lakshmi)',
      relation: 'Mother',
      date: 'Oct 25, 2025',
      type: 'birthday' as const,
      daysLeft: 4,
    },
    {
      name: 'Ravi & Priya',
      relation: 'Brother & Sister-in-law',
      date: 'Oct 28, 2025',
      type: 'anniversary' as const,
      daysLeft: 7,
    },
    {
      name: 'Diwali Puja',
      relation: 'Family Celebration',
      date: 'Nov 1, 2025',
      type: 'puja' as const,
      daysLeft: 11,
    },
  ];

  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section with Pattern */}
      <div className="relative bg-gradient-to-br from-primary/10 via-secondary/10 to-background border-b border-border overflow-hidden">
        {/* Decorative Pattern */}
        <div className="absolute inset-0 opacity-5">
          <div className="absolute top-0 left-0 w-96 h-96">
            <svg viewBox="0 0 200 200" className="w-full h-full">
              {[...Array(5)].map((_, i) => (
                <circle
                  key={i}
                  cx="100"
                  cy="100"
                  r={80 - i * 15}
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="1"
                />
              ))}
            </svg>
          </div>
        </div>

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <h1 className="text-4xl mb-2">Welcome Back! 🙏</h1>
            <p className="text-muted-foreground text-lg">
              Your family connection dashboard
            </p>
          </motion.div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
          {stats.map((stat, index) => (
            <StatsCard key={stat.title} {...stat} delay={index * 0.1} />
          ))}
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Upcoming Events */}
          <div className="lg:col-span-2 space-y-6">
            <div className="flex items-center justify-between">
              <div>
                <h2>Upcoming Celebrations</h2>
                <p className="text-sm text-muted-foreground mt-1">
                  Never miss a special moment
                </p>
              </div>
              <Button
                variant="outline"
                onClick={() => onNavigate('events')}
                className="border-border hover:bg-muted"
              >
                View All
              </Button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {upcomingEvents.map((event, index) => (
                <motion.div
                  key={event.name}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ duration: 0.5, delay: 0.2 + index * 0.1 }}
                >
                  <EventCard {...event} />
                </motion.div>
              ))}
            </div>
          </div>

          {/* Family Tree Preview */}
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.6, delay: 0.4 }}
            className="space-y-6"
          >
            <div>
              <h2>Family Tree</h2>
              <p className="text-sm text-muted-foreground mt-1">
                Quick preview of your roots
              </p>
            </div>

            <Card className="p-6 bg-card border-border shadow-lg hover:shadow-xl transition-all duration-300 cursor-pointer" onClick={() => onNavigate('family-tree')}>
              {/* Mini Family Tree Visualization */}
              <div className="relative h-64 flex items-center justify-center">
                {/* Root Node */}
                <div className="absolute top-4 left-1/2 transform -translate-x-1/2">
                  <div className="relative">
                    <div className="w-16 h-16 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center shadow-lg border-2 border-background">
                      <span className="text-primary-foreground">GF</span>
                    </div>
                    <p className="text-xs text-center mt-2 text-foreground">Grandfather</p>
                  </div>
                </div>

                {/* Connection Lines */}
                <svg className="absolute inset-0 w-full h-full pointer-events-none">
                  <line x1="50%" y1="60" x2="30%" y2="120" stroke="url(#gradient1)" strokeWidth="2" />
                  <line x1="50%" y1="60" x2="70%" y2="120" stroke="url(#gradient1)" strokeWidth="2" />
                  <line x1="30%" y1="140" x2="20%" y2="200" stroke="url(#gradient2)" strokeWidth="2" />
                  <line x1="30%" y1="140" x2="40%" y2="200" stroke="url(#gradient2)" strokeWidth="2" />
                  
                  <defs>
                    <linearGradient id="gradient1" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" stopColor="#C9A961" stopOpacity="0.6" />
                      <stop offset="100%" stopColor="#A8C69F" stopOpacity="0.6" />
                    </linearGradient>
                    <linearGradient id="gradient2" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" stopColor="#A8C69F" stopOpacity="0.4" />
                      <stop offset="100%" stopColor="#C9A961" stopOpacity="0.4" />
                    </linearGradient>
                  </defs>
                </svg>

                {/* Second Generation */}
                <div className="absolute top-24 left-[30%] transform -translate-x-1/2">
                  <div className="w-12 h-12 rounded-full bg-gradient-to-br from-secondary to-primary flex items-center justify-center shadow-md border-2 border-background">
                    <span className="text-xs text-primary-foreground">P</span>
                  </div>
                  <p className="text-xs text-center mt-1 text-muted-foreground">Parent</p>
                </div>
                <div className="absolute top-24 right-[30%] transform translate-x-1/2">
                  <div className="w-12 h-12 rounded-full bg-gradient-to-br from-secondary to-primary flex items-center justify-center shadow-md border-2 border-background">
                    <span className="text-xs text-primary-foreground">U</span>
                  </div>
                  <p className="text-xs text-center mt-1 text-muted-foreground">Uncle</p>
                </div>

                {/* Third Generation */}
                <div className="absolute bottom-8 left-[20%] transform -translate-x-1/2">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-rose-400 to-pink-400 flex items-center justify-center shadow-sm border-2 border-background">
                    <span className="text-xs text-white">Y</span>
                  </div>
                  <p className="text-xs text-center mt-1 text-muted-foreground">You</p>
                </div>
                <div className="absolute bottom-8 left-[40%] transform -translate-x-1/2">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-rose-400 to-pink-400 flex items-center justify-center shadow-sm border-2 border-background">
                    <span className="text-xs text-white">S</span>
                  </div>
                  <p className="text-xs text-center mt-1 text-muted-foreground">Sibling</p>
                </div>
              </div>

              <Button className="w-full mt-4 bg-gradient-to-r from-primary to-secondary hover:opacity-90 text-primary-foreground">
                Explore Full Tree
              </Button>
            </Card>

            {/* Quick Actions */}
            <Card className="p-6 bg-card border-border shadow-lg">
              <h3 className="mb-4">Quick Actions</h3>
              <div className="space-y-3">
                <Button variant="outline" className="w-full justify-start border-border hover:bg-muted">
                  <Users className="w-4 h-4 mr-2" />
                  Add Family Member
                </Button>
                <Button variant="outline" className="w-full justify-start border-border hover:bg-muted">
                  <Calendar className="w-4 h-4 mr-2" />
                  Schedule Event
                </Button>
                <Button variant="outline" className="w-full justify-start border-border hover:bg-muted">
                  <Heart className="w-4 h-4 mr-2" />
                  Send Greeting
                </Button>
              </div>
            </Card>
          </motion.div>
        </div>
      </div>
    </div>
  );
}
