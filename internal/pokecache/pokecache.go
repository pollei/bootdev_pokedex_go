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
	Entries map[string]CacheEntry
	mutx    sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	return Cache{}
}

func (c Cache) Add(key string, val []byte) {
	c.mutx.Lock()
	defer c.mutx.Unlock()

}
func (c Cache) Get(key string) ([]byte, bool) {
	ret := []byte{}
	c.mutx.Lock()
	defer c.mutx.Unlock()
	return ret, false

}

func (c Cache) ReapLoop() {

}
