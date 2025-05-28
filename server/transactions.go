package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type Transaction struct {
	ID       string `json:"id"`
	ClientID string `json:"clientid"`
	Status   string `json:"status"`
	Details  string `json:"details,omitempty"`
	Message  string `json:"message,omitempty"`
}

func (s *Server) HandleTransactions(w http.ResponseWriter, r *http.Request) {

	// Ensure the body is closed after reading
	defer r.Body.Close()

	// Read the entire request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	txn := new(Transaction)
	err = json.Unmarshal(body, txn)
	log.Println(txn.ClientID, txn.ID)

	if err != nil {
		log.Println("unmarshalling err", err)
	}

	ClientID := strings.TrimSpace(txn.ClientID)

	if ClientID != "" {
		if client, ok := s.clients[ClientID]; ok {
			log.Println("clientid", ClientID)
			client.mu.Lock()
			client.send <- body
			client.mu.Unlock()
		} else {
			log.Println("client not registerd", ClientID)
		}
	}
}
