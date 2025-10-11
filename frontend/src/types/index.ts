// User and Authentication Types
export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  phone?: string;
  profilePicture?: string;
  role: UserRole;
  createdAt: string;
  updatedAt: string;
}

export enum UserRole {
  ADMIN = 'ADMIN',
  MEMBER = 'MEMBER',
  GUEST = 'GUEST'
}

export interface AuthContextType {
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  register: (userData: RegisterData) => Promise<void>;
  logout: () => void;
  loading: boolean;
}

export interface RegisterData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  phone?: string;
}

// Family Tree Types
export interface FamilyMember {
  id: string;
  firstName: string;
  lastName: string;
  middleName?: string;
  dateOfBirth?: string;
  dateOfDeath?: string;
  gender: Gender;
  profilePicture?: string;
  phone?: string;
  email?: string;
  address?: Address;
  occupation?: string;
  notes?: string;
  isAlive: boolean;
  createdAt: string;
  updatedAt: string;
}

export enum Gender {
  MALE = 'MALE',
  FEMALE = 'FEMALE',
  OTHER = 'OTHER'
}

export interface Address {
  street?: string;
  city?: string;
  state?: string;
  country?: string;
  zipCode?: string;
}

export interface Relationship {
  id: string;
  fromMemberId: string;
  toMemberId: string;
  relationshipType: RelationshipType;
  startDate?: string;
  endDate?: string;
  isActive: boolean;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export enum RelationshipType {
  // Immediate Family
  PARENT = 'PARENT',
  CHILD = 'CHILD',
  SPOUSE = 'SPOUSE',
  SIBLING = 'SIBLING',
  
  // Extended Family
  GRANDPARENT = 'GRANDPARENT',
  GRANDCHILD = 'GRANDCHILD',
  UNCLE = 'UNCLE',
  AUNT = 'AUNT',
  NEPHEW = 'NEPHEW',
  NIECE = 'NIECE',
  COUSIN = 'COUSIN',
  
  // In-laws
  FATHER_IN_LAW = 'FATHER_IN_LAW',
  MOTHER_IN_LAW = 'MOTHER_IN_LAW',
  SON_IN_LAW = 'SON_IN_LAW',
  DAUGHTER_IN_LAW = 'DAUGHTER_IN_LAW',
  BROTHER_IN_LAW = 'BROTHER_IN_LAW',
  SISTER_IN_LAW = 'SISTER_IN_LAW',
  
  // Indian Specific Relations
  MATERNAL_UNCLE = 'MATERNAL_UNCLE',
  PATERNAL_UNCLE = 'PATERNAL_UNCLE',
  MATERNAL_AUNT = 'MATERNAL_AUNT',
  PATERNAL_AUNT = 'PATERNAL_AUNT',
  MATERNAL_GRANDFATHER = 'MATERNAL_GRANDFATHER',
  PATERNAL_GRANDFATHER = 'PATERNAL_GRANDFATHER',
  MATERNAL_GRANDMOTHER = 'MATERNAL_GRANDMOTHER',
  PATERNAL_GRANDMOTHER = 'PATERNAL_GRANDMOTHER'
}

// Event Types
export interface Event {
  id: string;
  title: string;
  description?: string;
  eventType: EventType;
  date: string;
  time?: string;
  location?: string;
  isRecurring: boolean;
  recurrencePattern?: RecurrencePattern;
  reminderSettings: ReminderSettings;
  attendees: EventAttendee[];
  createdBy: string;
  createdAt: string;
  updatedAt: string;
}

export enum EventType {
  BIRTHDAY = 'BIRTHDAY',
  ANNIVERSARY = 'ANNIVERSARY',
  WEDDING = 'WEDDING',
  FUNERAL = 'FUNERAL',
  FESTIVAL = 'FESTIVAL',
  CUSTOM = 'CUSTOM'
}

export interface RecurrencePattern {
  frequency: 'YEARLY' | 'MONTHLY' | 'WEEKLY' | 'DAILY';
  interval: number;
  endDate?: string;
  occurrences?: number;
}

export interface ReminderSettings {
  enabled: boolean;
  methods: NotificationMethod[];
  advanceTime: number; // in hours
  customMessages?: { [key in NotificationMethod]?: string };
}

export enum NotificationMethod {
  EMAIL = 'EMAIL',
  SMS = 'SMS',
  WHATSAPP = 'WHATSAPP',
  PUSH = 'PUSH'
}

export interface EventAttendee {
  memberId: string;
  memberName: string;
  role: 'HOST' | 'GUEST' | 'ORGANIZER';
  rsvpStatus: 'PENDING' | 'ACCEPTED' | 'DECLINED' | 'MAYBE';
}

// Tree Visualization Types
export interface TreeNode {
  id: string;
  type: 'person' | 'event';
  data: FamilyMember | Event;
  position: { x: number; y: number };
  style?: React.CSSProperties;
}

export interface TreeEdge {
  id: string;
  source: string;
  target: string;
  type: RelationshipType;
  label?: string;
  style?: React.CSSProperties;
}

// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  errors?: string[];
}

export interface PaginatedResponse<T> {
  content: T[];
  totalElements: number;
  totalPages: number;
  size: number;
  number: number;
  first: boolean;
  last: boolean;
}

// Search and Filter Types
export interface SearchFilters {
  query?: string;
  relationshipType?: RelationshipType;
  generation?: 'ancestors' | 'descendants' | 'siblings' | 'cousins';
  isAlive?: boolean;
  gender?: Gender;
  eventType?: EventType;
  dateRange?: {
    start: string;
    end: string;
  };
}

// Notification Types
export interface Notification {
  id: string;
  type: 'EVENT_REMINDER' | 'RELATIONSHIP_UPDATE' | 'FAMILY_UPDATE';
  title: string;
  message: string;
  isRead: boolean;
  createdAt: string;
  data?: any;
}

// Family Tree Context Types
export interface FamilyTreeContextType {
  currentTree: FamilyMember[];
  relationships: Relationship[];
  events: Event[];
  selectedMember: FamilyMember | null;
  setSelectedMember: (member: FamilyMember | null) => void;
  addMember: (member: Omit<FamilyMember, 'id' | 'createdAt' | 'updatedAt'>) => Promise<void>;
  updateMember: (id: string, updates: Partial<FamilyMember>) => Promise<void>;
  deleteMember: (id: string) => Promise<void>;
  addRelationship: (relationship: Omit<Relationship, 'id' | 'createdAt' | 'updatedAt'>) => Promise<void>;
  loading: boolean;
  error: string | null;
}
