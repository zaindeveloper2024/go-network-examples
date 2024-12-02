package main

import (
	"fmt"
	"net"
	"sync"
)

type Rule struct {
	SourceIP      net.IP
	DestinationIP net.IP
	Protocol      string
	Port          int
	Action        string
}

type Firewall struct {
	rules []Rule
	mu    sync.RWMutex
}

func NewFirewall() *Firewall {
	return &Firewall{
		rules: make([]Rule, 0),
	}
}

func (fw *Firewall) AddRule(rule Rule) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.rules = append(fw.rules, rule)
}

func (fw *Firewall) CheckPacket(sourceIP net.IP, destIP net.IP, protocol string, port int) bool {
	fw.mu.RLock()
	defer fw.mu.RUnlock()

	for _, rule := range fw.rules {
		if matchRule(rule, sourceIP, destIP, protocol, port) {
			return rule.Action == "allow"
		}
	}

	return false
}

func matchRule(rule Rule, sourceIP net.IP, destIP net.IP, protocol string, port int) bool {
	if rule.SourceIP != nil && !rule.SourceIP.Equal(sourceIP) {
		return false
	}

	if rule.DestinationIP != nil && !rule.DestinationIP.Equal(destIP) {
		return false
	}

	if rule.Protocol != "" && rule.Protocol != protocol {
		return false
	}

	if rule.Port != 0 && rule.Port != port {
		return false
	}

	return true
}

func main() {
	fw := NewFirewall()

	fw.AddRule(Rule{
		SourceIP:      net.ParseIP("192.168.1.100"),
		DestinationIP: net.ParseIP("10.0.0.1"),
		Protocol:      "tcp",
		Port:          80,
		Action:        "allow",
	})

	fw.AddRule(Rule{
		SourceIP: net.ParseIP("192.168.1.0"),
		Protocol: "udp",
		Port:     53,
		Action:   "deny",
	})

	testCases := []struct {
		sourceIP string
		destIP   string
		protocol string
		port     int
	}{
		{"192.168.1.100", "10.0.0.1", "tcp", 80},
		{"192.168.1.100", "10.0.0.1", "tcp", 443},
		{"192.168.1.101", "10.0.0.1", "udp", 53},
	}

	for _, tc := range testCases {
		allowed := fw.CheckPacket(
			net.ParseIP(tc.sourceIP),
			net.ParseIP(tc.destIP),
			tc.protocol,
			tc.port,
		)
		fmt.Printf("Packet from %s to %s (%s:%d): %v\n",
			tc.sourceIP, tc.destIP, tc.protocol, tc.port,
			allowed,
		)
	}
}
