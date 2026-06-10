package main

import (
	"log"
	"net"
)

func main() {
	contactBook := NewContactBookMap()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to listen", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept() // performs the three-way handshake
		if err != nil {
			log.Println("Failed to secure connection", err)
			continue
		}
		go process(conn, contactBook)
	}
}
