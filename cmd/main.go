package main

import (
	"fmt"
	"go-websocket-multi-client-server/server"
	"log"
	"net/http"
)

func main() {

	server := server.NewServer()

	// Set up HTTP handler to handle WebSocket connections
	http.HandleFunc("/ws", server.RegisterClient)
	http.HandleFunc("/transaction", server.HandleTransactions)

	// Start the WebSocket server in a goroutine
	go server.Run()

	// Start HTTP server
	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
