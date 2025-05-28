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
	if err := http.ListenAndServeTLS(":8080", "localhost.pem", "localhost-key.pem", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
