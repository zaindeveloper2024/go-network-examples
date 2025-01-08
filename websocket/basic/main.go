package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("ReadMessage: %v", err)
			break
		}

		log.Printf("Received: %s", message)

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Fatalf("WriteMessage: %v", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("WebSocket server is running on localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
