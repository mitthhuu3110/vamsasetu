// API Configuration
export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';
export const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';

// Relationship Types
export const RELATIONSHIP_TYPES = {
  SPOUSE_OF: 'SPOUSE_OF',
  PARENT_OF: 'PARENT_OF',
  SIBLING_OF: 'SIBLING_OF',
} as const;

// Event Types
export const EVENT_TYPES = {
  BIRTHDAY: 'birthday',
  ANNIVERSARY: 'anniversary',
  CEREMONY: 'ceremony',
  CUSTOM: 'custom',
} as const;

// User Roles
export const USER_ROLES = {
  OWNER: 'owner',
  VIEWER: 'viewer',
  ADMIN: 'admin',
} as const;

// Notification Channels
export const NOTIFICATION_CHANNELS = {
  WHATSAPP: 'whatsapp',
  SMS: 'sms',
  EMAIL: 'email',
} as const;

// Gender Options
export const GENDER_OPTIONS = {
  MALE: 'male',
  FEMALE: 'female',
  OTHER: 'other',
} as const;

// Color Palette
export const COLORS = {
  SAFFRON: '#E8650A',
  TURMERIC: '#F5A623',
  IVORY: '#FBF5E6',
  TEAL: '#0D4A52',
  TEAL_LIGHT: '#0D9488',
  CHARCOAL: '#2C2420',
  ROSE: '#E11D48',
  AMBER: '#F59E0B',
} as const;

// Edge Colors by Relationship Type
export const EDGE_COLORS = {
  [RELATIONSHIP_TYPES.SPOUSE_OF]: COLORS.ROSE,
  [RELATIONSHIP_TYPES.PARENT_OF]: COLORS.TEAL_LIGHT,
  [RELATIONSHIP_TYPES.SIBLING_OF]: COLORS.AMBER,
} as const;

// Cache TTL (in milliseconds)
export const CACHE_TTL = {
  FAMILY_TREE: 5 * 60 * 1000, // 5 minutes
  MEMBER: 10 * 60 * 1000, // 10 minutes
  SEARCH: 2 * 60 * 1000, // 2 minutes
  EVENTS: 5 * 60 * 1000, // 5 minutes
} as const;

// Responsive Breakpoints
export const BREAKPOINTS = {
  MOBILE: 768,
  TABLET: 1024,
} as const;

// Touch Target Size (for mobile accessibility)
export const MIN_TOUCH_TARGET = 44; // pixels
