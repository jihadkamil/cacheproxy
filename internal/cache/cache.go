package cache

import (
	"sync"
	"time"
)

// Cache defines the interface for a generic caching system.
type Cache interface {
	// Define cache interface methods here
	Set(key string, val []byte, ttl time.Duration) error
	Get(key string) ([]byte, bool)
	Delete(key string) error
}

type item struct {
	value      []byte
	expireDate time.Time
}

// inMemoryCache is a thread-safe implementation of the Cache interface.
type inMemoryCache struct {
	// items stores the cache data. Note: Go maps are NOT thread-safe by default.
	items map[string]item
	// mu is a Read/Write Mutex. It protects the 'items' map from concurrent access issues.
	// RWMutex is optimized for scenarios with many readers and few writers.
	mu sync.RWMutex
}

// NewInMemoryCache initializes the cache storage.
func NewInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		items: make(map[string]item),
	}
}

// Get retrieves an item from the cache if it exists and hasn't expired.
func (c *inMemoryCache) Get(key string) ([]byte, bool) {
	// RLock (Read Lock) allows multiple goroutines to read the map simultaneously.
	// We use this because we are only looking up a value, not changing the map.
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock() // We must release the read lock as soon as we are done reading.

	if !ok {
		return nil, false
	}

	// Check if the item has expired.
	if time.Now().After(it.expireDate) {
		// The item is expired, so we need to delete it.
		// Deleting modifies the map, so we need a full Lock (Write Lock).
		// Lock() ensures EXCLUSIVE access. No other goroutine can Read OR Write
		// while we hold this lock.
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock() // Release the write lock immediately after modification.
		return nil, false
	}
	return it.value, true
}

// Set adds or updates an item in the cache.
func (c *inMemoryCache) Set(key string, val []byte, ttl time.Duration) error {
	// We are writing to the map, so we use Lock (Write Lock).
	// This blocks all other readers (RLock) and writers (Lock) until we Unlock.
	c.mu.Lock()
	c.items[key] = item{
		value:      val,
		expireDate: time.Now().Add(ttl),
	}
	c.mu.Unlock() // Release the lock so others can access the map again.
	return nil
}

func (c *inMemoryCache) Delete(key string) error {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
	return nil
}
