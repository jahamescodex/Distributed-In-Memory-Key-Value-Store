package main

import (
	"log"
	"net"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		buff := make([]byte, 1024)
		return &buff
	},
}

func process(conn net.Conn, contactBook *contactBookMap) {
	handleClient(conn, contactBook, &bufferPool) // do NOT spawn another goroutine, we are already running async
}

func handleClient(conn net.Conn, contactBook *contactBookMap, bufferPool *sync.Pool) {
	defer conn.Close()
	log.Println("New Client Just connected:", conn.RemoteAddr())
	bufPtr := bufferPool.Get().(*[]byte)

	defer func() {
		bufferPool.Put(bufPtr)
		log.Println("Buffer Slice back into the pool")
	}()

}
