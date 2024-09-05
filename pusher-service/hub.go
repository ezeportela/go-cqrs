package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.onConnect(client)
		case client := <-h.unregister:
			h.onDisconnect(client)
		}
	}
}

func (h *Hub) onConnect(client *Client) {
	log.Println("Client connected", client.socket.RemoteAddr())

	h.mutex.Lock()
	defer h.mutex.Unlock()
	client.id = client.socket.RemoteAddr().String()
	h.clients = append(h.clients, client)
}

func (h *Hub) onDisconnect(client *Client) {
	log.Println("Client disconnected", client.socket.RemoteAddr())

	client.Close()
	h.mutex.Lock()
	defer h.mutex.Unlock()

	i := -1
	for j, c := range h.clients {
		if c.id == client.id {
			i = j
			break
		}
	}

	copy(h.clients[i:], h.clients[i+1:])
	h.clients[len(h.clients)-1] = nil
	h.clients = h.clients[:len(h.clients)-1]
}

func (h *Hub) Broadcast(message any, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range h.clients {
		if client != ignore {
			client.outbound <- data
		}
	}
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error upgrading connection", http.StatusInternalServerError)
		return
	}

	client := NewClient(h, socket)
	h.register <- client

	go client.Write()
}
