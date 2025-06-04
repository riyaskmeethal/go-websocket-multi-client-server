package server

import (
	"context"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	wid    string
	cid    string
	cancel context.CancelFunc
	conn   *websocket.Conn
	send   chan []byte
	mu     sync.Mutex
}

func NewClient(wid, cid string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		wid:    wid,
		cid:    cid,
		cancel: cancel,
		send:   make(chan []byte, 10),
		mu:     sync.Mutex{},
	}
	go client.SendMessages(ctx)
	return client

}

func (c *Client) SendMessages(ctx context.Context) {

	log.Println("client messsage manager running")

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down send message listner for client : ", c.wid, " : ", c.cid)
			return
		case msg := <-c.send:
			// Send the message to the client
			c.mu.Lock()
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			c.mu.Unlock()
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

}
