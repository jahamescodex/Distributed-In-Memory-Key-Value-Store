package main

import (
	"net"
	"strconv"
	"sync"
)

type contactBookMap struct {
	contactBook map[string]Record
	lock        sync.RWMutex
	HWM         int
	counterOPS  int
}

type Record struct {
	data []byte
	ID   int
}

func (c *contactBookMap) Set(key []byte, val []byte) { //recieve a struct which has a pointer
	c.lock.Lock()
	defer c.lock.Unlock()

	c.counterOPS++ // 1 - Key : Value

	copyVal := make([]byte, len(val)) // this is the 24byte struct, allocated on the function execution stack frame; this contains a pointer that points to the backing array that is on the heap
	copy(copyVal, val)                //dst, src
	c.contactBook[string(key)] = Record{
		data: copyVal,
		ID:   c.counterOPS,
	} // increased the length of the map

	if len(c.contactBook) > c.HWM {
		c.HWM = len(c.contactBook)
	}
}

func (c *contactBookMap) Get(key []byte) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	rec, ok := c.contactBook[string(key)]
	if !ok {
		return nil, false
	}
	copyVal := make([]byte, len(rec.data))
	copy(copyVal, rec.data)
	return copyVal, true
}

func (c *contactBookMap) Delete(key []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	currentHWM := c.HWM
	delete(c.contactBook, string(key))
	if len(c.contactBook) < (currentHWM/4) && len(c.contactBook) > 8 {
		copyMap := make(map[string]Record, len(c.contactBook))
		for k, v := range c.contactBook {
			copyMap[k] = v
		}
		c.contactBook = copyMap
		c.HWM = len(c.contactBook)
	}
}

func (c *contactBookMap) List(conn net.Conn, buffer *[]byte) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	window := (*buffer)
	window = window[:cap(window)]
	window = window[:0]
	for k, v := range c.contactBook { // ID_-_k_:_v\n
		digitLength := digitLength(v.ID)
		currentLength := len(window) + digitLength + 7 + len(v.data) + len(k)
		if currentLength > 1024 {
			conn.Write(window)
			window = window[:0]
		}
		window = strconv.AppendInt(window, int64(v.ID), 10)
		window = append(window, " - "...)
		window = append(window, k...)
		window = append(window, " : "...)
		window = append(window, v.data...)
		window = append(window, "\n"...)
	}
	if len(window) != 0 {
		conn.Write(window)
	}
}

func digitLength(input int) int {
	if input == 0 {
		return 1
	}

	length := 0

	for input != 0 {
		length++
		input = input / 10
	}
	return length
}

func NewContactBookMap() *contactBookMap {
	return &contactBookMap{
		contactBook: make(map[string]Record),
		HWM:         0,
		counterOPS:  0,
	}
}
