export interface Member {
  id: string;
  name: string;
  dateOfBirth: string;
  gender: 'male' | 'female' | 'other';
  email: string;
  phone: string;
  avatarUrl: string;
  createdAt: string;
  updatedAt: string;
  isDeleted: boolean;
}

export interface CreateMemberData {
  name: string;
  dateOfBirth: string;
  gender: 'male' | 'female' | 'other';
  email?: string;
  phone?: string;
  avatarUrl?: string;
}

// Alias for task requirement
export type CreateMemberRequest = CreateMemberData;

export interface UpdateMemberData extends Partial<CreateMemberData> {
  id: string;
}

export interface MemberNodeData {
  id: string;
  name: string;
  avatarUrl: string;
  relationBadge: string;
  hasUpcomingEvent: boolean;
  gender: string;
}
