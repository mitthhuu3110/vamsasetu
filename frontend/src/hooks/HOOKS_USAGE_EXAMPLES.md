# React Query Hooks Usage Examples

This document provides examples of how to use the React Query hooks in VamsaSetu.

## Authentication Hooks

### useLogin

```typescript
import { useLogin } from './hooks';

function LoginForm() {
  const login = useLogin();

  const handleSubmit = async (email: string, password: string) => {
    const result = await login.mutateAsync({ email, password });
    
    if (result.success) {
      // User is automatically logged in via authStore
      console.log('Logged in:', result.data?.user);
    } else {
      console.error('Login failed:', result.error);
    }
  };

  return (
    <form onSubmit={(e) => {
      e.preventDefault();
      handleSubmit('user@example.com', 'password');
    }}>
      {login.isPending && <p>Logging in...</p>}
      {login.isError && <p>Error: {login.error.message}</p>}
      <button type="submit" disabled={login.isPending}>Login</button>
    </form>
  );
}
```

### useRegister

```typescript
import { useRegister } from './hooks';

function RegisterForm() {
  const register = useRegister();

  const handleSubmit = async (data: RegisterRequest) => {
    const result = await register.mutateAsync(data);
    
    if (result.success) {
      // User is automatically registered and logged in
      console.log('Registered:', result.data?.user);
    }
  };

  return (
    <button onClick={() => handleSubmit({
      email: 'new@example.com',
      password: 'password',
      name: 'New User'
    })}>
      Register
    </button>
  );
}
```

### useProfile

```typescript
import { useProfile } from './hooks';

function ProfilePage() {
  const { data, isLoading, error } = useProfile();

  if (isLoading) return <div>Loading profile...</div>;
  if (error) return <div>Error: {error.message}</div>;
  if (!data?.success) return <div>Failed to load profile</div>;

  return (
    <div>
      <h1>{data.data?.name}</h1>
      <p>{data.data?.email}</p>
      <p>Role: {data.data?.role}</p>
    </div>
  );
}
```

## Member Hooks

### useMembers

```typescript
import { useMembers } from './hooks';

function MemberList() {
  const { data, isLoading } = useMembers({ page: 1, limit: 10 });

  if (isLoading) return <div>Loading members...</div>;

  return (
    <ul>
      {data?.data?.members.map(member => (
        <li key={member.id}>{member.name}</li>
      ))}
    </ul>
  );
}
```

### useCreateMember

```typescript
import { useCreateMember } from './hooks';

function AddMemberForm() {
  const createMember = useCreateMember();

  const handleSubmit = async () => {
    const result = await createMember.mutateAsync({
      name: 'John Doe',
      dateOfBirth: '1990-01-01',
      gender: 'male',
      email: 'john@example.com',
    });

    if (result.success) {
      // Cache is automatically invalidated
      console.log('Member created:', result.data);
    }
  };

  return (
    <button onClick={handleSubmit} disabled={createMember.isPending}>
      Add Member
    </button>
  );
}
```

### useUpdateMember

```typescript
import { useUpdateMember } from './hooks';

function EditMemberForm({ memberId }: { memberId: string }) {
  const updateMember = useUpdateMember();

  const handleUpdate = async () => {
    await updateMember.mutateAsync({
      id: memberId,
      data: { name: 'Updated Name' }
    });
  };

  return <button onClick={handleUpdate}>Update</button>;
}
```

### useDeleteMember

```typescript
import { useDeleteMember } from './hooks';

function DeleteMemberButton({ memberId }: { memberId: string }) {
  const deleteMember = useDeleteMember();

  return (
    <button onClick={() => deleteMember.mutate(memberId)}>
      Delete
    </button>
  );
}
```

### useSearchMembers

```typescript
import { useState } from 'react';
import { useSearchMembers } from './hooks';

function MemberSearch() {
  const [query, setQuery] = useState('');
  const { data, isLoading } = useSearchMembers(query);

  return (
    <div>
      <input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search members..."
      />
      {isLoading && <p>Searching...</p>}
      {data?.data?.members.map(member => (
        <div key={member.id}>{member.name}</div>
      ))}
    </div>
  );
}
```

## Relationship Hooks

### useRelationships

```typescript
import { useRelationships } from './hooks';

function RelationshipList() {
  const { data, isLoading } = useRelationships();

  if (isLoading) return <div>Loading...</div>;

  return (
    <ul>
      {data?.data?.map((rel, idx) => (
        <li key={idx}>
          {rel.fromId} → {rel.type} → {rel.toId}
        </li>
      ))}
    </ul>
  );
}
```

### useCreateRelationship

```typescript
import { useCreateRelationship } from './hooks';

function AddRelationshipForm() {
  const createRelationship = useCreateRelationship();

  const handleSubmit = async () => {
    await createRelationship.mutateAsync({
      type: 'PARENT_OF',
      fromId: 'parent-id',
      toId: 'child-id',
    });
  };

  return <button onClick={handleSubmit}>Add Relationship</button>;
}
```

### useDeleteRelationship

```typescript
import { useDeleteRelationship } from './hooks';

function DeleteRelationshipButton() {
  const deleteRelationship = useDeleteRelationship();

  return (
    <button onClick={() => deleteRelationship.mutate({
      id: 'rel-id',
      fromId: 'from-id',
      toId: 'to-id',
      type: 'SPOUSE_OF'
    })}>
      Delete
    </button>
  );
}
```

### useFindPath

```typescript
import { useFindPath } from './hooks';

function RelationshipPath({ fromId, toId }: { fromId: string; toId: string }) {
  const { data, isLoading } = useFindPath(fromId, toId);

  if (isLoading) return <div>Finding path...</div>;
  if (!data?.success) return <div>No path found</div>;

  return (
    <div>
      <h3>{data.data?.kinshipTerm}</h3>
      <p>{data.data?.description}</p>
      <div>
        {data.data?.path.map((node, idx) => (
          <span key={node.id}>
            {node.name}
            {idx < data.data!.path.length - 1 && ' → '}
          </span>
        ))}
      </div>
    </div>
  );
}
```

## Event Hooks

### useEvents

```typescript
import { useEvents } from './hooks';

function EventList() {
  const { data, isLoading } = useEvents({ 
    page: 1, 
    limit: 10,
    type: 'birthday' 
  });

  if (isLoading) return <div>Loading events...</div>;

  return (
    <ul>
      {data?.data?.events.map(event => (
        <li key={event.id}>
          {event.title} - {event.eventDate}
        </li>
      ))}
    </ul>
  );
}
```

### useCreateEvent

```typescript
import { useCreateEvent } from './hooks';

function AddEventForm() {
  const createEvent = useCreateEvent();

  const handleSubmit = async () => {
    await createEvent.mutateAsync({
      title: 'Birthday Party',
      eventDate: '2024-12-25',
      eventType: 'birthday',
      memberIds: ['member-1', 'member-2'],
      description: 'Annual birthday celebration',
    });
  };

  return <button onClick={handleSubmit}>Create Event</button>;
}
```

### useUpdateEvent

```typescript
import { useUpdateEvent } from './hooks';

function EditEventForm({ eventId }: { eventId: number }) {
  const updateEvent = useUpdateEvent();

  return (
    <button onClick={() => updateEvent.mutate({
      id: eventId,
      data: { title: 'Updated Title' }
    })}>
      Update Event
    </button>
  );
}
```

### useDeleteEvent

```typescript
import { useDeleteEvent } from './hooks';

function DeleteEventButton({ eventId }: { eventId: number }) {
  const deleteEvent = useDeleteEvent();

  return (
    <button onClick={() => deleteEvent.mutate(eventId)}>
      Delete Event
    </button>
  );
}
```

### useUpcomingEvents

```typescript
import { useUpcomingEvents } from './hooks';

function UpcomingEventsList() {
  const { data, isLoading } = useUpcomingEvents(7); // Next 7 days

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h2>Upcoming Events</h2>
      {data?.data?.map(event => (
        <div key={event.id}>
          {event.title} - {event.eventDate}
        </div>
      ))}
    </div>
  );
}
```

## Family Tree Hook

### useFamilyTree

```typescript
import { useFamilyTree } from './hooks';
import ReactFlow from 'reactflow';

function FamilyTreeVisualization() {
  const { data, isLoading } = useFamilyTree();

  if (isLoading) return <div>Loading family tree...</div>;
  if (!data?.success) return <div>Failed to load tree</div>;

  return (
    <div style={{ width: '100%', height: '600px' }}>
      <ReactFlow
        nodes={data.data?.nodes || []}
        edges={data.data?.edges || []}
        fitView
      />
    </div>
  );
}
```

## Cache Invalidation

All mutation hooks automatically invalidate relevant cache entries:

- **Member mutations** → Invalidates `['members']` and `['familyTree']`
- **Relationship mutations** → Invalidates `['relationships']` and `['familyTree']`
- **Event mutations** → Invalidates `['events']`
- **Auth mutations** → Invalidates `['profile']`

This ensures the UI always displays fresh data after mutations.

## Error Handling

All hooks return error states that can be used for error handling:

```typescript
function Component() {
  const { data, isLoading, error, isError } = useMembers();

  if (isError) {
    return <div>Error: {error.message}</div>;
  }

  // For mutations
  const createMember = useCreateMember();
  
  if (createMember.isError) {
    return <div>Failed to create member: {createMember.error.message}</div>;
  }

  // ...
}
```

## Loading States

All hooks provide loading states:

```typescript
function Component() {
  const { isLoading, isPending, isFetching } = useMembers();
  const createMember = useCreateMember();

  return (
    <div>
      {isLoading && <Spinner />}
      {createMember.isPending && <p>Creating member...</p>}
    </div>
  );
}
```
