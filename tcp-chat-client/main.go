package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Printf("Error dialing: %s", err.Error())
		return
	}

	defer conn.Close()

	go receiveMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter message (or 'quit' to exit):")
		if !scanner.Scan() {
			log.Printf("Error reading from stdin: %s", scanner.Err())
			break
		}

		text := scanner.Text()
		if strings.ToLower(text) == "quit" {
			log.Println("Exiting...")
			return
		}

		msg := Message{
			Type:    "chat",
			Content: text,
		}

		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(msg); err != nil {
			fmt.Println("Error encoding message:", err)
			return
		}
	}
}

func receiveMessages(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			fmt.Println("Connection closed by server")
			os.Exit(0)
		}

		fmt.Printf("\nReceived: %+v\n", msg)
		fmt.Print("Enter message (or 'quit' to exit): ")
	}
}
