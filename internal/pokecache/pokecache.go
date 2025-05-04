package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if entry, ok := c.entries[key]; ok {
		return entry.val, true
	}
	return nil, false
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{entries: map[string]cacheEntry{}, mu: &sync.Mutex{}}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			for k, v := range c.entries {
				if time.Since(v.createdAt) > interval {
					c.mu.Lock()
					delete(c.entries, k)
					c.mu.Unlock()
				}
			}
		}
	}
}
