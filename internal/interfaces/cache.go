package interfaces

import "time"

// CacheRepository defines the interface for caching operations
type CacheRepository interface {
	// Get retrieves a value from the cache
	Get(key string) (interface{}, bool)
	
	// Set stores a value in the cache with a TTL
	Set(key string, value interface{}, ttl time.Duration)
	
	// Delete removes a value from the cache
	Delete(key string)
	
	// Clear removes all items from the cache
	Clear()
	
	// Size returns the number of items in the cache
	Size() int
	
	// Keys returns all cache keys
	Keys() []string
	
	// Has checks if a key exists in the cache
	Has(key string) bool
	
	// Invalidate removes expired entries from the cache
	Invalidate()
}