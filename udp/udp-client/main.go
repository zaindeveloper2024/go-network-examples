package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		fmt.Printf("Failed to resolve server address: %v\n", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Printf("Failed to create connection: %v\n", err)
		return
	}
	defer conn.Close()

	message := fmt.Sprintf("Hello from client at %v", time.Now())
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
		os.Exit(1)
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Printf("Failed to receive response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Server response: %s\n", string(buffer[:n]))
}
