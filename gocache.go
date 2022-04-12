package gocache

import (
	"sync"
	"time"
)

var never = time.Now().Add(time.Hour * 24 * 365 * 100)

// Cache stores cached values.
type Cache[K comparable, V any] struct {
	store  map[K]*entry[V]
	mu     sync.Mutex
	config *Config
}

// New creates a new cache. A Config can be passed in.
func New[K comparable, V any](config ...*Config) *Cache[K, V] {
	c := &Cache[K, V]{}
	if len(config) > 0 {
		c.config = config[0]
	} else {
		c.config = NewConfig()
	}
	c.store = make(map[K]*entry[V], c.config.capacityHint)
	go c.watch()
	return c
}

func (c *Cache[K, V]) watch() {
	if c.config.evictionPolicy.flags&evictionPolicyActivePruningMask == 0 {
		return
	}
	for {
		time.Sleep(time.Minute)
		c.mu.Lock()
		c.expireUnsafe()
		c.mu.Unlock()
	}
}

func (c *Cache[K, V]) expireUnsafe() {
	ref := time.Now()
	for k, e := range c.store {
		if e.isExpired(ref) {
			delete(c.store, k)
		}
	}
}

func (c *Cache[K, V]) evict() {
	if c.Size() < c.config.maxCapacity {
		return
	}
	var someExpire bool
	var evict K
	ts := never
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, e := range c.store {
		if someExpire && !e.exp {
			continue
		}
		if !someExpire && e.exp {
			someExpire = true
			ts = e.ts
			evict = k
			continue
		}
		if e.ts.Before(ts) {
			ts = e.ts
			evict = k
		}
	}
	delete(c.store, evict)
}

// Set sets a value in the cache with the given key. A ttl can be
// passed in as a third argument to override the default ttl.
func (c *Cache[K, V]) Set(key K, val V, ttl ...time.Duration) {
	c.evict()
	e := newEntry(val)
	if len(ttl) > 0 {
		e.ts = e.ts.Add(ttl[0])
		e.exp = true
	} else if c.config.evictionPolicy.flags&evictionPolicySetMask > 0 {
		e.ts = e.ts.Add(c.config.evictionPolicy.ttl)
		e.exp = true
	}
	c.mu.Lock()
	c.store[key] = e
	c.mu.Unlock()
}

// Get gets a value from the cache with the given key. If a value
// with the given key is not present in the cache, the zero value
// is returned.
func (c *Cache[K, V]) Get(key K) V {
	c.mu.Lock()
	defer c.mu.Unlock()
	e := c.store[key]
	if e == nil || e.isExpired(time.Now()) {
		c.expireUnsafe()
		var v V
		return v
	}
	if c.config.evictionPolicy.flags&evictionPolicyGetMask > 0 {
		e.ts = time.Now().Add(c.config.evictionPolicy.ttl)
	}
	return e.val
}

// Has checks if there is a value in the cache with the given key.
func (c *Cache[K, V]) Has(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.store[key]
	if !ok {
		return false
	}
	if e.isExpired(time.Now()) {
		c.expireUnsafe()
		return false
	}
	if c.config.evictionPolicy.flags&evictionPolicyHasMask > 0 {
		e.ts = time.Now().Add(c.config.evictionPolicy.ttl)
	}
	return true
}

// Delete deletes a value from the cache with the given key.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// Size returns the number of values in the cache that have not
// expired.
func (c *Cache[K, V]) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.expireUnsafe()
	return len(c.store)
}

// Clear removes all values from the cache.
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[K]*entry[V], c.config.capacityHint)
}

// Keys returns all the keys in the cache for values that have not
// expired.
func (c *Cache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.expireUnsafe()
	keys := make([]K, 0, len(c.store))
	for k := range c.store {
		keys = append(keys, k)
	}
	return keys
}
