package server

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	id   string
	conn *websocket.Conn
	send chan []byte
	mu   sync.Mutex
}

func NewClient() *Client {
	return &Client{
		send: make(chan []byte),
		mu:   sync.Mutex{},
	}
}

type Server struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[string]*Client, 0),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

}

// Start the WebSocket server to handle new connections
func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			log.Println("registering", client.id)

			// Add new client to the client map
			s.clients[client.id] = client
			log.Println("registered", s.clients)

		case client := <-s.unregister:
			// Remove the client from the map and close the connection
			delete(s.clients, client.id)
			client.mu.Lock()
			close(client.send)
			client.conn.Close()
			client.mu.Unlock()

		case msg := <-s.broadcast:
			// Send the broadcast message to all connected clients
			for _, client := range s.clients {
				client.mu.Lock()
				client.send <- msg
				client.mu.Unlock()
			}
		}
	}
}

func (s *Server) SendMessages(client *Client) {

	for msg := range client.send {
		// Send the message to the client
		client.mu.Lock()
		err := client.conn.WriteMessage(websocket.TextMessage, msg)
		client.mu.Unlock()
		if err != nil {
			log.Println(err)
			return
		}

	}
}
