package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	serverAddr := "localhost:8080"

	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server is listening on %s\n", serverAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v\n", err)
		}
	}()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Accepted connection from %s\n", clientAddr)

	const (
		readTimeout  = 30 * time.Second
		writeTimeout = 10 * time.Second
		maxMsgSize   = 1024 * 1024
	)

	reader := bufio.NewReader(conn)

	for {
		if err := conn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
			log.Printf("Error setting read deadline for %s: %v\n", clientAddr, err)
			return
		}

		message, err := reader.ReadString('\n')
		if err != nil {
			handleReadError(err, clientAddr)
			return
		}

		if len(message) > maxMsgSize {
			log.Printf("Message too big from %s\n", clientAddr)
			return
		}

		response := processMessage(message)

		if err := conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
			log.Printf("Error setting write deadline for %s: %v\n", clientAddr, err)
			return
		}

		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Error writing to %s: %v\n", clientAddr, err)
			return
		}
	}
}

func handleReadError(err error, clientAddr string) error {
	switch {
	case errors.Is(err, io.EOF):
		log.Printf("Client %s disconnected normally\n", clientAddr)
		return err
	case errors.Is(err, io.ErrUnexpectedEOF):
		return fmt.Errorf("Unexpected EOF from %s", clientAddr)
	default:
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("Connection timeout from %s", clientAddr)
		}
		return fmt.Errorf("Error reading from %s: %v", clientAddr, err)
	}
}

func processMessage(message string) string {
	return fmt.Sprintf("Message received: %s\n", message)
}
