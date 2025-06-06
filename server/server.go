package server

import (
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	Workers      *sync.Map
	removeWorker chan string
}

func NewServer() *Server {
	s := &Server{
		Workers:      &sync.Map{},
		removeWorker: make(chan string, 10),
	}

	go s.RemoveWorker()
	return s
}

func (s *Server) RemoveWorker() {

	log.Println("Worker observer running")
	for wid := range s.removeWorker {
		if w := s.GetRegisteredWorker(wid); w != nil {
			w.cancel()
			s.Workers.Delete(wid)
			log.Println("Worker Removed.", wid)
		}
	}
}

func (s *Server) GetWorker(wid string) *Worker {
	if v, workerExist := s.Workers.Load(wid); workerExist {
		if w, ok := v.(*Worker); ok {
			log.Println("existing worker")
			return w
		}
	}
	newWorker := NewWorker(wid)
	s.Workers.Store(wid, newWorker)
	log.Println("New worker created.")
	return newWorker
}

func (s *Server) GetRegisteredWorker(wid string) (worker *Worker) {
	if v, workerExist := s.Workers.Load(wid); workerExist {
		if worker, ok := v.(*Worker); ok {
			return worker
		}
	}
	return
}

func StatsPrinter() {
	for {
		log.Printf("Active connections: %d | Goroutines: %d\n",
			atomic.LoadInt64(&activeConnections),
			runtime.NumGoroutine())
		time.Sleep(30 * time.Second)
	}
}

func (s *Server) CleanWorker(wid string, c *Client) {

	log.Println("Cleaning worker : ", wid)

	if w := s.GetRegisteredWorker(wid); w != nil {
		atomic.AddInt64(&w.clientsCount, -1)
		w.unregister <- c
		if w.clientsCount < 1 {
			log.Println("removing worker")
			s.removeWorker <- wid
		}
	}
}

type Message struct {
	Status  string
	Content any
}
