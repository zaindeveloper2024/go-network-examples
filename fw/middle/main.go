package main

import (
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
	mu    sync.Mutex
}

func NewFirewall() *Firewall {
	return &Firewall{
		rules: make([]Rule, 0),
	}
}

func (f *Firewall) AddRule(rule Rule) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.rules = append(f.rules, rule)
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
}
