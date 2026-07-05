package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

var empty = []byte("Command Line cannot be empty\n")
var invalid = []byte("Invalid: Does not have enough arguments\n")
var size = []byte("ERROR: Arguments too big\n")
var emptyVal = []byte("Null\n")
var success = []byte("Success\n")

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

		commandLine := bytes.TrimSpace((*buffHeaderPtr)[:n]) // clears leading /n /r or white spaces
		if len(commandLine) == 0 {
			conn.Write(empty)
			continue
		}

		execute(conn, commandLine, c)
	}
}

func execute(conn net.Conn, commandLine []byte, c *contactBookMap) {

	split := bytes.SplitN(commandLine, []byte(" "), 3) // slice of byte slices [ []byte, []byte ]
	command := split[0]
	args := split[1:]

	for i := range command {
		if command[i] >= 'a' && command[i] <= 'z' {
			command[i] -= 32
		}
	}

	switch string(command) {
	case "SET":
		if len(args) != 2 {
			conn.Write(invalid)
			return
		}
		if len(args[1]) > 1024 || len(args[0]) > 1024 {
			conn.Write(size)
			return
		}
		c.Set(string(args[0]), args[1])
		conn.Write(success)
	case "GET":
		if len(args) != 1 {
			conn.Write(invalid)
			return
		}
		if len(args[0]) > 1024 {
			conn.Write(size)
			return
		}
		output, ok := c.Get(string(args[0]))
		if !ok {
			conn.Write(emptyVal)
			return
		}
		conn.Write(output)
		conn.Write(success)
	case "DELETE":
		if len(args) != 1 {
			conn.Write(invalid)
			return
		}
		if len(args[0]) > 1024 {
			conn.Write(size)
			return
		}
		c.Delete(string(args[0]))
		conn.Write(success)
	case "LIST":
	}
}
