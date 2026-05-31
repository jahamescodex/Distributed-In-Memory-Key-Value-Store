package main

import (
	"log"
	"net"
	"strings"
)

type cmd struct {
	opCode string
	args   []string
}

func process(jobs <-chan net.Conn, contactBook *contactBookMap) {
	for conn := range jobs {
		handleClient(conn, contactBook)
	}
}

func handleClient(connection net.Conn, contactBook *contactBookMap) {
	defer connection.Close()
	var buff [1024]byte
	writePos := 0

	for {
		n, err := connection.Read(buff[writePos : writePos+1])
		if err != nil {
			return
		}
		if n == 0 {
			continue
		}
		currentByte := buff[writePos]

		if currentByte == '\n' {
			line := buff[:writePos]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}

			commandLine := strings.Fields(string(line))

			if len(commandLine) == 0 {
				continue
			}

			parser := &cmd{
				opCode: commandLine[0],
				args:   commandLine[1:],
			}

			processRequest(parser, contactBook)

			writePos = 0
			continue
		}
		writePos++

		if writePos >= len(buff) {
			log.Printf("Oversized payload from %s - closing wire", connection.RemoteAddr())
			connection.Write([]byte("ERR Command too long\n"))
			return
		}
	}
}

func processRequest(input *cmd, contactBook *contactBookMap) {

}
