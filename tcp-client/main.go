package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error dialing: %s", err.Error())
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Fatalf("Error setting deadline: %s", err.Error())
	}

	msg := []byte("Hello from client\n")

	_, err = conn.Write(msg)
	if err != nil {
		log.Fatalf("Error writing: %s", err.Error())
	}

	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		log.Fatalf("Error reading: %s", err.Error())
	}

	fmt.Printf("Response from server: %s", string(response))
}
