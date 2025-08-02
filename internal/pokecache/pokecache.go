package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}
type Cache struct {
	Entries  map[string]CacheEntry
	tick     *time.Ticker
	interval time.Duration
	done     chan bool
	mutx     sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	ret := Cache{}
	ret.Entries = make(map[string]CacheEntry, 50)
	ret.interval = interval
	ret.tick = time.NewTicker(interval)
	ret.done = make(chan bool)
	go ret.ReapLoop()
	return &ret
}

func (c *Cache) Add(key string, val []byte) {
	c.mutx.Lock()
	defer c.mutx.Unlock()
	//elem, ok := c.Entries[key]
	c.Entries[key] = CacheEntry{time.Now(), val}
}
func (c *Cache) Get(key string) ([]byte, bool) {
	//ret := []byte{}
	c.mutx.Lock()
	defer c.mutx.Unlock()
	elem, ok := c.Entries[key]
	if ok {
		elem.createdAt = time.Now()
		return elem.val, true
	}
	return nil, false

}

func (c *Cache) reap() {
	c.mutx.Lock()
	defer c.mutx.Unlock()
	t := time.Now()
	var delList []string
	//expTime := t.Sub(c.interval)
	for keyName, entry := range c.Entries {
		expTime := entry.createdAt.Add(c.interval)
		if t.After(expTime) {
			delList = append(delList, keyName)
		}
	}
	for _, keyName := range delList {
		delete(c.Entries, keyName)
	}
}

func (c *Cache) ReapLoop() {
	for {
		select {
		case <-c.done:
			return
		case <-c.tick.C:
			c.reap()
		}
	}
}
