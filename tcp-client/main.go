package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Printf("Error dialing: %s", err.Error())
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("Hello from client\n"))
	if err != nil {
		log.Printf("Error writing: %s", err.Error())
		return
	}

	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		log.Printf("Error reading: %s", err.Error())
		return
	}

	fmt.Printf("Response from server: %s", string(response))
}
