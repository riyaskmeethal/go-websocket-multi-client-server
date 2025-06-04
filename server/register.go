package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var activeConnections int64

func (s *Server) RegisterClient(c *gin.Context) {

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	// Optional: ping client periodically to keep alive
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wid := c.Param("wid")
	cid := c.Param("cid")

	if wid == "" || cid == "" {
		log.Println("Worker id and Cleint id needed.")
		conn.WriteJSON("Worker id and Cleint id needed.")
		conn.Close()
		return
	}

	worker := s.GetWorker(wid)
	if client := worker.GetClient(cid); client == nil {

		client = NewClient(wid, cid)
		client.conn = conn

		worker.register <- client
		atomic.AddInt64(&activeConnections, 1)
		defer s.CleanWorker(wid, client)

		go PingClient(ctx, client)

		log.Println("Listening for client communications")
		for {
			_, m, err := client.conn.ReadMessage()
			if err != nil {
				break
			}
			fmt.Println(string(m))
		}

	} else {
		log.Println("Allready registered")
		conn.WriteJSON("Allready registered")
		conn.Close()
		return
	}
}

func PingClient(ctx context.Context, c *Client) {

	log.Println("client ping service running")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down ping service for client : ", c.wid, " : ", c.cid)
			return
		case <-ticker.C:
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				log.Println("Ping failed:", err)
				return // connection likely dead
			}
		}
	}
}
