package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

var empty = []byte("Command Line cannot be empty")

var bufferPool = sync.Pool{
	New: func() any {
		bufHeader := make([]byte, 1024)
		return &bufHeader
	},
}

func process(conn net.Conn, c *contactBookMap) {
	handleClient(conn, c, &bufferPool)
}

func handleClient(conn net.Conn, c *contactBookMap, bufferPool *sync.Pool) {
	defer conn.Close()
	log.Printf("Client: %s just connected\n", conn.RemoteAddr())
	buffHeaderPtr := bufferPool.Get().(*[]byte) // pointing to the 24-byte struct byte slice-header

	defer func() {
		log.Printf("Client: %s just disconnected, buffer put back into pool", conn.RemoteAddr())
		clear(*buffHeaderPtr)         // dereferences to gain access to the underlying back array that points to the actual information
		bufferPool.Put(buffHeaderPtr) // returning the pointer of the 24-byte struct back into the pool
	}()

	for {
		*buffHeaderPtr = (*buffHeaderPtr)[:cap(*buffHeaderPtr)]

		n, err := conn.Read(*buffHeaderPtr)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Println("Error:", err)
			return
		}

		commandLine := bytes.TrimSpace((*buffHeaderPtr)[:n])
		if len(commandLine) == 0 {
			conn.Write(empty)
			continue
		}

		token(commandLine, c)
	}
}

func token(commandLine []byte, c *contactBookMap) {

}
