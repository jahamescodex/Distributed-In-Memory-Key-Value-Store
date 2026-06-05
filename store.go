package main

import (
	"sync"
)

type contactBookMap struct {
	contactBook map[string]Record
	lock        sync.RWMutex
	ID          uint
}

type Record struct {
	UID uint
	val []byte
}

func (c *contactBookMap) Set(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.ID++
	reserve := make([]byte, len(val))
	copy(reserve, val)
	c.contactBook[key] = Record{
		UID: c.ID,
		val: reserve,
	}
}

func (c *contactBookMap) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.contactBook, key)
}

func (c *contactBookMap) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	val, ok := c.contactBook[key]
	if !ok {
		return nil, false
	}
	return val.val, true
}

func NewContactBookMap() *contactBookMap {
	return &contactBookMap{
		contactBook: make(map[string]Record),
	}
}
