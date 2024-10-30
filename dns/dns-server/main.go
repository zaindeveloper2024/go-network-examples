package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type DNSRecord struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32
	IP    net.IP
}

type RecordData struct {
	Name string
	IP   string
}

type DNSServer struct {
	records map[string]DNSRecord
}

func NewDNSServer() *DNSServer {
	return &DNSServer{
		records: make(map[string]DNSRecord),
	}
}

func (s *DNSServer) AddRecord(name string, ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address: %s", ip)
	}

	s.records[name] = DNSRecord{
		Name:  name,
		Type:  dns.TypeA,
		Class: dns.ClassINET,
		TTL:   300,
		IP:    parsedIP,
	}

	return nil
}

func (s *DNSServer) AddRecords(records []RecordData) error {
	for _, record := range records {
		err := s.AddRecord(record.Name, record.IP)
		if err != nil {
			return fmt.Errorf("Failed to add record: %s\n", err.Error())
		}
	}
	return nil
}

func (s *DNSServer) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	log.Printf("received query for %v\n", r.Question)

	switch r.Opcode {
	case dns.OpcodeQuery:
		for _, q := range r.Question {
			log.Printf("Query for %s\n", q.Name)
			name := strings.ToLower(q.Name)
			if record, exists := s.records[name]; exists {
				if q.Qtype == dns.TypeA {
					rr := &dns.A{
						Hdr: dns.RR_Header{
							Name:   q.Name,
							Rrtype: dns.TypeA,
							Class:  dns.ClassINET,
							Ttl:    record.TTL,
						},
						A: record.IP,
					}
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}

	w.WriteMsg(m)
}
func main() {
	dnsServer := NewDNSServer()

	records := []RecordData{
		{Name: "example.com.", IP: "93.184.216.34"},
		{Name: "test.com.", IP: "192.168.1.1"},
		{Name: "demo.com.", IP: "10.0.0.1"},
	}

	err := dnsServer.AddRecords(records)
	if err != nil {
		log.Fatalf("Failed to add records: %s\n", err.Error())
	}

	dns.HandleFunc(".", dnsServer.handleDNSRequest)

	port := 53
	log.Printf("Starting DNS server on port %d\n", port)

	srv := &dns.Server{
		Addr: fmt.Sprintf(":%d", port),
		Net:  "udp",
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
