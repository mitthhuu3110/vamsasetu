import React from 'react';
import { Link } from 'react-router-dom';
import { useAuthStore } from '../stores/authStore';
import { useUpcomingEvents } from '../hooks/useEvents';
import Card, { CardHeader, CardTitle, CardContent } from '../components/ui/Card';
import Button from '../components/ui/Button';
import LoadingSpinner from '../components/common/LoadingSpinner';
import EmptyState from '../components/common/EmptyState';

const DashboardPage: React.FC = () => {
  const { user } = useAuthStore();
  const { data: upcomingEventsResponse, isLoading } = useUpcomingEvents();
  
  const upcomingEvents = upcomingEventsResponse?.data || [];

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Welcome Section */}
      <div className="bg-gradient-to-r from-saffron to-turmeric text-white rounded-lg p-8 shadow-lg">
        <h1 className="font-display text-3xl md:text-4xl font-bold mb-2">
          Welcome back, {user?.name}!
        </h1>
        <p className="text-white/90">
          Manage your family tree, track events, and stay connected with your heritage.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {/* Quick Actions */}
        <Card variant="elevated" className="hover:shadow-xl transition-shadow">
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <Link to="/family-tree">
              <Button variant="outline" fullWidth className="justify-start">
                🌳 View Family Tree
              </Button>
            </Link>
            <Link to="/members">
              <Button variant="outline" fullWidth className="justify-start">
                👥 Manage Members
              </Button>
            </Link>
            <Link to="/events">
              <Button variant="outline" fullWidth className="justify-start">
                📅 View Events
              </Button>
            </Link>
          </CardContent>
        </Card>

        {/* Upcoming Events */}
        <Card variant="elevated" className="md:col-span-2">
          <CardHeader>
            <CardTitle>Upcoming Events</CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <LoadingSpinner />
            ) : upcomingEvents.length > 0 ? (
              <div className="space-y-3">
                {upcomingEvents.slice(0, 3).map((event) => (
                  <div
                    key={event.id}
                    className="flex items-center justify-between p-3 bg-ivory rounded-lg"
                  >
                    <div>
                      <p className="font-medium text-charcoal">{event.title}</p>
                      <p className="text-sm text-charcoal/60">
                        {new Date(event.eventDate).toLocaleDateString()}
                      </p>
                    </div>
                    <span className="text-2xl">
                      {event.eventType === 'birthday' && '🎂'}
                      {event.eventType === 'anniversary' && '💍'}
                      {event.eventType === 'ceremony' && '🎉'}
                      {event.eventType === 'custom' && '📌'}
                    </span>
                  </div>
                ))}
                <Link to="/events">
                  <Button variant="outline" fullWidth className="mt-2">
                    View All Events
                  </Button>
                </Link>
              </div>
            ) : (
              <EmptyState
                icon="📅"
                title="No upcoming events"
                description="Add events to keep track of important family dates"
                actionLabel="Add Event"
                onAction={() => window.location.href = '/events'}
              />
            )}
          </CardContent>
        </Card>
      </div>

      {/* Family Tree Summary */}
      <Card variant="elevated">
        <CardHeader>
          <CardTitle>Family Tree Overview</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="text-center p-4 bg-ivory rounded-lg">
              <div className="text-3xl font-bold text-saffron mb-1">-</div>
              <div className="text-sm text-charcoal/60">Total Members</div>
            </div>
            <div className="text-center p-4 bg-ivory rounded-lg">
              <div className="text-3xl font-bold text-teal mb-1">-</div>
              <div className="text-sm text-charcoal/60">Relationships</div>
            </div>
            <div className="text-center p-4 bg-ivory rounded-lg">
              <div className="text-3xl font-bold text-turmeric mb-1">-</div>
              <div className="text-sm text-charcoal/60">Generations</div>
            </div>
            <div className="text-center p-4 bg-ivory rounded-lg">
              <div className="text-3xl font-bold text-rose mb-1">-</div>
              <div className="text-sm text-charcoal/60">Events</div>
            </div>
          </div>
          <Link to="/family-tree">
            <Button variant="primary" fullWidth className="mt-4">
              Explore Family Tree
            </Button>
          </Link>
        </CardContent>
      </Card>
    </div>
  );
};

export default DashboardPage;
