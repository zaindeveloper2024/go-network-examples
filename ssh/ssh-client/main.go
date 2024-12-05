package main

import (
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	config := &ssh.ClientConfig{
		User: "test",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "localhost:2222", config)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %v", err)
	}

	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %v", err)
	}

	if err := session.Wait(); err != nil {
		log.Fatalf("command failed: %v", err)
	}
}
