package websockets

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	mu   sync.Mutex // guards writes
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client // keyed by user ID
}

func NewHub() *Hub {
	return &Hub{clients: make(map[string]*Client)}
}

func (h *Hub) Register(id string, conn *websocket.Conn) *Client {
	c := &Client{ID: id, Conn: conn}
	h.mu.Lock()
	h.clients[id] = c
	h.mu.Unlock()
	return c
}

func (h *Hub) Unregister(id string) {
	h.mu.Lock()
	delete(h.clients, id)
	h.mu.Unlock()
}

// Send to a specific list of user IDs
func (h *Hub) SendToIDs(ids []string, msgType int, msg []byte) {
	h.mu.RLock()
	targets := make([]*Client, 0, len(ids))
	for _, id := range ids {
		if c, ok := h.clients[id]; ok {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()

	for _, c := range targets {
		c.mu.Lock()
		c.Conn.WriteMessage(msgType, msg)
		c.mu.Unlock()
	}
}
