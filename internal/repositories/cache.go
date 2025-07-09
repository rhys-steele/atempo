package repositories

import (
	"sync"
	"time"

	"atempo/internal/interfaces"
)

// CacheEntry represents a cached item with expiration
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired checks if the cache entry has expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// MemoryCacheRepository implements CacheRepository with in-memory storage
type MemoryCacheRepository struct {
	mu    sync.RWMutex
	items map[string]*CacheEntry
}

// NewMemoryCacheRepository creates a new in-memory cache repository
func NewMemoryCacheRepository() interfaces.CacheRepository {
	cache := &MemoryCacheRepository{
		items: make(map[string]*CacheEntry),
	}
	
	// Start background cleanup goroutine
	go cache.cleanupExpired()
	
	return cache
}

// Get retrieves a value from the cache
func (c *MemoryCacheRepository) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, exists := c.items[key]
	if !exists {
		return nil, false
	}
	
	if entry.IsExpired() {
		// Entry has expired, remove it
		delete(c.items, key)
		return nil, false
	}
	
	return entry.Value, true
}

// Set stores a value in the cache with a TTL
func (c *MemoryCacheRepository) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items[key] = &CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *MemoryCacheRepository) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *MemoryCacheRepository) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items = make(map[string]*CacheEntry)
}

// Size returns the number of items in the cache
func (c *MemoryCacheRepository) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return len(c.items)
}

// Keys returns all cache keys
func (c *MemoryCacheRepository) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	
	return keys
}

// Has checks if a key exists in the cache (without checking expiration)
func (c *MemoryCacheRepository) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	_, exists := c.items[key]
	return exists
}

// Invalidate removes expired entries from the cache
func (c *MemoryCacheRepository) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	now := time.Now()
	for key, entry := range c.items {
		if now.After(entry.ExpiresAt) {
			delete(c.items, key)
		}
	}
}

// cleanupExpired runs in the background to periodically clean up expired entries
func (c *MemoryCacheRepository) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute) // Cleanup every 5 minutes
	defer ticker.Stop()
	
	for range ticker.C {
		c.Invalidate()
	}
}

// Stats returns cache statistics
func (c *MemoryCacheRepository) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	now := time.Now()
	expired := 0
	
	for _, entry := range c.items {
		if now.After(entry.ExpiresAt) {
			expired++
		}
	}
	
	return CacheStats{
		TotalItems:   len(c.items),
		ExpiredItems: expired,
		ActiveItems:  len(c.items) - expired,
	}
}

// CacheStats represents cache statistics
type CacheStats struct {
	TotalItems   int
	ExpiredItems int
	ActiveItems  int
}

// TTLCacheRepository implements CacheRepository with configurable TTL defaults
type TTLCacheRepository struct {
	cache      interfaces.CacheRepository
	defaultTTL time.Duration
}

// NewTTLCacheRepository creates a new cache repository with default TTL
func NewTTLCacheRepository(defaultTTL time.Duration) interfaces.CacheRepository {
	return &TTLCacheRepository{
		cache:      NewMemoryCacheRepository(),
		defaultTTL: defaultTTL,
	}
}

// Get retrieves a value from the cache
func (c *TTLCacheRepository) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

// Set stores a value in the cache with default TTL
func (c *TTLCacheRepository) Set(key string, value interface{}, ttl time.Duration) {
	if ttl == 0 {
		ttl = c.defaultTTL
	}
	c.cache.Set(key, value, ttl)
}

// SetWithDefaultTTL stores a value in the cache with the default TTL
func (c *TTLCacheRepository) SetWithDefaultTTL(key string, value interface{}) {
	c.cache.Set(key, value, c.defaultTTL)
}

// Delete removes a value from the cache
func (c *TTLCacheRepository) Delete(key string) {
	c.cache.Delete(key)
}

// Clear removes all items from the cache
func (c *TTLCacheRepository) Clear() {
	c.cache.Clear()
}

// Size returns the number of items in the cache
func (c *TTLCacheRepository) Size() int {
	return c.cache.Size()
}

// Keys returns all cache keys
func (c *TTLCacheRepository) Keys() []string {
	return c.cache.Keys()
}

// Has checks if a key exists in the cache
func (c *TTLCacheRepository) Has(key string) bool {
	return c.cache.Has(key)
}

// Invalidate removes expired entries from the cache
func (c *TTLCacheRepository) Invalidate() {
	c.cache.Invalidate()
}