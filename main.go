package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	contactBook := NewContactBookMap()
	wg := sync.WaitGroup{}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	go func() {
		<-shutdown
		cancel()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept() // performs the three-way handshake under the hood
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Printf("Error: %v", err)
				break
			}
			log.Println("Failed to secure connection", err)
			continue
		}
		wg.Add(1)
		go process(conn, contactBook, ctx, &wg)
	}

	wg.Wait()
}
