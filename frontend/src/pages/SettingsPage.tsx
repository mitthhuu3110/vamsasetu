import React from 'react';
import Card, { CardHeader, CardTitle, CardContent } from '../components/ui/Card';
import { useAuthStore } from '../stores/authStore';

const SettingsPage: React.FC = () => {
  const { user } = useAuthStore();

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="font-display text-3xl font-bold text-charcoal">
          Settings
        </h1>
        <p className="text-charcoal/70 mt-1">
          Manage your account and preferences
        </p>
      </div>

      {/* Profile Settings */}
      <Card variant="elevated">
        <CardHeader>
          <CardTitle>Profile Information</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-charcoal mb-1">
              Name
            </label>
            <p className="text-charcoal/70">{user?.name || 'Not set'}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-charcoal mb-1">
              Email
            </label>
            <p className="text-charcoal/70">{user?.email || 'Not set'}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-charcoal mb-1">
              Role
            </label>
            <p className="text-charcoal/70 capitalize">{user?.role || 'Not set'}</p>
          </div>
        </CardContent>
      </Card>

      {/* Notification Preferences */}
      <Card variant="elevated">
        <CardHeader>
          <CardTitle>Notification Preferences</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-charcoal/60 text-sm">
            Notification settings coming soon...
          </p>
        </CardContent>
      </Card>

      {/* Theme Settings */}
      <Card variant="elevated">
        <CardHeader>
          <CardTitle>Theme Settings</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-charcoal/60 text-sm">
            Theme customization coming soon...
          </p>
        </CardContent>
      </Card>
    </div>
  );
};

export default SettingsPage;
