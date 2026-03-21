import { useEffect, useRef, useCallback } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useAuthStore } from '../stores/authStore';
import { WS_URL } from '../utils/constants';

interface WebSocketMessage {
  type: string;
  data: any;
  timestamp: number;
}

interface UseWebSocketOptions {
  enabled?: boolean;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
  onMessage?: (message: WebSocketMessage) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  onError?: (error: Event) => void;
}

/**
 * WebSocket hook for real-time updates
 * Automatically connects when authenticated and invalidates React Query cache
 * based on event types
 */
export function useWebSocket(options: UseWebSocketOptions = {}) {
  const {
    enabled = true,
    reconnectInterval = 3000,
    maxReconnectAttempts = 5,
    onMessage,
    onConnect,
    onDisconnect,
    onError,
  } = options;

  const { accessToken, isAuthenticated } = useAuthStore();
  const queryClient = useQueryClient();
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const shouldConnectRef = useRef(enabled && isAuthenticated);

  // Handle incoming WebSocket messages
  const handleMessage = useCallback(
    (event: MessageEvent) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        
        // Call custom message handler if provided
        if (onMessage) {
          onMessage(message);
        }

        // Invalidate React Query cache based on event type
        switch (message.type) {
          case 'member_created':
          case 'member_updated':
          case 'member_deleted':
            queryClient.invalidateQueries({ queryKey: ['members'] });
            queryClient.invalidateQueries({ queryKey: ['familyTree'] });
            break;

          case 'relationship_created':
          case 'relationship_deleted':
            queryClient.invalidateQueries({ queryKey: ['relationships'] });
            queryClient.invalidateQueries({ queryKey: ['familyTree'] });
            break;

          case 'event_created':
          case 'event_updated':
          case 'event_deleted':
            queryClient.invalidateQueries({ queryKey: ['events'] });
            break;

          default:
            console.log('Unknown WebSocket event type:', message.type);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    },
    [queryClient, onMessage]
  );

  // Connect to WebSocket
  const connect = useCallback(() => {
    if (!shouldConnectRef.current || !accessToken) {
      return;
    }

    // Close existing connection if any
    if (wsRef.current) {
      wsRef.current.close();
    }

    try {
      // Connect with JWT token as query parameter
      const wsUrl = `${WS_URL}?token=${accessToken}`;
      const ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        console.log('WebSocket connected');
        reconnectAttemptsRef.current = 0;
        if (onConnect) {
          onConnect();
        }
      };

      ws.onmessage = handleMessage;

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        if (onError) {
          onError(error);
        }
      };

      ws.onclose = () => {
        console.log('WebSocket disconnected');
        if (onDisconnect) {
          onDisconnect();
        }

        // Attempt to reconnect if enabled and within max attempts
        if (
          shouldConnectRef.current &&
          reconnectAttemptsRef.current < maxReconnectAttempts
        ) {
          reconnectAttemptsRef.current += 1;
          console.log(
            `Attempting to reconnect (${reconnectAttemptsRef.current}/${maxReconnectAttempts})...`
          );
          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, reconnectInterval);
        }
      };

      wsRef.current = ws;
    } catch (error) {
      console.error('Error creating WebSocket connection:', error);
    }
  }, [
    accessToken,
    handleMessage,
    reconnectInterval,
    maxReconnectAttempts,
    onConnect,
    onDisconnect,
    onError,
  ]);

  // Disconnect from WebSocket
  const disconnect = useCallback(() => {
    shouldConnectRef.current = false;
    
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }

    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
  }, []);

  // Effect to manage WebSocket connection
  useEffect(() => {
    shouldConnectRef.current = enabled && isAuthenticated;

    if (shouldConnectRef.current) {
      connect();
    } else {
      disconnect();
    }

    // Cleanup on unmount
    return () => {
      disconnect();
    };
  }, [enabled, isAuthenticated, connect, disconnect]);

  return {
    isConnected: wsRef.current?.readyState === WebSocket.OPEN,
    disconnect,
    reconnect: connect,
  };
}
