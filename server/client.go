package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func (s *Server) RegisterClient(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientID := strings.TrimSpace(r.URL.Query().Get("id"))

	if clientID == "" {
		log.Println("Cleint id needed.")
		return
	}
	if _, ok := s.clients[clientID]; !ok {

		Client := NewClient()
		Client.id = clientID
		Client.conn = conn

		s.register <- Client

		log.Println("clients", len(s.clients))

		// routine for sending messages to the client
		go s.SendMessages(Client)

		// Listen for messages from the client
		for {
			t, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			log.Println(string(msg))
			conn.WriteMessage(t, msg)
		}

		s.unregister <- Client

	} else {
		log.Println("Client already registerd with id ", clientID)
	}

}
