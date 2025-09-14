package service

import (
	"sync"
	"time"
)

type CacheEntry struct {
	OriginalURL string
	ExpiresAt   *time.Time
	MaxVisits   int64
}

type InMemoryCache struct {
	mu   sync.RWMutex
	data map[string]CacheEntry
}

func newInMemoryCache() *InMemoryCache {
	return &InMemoryCache{data: make(map[string]CacheEntry)}
}

func (c *InMemoryCache) Set(key string, val CacheEntry) {
	c.mu.Lock()
	c.data[key] = val
	c.mu.Unlock()
}

func (c *InMemoryCache) Get(key string) (CacheEntry, bool) {
	c.mu.RLock()
	v, ok := c.data[key]
	c.mu.RUnlock()
	return v, ok
}

func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}
