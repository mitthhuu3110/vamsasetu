import { motion } from 'motion/react';
import { Mail, Phone, MapPin, Calendar, Users, Edit, Heart, Gift, Plus } from 'lucide-react';
import { Button } from './ui/button';
import { Card } from './ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Badge } from './ui/badge';
import { Input } from './ui/input';
import { Label } from './ui/label';

export function ProfilePage() {
  const profileData = {
    name: 'Arjun Krishna',
    relation: 'Self',
    dateOfBirth: 'March 15, 1995',
    email: 'arjun.krishna@email.com',
    phone: '+91 98765 43210',
    location: 'Hyderabad, Telangana',
    avatar: '',
  };

  const relationships = [
    { name: 'Krishna', relation: 'Father', connection: 'direct' },
    { name: 'Lakshmi', relation: 'Mother', connection: 'direct' },
    { name: 'Ravi', relation: 'Brother', connection: 'sibling' },
    { name: 'Sita', relation: 'Sister', connection: 'sibling' },
    { name: 'Venkata Rao', relation: 'Grandfather', connection: 'grandparent' },
    { name: 'Saraswati', relation: 'Grandmother', connection: 'grandparent' },
  ];

  const notes = [
    { id: 1, text: 'Father loves traditional Telugu cuisine', date: 'Oct 15, 2025' },
    { id: 2, text: "Mother's favorite flower is jasmine", date: 'Oct 10, 2025' },
    { id: 3, text: 'Plan family gathering for Diwali', date: 'Oct 5, 2025' },
  ];

  return (
    <div className="min-h-screen bg-background">
      {/* Cover Section */}
      <div className="relative h-48 bg-gradient-to-r from-primary via-secondary to-primary overflow-hidden">
        <div className="absolute inset-0 opacity-20">
          <svg className="w-full h-full" viewBox="0 0 400 200">
            {[...Array(8)].map((_, i) => (
              <circle
                key={i}
                cx={50 + i * 50}
                cy="100"
                r="40"
                fill="none"
                stroke="white"
                strokeWidth="1"
              />
            ))}
          </svg>
        </div>
      </div>

      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Profile Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="relative -mt-24"
        >
          <Card className="p-6 bg-card border-border shadow-xl">
            <div className="flex flex-col md:flex-row items-center md:items-start gap-6">
              {/* Avatar */}
              <div className="relative">
                <div className="w-32 h-32 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center overflow-hidden shadow-2xl border-4 border-background">
                  {profileData.avatar ? (
                    <img src={profileData.avatar} alt={profileData.name} className="w-full h-full object-cover" />
                  ) : (
                    <span className="text-5xl text-primary-foreground">
                      {profileData.name.charAt(0)}
                    </span>
                  )}
                </div>
                <button className="absolute bottom-0 right-0 w-10 h-10 rounded-full bg-primary text-primary-foreground shadow-lg flex items-center justify-center hover:bg-primary/90 transition-colors">
                  <Edit className="w-4 h-4" />
                </button>
              </div>

              {/* Profile Info */}
              <div className="flex-1 text-center md:text-left">
                <div className="flex flex-col md:flex-row md:items-center md:justify-between mb-4">
                  <div>
                    <h1 className="text-3xl mb-1">{profileData.name}</h1>
                    <Badge className="bg-primary/20 text-primary border-primary/30">
                      {profileData.relation}
                    </Badge>
                  </div>
                  <Button className="mt-4 md:mt-0 bg-gradient-to-r from-primary to-secondary hover:opacity-90 text-primary-foreground">
                    <Edit className="w-4 h-4 mr-2" />
                    Edit Profile
                  </Button>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
                  <div className="flex items-center justify-center md:justify-start text-muted-foreground">
                    <Calendar className="w-4 h-4 mr-2" />
                    <span>Born {profileData.dateOfBirth}</span>
                  </div>
                  <div className="flex items-center justify-center md:justify-start text-muted-foreground">
                    <Mail className="w-4 h-4 mr-2" />
                    <span>{profileData.email}</span>
                  </div>
                  <div className="flex items-center justify-center md:justify-start text-muted-foreground">
                    <Phone className="w-4 h-4 mr-2" />
                    <span>{profileData.phone}</span>
                  </div>
                  <div className="flex items-center justify-center md:justify-start text-muted-foreground">
                    <MapPin className="w-4 h-4 mr-2" />
                    <span>{profileData.location}</span>
                  </div>
                </div>
              </div>
            </div>
          </Card>
        </motion.div>

        {/* Tabbed Content */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.2 }}
          className="mt-8 mb-8"
        >
          <Tabs defaultValue="details" className="w-full">
            <TabsList className="grid w-full max-w-lg mx-auto grid-cols-3 mb-8 bg-muted">
              <TabsTrigger value="details">Personal Details</TabsTrigger>
              <TabsTrigger value="relationships">Relationships</TabsTrigger>
              <TabsTrigger value="notes">Notes</TabsTrigger>
            </TabsList>

            {/* Personal Details Tab */}
            <TabsContent value="details">
              <Card className="p-6 bg-card border-border shadow-lg">
                <h2 className="mb-6">Personal Information</h2>
                <div className="space-y-6">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="space-y-2">
                      <Label htmlFor="fullName">Full Name</Label>
                      <Input
                        id="fullName"
                        defaultValue={profileData.name}
                        className="bg-input-background border-border"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="dob">Date of Birth</Label>
                      <Input
                        id="dob"
                        defaultValue={profileData.dateOfBirth}
                        className="bg-input-background border-border"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="email">Email</Label>
                      <Input
                        id="email"
                        type="email"
                        defaultValue={profileData.email}
                        className="bg-input-background border-border"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="phone">Phone</Label>
                      <Input
                        id="phone"
                        defaultValue={profileData.phone}
                        className="bg-input-background border-border"
                      />
                    </div>
                    <div className="space-y-2 md:col-span-2">
                      <Label htmlFor="location">Location</Label>
                      <Input
                        id="location"
                        defaultValue={profileData.location}
                        className="bg-input-background border-border"
                      />
                    </div>
                  </div>
                  <Button className="bg-gradient-to-r from-primary to-secondary hover:opacity-90 text-primary-foreground">
                    Save Changes
                  </Button>
                </div>
              </Card>
            </TabsContent>

            {/* Relationships Tab */}
            <TabsContent value="relationships">
              <Card className="p-6 bg-card border-border shadow-lg">
                <div className="flex items-center justify-between mb-6">
                  <h2>Family Relationships</h2>
                  <Button className="bg-gradient-to-r from-primary to-secondary hover:opacity-90 text-primary-foreground">
                    <Plus className="w-4 h-4 mr-2" />
                    Add Relation
                  </Button>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {relationships.map((rel, index) => (
                    <motion.div
                      key={rel.name}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ duration: 0.4, delay: index * 0.1 }}
                    >
                      <Card className="p-4 bg-muted/50 border-border hover:shadow-md transition-all duration-300">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center space-x-3">
                            <div className="w-12 h-12 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center">
                              <span className="text-primary-foreground">
                                {rel.name.charAt(0)}
                              </span>
                            </div>
                            <div>
                              <p className="text-foreground">{rel.name}</p>
                              <p className="text-sm text-muted-foreground">{rel.relation}</p>
                            </div>
                          </div>
                          <Badge
                            variant="outline"
                            className={`
                              ${rel.connection === 'direct' ? 'bg-rose-100 text-rose-700 border-rose-300' : ''}
                              ${rel.connection === 'sibling' ? 'bg-amber-100 text-amber-700 border-amber-300' : ''}
                              ${rel.connection === 'grandparent' ? 'bg-emerald-100 text-emerald-700 border-emerald-300' : ''}
                            `}
                          >
                            {rel.connection === 'direct' && <Heart className="w-3 h-3 mr-1" />}
                            {rel.connection === 'sibling' && <Users className="w-3 h-3 mr-1" />}
                            {rel.connection === 'grandparent' && <Gift className="w-3 h-3 mr-1" />}
                            {rel.connection}
                          </Badge>
                        </div>
                      </Card>
                    </motion.div>
                  ))}
                </div>
              </Card>
            </TabsContent>

            {/* Notes Tab */}
            <TabsContent value="notes">
              <Card className="p-6 bg-card border-border shadow-lg">
                <div className="flex items-center justify-between mb-6">
                  <h2>Personal Notes</h2>
                  <Button className="bg-gradient-to-r from-primary to-secondary hover:opacity-90 text-primary-foreground">
                    <Plus className="w-4 h-4 mr-2" />
                    Add Note
                  </Button>
                </div>

                <div className="space-y-4">
                  {notes.map((note, index) => (
                    <motion.div
                      key={note.id}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ duration: 0.4, delay: index * 0.1 }}
                    >
                      <Card className="p-4 bg-muted/50 border-border hover:bg-muted transition-all duration-300">
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <p className="text-foreground mb-2">{note.text}</p>
                            <p className="text-xs text-muted-foreground">{note.date}</p>
                          </div>
                          <Button variant="ghost" size="sm">
                            <Edit className="w-4 h-4" />
                          </Button>
                        </div>
                      </Card>
                    </motion.div>
                  ))}
                </div>
              </Card>
            </TabsContent>
          </Tabs>
        </motion.div>
      </div>
    </div>
  );
}
