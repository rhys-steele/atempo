package repositories

import (
	"sync"
	"testing"
	"time"
)

func TestCacheSetAndGet(t *testing.T) {
	c := NewMemoryCacheRepository()

	c.Set("key", "value", time.Minute)

	v, ok := c.Get("key")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if v.(string) != "value" {
		t.Fatalf("got %v, want %q", v, "value")
	}
}

func TestCacheMiss(t *testing.T) {
	c := NewMemoryCacheRepository()
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected miss for unknown key")
	}
}

func TestCacheExpiry(t *testing.T) {
	c := NewMemoryCacheRepository()
	c.Set("key", "value", time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected key to be expired")
	}
}

// TestCacheGetDeletesExpired explicitly covers the data-race fix:
// Get must hold a write lock when deleting expired entries, not a read lock.
// Run with -race to catch any regression.
func TestCacheGetDeletesExpiredUnderWriteLock(t *testing.T) {
	c := NewMemoryCacheRepository()
	c.Set("k", 1, time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Get("k") // concurrent Get on an expired key
		}()
	}
	wg.Wait()
}

func TestCacheDelete(t *testing.T) {
	c := NewMemoryCacheRepository()
	c.Set("key", "value", time.Minute)
	c.Delete("key")

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestCacheHas(t *testing.T) {
	c := NewMemoryCacheRepository()
	c.Set("key", "value", time.Minute)

	if !c.Has("key") {
		t.Fatal("expected Has to return true")
	}
	c.Delete("key")
	if c.Has("key") {
		t.Fatal("expected Has to return false after delete")
	}
}

func TestCacheClear(t *testing.T) {
	c := NewMemoryCacheRepository()
	c.Set("a", 1, time.Minute)
	c.Set("b", 2, time.Minute)
	c.Clear()

	if c.Has("a") || c.Has("b") {
		t.Fatal("expected cache to be empty after Clear")
	}
}

func TestCacheInvalidate(t *testing.T) {
	c := NewMemoryCacheRepository()
	c.Set("live", 1, time.Minute)
	c.Set("dead", 2, time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	c.Invalidate()

	if !c.Has("live") {
		t.Error("expected live key to survive Invalidate")
	}
	if c.Has("dead") {
		t.Error("expected dead key to be removed by Invalidate")
	}
}
