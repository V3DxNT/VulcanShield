package websocket

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}


type EventPayload struct {
	EventType string    `json:"event_type"` 
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}


type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
	log     *slog.Logger
}


func NewHub(log *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
		log:     log,
	}
}


func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Warn("websocket upgrade failed", "error", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	h.log.Info("websocket client connected", "remote_addr", r.RemoteAddr)

	
	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		_ = conn.Close()
		h.log.Info("websocket client disconnected", "remote_addr", r.RemoteAddr)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}


func (h *Hub) Broadcast(eventType string, data any) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.clients) == 0 {
		return
	}

	event := EventPayload{
		EventType: eventType,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}

	msg, err := json.Marshal(event)
	if err != nil {
		return
	}

	for conn := range h.clients {
		_ = conn.WriteMessage(websocket.TextMessage, msg)
	}
}
