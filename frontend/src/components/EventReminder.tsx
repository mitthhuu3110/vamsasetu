import { useState } from 'react';
import { motion } from 'motion/react';
import { Calendar as CalendarIcon, Heart, Sparkles, Filter, Plus } from 'lucide-react';
import { EventCard } from './EventCard';
import { Button } from './ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Card } from './ui/card';

export function EventReminder() {
  const [selectedFilter, setSelectedFilter] = useState('all');

  const events = [
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
      name: 'Nanna (Krishna)',
      relation: 'Father',
      date: 'Nov 15, 2025',
      type: 'birthday' as const,
      daysLeft: 25,
    },
    {
      name: 'Diwali Puja',
      relation: 'Family Celebration',
      date: 'Nov 1, 2025',
      type: 'puja' as const,
      daysLeft: 11,
    },
    {
      name: 'Sita',
      relation: 'Sister',
      date: 'Nov 8, 2025',
      type: 'birthday' as const,
      daysLeft: 18,
    },
    {
      name: 'Sankranti Puja',
      relation: 'Family Celebration',
      date: 'Jan 14, 2026',
      type: 'puja' as const,
      daysLeft: 85,
    },
    {
      name: 'Parents Anniversary',
      relation: 'Krishna & Lakshmi',
      date: 'Dec 5, 2025',
      type: 'anniversary' as const,
      daysLeft: 45,
    },
    {
      name: 'Arun',
      relation: 'Cousin',
      date: 'Nov 22, 2025',
      type: 'birthday' as const,
      daysLeft: 32,
    },
  ];

  const filteredEvents = events
    .filter((event) => selectedFilter === 'all' || event.type === selectedFilter)
    .sort((a, b) => a.daysLeft - b.daysLeft);

  const upcomingThisWeek = filteredEvents.filter((e) => e.daysLeft <= 7);
  const upcomingThisMonth = filteredEvents.filter((e) => e.daysLeft > 7 && e.daysLeft <= 30);
  const upcomingLater = filteredEvents.filter((e) => e.daysLeft > 30);

  return (
    <div className="min-h-screen bg-background">
      {/* Header Section */}
      <div className="relative bg-gradient-to-br from-rose-50 via-amber-50 to-emerald-50 border-b border-border overflow-hidden">
        {/* Decorative Elements */}
        <div className="absolute top-10 right-10 w-32 h-32 opacity-10">
          <svg viewBox="0 0 100 100" className="w-full h-full">
            <path
              d="M50,10 L60,40 L90,40 L65,60 L75,90 L50,70 L25,90 L35,60 L10,40 L40,40 Z"
              fill="currentColor"
            />
          </svg>
        </div>
        <div className="absolute bottom-10 left-10 w-24 h-24 opacity-10">
          <Heart className="w-full h-full" />
        </div>

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="flex items-center justify-between"
          >
            <div>
              <h1 className="text-4xl mb-2">Event Reminders</h1>
              <p className="text-muted-foreground text-lg">
                Never miss a celebration with your loved ones
              </p>
            </div>
            <Button className="bg-gradient-to-r from-primary to-secondary hover:opacity-90 text-primary-foreground shadow-lg">
              <Plus className="w-4 h-4 mr-2" />
              Add Event
            </Button>
          </motion.div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Filters */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
          className="mb-8"
        >
          <Card className="p-4 bg-card border-border shadow-md">
            <div className="flex items-center justify-between flex-wrap gap-4">
              <div className="flex items-center space-x-2">
                <Filter className="w-5 h-5 text-muted-foreground" />
                <span className="text-muted-foreground">Filter by:</span>
              </div>
              <div className="flex items-center space-x-2">
                <Button
                  variant={selectedFilter === 'all' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setSelectedFilter('all')}
                  className={selectedFilter === 'all' ? 'bg-primary text-primary-foreground' : 'border-border'}
                >
                  All Events
                </Button>
                <Button
                  variant={selectedFilter === 'birthday' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setSelectedFilter('birthday')}
                  className={selectedFilter === 'birthday' ? 'bg-rose-500 text-white hover:bg-rose-600' : 'border-border'}
                >
                  <Heart className="w-3 h-3 mr-1" />
                  Birthdays
                </Button>
                <Button
                  variant={selectedFilter === 'anniversary' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setSelectedFilter('anniversary')}
                  className={selectedFilter === 'anniversary' ? 'bg-amber-500 text-white hover:bg-amber-600' : 'border-border'}
                >
                  <Sparkles className="w-3 h-3 mr-1" />
                  Anniversaries
                </Button>
                <Button
                  variant={selectedFilter === 'puja' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setSelectedFilter('puja')}
                  className={selectedFilter === 'puja' ? 'bg-emerald-500 text-white hover:bg-emerald-600' : 'border-border'}
                >
                  <CalendarIcon className="w-3 h-3 mr-1" />
                  Puja
                </Button>
              </div>
            </div>
          </Card>
        </motion.div>

        {/* Stats Cards */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
          className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8"
        >
          <Card className="p-6 bg-gradient-to-br from-rose-50 to-rose-100 border-rose-200 shadow-md">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-rose-700 text-sm mb-1">This Week</p>
                <p className="text-3xl text-rose-900">{upcomingThisWeek.length}</p>
              </div>
              <div className="w-12 h-12 rounded-full bg-rose-200 flex items-center justify-center">
                <Heart className="w-6 h-6 text-rose-700" />
              </div>
            </div>
          </Card>

          <Card className="p-6 bg-gradient-to-br from-amber-50 to-amber-100 border-amber-200 shadow-md">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-amber-700 text-sm mb-1">This Month</p>
                <p className="text-3xl text-amber-900">{upcomingThisMonth.length}</p>
              </div>
              <div className="w-12 h-12 rounded-full bg-amber-200 flex items-center justify-center">
                <CalendarIcon className="w-6 h-6 text-amber-700" />
              </div>
            </div>
          </Card>

          <Card className="p-6 bg-gradient-to-br from-emerald-50 to-emerald-100 border-emerald-200 shadow-md">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-emerald-700 text-sm mb-1">Later</p>
                <p className="text-3xl text-emerald-900">{upcomingLater.length}</p>
              </div>
              <div className="w-12 h-12 rounded-full bg-emerald-200 flex items-center justify-center">
                <Sparkles className="w-6 h-6 text-emerald-700" />
              </div>
            </div>
          </Card>
        </motion.div>

        {/* Tabbed Event Lists */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.3 }}
        >
          <Tabs defaultValue="week" className="w-full">
            <TabsList className="grid w-full max-w-md mx-auto grid-cols-3 mb-8 bg-muted">
              <TabsTrigger value="week">This Week</TabsTrigger>
              <TabsTrigger value="month">This Month</TabsTrigger>
              <TabsTrigger value="later">Later</TabsTrigger>
            </TabsList>

            <TabsContent value="week">
              {upcomingThisWeek.length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {upcomingThisWeek.map((event, index) => (
                    <motion.div
                      key={event.name + event.date}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ duration: 0.4, delay: index * 0.1 }}
                    >
                      <EventCard {...event} />
                    </motion.div>
                  ))}
                </div>
              ) : (
                <Card className="p-12 text-center bg-card border-border">
                  <CalendarIcon className="w-12 h-12 mx-auto mb-4 text-muted-foreground" />
                  <p className="text-muted-foreground">No events this week</p>
                </Card>
              )}
            </TabsContent>

            <TabsContent value="month">
              {upcomingThisMonth.length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {upcomingThisMonth.map((event, index) => (
                    <motion.div
                      key={event.name + event.date}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ duration: 0.4, delay: index * 0.1 }}
                    >
                      <EventCard {...event} />
                    </motion.div>
                  ))}
                </div>
              ) : (
                <Card className="p-12 text-center bg-card border-border">
                  <CalendarIcon className="w-12 h-12 mx-auto mb-4 text-muted-foreground" />
                  <p className="text-muted-foreground">No events this month</p>
                </Card>
              )}
            </TabsContent>

            <TabsContent value="later">
              {upcomingLater.length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {upcomingLater.map((event, index) => (
                    <motion.div
                      key={event.name + event.date}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ duration: 0.4, delay: index * 0.1 }}
                    >
                      <EventCard {...event} />
                    </motion.div>
                  ))}
                </div>
              ) : (
                <Card className="p-12 text-center bg-card border-border">
                  <CalendarIcon className="w-12 h-12 mx-auto mb-4 text-muted-foreground" />
                  <p className="text-muted-foreground">No upcoming events</p>
                </Card>
              )}
            </TabsContent>
          </Tabs>
        </motion.div>
      </div>
    </div>
  );
}
