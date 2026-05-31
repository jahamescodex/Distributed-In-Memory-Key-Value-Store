package main

import (
	"log"
	"net"
)

func main() {
	contactBook := NewContactBookMap()
	jobsPipe := make(chan net.Conn, 100)

	for w := 1; w <= 10; w++ { // starting 10 workers
		go process(jobsPipe, contactBook)
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}
		jobsPipe <- conn
	}
}
