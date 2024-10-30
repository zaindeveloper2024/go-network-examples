package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// go run .
// go run . localhost:9090
func main() {
	serverAddr := "localhost:8080"
	if len(os.Args) > 1 {
		serverAddr = os.Args[1]
	}

	err := serverRun(serverAddr)
	if err != nil {
		log.Fatalf("Error running server: %v\n", err)
	}
}

func serverRun(serverAddr string) error {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("error listening: %w", err)
	}
	defer listener.Close()

	log.Printf("Server is listening on %s\n", serverAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting a connection request: %v\n", err)
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

		msg, err := reader.ReadString('\n')
		if err != nil {
			handleReadError(err, clientAddr)
			return
		}

		if len(msg) > maxMsgSize {
			log.Printf("Message too big from %s\n", clientAddr)
			return
		}

		response := processMessage(msg)

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
		return fmt.Errorf("unexpected EOF from %s", clientAddr)
	default:
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("connection timeout from %s", clientAddr)
		}
		return fmt.Errorf("error reading from %s: %v", clientAddr, err)
	}
}

func processMessage(msg string) string {
	return fmt.Sprintf("Message received: %s\n", msg)
}
