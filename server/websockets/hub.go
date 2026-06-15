package websockets

import (
	"log"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	IDs   []string
	Conns []*websocket.Conn
	mu    sync.Mutex // guards writes
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client // keyed by chat ID
}

func NewHub() *Hub {
	return &Hub{clients: make(map[string]*Client)}
}

func (h *Hub) Register(chatID, userID string, conn *websocket.Conn) *Client {
	h.mu.Lock()
	if _, ok := h.clients[chatID]; ok {
		if slices.Contains(h.clients[chatID].IDs, userID) {
			h.mu.Unlock()
			return h.clients[chatID]
		}
		// check if id in IDs
		h.clients[chatID].mu.Lock()
		index := slices.Index(h.clients[chatID].IDs, userID)
		if index != -1 {
			additionalConnection := h.clients[chatID].Conns[index]
			additionalConnection.Close()
			h.clients[chatID].Conns[index] = conn
			h.clients[chatID].mu.Unlock()
			return h.clients[chatID]
		}

		h.clients[chatID].IDs = append(h.clients[chatID].IDs, userID)
		h.clients[chatID].Conns = append(h.clients[chatID].Conns, conn)
		h.clients[chatID].mu.Unlock()
		h.mu.Unlock()
		return h.clients[chatID]
	}
	c := &Client{IDs: []string{userID}, Conns: []*websocket.Conn{conn}}
	h.clients[chatID] = c
	h.mu.Unlock()
	return c
}

func (h *Hub) Unregister(chatID, id string) {
	h.mu.Lock()
	if _, ok := h.clients[chatID]; !ok {
		h.mu.Unlock()
		log.Println("Could not find chatID in chat connections")
		return
	}
	h.clients[chatID].mu.Lock()
	index := slices.Index(h.clients[chatID].IDs, id)
	if index == -1 {
		h.clients[chatID].mu.Unlock()
		h.mu.Unlock()
		log.Println("Could not find userID in chat connections")
		return
	}
	h.clients[chatID].IDs = slices.Delete(h.clients[chatID].IDs, index, index+1)
	h.clients[chatID].Conns = slices.Delete(h.clients[chatID].Conns, index, index+1)
	h.clients[chatID].mu.Unlock()
	log.Println(h.clients[chatID].IDs)
	h.mu.Unlock()
}

// Send to a specific list of user IDs
func (h *Hub) SendToIDs(chatID string, msgType int, msg []byte) {
	target := h.clients[chatID]
	target.mu.Lock() // lock once, outside the loop
	defer target.mu.Unlock()
	for _, conn := range target.Conns {
		err := conn.WriteMessage(msgType, msg)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
