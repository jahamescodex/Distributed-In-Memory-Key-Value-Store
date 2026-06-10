package main

import (
	"sync"
)

type contactBookMap struct {
	contactBook map[string]Record
	lock        sync.RWMutex
	countOPS    uint
	HWM         uint
}

type Record struct {
	Data []byte
	ID   uint
}

func (c *contactBookMap) Set(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.countOPS++

	cloneVal := make([]byte, len(val))
	copy(cloneVal, val)
	c.contactBook[key] = Record{
		Data: cloneVal,
		ID:   c.countOPS,
	}

	if uint(len(c.contactBook)) > c.HWM {
		c.HWM = uint(len(c.contactBook))
	}

}

// Maps do NOT shrink
func (c *contactBookMap) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	currentHWM := c.HWM
	delete(c.contactBook, key)
	if len(c.contactBook) < int(currentHWM/4) && len(c.contactBook) > 8 {
		copyMap := make(map[string]Record, len(c.contactBook))
		for k, v := range c.contactBook {
			copyMap[k] = v
		}
		c.contactBook = copyMap
		c.HWM = uint(len(c.contactBook))
	}
}

func (c *contactBookMap) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	val, ok := c.contactBook[key]
	if !ok {
		return nil, false
	} else {
		return val.Data, true
	}

}

func NewContactBookMap() *contactBookMap {
	return &contactBookMap{
		contactBook: make(map[string]Record),
		HWM:         0,
	}
}
