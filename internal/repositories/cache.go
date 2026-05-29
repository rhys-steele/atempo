package repositories

import (
	"sync"
	"time"

	"atempo/internal/interfaces"
)

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

func (e *cacheEntry) expired() bool {
	return time.Now().After(e.expiresAt)
}

// MemoryCacheRepository is a thread-safe in-memory cache with TTL expiry.
type MemoryCacheRepository struct {
	mu    sync.Mutex
	items map[string]*cacheEntry
}

// NewMemoryCacheRepository creates a new cache and starts a background cleanup ticker.
func NewMemoryCacheRepository() interfaces.CacheRepository {
	c := &MemoryCacheRepository{
		items: make(map[string]*cacheEntry),
	}
	go c.runCleanup()
	return c
}

func (c *MemoryCacheRepository) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if entry.expired() {
		delete(c.items, key)
		return nil, false
	}
	return entry.value, true
}

func (c *MemoryCacheRepository) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *MemoryCacheRepository) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *MemoryCacheRepository) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*cacheEntry)
}

func (c *MemoryCacheRepository) Has(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.items[key]
	return ok && !entry.expired()
}

func (c *MemoryCacheRepository) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, e := range c.items {
		if now.After(e.expiresAt) {
			delete(c.items, k)
		}
	}
}

func (c *MemoryCacheRepository) runCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.Invalidate()
	}
}
