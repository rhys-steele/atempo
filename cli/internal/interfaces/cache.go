package interfaces

import "time"

// CacheRepository defines in-memory caching operations.
type CacheRepository interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
	Has(key string) bool
	Invalidate()
}
