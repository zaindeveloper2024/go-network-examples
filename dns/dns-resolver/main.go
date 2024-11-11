package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

type DNSQuestion struct {
	Name  string
	Type  uint16
	Class uint16
}

type DNSMessage struct {
	Header   DNSHeader
	Question DNSQuestion
}

func encodeDNSName(domain string) []byte {
	var encoded []byte
	parts := []byte(domain)
	start := 0
	for i := 0; i < len(parts); i++ {
		if parts[i] == '.' || i == len(parts)-1 {
			length := i - start
			if i == len(parts)-1 && parts[i] != '.' {
				length++
			}
			encoded = append(encoded, byte(length))
			encoded = append(encoded, parts[start:start+length]...)
			start = i + 1
		}
	}
	encoded = append(encoded, 0)
	return encoded
}

func createDNSQuery(domain string) []byte {
	message := DNSMessage{
		Header: DNSHeader{
			ID:      0x1234,
			Flags:   0x0100,
			QDCount: 1,
			ANCount: 0,
			NSCount: 0,
			ARCount: 0,
		},
		Question: DNSQuestion{
			Name:  domain,
			Type:  1,
			Class: 1,
		},
	}

	query := make([]byte, 12)

	binary.BigEndian.PutUint16(query[0:2], message.Header.ID)
	binary.BigEndian.PutUint16(query[2:4], message.Header.Flags)
	binary.BigEndian.PutUint16(query[4:6], message.Header.QDCount)
	binary.BigEndian.PutUint16(query[6:8], message.Header.ANCount)
	binary.BigEndian.PutUint16(query[8:10], message.Header.NSCount)
	binary.BigEndian.PutUint16(query[10:12], message.Header.ARCount)

	query = append(query, encodeDNSName(message.Question.Name)...)

	typeBytes := make([]byte, 2)
	classBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBytes, message.Question.Type)
	binary.BigEndian.PutUint16(classBytes, message.Question.Class)
	query = append(query, typeBytes...)
	query = append(query, classBytes...)

	return query
}

func parseDNSResponse(response []byte) {
	if len(response) < 12 {
		fmt.Println("Response too short")
		return
	}

	header := DNSHeader{
		ID:      binary.BigEndian.Uint16(response[0:2]),
		Flags:   binary.BigEndian.Uint16(response[2:4]),
		QDCount: binary.BigEndian.Uint16(response[4:6]),
		ANCount: binary.BigEndian.Uint16(response[6:8]),
		NSCount: binary.BigEndian.Uint16(response[8:10]),
		ARCount: binary.BigEndian.Uint16(response[10:12]),
	}

	fmt.Printf("Response Header: %+v\n", header)

	offset := 12
	for i := 0; i < int(header.QDCount); i++ {
		for response[offset] != 0 {
			offset += int(response[offset]) + 1
		}
		offset += 5
	}

	for i := 0; i < int(header.ANCount); i++ {
		if (response[offset] & 0xC0) == 0xC0 {
			offset += 2
		} else {
			for response[offset] != 0 {
				offset += int(response[offset]) + 1
			}
			offset++
		}

		rtype := binary.BigEndian.Uint16(response[offset:])
		offset += 8
		rdLength := binary.BigEndian.Uint16(response[offset:])
		offset += 2

		if rtype == 1 && rdLength == 4 {
			ip := net.IPv4(response[offset], response[offset+1], response[offset+2], response[offset+3])
			fmt.Printf("Found A Record: %s\n", ip.String())
		}
		offset += int(rdLength)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go example.com")
		return
	}

	domain := os.Args[1]
	server := "8.8.8.8:53"

	conn, err := net.Dial("udp", server)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	query := createDNSQuery(domain)
	_, err = conn.Write(query)
	if err != nil {
		fmt.Printf("Failed to send query: %v\n", err)
		return
	}

	response := make([]byte, 512)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return
	}

	parseDNSResponse(response[:n])
}
