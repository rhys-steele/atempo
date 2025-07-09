package repositories

import (
	"testing"
	"time"
)

func TestMemoryCacheRepository_GetSet(t *testing.T) {
	cache := NewMemoryCacheRepository()

	// Test setting and getting a value
	cache.Set("test-key", "test-value", 1*time.Minute)

	value, found := cache.Get("test-key")
	if !found {
		t.Error("Expected to find cached value")
	}

	if value != "test-value" {
		t.Errorf("Expected 'test-value', got '%v'", value)
	}
}

func TestMemoryCacheRepository_Expiration(t *testing.T) {
	cache := NewMemoryCacheRepository()

	// Set value with very short TTL
	cache.Set("test-key", "test-value", 1*time.Millisecond)

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Should not find expired value
	_, found := cache.Get("test-key")
	if found {
		t.Error("Should not find expired value")
	}
}

func TestMemoryCacheRepository_Delete(t *testing.T) {
	cache := NewMemoryCacheRepository()

	// Set and then delete
	cache.Set("test-key", "test-value", 1*time.Minute)
	cache.Delete("test-key")

	// Should not find deleted value
	_, found := cache.Get("test-key")
	if found {
		t.Error("Should not find deleted value")
	}
}

func TestMemoryCacheRepository_Clear(t *testing.T) {
	cache := NewMemoryCacheRepository()

	// Set multiple values
	cache.Set("key1", "value1", 1*time.Minute)
	cache.Set("key2", "value2", 1*time.Minute)

	// Clear cache
	cache.Clear()

	// Should not find any values
	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")

	if found1 || found2 {
		t.Error("Should not find any values after clear")
	}
}

func TestMemoryCacheRepository_Size(t *testing.T) {
	cache := NewMemoryCacheRepository()

	// Check initial size
	if cache.Size() != 0 {
		t.Errorf("Expected size 0, got %d", cache.Size())
	}

	// Add items
	cache.Set("key1", "value1", 1*time.Minute)
	cache.Set("key2", "value2", 1*time.Minute)

	if cache.Size() != 2 {
		t.Errorf("Expected size 2, got %d", cache.Size())
	}
}

func TestMemoryCacheRepository_Keys(t *testing.T) {
	cache := NewMemoryCacheRepository()

	// Add items
	cache.Set("key1", "value1", 1*time.Minute)
	cache.Set("key2", "value2", 1*time.Minute)

	keys := cache.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}

	// Check if keys are present
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	if !keyMap["key1"] || !keyMap["key2"] {
		t.Error("Expected to find both keys")
	}
}

func TestMemoryCacheRepository_Has(t *testing.T) {
	cache := NewMemoryCacheRepository()

	// Should not have key initially
	if cache.Has("test-key") {
		t.Error("Should not have key initially")
	}

	// Add key
	cache.Set("test-key", "test-value", 1*time.Minute)

	// Should have key now
	if !cache.Has("test-key") {
		t.Error("Should have key after setting")
	}
}

func TestMemoryCacheRepository_Invalidate(t *testing.T) {
	cache := NewMemoryCacheRepository()

	// Add items with different expiration times
	cache.Set("key1", "value1", 1*time.Millisecond)  // Will expire quickly
	cache.Set("key2", "value2", 1*time.Hour)         // Will not expire

	// Wait for first item to expire
	time.Sleep(10 * time.Millisecond)

	// Manually invalidate
	cache.Invalidate()

	// Check that expired item is removed and non-expired is kept
	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")

	if found1 {
		t.Error("Expired item should be removed")
	}

	if !found2 {
		t.Error("Non-expired item should be kept")
	}
}

func TestMemoryCacheRepository_Stats(t *testing.T) {
	cache := NewMemoryCacheRepository().(*MemoryCacheRepository)

	// Add items with different expiration times
	cache.Set("key1", "value1", 1*time.Millisecond)  // Will expire quickly
	cache.Set("key2", "value2", 1*time.Hour)         // Will not expire

	// Wait for first item to expire
	time.Sleep(10 * time.Millisecond)

	stats := cache.Stats()

	if stats.TotalItems != 2 {
		t.Errorf("Expected 2 total items, got %d", stats.TotalItems)
	}

	if stats.ExpiredItems != 1 {
		t.Errorf("Expected 1 expired item, got %d", stats.ExpiredItems)
	}

	if stats.ActiveItems != 1 {
		t.Errorf("Expected 1 active item, got %d", stats.ActiveItems)
	}
}

func TestTTLCacheRepository_DefaultTTL(t *testing.T) {
	defaultTTL := 5 * time.Minute
	cache := NewTTLCacheRepository(defaultTTL).(*TTLCacheRepository)

	// Set value without specifying TTL (should use default)
	cache.Set("test-key", "test-value", 0)

	// Verify the value exists
	value, found := cache.Get("test-key")
	if !found {
		t.Error("Expected to find cached value")
	}

	if value != "test-value" {
		t.Errorf("Expected 'test-value', got '%v'", value)
	}
}

func TestTTLCacheRepository_SetWithDefaultTTL(t *testing.T) {
	defaultTTL := 5 * time.Minute
	cache := NewTTLCacheRepository(defaultTTL).(*TTLCacheRepository)

	// Set value with default TTL
	cache.SetWithDefaultTTL("test-key", "test-value")

	// Verify the value exists
	value, found := cache.Get("test-key")
	if !found {
		t.Error("Expected to find cached value")
	}

	if value != "test-value" {
		t.Errorf("Expected 'test-value', got '%v'", value)
	}
}

func TestCacheEntry_IsExpired(t *testing.T) {
	// Create expired entry
	expiredEntry := &CacheEntry{
		Value:     "test",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if !expiredEntry.IsExpired() {
		t.Error("Entry should be expired")
	}

	// Create non-expired entry
	validEntry := &CacheEntry{
		Value:     "test",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if validEntry.IsExpired() {
		t.Error("Entry should not be expired")
	}
}

func TestMemoryCacheRepository_ConcurrentAccess(t *testing.T) {
	cache := NewMemoryCacheRepository()

	// Test concurrent writes and reads
	done := make(chan bool, 2)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set("key", i, 1*time.Minute)
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			cache.Get("key")
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// Test should complete without panicking
}

func BenchmarkMemoryCacheRepository_Set(b *testing.B) {
	cache := NewMemoryCacheRepository()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("bench-key", "bench-value", 1*time.Minute)
	}
}

func BenchmarkMemoryCacheRepository_Get(b *testing.B) {
	cache := NewMemoryCacheRepository()
	cache.Set("bench-key", "bench-value", 1*time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("bench-key")
	}
}