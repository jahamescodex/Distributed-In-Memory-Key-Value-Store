//go:build ignore

package main

import (
	"strconv"
	"sync"
	"sync/atomic"
)

type contactBookMap struct {
	contactBook map[string]string
	lock        sync.RWMutex
	ops         atomic.Uint64
}

func (c *contactBookMap) Set(key string, val string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	UID := strconv.FormatUint(c.ops.Add(1), 10)
	c.contactBook[key] = (UID + " - " + val)
}

func (c *contactBookMap) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.contactBook, key)
}

func (c *contactBookMap) Get(key string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	val, ok := c.contactBook[key]
	if !ok {
		return "Key-value pair does not exit"
	}
	return val
}

func NewContactBookMap() *contactBookMap {
	return &contactBookMap{
		contactBook: make(map[string]string),
	}
}
