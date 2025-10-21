import React from 'react';
import { useState } from 'react';
import { Navigation } from './components/Navigation';
import { LoginPage } from './components/LoginPage';
import { Dashboard } from './components/Dashboard';
import { FamilyTree } from './components/FamilyTree';
import { EventReminder } from './components/EventReminder';
import { ProfilePage } from './components/ProfilePage';
import { Card } from './components/ui/card';
import { Settings as SettingsIcon, Bell, Shield, Palette, Users } from 'lucide-react';
import { Button } from './components/ui/button';
import { Switch } from './components/ui/switch';
import { Label } from './components/ui/label';

export default function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [currentPage, setCurrentPage] = useState('dashboard');

  const handleLogin = () => {
    setIsLoggedIn(true);
  };

  const handleNavigate = (page: string) => {
    setCurrentPage(page);
  };

  // Settings Page Component
  const SettingsPage = () => (
    <div className="min-h-screen bg-background">
      <div className="relative bg-gradient-to-br from-primary/10 via-secondary/10 to-background border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <h1 className="text-4xl mb-2">Settings</h1>
          <p className="text-muted-foreground text-lg">
            Manage your account and preferences
          </p>
        </div>
      </div>

      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-6">
          {/* Notifications Settings */}
          <Card className="p-6 bg-card border-border shadow-lg">
            <div className="flex items-start space-x-4 mb-6">
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-primary to-secondary flex items-center justify-center">
                <Bell className="w-5 h-5 text-primary-foreground" />
              </div>
              <div>
                <h2>Notifications</h2>
                <p className="text-sm text-muted-foreground">
                  Manage how you receive event reminders
                </p>
              </div>
            </div>
            <div className="space-y-4">
              <div className="flex items-center justify-between py-3 border-b border-border">
                <div>
                  <Label htmlFor="email-notifications">Email Notifications</Label>
                  <p className="text-sm text-muted-foreground">Receive reminders via email</p>
                </div>
                <Switch id="email-notifications" defaultChecked />
              </div>
              <div className="flex items-center justify-between py-3 border-b border-border">
                <div>
                  <Label htmlFor="sms-notifications">SMS Notifications</Label>
                  <p className="text-sm text-muted-foreground">Receive reminders via SMS</p>
                </div>
                <Switch id="sms-notifications" />
              </div>
              <div className="flex items-center justify-between py-3">
                <div>
                  <Label htmlFor="push-notifications">Push Notifications</Label>
                  <p className="text-sm text-muted-foreground">Receive browser notifications</p>
                </div>
                <Switch id="push-notifications" defaultChecked />
              </div>
            </div>
          </Card>

          {/* Privacy Settings */}
          <Card className="p-6 bg-card border-border shadow-lg">
            <div className="flex items-start space-x-4 mb-6">
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-emerald-400 to-teal-500 flex items-center justify-center">
                <Shield className="w-5 h-5 text-white" />
              </div>
              <div>
                <h2>Privacy</h2>
                <p className="text-sm text-muted-foreground">
                  Control your family tree visibility
                </p>
              </div>
            </div>
            <div className="space-y-4">
              <div className="flex items-center justify-between py-3 border-b border-border">
                <div>
                  <Label htmlFor="public-tree">Public Family Tree</Label>
                  <p className="text-sm text-muted-foreground">Allow others to view your tree</p>
                </div>
                <Switch id="public-tree" />
              </div>
              <div className="flex items-center justify-between py-3">
                <div>
                  <Label htmlFor="share-events">Share Events</Label>
                  <p className="text-sm text-muted-foreground">Let family members see events</p>
                </div>
                <Switch id="share-events" defaultChecked />
              </div>
            </div>
          </Card>

          {/* Appearance Settings */}
          <Card className="p-6 bg-card border-border shadow-lg">
            <div className="flex items-start space-x-4 mb-6">
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-violet-400 to-purple-500 flex items-center justify-center">
                <Palette className="w-5 h-5 text-white" />
              </div>
              <div>
                <h2>Appearance</h2>
                <p className="text-sm text-muted-foreground">
                  Customize your interface theme
                </p>
              </div>
            </div>
            <div className="space-y-4">
              <div className="flex items-center justify-between py-3">
                <div>
                  <Label htmlFor="dark-mode">Dark Mode</Label>
                  <p className="text-sm text-muted-foreground">Use darker colors (coming soon)</p>
                </div>
                <Switch id="dark-mode" disabled />
              </div>
            </div>
          </Card>

          {/* Family Settings */}
          <Card className="p-6 bg-card border-border shadow-lg">
            <div className="flex items-start space-x-4 mb-6">
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-rose-400 to-pink-500 flex items-center justify-center">
                <Users className="w-5 h-5 text-white" />
              </div>
              <div>
                <h2>Family Management</h2>
                <p className="text-sm text-muted-foreground">
                  Manage your family tree settings
                </p>
              </div>
            </div>
            <div className="space-y-3">
              <Button variant="outline" className="w-full justify-start border-border hover:bg-muted">
                Export Family Data
              </Button>
              <Button variant="outline" className="w-full justify-start border-border hover:bg-muted">
                Invite Family Members
              </Button>
              <Button variant="outline" className="w-full justify-start border-border hover:bg-muted text-destructive hover:text-destructive">
                Delete Account
              </Button>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );

  // Render login page if not logged in
  if (!isLoggedIn) {
    return <LoginPage onLogin={handleLogin} />;
  }

  // Render main app with navigation
  return (
    <div className="min-h-screen bg-background">
      <Navigation currentPage={currentPage} onNavigate={handleNavigate} />
      
      {currentPage === 'dashboard' && <Dashboard onNavigate={handleNavigate} />}
      {currentPage === 'family-tree' && <FamilyTree />}
      {currentPage === 'events' && <EventReminder />}
      {currentPage === 'profile' && <ProfilePage />}
      {currentPage === 'settings' && <SettingsPage />}
    </div>
  );
}
