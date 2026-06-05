package main

import (
	"log"
	"net"
	"strings"
)

type opCode int

const (
	OpInvalidCode opCode = iota
	OpSet
	OpDelete
	OpGet
	OpList
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

			if strings.ToUpper(commandLine[0]) == "EXIT" {
				connection.Write([]byte("Ending Connection\n"))
				log.Printf("Connection: %v closed", connection)
				connection.Close()
			}

			parser := &cmd{
				opCode: commandLine[0],
				args:   commandLine[1:],
			}

			processRequest(parser, contactBook, connection)

			writePos = 0
			continue
		}
		writePos++

		if writePos >= len(buff) {
			log.Printf("Oversized payload from %s - closing wire\n", connection.RemoteAddr())
			connection.Write([]byte("ERR Command too long\n"))
			return
		}
	}
}

func processRequest(input *cmd, contactBook *contactBookMap, connection net.Conn) {
	var key, val string
	op := strings.ToUpper(input.opCode)

	code := validator(op, input.args)

	switch code {
	case OpInvalidCode:
		connection.Write([]byte("Command issue\n"))
		return
	case OpSet:
		key = input.args[0]
		val = strings.Join(input.args[1:], " ")
		contactBook.Set(key, val)
		connection.Write([]byte("Set Complete\n"))
	case OpDelete:
		key = input.args[0]
		contactBook.Delete(key)
		connection.Write([]byte("Delete Complete\n"))
	case OpGet:
		key = input.args[0]
		msg, ok := contactBook.Get(key)
		if !ok {
			connection.Write([]byte(msg + "\n"))
		} else {
			connection.Write([]byte(msg + "\n"))
		}
	case OpList:
		contactBook.lock.RLock()
		defer contactBook.lock.RUnlock()
		if len(contactBook.contactBook) == 0 {
			connection.Write([]byte("Empty List\n"))
		} else {
			for _, val := range contactBook.contactBook {
				connection.Write([]byte(val + "\n"))
			}
		}
	default:
		connection.Write([]byte("Error\n"))
		return
	}
}

func validator(op string, args []string) opCode {
	if op == "SET" && len(args) < 2 {
		return OpInvalidCode
	} else if (op == "DELETE" || op == "GET") && len(args) < 1 {
		return OpInvalidCode
	}
	if op == "SET" {
		return OpSet
	}
	if op == "DELETE" {
		return OpDelete
	}
	if op == "GET" {
		return OpGet
	}
	if op == "LIST" {
		return OpList
	}
	return OpInvalidCode
}
