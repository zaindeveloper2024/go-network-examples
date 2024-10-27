package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type Server struct {
	clients    map[net.Conn]bool
	clientsMux sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		clients: make(map[net.Conn]bool),
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	server := NewServer()

	fmt.Printf("Server is listening on %s\n", ":8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}

		go server.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v\n", err)
		}
	}()

	s.clientsMux.Lock()
	s.clients[conn] = true
	s.clientsMux.Unlock()

	clientAddr := conn.RemoteAddr().String()

	log.Printf("Accepted connection from %s\n", clientAddr)
	defer func() {
		s.clientsMux.Lock()
		delete(s.clients, conn)
		s.clientsMux.Unlock()
	}()

	decoder := json.NewDecoder(conn)
	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			fmt.Printf("Client disconnected: %v\n", err)
			return
		}

		fmt.Printf("Received message: %v\n", msg)

		s.broadcast(msg, conn)
	}
}

func (s *Server) broadcast(msg Message, sender net.Conn) {
	s.clientsMux.RLock()
	defer s.clientsMux.RUnlock()

	for client := range s.clients {
		if client != sender {
			encoder := json.NewEncoder(client)
			if err := encoder.Encode(msg); err != nil {
				fmt.Printf("Error broadcasting to client: %v\n", err)
			}
		}
	}
}
