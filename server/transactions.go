package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Transaction struct {
	ID      string `json:"id"`
	WID     string `json:"workerid"`
	CID     string `json:"clientid"`
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
	Message string `json:"message,omitempty"`
}
type BroadCast struct {
	ID      string `json:"id"`
	WID     string `json:"workerid"`
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
	Message string `json:"message,omitempty"`
}

func (s *Server) HandleTransactions(c *gin.Context) {

	// Ensure the body is closed after reading
	defer c.Request.Body.Close()

	// Read the entire request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		http.Error(c.Writer, "Unable to read request body", http.StatusBadRequest)
		return
	}

	txn := new(Transaction)
	err = json.Unmarshal(body, txn)
	log.Println(txn.CID, txn.WID)

	if err != nil {
		log.Println("unmarshalling err", err)
		return
	}

	if txn.WID == "" || txn.CID == "" {
		log.Println("Worker id and Cleint id needed.")
		return
	}

	client := s.GetRegisteredWorker(txn.WID).GetClient(txn.CID)

	if client == nil {
		log.Println("client not registerd", txn.CID)
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	client.send <- body
}

func (s *Server) HandleBroadCast(c *gin.Context) {

	// Ensure the body is closed after reading
	defer c.Request.Body.Close()

	// Read the entire request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		http.Error(c.Writer, "Unable to read request body", http.StatusBadRequest)
		return
	}

	bx := new(BroadCast)
	err = json.Unmarshal(body, bx)
	log.Println(bx.WID)

	if err != nil {
		log.Println("unmarshalling err", err)
		return
	}

	if bx.WID == "" {
		log.Println("Worker id needed.")
		return
	}

	worker := s.GetRegisteredWorker(bx.WID)
	if worker == nil {
		log.Println("worker not registerd", bx.WID)
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	worker.broadcast <- body
}
