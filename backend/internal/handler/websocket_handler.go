package handler

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// Client represents a WebSocket client
type Client struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s (User: %s)", client.ID, client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Printf("Client unregistered: %s (User: %s)", client.ID, client.UserID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastUpdate sends an update to all connected clients
func (h *Hub) BroadcastUpdate(eventType string, data interface{}) {
	message := map[string]interface{}{
		"type":      eventType,
		"data":      data,
		"timestamp": time.Now().Unix(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	h.broadcast <- jsonMessage
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(hub *Hub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// Extract user ID from query params or locals (set by auth middleware)
		userID := c.Query("userId", "anonymous")
		if localUserID := c.Locals("userID"); localUserID != nil {
			if uid, ok := localUserID.(string); ok {
				userID = uid
			}
		}

		client := &Client{
			ID:     c.Params("id", "unknown"),
			UserID: userID,
			Conn:   c,
			Send:   make(chan []byte, 256),
		}

		hub.register <- client

		// Start goroutines for reading and writing
		go client.writePump()
		client.readPump(hub)
	})
}

// readPump reads messages from the WebSocket connection
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Echo message back (can be extended for client-to-server messages)
		log.Printf("Received message from client %s: %s", c.ID, string(message))
	}
}

// writePump writes messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
