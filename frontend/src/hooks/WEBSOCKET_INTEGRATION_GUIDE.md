# WebSocket Integration Guide

## Quick Start

### 1. Basic Setup (App-Level)
Add WebSocket to your main App component for global real-time updates:

```typescript
// src/App.tsx
import { useWebSocket } from './hooks';

function App() {
  // Enable WebSocket for all authenticated users
  useWebSocket();
  
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        {/* Your routes */}
      </BrowserRouter>
    </QueryClientProvider>
  );
}
```

### 2. Conditional Connection
Only connect on specific pages to save resources:

```typescript
// src/pages/DashboardPage.tsx
import { useWebSocket } from '../hooks';

function DashboardPage() {
  // Only connect when on dashboard
  useWebSocket({ enabled: true });
  
  return <div>Dashboard content</div>;
}
```

### 3. With Status Indicator
Show connection status to users:

```typescript
function Header() {
  const { isConnected } = useWebSocket();
  
  return (
    <header>
      <h1>VamsaSetu</h1>
      <div>
        {isConnected ? (
          <span>🟢 Live</span>
        ) : (
          <span>🔴 Offline</span>
        )}
      </div>
    </header>
  );
}
```

## Advanced Usage

### Custom Message Handling
Handle specific events with custom logic:

```typescript
function FamilyTreePage() {
  const [notification, setNotification] = React.useState('');
  
  useWebSocket({
    onMessage: (message) => {
      switch (message.type) {
        case 'member_created':
          setNotification(`New member: ${message.data.name}`);
          break;
        case 'relationship_created':
          setNotification('New relationship added');
          break;
        case 'event_created':
          setNotification(`New event: ${message.data.title}`);
          break;
      }
    },
  });
  
  return (
    <div>
      {notification && <Toast message={notification} />}
      {/* Family tree content */}
    </div>
  );
}
```

### Connection Lifecycle Tracking
Monitor connection state changes:

```typescript
function ConnectionMonitor() {
  const [status, setStatus] = React.useState('Connecting...');
  
  useWebSocket({
    onConnect: () => setStatus('Connected'),
    onDisconnect: () => setStatus('Disconnected'),
    onError: () => setStatus('Error'),
  });
  
  return <div>Status: {status}</div>;
}
```

### Manual Connection Control
Control connection manually:

```typescript
function SettingsPage() {
  const [enabled, setEnabled] = React.useState(true);
  const { isConnected, disconnect, reconnect } = useWebSocket({
    enabled,
  });
  
  return (
    <div>
      <label>
        <input
          type="checkbox"
          checked={enabled}
          onChange={(e) => setEnabled(e.target.checked)}
        />
        Enable real-time updates
      </label>
      
      <button onClick={disconnect} disabled={!isConnected}>
        Disconnect
      </button>
      
      <button onClick={reconnect} disabled={isConnected}>
        Reconnect
      </button>
    </div>
  );
}
```

## Configuration Options

### Reconnection Settings
Customize reconnection behavior:

```typescript
useWebSocket({
  reconnectInterval: 5000,      // Wait 5 seconds before reconnecting
  maxReconnectAttempts: 10,     // Try up to 10 times
  onDisconnect: () => {
    console.log('Disconnected, will retry...');
  },
});
```

### Conditional Connection
Enable/disable based on conditions:

```typescript
const { user } = useAuthStore();
const isOnRealtimePage = location.pathname.includes('/dashboard');

useWebSocket({
  enabled: user?.role === 'admin' && isOnRealtimePage,
});
```

## Cache Invalidation

The hook automatically invalidates React Query caches based on event types:

| Event Type | Invalidated Queries | Effect |
|-----------|-------------------|--------|
| `member_created` | `members`, `familyTree` | Member list and tree refresh |
| `member_updated` | `members`, `familyTree` | Member data updates |
| `member_deleted` | `members`, `familyTree` | Member removed from UI |
| `relationship_created` | `relationships`, `familyTree` | New connection shown |
| `relationship_deleted` | `relationships`, `familyTree` | Connection removed |
| `event_created` | `events` | Event list refreshes |
| `event_updated` | `events` | Event data updates |
| `event_deleted` | `events` | Event removed from list |

### How It Works
1. WebSocket message received
2. Message type parsed
3. Relevant query keys invalidated
4. React Query refetches data
5. Components re-render with new data

### No Manual Invalidation Needed
```typescript
// ❌ Don't do this - it's automatic!
const queryClient = useQueryClient();
useWebSocket({
  onMessage: (message) => {
    if (message.type === 'member_created') {
      queryClient.invalidateQueries(['members']); // Already done!
    }
  },
});

// ✅ Just use the hook
useWebSocket();
```

## Best Practices

### 1. Single Connection Per App
Use WebSocket once at the app level, not in every component:

```typescript
// ✅ Good - Single connection
function App() {
  useWebSocket();
  return <YourApp />;
}

// ❌ Bad - Multiple connections
function ComponentA() {
  useWebSocket();
  return <div>A</div>;
}
function ComponentB() {
  useWebSocket();
  return <div>B</div>;
}
```

### 2. Conditional Connection for Performance
Only connect when needed:

```typescript
// ✅ Good - Only on relevant pages
const needsRealtime = ['/dashboard', '/family-tree'].includes(location.pathname);
useWebSocket({ enabled: needsRealtime });

// ❌ Bad - Always connected
useWebSocket();
```

### 3. Handle Connection Errors
Provide feedback to users:

```typescript
const [error, setError] = React.useState('');

useWebSocket({
  onError: () => {
    setError('Connection lost. Retrying...');
  },
  onConnect: () => {
    setError('');
  },
});

{error && <ErrorBanner message={error} />}
```

### 4. Clean Up on Logout
The hook automatically disconnects when user logs out (no manual cleanup needed):

```typescript
// ✅ Automatic cleanup
const { clearAuth } = useAuthStore();
clearAuth(); // WebSocket disconnects automatically

// ❌ No need for manual disconnect
const { disconnect } = useWebSocket();
clearAuth();
disconnect(); // Redundant!
```

## Troubleshooting

### Connection Not Establishing
1. Check authentication: `useAuthStore().accessToken` should be set
2. Verify WebSocket URL in `src/utils/constants.ts`
3. Check backend is running and WebSocket endpoint is available
4. Inspect browser console for connection errors

### Messages Not Received
1. Verify backend is broadcasting messages
2. Check message format matches expected structure
3. Enable debug logging:
```typescript
useWebSocket({
  onMessage: (message) => {
    console.log('Received:', message);
  },
});
```

### Cache Not Invalidating
1. Verify event type matches expected values
2. Check React Query DevTools to see invalidation
3. Ensure queries are using correct query keys

### Reconnection Not Working
1. Check `maxReconnectAttempts` setting
2. Verify backend accepts reconnections
3. Monitor reconnection attempts:
```typescript
let attempts = 0;
useWebSocket({
  onDisconnect: () => {
    attempts++;
    console.log(`Reconnect attempt ${attempts}`);
  },
  onConnect: () => {
    attempts = 0;
  },
});
```

## Environment Configuration

### Development
```env
# .env.development
VITE_WS_URL=ws://localhost:8080/ws
```

### Production
```env
# .env.production
VITE_WS_URL=wss://api.vamsasetu.com/ws
```

### Using Environment Variables
```typescript
// src/utils/constants.ts
export const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';
```

## Security Considerations

### 1. JWT Token Authentication
Token is passed as query parameter:
```
ws://localhost:8080/ws?token=<jwt_token>
```

### 2. Secure WebSocket (WSS)
Always use `wss://` in production:
```typescript
const WS_URL = process.env.NODE_ENV === 'production'
  ? 'wss://api.vamsasetu.com/ws'
  : 'ws://localhost:8080/ws';
```

### 3. Token Refresh
Hook automatically disconnects and reconnects when token refreshes:
```typescript
// Handled automatically by authStore
const { setAuth } = useAuthStore();
setAuth(user, newAccessToken, refreshToken); // WebSocket reconnects
```

## Performance Tips

### 1. Debounce Rapid Updates
If receiving many updates, debounce cache invalidation:
```typescript
const debouncedInvalidate = debounce(() => {
  queryClient.invalidateQueries(['members']);
}, 1000);

useWebSocket({
  onMessage: (message) => {
    if (message.type === 'member_updated') {
      debouncedInvalidate();
    }
  },
});
```

### 2. Selective Invalidation
Invalidate specific items instead of entire lists:
```typescript
useWebSocket({
  onMessage: (message) => {
    if (message.type === 'member_updated') {
      // Invalidate specific member
      queryClient.invalidateQueries(['member', message.data.id]);
    }
  },
});
```

### 3. Conditional Connection
Only connect when user is active:
```typescript
const [isActive, setIsActive] = React.useState(true);

React.useEffect(() => {
  const handleVisibilityChange = () => {
    setIsActive(!document.hidden);
  };
  
  document.addEventListener('visibilitychange', handleVisibilityChange);
  return () => {
    document.removeEventListener('visibilitychange', handleVisibilityChange);
  };
}, []);

useWebSocket({ enabled: isActive });
```

## Testing

### Mocking WebSocket in Tests
```typescript
// Mock WebSocket
global.WebSocket = jest.fn().mockImplementation(() => ({
  readyState: WebSocket.OPEN,
  close: jest.fn(),
  send: jest.fn(),
}));

// Test component
render(<ComponentWithWebSocket />);
```

### Testing Message Handling
```typescript
const mockMessage = {
  type: 'member_created',
  data: { id: '1', name: 'Test' },
  timestamp: Date.now(),
};

// Simulate message
const ws = new WebSocket('ws://test');
ws.onmessage?.(new MessageEvent('message', {
  data: JSON.stringify(mockMessage),
}));
```

## Related Documentation
- [useAuth Hook](./useAuth.ts) - Authentication management
- [useMembers Hook](./useMembers.ts) - Member data fetching
- [React Query Docs](https://tanstack.com/query/latest) - Query client usage
- [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket) - Browser WebSocket API
