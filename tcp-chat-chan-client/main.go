package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

type Client struct {
	conn     net.Conn
	nickname string
}

type ChatServer struct {
	clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

func NewServer() *ChatServer {
	return &ChatServer{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (server *ChatServer) run() {
	for {
		select {
		case client := <-server.register:
			server.mutex.Lock()
			server.clients[client] = true
			server.mutex.Unlock()
			server.broadcast <- fmt.Sprintf("%s joined the chat", client.nickname)
		case client := <-server.unregister:
			server.mutex.Lock()
			if _, ok := server.clients[client]; ok {
				delete(server.clients, client)
				close := fmt.Sprintf("%s left the chat", client.nickname)
				server.broadcast <- close
			}
			server.mutex.Unlock()
		case message := <-server.broadcast:
			server.mutex.Lock()
			for client := range server.clients {
				go func(c *Client, msg string) {
					fmt.Fprintf(c.conn, "%s\n", msg)
				}(client, message)
			}
			server.mutex.Unlock()
		}
	}
}

func (server *ChatServer) handleClient(client *Client) {
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			server.unregister <- client
			client.conn.Close()
			return
		}
		server.broadcast <- fmt.Sprintf("%s: %s", client.nickname, message)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	defer listener.Close()

	server := NewServer()
	go server.run()

	fmt.Println("Server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		client := &Client{
			conn:     conn,
			nickname: fmt.Sprintf("Anonymous-%d", len(server.clients)+1),
		}

		server.register <- client
		go server.handleClient(client)
	}
}
