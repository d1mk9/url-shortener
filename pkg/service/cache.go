package service

import (
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

type InMemoryCache struct {
	cache *lru.Cache[string, CacheEntry]
}

type CacheEntry struct {
	OriginalURL string
	ExpiresAt   *time.Time
	MaxVisits   int64
}

func newInMemoryCache() *InMemoryCache {
	c, _ := lru.New[string, CacheEntry](10000)
	return &InMemoryCache{cache: c}
}

func (c *InMemoryCache) Set(key string, val CacheEntry) {
	c.cache.Add(key, val)
}

func (c *InMemoryCache) Get(key string) (CacheEntry, bool) {
	return c.cache.Get(key)
}

func (c *InMemoryCache) Delete(key string) {
	c.cache.Remove(key)
}
