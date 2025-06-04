package server

import (
	"context"
	"log"
	"sync"
)

type Worker struct {
	wid          string
	clients      *sync.Map
	cancel       context.CancelFunc
	clientsCount int
	broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
}

func NewWorker(wid string) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	worker := &Worker{
		wid:        wid,
		cancel:     cancel,
		clients:    &sync.Map{},
		broadcast:  make(chan []byte, 10),
		register:   make(chan *Client, 1000),
		unregister: make(chan *Client, 1000),
	}
	log.Println("3 Worker jobs running ------------- job")
	go worker.AddClient(ctx)
	go worker.RemoveClient(ctx)
	go worker.BroadCast(ctx)
	return worker

}

func (w *Worker) GetClient(cid string) *Client {

	if c, clientExist := w.clients.Load(cid); clientExist {
		if client, ok := c.(*Client); ok {
			return client
		}
	}
	return nil
}

func (w *Worker) AddClient(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down add client listner for worker : ", w.wid)
			return
		case client := <-w.register:
			w.clients.Store(client.cid, client)
			w.clientsCount++
			log.Println("Client Registerd", client.wid, ":", client.cid)
		}
	}
}

func (w *Worker) RemoveClient(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down remove client listner for worker : ", w.wid)
			return
		case client := <-w.unregister:
			client.cancel()
			client.conn.Close()
			w.clients.Delete(client.cid)
			log.Println("Client Removed.", client.wid, ":", client.cid)
		}
	}
}

func (w *Worker) BroadCast(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down broadcast listner for worker : ", w.wid)
			return
		case msg := <-w.broadcast:
			w.clients.Range(func(key, value any) bool {
				if c, ok := value.(*Client); ok {
					c.mu.Lock()
					c.send <- msg
					c.mu.Unlock()
				}
				return true
			})
		}
	}

}
