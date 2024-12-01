package main

import "sync"

type Rule struct {
	SourceIP      string
	DestinationIP string
	Protocol      string
	Port          int
	Action        string
}

type Firewall struct {
	rules []Rule
	mu    sync.Mutex
}

func NewFireWall() *Firewall {
	return &Firewall{
		rules: make([]Rule, 0),
	}
}

func (fw *Firewall) AddRule(rule Rule) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.rules = append(fw.rules, rule)
}

func main() {
	fw := NewFireWall()

	fw.AddRule(Rule{
		SourceIP:      "192.168.1.100",
		DestinationIP: "*",
		Protocol:      "tcp",
		Port:          80,
		Action:        "allow",
	})

	fw.AddRule(Rule{
		SourceIP:      "*",
		DestinationIP: "*",
		Protocol:      "tcp",
		Port:          443,
		Action:        "allow",
	})

	fw.AddRule(Rule{
		SourceIP:      "10.0.0.100",
		DestinationIP: "*",
		Protocol:      "*",
		Port:          0,
		Action:        "deny",
	})

	select {}
}
