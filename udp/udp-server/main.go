package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer conn.Close()

	fmt.Println("Server is running on port 8080")

	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		message := string(buffer[:n])
		fmt.Printf("Received %s from %s\n", message, remoteAddr)

		response := fmt.Sprintf("Received %v", time.Now())
		_, err = conn.WriteToUDP([]byte(response), remoteAddr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
