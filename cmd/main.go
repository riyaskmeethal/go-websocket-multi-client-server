package main

import (
	"fmt"
	"go-websocket-multi-client-server/server"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	s := server.NewServer()

	router := gin.Default()

	router.GET("/subscribe/:wid/:cid", s.RegisterClient)
	router.POST("/transaction", s.HandleTransactions)
	router.POST("/broadcast", s.HandleBroadCast)

	go server.StatsPrinter()

	// Start HTTP server
	fmt.Println("Server started on :80")
	if err := router.Run(":80"); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
