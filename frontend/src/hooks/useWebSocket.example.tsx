import React from 'react';
import { useWebSocket } from './useWebSocket';

/**
 * Example 1: Basic WebSocket connection
 * Automatically connects when authenticated and invalidates cache
 */
export function BasicWebSocketExample() {
  const { isConnected } = useWebSocket();

  return (
    <div>
      <p>WebSocket Status: {isConnected ? 'Connected' : 'Disconnected'}</p>
    </div>
  );
}

/**
 * Example 2: WebSocket with custom message handler
 * Handle specific message types with custom logic
 */
export function CustomMessageHandlerExample() {
  const [lastMessage, setLastMessage] = React.useState<string>('');

  const { isConnected } = useWebSocket({
    onMessage: (message) => {
      console.log('Received message:', message);
      setLastMessage(`${message.type} at ${new Date(message.timestamp * 1000).toLocaleTimeString()}`);
    },
  });

  return (
    <div>
      <p>Status: {isConnected ? '🟢 Connected' : '🔴 Disconnected'}</p>
      {lastMessage && <p>Last update: {lastMessage}</p>}
    </div>
  );
}

/**
 * Example 3: WebSocket with connection callbacks
 * Track connection state changes
 */
export function ConnectionCallbacksExample() {
  const [connectionLog, setConnectionLog] = React.useState<string[]>([]);

  const addLog = (message: string) => {
    setConnectionLog((prev) => [...prev, `${new Date().toLocaleTimeString()}: ${message}`]);
  };

  const { isConnected, disconnect, reconnect } = useWebSocket({
    onConnect: () => addLog('Connected to WebSocket'),
    onDisconnect: () => addLog('Disconnected from WebSocket'),
    onError: (error) => addLog(`Error: ${error.type}`),
  });

  return (
    <div>
      <div>
        <p>Status: {isConnected ? 'Connected' : 'Disconnected'}</p>
        <button onClick={disconnect} disabled={!isConnected}>
          Disconnect
        </button>
        <button onClick={reconnect} disabled={isConnected}>
          Reconnect
        </button>
      </div>
      <div>
        <h3>Connection Log:</h3>
        <ul>
          {connectionLog.map((log, index) => (
            <li key={index}>{log}</li>
          ))}
        </ul>
      </div>
    </div>
  );
}

/**
 * Example 4: Conditional WebSocket connection
 * Only connect when needed (e.g., on specific pages)
 */
export function ConditionalConnectionExample() {
  const [enableWebSocket, setEnableWebSocket] = React.useState(false);

  const { isConnected } = useWebSocket({
    enabled: enableWebSocket,
  });

  return (
    <div>
      <label>
        <input
          type="checkbox"
          checked={enableWebSocket}
          onChange={(e) => setEnableWebSocket(e.target.checked)}
        />
        Enable real-time updates
      </label>
      <p>Status: {isConnected ? 'Connected' : 'Disconnected'}</p>
    </div>
  );
}

/**
 * Example 5: WebSocket with custom reconnect settings
 * Configure reconnection behavior
 */
export function CustomReconnectExample() {
  const [reconnectAttempts, setReconnectAttempts] = React.useState(0);

  const { isConnected } = useWebSocket({
    reconnectInterval: 5000, // 5 seconds
    maxReconnectAttempts: 10,
    onMessage: () => {
      // Handle messages if needed
    },
    onDisconnect: () => {
      setReconnectAttempts((prev) => prev + 1);
    },
    onConnect: () => {
      setReconnectAttempts(0);
    },
  });

  return (
    <div>
      <p>Status: {isConnected ? 'Connected' : 'Disconnected'}</p>
      {reconnectAttempts > 0 && (
        <p>Reconnect attempts: {reconnectAttempts}</p>
      )}
    </div>
  );
}

/**
 * Example 6: Real-time notification indicator
 * Show a visual indicator when updates are received
 */
export function RealtimeNotificationExample() {
  const [hasUpdate, setHasUpdate] = React.useState(false);

  const { isConnected } = useWebSocket({
    onMessage: () => {
      setHasUpdate(true);
      setTimeout(() => setHasUpdate(false), 3000); // Clear after 3 seconds
    },
  });

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
        <span>{isConnected ? '🟢' : '🔴'}</span>
        {hasUpdate && (
          <span style={{ animation: 'pulse 1s ease-in-out' }}>
            ✨ New update received!
          </span>
        )}
      </div>
    </div>
  );
}

/**
 * Example 7: Integration with App component
 * Set up WebSocket at the app level for global real-time updates
 */
export function AppLevelWebSocketExample() {
  // In your App.tsx or main layout component:
  useWebSocket({
    onMessage: (message) => {
      // Optional: Show toast notifications for specific events
      if (message.type === 'member_created') {
        console.log('New family member added!');
      }
    },
  });

  return (
    <div>
      {/* Your app content */}
      <p>WebSocket is active in the background</p>
    </div>
  );
}
